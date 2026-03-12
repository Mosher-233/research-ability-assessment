package models

import (
	"time"
)

type Task struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	CourseID    string    `json:"course_id" gorm:"not null"`
	TeacherID   string    `json:"teacher_id" gorm:"not null"`
	StartDate   time.Time `json:"start_date" gorm:"not null"`
	EndDate     time.Time `json:"end_date" gorm:"not null"`
	Status      string    `json:"status" gorm:"not null"` // active, completed, archived
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// 移除关联关系，避免迁移问题
	// Teacher     Teacher    `json:"teacher" gorm:"foreignKey:TeacherID"`
	// StudentTasks []StudentTask `json:"student_tasks" gorm:"foreignKey:TaskID"`
}

type StudentTask struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	TaskID    string    `json:"task_id" gorm:"not null"`
	StudentID string    `json:"student_id" gorm:"not null"`
	Status    string    `json:"status" gorm:"not null"` // pending, processing, completed
	Progress  int       `json:"progress" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// 移除关联关系，避免迁移问题
	// Task      Task      `json:"task" gorm:"foreignKey:TaskID"`
	// Student   Student   `json:"student" gorm:"foreignKey:StudentID"`
	// Evidences []Evidence `json:"evidences" gorm:"foreignKey:StudentTaskID"`
}