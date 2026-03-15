package postgres

import (
	"context"
	"log"
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
	log.Printf("ResultRepo: 创建推理结果, StudentID=%s, TaskID=%s", result.StudentID, result.TaskID)
	return r.db.WithContext(ctx).Create(result).Error
}

func (r *ResultRepo) GetInferenceResultByID(ctx context.Context, id string) (*models.InferenceResult, error) {
	var result models.InferenceResult
	if err := r.db.WithContext(ctx).First(&result, "id = ?", id).Error; err != nil {
		log.Printf("ResultRepo: 未找到推理结果, ID=%s, err=%v", id, err)
		return nil, err
	}
	return &result, nil
}

func (r *ResultRepo) GetInferenceResultsByTaskID(ctx context.Context, taskID string) ([]models.InferenceResult, error) {
	var results []models.InferenceResult
	if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Find(&results).Error; err != nil {
		log.Printf("ResultRepo: 获取任务推理结果失败, TaskID=%s, err=%v", taskID, err)
		return nil, err
	}
	log.Printf("ResultRepo: 获取任务推理结果成功, TaskID=%s, 数量=%d", taskID, len(results))
	return results, nil
}

func (r *ResultRepo) GetInferenceResultByStudentAndTask(ctx context.Context, studentID string, taskID string) (*models.InferenceResult, error) {
	var result models.InferenceResult
	if err := r.db.WithContext(ctx).First(&result, "student_id = ? AND task_id = ?", studentID, taskID).Error; err != nil {
		log.Printf("ResultRepo: 未找到学生任务推理结果, StudentID=%s, TaskID=%s, err=%v", studentID, taskID, err)
		return nil, err
	}
	return &result, nil
}

func (r *ResultRepo) GetAllInferenceResults(ctx context.Context) ([]models.InferenceResult, error) {
	var results []models.InferenceResult
	if err := r.db.WithContext(ctx).Find(&results).Error; err != nil {
		log.Printf("ResultRepo: 获取所有推理结果失败, err=%v", err)
		return nil, err
	}
	log.Printf("ResultRepo: 获取所有推理结果成功, 数量=%d", len(results))
	return results, nil
}

func (r *ResultRepo) GetInferenceResultsByStudentID(ctx context.Context, studentID string) ([]models.InferenceResult, error) {
	var results []models.InferenceResult
	if err := r.db.WithContext(ctx).Where("student_id = ?", studentID).Find(&results).Error; err != nil {
		log.Printf("ResultRepo: 获取学生推理结果失败, StudentID=%s, err=%v", studentID, err)
		return nil, err
	}
	log.Printf("ResultRepo: 获取学生推理结果成功, StudentID=%s, 数量=%d", studentID, len(results))
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