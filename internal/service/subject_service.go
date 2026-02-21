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

type subjectService struct {
	subjectRepo  repository.SubjectRepository
	topicRepo    repository.TopicRepository
	progressRepo repository.ProgressRepository
	log          *zap.Logger
}

func NewSubjectService(
	subjectRepo repository.SubjectRepository,
	topicRepo repository.TopicRepository,
	progressRepo repository.ProgressRepository,
	log *zap.Logger,
) SubjectService {
	return &subjectService{
		subjectRepo:  subjectRepo,
		topicRepo:    topicRepo,
		progressRepo: progressRepo,
		log:          log,
	}
}

func (s *subjectService) ListSubjects(ctx context.Context, examType *entity.ExamType) ([]dto.SubjectRes, error) {
	var subjects []entity.Subject
	var err error

	if examType != nil {
		subjects, err = s.subjectRepo.FindByExamType(ctx, *examType)
	} else {
		subjects, err = s.subjectRepo.FindAll(ctx)
	}
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	res := make([]dto.SubjectRes, len(subjects))
	for i, sub := range subjects {
		res[i] = dto.SubjectRes{
			ID:           sub.ID,
			Name:         sub.Name,
			ExamType:     string(sub.ExamType),
			DisplayOrder: sub.DisplayOrder,
		}
	}
	return res, nil
}

func (s *subjectService) ListTopics(ctx context.Context, subjectID uuid.UUID, userID uuid.UUID) ([]dto.TopicRes, error) {
	if _, err := s.subjectRepo.FindByID(ctx, subjectID); err != nil {
		return nil, err
	}

	topics, err := s.topicRepo.FindBySubjectID(ctx, subjectID)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	completedIDs, err := s.progressRepo.FindCompletedByUserAndSubject(ctx, userID, subjectID)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	completedSet := make(map[uuid.UUID]bool, len(completedIDs))
	for _, id := range completedIDs {
		completedSet[id] = true
	}

	res := make([]dto.TopicRes, len(topics))
	for i, t := range topics {
		res[i] = dto.TopicRes{
			ID:           t.ID,
			SubjectID:    t.SubjectID,
			Name:         t.Name,
			DisplayOrder: t.DisplayOrder,
			IsCompleted:  completedSet[t.ID],
		}
	}
	return res, nil
}

func (s *subjectService) MarkTopic(ctx context.Context, userID, topicID uuid.UUID, isCompleted bool) error {
	if _, err := s.topicRepo.FindByID(ctx, topicID); err != nil {
		return err
	}
	if isCompleted {
		return s.progressRepo.Upsert(ctx, userID, topicID)
	}
	return s.progressRepo.Delete(ctx, userID, topicID)
}

func (s *subjectService) GetSubjectProgress(ctx context.Context, userID, subjectID uuid.UUID) (*dto.SubjectProgressRes, error) {
	subject, err := s.subjectRepo.FindByID(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	total, completed, err := s.progressRepo.CountBySubject(ctx, userID, subjectID)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	var pct float64
	if total > 0 {
		pct = float64(completed) / float64(total) * 100
	}

	return &dto.SubjectProgressRes{
		SubjectID:       subject.ID,
		SubjectName:     subject.Name,
		TotalTopics:     int(total),
		CompletedTopics: int(completed),
		Percentage:      pct,
	}, nil
}

func (s *subjectService) GetAllProgress(ctx context.Context, userID uuid.UUID) ([]dto.SubjectProgressRes, error) {
	subjects, err := s.subjectRepo.FindAll(ctx)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	res := make([]dto.SubjectProgressRes, 0, len(subjects))
	for _, sub := range subjects {
		total, completed, err := s.progressRepo.CountBySubject(ctx, userID, sub.ID)
		if err != nil {
			continue
		}
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
