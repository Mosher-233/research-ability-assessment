package extractor

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testdataWordsDir = "../../testdata/words"


func TestDOCXExtractor_SupportedFormats(t *testing.T) {
	e := NewDOCXExtractor()
	formats := e.SupportedFormats()
	require.Len(t, formats, 1)
	assert.Equal(t, ".docx", formats[0])
}

func TestDOCXExtractor_Extract_NonExistentFile(t *testing.T) {
	e := NewDOCXExtractor()
	content, err := e.Extract("/nonexistent/path/file.docx")
	assert.Error(t, err)
	assert.Empty(t, content)
}

func TestDOCXExtractor_Extract_QualityA(t *testing.T) {
	skipIfNoDir(t, testdataWordsDir)
	e := NewDOCXExtractor()
	path := filepath.Join(testdataWordsDir, "A_优秀_S001_张明_语音情感识别_R001.docx")
	content, err := e.Extract(path)
	require.NoError(t, err)
	require.NotEmpty(t, content)
	assert.True(t, len([]rune(content)) > 100, "A-quality docx should have substantial text, got %d chars", len([]rune(content)))
	t.Logf("A-quality docx: %d chars, first 100: %s", len([]rune(content)), truncateRunes(content, 100))
}

func TestDOCXExtractor_Extract_QualityD(t *testing.T) {
	skipIfNoDir(t, testdataWordsDir)
	e := NewDOCXExtractor()
	path := filepath.Join(testdataWordsDir, "D_不合格_S013_马超_联邦学习推荐_R013.docx")
	content, err := e.Extract(path)
	require.NoError(t, err)
	require.NotEmpty(t, content)
	t.Logf("D-quality docx: %d chars", len([]rune(content)))
}

func TestDOCXExtractor_Extract_QualityE_OffTopic(t *testing.T) {
	skipIfNoDir(t, testdataWordsDir)
	e := NewDOCXExtractor()
	path := filepath.Join(testdataWordsDir, "E_答非所问_S016_徐婷_语音情感识别_R016.docx")
	content, err := e.Extract(path)
	require.NoError(t, err)
	require.NotEmpty(t, content)
	assert.True(t, len([]rune(content)) > 0, "off-topic docx should still extract text")
}

func TestDOCXExtractor_Extract_QualityF_Empty(t *testing.T) {
	skipIfNoDir(t, testdataWordsDir)
	e := NewDOCXExtractor()
	path := filepath.Join(testdataWordsDir, "F_内容空洞_S018_曹雪_GNN影响力分析_R018.docx")
	content, err := e.Extract(path)
	require.NoError(t, err)
	require.NotEmpty(t, content)
	t.Logf("F-quality (empty content) docx: %d chars", len([]rune(content)))
}

func TestDOCXExtractor_Extract_QualityG_Malformed(t *testing.T) {
	skipIfNoDir(t, testdataWordsDir)
	e := NewDOCXExtractor()
	path := filepath.Join(testdataWordsDir, "G_格式混乱_S020_潘丽_联邦学习推荐_R020.docx")
	content, err := e.Extract(path)
	require.NoError(t, err)
	require.NotEmpty(t, content)
	t.Logf("G-quality (malformed) docx: %d chars", len([]rune(content)))
}

func TestDOCXExtractor_Extract_MixedQualityClasses(t *testing.T) {
	skipIfNoDir(t, testdataWordsDir)
	samples := map[string]string{
		"A-优秀":   "A_优秀_S001_张明_语音情感识别_R001.docx",
		"B-良好":   "B_良好_S004_赵强_语音情感识别_R004.docx",
		"C-合格":   "C_合格_S009_吴涛_语音情感识别_R009.docx",
		"D-不合格":  "D_不合格_S013_马超_联邦学习推荐_R013.docx",
		"E-答非所问": "E_答非所问_S016_徐婷_语音情感识别_R016.docx",
		"F-内容空洞": "F_内容空洞_S018_曹雪_GNN影响力分析_R018.docx",
		"G-格式混乱": "G_格式混乱_S020_潘丽_联邦学习推荐_R020.docx",
	}

	e := NewDOCXExtractor()
	for label, name := range samples {
		t.Run(label, func(t *testing.T) {
			content, err := e.Extract(filepath.Join(testdataWordsDir, name))
			require.NoError(t, err, "extraction must succeed for %s", label)
			assert.NotEmpty(t, content, "content must not be empty for %s", label)
			t.Logf("%s: %d chars", label, len([]rune(content)))
		})
	}
}

func TestDOCXExtractor_Extract_AllFiles(t *testing.T) {
	skipIfNoDir(t, testdataWordsDir)
	if testing.Short() {
		t.Skip("skipping batch test in short mode")
	}

	matches, err := filepath.Glob(filepath.Join(testdataWordsDir, "*.docx"))
	require.NoError(t, err)
	require.NotEmpty(t, matches, "no DOCX files found in testdata/words")

	t.Logf("Found %d DOCX files", len(matches))

	e := NewDOCXExtractor()
	failed := make([]string, 0)
	var totalChars int
	var minChars, maxChars int = int(^uint(0) >> 1), 0

	for _, path := range matches {
		content, err := e.Extract(path)
		if err != nil {
			failed = append(failed, filepath.Base(path)+": "+err.Error())
			continue
		}
		n := len([]rune(content))
		totalChars += n
		if n < minChars {
			minChars = n
		}
		if n > maxChars {
			maxChars = n
		}
	}

	successCount := len(matches) - len(failed)
	t.Logf("Total: %d, Success: %d, Failed: %d", len(matches), successCount, len(failed))
	t.Logf("Char count — min: %d, max: %d, avg: %d", minChars, maxChars, totalChars/max(successCount, 1))

	for _, f := range failed {
		t.Errorf("FAIL: %s", f)
	}

	assert.Equal(t, len(matches), successCount, "all DOCX files must extract successfully")
}

func TestExtractorChain_Extract_DOCX(t *testing.T) {
	skipIfNoDir(t, testdataWordsDir)
	chain := NewExtractorChain()
	path := filepath.Join(testdataWordsDir, "B_良好_S005_刘洋_联邦学习推荐_R005.docx")
	result := chain.Extract(path)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Content)
	assert.True(t, result.Confidence > 0.9, "confidence should be high for matched extension")
	assert.Contains(t, result.Metadata, "DOCXExtractor")
}

func TestExtractorChain_AllDOCXs(t *testing.T) {
	skipIfNoDir(t, testdataWordsDir)
	if testing.Short() {
		t.Skip("skipping chain batch test in short mode")
	}

	matches, err := filepath.Glob(filepath.Join(testdataWordsDir, "*.docx"))
	require.NoError(t, err)

	chain := NewExtractorChain()
	failed := 0

	for _, path := range matches {
		result := chain.Extract(path)
		if result == nil || strings.TrimSpace(result.Content) == "" {
			t.Errorf("chain failed: %s", filepath.Base(path))
			failed++
		}
	}
	assert.Zero(t, failed, "chain must extract all DOCX files")
}

func TestExtractorChain_MixedFormats(t *testing.T) {
	skipIfNoDir(t, testdataWordsDir)
	chain := NewExtractorChain()

	// DOCX
	docxResult := chain.Extract(filepath.Join(testdataWordsDir, "A_优秀_S001_张明_语音情感识别_R001.docx"))
	require.NotNil(t, docxResult)
	assert.Contains(t, docxResult.Metadata, "DOCXExtractor")

	// PDF (use the pdf testdata dir from pdf_test.go)
	pdfResult := chain.Extract(filepath.Join("../../testdata/pdfs", "A_S001_张明化_语音情感识别_R001.pdf"))
	require.NotNil(t, pdfResult)
	assert.Contains(t, pdfResult.Metadata, "PDFExtractor")

	t.Logf("DOCX: %d chars (via %s)", len([]rune(docxResult.Content)), docxResult.Metadata)
	t.Logf("PDF: %d chars (via %s)", len([]rune(pdfResult.Content)), pdfResult.Metadata)
}

func BenchmarkDOCXExtractor_Single(b *testing.B) {
	e := NewDOCXExtractor()
	path := filepath.Join(testdataWordsDir, "A_优秀_S001_张明_语音情感识别_R001.docx")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Extract(path)
	}
}

func BenchmarkDOCXExtractor_SmallFile(b *testing.B) {
	e := NewDOCXExtractor()
	path := filepath.Join(testdataWordsDir, "F_内容空洞_S019_朱明_语音情感识别_R019.docx")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Extract(path)
	}
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}
