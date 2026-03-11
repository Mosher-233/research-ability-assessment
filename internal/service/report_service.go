package service

import (
	"context"
	"research-ability-assessment/internal/models"
)

type ReportService struct {
	inferenceService *InferenceService
}

func NewReportService(inferenceService *InferenceService) *ReportService {
	return &ReportService{
		inferenceService: inferenceService,
	}
}

func (s *ReportService) GenerateStudentReport(ctx context.Context, studentID string, taskID string) (*models.InferenceResult, error) {
	return s.inferenceService.GetInferenceResultByStudentAndTask(ctx, studentID, taskID)
}

func (s *ReportService) GenerateTaskReport(ctx context.Context, taskID string) ([]models.InferenceResult, error) {
	return s.inferenceService.GetInferenceResultsByTaskID(ctx, taskID)
}