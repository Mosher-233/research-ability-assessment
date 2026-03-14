package handler

import (
	"log"
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
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	CourseID    string `json:"course_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
}

type AssignTaskRequest struct {
	StudentIDs []string `json:"student_ids" binding:"required"`
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("创建任务参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("收到创建任务请求: Name=%s, CourseID=%s, StartDate=%s, EndDate=%s", 
		req.Name, req.CourseID, req.StartDate, req.EndDate)

	userID := c.GetString("userID")
	if userID == "" {
		log.Printf("创建任务失败: 未授权")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		log.Printf("解析开始日期失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "开始日期格式错误",
			"data":    nil,
		})
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		log.Printf("解析结束日期失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "结束日期格式错误",
			"data":    nil,
		})
		return
	}

	task := &models.Task{
		Name:        req.Name,
		Description: req.Description,
		CourseID:    req.CourseID,
		TeacherID:   userID,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      "active",
	}

	log.Printf("准备创建任务: ID=%s, TeacherID=%s", task.ID, task.TeacherID)

	if err := h.taskService.CreateTask(c.Request.Context(), task); err != nil {
		log.Printf("创建任务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建任务失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("任务创建成功: ID=%s", task.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建任务成功",
		"data":    task,
	})
}

type TaskWithTeacher struct {
	models.Task
	Teacher *models.User `json:"teacher"`
}

type StudentTaskWithDetails struct {
	models.StudentTask
	Student *models.User  `json:"student"`
	Task    *models.Task `json:"task"`
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

	var teacher *models.User
	if task.TeacherID != "" {
		teacher, _ = h.taskService.GetUserByID(c.Request.Context(), task.TeacherID)
	}

	taskWithTeacher := TaskWithTeacher{
		Task:    *task,
		Teacher: teacher,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取任务成功",
		"data":    taskWithTeacher,
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
		log.Printf("分配任务失败: 任务ID为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "任务ID不能为空",
			"data":    nil,
		})
		return
	}

	log.Printf("收到分配任务请求: TaskID=%s", taskID)

	var req AssignTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("分配任务参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("分配任务给学生: StudentIDs=%v", req.StudentIDs)

	if err := h.taskService.AssignTaskToStudents(c.Request.Context(), taskID, req.StudentIDs); err != nil {
		log.Printf("分配任务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "分配任务失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("任务分配成功: TaskID=%s", taskID)

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

func (h *TaskHandler) GetStudentTasksList(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	studentTasks, err := h.taskService.GetStudentTasksByStudentID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取学生任务失败",
			"data":    nil,
		})
		return
	}

	var results []StudentTaskWithDetails
	for _, st := range studentTasks {
		student, _ := h.taskService.GetUserByID(c.Request.Context(), st.StudentID)
		task, _ := h.taskService.GetTaskByID(c.Request.Context(), st.TaskID)
		
		results = append(results, StudentTaskWithDetails{
			StudentTask: st,
			Student:     student,
			Task:        task,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取学生任务成功",
		"data":    results,
	})
}