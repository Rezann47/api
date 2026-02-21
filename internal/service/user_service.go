package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/Rezann47/YksKoc/internal/domain/dto"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/Rezann47/YksKoc/internal/repository"
	"github.com/Rezann47/YksKoc/pkg/password"
)

type userService struct {
	userRepo repository.UserRepository
	log      *zap.Logger
}

func NewUserService(userRepo repository.UserRepository, log *zap.Logger) UserService {
	return &userService{userRepo: userRepo, log: log}
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserRes, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := mapUserToRes(user)
	return &res, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileReq) (*dto.UserRes, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	res := mapUserToRes(user)
	return &res, nil
}

func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordReq) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := password.Compare(user.PasswordHash, req.CurrentPassword); err != nil {
		return apperror.NewUnauthorized("mevcut şifre hatalı")
	}
	hash, err := password.Hash(req.NewPassword)
	if err != nil {
		return apperror.NewInternal(err)
	}
	user.PasswordHash = hash
	return s.userRepo.Update(ctx, user)
}

func (s *userService) GetPremiumStatus(ctx context.Context, userID uuid.UUID) (*dto.PremiumStatusRes, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &dto.PremiumStatusRes{IsPremium: user.IsPremium}, nil
}

func (s *userService) ActivatePremium(ctx context.Context, userID uuid.UUID) (*dto.PremiumStatusRes, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.IsPremium = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return &dto.PremiumStatusRes{IsPremium: true}, nil
}

func mapUserToRes(u *entity.User) dto.UserRes {
	return dto.UserRes{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Role:        string(u.Role),
		StudentCode: u.StudentCode,
		IsPremium:   u.IsPremium,
		CreatedAt:   u.CreatedAt,
	}
}
