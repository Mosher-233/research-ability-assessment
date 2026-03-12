package postgres

import (
	"context"
	"research-ability-assessment/internal/models"

	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) CreateTask(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *TaskRepo) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	var task models.Task
	if err := r.db.WithContext(ctx).First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepo) GetTasksByTeacherID(ctx context.Context, teacherID string) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.WithContext(ctx).Where("teacher_id = ?", teacherID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepo) UpdateTaskStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&models.Task{}).Where("id = ?", id).Update("status", status).Error
}

func (r *TaskRepo) CreateStudentTask(ctx context.Context, studentTask *models.StudentTask) error {
	return r.db.WithContext(ctx).Create(studentTask).Error
}

func (r *TaskRepo) GetStudentTaskByID(ctx context.Context, id string) (*models.StudentTask, error) {
	var studentTask models.StudentTask
	if err := r.db.WithContext(ctx).First(&studentTask, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &studentTask, nil
}

func (r *TaskRepo) GetStudentTasksByTaskID(ctx context.Context, taskID string) ([]models.StudentTask, error) {
	var studentTasks []models.StudentTask
	if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Find(&studentTasks).Error; err != nil {
		return nil, err
	}
	return studentTasks, nil
}

func (r *TaskRepo) UpdateStudentTaskStatus(ctx context.Context, id string, status string, progress int) error {
	return r.db.WithContext(ctx).Model(&models.StudentTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":   status,
		"progress": progress,
	}).Error
}

func (r *TaskRepo) GetAssignedTasks(ctx context.Context, studentID string) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.WithContext(ctx).Joins("JOIN student_tasks ON tasks.id = student_tasks.task_id").Where("student_tasks.student_id = ?", studentID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}