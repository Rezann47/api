package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/Rezann47/YksKoc/internal/domain/dto"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/Rezann47/YksKoc/internal/repository"
)

type instructorService struct {
	relRepo      repository.InstructorStudentRepository
	userRepo     repository.UserRepository
	pomodoroRepo repository.PomodoroRepository
	progressRepo repository.ProgressRepository
	subjectRepo  repository.SubjectRepository
	examRepo     repository.ExamResultRepository
	log          *zap.Logger
}

func NewInstructorService(
	relRepo repository.InstructorStudentRepository,
	userRepo repository.UserRepository,
	pomodoroRepo repository.PomodoroRepository,
	progressRepo repository.ProgressRepository,
	subjectRepo repository.SubjectRepository,
	examRepo repository.ExamResultRepository,
	log *zap.Logger,
) InstructorService {
	return &instructorService{
		relRepo:      relRepo,
		userRepo:     userRepo,
		pomodoroRepo: pomodoroRepo,
		progressRepo: progressRepo,
		subjectRepo:  subjectRepo,
		examRepo:     examRepo,
		log:          log,
	}
}

func (s *instructorService) AddStudent(ctx context.Context, instructorID uuid.UUID, req dto.AddStudentReq) error {
	student, err := s.userRepo.FindByStudentCode(ctx, req.StudentCode)
	if err != nil {
		return apperror.NewNotFound("öğrenci kodu bulunamadı", nil)
	}
	if student.Role != entity.RoleStudent {
		return apperror.NewValidation("bu kod bir öğrenciye ait değil")
	}
	return s.relRepo.Add(ctx, instructorID, student.ID)
}

func (s *instructorService) RemoveStudent(ctx context.Context, instructorID, studentID uuid.UUID) error {
	return s.relRepo.Remove(ctx, instructorID, studentID)
}

func (s *instructorService) ListStudents(ctx context.Context, instructorID uuid.UUID, pagination dto.PaginationReq) (*dto.PaginatedRes[dto.StudentListItemRes], error) {
	students, total, err := s.relRepo.ListStudents(ctx, instructorID, pagination.Offset(), pagination.Limit)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	res := make([]dto.StudentListItemRes, len(students))
	for i, st := range students {
		code := ""
		if st.StudentCode != nil {
			code = *st.StudentCode
		}
		res[i] = dto.StudentListItemRes{
			ID:          st.ID,
			Name:        st.Name,
			Email:       st.Email,
			StudentCode: code,
			AddedAt:     st.CreatedAt,
		}
	}

	paged := dto.NewPaginatedRes(res, total, pagination.Page, pagination.Limit)
	return &paged, nil
}

// mustBeStudent eğitmenin bu öğrenciye erişim hakkı olup olmadığını doğrular
func (s *instructorService) mustBeStudent(ctx context.Context, instructorID, studentID uuid.UUID) error {
	ok, err := s.relRepo.Exists(ctx, instructorID, studentID)
	if err != nil {
		return apperror.NewInternal(err)
	}
	if !ok {
		return apperror.NewForbidden("bu öğrenciye erişim izniniz yok")
	}
	return nil
}

func (s *instructorService) GetStudentPomodoros(ctx context.Context, instructorID, studentID uuid.UUID, from, to time.Time) (*dto.PomodoroStatsRes, error) {
	if err := s.mustBeStudent(ctx, instructorID, studentID); err != nil {
		return nil, err
	}

	totalMin, err := s.pomodoroRepo.SumMinutesByUser(ctx, studentID, from, to)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	daily, err := s.pomodoroRepo.DailyStatsByUser(ctx, studentID, from, to)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	var totalSessions int
	breakdown := make([]dto.DailyStats, len(daily))
	for i, d := range daily {
		totalSessions += d.Sessions
		breakdown[i] = dto.DailyStats{Date: d.Date, TotalMinutes: d.TotalMinutes, Sessions: d.Sessions}
	}

	return &dto.PomodoroStatsRes{
		TotalMinutes:   int(totalMin),
		TotalSessions:  totalSessions,
		DailyBreakdown: breakdown,
	}, nil
}

func (s *instructorService) GetStudentProgress(ctx context.Context, instructorID, studentID uuid.UUID) ([]dto.SubjectProgressRes, error) {
	if err := s.mustBeStudent(ctx, instructorID, studentID); err != nil {
		return nil, err
	}

	subjects, err := s.subjectRepo.FindAll(ctx)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	res := make([]dto.SubjectProgressRes, 0, len(subjects))
	for _, sub := range subjects {
		total, completed, _ := s.progressRepo.CountBySubject(ctx, studentID, sub.ID)
		var pct float64
		if total > 0 {
			pct = float64(completed) / float64(total) * 100
		}
		res = append(res, dto.SubjectProgressRes{
			SubjectID:       sub.ID,
			SubjectName:     sub.Name,
			TotalTopics:     int(total),
			CompletedTopics: int(completed),
			Percentage:      pct,
		})
	}
	return res, nil
}

func (s *instructorService) GetStudentExamResults(ctx context.Context, instructorID, studentID uuid.UUID, examType *entity.ExamType, pagination dto.PaginationReq) (*dto.PaginatedRes[dto.ExamResultRes], error) {
	if err := s.mustBeStudent(ctx, instructorID, studentID); err != nil {
		return nil, err
	}

	items, total, err := s.examRepo.ListByUser(ctx, studentID, examType, pagination.Offset(), pagination.Limit)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	res := make([]dto.ExamResultRes, len(items))
	for i, item := range items {
		res[i] = dto.ExamResultRes{
			ID:        item.ID,
			ExamType:  string(item.ExamType),
			ExamDate:  item.ExamDate,
			Scores:    item.Scores,
			TotalNet:  item.TotalNet,
			Note:      item.Note,
			CreatedAt: item.CreatedAt,
		}
	}

	paged := dto.NewPaginatedRes(res, total, pagination.Page, pagination.Limit)
	return &paged, nil
}
