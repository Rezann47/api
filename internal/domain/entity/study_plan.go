package entity

import (
	"time"

	"github.com/google/uuid"
)

type StudyPlan struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
	Title     string    `gorm:"type:varchar(200);not null;default:'Çalışma Planı'"`
	PlanDate  time.Time `gorm:"type:date;not null;index"`
	Note      *string   `gorm:"type:text"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`

	User    User            `gorm:"foreignKey:UserID"`
	Creator User            `gorm:"foreignKey:CreatedBy"`
	Items   []StudyPlanItem `gorm:"foreignKey:PlanID"`
}

func (StudyPlan) TableName() string { return "study_plans" }

type StudyPlanItem struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PlanID          uuid.UUID  `gorm:"type:uuid;not null;index"`
	SubjectID       uuid.UUID  `gorm:"type:uuid;not null"`
	TopicID         *uuid.UUID `gorm:"type:uuid"`
	DurationMinutes int        `gorm:"not null;default:30"`
	DisplayOrder    int16      `gorm:"not null;default:0"`
	IsCompleted     bool       `gorm:"not null;default:false"`
	CompletedAt     *time.Time
	CreatedAt       time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"not null;autoUpdateTime"`

	Plan    StudyPlan `gorm:"foreignKey:PlanID"`
	Subject Subject   `gorm:"foreignKey:SubjectID"`
	Topic   *Topic    `gorm:"foreignKey:TopicID"`
}

func (StudyPlanItem) TableName() string { return "study_plan_items" }
