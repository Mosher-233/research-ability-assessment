package models

import (
	"time"
)

type Evidence struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	StudentTaskID  string    `json:"student_task_id" gorm:"not null"`
	Type           string    `json:"type" gorm:"not null"` // document, code, presentation, etc.
	Content        string    `json:"content" gorm:"type:text"`
	FileName       string    `json:"file_name" gorm:"size:255"`    // 文件名
	FilePath       string    `json:"file_path" gorm:"size:500"`    // 文件路径
	FileType       string    `json:"file_type" gorm:"size:100"`    // 文件类型
	FileSize       int64     `json:"file_size"`                      // 文件大小
	KBMName        string    `json:"kbm_name" gorm:"not null"` // 知识、能力、方法名称
	KBMLevel       int       `json:"kbm_level" gorm:"default:0"` // 知识、能力、方法水平
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}