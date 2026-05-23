package agent

import (
	"math"
	"github.com/Mosher-233/research-ability-assessment/internal/models"
)

type LogicUnit struct{}

func NewLogicUnit() *LogicUnit {
	return &LogicUnit{}
}

func (l *LogicUnit) EvaluateEvidence(evidence *models.Evidence) float64 {
	if evidence.KBMLevel > 0 {
		score := float64(evidence.KBMLevel) * 20
		return math.Min(score, 100)
	}
	return 60
}

func (l *LogicUnit) CalculateDimensionScore(evidences []*models.Evidence) (float64, error) {
	if len(evidences) == 0 {
		return 50, nil
	}

	totalScore := 0.0
	for _, evidence := range evidences {
		totalScore += l.EvaluateEvidence(evidence)
	}

	averageScore := totalScore / float64(len(evidences))
	return math.Min(averageScore, 100), nil
}

func (l *LogicUnit) CalculateOverallScore(dimensionScores map[string]models.DimensionScore) float64 {
	totalWeightedScore := 0.0
	totalWeight := 0.0

	for dimID, score := range dimensionScores {
		weight := models.GetDimensionWeight(dimID)
		totalWeightedScore += score.Score * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0
	}

	return math.Round(totalWeightedScore/totalWeight*100) / 100
}
