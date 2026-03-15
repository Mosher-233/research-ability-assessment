package models

import (
	"time"

	"gorm.io/datatypes"
)

type InferenceResult struct {
	ID              string            `json:"id" gorm:"primaryKey"`
	StudentID       string            `json:"student_id" gorm:"not null"`
	TaskID          string            `json:"task_id" gorm:"not null"`
	OverallScore    float64           `json:"overall_score" gorm:"not null"`
	OverallLevel    string            `json:"overall_level" gorm:"not null"`
	DimensionScores datatypes.JSON    `json:"dimension_scores" gorm:"type:json"`
	Reasoning       string            `json:"reasoning" gorm:"not null"`
	CreatedAt       time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	
	// 移除关联关系，避免迁移问题
	// Student         Student           `json:"student" gorm:"foreignKey:StudentID"`
	// Task            Task              `json:"task" gorm:"foreignKey:TaskID"`
}

type DimensionScore struct {
	Name        string   `json:"name"`
	Score       float64  `json:"score"`
	Level       string   `json:"level"`
	Details     string   `json:"details"`
	EvidenceIDs []string `json:"evidence_ids"`
}