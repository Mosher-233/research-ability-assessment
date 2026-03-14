package models

import (
	"time"
)

type Report struct {
	ID              string                 `json:"id" gorm:"primaryKey"`
	StudentTaskID   string                 `json:"student_task_id" gorm:"not null;index"`
	StudentID       string                 `json:"student_id" gorm:"not null;index"`
	TaskID          string                 `json:"task_id" gorm:"not null;index"`
	OverallScore    float64                `json:"overall_score" gorm:"not null"`
	OverallLevel    string                 `json:"overall_level" gorm:"not null"`
	DimensionScores map[string]DimensionScore `json:"dimension_scores" gorm:"type:json"`
	ClassComparison *ClassComparisonData    `json:"class_comparison" gorm:"type:json"`
	Rank            int                    `json:"rank" gorm:"default:0"`
	Percentile      float64                `json:"percentile" gorm:"default:0"`
	Strengths       []string               `json:"strengths" gorm:"type:json"`
	Weaknesses      []string               `json:"weaknesses" gorm:"type:json"`
	DetailedAnalysis map[string]string     `json:"detailed_analysis" gorm:"type:json"`
	Suggestions     []ImprovementSuggestion `json:"suggestions" gorm:"type:json"`
	RadarChartData  *RadarChartData        `json:"radar_chart_data" gorm:"type:json"`
	ReportPath      string                 `json:"report_path" gorm:"size:500"`
	CreatedAt       time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

type ClassComparisonData struct {
	ClassSize          int                `json:"class_size"`
	ClassAverage       float64            `json:"class_average"`
	ClassMaxScore      float64            `json:"class_max_score"`
	ClassMinScore      float64            `json:"class_min_score"`
	DimensionAverages  map[string]float64 `json:"dimension_averages"`
}

type ImprovementSuggestion struct {
	ID            string            `json:"id"`
	Dimension     string            `json:"dimension"`
	DimensionName string            `json:"dimension_name"`
	CurrentScore  float64           `json:"current_score"`
	TargetScore   float64           `json:"target_score"`
	Suggestion    string            `json:"suggestion"`
	ActionItems   []string          `json:"action_items"`
	Resources     []LearningResource `json:"resources"`
	Priority      int               `json:"priority"`
}

type LearningResource struct {
	Type     string `json:"type"`     // course/book/online
	Title    string `json:"title"`    // 资源标题
	URL      string `json:"url"`      // 资源链接
	Duration string `json:"duration"` // 学习时长
	Author   string `json:"author"`   // 作者（书籍）
	ISBN     string `json:"isbn"`     // ISBN（书籍）
}

type RadarChartData struct {
	Labels        []string  `json:"labels"`         // 维度名称
	Values        []float64 `json:"values"`         // 学生得分
	ClassAverages []float64 `json:"class_averages"` // 班级平均分（可选）
}

type InferenceRequest struct {
	StudentTaskID string     `json:"student_task_id"`
	Evidences     []Evidence `json:"evidences"`
	Feedbacks     []Feedback `json:"feedbacks"`
	Dimensions    []Dimension `json:"dimensions"`
	StudentInfo   *User      `json:"student_info"`
	TaskInfo      *Task      `json:"task_info"`
}

type Dimension struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Weight      float64 `json:"weight"`
}
