package dto

import (
	"time"

	"github.com/google/uuid"
)

type AddStudentReq struct {
	StudentCode string `json:"student_code" binding:"required,min=4,max=20"`
}
type InstructorRes struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	AvatarID   int16      `json:"avatar_id"`
	IsOnline   bool       `json:"is_online"`
	LastSeenAt *time.Time `json:"last_seen_at"`
}

type StudentListItemRes struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	StudentCode string     `json:"student_code"`
	AvatarID    int16      `json:"avatar_id"`
	IsOnline    bool       `json:"is_online"`
	LastSeenAt  *time.Time `json:"last_seen_at"`
	AddedAt     time.Time  `json:"added_at"`
}

type RemoveStudentReq struct {
	StudentID uuid.UUID `json:"student_id" binding:"required"`
}
