package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Rezann47/YksKoc/internal/domain/entity"
)

type StreakRepository struct {
	db *gorm.DB
}

func NewStreakRepository(db *gorm.DB) *StreakRepository {
	return &StreakRepository{db: db}
}

// GetStreakInfo — kullanıcının streak kolonlarını çeker
func (r *StreakRepository) GetStreakInfo(ctx context.Context, userID uuid.UUID) (currentStreak, longestStreak int, lastStudyDate *time.Time, err error) {
	var result struct {
		CurrentStreak int        `gorm:"column:current_streak"`
		LongestStreak int        `gorm:"column:longest_streak"`
		LastStudyDate *time.Time `gorm:"column:last_study_date"`
	}
	err = r.db.WithContext(ctx).
		Table("users").
		Select("current_streak, longest_streak, last_study_date").
		Where("id = ? AND deleted_at IS NULL", userID).
		Scan(&result).Error
	return result.CurrentStreak, result.LongestStreak, result.LastStudyDate, err
}

// UpdateStreak — streak ve last_study_date günceller
func (r *StreakRepository) UpdateStreak(ctx context.Context, tx *gorm.DB, userID uuid.UUID, currentStreak, longestStreak int, lastStudyDate time.Time) error {
	return tx.WithContext(ctx).
		Table("users").
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"current_streak":  currentStreak,
			"longest_streak":  longestStreak,
			"last_study_date": lastStudyDate,
		}).Error
}

// BadgeExists — bu rozeti daha önce kazanmış mı?
func (r *StreakRepository) BadgeExists(ctx context.Context, tx *gorm.DB, userID uuid.UUID, badgeKey string) (bool, error) {
	var count int64
	err := tx.WithContext(ctx).
		Model(&entity.Badge{}).
		Where("user_id = ? AND badge_key = ?", userID, badgeKey).
		Count(&count).Error
	return count > 0, err
}

// InsertBadge — yeni rozet ekler
func (r *StreakRepository) InsertBadge(ctx context.Context, tx *gorm.DB, userID uuid.UUID, def entity.BadgeDef) (*entity.Badge, error) {
	badge := &entity.Badge{
		UserID:    userID,
		BadgeKey:  def.Key,
		BadgeName: def.Name,
		BadgeIcon: def.Icon,
	}
	err := tx.WithContext(ctx).Create(badge).Error
	return badge, err
}

// GetUserBadges — kullanıcının tüm rozetleri
func (r *StreakRepository) GetUserBadges(ctx context.Context, userID uuid.UUID) ([]entity.Badge, error) {
	var badges []entity.Badge
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("earned_at DESC").
		Find(&badges).Error
	return badges, err
}

// GetLeaderboard — haftalık sıralama (view'dan)
func (r *StreakRepository) GetLeaderboard(ctx context.Context, limit int) ([]entity.LeaderboardEntry, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var entries []entity.LeaderboardEntry
	err := r.db.WithContext(ctx).
		Limit(limit).
		Find(&entries).Error
	return entries, err
}

// GetMyRank — kullanıcının sırası
func (r *StreakRepository) GetMyRank(ctx context.Context, userID uuid.UUID) int {
	var rank int
	r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(rank, 0) FROM (
			SELECT id, RANK() OVER (ORDER BY total_minutes DESC) AS rank
			FROM leaderboard_weekly
		) ranked WHERE id = ?
	`, userID).Scan(&rank)
	return rank
}
