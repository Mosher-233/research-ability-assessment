package service

import (
	"context"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/internal/repository/postgres"

	"github.com/google/uuid"
)

type InferenceService struct {
	resultRepo    *postgres.ResultRepo
	evidenceService *EvidenceService
}

func NewInferenceService(resultRepo *postgres.ResultRepo, evidenceService *EvidenceService) *InferenceService {
	return &InferenceService{
		resultRepo:    resultRepo,
		evidenceService: evidenceService,
	}
}

func (s *InferenceService) CreateInferenceResult(ctx context.Context, result *models.InferenceResult) error {
	result.ID = uuid.New().String()
	return s.resultRepo.CreateInferenceResult(ctx, result)
}

func (s *InferenceService) GetInferenceResultByID(ctx context.Context, id string) (*models.InferenceResult, error) {
	return s.resultRepo.GetInferenceResultByID(ctx, id)
}

func (s *InferenceService) GetInferenceResultsByTaskID(ctx context.Context, taskID string) ([]models.InferenceResult, error) {
	return s.resultRepo.GetInferenceResultsByTaskID(ctx, taskID)
}

func (s *InferenceService) GetInferenceResultByStudentAndTask(ctx context.Context, studentID string, taskID string) (*models.InferenceResult, error) {
	return s.resultRepo.GetInferenceResultByStudentAndTask(ctx, studentID, taskID)
}

func (s *InferenceService) GetEvidencesByStudentAndTask(ctx context.Context, studentID string, taskID string) ([]models.Evidence, error) {
	return s.evidenceService.GetEvidencesByStudentAndTask(ctx, studentID, taskID)
}