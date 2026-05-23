package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Mosher-233/research-ability-assessment/internal/llm"
	"github.com/Mosher-233/research-ability-assessment/internal/models"
)

type FeedbackAgent struct {
	llmClient *llm.Client
}

func NewFeedbackAgent(llmClient *llm.Client) *FeedbackAgent {
	return &FeedbackAgent{llmClient: llmClient}
}

type FeedbackResult struct {
	OverallFeedback    string                        `json:"overall_feedback"`
	DimensionFeedbacks map[string]string             `json:"dimension_feedbacks"`
	Suggestions        []models.ImprovementSuggestion `json:"suggestions"`
}

func (a *FeedbackAgent) GenerateFeedback(ctx context.Context, result *InferenceResult) (*FeedbackResult, error) {
	if a.llmClient != nil {
		llmFeedback, err := a.generateLLMFeedback(ctx, result)
		if err == nil {
			return llmFeedback, nil
		}
		log.Printf("FeedbackAgent: LLM生成反馈失败，使用规则引擎: %v", err)
	}

	return a.generateRuleFeedback(result), nil
}

func (a *FeedbackAgent) generateLLMFeedback(ctx context.Context, result *InferenceResult) (*FeedbackResult, error) {
	dimensionJSON, _ := json.Marshal(result.DimensionScores)

	systemPrompt := `你是一个专业的大学生研究能力评价专家，负责基于评估结果生成个性化改进建议。

你需要基于以下四个评估维度生成反馈：
1. 文献综述 (literature_review)：文献检索与综述撰写能力
2. 研究设计 (research_design)：研究方案设计与实验规划能力
3. 数据分析 (data_analysis)：数据处理与统计分析能力
4. 批判性思维 (critical_thinking)：批判性思考与创新思维能力

## 要求
- 对得分 >= 80 的维度给予肯定和保持建议
- 对得分 < 70 的维度给出具体的改进路径和可执行行动项
- 建议要具体、可操作，不要空泛
- 每个待提升维度至少给出3个行动项
- 建议必须与学生提交的具体证据内容相关联，而非泛泛而谈
- 引用证据中的具体内容来说明为什么某个维度得分低以及如何改进

## 输出格式（严格JSON，不要包含markdown代码块标记）
{
  "overall_feedback": "总体反馈描述，结合学生具体证据内容",
  "dimension_feedbacks": {
    "literature_review": "针对该学生文献综述能力的具体反馈",
    "research_design": "针对该学生研究设计能力的具体反馈",
    "data_analysis": "针对该学生数据分析能力的具体反馈",
    "critical_thinking": "针对该学生批判性思维能力的具体反馈"
  },
  "suggestions": [
    {
      "dimension": "dimension_id",
      "dimension_name": "维度中文名",
      "current_score": 75.0,
      "target_score": 85.0,
      "suggestion": "基于学生具体证据问题的改进建议",
      "action_items": ["具体行动1", "具体行动2", "具体行动3"],
      "priority": 1
    }
  ]
}`

	userPrompt := fmt.Sprintf("评估结果：\n总体得分：%.1f，等级：%s\n各维度得分：%s\n综合推理：%s\n\n请基于以上结果生成个性化的改进反馈和建议。务必引用评估推理中提到的具体证据问题。",
		result.OverallScore, result.OverallLevel, string(dimensionJSON), result.Reasoning)

	messages := []llm.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	response, err := a.llmClient.Chat(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("LLM调用失败: %w", err)
	}

	return a.parseLLMFeedbackResponse(response)
}

func (a *FeedbackAgent) parseLLMFeedbackResponse(response string) (*FeedbackResult, error) {
	jsonStr := extractJSON(response)
	if jsonStr == "" {
		return nil, fmt.Errorf("未在LLM响应中找到JSON")
	}

	var feedbackResult FeedbackResult
	if err := json.Unmarshal([]byte(jsonStr), &feedbackResult); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}

	if feedbackResult.DimensionFeedbacks == nil {
		feedbackResult.DimensionFeedbacks = make(map[string]string)
	}

	return &feedbackResult, nil
}

func (a *FeedbackAgent) generateRuleFeedback(result *InferenceResult) *FeedbackResult {
	feedbackResult := &FeedbackResult{
		DimensionFeedbacks: make(map[string]string),
	}

	feedbackResult.OverallFeedback = fmt.Sprintf("您的研究能力总体得分为%.1f分，等级为%s。", result.OverallScore, result.OverallLevel)

	if result.OverallScore >= 80 {
		feedbackResult.OverallFeedback += "整体表现优秀，研究能力均衡发展。建议在保持现有水平的基础上，探索更具挑战性的研究方向。"
	} else if result.OverallScore >= 70 {
		feedbackResult.OverallFeedback += "整体表现良好，具备基本的研究能力，部分维度有进一步提升空间。"
	} else if result.OverallScore >= 60 {
		feedbackResult.OverallFeedback += "整体表现中等，建议系统性地加强研究能力训练，重点提升薄弱维度。"
	} else {
		feedbackResult.OverallFeedback += "整体表现有待提高，建议从基础研究方法入手，逐步建立完整的研究能力体系。"
	}

	// Score-specific dimension feedback (varies by score range within each band)
	dimensionFeedbackMap := map[string]struct {
		excellent string
		good      string
		average   string
		weak      string
	}{
		models.DimLiteratureReview: {
			excellent: "文献综述能力出色，能够系统检索相关文献并进行批判性分析。建议尝试撰写学术综述文章。",
			good:      "文献综述能力良好，文献检索覆盖面广，综述组织较为清晰。可加强批判性分析深度。",
			average:   "文献综述能力一般，建议加强文献检索策略的学习，提高综述的组织逻辑性和批判分析能力。",
			weak:      "文献综述能力有待提升，建议从基础文献检索方法学起，逐步提高文献筛选和批判性分析能力。",
		},
		models.DimResearchDesign: {
			excellent: "研究设计能力出色，实验方案科学合理，变量控制严密。建议尝试设计更复杂的多变量研究。",
			good:      "研究设计能力良好，方案设计合理，变量控制基本到位。可加强实验实施的质量控制。",
			average:   "研究设计能力一般，建议加强对实验设计方法论的学习，提升方案的科学性和可行性。",
			weak:      "研究设计能力有待提升，建议系统学习研究设计方法论，加强变量控制和实验规范。",
		},
		models.DimDataAnalysis: {
			excellent: "数据分析能力出色，方法选择恰当，结果解释准确深入。建议探索更高级的统计分析方法。",
			good:      "数据分析能力良好，能够选择合适的分析方法并进行合理解释。可加强数据可视化和深层分析。",
			average:   "数据分析能力一般，建议学习更多数据分析方法，加强统计软件的使用能力和结果解释训练。",
			weak:      "数据分析能力有待提升，建议从基础统计学和数据分析方法开始学习，系统提升分析能力。",
		},
		models.DimCriticalThinking: {
			excellent: "批判性思维能力出色，能够独立提出新颖问题并从多角度进行分析。建议指导低年级同学。",
			good:      "批判性思维能力良好，具备独立思考和问题分析能力。可加强创新思维的训练。",
			average:   "批判性思维能力一般，建议多进行批判性阅读训练，培养从多角度分析问题的习惯。",
			weak:      "批判性思维能力有待提升，建议从批判性思维基本训练开始，培养独立思考和创新能力。",
		},
	}

	actionItemMap := map[string][]string{
		models.DimLiteratureReview: {
			"每周阅读3-5篇核心期刊论文,并做批判性笔记",
			"使用文献管理工具建立个人文献数据库",
			"练习撰写文献综述段落，注意逻辑衔接",
		},
		models.DimResearchDesign: {
			"分析5个优秀研究案例的实验设计",
			"在设计新实验前列出所有潜在变量和控制方法",
			"定期与导师讨论研究方案的科学性",
		},
		models.DimDataAnalysis: {
			"完成3个数据分析实战项目(可参考Kaggle)",
			"学习使用专业统计软件(如SPSS、R或Python)",
			"练习撰写数据分析报告，注重结果的合理解释",
		},
		models.DimCriticalThinking: {
			"每日阅读一篇学术评论文章，练习批判性分析",
			"参加学术讨论或辩论活动，锻炼多角度思维",
			"写反思日记，记录对所学知识的批判性思考",
		},
	}

	priority := 1
	for dimID, score := range result.DimensionScores {
		dimName := models.GetDimensionName(dimID)

		var fb string
		specificHint := fmt.Sprintf("(当前得分%.1f)", score.Score)
		switch {
		case score.Score >= 90:
			fb = dimensionFeedbackMap[dimID].excellent + " " + specificHint
		case score.Score >= 75:
			fb = dimensionFeedbackMap[dimID].good + " " + specificHint
		case score.Score >= 70:
			fb = dimensionFeedbackMap[dimID].good + " " + specificHint + " 建议针对薄弱环节专项提升。"
		case score.Score >= 60:
			fb = dimensionFeedbackMap[dimID].average + " " + specificHint
		default:
			fb = dimensionFeedbackMap[dimID].weak + " " + specificHint
		}
		feedbackResult.DimensionFeedbacks[dimID] = fb

		if score.Score < 75 {
			targetScore := score.Score + 15
			if targetScore > 100 {
				targetScore = 100
			}
			suggestion := models.ImprovementSuggestion{
				ID:            dimID,
				Dimension:     dimID,
				DimensionName: dimName,
				CurrentScore:  score.Score,
				TargetScore:   targetScore,
				Suggestion:    fb,
				ActionItems:   actionItemMap[dimID],
				Priority:      priority,
			}
			feedbackResult.Suggestions = append(feedbackResult.Suggestions, suggestion)
			priority++
		}
	}

	if len(feedbackResult.Suggestions) == 0 {
		feedbackResult.Suggestions = []models.ImprovementSuggestion{
			{
				ID:            "general",
				Dimension:     "general",
				DimensionName: "综合提升",
				CurrentScore:  result.OverallScore,
				TargetScore:   result.OverallScore + 5,
				Suggestion:    "各维度表现均衡优异，建议挑战更高层次的研究任务。",
				ActionItems: []string{
					"参与导师的科研项目",
					"尝试独立设计并完成一个研究课题",
					"将研究成果整理成论文投稿",
				},
				Priority: 1,
			},
		}
	}

	return feedbackResult
}
