package llm

import (
	"strconv"
	"strings"
)

// ParsedFeedback 解析后的反馈信息
type ParsedFeedback struct {
	KBMLevel    int
	Strengths   string
	Weaknesses  string
	Suggestions string
	Feedback    string
}

// ParseInferenceResponse 解析LLM的能力推理响应（暂未使用）
// func ParseInferenceResponse(response string) (*models.InferenceResult, error) {
// 	// 这里可以实现对LLM响应的解析逻辑
// 	// 例如从文本响应中提取评估结果
//
// 	// 简单的示例实现
// 	result := &models.InferenceResult{
// 		OverallScore:    0.75,
// 		OverallLevel:    "良好",
// 		DimensionScores: make(map[string]models.DimensionScore),
// 		Reasoning:       response,
// 	}
//
// 	// 添加示例维度得分
// 	result.DimensionScores["literature"] = models.DimensionScore{
// 		Name:        "literature",
// 		Score:       0.8,
// 		Level:       "优秀",
// 		Details:     "文献检索策略合理，综述质量高",
// 		EvidenceIDs: []string{},
// 	}
//
// 	result.DimensionScores["experiment_design"] = models.DimensionScore{
// 		Name:        "experiment_design",
// 		Score:       0.7,
// 		Level:       "良好",
// 		Details:     "实验方案合理，变量控制较好",
// 		EvidenceIDs: []string{},
// 	}
//
// 	result.DimensionScores["data_processing"] = models.DimensionScore{
// 		Name:        "data_processing",
// 		Score:       0.75,
// 		Level:       "良好",
// 		Details:     "数据分析方法选择恰当，结果解释准确",
// 		EvidenceIDs: []string{},
// 	}
//
// 	result.DimensionScores["innovation"] = models.DimensionScore{
// 		Name:        "innovation",
// 		Score:       0.65,
// 		Level:       "良好",
// 		Details:     "问题提出有一定新颖性，解决方案有一定原创性",
// 		EvidenceIDs: []string{},
// 	}
//
// 	return result, nil
// }

// ParseFeedbackResponse 解析LLM的反馈响应
func ParseFeedbackResponse(response string) (*ParsedFeedback, error) {
	result := &ParsedFeedback{
		KBMLevel:    3,
		Strengths:   "提交的证据内容较为完整",
		Weaknesses:  "可进一步提升深度",
		Suggestions: "继续深入研究",
		Feedback:    response,
	}

	lines := strings.Split(response, "\n")

	var currentSection string
	var strengthsBuilder strings.Builder
	var weaknessesBuilder strings.Builder
	var suggestionsBuilder strings.Builder
	var overallBuilder strings.Builder
	foundLevel := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 检查是否是新的章节标题
		lowerLine := strings.ToLower(line)

		if !foundLevel && (strings.Contains(lowerLine, "级别") || strings.Contains(lowerLine, "level:")) {
			// 只解析第一个出现的级别
			for _, r := range line {
				if r >= '1' && r <= '5' {
					if level, err := strconv.Atoi(string(r)); err == nil {
						result.KBMLevel = level
						foundLevel = true
					}
					break
				}
			}
			currentSection = ""
		} else if strings.Contains(lowerLine, "优点") || strings.HasPrefix(lowerLine, "strengths") {
			currentSection = "strengths"
		} else if strings.Contains(lowerLine, "不足") || strings.HasPrefix(lowerLine, "weaknesses") {
			currentSection = "weaknesses"
		} else if strings.Contains(lowerLine, "建议") || strings.HasPrefix(lowerLine, "suggestions") {
			currentSection = "suggestions"
		} else if strings.Contains(lowerLine, "总体评价") || strings.HasPrefix(lowerLine, "overall") {
			currentSection = "overall"
		} else if line != "" {
			// 不是章节标题，添加到当前章节
			switch currentSection {
			case "strengths":
				if strengthsBuilder.Len() > 0 {
					strengthsBuilder.WriteString("\n")
				}
				strengthsBuilder.WriteString(line)
			case "weaknesses":
				if weaknessesBuilder.Len() > 0 {
					weaknessesBuilder.WriteString("\n")
				}
				weaknessesBuilder.WriteString(line)
			case "suggestions":
				if suggestionsBuilder.Len() > 0 {
					suggestionsBuilder.WriteString("\n")
				}
				suggestionsBuilder.WriteString(line)
			case "overall":
				if overallBuilder.Len() > 0 {
					overallBuilder.WriteString("\n")
				}
				overallBuilder.WriteString(line)
			}
		}
	}

	// 更新结果
	if s := strings.TrimSpace(strengthsBuilder.String()); s != "" {
		result.Strengths = s
	}
	if w := strings.TrimSpace(weaknessesBuilder.String()); w != "" {
		result.Weaknesses = w
	}
	if s := strings.TrimSpace(suggestionsBuilder.String()); s != "" {
		result.Suggestions = s
	}

	return result, nil
}
