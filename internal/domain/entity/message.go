package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SenderID   uuid.UUID `gorm:"type:uuid;not null;index"`
	ReceiverID uuid.UUID `gorm:"type:uuid;not null;index"`
	Content    string    `gorm:"type:text;not null"`
	IsRead     bool      `gorm:"not null;default:false"`
	ReadAt     *time.Time
	CreatedAt  time.Time `gorm:"not null;autoCreateTime"`

	Sender   *User `gorm:"foreignKey:SenderID"`
	Receiver *User `gorm:"foreignKey:ReceiverID"`
}
