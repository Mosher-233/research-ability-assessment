package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role      string    `json:"role" gorm:"not null"` // teacher, student
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Teacher struct {
	User
	Department string `json:"department" gorm:"not null"`
	Title      string `json:"title"`
}

type Student struct {
	User
	StudentID string `json:"student_id" gorm:"unique;not null"`
	Major     string `json:"major" gorm:"not null"`
	Grade     string `json:"grade" gorm:"not null"`
}