package models

import (
	"time"
)

type Feedback struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	EvidenceID  string    `json:"evidence_id" gorm:"not null;index"`
	Content     string    `json:"content" gorm:"type:text;not null"`
	KBMLevel    int       `json:"kbm_level" gorm:"not null"`
	Strengths   string    `json:"strengths" gorm:"type:text"`
	Weaknesses  string    `json:"weaknesses" gorm:"type:text"`
	Suggestions string    `json:"suggestions" gorm:"type:text"`
	FileName    string    `json:"file_name" gorm:"size:255"`    // 反馈文件名
	FilePath    string    `json:"file_path" gorm:"size:500"`    // 反馈文件路径
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
