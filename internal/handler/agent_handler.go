package handler

import (
	"log"
	"net/http"
	"github.com/Mosher-233/research-ability-assessment/internal/agent"
	"github.com/Mosher-233/research-ability-assessment/internal/service"

	"github.com/gin-gonic/gin"
)

type AgentHandler struct {
	controlUnit      *agent.ControlUnit
	inferenceService *service.InferenceService
	reportService    *service.ReportService
}

func NewAgentHandler(controlUnit *agent.ControlUnit, inferenceService *service.InferenceService, reportService *service.ReportService) *AgentHandler {
	return &AgentHandler{
		controlUnit:      controlUnit,
		inferenceService: inferenceService,
		reportService:    reportService,
	}
}

type EvaluateRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	TaskID    string `json:"task_id" binding:"required"`
}

type EvaluateAllRequest struct {
	TaskID string `json:"task_id" binding:"required"`
}

type FeedbackRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	TaskID    string `json:"task_id" binding:"required"`
}

func (h *AgentHandler) EvaluateStudent(c *gin.Context) {
	var req EvaluateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("AgentHandler: 开始执行学生评估, StudentID=%s, TaskID=%s", req.StudentID, req.TaskID)

	result, err := h.controlUnit.ExecuteEvaluation(c.Request.Context(), req.TaskID, req.StudentID)
	if err != nil {
		log.Printf("AgentHandler: 评估失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "评估执行失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "评估执行成功",
		"data":    result,
	})
}

func (h *AgentHandler) EvaluateAllStudents(c *gin.Context) {
	var req EvaluateAllRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("AgentHandler: 开始批量评估, TaskID=%s", req.TaskID)

	studentTasks, err := h.controlUnit.GetStudentTasks(c.Request.Context(), req.TaskID)
	if err != nil {
		log.Printf("AgentHandler: 获取学生任务列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取学生任务列表失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	var results []*agent.EvaluationResult
	var errors []string

	for _, st := range studentTasks {
		result, err := h.controlUnit.ExecuteEvaluation(c.Request.Context(), req.TaskID, st.StudentID)
		if err != nil {
			log.Printf("AgentHandler: 学生 %s 评估失败: %v", st.StudentID, err)
			errors = append(errors, "学生 "+st.StudentID+" 评估失败: "+err.Error())
			continue
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "批量评估完成",
		"data": gin.H{
			"results": results,
			"errors":  errors,
			"total":   len(studentTasks),
			"success": len(results),
			"failed":  len(errors),
		},
	})
}

func (h *AgentHandler) GetEvaluationStatus(c *gin.Context) {
	studentID := c.Query("student_id")
	taskID := c.Query("task_id")

	if studentID == "" || taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "student_id 和 task_id 不能为空",
			"data":    nil,
		})
		return
	}

	result, err := h.inferenceService.GetInferenceResultByStudentAndTask(c.Request.Context(), studentID, taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "未找到评估结果",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取评估状态成功",
		"data":    result,
	})
}

func (h *AgentHandler) GenerateFeedback(c *gin.Context) {
	var req FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	inferenceResult, err := h.inferenceService.GetInferenceResultByStudentAndTask(c.Request.Context(), req.StudentID, req.TaskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "未找到推理结果，请先生成评估结果",
			"data":    nil,
		})
		return
	}

	feedback, err := h.controlUnit.GenerateFeedback(c.Request.Context(), inferenceResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成反馈失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "反馈生成成功",
		"data":    feedback,
	})
}
