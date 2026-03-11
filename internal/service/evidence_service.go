package service

import (
	"context"
	"research-ability-assessment/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EvidenceService struct {
	db *gorm.DB
}

func NewEvidenceService(db *gorm.DB) *EvidenceService {
	return &EvidenceService{db: db}
}

func (s *EvidenceService) CreateEvidence(ctx context.Context, evidence *models.Evidence) error {
	evidence.ID = uuid.New().String()
	return s.db.WithContext(ctx).Create(evidence).Error
}

func (s *EvidenceService) GetEvidenceByID(ctx context.Context, id string) (*models.Evidence, error) {
	var evidence models.Evidence
	if err := s.db.WithContext(ctx).First(&evidence, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &evidence, nil
}

func (s *EvidenceService) GetEvidencesByStudentTaskID(ctx context.Context, studentTaskID string) ([]models.Evidence, error) {
	var evidences []models.Evidence
	if err := s.db.WithContext(ctx).Where("student_task_id = ?", studentTaskID).Find(&evidences).Error; err != nil {
		return nil, err
	}
	return evidences, nil
}

func (s *EvidenceService) GetEvidencesByStudentAndTask(ctx context.Context, studentID string, taskID string) ([]models.Evidence, error) {
	var evidences []models.Evidence
	if err := s.db.WithContext(ctx).Joins("JOIN student_tasks ON evidences.student_task_id = student_tasks.id").Where("student_tasks.student_id = ? AND student_tasks.task_id = ?", studentID, taskID).Find(&evidences).Error; err != nil {
		return nil, err
	}
	return evidences, nil
}