package agent

import (
	"context"
	"fmt"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/internal/service"
)

type InferenceAgent struct {
	evidanceAgent    *EvidenceAgent
	logicUnit        *LogicUnit
	inferenceService *service.InferenceService
}

type InferenceResult struct {
	OverallScore    float64                       `json:"overall_score"`
	OverallLevel    string                        `json:"overall_level"`
	DimensionScores map[string]models.DimensionScore `json:"dimension_scores"`
	Reasoning       string                        `json:"reasoning"`
}

func NewInferenceAgent(
	evidanceAgent *EvidenceAgent,
	logicUnit *LogicUnit,
	inferenceService *service.InferenceService,
) *InferenceAgent {
	return &InferenceAgent{
		evidanceAgent:    evidanceAgent,
		logicUnit:        logicUnit,
		inferenceService: inferenceService,
	}
}

func (a *InferenceAgent) InferAbility(ctx context.Context, studentID string, taskID string) (*InferenceResult, error) {
	// 收集证据
	evidences, err := a.evidanceAgent.CollectEvidence(ctx, studentID, taskID)
	if err != nil {
		return nil, fmt.Errorf("收集证据失败: %w", err)
	}

	// 按维度分组证据
	dimensionEvidences := a.groupEvidencesByDimension(evidences)

	// 计算各维度得分
	dimensionScores := make(map[string]models.DimensionScore)
	for dimension, dimEvidences := range dimensionEvidences {
		score, err := a.logicUnit.CalculateDimensionScore(dimEvidences)
		if err != nil {
			return nil, fmt.Errorf("计算维度得分失败: %w", err)
		}

		// 收集证据ID
		evidenceIDs := make([]string, len(dimEvidences))
		for i, evidence := range dimEvidences {
			evidenceIDs[i] = evidence.ID
		}

		level := a.logicUnit.DetermineLevel(score)
		dimensionScores[dimension] = models.DimensionScore{
			Name:        dimension,
			Score:       score,
			Level:       level,
			Details:     fmt.Sprintf("基于%d个证据的评估", len(dimEvidences)),
			EvidenceIDs: evidenceIDs,
		}
	}

	// 计算总体得分
	overallScore := a.logicUnit.CalculateOverallScore(dimensionScores)
	overallLevel := a.logicUnit.DetermineLevel(overallScore)

	// 生成推理理由
	reasoning := a.generateReasoning(dimensionScores, overallScore, overallLevel)

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
		dimension := a.mapKBMToDimension(evidence.KBMName)
		groups[dimension] = append(groups[dimension], evidence)
	}
	
	return groups
}

func (a *InferenceAgent) mapKBMToDimension(kbmName string) string {
	// KBM到维度的映射
	mapping := map[string]string{
		"文献检索策略":       "literature",
		"文献综述质量":       "literature",
		"文献批判性分析":     "literature",
		"实验方案合理性":     "experiment_design",
		"变量控制":         "experiment_design",
		"实验实施质量":       "experiment_design",
		"数据分析方法选择":   "data_processing",
		"结果解释准确性":     "data_processing",
		"问题提出新颖性":     "innovation",
		"解决方案原创性":     "innovation",
	}
	
	if dimension, ok := mapping[kbmName]; ok {
		return dimension
	}
	
	return "innovation" // 默认映射到创新维度
}

func (a *InferenceAgent) generateReasoning(dimensionScores map[string]models.DimensionScore, overallScore float64, overallLevel string) string {
	// 生成推理理由
	reasoning := fmt.Sprintf("基于收集到的证据，对学生的研究能力进行了综合评估。")
	reasoning += fmt.Sprintf("\n总体得分为%.2f，等级为%s。", overallScore, overallLevel)
	
	reasoning += "\n各维度评估结果："
	for dimension, score := range dimensionScores {
		reasoning += fmt.Sprintf("\n- %s: 得分%.2f，等级%s", dimension, score.Score, score.Level)
	}
	
	return reasoning
}