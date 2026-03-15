package handler

import (
	"context"
	"log"
	"net/http"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/internal/service"

	"github.com/gin-gonic/gin"
)

type ResultHandler struct {
	inferenceService *service.InferenceService
	reportService    *service.ReportService
	resultRepo       interface{}
	taskService      *service.TaskService
	userRepo         interface {
		GetUserByID(ctx context.Context, id string) (*models.User, error)
	}
	taskRepo         interface {
		GetTaskByID(ctx context.Context, id string) (*models.Task, error)
	}
}

func NewResultHandler(inferenceService *service.InferenceService, reportService *service.ReportService, resultRepo interface{}, taskService *service.TaskService, userRepo interface{}, taskRepo interface{}) *ResultHandler {
	ur, _ := userRepo.(interface {
		GetUserByID(ctx context.Context, id string) (*models.User, error)
	})
	tr, _ := taskRepo.(interface {
		GetTaskByID(ctx context.Context, id string) (*models.Task, error)
	})
	return &ResultHandler{
		inferenceService: inferenceService,
		reportService:    reportService,
		resultRepo:       resultRepo,
		taskService:      taskService,
		userRepo:         ur,
		taskRepo:         tr,
	}
}

type GenerateInferenceRequest struct {
	StudentTaskID string `json:"student_task_id" binding:"required"`
	StudentID     string `json:"student_id" binding:"required"`
	TaskID        string `json:"task_id" binding:"required"`
}

type GenerateReportRequest struct {
	StudentTaskID string `json:"student_task_id" binding:"required"`
	StudentID     string `json:"student_id" binding:"required"`
	TaskID        string `json:"task_id" binding:"required"`
}

func (h *ResultHandler) GenerateInferenceResult(c *gin.Context) {
	var req GenerateInferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	result, err := h.inferenceService.GenerateInferenceWithLLM(c.Request.Context(), &service.GenerateInferenceRequest{
		StudentTaskID: req.StudentTaskID,
		StudentID:     req.StudentID,
		TaskID:        req.TaskID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成推理结果失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "生成推理结果成功",
		"data":    result,
	})
}

func (h *ResultHandler) GenerateReport(c *gin.Context) {
	var req GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	report, err := h.reportService.GenerateReport(c.Request.Context(), req.StudentTaskID, req.StudentID, req.TaskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成报告失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "生成报告成功",
		"data":    report,
	})
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

	// 首先获取学生任务ID
	log.Printf("GenerateStudentReport: 正在获取学生任务, StudentID=%s, TaskID=%s", studentID, taskID)
	studentTask, err := h.taskService.GetStudentTaskByStudentAndTask(c.Request.Context(), studentID, taskID)
	if err != nil {
		log.Printf("GenerateStudentReport: 未找到学生任务: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "未找到对应的学生任务",
			"data":    nil,
		})
		return
	}

	log.Printf("GenerateStudentReport: 开始生成报告, StudentTaskID=%s, StudentID=%s, TaskID=%s", studentTask.ID, studentID, taskID)

	report, err := h.reportService.GenerateReport(c.Request.Context(), studentTask.ID, studentID, taskID)
	if err != nil {
		log.Printf("GenerateStudentReport: 生成报告失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成学生报告失败: " + err.Error(),
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

func (h *ResultHandler) GetResults(c *gin.Context) {
	results, err := h.inferenceService.GetAllInferenceResults(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取结果失败",
			"data":    nil,
		})
		return
	}

	enrichedResults := h.enrichResultsWithDetails(c.Request.Context(), results)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取结果成功",
		"data":    enrichedResults,
	})
}

func (h *ResultHandler) GetStudentResults(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	results, err := h.inferenceService.GetInferenceResultsByStudentID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取结果失败",
			"data":    nil,
		})
		return
	}

	enrichedResults := h.enrichResultsWithDetails(c.Request.Context(), results)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取结果成功",
		"data":    enrichedResults,
	})
}

type EnrichedInferenceResult struct {
	*models.InferenceResult
	Student interface{} `json:"student"`
	Task    interface{} `json:"task"`
}

func (h *ResultHandler) enrichResultsWithDetails(ctx context.Context, results []models.InferenceResult) []map[string]interface{} {
	enrichedResults := make([]map[string]interface{}, 0, len(results))
	
	for _, result := range results {
		resultMap := make(map[string]interface{})
		
		resultMap["id"] = result.ID
		resultMap["student_id"] = result.StudentID
		resultMap["task_id"] = result.TaskID
		resultMap["overall_score"] = result.OverallScore
		resultMap["overall_level"] = result.OverallLevel
		resultMap["dimension_scores"] = result.DimensionScores
		resultMap["reasoning"] = result.Reasoning
		resultMap["created_at"] = result.CreatedAt
		resultMap["updated_at"] = result.UpdatedAt
		
		if h.userRepo != nil {
			student, _ := h.userRepo.GetUserByID(ctx, result.StudentID)
			if student != nil {
				studentMap := map[string]interface{}{
					"id":         student.ID,
					"name":       student.Name,
					"student_id": student.ID,
				}
				resultMap["student"] = studentMap
			}
		}
		
		if h.taskRepo != nil {
			task, _ := h.taskRepo.GetTaskByID(ctx, result.TaskID)
			if task != nil {
				resultMap["task"] = task
			}
		}
		
		enrichedResults = append(enrichedResults, resultMap)
	}
	
	return enrichedResults
}

func (h *ResultHandler) GetReports(c *gin.Context) {
	reports, err := h.reportService.GetAllReports(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取报告失败",
			"data":    nil,
		})
		return
	}

	enrichedReports := h.enrichReportsWithDetails(c.Request.Context(), reports)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取报告成功",
		"data":    enrichedReports,
	})
}

func (h *ResultHandler) GetStudentReports(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	reports, err := h.reportService.GetReportsByStudentID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取报告失败",
			"data":    nil,
		})
		return
	}

	enrichedReports := h.enrichReportsWithDetails(c.Request.Context(), reports)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取报告成功",
		"data":    enrichedReports,
	})
}

func (h *ResultHandler) enrichReportsWithDetails(ctx context.Context, reports []models.Report) []map[string]interface{} {
	enrichedReports := make([]map[string]interface{}, 0, len(reports))
	
	for _, report := range reports {
		reportMap := make(map[string]interface{})
		
		reportMap["id"] = report.ID
		reportMap["student_task_id"] = report.StudentTaskID
		reportMap["student_id"] = report.StudentID
		reportMap["task_id"] = report.TaskID
		reportMap["overall_score"] = report.OverallScore
		reportMap["overall_level"] = report.OverallLevel
		reportMap["dimension_scores"] = report.DimensionScores
		reportMap["class_comparison"] = report.ClassComparison
		reportMap["rank"] = report.Rank
		reportMap["percentile"] = report.Percentile
		reportMap["strengths"] = report.Strengths
		reportMap["weaknesses"] = report.Weaknesses
		reportMap["detailed_analysis"] = report.DetailedAnalysis
		reportMap["suggestions"] = report.Suggestions
		reportMap["radar_chart_data"] = report.RadarChartData
		reportMap["report_path"] = report.ReportPath
		reportMap["created_at"] = report.CreatedAt
		reportMap["updated_at"] = report.UpdatedAt
		
		if h.userRepo != nil {
			student, _ := h.userRepo.GetUserByID(ctx, report.StudentID)
			if student != nil {
				reportMap["student_name"] = student.Name
			}
		}
		
		if h.taskRepo != nil {
			task, _ := h.taskRepo.GetTaskByID(ctx, report.TaskID)
			if task != nil {
				reportMap["task_name"] = task.Name
			}
		}
		
		enrichedReports = append(enrichedReports, reportMap)
	}
	
	return enrichedReports
}

func (h *ResultHandler) GenerateStudentInference(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	taskID := c.Query("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "任务ID不能为空",
			"data":    nil,
		})
		return
	}

	log.Printf("GenerateStudentInference: 开始为学生生成推理结果, StudentID=%s, TaskID=%s", userID, taskID)

	// 获取学生任务
	studentTask, err := h.taskService.GetStudentTaskByStudentAndTask(c.Request.Context(), userID, taskID)
	if err != nil {
		log.Printf("GenerateStudentInference: 未找到学生任务: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "未找到对应的学生任务",
			"data":    nil,
		})
		return
	}

	log.Printf("GenerateStudentInference: 找到学生任务, StudentTaskID=%s", studentTask.ID)

	// 检查是否已经存在推理结果
	existingResult, _ := h.inferenceService.GetInferenceResultByStudentAndTask(c.Request.Context(), userID, taskID)
	if existingResult != nil {
		log.Printf("GenerateStudentInference: 推理结果已存在, ResultID=%s", existingResult.ID)
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "推理结果已存在",
			"data":    existingResult,
		})
		return
	}

	// 生成推理结果
	result, err := h.inferenceService.GenerateInferenceWithLLM(c.Request.Context(), &service.GenerateInferenceRequest{
		StudentTaskID: studentTask.ID,
		StudentID:     userID,
		TaskID:        taskID,
	})
	if err != nil {
		log.Printf("GenerateStudentInference: 生成推理结果失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成推理结果失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("GenerateStudentInference: 推理结果生成成功, ResultID=%s", result.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "生成推理结果成功",
		"data":    result,
	})
}