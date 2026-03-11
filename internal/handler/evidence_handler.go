package handler

import (
	"net/http"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/internal/service"

	"github.com/gin-gonic/gin"
)

type EvidenceHandler struct {
	evidenceService *service.EvidenceService
}

func NewEvidenceHandler(evidenceService *service.EvidenceService) *EvidenceHandler {
	return &EvidenceHandler{evidenceService: evidenceService}
}

type CreateEvidenceRequest struct {
	StudentTaskID string `json:"student_task_id" binding:"required"`
	Type         string `json:"type" binding:"required"`
	Content      string `json:"content" binding:"required"`
	KBMName      string `json:"kbm_name" binding:"required"`
	KBMLevel     int    `json:"kbm_level" binding:"required,min=1,max=5"`
}

func (h *EvidenceHandler) CreateEvidence(c *gin.Context) {
	var req CreateEvidenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	evidence := &models.Evidence{
		StudentTaskID: req.StudentTaskID,
		Type:         req.Type,
		Content:      req.Content,
		KBMName:      req.KBMName,
		KBMLevel:     req.KBMLevel,
	}

	if err := h.evidenceService.CreateEvidence(c.Request.Context(), evidence); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建证据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建证据成功",
		"data":    evidence,
	})
}

func (h *EvidenceHandler) GetEvidenceByID(c *gin.Context) {
	evidenceID := c.Param("evidence_id")
	if evidenceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "证据ID不能为空",
			"data":    nil,
		})
		return
	}

	evidence, err := h.evidenceService.GetEvidenceByID(c.Request.Context(), evidenceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取证据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取证据成功",
		"data":    evidence,
	})
}

func (h *EvidenceHandler) GetEvidencesByStudentTaskID(c *gin.Context) {
	studentTaskID := c.Param("student_task_id")
	if studentTaskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "学生任务ID不能为空",
			"data":    nil,
		})
		return
	}

	evidences, err := h.evidenceService.GetEvidencesByStudentTaskID(c.Request.Context(), studentTaskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取证据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取证据成功",
		"data":    evidences,
	})
}

func (h *EvidenceHandler) GetEvidencesByStudentAndTask(c *gin.Context) {
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

	evidences, err := h.evidenceService.GetEvidencesByStudentAndTask(c.Request.Context(), studentID, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取证据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取证据成功",
		"data":    evidences,
	})
}