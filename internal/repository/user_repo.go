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

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isDuplicateKeyError(err) {
			return apperror.NewConflict("bu e-posta adresi zaten kullanılıyor", err)
		}
		return apperror.NewInternal(err)
	}
	return nil
}

func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFound("kullanıcı", err)
		}
		return nil, apperror.NewInternal(err)
	}
	return &user, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFound("kullanıcı", err)
		}
		return nil, apperror.NewInternal(err)
	}
	return &user, nil
}

func (r *userRepo) FindByStudentCode(ctx context.Context, code string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Where("student_code = ? AND role = 'student'", code).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFound("öğrenci", err)
		}
		return nil, apperror.NewInternal(err)
	}
	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return apperror.NewInternal(err)
	}
	return nil
}

func (r *userRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&entity.User{}, id).Error; err != nil {
		return apperror.NewInternal(err)
	}
	return nil
}

func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, apperror.NewInternal(err)
	}
	return count > 0, nil
}

// 🔹 Ping için kullanılan method
func (r *userRepo) UpdateLastSeen(ctx context.Context, id uuid.UUID, t time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		UpdateColumn("last_seen_at", t).
		Error
}

// --- helpers ---

func isDuplicateKeyError(err error) bool {
	return err != nil && (containsString(err.Error(), "duplicate key") ||
		containsString(err.Error(), "unique constraint"))
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
