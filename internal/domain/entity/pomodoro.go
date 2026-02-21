package entity

import (
	"time"

	"github.com/google/uuid"
)

type Pomodoro struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index"`
	SubjectID       *uuid.UUID `gorm:"type:uuid"`
	DurationMinutes int16      `gorm:"not null"`
	StartedAt       time.Time  `gorm:"not null;default:now()"`
	CreatedAt       time.Time  `gorm:"not null;autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"not null;autoUpdateTime"`

	User    User     `gorm:"foreignKey:UserID"`
	Subject *Subject `gorm:"foreignKey:SubjectID"`
}
