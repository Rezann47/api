package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
)

// ─── Subject ──────────────────────────────────────────────

type subjectRepo struct{ db *gorm.DB }

func NewSubjectRepository(db *gorm.DB) SubjectRepository {
	return &subjectRepo{db: db}
}

func (r *subjectRepo) FindAll(ctx context.Context) ([]entity.Subject, error) {
	var subjects []entity.Subject
	err := r.db.WithContext(ctx).Order("exam_type, display_order").Find(&subjects).Error
	return subjects, err
}

func (r *subjectRepo) FindByExamType(ctx context.Context, examType entity.ExamType) ([]entity.Subject, error) {
	var subjects []entity.Subject
	err := r.db.WithContext(ctx).
		Where("exam_type = ?", examType).
		Order("display_order").
		Find(&subjects).Error
	return subjects, err
}

func (r *subjectRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Subject, error) {
	var subject entity.Subject
	err := r.db.WithContext(ctx).First(&subject, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFound("ders", err)
		}
		return nil, apperror.NewInternal(err)
	}
	return &subject, nil
}

// ─── Topic ────────────────────────────────────────────────

type topicRepo struct{ db *gorm.DB }

func NewTopicRepository(db *gorm.DB) TopicRepository {
	return &topicRepo{db: db}
}

func (r *topicRepo) FindBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]entity.Topic, error) {
	var topics []entity.Topic
	err := r.db.WithContext(ctx).
		Where("subject_id = ?", subjectID).
		Order("display_order").
		Find(&topics).Error
	return topics, err
}

func (r *topicRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Topic, error) {
	var topic entity.Topic
	err := r.db.WithContext(ctx).First(&topic, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFound("konu", err)
		}
		return nil, apperror.NewInternal(err)
	}
	return &topic, nil
}

// ─── Progress ─────────────────────────────────────────────

type progressRepo struct{ db *gorm.DB }

func NewProgressRepository(db *gorm.DB) ProgressRepository {
	return &progressRepo{db: db}
}

func (r *progressRepo) Upsert(ctx context.Context, userID, topicID uuid.UUID) error {
	progress := entity.StudentTopicProgress{
		UserID:  userID,
		TopicID: topicID,
	}
	// ON CONFLICT DO NOTHING — zaten tamamlanmışsa tekrar ekleme
	return r.db.WithContext(ctx).
		Where(entity.StudentTopicProgress{UserID: userID, TopicID: topicID}).
		FirstOrCreate(&progress).Error
}

func (r *progressRepo) Delete(ctx context.Context, userID, topicID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND topic_id = ?", userID, topicID).
		Delete(&entity.StudentTopicProgress{}).Error
}

func (r *progressRepo) FindCompletedByUserAndSubject(ctx context.Context, userID, subjectID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := r.db.WithContext(ctx).
		Model(&entity.StudentTopicProgress{}).
		Joins("JOIN topics ON topics.id = student_topic_progress.topic_id").
		Where("student_topic_progress.user_id = ? AND topics.subject_id = ?", userID, subjectID).
		Pluck("student_topic_progress.topic_id", &ids).Error
	return ids, err
}

func (r *progressRepo) FindAllCompletedByUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := r.db.WithContext(ctx).
		Model(&entity.StudentTopicProgress{}).
		Where("user_id = ?", userID).
		Pluck("topic_id", &ids).Error
	return ids, err
}

func (r *progressRepo) CountBySubject(ctx context.Context, userID, subjectID uuid.UUID) (total int64, completed int64, err error) {
	// Toplam topic sayısı
	r.db.WithContext(ctx).Model(&entity.Topic{}).
		Where("subject_id = ?", subjectID).Count(&total)

	// Tamamlanan topic sayısı
	r.db.WithContext(ctx).Model(&entity.StudentTopicProgress{}).
		Joins("JOIN topics ON topics.id = student_topic_progress.topic_id").
		Where("student_topic_progress.user_id = ? AND topics.subject_id = ?", userID, subjectID).
		Count(&completed)

	return total, completed, nil
}
