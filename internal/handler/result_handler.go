package handler

import (
	"net/http"
	"research-ability-assessment/internal/service"

	"github.com/gin-gonic/gin"
)

type ResultHandler struct {
	inferenceService *service.InferenceService
	reportService    *service.ReportService
}

func NewResultHandler(inferenceService *service.InferenceService, reportService *service.ReportService) *ResultHandler {
	return &ResultHandler{
		inferenceService: inferenceService,
		reportService:    reportService,
	}
}

func (h *ResultHandler) GetInferenceResultByID(c *gin.Context) {
	resultID := c.Param("result_id")
	if resultID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "结果ID不能为空",
			"data":    nil,
		})
		return
	}

	result, err := h.inferenceService.GetInferenceResultByID(c.Request.Context(), resultID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取推理结果失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取推理结果成功",
		"data":    result,
	})
}

func (h *ResultHandler) GetInferenceResultsByTaskID(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "任务ID不能为空",
			"data":    nil,
		})
		return
	}

	results, err := h.inferenceService.GetInferenceResultsByTaskID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取推理结果失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取推理结果成功",
		"data":    results,
	})
}

func (h *ResultHandler) GetInferenceResultByStudentAndTask(c *gin.Context) {
	studentID := c.Query("student_id")
	taskID := c.Query("task_id")

	if studentID == "" || taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "学生ID和任务ID不能为空",
			"data":    nil,
		})
		return
	}

	result, err := h.inferenceService.GetInferenceResultByStudentAndTask(c.Request.Context(), studentID, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取推理结果失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取推理结果成功",
		"data":    result,
	})
}

func (h *ResultHandler) GenerateStudentReport(c *gin.Context) {
	studentID := c.Query("student_id")
	taskID := c.Query("task_id")

	if studentID == "" || taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "学生ID和任务ID不能为空",
			"data":    nil,
		})
		return
	}

	report, err := h.reportService.GenerateStudentReport(c.Request.Context(), studentID, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成学生报告失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "生成学生报告成功",
		"data":    report,
	})
}

func (h *ResultHandler) GenerateTaskReport(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "任务ID不能为空",
			"data":    nil,
		})
		return
	}

	report, err := h.reportService.GenerateTaskReport(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成任务报告失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "生成任务报告成功",
		"data":    report,
	})
}