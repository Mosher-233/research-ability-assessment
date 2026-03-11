package agent

import (
	"context"
	"fmt"
	"research-ability-assessment/internal/models"
)

type FeedbackAgent struct {
	// 可以添加LLM客户端或其他反馈生成组件
}

func NewFeedbackAgent() *FeedbackAgent {
	return &FeedbackAgent{}
}

func (a *FeedbackAgent) GenerateFeedback(ctx context.Context, result *InferenceResult) string {
	// 生成反馈
	feedback := "# 研究能力评估反馈\n\n"
	feedback += fmt.Sprintf("## 总体评估\n您的研究能力总体得分为%.2f，等级为%s。\n\n", result.OverallScore, result.OverallLevel)
	
	feedback += "## 各维度评估\n"
	for dimension, score := range result.DimensionScores {
		feedback += fmt.Sprintf("### %s\n- 得分：%.2f\n- 等级：%s\n- 详情：%s\n\n", dimension, score.Score, score.Level, score.Details)
	}
	
	feedback += "## 改进建议\n"
	feedback += a.generateImprovementSuggestions(result.DimensionScores)
	
	return feedback
}

func (a *FeedbackAgent) generateImprovementSuggestions(dimensionScores map[string]models.DimensionScore) string {
	suggestions := ""
	
	for dimension, score := range dimensionScores {
		if score.Score < 0.6 {
			switch dimension {
			case "literature":
				suggestions += "- **文献能力**：建议加强文献检索策略的学习，提高文献综述的质量和批判性分析能力。\n"
			case "experiment_design":
				suggestions += "- **实验设计**：建议优化实验方案的合理性，加强变量控制，提高实验实施的质量。\n"
			case "data_processing":
				suggestions += "- **数据处理**：建议学习更多数据分析方法，提高结果解释的准确性。\n"
			case "innovation":
				suggestions += "- **创新能力**：建议培养问题提出的新颖性，提高解决方案的原创性。\n"
			default:
				suggestions += fmt.Sprintf("- **%s**：建议加强相关知识和技能的学习。\n", dimension)
			}
		}
	}
	
	if suggestions == "" {
		suggestions = "您的各维度表现良好，继续保持！\n"
	}
	
	return suggestions
}