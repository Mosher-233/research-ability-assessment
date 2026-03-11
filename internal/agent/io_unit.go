package agent

import (
	"context"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/internal/service"
)

type IOUnit struct {
	taskService     *service.TaskService
	evidenceService *service.EvidenceService
}

func NewIOUnit(taskService *service.TaskService, evidenceService *service.EvidenceService) *IOUnit {
	return &IOUnit{
		taskService:     taskService,
		evidenceService: evidenceService,
	}
}

func (io *IOUnit) GetTask(ctx context.Context, taskID string) (*models.Task, error) {
	return io.taskService.GetTaskByID(ctx, taskID)
}

func (io *IOUnit) GetStudentTask(ctx context.Context, taskID string, studentID string) (*models.StudentTask, error) {
	studentTasks, err := io.taskService.GetStudentTasksByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	
	for _, st := range studentTasks {
		if st.StudentID == studentID {
			return &st, nil
		}
	}
	
	return nil, nil
}

func (io *IOUnit) GetEvidences(ctx context.Context, studentID string, taskID string) ([]*models.Evidence, error) {
	evidences, err := io.evidenceService.GetEvidencesByStudentAndTask(ctx, studentID, taskID)
	if err != nil {
		return nil, err
	}
	
	// 转换为指针切片
	evidencePtrs := make([]*models.Evidence, len(evidences))
	for i, evidence := range evidences {
		evidencePtrs[i] = &evidence
	}
	
	return evidencePtrs, nil
}

func (io *IOUnit) UpdateTaskStatus(ctx context.Context, taskID string, studentID string, status string, progress int) error {
	return io.taskService.UpdateStudentTaskStatus(ctx, taskID, studentID, status, progress)
}