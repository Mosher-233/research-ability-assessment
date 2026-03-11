package agent

import (
	"research-ability-assessment/internal/models"
)

type LogicUnit struct {
	// 可以添加规则引擎或其他逻辑处理组件
}

func NewLogicUnit() *LogicUnit {
	return &LogicUnit{}
}

func (l *LogicUnit) EvaluateEvidence(evidence *models.Evidence) (float64, error) {
	// 这里可以实现证据评估的逻辑
	// 例如根据证据类型、内容和KBM级别计算得分
	
	// 简单的示例实现
	score := float64(evidence.KBMLevel) * 0.2
	return score, nil
}

func (l *LogicUnit) CalculateDimensionScore(evidences []*models.Evidence) (float64, error) {
	// 这里可以实现维度得分计算的逻辑
	// 例如对该维度下的所有证据得分进行加权平均
	
	if len(evidences) == 0 {
		return 0, nil
	}
	
	totalScore := 0.0
	for _, evidence := range evidences {
		score, err := l.EvaluateEvidence(evidence)
		if err != nil {
			return 0, err
		}
		totalScore += score
	}
	
	averageScore := totalScore / float64(len(evidences))
	return averageScore, nil
}

func (l *LogicUnit) CalculateOverallScore(dimensionScores map[string]models.DimensionScore) float64 {
	// 这里可以实现总体得分计算的逻辑
	// 例如对所有维度得分进行加权平均
	
	totalScore := 0.0
	count := 0
	
	for _, score := range dimensionScores {
		totalScore += score.Score
		count++
	}
	
	if count == 0 {
		return 0
	}
	
	averageScore := totalScore / float64(count)
	return averageScore
}

func (l *LogicUnit) DetermineLevel(score float64) string {
	// 根据得分确定等级
	switch {
	case score >= 0.8:
		return "优秀"
	case score >= 0.6:
		return "良好"
	case score >= 0.4:
		return "中等"
	default:
		return "待提高"
	}
}