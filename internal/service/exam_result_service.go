package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/Rezann47/YksKoc/internal/domain/dto"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/Rezann47/YksKoc/internal/repository"
)

type examResultService struct {
	repo repository.ExamResultRepository
	log  *zap.Logger
}

func NewExamResultService(repo repository.ExamResultRepository, log *zap.Logger) ExamResultService {
	return &examResultService{repo: repo, log: log}
}

func (s *examResultService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateExamResultReq) (*dto.ExamResultRes, error) {
	scores := make(entity.ExamScores, len(req.Scores))
	var totalNet float64

	for subject, sc := range req.Scores {
		net := float64(sc.Correct) - float64(sc.Wrong)/4.0
		scores[subject] = entity.SubjectScores{
			Correct: sc.Correct,
			Wrong:   sc.Wrong,
			Net:     net,
		}
		totalNet += net
	}

	result := &entity.ExamResult{
		UserID:   userID,
		ExamType: entity.ExamType(req.ExamType),
		ExamDate: req.ExamDate,
		Scores:   scores,
		TotalNet: totalNet,
		Note:     req.Note,
	}

	if err := s.repo.Create(ctx, result); err != nil {
		return nil, apperror.NewInternal(err)
	}

	return mapExamResultToRes(result), nil
}

func (s *examResultService) List(ctx context.Context, userID uuid.UUID, examType *entity.ExamType, pagination dto.PaginationReq) (*dto.PaginatedRes[dto.ExamResultRes], error) {
	items, total, err := s.repo.ListByUser(ctx, userID, examType, pagination.Offset(), pagination.Limit)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	res := make([]dto.ExamResultRes, len(items))
	for i, item := range items {
		res[i] = *mapExamResultToRes(&item)
	}
	paged := dto.NewPaginatedRes(res, total, pagination.Page, pagination.Limit)
	return &paged, nil
}

func (s *examResultService) GetStats(ctx context.Context, userID uuid.UUID, examType entity.ExamType) (*dto.ExamStatsRes, error) {
	items, total, err := s.repo.ListByUser(ctx, userID, &examType, 0, 1000)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}
	if total == 0 {
		return &dto.ExamStatsRes{ExamType: string(examType)}, nil
	}

	var sumNet, bestNet float64
	trend := make([]dto.TrendPoint, 0, len(items))

	for _, item := range items {
		sumNet += item.TotalNet
		if item.TotalNet > bestNet {
			bestNet = item.TotalNet
		}
		trend = append(trend, dto.TrendPoint{Date: item.ExamDate, TotalNet: item.TotalNet})
	}

	averages, _ := s.repo.AverageNetByUser(ctx, userID, examType)

	return &dto.ExamStatsRes{
		ExamType:        string(examType),
		TotalExams:      int(total),
		AverageTotalNet: sumNet / float64(total),
		BestTotalNet:    bestNet,
		SubjectAverages: averages,
		Trend:           trend,
	}, nil
}

func (s *examResultService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.Delete(ctx, id, userID)
}

func mapExamResultToRes(r *entity.ExamResult) *dto.ExamResultRes {
	return &dto.ExamResultRes{
		ID:        r.ID,
		ExamType:  string(r.ExamType),
		ExamDate:  r.ExamDate,
		Scores:    r.Scores,
		TotalNet:  r.TotalNet,
		Note:      r.Note,
		CreatedAt: r.CreatedAt,
	}
}
