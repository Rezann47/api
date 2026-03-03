package repository

import (
	"context"
	"time"

	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/google/uuid"
)

// ─── UserRepository ───────────────────────────────────────

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByStudentCode(ctx context.Context, code string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateLastSeen(ctx context.Context, id uuid.UUID, t time.Time) error
	UpdatePremium(ctx context.Context, id uuid.UUID, isPremium bool, expiresAt *time.Time, txID string) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

// ─── RefreshTokenRepository ───────────────────────────────

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*entity.RefreshToken, error)
	RevokeByHash(ctx context.Context, hash string) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

// ─── SubjectRepository ────────────────────────────────────

type SubjectRepository interface {
	FindAll(ctx context.Context) ([]entity.Subject, error)
	FindByExamType(ctx context.Context, examType entity.ExamType) ([]entity.Subject, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Subject, error)
}

// ─── TopicRepository ──────────────────────────────────────

type TopicRepository interface {
	FindBySubjectID(ctx context.Context, subjectID uuid.UUID) ([]entity.Topic, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Topic, error)
}

// ─── ProgressRepository ───────────────────────────────────

type ProgressRepository interface {
	// Upsert: var olan kaydı günceller, yoksa ekler
	Upsert(ctx context.Context, userID, topicID uuid.UUID) error
	// Delete: konu işaretini kaldırır (tamamlanmadı)
	Delete(ctx context.Context, userID, topicID uuid.UUID) error
	// FindCompletedByUserAndSubject: bir dersteki tamamlanan topic ID'leri
	FindCompletedByUserAndSubject(ctx context.Context, userID, subjectID uuid.UUID) ([]uuid.UUID, error)
	// FindAllCompletedByUser: tüm tamamlanan topic ID'leri
	FindAllCompletedByUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	// CountBySubject: (total, completed) döner
	CountBySubject(ctx context.Context, userID, subjectID uuid.UUID) (total int64, completed int64, err error)
}

// ─── PomodoroRepository ───────────────────────────────────

type PomodoroRepository interface {
	Create(ctx context.Context, pomodoro *entity.Pomodoro) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Pomodoro, error)
	ListByUser(ctx context.Context, userID uuid.UUID, from, to *time.Time, offset, limit int) ([]entity.Pomodoro, int64, error)
	SumMinutesByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) (int64, error)
	DailyStatsByUser(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]entity.DailyStat, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

// ─── ExamResultRepository ─────────────────────────────────

type ExamResultRepository interface {
	Create(ctx context.Context, result *entity.ExamResult) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ExamResult, error)
	ListByUser(ctx context.Context, userID uuid.UUID, examType *entity.ExamType, offset, limit int) ([]entity.ExamResult, int64, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	// AverageNetByUser: ders bazlı ortalama + genel ortalama
	AverageNetByUser(ctx context.Context, userID uuid.UUID, examType entity.ExamType) (map[string]float64, error)
}

// ─── InstructorStudentRepository ─────────────────────────

type InstructorStudentRepository interface {
	Add(ctx context.Context, instructorID, studentID uuid.UUID) error
	Remove(ctx context.Context, instructorID, studentID uuid.UUID) error
	ListStudents(ctx context.Context, instructorID uuid.UUID, offset, limit int) ([]entity.User, int64, error)
	ListInstructors(ctx context.Context, studentID uuid.UUID) ([]entity.User, error)
	Exists(ctx context.Context, instructorID, studentID uuid.UUID) (bool, error)
}
