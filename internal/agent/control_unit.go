package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"github.com/Mosher-233/research-ability-assessment/internal/models"
	"github.com/Mosher-233/research-ability-assessment/internal/service"
	"time"

	"github.com/google/uuid"
)

type ControlUnit struct {
	taskService      *service.TaskService
	inferenceService *service.InferenceService
	inferenceAgent   *InferenceAgent
	feedbackAgent    *FeedbackAgent
	storage          *StorageUnit
}

type EvaluationTask struct {
	TaskID string
	StudentID string
	Progress  int
	Status    string
}

type EvaluationResult struct {
	StudentID       string                            `json:"student_id"`
	TaskID          string                            `json:"task_id"`
	OverallScore    float64                           `json:"overall_score"`
	OverallLevel    string                            `json:"overall_level"`
	DimensionScores map[string]models.DimensionScore  `json:"dimension_scores"`
	Reasoning       string                            `json:"reasoning"`
	Feedback        *FeedbackResult                   `json:"feedback,omitempty"`
}

func NewControlUnit(
	taskService *service.TaskService,
	inferenceService *service.InferenceService,
	inferenceAgent *InferenceAgent,
	feedbackAgent *FeedbackAgent,
	storage *StorageUnit,
) *ControlUnit {
	return &ControlUnit{
		taskService:      taskService,
		inferenceService: inferenceService,
		inferenceAgent:   inferenceAgent,
		feedbackAgent:    feedbackAgent,
		storage:          storage,
	}
}

func (c *ControlUnit) ExecuteEvaluation(ctx context.Context, taskID string, studentID string) (*EvaluationResult, error) {
	task := &EvaluationTask{
		TaskID:    taskID,
		StudentID: studentID,
		Progress:  0,
		Status:    "processing",
	}

	c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "processing", task.Progress)

	task.Progress = 50
	c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "processing", task.Progress)

	result, err := c.inferenceAgent.InferAbility(ctx, task.StudentID, task.TaskID)
	if err != nil {
		return nil, fmt.Errorf("能力推理失败: %w", err)
	}

	dimensionScoresJSON, _ := json.Marshal(result.DimensionScores)

	inferenceResult := &models.InferenceResult{
		ID:              uuid.New().String(),
		StudentID:       task.StudentID,
		TaskID:          task.TaskID,
		OverallScore:    result.OverallScore,
		OverallLevel:    result.OverallLevel,
		DimensionScores: dimensionScoresJSON,
		Reasoning:       result.Reasoning,
		CreatedAt:       time.Now(),
	}

	if err := c.storage.StoreInferenceResult(ctx, inferenceResult); err != nil {
		return nil, fmt.Errorf("存储推理结果失败: %w", err)
	}

	if len(result.Citations) > 0 {
		if err := c.storage.StoreCitations(ctx, inferenceResult.ID, result.Citations); err != nil {
			log.Printf("ControlUnit: 存储引用记录失败: %v", err)
		}
	}

	for dimension, score := range result.DimensionScores {
		if err := c.storage.UpdateKnowledgeGraph(ctx, task.StudentID, dimension, score.Score); err != nil {
			log.Printf("更新知识图谱失败: %v", err)
		}
	}

	task.Progress = 70
	c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "processing", task.Progress)

	feedbackResult, err := c.generateFeedback(ctx, inferenceResult)
	if err != nil {
		log.Printf("ControlUnit: 反馈生成失败: %v", err)
	}

	task.Progress = 100
	task.Status = "completed"
	c.taskService.UpdateStudentTaskStatus(ctx, task.TaskID, task.StudentID, "completed", task.Progress)

	evalResult := &EvaluationResult{
		StudentID:       inferenceResult.StudentID,
		TaskID:          inferenceResult.TaskID,
		OverallScore:    inferenceResult.OverallScore,
		OverallLevel:    inferenceResult.OverallLevel,
		DimensionScores: result.DimensionScores,
		Reasoning:       inferenceResult.Reasoning,
		Feedback:        feedbackResult,
	}

	return evalResult, nil
}

func (c *ControlUnit) generateFeedback(ctx context.Context, result *models.InferenceResult) (*FeedbackResult, error) {
	if c.feedbackAgent == nil {
		return nil, nil
	}

	var dimensionScores map[string]models.DimensionScore
	if len(result.DimensionScores) > 0 {
		json.Unmarshal(result.DimensionScores, &dimensionScores)
	}

	agentResult := &InferenceResult{
		OverallScore:    result.OverallScore,
		OverallLevel:    result.OverallLevel,
		DimensionScores: dimensionScores,
		Reasoning:       result.Reasoning,
	}

	return c.feedbackAgent.GenerateFeedback(ctx, agentResult)
}

type StudentTaskInfo struct {
	StudentID string
	TaskID    string
}

func (c *ControlUnit) GetStudentTasks(ctx context.Context, taskID string) ([]StudentTaskInfo, error) {
	tasks, err := c.taskService.GetStudentTasksByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	var result []StudentTaskInfo
	for _, t := range tasks {
		result = append(result, StudentTaskInfo{
			StudentID: t.StudentID,
			TaskID:    t.TaskID,
		})
	}

	return result, nil
}

func (c *ControlUnit) GenerateFeedback(ctx context.Context, result *models.InferenceResult) (*FeedbackResult, error) {
	return c.generateFeedback(ctx, result)
}
