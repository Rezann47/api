package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/Rezann47/YksKoc/internal/repository"
)

type StreakService struct {
	db   *gorm.DB
	repo *repository.StreakRepository
}

func NewStreakService(db *gorm.DB) *StreakService {
	return &StreakService{
		db:   db,
		repo: repository.NewStreakRepository(db),
	}
}

// UpdateStreak — pomodoro tamamlanınca çağrılır.
// Günlük seriyi günceller, yeni rozet kazanıldıysa kaydeder.
func (s *StreakService) UpdateStreak(ctx context.Context, userID uuid.UUID) (*entity.StreakUpdateResult, error) {
	// Mevcut streak bilgisi
	currentStreak, longestStreak, lastStudyDate, err := s.repo.GetStreakInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	newStreak := currentStreak

	if lastStudyDate == nil {
		newStreak = 1
	} else {
		lastDay := lastStudyDate.UTC().Truncate(24 * time.Hour)
		diffDays := int(today.Sub(lastDay).Hours() / 24)

		switch {
		case diffDays == 0:
			// Bugün zaten çalışmış, streak değişmez
		case diffDays == 1:
			// Dün çalışmış → seri devam
			newStreak = currentStreak + 1
		default:
			// Seri koptu → sıfırla
			newStreak = 1
		}
	}

	if newStreak > longestStreak {
		longestStreak = newStreak
	}

	var newBadges []entity.Badge

	// Transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Streak güncelle
		if err := s.repo.UpdateStreak(ctx, tx, userID, newStreak, longestStreak, today); err != nil {
			return err
		}

		// Yeni rozetleri kontrol et
		for _, def := range entity.AllBadges {
			if newStreak < def.StreakDay {
				continue
			}
			exists, err := s.repo.BadgeExists(ctx, tx, userID, def.Key)
			if err != nil || exists {
				continue
			}
			badge, err := s.repo.InsertBadge(ctx, tx, userID, def)
			if err != nil {
				continue
			}
			newBadges = append(newBadges, *badge)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &entity.StreakUpdateResult{
		CurrentStreak: newStreak,
		LongestStreak: longestStreak,
		NewBadges:     newBadges,
	}, nil
}

// GetMyStreak — streak bilgisi
func (s *StreakService) GetMyStreak(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	currentStreak, longestStreak, lastStudyDate, err := s.repo.GetStreakInfo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"current_streak":  currentStreak,
		"longest_streak":  longestStreak,
		"last_study_date": lastStudyDate,
	}, nil
}

// GetMyBadges — kazanılan rozetler
func (s *StreakService) GetMyBadges(ctx context.Context, userID uuid.UUID) ([]entity.Badge, error) {
	return s.repo.GetUserBadges(ctx, userID)
}

// GetLeaderboard — haftalık sıralama + kendi sıran
func (s *StreakService) GetLeaderboard(ctx context.Context, userID uuid.UUID) ([]entity.LeaderboardEntry, int, error) {
	entries, err := s.repo.GetLeaderboard(ctx, 50)
	if err != nil {
		return nil, 0, err
	}
	myRank := s.repo.GetMyRank(ctx, userID)
	return entries, myRank, nil
}
