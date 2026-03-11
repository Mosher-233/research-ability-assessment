package models

import (
	"time"
)

type Evidence struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	StudentTaskID  string    `json:"student_task_id" gorm:"not null"`
	Type           string    `json:"type" gorm:"not null"` // document, code, presentation, etc.
	Content        string    `json:"content" gorm:"not null"`
	KBMName        string    `json:"kbm_name" gorm:"not null"` // 知识、能力、方法名称
	KBMLevel       int       `json:"kbm_level" gorm:"not null"` // 知识、能力、方法水平
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	StudentTask    StudentTask `json:"student_task" gorm:"foreignKey:StudentTaskID"`
}