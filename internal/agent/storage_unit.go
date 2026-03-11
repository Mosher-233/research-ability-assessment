package agent

import (
	"context"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/internal/repository/neo4j"
	"research-ability-assessment/internal/repository/postgres"
)

type StorageUnit struct {
	resultRepo *postgres.ResultRepo
	graphRepo  *neo4j.GraphRepo
}

func NewStorageUnit(resultRepo *postgres.ResultRepo, graphRepo *neo4j.GraphRepo) *StorageUnit {
	return &StorageUnit{
		resultRepo: resultRepo,
		graphRepo:  graphRepo,
	}
}

func (s *StorageUnit) StoreInferenceResult(ctx context.Context, result *models.InferenceResult) error {
	return s.resultRepo.CreateInferenceResult(ctx, result)
}

func (s *StorageUnit) UpdateKnowledgeGraph(ctx context.Context, studentID string, dimension string, score float64) error {
	// 确保维度节点存在
	if err := s.graphRepo.CreateDimensionNode(dimension); err != nil {
		return err
	}
	
	// 更新学生和维度之间的关系
	return s.graphRepo.UpdateKnowledgeGraph(studentID, dimension, score)
}

func (s *StorageUnit) GetStudentScores(ctx context.Context, studentID string) (map[string]float64, error) {
	return s.graphRepo.GetStudentScores(studentID)
}
