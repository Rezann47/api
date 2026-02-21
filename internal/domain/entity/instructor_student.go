package entity

import (
	"time"

	"github.com/google/uuid"
)

type InstructorStudent struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	InstructorID uuid.UUID `gorm:"type:uuid;not null;index"`
	StudentID    uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`

	Instructor User `gorm:"foreignKey:InstructorID"`
	Student    User `gorm:"foreignKey:StudentID"`
}
