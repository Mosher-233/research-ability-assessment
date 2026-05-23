package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"

	"github.com/Mosher-233/research-ability-assessment/internal/llm"
	"github.com/Mosher-233/research-ability-assessment/internal/models"
	"github.com/Mosher-233/research-ability-assessment/internal/service"
)

type EvidenceAgent struct {
	evidenceService *service.EvidenceService
	llmClient       *llm.Client
}

func NewEvidenceAgent(evidenceService *service.EvidenceService) *EvidenceAgent {
	return &EvidenceAgent{evidenceService: evidenceService}
}

func NewEvidenceAgentWithLLM(evidenceService *service.EvidenceService, llmClient *llm.Client) *EvidenceAgent {
	return &EvidenceAgent{evidenceService: evidenceService, llmClient: llmClient}
}

func (a *EvidenceAgent) SetLLMClient(llmClient *llm.Client) {
	a.llmClient = llmClient
}

func (a *EvidenceAgent) CollectEvidence(ctx context.Context, studentID string, taskID string) ([]*models.Evidence, error) {
	evidences, err := a.evidenceService.GetEvidencesByStudentAndTask(ctx, studentID, taskID)
	if err != nil {
		return nil, err
	}

	evidencePtrs := make([]*models.Evidence, len(evidences))
	for i := range evidences {
		evidencePtrs[i] = &evidences[i]
	}

	return evidencePtrs, nil
}

type KBMInfo struct {
	KBMName     string
	Level       int
	Credibility float64
	Rationale   string
}

type EvidenceReport struct {
	KBMInfo
	CleanedText     string
	MatchedKeywords []string
	LengthScore     int
	KeywordRichness int
	StructureScore  int
	TotalScore      int
}

func (a *EvidenceAgent) PreprocessEvidence(text string) string {
	text = strings.TrimSpace(text)
	var b strings.Builder
	prevSpace := false
	for _, r := range text {
		switch {
		case r == '\r':
			continue
		case r == '\n':
			b.WriteRune('\n')
			prevSpace = false
		case unicode.IsSpace(r):
			if !prevSpace {
				b.WriteRune(' ')
				prevSpace = true
			}
		default:
			b.WriteRune(r)
			prevSpace = false
		}
	}
	cleaned := b.String()
	cleaned = regexp.MustCompile(`\n{3,}`).ReplaceAllString(cleaned, "\n\n")

	runes := []rune(cleaned)
	if len(runes) > 2000 {
		cleaned = string(runes[:2000])
	}
	return cleaned
}

var kbmKeywords = map[string][]string{
	"文献检索策略":   {"PICO", "PRISMA", "双人独立筛选", "纳入", "检索式", "核心文献", "筛选标准", "Kappa", "命中数"},
	"文献综述质量":   {"系统综述", "综述", "梳理", "归纳", "对比表格", "研究空白", "代表性论文", "研究现状", "文献综述", "多维度"},
	"文献批判性分析":  {"局限", "不足", "缺陷", "批判", "质疑", "方法论", "样本规模", "偏颇", "可进一步"},
	"实验方案合理性":  {"对照实验", "实验组", "对照组", "交叉验证", "数据划分", "训练集", "测试集", "验证集", "消融实验", "基线模型", "折交叉验证", "实验方案", "实验设计", "析因", "因素方差"},
	"变量控制":     {"变量控制", "控制变量", "随机种子", "统一配置", "重复实验", "保持一致性", "校准", "预热时间"},
	"实验实施质量":   {"成功率", "168小时", "连续运行", "断点续传", "部署", "传感器"},
	"数据分析方法选择": {"SHAP", "特征重要性", "混淆矩阵", "t检验", "方差分析", "ANOVA", "ROC-AUC", "F1-score", "McNemar", "loss曲线", "显著性", "p值", "相关系数", "回归", "聚类"},
	"结果解释准确性":  {"实验结果表明", "导致", "原因是", "分析发现", "讨论", "局限于", "偏差", "异常", "泛化"},
	"问题提出新颖性":  {"首次提出", "创新", "新颖", "新视角", "新指标", "与众不同", "原创", "独特", "研究问题", "研究假设", "前人未", "不同于现有"},
	"解决方案原创性":  {"架构", "框架", "编码器", "门控", "多模态", "融合", "流水线", "注意力机制", "知识蒸馏", "联邦学习"},
}

var kbmLabels = []string{
	"文献检索策略", "文献综述质量", "文献批判性分析",
	"实验方案合理性", "变量控制", "实验实施质量",
	"数据分析方法选择", "结果解释准确性",
	"问题提出新颖性", "解决方案原创性",
}

func (a *EvidenceAgent) ClassifyEvidence(text string) (string, []string) {
	// Try LLM classification first
	if a.llmClient != nil && len(strings.TrimSpace(text)) > 20 {
		kbmName, matched, err := a.classifyWithLLM(text)
		if err == nil && kbmName != "" {
			return kbmName, matched
		}
		log.Printf("EvidenceAgent: LLM分类失败，回退到关键词匹配: %v", err)
	}
	return a.classifyWithKeywords(text)
}

func (a *EvidenceAgent) classifyWithLLM(text string) (string, []string, error) {
	sb := &strings.Builder{}
	sb.WriteString("请将以下学生研究证据归类到最合适的KBM类别中。\n\n")
	sb.WriteString("可选类别：\n")
	for _, label := range kbmLabels {
		dimName := models.GetDimensionName(models.GetDimensionByKBM(label))
		sb.WriteString(fmt.Sprintf("- %s (所属维度: %s)\n", label, dimName))
	}
	sb.WriteString("\n输出格式（严格JSON）：\n")
	sb.WriteString(`{"kbm_name": "最匹配的类别", "rationale": "归类依据", "matched_keywords": ["关键概念1", "关键概念2"]}`)
	sb.WriteString(fmt.Sprintf("\n\n证据内容（截取前2000字符）：\n%s", text[:min(len(text), 2000)]))

	messages := []llm.Message{
		{Role: "system", Content: "你是一个研究证据分类专家。请根据证据内容准确分类到预定义的KBM类别。"},
		{Role: "user", Content: sb.String()},
	}

	response, err := a.llmClient.Chat(context.Background(), messages)
	if err != nil {
		return "", nil, err
	}

	jsonStr := extractJSON(response)
	if jsonStr == "" {
		return "", nil, fmt.Errorf("LLM分类响应中无有效JSON")
	}

	var result struct {
		KBMName         string   `json:"kbm_name"`
		Rationale       string   `json:"rationale"`
		MatchedKeywords []string `json:"matched_keywords"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return "", nil, err
	}

	// Validate the KBM name
	for _, label := range kbmLabels {
		if label == result.KBMName {
			return result.KBMName, result.MatchedKeywords, nil
		}
	}

	return "", nil, fmt.Errorf("LLM返回了无效的KBM类别: %s", result.KBMName)
}

func (a *EvidenceAgent) classifyWithKeywords(text string) (string, []string) {
	textLower := strings.ToLower(text)
	bestKBM := ""
	bestScore := 0
	var matched []string

	for kbm, keywords := range kbmKeywords {
		score := 0
		var kbmMatched []string
		for _, kw := range keywords {
			if strings.Contains(textLower, strings.ToLower(kw)) {
				score++
				kbmMatched = append(kbmMatched, kw)
			}
		}
		if score > bestScore {
			bestScore = score
			bestKBM = kbm
			matched = kbmMatched
		}
	}

	if bestScore == 0 || bestKBM == "" {
		return "文献检索策略", nil
	}
	return bestKBM, matched
}

func (a *EvidenceAgent) ExtractKBMInfo(text string) *KBMInfo {
	cleaned := a.PreprocessEvidence(text)

	// Try LLM-based assessment
	if a.llmClient != nil && len(strings.TrimSpace(cleaned)) > 20 {
		info, err := a.extractKBMWithLLM(cleaned)
		if err == nil {
			return info
		}
		log.Printf("EvidenceAgent: LLM评估失败，回退到规则引擎: %v", err)
	}

	kbmName, matched := a.ClassifyEvidence(cleaned)
	level, credibility, rationale := a.assessEvidenceWithRules(cleaned, kbmName, matched)
	return &KBMInfo{KBMName: kbmName, Level: level, Credibility: credibility, Rationale: rationale}
}

func (a *EvidenceAgent) extractKBMWithLLM(text string) (*KBMInfo, error) {
	sb := &strings.Builder{}
	sb.WriteString("请评估以下学生证据的质量，给出KBM分类和等级评定。\n\n")
	sb.WriteString("可选KBM类别：\n")
	for _, label := range kbmLabels {
		sb.WriteString(fmt.Sprintf("- %s\n", label))
	}
	sb.WriteString("\n等级标准：\n")
	sb.WriteString("1级(不合格): 内容浅薄，缺乏关键概念，论证不足\n")
	sb.WriteString("2级(合格): 内容基本完整，涉及一些关键概念，论证尚可\n")
	sb.WriteString("3级(良好): 内容充实，关键概念运用得当，论证较为严谨\n")
	sb.WriteString("4级(优秀): 内容深入，关键概念运用精准，论证严谨且有创新性\n\n")
	sb.WriteString("输出JSON格式：\n")
	sb.WriteString(`{"kbm_name": "类别", "level": 3, "credibility": 0.85, "rationale": "详细的评估理由，引用证据中的具体内容作为依据"}`)
	sb.WriteString(fmt.Sprintf("\n\n证据内容（截取前2000字符）：\n%s", text[:min(len(text), 2000)]))

	messages := []llm.Message{
		{Role: "system", Content: "你是一个研究证据质量评估专家。请根据证据内容和等级标准进行准确评估。"},
		{Role: "user", Content: sb.String()},
	}

	response, err := a.llmClient.Chat(context.Background(), messages)
	if err != nil {
		return nil, err
	}

	jsonStr := extractJSON(response)
	if jsonStr == "" {
		return nil, fmt.Errorf("LLM响应中无有效JSON")
	}

	var result struct {
		KBMName      string  `json:"kbm_name"`
		Level        int     `json:"level"`
		Credibility  float64 `json:"credibility"`
		Rationale    string  `json:"rationale"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, err
	}

	if result.Level < 1 {
		result.Level = 1
	}
	if result.Level > 4 {
		result.Level = 4
	}
	if result.Credibility < 0 {
		result.Credibility = 0.5
	}
	if result.Credibility > 1 {
		result.Credibility = 1.0
	}

	return &KBMInfo{
		KBMName:     result.KBMName,
		Level:       result.Level,
		Credibility: result.Credibility,
		Rationale:   result.Rationale,
	}, nil
}

func (a *EvidenceAgent) assessEvidenceWithRules(text string, kbmName string, matched []string) (level int, credibility float64, rationale string) {
	textLower := strings.ToLower(text)
	runes := []rune(text)
	textLen := len(runes)

	lengthScore := 1
	switch {
	case textLen >= 300:
		lengthScore = 4
	case textLen >= 150:
		lengthScore = 3
	case textLen >= 50:
		lengthScore = 2
	}

	keywordRichness := 0
	if kw, ok := kbmKeywords[kbmName]; ok {
		for _, k := range kw {
			if strings.Contains(textLower, strings.ToLower(k)) {
				keywordRichness++
			}
		}
	}

	structureScore := 1
	if strings.Contains(text, "：") || strings.Contains(text, ":") {
		structureScore = 2
	}
	if strings.Contains(text, "1") || strings.Contains(text, "（1）") || strings.Contains(text, "(1)") {
		structureScore = 3
	}

	keywordWeight := keywordRichness / 2
	if keywordWeight > 5 {
		keywordWeight = 5
	}
	total := lengthScore + keywordWeight + structureScore

	switch {
	case total >= 8:
		level = 4
	case total >= 5:
		level = 3
	case total >= 3:
		level = 2
	default:
		level = 1
	}

	credibility = a.calcCredibility(textLen > 0, keywordRichness, structureScore, total)
	rationale = a.buildRationale(kbmName, level, matched, lengthScore, keywordRichness, structureScore, total, credibility)
	return
}

func (a *EvidenceAgent) calcCredibility(hasContent bool, keywordRichness, structureScore, total int) float64 {
	if !hasContent {
		return 0.0
	}
	kwFactor := float64(keywordRichness) / 6.0
	if kwFactor > 1.0 {
		kwFactor = 1.0
	}
	strFactor := float64(structureScore) / 3.0
	totFactor := float64(total) / 11.0
	if totFactor > 1.0 {
		totFactor = 1.0
	}
	return (kwFactor*0.5 + strFactor*0.2 + totFactor*0.3)
}

func (a *EvidenceAgent) buildRationale(kbmName string, level int, matched []string, lengthScore, keywordRichness, structureScore, total int, credibility float64) string {
	levelLabels := []string{"", "不合格", "合格", "良好", "优秀"}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("分类为「%s」，等级评定为「%s」(L%d)[规则引擎]。", kbmName, levelLabels[level], level))
	if len(matched) > 0 {
		sb.WriteString(fmt.Sprintf("匹配关键词(%d个): ", len(matched)))
		for i, kw := range matched {
			if i > 0 {
				sb.WriteString(", ")
			}
			if i >= 4 {
				sb.WriteString(fmt.Sprintf("...共%d个", len(matched)))
				break
			}
			sb.WriteString(kw)
		}
		sb.WriteString("。")
	}
	sb.WriteString(fmt.Sprintf("评分因子: 长度%d + 关键词丰富度%d + 结构%d = 综合%d。", lengthScore, keywordRichness, structureScore, total))
	sb.WriteString(fmt.Sprintf("可信度: %.0f%%。", credibility*100))
	return sb.String()
}

func (a *EvidenceAgent) AutoClassifyAndUpdate(ctx context.Context, evidence *models.Evidence) error {
	if evidence.KBMName != "" {
		return nil
	}
	info := a.ExtractKBMInfo(evidence.Content)
	evidence.KBMName = info.KBMName
	evidence.KBMLevel = info.Level
	return nil
}
