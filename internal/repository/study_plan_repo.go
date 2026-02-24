package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Rezann47/YksKoc/internal/domain/entity"
)

type StudyPlanRepository interface {
	Create(ctx context.Context, plan *entity.StudyPlan) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.StudyPlan, error)
	ListByUserAndDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*entity.StudyPlan, error)
	ListByUserAndMonth(ctx context.Context, userID uuid.UUID, year, month int) ([]*entity.StudyPlan, error)
	Update(ctx context.Context, plan *entity.StudyPlan) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
	CompleteItem(ctx context.Context, itemID, planID uuid.UUID) error
	UncompleteItem(ctx context.Context, itemID, planID uuid.UUID) error
}

type studyPlanRepo struct{ db *gorm.DB }

func NewStudyPlanRepository(db *gorm.DB) StudyPlanRepository {
	return &studyPlanRepo{db: db}
}

func (r *studyPlanRepo) Create(ctx context.Context, plan *entity.StudyPlan) error {
	return r.db.WithContext(ctx).
		Create(plan).Error
}

func (r *studyPlanRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.StudyPlan, error) {
	var plan entity.StudyPlan
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Subject").
		Preload("Items.Topic").
		Preload("Creator").
		Where("id = ?", id).
		First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *studyPlanRepo) ListByUserAndDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*entity.StudyPlan, error) {
	var plans []*entity.StudyPlan
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Subject").
		Preload("Items.Topic").
		Preload("Creator").
		Where("user_id = ? AND plan_date = ?", userID, date.Format("2006-01-02")).
		Order("created_at ASC").
		Find(&plans).Error
	return plans, err
}

func (r *studyPlanRepo) ListByUserAndMonth(ctx context.Context, userID uuid.UUID, year, month int) ([]*entity.StudyPlan, error) {
	var plans []*entity.StudyPlan
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Subject").
		Where("user_id = ? AND plan_date >= ? AND plan_date < ?", userID, start, end).
		Order("plan_date ASC").
		Find(&plans).Error
	return plans, err
}

func (r *studyPlanRepo) Update(ctx context.Context, plan *entity.StudyPlan) error {
	return r.db.WithContext(ctx).Save(plan).Error
}

func (r *studyPlanRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&entity.StudyPlan{}).Error
}

func (r *studyPlanRepo) CompleteItem(ctx context.Context, itemID, planID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.StudyPlanItem{}).
		Where("id = ? AND plan_id = ?", itemID, planID).
		Updates(map[string]interface{}{
			"is_completed": true,
			"completed_at": now,
		}).Error
}

func (r *studyPlanRepo) UncompleteItem(ctx context.Context, itemID, planID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.StudyPlanItem{}).
		Where("id = ? AND plan_id = ?", itemID, planID).
		Updates(map[string]interface{}{
			"is_completed": false,
			"completed_at": nil,
		}).Error
}
