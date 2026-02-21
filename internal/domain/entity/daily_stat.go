package entity

// DailyStat pomodoro repository ham sorgu sonucu için yardımcı struct.
// GORM modeli değil — raw SQL scan için kullanılır.
type DailyStat struct {
	Date         string `gorm:"column:date"`
	TotalMinutes int    `gorm:"column:total_minutes"`
	Sessions     int    `gorm:"column:sessions"`
}
