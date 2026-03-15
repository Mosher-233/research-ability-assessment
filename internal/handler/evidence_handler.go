package handler

import (
	"net/http"
	"os"
	"path/filepath"
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
	Content      string `json:"content"`
	KBMName      string `json:"kbm_name" binding:"required"`
	KBMLevel     int    `json:"kbm_level" binding:"omitempty,min=1,max=5"`
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
	FileType     string `json:"file_type"`
	FileSize     int64  `json:"file_size"`
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
		FileName:     req.FileName,
		FilePath:     req.FilePath,
		FileType:     req.FileType,
		FileSize:     req.FileSize,
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

func (h *EvidenceHandler) UploadEvidenceFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择文件",
			"data":    nil,
		})
		return
	}

	studentTaskID := c.PostForm("student_task_id")
	evidenceType := c.PostForm("type")
	kbmName := c.PostForm("kbm_name")
	if studentTaskID == "" || evidenceType == "" || kbmName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数不完整",
			"data":    nil,
		})
		return
	}

	uploadDir := "./uploads/evidences"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建上传目录失败",
			"data":    nil,
		})
		return
	}

	fileName := file.Filename
	filePath := filepath.Join(uploadDir, fileName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存文件失败",
			"data":    nil,
		})
		return
	}

	content := ""
	ext := filepath.Ext(fileName)
	if ext == ".txt" || ext == ".md" {
		fileContent, err := os.ReadFile(filePath)
		if err == nil {
			content = string(fileContent)
		}
	}

	evidence := &models.Evidence{
		StudentTaskID: studentTaskID,
		Type:         evidenceType,
		Content:      content,
		KBMName:      kbmName,
		FileName:     fileName,
		FilePath:     filePath,
		FileType:     ext,
		FileSize:     file.Size,
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
		"message": "上传证据成功",
		"data":    evidence,
	})
}

func (h *EvidenceHandler) DownloadFile(c *gin.Context) {
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
	if err != nil || evidence.FilePath == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文件不存在",
			"data":    nil,
		})
		return
	}

	c.FileAttachment(evidence.FilePath, evidence.FileName)
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

func (h *EvidenceHandler) GetEvidences(c *gin.Context) {
	userID := c.GetString("userID")
	userRole := c.GetString("role")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	var evidences interface{}
	var err error

	if userRole == "teacher" {
		evidences, err = h.evidenceService.GetEvidencesWithDetailsByTeacherID(c.Request.Context(), userID)
	} else {
		evidences, err = h.evidenceService.GetEvidencesByUserID(c.Request.Context(), userID)
	}

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

func (h *EvidenceHandler) AnalyzeEvidence(c *gin.Context) {
	evidenceID := c.Param("evidence_id")
	if evidenceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "证据ID不能为空",
			"data":    nil,
		})
		return
	}

	result, err := h.evidenceService.AnalyzeEvidence(c.Request.Context(), evidenceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "分析证据失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "分析证据成功",
		"data":    result,
	})
}

func (h *EvidenceHandler) GetFeedbackByEvidenceID(c *gin.Context) {
	evidenceID := c.Param("evidence_id")
	if evidenceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "证据ID不能为空",
			"data":    nil,
		})
		return
	}

	feedback, err := h.evidenceService.GetFeedbackByEvidenceID(c.Request.Context(), evidenceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取反馈失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取反馈成功",
		"data":    feedback,
	})
}

func (h *EvidenceHandler) GetFeedbacks(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	feedbacks, err := h.evidenceService.GetFeedbacksByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取反馈失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取反馈成功",
		"data":    feedbacks,
	})
}

func (h *EvidenceHandler) DeleteEvidence(c *gin.Context) {
	evidenceID := c.Param("evidence_id")
	if evidenceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "证据ID不能为空",
			"data":    nil,
		})
		return
	}

	if err := h.evidenceService.DeleteEvidence(c.Request.Context(), evidenceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除证据失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除证据成功",
		"data":    nil,
	})
}