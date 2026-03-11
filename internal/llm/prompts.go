package llm

import (
	"fmt"
	"research-ability-assessment/internal/models"
)

// GetInferencePrompt 获取能力推理的提示词
func GetInferencePrompt(evidences []*models.Evidence) []Message {
	prompt := "你是一个研究能力评价专家，负责根据学生提供的证据评估其研究能力。"
	prompt += "请根据以下证据，对学生的研究能力进行综合评估，包括文献能力、实验设计能力、数据处理能力和创新能力四个维度。"
	prompt += "每个维度的评估结果应包括得分（0-1）、等级（优秀、良好、中等、待提高）和详细说明。"
	prompt += "最后给出总体得分和等级，并提供改进建议。"
	prompt += "\n\n证据：\n"

	for i, evidence := range evidences {
		prompt += fmt.Sprintf("%d. 类型：%s，内容：%s，KBM名称：%s，KBM级别：%d\n", i+1, evidence.Type, evidence.Content, evidence.KBMName, evidence.KBMLevel)
	}

	return []Message{
		{
			Role:    "system",
			Content: "你是一个研究能力评价专家，负责根据学生提供的证据评估其研究能力。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}
}

// GetFeedbackPrompt 获取反馈生成的提示词
func GetFeedbackPrompt(result *models.InferenceResult) []Message {
	prompt := "你是一个教育顾问，负责根据研究能力评估结果为学生提供详细的反馈和改进建议。"
	prompt += "请根据以下评估结果，生成一份详细的反馈报告，包括总体评估、各维度评估和具体的改进建议。"
	prompt += "反馈报告应该友好、专业，并且具有建设性。"
	prompt += "\n\n评估结果：\n"
	prompt += fmt.Sprintf("总体得分：%.2f，等级：%s\n", result.OverallScore, result.OverallLevel)
	prompt += "各维度得分：\n"

	for dimension, score := range result.DimensionScores {
		prompt += fmt.Sprintf("- %s：得分 %.2f，等级 %s\n", dimension, score.Score, score.Level)
	}

	return []Message{
		{
			Role:    "system",
			Content: "你是一个教育顾问，负责根据研究能力评估结果为学生提供详细的反馈和改进建议。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}
}