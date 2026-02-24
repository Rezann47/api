package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/Rezann47/YksKoc/internal/repository"
)

type CreateStudyPlanInput struct {
	UserID    uuid.UUID
	CreatedBy uuid.UUID
	Title     string
	PlanDate  time.Time
	Note      *string
	Items     []CreateStudyPlanItemInput
}

type CreateStudyPlanItemInput struct {
	SubjectID       uuid.UUID
	TopicID         *uuid.UUID
	DurationMinutes int
	DisplayOrder    int16
}

type StudyPlanService interface {
	Create(ctx context.Context, input CreateStudyPlanInput) (*entity.StudyPlan, error)
	GetByID(ctx context.Context, id, requesterID uuid.UUID) (*entity.StudyPlan, error)
	ListByDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*entity.StudyPlan, error)
	ListByMonth(ctx context.Context, userID uuid.UUID, year, month int) ([]*entity.StudyPlan, error)
	Delete(ctx context.Context, id, requesterID uuid.UUID) error
	CompleteItem(ctx context.Context, planID, itemID, requesterID uuid.UUID) error
	UncompleteItem(ctx context.Context, planID, itemID, requesterID uuid.UUID) error
}

type studyPlanService struct {
	repo repository.StudyPlanRepository
	log  *zap.Logger
}

func NewStudyPlanService(repo repository.StudyPlanRepository, log *zap.Logger) StudyPlanService {
	return &studyPlanService{repo: repo, log: log}
}

func (s *studyPlanService) Create(ctx context.Context, input CreateStudyPlanInput) (*entity.StudyPlan, error) {
	title := input.Title
	if title == "" {
		title = "Çalışma Planı"
	}

	items := make([]entity.StudyPlanItem, 0, len(input.Items))
	for i, it := range input.Items {
		dur := it.DurationMinutes
		if dur <= 0 {
			dur = 30
		}
		order := it.DisplayOrder
		if order == 0 {
			order = int16(i)
		}
		items = append(items, entity.StudyPlanItem{
			SubjectID:       it.SubjectID,
			TopicID:         it.TopicID,
			DurationMinutes: dur,
			DisplayOrder:    order,
		})
	}

	plan := &entity.StudyPlan{
		UserID:    input.UserID,
		CreatedBy: input.CreatedBy,
		Title:     title,
		PlanDate:  input.PlanDate,
		Note:      input.Note,
		Items:     items,
	}

	if err := s.repo.Create(ctx, plan); err != nil {
		return nil, apperror.NewInternal(err)
	}

	// Preload ile dolu döndür
	return s.repo.GetByID(ctx, plan.ID)
}

func (s *studyPlanService) GetByID(ctx context.Context, id, requesterID uuid.UUID) (*entity.StudyPlan, error) {
	plan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewNotFound("çalışma planı", err)
	}
	// Sadece öğrencinin kendisi veya koçu görebilir
	if plan.UserID != requesterID && plan.CreatedBy != requesterID {
		return nil, apperror.NewForbidden("bu plana erişim yetkiniz yok")
	}
	return plan, nil
}

func (s *studyPlanService) ListByDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*entity.StudyPlan, error) {
	plans, err := s.repo.ListByUserAndDate(ctx, userID, date)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}
	return plans, nil
}

func (s *studyPlanService) ListByMonth(ctx context.Context, userID uuid.UUID, year, month int) ([]*entity.StudyPlan, error) {
	plans, err := s.repo.ListByUserAndMonth(ctx, userID, year, month)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}
	return plans, nil
}

func (s *studyPlanService) Delete(ctx context.Context, id, requesterID uuid.UUID) error {
	plan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperror.NewNotFound("çalışma planı", err)
	}
	if plan.UserID != requesterID && plan.CreatedBy != requesterID {
		return apperror.NewForbidden("bu planı silme yetkiniz yok")
	}
	return s.repo.Delete(ctx, id, plan.UserID)
}

func (s *studyPlanService) CompleteItem(ctx context.Context, planID, itemID, requesterID uuid.UUID) error {
	plan, err := s.repo.GetByID(ctx, planID)
	if err != nil {
		return apperror.NewNotFound("çalışma planı", err)
	}
	if plan.UserID != requesterID {
		return apperror.NewForbidden("bu planı düzenleme yetkiniz yok")
	}
	return s.repo.CompleteItem(ctx, itemID, planID)
}

func (s *studyPlanService) UncompleteItem(ctx context.Context, planID, itemID, requesterID uuid.UUID) error {
	plan, err := s.repo.GetByID(ctx, planID)
	if err != nil {
		return apperror.NewNotFound("çalışma planı", err)
	}
	if plan.UserID != requesterID {
		return apperror.NewForbidden("bu planı düzenleme yetkiniz yok")
	}
	return s.repo.UncompleteItem(ctx, itemID, planID)
}
