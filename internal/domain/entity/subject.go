package entity

import (
	"time"

	"github.com/google/uuid"
)

type ExamType string

const (
	ExamTypeTYT ExamType = "TYT"
	ExamTypeAYT ExamType = "AYT"
)

type Subject struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name         string    `gorm:"type:varchar(100);not null"`
	ExamType     ExamType  `gorm:"type:exam_type;not null;index"`
	DisplayOrder int16     `gorm:"not null;default:0"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime"`

	Topics []Topic `gorm:"foreignKey:SubjectID"`
}

type Topic struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SubjectID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Name         string    `gorm:"type:varchar(200);not null"`
	DisplayOrder int16     `gorm:"not null;default:0"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime"`

	Subject Subject `gorm:"foreignKey:SubjectID"`
}

type StudentTopicProgress struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index"`
	TopicID        uuid.UUID `gorm:"type:uuid;not null"`
	CompletionDate time.Time `gorm:"not null;default:now()"`
	CreatedAt      time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"not null;autoUpdateTime"`

	User  User  `gorm:"foreignKey:UserID"`
	Topic Topic `gorm:"foreignKey:TopicID"`
}

func (StudentTopicProgress) TableName() string {
	return "student_topic_progress"
}
