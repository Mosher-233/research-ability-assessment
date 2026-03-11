package agent

import (
	"context"
	"fmt"
	"log"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/internal/service"
	"time"

	"github.com/google/uuid"
)

type ControlUnit struct {
	taskService      *service.TaskService
	evidenceService  *service.EvidenceService
	inferenceService *service.InferenceService
	inferenceAgent   *InferenceAgent
	storage          *StorageUnit
}

type EvaluationTask struct {
	TaskID     string
	StudentID  string
	Progress   int
	Status     string
}

func NewControlUnit(
	taskService *service.TaskService,
	evidenceService *service.EvidenceService,
	inferenceService *service.InferenceService,
	inferenceAgent *InferenceAgent,
	storage *StorageUnit,
) *ControlUnit {
	return &ControlUnit{
		taskService:      taskService,
		evidenceService:  evidenceService,
		inferenceService: inferenceService,
		inferenceAgent:   inferenceAgent,
		storage:          storage,
	}
}

func (c *ControlUnit) ExecuteEvaluation(ctx context.Context, taskID string, studentID string) error {
	task := &EvaluationTask{
		TaskID:     taskID,
		StudentID:  studentID,
		Progress:   0,
		Status:     "processing",
	}

	// 更新任务状态为处理中
	c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "processing", task.Progress)

	// 执行证据收集
	if err := c.executeEvidenceCollection(ctx, task); err != nil {
		return fmt.Errorf("证据收集失败: %w", err)
	}

	// 执行能力推理
	if err := c.executeInference(ctx, task); err != nil {
		return fmt.Errorf("能力推理失败: %w", err)
	}

	// 完成任务
	task.Progress = 100
	task.Status = "completed"
	c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "completed", task.Progress)

	return nil
}

func (c *ControlUnit) executeEvidenceCollection(ctx context.Context, task *EvaluationTask) error {
	// 更新进度
	task.Progress = 20
	c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "processing", task.Progress)

	// 这里可以实现证据收集的逻辑
	// 例如从学生提交的材料中提取证据

	// 模拟证据收集过程
	time.Sleep(1 * time.Second)

	// 更新进度
	task.Progress = 30
	c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "processing", task.Progress)

	return nil
}

func (c *ControlUnit) executeInference(ctx context.Context, task *EvaluationTask) error {
	// 更新进度
	task.Progress = 50
	c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "processing", task.Progress)

	// 执行能力推理
	result, err := c.inferenceAgent.InferAbility(ctx, task.StudentID, task.TaskID)
	if err != nil {
		return fmt.Errorf("能力推理失败: %w", err)
	}

	// 存储推理结果
	inferenceResult := &models.InferenceResult{
		ID:           uuid.New().String(),
		StudentID:    task.StudentID,
		TaskID:       task.TaskID,
		OverallScore: result.OverallScore,
		OverallLevel: result.OverallLevel,
		DimensionScores: result.DimensionScores,
		Reasoning:    result.Reasoning,
		CreatedAt:    time.Now(),
	}

	if err := c.storage.StoreInferenceResult(ctx, inferenceResult); err != nil {
		return fmt.Errorf("存储推理结果失败: %w", err)
	}

	// 更新知识图谱
	for dimension, score := range result.DimensionScores {
		if err := c.storage.UpdateKnowledgeGraph(ctx, task.StudentID, dimension, score.Score); err != nil {
			log.Printf("更新知识图谱失败: %v", err)
		}
	}

	// 更新进度
	task.Progress = 70
	c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "processing", task.Progress)

	return nil
}