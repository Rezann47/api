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

type pomodoroService struct {
	repo repository.PomodoroRepository
	log  *zap.Logger
}

func NewPomodoroService(repo repository.PomodoroRepository, log *zap.Logger) PomodoroService {
	return &pomodoroService{repo: repo, log: log}
}

// Create Pomodoro (test-ready, fallback ile)
func (s *pomodoroService) Create(ctx context.Context, userID uuid.UUID, req dto.CreatePomodoroReq) (*dto.PomodoroRes, error) {
	// Başlangıç zamanı pointer ile
	var startedAt *time.Time
	if req.StartedAt != nil {
		startedAt = req.StartedAt
	} else {
		now := time.Now()
		startedAt = &now
	}

	// Minimum 1 dakika
	duration := req.DurationMinutes
	if duration < 1 {
		duration = 1
	}

	p := &entity.Pomodoro{
		UserID:          userID,
		SubjectID:       req.SubjectID,
		DurationMinutes: duration,
		StartedAt:       *startedAt, // entity time.Time olduğu için dereference ediyoruz
	}

	// DB insert hatalarını logla
	if err := s.repo.Create(ctx, p); err != nil {
		s.log.Error("Pomodoro create failed",
			zap.Error(err),
			zap.Any("payload", p),
		)
		return nil, apperror.NewInternal(err)
	}

	return mapPomodoroToRes(p, nil), nil
}

// Listeleme
func (s *pomodoroService) List(ctx context.Context, userID uuid.UUID, filter dto.PomodoroListFilter) (*dto.PaginatedRes[dto.PomodoroRes], error) {
	items, total, err := s.repo.ListByUser(ctx, userID, filter.From, filter.To, filter.Offset(), filter.Limit)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	res := make([]dto.PomodoroRes, len(items))
	for i, p := range items {
		var subName *string
		if p.Subject != nil {
			n := p.Subject.Name
			subName = &n
		}
		res[i] = *mapPomodoroToRes(&p, subName)
	}

	paged := dto.NewPaginatedRes(res, total, filter.Page, filter.Limit)
	return &paged, nil
}

// İstatistik
func (s *pomodoroService) GetStats(ctx context.Context, userID uuid.UUID, from, to time.Time) (*dto.PomodoroStatsRes, error) {
	totalMin, err := s.repo.SumMinutesByUser(ctx, userID, from, to)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	daily, err := s.repo.DailyStatsByUser(ctx, userID, from, to)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	var totalSessions int
	breakdown := make([]dto.DailyStats, len(daily))
	for i, d := range daily {
		totalSessions += d.Sessions
		breakdown[i] = dto.DailyStats{
			Date:         d.Date,
			TotalMinutes: d.TotalMinutes,
			Sessions:     d.Sessions,
		}
	}

	return &dto.PomodoroStatsRes{
		TotalMinutes:   int(totalMin),
		TotalSessions:  totalSessions,
		DailyBreakdown: breakdown,
	}, nil
}

// Silme
func (s *pomodoroService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.Delete(ctx, id, userID)
}

// Mapper
func mapPomodoroToRes(p *entity.Pomodoro, subjectName *string) *dto.PomodoroRes {
	return &dto.PomodoroRes{
		ID:              p.ID,
		DurationMinutes: p.DurationMinutes,
		SubjectID:       p.SubjectID,
		SubjectName:     subjectName,
		StartedAt:       p.StartedAt,
		CreatedAt:       p.CreatedAt,
	}
}
