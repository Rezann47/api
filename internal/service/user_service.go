package service

import (
	"context"
	"time"

	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/Rezann47/YksKoc/internal/domain/dto"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/Rezann47/YksKoc/internal/repository"
	"github.com/Rezann47/YksKoc/pkg/password"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userService struct {
	userRepo    repository.UserRepository
	log         *zap.Logger
	appleSecret string // ✅ eklendi
}

func NewUserService(userRepo repository.UserRepository, log *zap.Logger, appleSecret string) UserService {
	return &userService{
		userRepo:    userRepo,
		log:         log,
		appleSecret: appleSecret, // ✅ eklendi
	}
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserRes, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.IsPremium && user.PremiumExpiresAt != nil && user.PremiumExpiresAt.Before(time.Now()) {
		user.IsPremium = false
		if err := s.userRepo.UpdatePremium(ctx, user.ID, false, user.PremiumExpiresAt, user.LastPremiumTransaction); err != nil {
			s.log.Warn("premium expire güncellenemedi", zap.Error(err))
		}
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
	if req.AvatarID != nil {
		user.AvatarID = *req.AvatarID
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

	if user.IsPremium && user.PremiumExpiresAt != nil && user.PremiumExpiresAt.Before(time.Now()) {
		user.IsPremium = false
		if err := s.userRepo.UpdatePremium(ctx, user.ID, false, user.PremiumExpiresAt, user.LastPremiumTransaction); err != nil {
			s.log.Warn("premium expire güncellenemedi", zap.Error(err))
		}
	}

	return &dto.PremiumStatusRes{
		IsPremium:        user.IsPremium,
		PremiumExpiresAt: user.PremiumExpiresAt,
	}, nil
}

func (s *userService) ActivatePremium(ctx context.Context, userID uuid.UUID, req dto.ActivatePremiumReq) (*dto.PremiumStatusRes, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var expiresAt *time.Time
	txID := req.TransactionID

	if req.Platform == "ios" && req.ReceiptData != "" {
		// Apple'a doğrulat
		valid, appleExpires, appleTxID, err := VerifyAppleReceipt(ctx, req.ReceiptData, s.appleSecret)
		if err != nil || !valid {
			s.log.Warn("apple receipt doğrulama başarısız", zap.Error(err))
			return nil, apperror.NewBadRequest("geçersiz receipt")
		}
		expiresAt = appleExpires
		txID = appleTxID
	} else {
		// Android veya receipt gelmedi
		if req.ExpiresDate != nil {
			expiresAt = req.ExpiresDate
		} else {
			now := time.Now()
			var t time.Time
			if user.PremiumExpiresAt != nil && user.PremiumExpiresAt.After(now) {
				t = user.PremiumExpiresAt.AddDate(0, 1, 0)
			} else {
				t = now.AddDate(0, 1, 0)
			}
			expiresAt = &t
		}
	}

	// Duplicate transaction kontrolü
	if txID != "" && user.LastPremiumTransaction == txID {
		return &dto.PremiumStatusRes{
			IsPremium:        user.IsPremium,
			PremiumExpiresAt: user.PremiumExpiresAt,
		}, nil
	}

	if err := s.userRepo.UpdatePremium(ctx, userID, true, expiresAt, txID); err != nil {
		return nil, err
	}

	return &dto.PremiumStatusRes{
		IsPremium:        true,
		PremiumExpiresAt: expiresAt,
	}, nil
}

func (s *userService) Ping(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.UpdateLastSeen(ctx, userID, time.Now())
}

func mapUserToRes(u *entity.User) dto.UserRes {
	return dto.UserRes{
		ID:               u.ID,
		Name:             u.Name,
		Email:            u.Email,
		Role:             string(u.Role),
		StudentCode:      u.StudentCode,
		IsPremium:        u.IsPremium,
		AvatarID:         u.AvatarID,
		LastSeenAt:       u.LastSeenAt,
		IsOnline:         u.IsOnline(),
		CreatedAt:        u.CreatedAt,
		PremiumExpiresAt: u.PremiumExpiresAt,
	}
}
func (s *userService) DeleteAccount(ctx context.Context, userID uuid.UUID, req dto.DeleteAccountReq) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := password.Compare(user.PasswordHash, req.Password); err != nil {
		return apperror.NewUnauthorized("şifre hatalı")
	}
	return s.userRepo.DeleteAccount(ctx, userID)
}
