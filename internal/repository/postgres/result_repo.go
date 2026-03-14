package postgres

import (
	"context"
	"research-ability-assessment/internal/models"

	"gorm.io/gorm"
)

type ResultRepo struct {
	db *gorm.DB
}

func NewResultRepo(db *gorm.DB) *ResultRepo {
	return &ResultRepo{db: db}
}

func (r *ResultRepo) CreateInferenceResult(ctx context.Context, result *models.InferenceResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}

func (r *ResultRepo) GetInferenceResultByID(ctx context.Context, id string) (*models.InferenceResult, error) {
	var result models.InferenceResult
	if err := r.db.WithContext(ctx).First(&result, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ResultRepo) GetInferenceResultsByTaskID(ctx context.Context, taskID string) ([]models.InferenceResult, error) {
	var results []models.InferenceResult
	if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Preload("Student").Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *ResultRepo) GetInferenceResultByStudentAndTask(ctx context.Context, studentID string, taskID string) (*models.InferenceResult, error) {
	var result models.InferenceResult
	if err := r.db.WithContext(ctx).First(&result, "student_id = ? AND task_id = ?", studentID, taskID).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ResultRepo) GetAllInferenceResults(ctx context.Context) ([]models.InferenceResult, error) {
	var results []models.InferenceResult
	if err := r.db.WithContext(ctx).Preload("Student").Preload("Task").Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *ResultRepo) GetInferenceResultsByStudentID(ctx context.Context, studentID string) ([]models.InferenceResult, error) {
	var results []models.InferenceResult
	if err := r.db.WithContext(ctx).Where("student_id = ?", studentID).Preload("Task").Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *ResultRepo) CreateReport(ctx context.Context, report *models.Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *ResultRepo) GetReportByID(ctx context.Context, id string) (*models.Report, error) {
	var report models.Report
	if err := r.db.WithContext(ctx).First(&report, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ResultRepo) GetReportByStudentAndTask(ctx context.Context, studentID string, taskID string) (*models.Report, error) {
	var report models.Report
	if err := r.db.WithContext(ctx).First(&report, "student_id = ? AND task_id = ?", studentID, taskID).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ResultRepo) GetReportsByTaskID(ctx context.Context, taskID string) ([]models.Report, error) {
	var reports []models.Report
	if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *ResultRepo) GetReportsByStudentID(ctx context.Context, studentID string) ([]models.Report, error) {
	var reports []models.Report
	if err := r.db.WithContext(ctx).Where("student_id = ?", studentID).Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *ResultRepo) GetAllReports(ctx context.Context) ([]models.Report, error) {
	var reports []models.Report
	if err := r.db.WithContext(ctx).Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}