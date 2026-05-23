package models

import "time"

// EvidenceCitation tracks which assessment conclusion cites which evidence passage.
// It provides an auditable chain from evaluation result → evidence → original text.
type EvidenceCitation struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	ResultID       string    `json:"result_id" gorm:"not null;index"`
	DimensionID    string    `json:"dimension_id" gorm:"not null;index"`
	EvidenceID     string    `json:"evidence_id" gorm:"not null;index"`
	ExcerptText    string    `json:"excerpt_text" gorm:"type:text"`
	ExcerptOffset  int       `json:"excerpt_offset"`
	RelevanceScore float64   `json:"relevance_score"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
}
