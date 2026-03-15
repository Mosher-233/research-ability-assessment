package service

import (
	"context"
	"log"
	"research-ability-assessment/internal/models"
	"research-ability-assessment/internal/repository/postgres"
	"research-ability-assessment/pkg/utils"
)

type TaskService struct {
	taskRepo *postgres.TaskRepo
	userRepo *postgres.UserRepo
}

func NewTaskService(taskRepo *postgres.TaskRepo, userRepo *postgres.UserRepo) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task) error {
	task.ID = utils.GenerateTaskID()
	log.Printf("TaskService: 生成任务ID: %s", task.ID)
	err := s.taskRepo.CreateTask(ctx, task)
	if err != nil {
		log.Printf("TaskService: 创建任务失败: %v", err)
	} else {
		log.Printf("TaskService: 任务创建成功: %s", task.ID)
	}
	return err
}

func (s *TaskService) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	return s.taskRepo.GetTaskByID(ctx, id)
}

func (s *TaskService) GetTasksByTeacherID(ctx context.Context, teacherID string) ([]models.Task, error) {
	return s.taskRepo.GetTasksByTeacherID(ctx, teacherID)
}

func (s *TaskService) UpdateTaskStatus(ctx context.Context, id string, status string) error {
	return s.taskRepo.UpdateTaskStatus(ctx, id, status)
}

func (s *TaskService) AssignTaskToStudents(ctx context.Context, taskID string, studentIDs []string) error {
	log.Printf("TaskService: 开始分配任务 %s 给 %d 个学生", taskID, len(studentIDs))
	
	// 获取已分配的学生
	existingStudentTasks, err := s.taskRepo.GetStudentTasksByTaskID(ctx, taskID)
	if err != nil {
		log.Printf("TaskService: 获取已分配学生失败: %v", err)
		return err
	}
	
	log.Printf("TaskService: 已分配学生数量: %d", len(existingStudentTasks))
	
	// 创建已分配学生ID集合
	assignedStudentIDs := make(map[string]bool)
	for _, st := range existingStudentTasks {
		assignedStudentIDs[st.StudentID] = true
	}
	
	// 只分配尚未分配的学生
	newlyAssigned := 0
	for _, studentID := range studentIDs {
		if !assignedStudentIDs[studentID] {
			studentTask := &models.StudentTask{
				ID:        utils.GenerateStudentTaskID(),
				TaskID:    taskID,
				StudentID: studentID,
				Status:    "pending",
				Progress:  0,
			}
			log.Printf("TaskService: 为学生 %s 创建学生任务 %s", studentID, studentTask.ID)
			if err := s.taskRepo.CreateStudentTask(ctx, studentTask); err != nil {
				log.Printf("TaskService: 创建学生任务失败: %v", err)
				return err
			}
			newlyAssigned++
		} else {
			log.Printf("TaskService: 学生 %s 已分配，跳过", studentID)
		}
	}
	
	log.Printf("TaskService: 任务分配完成，新分配 %d 个学生", newlyAssigned)
	return nil
}

func (s *TaskService) GetStudentTasksByTaskID(ctx context.Context, taskID string) ([]models.StudentTask, error) {
	return s.taskRepo.GetStudentTasksByTaskID(ctx, taskID)
}

func (s *TaskService) UpdateStudentTaskStatus(ctx context.Context, taskID string, studentID string, status string, progress int) error {
	// 这里需要先找到对应的 StudentTask ID
	studentTasks, err := s.taskRepo.GetStudentTasksByTaskID(ctx, taskID)
	if err != nil {
		return err
	}
	
	for _, st := range studentTasks {
		if st.StudentID == studentID {
			return s.taskRepo.UpdateStudentTaskStatus(ctx, st.ID, status, progress)
		}
	}
	
	return nil
}

func (s *TaskService) GetStudents(ctx context.Context) ([]models.Student, error) {
	return s.userRepo.GetStudents(ctx)
}

func (s *TaskService) GetAssignedTasks(ctx context.Context, studentID string) ([]models.Task, error) {
	return s.taskRepo.GetAssignedTasks(ctx, studentID)
}

func (s *TaskService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *TaskService) GetStudentTasksByStudentID(ctx context.Context, studentID string) ([]models.StudentTask, error) {
	return s.taskRepo.GetStudentTasksByStudentID(ctx, studentID)
}

func (s *TaskService) GetStudentTaskByStudentAndTask(ctx context.Context, studentID string, taskID string) (*models.StudentTask, error) {
	return s.taskRepo.GetStudentTaskByStudentAndTask(ctx, studentID, taskID)
}