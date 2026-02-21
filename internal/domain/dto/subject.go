package dto

import "github.com/google/uuid"

type SubjectRes struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	ExamType     string     `json:"exam_type"`
	DisplayOrder int16      `json:"display_order"`
	TopicCount   int        `json:"topic_count,omitempty"`
}

type TopicRes struct {
	ID           uuid.UUID `json:"id"`
	SubjectID    uuid.UUID `json:"subject_id"`
	Name         string    `json:"name"`
	DisplayOrder int16     `json:"display_order"`
	IsCompleted  bool      `json:"is_completed"` // öğrenciye göre doldurulur
}

type SubjectProgressRes struct {
	SubjectID       uuid.UUID `json:"subject_id"`
	SubjectName     string    `json:"subject_name"`
	TotalTopics     int       `json:"total_topics"`
	CompletedTopics int       `json:"completed_topics"`
	Percentage      float64   `json:"percentage"`
}

type MarkTopicReq struct {
	IsCompleted bool `json:"is_completed"`
}
