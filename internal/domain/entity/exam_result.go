package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SubjectScores bir dersin doğru/yanlış/net değerleri
type SubjectScores struct {
	Correct int     `json:"correct"`
	Wrong   int     `json:"wrong"`
	Net     float64 `json:"net"`
}

// ExamScores JSONB kolonuna eşlenir
// Örn: {"turkish":{"correct":28,"wrong":4,"net":27.0},"math":{...}}
type ExamScores map[string]SubjectScores

func (e ExamScores) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *ExamScores) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("ExamScores.Scan: expected []byte, got %T", value)
	}
	return json.Unmarshal(b, e)
}

type ExamResult struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	ExamType  ExamType   `gorm:"type:exam_type;not null;index"`
	ExamDate  time.Time  `gorm:"type:date;not null"`
	Scores    ExamScores `gorm:"type:jsonb;not null;default:'{}'"`
	TotalNet  float64    `gorm:"type:numeric(6,2);not null;default:0"`
	Note      *string    `gorm:"type:text"`
	CreatedAt time.Time  `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"not null;autoUpdateTime"`

	User User `gorm:"foreignKey:UserID"`
}
