package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePomodoroReq struct {
	DurationMinutes int16      `json:"duration_minutes" binding:"required,min=1,max=480"`
	SubjectID       *uuid.UUID `json:"subject_id"`
	StartedAt       *time.Time `json:"started_at"` // nil ise now() kullanılır
}

type PomodoroRes struct {
	ID              uuid.UUID  `json:"id"`
	DurationMinutes int16      `json:"duration_minutes"`
	SubjectID       *uuid.UUID `json:"subject_id,omitempty"`
	SubjectName     *string    `json:"subject_name,omitempty"`
	StartedAt       time.Time  `json:"started_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type PomodoroListFilter struct {
	From *time.Time `form:"from"`
	To   *time.Time `form:"to"`
	PaginationReq
}

// PomodoroStatsRes günlük/haftalık özet
type PomodoroStatsRes struct {
	TotalMinutes    int            `json:"total_minutes"`
	TotalSessions   int            `json:"total_sessions"`
	DailyBreakdown  []DailyStats   `json:"daily_breakdown"`
}

type DailyStats struct {
	Date         string `json:"date"` // YYYY-MM-DD
	TotalMinutes int    `json:"total_minutes"`
	Sessions     int    `json:"sessions"`
}
