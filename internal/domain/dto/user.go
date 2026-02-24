package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserRes struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	StudentCode *string    `json:"student_code,omitempty"`
	IsPremium   bool       `json:"is_premium"`
	AvatarID    int16      `json:"avatar_id"`
	LastSeenAt  *time.Time `json:"last_seen_at"`
	IsOnline    bool       `json:"is_online"`
	CreatedAt   time.Time  `json:"created_at"`
}

type UpdateProfileReq struct {
	Name     string `json:"name"      binding:"omitempty,min=2,max=100"`
	AvatarID *int16 `json:"avatar_id" binding:"omitempty,min=1,max=100"`
}

type ChangePasswordReq struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password"     binding:"required,min=8,max=72"`
}

type PremiumStatusRes struct {
	IsPremium bool `json:"is_premium"`
}
