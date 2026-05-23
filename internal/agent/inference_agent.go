package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Mosher-233/research-ability-assessment/internal/llm"
	"github.com/Mosher-233/research-ability-assessment/internal/models"
)

type InferenceAgent struct {
	evidenceAgent *EvidenceAgent
	logicUnit     *LogicUnit
	llmClient     *llm.Client
}

type InferenceResult struct {
	OverallScore    float64                          `json:"overall_score"`
	OverallLevel    string                           `json:"overall_level"`
	DimensionScores map[string]models.DimensionScore `json:"dimension_scores"`
	Reasoning       string                           `json:"reasoning"`
	Citations       []CitationInfo                   `json:"citations,omitempty"`
}

// CitationInfo holds a single evidence citation for an evaluation result.
type CitationInfo struct {
	DimensionID    string  `json:"dimension_id"`
	EvidenceID     string  `json:"evidence_id"`
	ExcerptText    string  `json:"excerpt_text"`
	RelevanceScore float64 `json:"relevance_score"`
}

func NewInferenceAgent(evidenceAgent *EvidenceAgent, logicUnit *LogicUnit) *InferenceAgent {
	return &InferenceAgent{
		evidenceAgent: evidenceAgent,
		logicUnit:     logicUnit,
	}
}

func NewInferenceAgentWithLLM(evidenceAgent *EvidenceAgent, logicUnit *LogicUnit, llmClient *llm.Client) *InferenceAgent {
	return &InferenceAgent{
		evidenceAgent: evidenceAgent,
		logicUnit:     logicUnit,
		llmClient:     llmClient,
	}
}

func (a *InferenceAgent) SetLLMClient(llmClient *llm.Client) {
	a.llmClient = llmClient
}

func (a *InferenceAgent) InferAbility(ctx context.Context, studentID string, taskID string) (*InferenceResult, error) {
	evidences, err := a.evidenceAgent.CollectEvidence(ctx, studentID, taskID)
	if err != nil {
		return nil, fmt.Errorf("收集证据失败: %w", err)
	}

	if a.llmClient != nil {
		result, err := a.inferWithLLM(ctx, evidences)
		if err == nil {
			log.Printf("InferenceAgent: LLM推理成功")
			return result, nil
		}
		log.Printf("InferenceAgent: LLM推理失败，回退到规则引擎: %v", err)
	}

	return a.inferWithRules(evidences)
}

func (a *InferenceAgent) inferWithLLM(ctx context.Context, evidences []*models.Evidence) (*InferenceResult, error) {
	systemPrompt := a.buildInferenceSystemPrompt()
	userPrompt := a.buildInferenceUserPrompt(evidences)

	messages := []llm.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	response, err := a.llmClient.Chat(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("LLM调用失败: %w", err)
	}

	return a.parseLLMInferenceResponse(response, evidences)
}

func (a *InferenceAgent) buildInferenceSystemPrompt() string {
	return `你是一个专业的大学生研究能力评价专家。你需要基于学生提交的研究证据，对其研究能力进行多维度综合评估。

## 评估维度与评分量规 (Rubrics)

### 1. 文献综述 (literature_review) - 权重 0.25
- 优秀(90-100): 文献检索策略完善(使用PICO/PRISMA等方法)，综述结构严谨，批判性分析深入，能识别研究空白
- 良好(80-89): 文献检索覆盖面广，综述组织清晰，有一定批判性分析
- 中等(70-79): 文献检索基本完整，综述结构合理，批判性分析较浅
- 及格(60-69): 文献检索范围有限，综述逻辑性不足，缺乏批判性视角
- 不及格(<60): 文献检索随意，综述不成体系，无明显批判性分析

### 2. 研究设计 (research_design) - 权重 0.25
- 优秀(90-100): 实验方案科学严谨，变量控制完善，包含对照/重复/随机化，实施细节清晰
- 良好(80-89): 实验方案合理，变量控制基本到位，实施步骤明确
- 中等(70-79): 实验方案基本可行，部分变量控制不足，实施步骤有待细化
- 及格(60-69): 实验方案存在明显缺陷，变量控制不充分
- 不及格(<60): 实验方案不合理，缺乏基本的变量控制

### 3. 数据分析 (data_analysis) - 权重 0.25
- 优秀(90-100): 分析方法选择恰当且有理论支撑，结果解释准确深入，数据可视化有效
- 良好(80-89): 分析方法选择合理，结果解释清晰，有基本的数据可视化
- 中等(70-79): 分析方法基本正确，结果解释略显浅层或不够严谨
- 及格(60-69): 分析方法选择有待商榷，结果解释存在偏差
- 不及格(<60): 分析方法不当，结果解释错误或严重不足

### 4. 批判性思维 (critical_thinking) - 权重 0.25
- 优秀(90-100): 问题提出新颖有深度，多角度分析透彻，解决方案原创性强
- 良好(80-89): 能独立提出问题，有基本的多角度分析，解决方案有一定创新
- 中等(70-79): 能复述他人观点，初步尝试批判性分析，创新性一般
- 及格(60-69): 以描述为主，缺乏独立思考和批判性分析
- 不及格(<60): 完全没有批判性思考，仅简单罗列信息

## 输出格式
必须严格输出以下JSON格式（不要包含markdown代码块标记）：
{
  "overall_reasoning": "总体评价理由，每个维度至少引用一条证据中的具体内容作为依据",
  "dimension_scores": {
    "literature_review": {
      "score": 85.5,
      "level": "良好",
      "reasoning": "详细评分理由，必须引用证据中的具体段落或内容作为判据",
      "evidence_quotes": ["证据原文引用段落1", "证据原文引用段落2"]
    },
    "research_design": {
      "score": 80.0,
      "level": "良好",
      "reasoning": "详细评分理由，引用具体证据内容",
      "evidence_quotes": ["证据原文引用段落"]
    },
    "data_analysis": {
      "score": 75.5,
      "level": "中等",
      "reasoning": "详细评分理由，引用具体证据内容",
      "evidence_quotes": ["证据原文引用段落"]
    },
    "critical_thinking": {
      "score": 70.0,
      "level": "中等",
      "reasoning": "详细评分理由，引用具体证据内容",
      "evidence_quotes": ["证据原文引用段落"]
    }
  }
}`
}

func (a *InferenceAgent) buildInferenceUserPrompt(evidences []*models.Evidence) string {
	var sb strings.Builder
	sb.WriteString("## 学生提交的研究证据\n\n")
	for i, ev := range evidences {
		dimName := models.GetDimensionName(models.GetDimensionByKBM(ev.KBMName))
		sb.WriteString(fmt.Sprintf("### 证据%d (ID: %s)\n", i+1, ev.ID))
		sb.WriteString(fmt.Sprintf("- 类型: %s\n", ev.Type))
		sb.WriteString(fmt.Sprintf("- KBM分类: %s → 对应维度: %s\n", ev.KBMName, dimName))
		sb.WriteString(fmt.Sprintf("- 来源: %s\n", ev.SourceType))
		if ev.FileName != "" {
			sb.WriteString(fmt.Sprintf("- 文件名: %s (%s)\n", ev.FileName, ev.FileType))
		}
		sb.WriteString(fmt.Sprintf("\n**内容:**\n%s\n\n", ev.Content))
	}
	sb.WriteString("请基于以上证据，严格按照评分量规进行各维度评估，并务必引用证据原文作为评分依据。")
	return sb.String()
}

type llmInferenceResponse struct {
	OverallReasoning string `json:"overall_reasoning"`
	DimensionScores  map[string]struct {
		Score          float64  `json:"score"`
		Level          string   `json:"level"`
		Reasoning      string   `json:"reasoning"`
		EvidenceQuotes []string `json:"evidence_quotes"`
	} `json:"dimension_scores"`
}

func (a *InferenceAgent) parseLLMInferenceResponse(response string, evidences []*models.Evidence) (*InferenceResult, error) {
	jsonStr := extractJSON(response)
	if jsonStr == "" {
		return nil, fmt.Errorf("未在LLM响应中找到有效JSON")
	}

	var parsed llmInferenceResponse
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, fmt.Errorf("解析LLM响应JSON失败: %w", err)
	}

	dimensionScores := make(map[string]models.DimensionScore)
	var totalWeightedScore float64
	var citations []CitationInfo

	dimensions := models.DefaultDimensions
	weightMap := make(map[string]float64)
	nameMap := make(map[string]string)
	for _, dim := range dimensions {
		weightMap[dim.ID] = dim.Weight
		nameMap[dim.ID] = dim.Name
	}

	// Build evidence ID lookup for citation matching
	evidenceContentMap := make(map[string]string)
	for _, ev := range evidences {
		evidenceContentMap[ev.ID] = ev.Content
	}

	for dimID, score := range parsed.DimensionScores {
		weight := weightMap[dimID]
		if weight == 0 {
			weight = 0.25
		}
		dimName := nameMap[dimID]
		if dimName == "" {
			dimName = dimID
		}

		dimensionScores[dimID] = models.DimensionScore{
			Name:    dimName,
			Score:   score.Score,
			Level:   score.Level,
			Details: score.Reasoning,
		}
		totalWeightedScore += score.Score * weight

		// Build citations from evidence_quotes
		for _, quote := range score.EvidenceQuotes {
			matchedEvID := findEvidenceByExcerpt(quote, evidences)
			citations = append(citations, CitationInfo{
				DimensionID:    dimID,
				EvidenceID:     matchedEvID,
				ExcerptText:    quote,
				RelevanceScore: 0.8,
			})
		}
	}

	overallScore := totalWeightedScore
	overallLevel := models.GetLevelFromScore(overallScore)

	return &InferenceResult{
		OverallScore:    overallScore,
		OverallLevel:    overallLevel,
		DimensionScores: dimensionScores,
		Reasoning:       parsed.OverallReasoning,
		Citations:       citations,
	}, nil
}

// findEvidenceByExcerpt matches a quoted excerpt to the most likely evidence source.
func findEvidenceByExcerpt(excerpt string, evidences []*models.Evidence) string {
	if len(excerpt) < 10 {
		return "unknown"
	}
	bestID := "unknown"
	bestLen := 0
	for _, ev := range evidences {
		if strings.Contains(ev.Content, excerpt[:min(len(excerpt), 50)]) {
			if len(ev.Content) > bestLen {
				bestLen = len(ev.Content)
				bestID = ev.ID
			}
		}
	}
	return bestID
}

// extractJSON extracts a JSON object from text that may contain markdown wrappers.
func extractJSON(s string) string {
	// Try to extract from ```json ... ``` blocks first
	if idx := findJSONInMarkdown(s); idx >= 0 {
		s = s[idx:]
	}
	// Find the outermost JSON object
	start := strings.Index(s, "{")
	if start < 0 {
		return ""
	}
	depth := 0
	for i := start; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}
	return ""
}

func findJSONInMarkdown(s string) int {
	markers := []string{"```json\n", "```json\r\n", "```\n", "```\r\n"}
	for _, m := range markers {
		if idx := strings.Index(s, m); idx >= 0 {
			return idx + len(m)
		}
	}
	return -1
}

func (a *InferenceAgent) inferWithRules(evidences []*models.Evidence) (*InferenceResult, error) {
	dimensionEvidences := a.groupEvidencesByDimension(evidences)

	dimensionScores := make(map[string]models.DimensionScore)
	for dimension, dimEvidences := range dimensionEvidences {
		score, err := a.logicUnit.CalculateDimensionScore(dimEvidences)
		if err != nil {
			return nil, fmt.Errorf("计算维度得分失败: %w", err)
		}

		evidenceIDs := make([]string, len(dimEvidences))
		for i, evidence := range dimEvidences {
			evidenceIDs[i] = evidence.ID
		}

		level := models.GetLevelFromScore(score)
		dimName := models.GetDimensionName(dimension)
		dimensionScores[dimension] = models.DimensionScore{
			Name:        dimName,
			Score:       score,
			Level:       level,
			Details:     fmt.Sprintf("基于%d个证据的规则评估", len(dimEvidences)),
			EvidenceIDs: evidenceIDs,
		}
	}

	overallScore := a.logicUnit.CalculateOverallScore(dimensionScores)
	overallLevel := models.GetLevelFromScore(overallScore)
	reasoning := a.generateRuleReasoning(dimensionScores, overallScore, overallLevel)

	return &InferenceResult{
		OverallScore:    overallScore,
		OverallLevel:    overallLevel,
		DimensionScores: dimensionScores,
		Reasoning:       reasoning,
	}, nil
}

func (a *InferenceAgent) groupEvidencesByDimension(evidences []*models.Evidence) map[string][]*models.Evidence {
	groups := make(map[string][]*models.Evidence)
	for _, evidence := range evidences {
		dimension := models.GetDimensionByKBM(evidence.KBMName)
		groups[dimension] = append(groups[dimension], evidence)
	}
	return groups
}

func (a *InferenceAgent) generateRuleReasoning(dimensionScores map[string]models.DimensionScore, overallScore float64, overallLevel string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("基于收集到的%d个维度的证据，对学生的研究能力进行了综合评估（规则引擎）。", len(dimensionScores)))
	sb.WriteString(fmt.Sprintf("\n总体得分为%.2f，等级为%s。", overallScore, overallLevel))

	var strengths []string
	var weaknesses []string
	for _, score := range dimensionScores {
		if score.Score >= 80 {
			strengths = append(strengths, fmt.Sprintf("%s(%.1f分)", score.Name, score.Score))
		} else if score.Score < 70 {
			weaknesses = append(weaknesses, fmt.Sprintf("%s(%.1f分)", score.Name, score.Score))
		}
	}
	if len(strengths) > 0 {
		sb.WriteString(fmt.Sprintf(" 优势维度：%s。", strings.Join(strengths, "、")))
	}
	if len(weaknesses) > 0 {
		sb.WriteString(fmt.Sprintf(" 待提升维度：%s。", strings.Join(weaknesses, "、")))
	}
	sb.WriteString(" 建议学生在保持优势的同时，针对待提升维度进行重点改进。")
	return sb.String()
}
