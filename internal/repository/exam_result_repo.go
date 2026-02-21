package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
)

type examResultRepo struct{ db *gorm.DB }

func NewExamResultRepository(db *gorm.DB) ExamResultRepository {
	return &examResultRepo{db: db}
}

func (r *examResultRepo) Create(ctx context.Context, result *entity.ExamResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}

func (r *examResultRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.ExamResult, error) {
	var result entity.ExamResult
	err := r.db.WithContext(ctx).First(&result, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFound("deneme sonucu", err)
		}
		return nil, apperror.NewInternal(err)
	}
	return &result, nil
}

func (r *examResultRepo) ListByUser(ctx context.Context, userID uuid.UUID, examType *entity.ExamType, offset, limit int) ([]entity.ExamResult, int64, error) {
	q := r.db.WithContext(ctx).Model(&entity.ExamResult{}).Where("user_id = ?", userID)
	if examType != nil {
		q = q.Where("exam_type = ?", *examType)
	}

	var total int64
	q.Count(&total)

	var results []entity.ExamResult
	err := q.Order("exam_date DESC").Offset(offset).Limit(limit).Find(&results).Error
	return results, total, err
}

func (r *examResultRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&entity.ExamResult{})
	if res.RowsAffected == 0 {
		return apperror.NewNotFound("deneme sonucu", nil)
	}
	return res.Error
}

func (r *examResultRepo) AverageNetByUser(ctx context.Context, userID uuid.UUID, examType entity.ExamType) (map[string]float64, error) {
	// JSONB üzerinden ders bazlı ortalama — PostgreSQL özel sorgu
	type row struct {
		Subject string  `gorm:"column:subject"`
		AvgNet  float64 `gorm:"column:avg_net"`
	}
	var rows []row

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			kv.key                          AS subject,
			AVG((kv.value->>'net')::numeric) AS avg_net
		FROM exam_results,
		     jsonb_each(scores) AS kv
		WHERE user_id  = ?
		  AND exam_type = ?
		GROUP BY kv.key
	`, userID, examType).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]float64, len(rows))
	for _, r := range rows {
		result[r.Subject] = r.AvgNet
	}
	return result, nil
}
