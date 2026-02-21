package dto

import (
	"time"

	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/google/uuid"
)

type SubjectScoresReq struct {
	Correct int `json:"correct" binding:"min=0"`
	Wrong   int `json:"wrong"   binding:"min=0"`
}

type CreateExamResultReq struct {
	ExamType string                      `json:"exam_type" binding:"required,oneof=TYT AYT"`
	ExamDate time.Time                   `json:"exam_date"  binding:"required"`
	Scores   map[string]SubjectScoresReq `json:"scores"     binding:"required,min=1"`
	Note     *string                     `json:"note"`
}

type ExamResultRes struct {
	ID        uuid.UUID         `json:"id"`
	ExamType  string            `json:"exam_type"`
	ExamDate  time.Time         `json:"exam_date"`
	Scores    entity.ExamScores `json:"scores"`
	TotalNet  float64           `json:"total_net"`
	Note      *string           `json:"note,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

type ExamStatsRes struct {
	ExamType        string             `json:"exam_type"`
	TotalExams      int                `json:"total_exams"`
	AverageTotalNet float64            `json:"average_total_net"`
	BestTotalNet    float64            `json:"best_total_net"`
	SubjectAverages map[string]float64 `json:"subject_averages"`
	Trend           []TrendPoint       `json:"trend"` // tarih sıralı net listesi
}

type TrendPoint struct {
	Date     time.Time `json:"date"`
	TotalNet float64   `json:"total_net"`
}
