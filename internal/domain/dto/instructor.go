package dto

import (
	"time"

	"github.com/google/uuid"
)

type AddStudentReq struct {
	StudentCode string `json:"student_code" binding:"required,min=4,max=20"`
}

type StudentListItemRes struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	StudentCode string    `json:"student_code"`
	AddedAt     time.Time `json:"added_at"`
}

type RemoveStudentReq struct {
	StudentID uuid.UUID `json:"student_id" binding:"required"`
}
