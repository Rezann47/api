package service

import (
	"context"
	"time"

	"github.com/Rezann47/YksKoc/internal/domain/dto"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterReq) (*dto.LoginRes, error)
	Login(ctx context.Context, req dto.LoginReq, userAgent, ip string) (*dto.LoginRes, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.TokenRes, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID uuid.UUID) error
}

type UserService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserRes, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileReq) (*dto.UserRes, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordReq) error
	GetPremiumStatus(ctx context.Context, userID uuid.UUID) (*dto.PremiumStatusRes, error)
	ActivatePremium(ctx context.Context, userID uuid.UUID) (*dto.PremiumStatusRes, error)
	Ping(ctx context.Context, userID uuid.UUID) error
}

type SubjectService interface {
	ListSubjects(ctx context.Context, examType *entity.ExamType) ([]dto.SubjectRes, error)
	ListTopics(ctx context.Context, subjectID uuid.UUID, userID uuid.UUID) ([]dto.TopicRes, error)
	MarkTopic(ctx context.Context, userID, topicID uuid.UUID, isCompleted bool) error
	GetSubjectProgress(ctx context.Context, userID, subjectID uuid.UUID) (*dto.SubjectProgressRes, error)
	GetAllProgress(ctx context.Context, userID uuid.UUID) ([]dto.SubjectProgressRes, error)
}

type PomodoroService interface {
	Create(ctx context.Context, userID uuid.UUID, req dto.CreatePomodoroReq) (*dto.PomodoroRes, error)
	List(ctx context.Context, userID uuid.UUID, filter dto.PomodoroListFilter) (*dto.PaginatedRes[dto.PomodoroRes], error)
	GetStats(ctx context.Context, userID uuid.UUID, from, to time.Time) (*dto.PomodoroStatsRes, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

type ExamResultService interface {
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateExamResultReq) (*dto.ExamResultRes, error)
	List(ctx context.Context, userID uuid.UUID, examType *entity.ExamType, pagination dto.PaginationReq) (*dto.PaginatedRes[dto.ExamResultRes], error)
	GetStats(ctx context.Context, userID uuid.UUID, examType entity.ExamType) (*dto.ExamStatsRes, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

type InstructorService interface {
	AddStudent(ctx context.Context, instructorID uuid.UUID, req dto.AddStudentReq) error
	RemoveStudent(ctx context.Context, instructorID, studentID uuid.UUID) error
	ListStudents(ctx context.Context, instructorID uuid.UUID, pagination dto.PaginationReq) (*dto.PaginatedRes[dto.StudentListItemRes], error)
	// Öğrenci verisini görüntüleme (salt okunur)
	GetStudentPomodoros(ctx context.Context, instructorID, studentID uuid.UUID, from, to time.Time) (*dto.PomodoroStatsRes, error)
	GetStudentProgress(ctx context.Context, instructorID, studentID uuid.UUID) ([]dto.SubjectProgressRes, error)
	GetStudentExamResults(ctx context.Context, instructorID, studentID uuid.UUID, examType *entity.ExamType, pagination dto.PaginationReq) (*dto.PaginatedRes[dto.ExamResultRes], error)
	ListMyInstructors(ctx context.Context, studentID uuid.UUID) ([]dto.InstructorRes, error)
}
