package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleStudent    Role = "student"
	RoleInstructor Role = "instructor"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name         string         `gorm:"type:varchar(100);not null"`
	Email        string         `gorm:"type:citext;not null"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`
	Role         Role           `gorm:"type:user_role;not null;default:'student'"`
	StudentCode  *string        `gorm:"type:varchar(20)"`
	IsPremium    bool           `gorm:"not null;default:false"`
	IsActive     bool           `gorm:"not null;default:true"`
	CreatedAt    time.Time      `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"not null;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
