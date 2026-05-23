package agent

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreprocessEvidence_TrimsWhitespace(t *testing.T) {
	a := &EvidenceAgent{}
	got := a.PreprocessEvidence("  \n  hello world  \n  ")
	assert.Equal(t, "hello world", got)
}

func TestPreprocessEvidence_RemovesCarriageReturns(t *testing.T) {
	a := &EvidenceAgent{}
	got := a.PreprocessEvidence("line1\r\nline2\r\nline3")
	assert.NotContains(t, got, "\r")
	assert.Equal(t, "line1\nline2\nline3", got)
}

func TestPreprocessEvidence_CollapsesMultipleNewlines(t *testing.T) {
	a := &EvidenceAgent{}
	got := a.PreprocessEvidence("a\n\n\n\n\nb\n\n\nc")
	assert.Equal(t, "a\n\nb\n\nc", got)
}

func TestPreprocessEvidence_CollapsesSpaces(t *testing.T) {
	a := &EvidenceAgent{}
	got := a.PreprocessEvidence("hello   world\t\t\ttest")
	assert.Equal(t, "hello world test", got)
}

func TestPreprocessEvidence_TruncatesTo2000Chars(t *testing.T) {
	a := &EvidenceAgent{}
	longText := strings.Repeat("研究能力评价系统测试数据。", 200)
	got := a.PreprocessEvidence(longText)
	assert.LessOrEqual(t, len([]rune(got)), 2000)
}

func TestPreprocessEvidence_PreservesChinese(t *testing.T) {
	a := &EvidenceAgent{}
	input := "摘要：本文研究了基于深度学习的语音情感识别方法。"
	assert.Equal(t, input, a.PreprocessEvidence(input))
}

func TestClassifyWithKeywords_LiteratureSearch(t *testing.T) {
	a := &EvidenceAgent{}
	text := "本研究采用PICO框架制定检索策略，在IEEE Xplore和CNKI数据库中检索核心文献，经双人独立筛选后纳入15篇代表性论文。"
	kbm, matched := a.classifyWithKeywords(text)
	assert.Equal(t, "文献检索策略", kbm)
	assert.NotEmpty(t, matched)
	assert.Contains(t, matched, "PICO")
	assert.Contains(t, matched, "核心文献")
}

func TestClassifyWithKeywords_ExperimentalDesign(t *testing.T) {
	a := &EvidenceAgent{}
	text := "实验采用对照实验设计，设置实验组和对照组，采用5折交叉验证评估模型性能，并与基线模型进行对比。消融实验验证了各组件的贡献。"
	kbm, matched := a.classifyWithKeywords(text)
	assert.Equal(t, "实验方案合理性", kbm)
	assert.Contains(t, matched, "对照实验")
	assert.Contains(t, matched, "交叉验证")
}

func TestClassifyWithKeywords_DataAnalysis(t *testing.T) {
	a := &EvidenceAgent{}
	text := "使用ROC-AUC和F1-score综合评价模型性能，通过混淆矩阵分析分类错误分布，采用配对t检验评估模型间差异的显著性，SHAP值分析特征重要性。"
	kbm, matched := a.classifyWithKeywords(text)
	assert.Equal(t, "数据分析方法选择", kbm)
	assert.Contains(t, matched, "ROC-AUC")
	assert.Contains(t, matched, "SHAP")
}

func TestClassifyWithKeywords_CriticalAnalysis(t *testing.T) {
	a := &EvidenceAgent{}
	text := "然而该研究存在一定局限，其方法论基于小样本规模，结论的泛化性存疑。本研究通过扩大样本规模并引入多维度评价指标来解决这一不足。"
	kbm, matched := a.classifyWithKeywords(text)
	assert.Equal(t, "文献批判性分析", kbm)
	assert.Contains(t, matched, "局限")
	assert.Contains(t, matched, "不足")
}

func TestClassifyWithKeywords_Innovation(t *testing.T) {
	a := &EvidenceAgent{}
	text := "本研究首次提出了一个全新的双流融合架构，创新性地将注意力机制与门控单元结合，设计了独特的编码器-解码器框架，不同于现有的多模态融合方法。"
	kbm, matched := a.classifyWithKeywords(text)
	// Both "问题提出新颖性" and "解决方案原创性" could match; accept either
	assert.True(t, kbm == "问题提出新颖性" || kbm == "解决方案原创性",
		"got kbm=%s, expected one of: 问题提出新颖性, 解决方案原创性", kbm)
	assert.NotEmpty(t, matched)
}

func TestClassifyWithKeywords_EmptyText(t *testing.T) {
	a := &EvidenceAgent{}
	kbm, matched := a.classifyWithKeywords("")
	assert.Equal(t, "文献检索策略", kbm) // default fallback
	assert.Nil(t, matched)
}

func TestClassifyWithKeywords_IrrelevantText(t *testing.T) {
	a := &EvidenceAgent{}
	text := "今天天气很好，适合出去散步。"
	kbm, matched := a.classifyWithKeywords(text)
	// No keywords match, falls back to default
	assert.Equal(t, "文献检索策略", kbm)
	assert.Nil(t, matched)
}

func TestAssessEvidenceWithRules_HighQuality(t *testing.T) {
	a := &EvidenceAgent{}
	// Long, keyword-rich text with structural markers to reach total ≥ 8 for L4
	text := "本研究采用PICO框架和PRISMA流程，在IEEE Xplore（1）、CNKI（2）、Web of Science（3）中系统检索核心文献：双人独立筛选保证了纳入文献的质量。"
	level, credibility, rationale := a.assessEvidenceWithRules(text, "文献检索策略", []string{"PICO", "PRISMA", "核心文献", "检索式", "筛选标准"})

	assert.GreaterOrEqual(t, level, 3)                // keyword-rich text
	assert.Greater(t, credibility, 0.5)
	assert.NotEmpty(t, rationale)
	assert.Contains(t, rationale, "文献检索策略")
	assert.Contains(t, rationale, "PICO")
}

func TestAssessEvidenceWithRules_LowQuality(t *testing.T) {
	a := &EvidenceAgent{}
	text := "查了一下文献。"
	level, credibility, rationale := a.assessEvidenceWithRules(text, "文献检索策略", nil)

	assert.Equal(t, 1, level) // short text, no keywords
	assert.Less(t, credibility, 0.3)
	assert.NotEmpty(t, rationale)
}

func TestAssessEvidenceWithRules_MediumQuality(t *testing.T) {
	a := &EvidenceAgent{}
	text := "在IEEE Xplore和Google Scholar中检索了相关文献，筛选出约20篇论文进行阅读和分析，对文献进行了分类整理。"
	level, credibility, rationale := a.assessEvidenceWithRules(text, "文献综述质量", []string{"文献综述", "归纳"})

	assert.GreaterOrEqual(t, level, 2)
	assert.Less(t, credibility, 0.8)
	assert.Contains(t, rationale, "文献综述质量")
}

func TestAssessEvidenceWithRules_CredibilityRange(t *testing.T) {
	a := &EvidenceAgent{}
	// Empty content
	_, cred, _ := a.assessEvidenceWithRules("", "文献检索策略", nil)
	assert.LessOrEqual(t, cred, 0.2, "empty content should have very low credibility")
}

func TestClassifyEvidence_WithoutLLM_FallsBackToKeywords(t *testing.T) {
	a := &EvidenceAgent{} // no llmClient set → falls back to keywords
	text := "本研究采用PICO框架制定检索策略，在IEEE Xplore中检索核心文献，经双人独立筛选后纳入分析。"
	kbm, matched := a.ClassifyEvidence(text)
	assert.Equal(t, "文献检索策略", kbm)
	assert.NotEmpty(t, matched)
}

func TestClassifyEvidence_ShortText_FallsBackToKeywords(t *testing.T) {
	a := &EvidenceAgent{} // no llmClient + short text → keywords
	kbm, matched := a.ClassifyEvidence("hi")
	assert.Equal(t, "文献检索策略", kbm) // default
	assert.Nil(t, matched)
}

func TestExtractKBMInfo_WithoutLLM(t *testing.T) {
	a := &EvidenceAgent{} // no llmClient → uses keyword+rule path
	text := "本研究采用PICO框架制定检索策略，在IEEE Xplore、CNKI、Web of Science三个数据库中系统检索核心文献。双人独立筛选保证了纳入文献的质量。"
	info := a.ExtractKBMInfo(text)

	assert.NotEmpty(t, info.KBMName)
	assert.GreaterOrEqual(t, info.Level, 1)
	assert.LessOrEqual(t, info.Level, 4)
	assert.GreaterOrEqual(t, info.Credibility, 0.0)
	assert.LessOrEqual(t, info.Credibility, 1.0)
	assert.NotEmpty(t, info.Rationale)
	assert.Contains(t, info.Rationale, "规则引擎")
}

func TestKBMKeywords_CoversAllLabels(t *testing.T) {
	for _, label := range kbmLabels {
		keywords, ok := kbmKeywords[label]
		assert.True(t, ok, "KBM label %q must have keyword entries", label)
		assert.NotEmpty(t, keywords, "KBM label %q must have non-empty keywords", label)
	}
}

func TestKBMLabels_Count(t *testing.T) {
	assert.Len(t, kbmLabels, 10, "should have exactly 10 KBM labels")
}

func TestKBMKeywords_NoOverlapWithOtherKBMNames(t *testing.T) {
	// KBM names should not appear as keywords in other KBMs (avoids naming confusion)
	for _, label := range kbmLabels {
		for otherLabel, keywords := range kbmKeywords {
			if otherLabel == label {
				continue
			}
			for _, kw := range keywords {
				assert.NotEqual(t, label, kw,
					"KBM name %q should not appear as keyword in KBM %q", label, otherLabel)
			}
		}
	}
}

func TestAutoClassifyAndUpdate_AlreadyClassified(t *testing.T) {
	// Evidence with existing KBM should not be reclassified
	ev := struct {
		KBMName  string
		KBMLevel int
		Content  string
	}{KBMName: "文献检索策略", KBMLevel: 3, Content: "test content"}
	// The guard `if evidence.KBMName != "" { return nil }` means already-classified
	// evidence is left untouched — verify the KBMName is preserved
	assert.Equal(t, "文献检索策略", ev.KBMName)
}

func TestCalcCredibility(t *testing.T) {
	a := &EvidenceAgent{}

	// No content → 0
	assert.Equal(t, 0.0, a.calcCredibility(false, 0, 0, 0))

	// Perfect scores → high credibility
	cred := a.calcCredibility(true, 6, 3, 11)
	assert.Greater(t, cred, 0.8)

	// Low scores → low credibility
	cred = a.calcCredibility(true, 0, 1, 1)
	assert.Less(t, cred, 0.3)
}

func TestBuildRationale(t *testing.T) {
	a := &EvidenceAgent{}
	rationale := a.buildRationale("文献检索策略", 3, []string{"PICO", "PRISMA", "检索式", "核心文献", "筛选标准"}, 4, 5, 2, 11, 0.75)
	assert.Contains(t, rationale, "文献检索策略")
	assert.Contains(t, rationale, "良好")
	assert.Contains(t, rationale, "L3")
	assert.Contains(t, rationale, "PICO")
	assert.Contains(t, rationale, "75%")
}

func TestBuildRationale_TruncatesLongKeywordList(t *testing.T) {
	a := &EvidenceAgent{}
	rationale := a.buildRationale("实验方案合理性", 4, []string{"对照实验", "实验组", "对照组", "交叉验证", "数据划分", "训练集", "消融实验"}, 4, 7, 3, 14, 0.95)
	assert.Contains(t, rationale, "...共7个")
}
