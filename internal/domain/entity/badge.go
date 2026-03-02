package entity

import (
	"time"

	"github.com/google/uuid"
)

// Badge — kazanılan rozet kaydı (badges tablosu)
type Badge struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"                       json:"user_id"`
	BadgeKey  string    `gorm:"size:50;not null"                               json:"badge_key"`
	BadgeName string    `gorm:"size:100;not null"                              json:"badge_name"`
	BadgeIcon string    `gorm:"size:10;not null"                               json:"badge_icon"`
	EarnedAt  time.Time `gorm:"not null;default:now()"                         json:"earned_at"`
}

func (Badge) TableName() string { return "badges" }

// LeaderboardEntry — leaderboard_weekly view'ından gelen satır
// (GORM ile view okuma için model)
type LeaderboardEntry struct {
	ID            uuid.UUID `gorm:"type:uuid" json:"id"`
	FullName      string    `json:"full_name"`
	AvatarID      string    `json:"avatar_id"`
	CurrentStreak int       `json:"current_streak"`
	TotalMinutes  int       `json:"total_minutes"`
	SessionCount  int       `json:"session_count"`
}

func (LeaderboardEntry) TableName() string { return "leaderboard_weekly" }

// BadgeDef — statik rozet tanımı
type BadgeDef struct {
	Key       string
	Name      string
	Icon      string
	StreakDay int
}

// AllBadges — kazanılabilecek tüm rozetler
var AllBadges = []BadgeDef{
	{Key: "streak_1", Name: "İlk Adım", Icon: "🌱", StreakDay: 1},
	{Key: "streak_3", Name: "3 Gün Serisi", Icon: "🔥", StreakDay: 3},
	{Key: "streak_7", Name: "Hafta Şampiyonu", Icon: "⚡", StreakDay: 7},
	{Key: "streak_14", Name: "2 Hafta Fırtınası", Icon: "🌪️", StreakDay: 14},
	{Key: "streak_30", Name: "Aylık Efsane", Icon: "👑", StreakDay: 30},
	{Key: "streak_60", Name: "Demir İrade", Icon: "💎", StreakDay: 60},
	{Key: "streak_100", Name: "Efsane", Icon: "🏆", StreakDay: 100},
}

// StreakUpdateResult — UpdateStreak'ten frontend'e dönen sonuç
type StreakUpdateResult struct {
	CurrentStreak int     `json:"current_streak"`
	LongestStreak int     `json:"longest_streak"`
	NewBadges     []Badge `json:"new_badges"`
}
