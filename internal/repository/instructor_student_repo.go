package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
)

type instructorStudentRepo struct{ db *gorm.DB }

func NewInstructorStudentRepository(db *gorm.DB) InstructorStudentRepository {
	return &instructorStudentRepo{db: db}
}

func (r *instructorStudentRepo) Add(ctx context.Context, instructorID, studentID uuid.UUID) error {
	rel := entity.InstructorStudent{
		InstructorID: instructorID,
		StudentID:    studentID,
	}
	if err := r.db.WithContext(ctx).Create(&rel).Error; err != nil {
		if isDuplicateKeyError(err) {
			return apperror.NewConflict("bu öğrenci zaten ekli", err)
		}
		return apperror.NewInternal(err)
	}
	return nil
}

func (r *instructorStudentRepo) Remove(ctx context.Context, instructorID, studentID uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Where("instructor_id = ? AND student_id = ?", instructorID, studentID).
		Delete(&entity.InstructorStudent{})
	if res.RowsAffected == 0 {
		return apperror.NewNotFound("öğrenci ilişkisi", nil)
	}
	return res.Error
}

func (r *instructorStudentRepo) ListStudents(ctx context.Context, instructorID uuid.UUID, offset, limit int) ([]entity.User, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&entity.InstructorStudent{}).
		Where("instructor_id = ?", instructorID).Count(&total)

	var students []entity.User
	err := r.db.WithContext(ctx).
		Joins("JOIN instructor_students ON instructor_students.student_id = users.id").
		Where("instructor_students.instructor_id = ?", instructorID).
		Where("users.deleted_at IS NULL").
		Order("instructor_students.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&students).Error

	return students, total, err
}

func (r *instructorStudentRepo) ListInstructors(ctx context.Context, studentID uuid.UUID) ([]entity.User, error) {
	var instructors []entity.User
	err := r.db.WithContext(ctx).
		Joins("JOIN instructor_students ON instructor_students.instructor_id = users.id").
		Where("instructor_students.student_id = ?", studentID).
		Where("users.deleted_at IS NULL").
		Find(&instructors).Error
	return instructors, err
}

func (r *instructorStudentRepo) Exists(ctx context.Context, instructorID, studentID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.InstructorStudent{}).
		Where("instructor_id = ? AND student_id = ?", instructorID, studentID).
		Count(&count).Error
	return count > 0, err
}
