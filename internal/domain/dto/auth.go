package dto

import "github.com/google/uuid"

type RegisterReq struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Role     string `json:"role"     binding:"required,oneof=student instructor"`
}

type LoginReq struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"` // Bearer
	ExpiresIn    int64  `json:"expires_in"` // saniye
}

type AuthUserRes struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	StudentCode *string   `json:"student_code,omitempty"`
	IsPremium   bool      `json:"is_premium"`
	AvatarID    int16     `json:"avatar_id"`
}

type LoginRes struct {
	User  AuthUserRes `json:"user"`
	Token TokenRes    `json:"token"`
}
