package handler

import (
	"net/http"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

type CreateTaskRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	CourseID    string    `json:"course_id" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
}

type AssignTaskRequest struct {
	StudentIDs []string `json:"student_ids" binding:"required"`
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	task := &models.Task{
		Name:        req.Name,
		Description: req.Description,
		CourseID:    req.CourseID,
		TeacherID:   userID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      "active",
	}

	if err := h.taskService.CreateTask(c.Request.Context(), task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建任务失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建任务成功",
		"data":    task,
	})
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "任务ID不能为空",
			"data":    nil,
		})
		return
	}

	task, err := h.taskService.GetTaskByID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取任务失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取任务成功",
		"data":    task,
	})
}

func (h *TaskHandler) GetTasksByTeacherID(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	tasks, err := h.taskService.GetTasksByTeacherID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取任务失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取任务成功",
		"data":    tasks,
	})
}

func (h *TaskHandler) AssignTask(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "任务ID不能为空",
			"data":    nil,
		})
		return
	}

	var req AssignTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	if err := h.taskService.AssignTaskToStudents(c.Request.Context(), taskID, req.StudentIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "分配任务失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "分配任务成功",
		"data":    nil,
	})
}

func (h *TaskHandler) GetStudentTasks(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "任务ID不能为空",
			"data":    nil,
		})
		return
	}

	studentTasks, err := h.taskService.GetStudentTasksByTaskID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取学生任务失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取学生任务成功",
		"data":    studentTasks,
	})
}

func (h *TaskHandler) GetStudents(c *gin.Context) {
	students, err := h.taskService.GetStudents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取学生列表失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取学生列表成功",
		"data":    students,
	})
}

func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "任务ID不能为空",
			"data":    nil,
		})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	if err := h.taskService.UpdateTaskStatus(c.Request.Context(), taskID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新任务状态失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新任务状态成功",
		"data":    nil,
	})
}

func (h *TaskHandler) GetAssignedTasks(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	tasks, err := h.taskService.GetAssignedTasks(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取任务失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取任务成功",
		"data":    tasks,
	})
}