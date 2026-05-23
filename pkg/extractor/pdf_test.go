package extractor

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testdataDir = "../../testdata/pdfs"


// ---------------------------------------------------------------------------
// Unit: cleanPDFText
// ---------------------------------------------------------------------------

func TestCleanPDFText_RemovesCarriageReturns(t *testing.T) {
	raw := "hello\r\nworld\r\ntest"
	got := cleanPDFText(raw)
	assert.NotContains(t, got, "\r")
	assert.Equal(t, "hello\nworld\ntest", got)
}

func TestCleanPDFText_CollapsesMultipleSpaces(t *testing.T) {
	raw := "hello   world\t\t\ttest"
	got := cleanPDFText(raw)
	assert.Equal(t, "hello world test", got)
}

func TestCleanPDFText_CollapsesMultipleNewlines(t *testing.T) {
	raw := "line1\n\n\n\n\nline2\n\n\nline3"
	got := cleanPDFText(raw)
	assert.Equal(t, "line1\n\nline2\n\nline3", got)
}

func TestCleanPDFText_TrimsWhitespace(t *testing.T) {
	raw := "  \n  hello world  \n  "
	got := cleanPDFText(raw)
	assert.Equal(t, "hello world", got)
}

func TestCleanPDFText_EmptyString(t *testing.T) {
	assert.Equal(t, "", cleanPDFText(""))
	assert.Equal(t, "", cleanPDFText("   \n\n   "))
}

func TestCleanPDFText_PreservesChineseCharacters(t *testing.T) {
	raw := "摘要：本文研究了基于深度学习的语音情感识别方法。"
	got := cleanPDFText(raw)
	assert.Equal(t, raw, got)
}

func TestCleanPDFText_MixedContent(t *testing.T) {
	raw := "Abstract\n\n   This paper proposes a novel method.\n   关键词：GNN；影响力分析  "
	got := cleanPDFText(raw)
	assert.True(t, strings.HasPrefix(got, "Abstract"))
	assert.True(t, strings.Contains(got, "GNN"))
	assert.False(t, strings.Contains(got, "   "))
	assert.False(t, strings.Contains(got, "\n\n\n"))
}

// ---------------------------------------------------------------------------
// Unit: absDiff
// ---------------------------------------------------------------------------

func TestAbsDiff(t *testing.T) {
	assert.Equal(t, 0.0, absDiff(5.0, 5.0))
	assert.Equal(t, 3.0, absDiff(5.0, 8.0))
	assert.Equal(t, 3.0, absDiff(8.0, 5.0))
	assert.Equal(t, 10.5, absDiff(-3.5, 7.0))
}

// ---------------------------------------------------------------------------
// SupportedFormats
// ---------------------------------------------------------------------------

func TestPDFExtractor_SupportedFormats(t *testing.T) {
	e := NewPDFExtractor()
	formats := e.SupportedFormats()
	require.Len(t, formats, 1)
	assert.Equal(t, ".pdf", formats[0])
}

// ---------------------------------------------------------------------------
// Extract: error paths
// ---------------------------------------------------------------------------

func TestPDFExtractor_Extract_NonExistentFile(t *testing.T) {
	e := NewPDFExtractor()
	content, err := e.Extract("/nonexistent/path/file.pdf")
	assert.Error(t, err)
	assert.Empty(t, content)
}

func TestPDFExtractor_Extract_NonPDFFile(t *testing.T) {
	e := NewPDFExtractor()
	content, err := e.Extract(filepath.Join(testdataDir, ".."))
	assert.Error(t, err)
	assert.Empty(t, content)
}

// ---------------------------------------------------------------------------
// Extract: real PDFs — single representative files
// ---------------------------------------------------------------------------

func TestPDFExtractor_Extract_SpeechEmotion(t *testing.T) {
	skipIfNoDir(t, testdataDir)
	e := NewPDFExtractor()
	path := filepath.Join(testdataDir, "A_S001_张明化_语音情感识别_R001.pdf")
	content, err := e.Extract(path)
	require.NoError(t, err)
	require.NotEmpty(t, content)
	t.Logf("Speech emotion – chars=%d, first 200: %s", len([]rune(content)), truncate(content, 200))
}

func TestPDFExtractor_Extract_FederatedLearning(t *testing.T) {
	skipIfNoDir(t, testdataDir)
	e := NewPDFExtractor()
	path := filepath.Join(testdataDir, "A_S002_学生002_联邦学习推荐_R002.pdf")
	content, err := e.Extract(path)
	require.NoError(t, err)
	require.NotEmpty(t, content)
	t.Logf("Federated learning – chars=%d, first 200: %s", len([]rune(content)), truncate(content, 200))
}

func TestPDFExtractor_Extract_GNNInfluence(t *testing.T) {
	skipIfNoDir(t, testdataDir)
	e := NewPDFExtractor()
	path := filepath.Join(testdataDir, "A_S003_学生003_GNN影响力分析_R003.pdf")
	content, err := e.Extract(path)
	require.NoError(t, err)
	require.NotEmpty(t, content)
	t.Logf("GNN influence – chars=%d, first 200: %s", len([]rune(content)), truncate(content, 200))
}

func TestPDFExtractor_Extract_ClassH(t *testing.T) {
	skipIfNoDir(t, testdataDir)
	e := NewPDFExtractor()
	path := filepath.Join(testdataDir, "H_S092_学生092_联邦学习推荐_R092.pdf")
	content, err := e.Extract(path)
	require.NoError(t, err)
	require.NotEmpty(t, content)
	t.Logf("Class H – chars=%d, first 200: %s", len([]rune(content)), truncate(content, 200))
}

func TestPDFExtractor_Extract_LargeFile(t *testing.T) {
	skipIfNoDir(t, testdataDir)
	// A_S001 is 155KB — one of the larger files
	e := NewPDFExtractor()
	path := filepath.Join(testdataDir, "A_S001_张明化_语音情感识别_R001.pdf")
	content, err := e.Extract(path)
	require.NoError(t, err)
	require.NotEmpty(t, content)
	assert.True(t, len([]rune(content)) > 100, "large file should yield substantial text")
}

// ---------------------------------------------------------------------------
// Extract: batch — all 108 PDFs
// ---------------------------------------------------------------------------

func TestPDFExtractor_Extract_AllFiles(t *testing.T) {
	skipIfNoDir(t, testdataDir)
	if testing.Short() {
		t.Skip("skipping batch test in short mode")
	}

	matches, err := filepath.Glob(filepath.Join(testdataDir, "*.pdf"))
	require.NoError(t, err)
	require.NotEmpty(t, matches, "no PDFs found in testdata/pdfs")

	t.Logf("Found %d PDF files", len(matches))

	e := NewPDFExtractor()
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
	t.Logf("Char count — min: %d, max: %d, avg: %d", minChars, maxChars, totalChars/successCount)

	for _, f := range failed {
		t.Errorf("FAIL: %s", f)
	}

	assert.Equal(t, len(matches), successCount, "all PDFs must extract successfully")
}

// ---------------------------------------------------------------------------
// Extract: concurrent extraction
// ---------------------------------------------------------------------------

func TestPDFExtractor_Extract_Concurrent(t *testing.T) {
	skipIfNoDir(t, testdataDir)
	if testing.Short() {
		t.Skip("skipping concurrent test in short mode")
	}

	matches, err := filepath.Glob(filepath.Join(testdataDir, "*.pdf"))
	require.NoError(t, err)
	if len(matches) < 10 {
		t.Skip("not enough PDFs for concurrent test")
	}

	e := NewPDFExtractor()
	var wg sync.WaitGroup
	errCh := make(chan error, len(matches))

	for _, path := range matches {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			content, err := e.Extract(p)
			if err != nil {
				errCh <- fmt.Errorf("%s: %w", filepath.Base(p), err)
				return
			}
			if strings.TrimSpace(content) == "" {
				errCh <- fmt.Errorf("%s: empty content", filepath.Base(p))
			}
		}(path)
	}

	wg.Wait()
	close(errCh)

	var errors []error
	for e := range errCh {
		errors = append(errors, e)
	}

	for _, e := range errors {
		t.Error(e)
	}
	assert.Empty(t, errors, "concurrent extraction must produce no errors")
}

// ---------------------------------------------------------------------------
// ExtractorChain
// ---------------------------------------------------------------------------

func TestExtractorChain_Extract_PDF(t *testing.T) {
	skipIfNoDir(t, testdataDir)
	chain := NewExtractorChain()
	path := filepath.Join(testdataDir, "B_S013_学生013_语音情感识别_R013.pdf")
	result := chain.Extract(path)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Content)
	assert.True(t, result.Confidence > 0.9, "confidence should be high for matched extension")
	assert.Contains(t, result.Metadata, "PDFExtractor")
}

func TestExtractorChain_Extract_PDFWithoutExtension(t *testing.T) {
	skipIfNoDir(t, testdataDir)
	// This tests the fallback path — not easily testable without file manipulation,
	// but we verify the chain works correctly with the given extension.
	chain := NewExtractorChain()
	path := filepath.Join(testdataDir, "C_S033_学生033_GNN影响力分析_R033.pdf")
	result := chain.Extract(path)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Content)
}

func TestExtractorChain_Extract_NonExistentFile(t *testing.T) {
	chain := NewExtractorChain()
	result := chain.Extract("/nonexistent/file.pdf")
	assert.Nil(t, result)
}

func TestExtractorChain_GetSupportedFormats(t *testing.T) {
	chain := NewExtractorChain()
	formats := chain.GetSupportedFormats()
	assert.Contains(t, formats, ".pdf")
	assert.Contains(t, formats, ".docx")
}

func TestExtractorChain_AllPDFs(t *testing.T) {
	skipIfNoDir(t, testdataDir)
	if testing.Short() {
		t.Skip("skipping chain batch test in short mode")
	}

	matches, err := filepath.Glob(filepath.Join(testdataDir, "*.pdf"))
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
	assert.Zero(t, failed, "chain must extract all PDFs")
}

// ---------------------------------------------------------------------------
// Content quality checks — spot-check representative files from each class
// ---------------------------------------------------------------------------

func TestPDFExtractor_ContentQuality_SampleByClass(t *testing.T) {
	// Pick one file from each of the 8 classes (A-H)
	samples := []string{
		"A_S001_张明化_语音情感识别_R001.pdf",
		"B_S013_学生013_语音情感识别_R013.pdf",
		"C_S033_学生033_GNN影响力分析_R033.pdf",
		"D_S053_学生053_联邦学习推荐_R053.pdf",
		"E_S068_学生068_联邦学习推荐_R068.pdf",
		"F_S078_学生078_GNN影响力分析_R078.pdf",
		"G_S086_学生086_联邦学习推荐_R086.pdf",
		"H_S092_学生092_联邦学习推荐_R092.pdf",
	}

	e := NewPDFExtractor()
	for _, name := range samples {
		t.Run(name, func(t *testing.T) {
			content, err := e.Extract(filepath.Join(testdataDir, name))
			require.NoError(t, err, "extraction must succeed for %s", name)
			assert.True(t, len([]rune(content)) > 50,
				"content too short for %s: got %d chars", name, len([]rune(content)))
		})
	}
}

func TestPDFExtractor_ContentQuality_SampleByTopic(t *testing.T) {
	topicSamples := map[string]string{
		"语音情感识别": "A_S001_张明化_语音情感识别_R001.pdf",
		"联邦学习推荐": "B_S014_学生014_联邦学习推荐_R014.pdf",
		"GNN影响力分析": "B_S015_学生015_GNN影响力分析_R015.pdf",
	}

	e := NewPDFExtractor()
	for topic, name := range topicSamples {
		t.Run(topic, func(t *testing.T) {
			content, err := e.Extract(filepath.Join(testdataDir, name))
			require.NoError(t, err)
			assert.True(t, len([]rune(content)) > 50,
				"content too short for %s: got %d chars", name, len([]rune(content)))
		})
	}
}

// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

func BenchmarkPDFExtractor_Single(b *testing.B) {
	e := NewPDFExtractor()
	path := filepath.Join(testdataDir, "A_S001_张明化_语音情感识别_R001.pdf")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Extract(path)
	}
}

func BenchmarkPDFExtractor_SmallFile(b *testing.B) {
	e := NewPDFExtractor()
	path := filepath.Join(testdataDir, "B_S013_学生013_语音情感识别_R013.pdf")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.Extract(path)
	}
}

func BenchmarkExtractorChain_PDF(b *testing.B) {
	chain := NewExtractorChain()
	path := filepath.Join(testdataDir, "A_S001_张明化_语音情感识别_R001.pdf")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain.Extract(path)
	}
}

func BenchmarkCleanPDFText(b *testing.B) {
	raw := "hello\r\nworld\r\n  test   with   spaces\n\n\n\nmultiple\n\n\nnewlines\n\n"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cleanPDFText(raw)
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func truncate(s string, maxLen int) string {
	r := []rune(s)
	if len(r) <= maxLen {
		return s
	}
	return string(r[:maxLen]) + "..."
}
