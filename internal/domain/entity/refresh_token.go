package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	TokenHash string     `gorm:"type:char(64);not null;uniqueIndex"`
	UserAgent *string    `gorm:"type:varchar(500)"`
	IPAddress *string    `gorm:"type:varchar(45)"`
	ExpiresAt time.Time  `gorm:"not null"`
	RevokedAt *time.Time
	CreatedAt time.Time  `gorm:"not null;autoCreateTime"`

	User User `gorm:"foreignKey:UserID"`
}

func (r *RefreshToken) IsExpired() bool  { return time.Now().After(r.ExpiresAt) }
func (r *RefreshToken) IsRevoked() bool  { return r.RevokedAt != nil }
func (r *RefreshToken) IsValid() bool    { return !r.IsExpired() && !r.IsRevoked() }
