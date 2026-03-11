package service

import (
	"context"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/internal/repository/postgres"

	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo *postgres.TaskRepo
	userRepo *postgres.UserRepo
}

func NewTaskService(taskRepo *postgres.TaskRepo, userRepo *postgres.UserRepo) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task) error {
	task.ID = uuid.New().String()
	return s.taskRepo.CreateTask(ctx, task)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	return s.taskRepo.GetTaskByID(ctx, id)
}

func (s *TaskService) GetTasksByTeacherID(ctx context.Context, teacherID string) ([]models.Task, error) {
	return s.taskRepo.GetTasksByTeacherID(ctx, teacherID)
}

func (s *TaskService) UpdateTaskStatus(ctx context.Context, id string, status string) error {
	return s.taskRepo.UpdateTaskStatus(ctx, id, status)
}

func (s *TaskService) AssignTaskToStudents(ctx context.Context, taskID string, studentIDs []string) error {
	for _, studentID := range studentIDs {
		studentTask := &models.StudentTask{
			ID:        uuid.New().String(),
			TaskID:    taskID,
			StudentID: studentID,
			Status:    "pending",
			Progress:  0,
		}
		if err := s.taskRepo.CreateStudentTask(ctx, studentTask); err != nil {
			return err
		}
	}
	return nil
}

func (s *TaskService) GetStudentTasksByTaskID(ctx context.Context, taskID string) ([]models.StudentTask, error) {
	return s.taskRepo.GetStudentTasksByTaskID(ctx, taskID)
}

func (s *TaskService) UpdateStudentTaskStatus(ctx context.Context, taskID string, studentID string, status string, progress int) error {
	// 这里需要先找到对应的 StudentTask ID
	studentTasks, err := s.taskRepo.GetStudentTasksByTaskID(ctx, taskID)
	if err != nil {
		return err
	}
	
	for _, st := range studentTasks {
		if st.StudentID == studentID {
			return s.taskRepo.UpdateStudentTaskStatus(ctx, st.ID, status, progress)
		}
	}
	
	return nil
}

func (s *TaskService) GetStudents(ctx context.Context) ([]models.Student, error) {
	return s.userRepo.GetStudents(ctx)
}