package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
)

type pomodoroRepo struct{ db *gorm.DB }

func NewPomodoroRepository(db *gorm.DB) PomodoroRepository {
	return &pomodoroRepo{db: db}
}

func (r *pomodoroRepo) Create(ctx context.Context, p *entity.Pomodoro) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *pomodoroRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Pomodoro, error) {
	var p entity.Pomodoro
	err := r.db.WithContext(ctx).Preload("Subject").First(&p, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFound("pomodoro", err)
		}
		return nil, apperror.NewInternal(err)
	}
	return &p, nil
}

func (r *pomodoroRepo) ListByUser(ctx context.Context, userID uuid.UUID, from, to *time.Time, offset, limit int) ([]entity.Pomodoro, int64, error) {
	q := r.db.WithContext(ctx).Model(&entity.Pomodoro{}).
		Where("user_id = ?", userID)
	if from != nil {
		q = q.Where("started_at >= ?", from)
	}
	if to != nil {
		q = q.Where("started_at <= ?", to)
	}

	var total int64
	q.Count(&total)

	var pomodoros []entity.Pomodoro
	err := q.Preload("Subject").
		Order("started_at DESC").
		Offset(offset).Limit(limit).
		Find(&pomodoros).Error

	return pomodoros, total, err
}

func (r *pomodoroRepo) SumMinutesByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&entity.Pomodoro{}).
		Where("user_id = ? AND started_at BETWEEN ? AND ?", userID, from, to).
		Select("COALESCE(SUM(duration_minutes), 0)").
		Scan(&total).Error
	return total, err
}

func (r *pomodoroRepo) DailyStatsByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]entity.DailyStat, error) {
	var stats []entity.DailyStat
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			DATE(started_at AT TIME ZONE 'Europe/Istanbul') AS date,
			COALESCE(SUM(duration_minutes), 0)             AS total_minutes,
			COUNT(*)                                        AS sessions
		FROM pomodoros
		WHERE user_id = ?
		  AND started_at BETWEEN ? AND ?
		GROUP BY DATE(started_at AT TIME ZONE 'Europe/Istanbul')
		ORDER BY date
	`, userID, from, to).Scan(&stats).Error
	return stats, err
}

func (r *pomodoroRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&entity.Pomodoro{})
	if result.RowsAffected == 0 {
		return apperror.NewNotFound("pomodoro", nil)
	}
	return result.Error
}
