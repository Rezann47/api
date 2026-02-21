package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Rezann47/YksKoc/internal/config"
	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/Rezann47/YksKoc/internal/domain/dto"
	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/Rezann47/YksKoc/internal/repository"
	"github.com/Rezann47/YksKoc/pkg/jwt"
	"github.com/Rezann47/YksKoc/pkg/password"
	"github.com/Rezann47/YksKoc/pkg/studentcode"
)

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.RefreshTokenRepository
	jwtCfg    config.JWTConfig
	log       *zap.Logger
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.RefreshTokenRepository,
	jwtCfg config.JWTConfig,
	log *zap.Logger,
) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtCfg:    jwtCfg,
		log:       log,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterReq) (*dto.LoginRes, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperror.NewConflict("bu e-posta zaten kayıtlı", nil)
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		return nil, apperror.NewInternal(fmt.Errorf("hash password: %w", err))
	}

	user := &entity.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         entity.Role(req.Role),
	}

	// Öğrencilere benzersiz kod üret
	if user.Role == entity.RoleStudent {
		code := studentcode.Generate()
		user.StudentCode = &code
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.buildLoginRes(ctx, user, nil, nil)
}

func (s *authService) Login(ctx context.Context, req dto.LoginReq, userAgent, ip string) (*dto.LoginRes, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		// e-posta bulunamadı → güvenlik için aynı mesaj
		return nil, apperror.NewUnauthorized("e-posta veya şifre hatalı")
	}

	if !user.IsActive {
		return nil, apperror.NewUnauthorized("hesap devre dışı")
	}

	if err := password.Compare(user.PasswordHash, req.Password); err != nil {
		return nil, apperror.NewUnauthorized("e-posta veya şifre hatalı")
	}

	return s.buildLoginRes(ctx, user, &userAgent, &ip)
}

func (s *authService) Refresh(ctx context.Context, rawToken string) (*dto.TokenRes, error) {
	claims, err := jwt.ValidateRefresh(rawToken, s.jwtCfg.RefreshSecret)
	if err != nil {
		return nil, apperror.NewUnauthorized("geçersiz veya süresi dolmuş token")
	}

	hash := jwt.HashToken(rawToken)
	stored, err := s.tokenRepo.FindByHash(ctx, hash)
	if err != nil || !stored.IsValid() {
		return nil, apperror.NewUnauthorized("token geçersiz veya iptal edilmiş")
	}

	// Token rotate: eskiyi iptal et, yeni çift üret
	if err := s.tokenRepo.RevokeByHash(ctx, hash); err != nil {
		return nil, apperror.NewInternal(err)
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user, stored.UserAgent, stored.IPAddress)
}

func (s *authService) Logout(ctx context.Context, rawToken string) error {
	hash := jwt.HashToken(rawToken)
	return s.tokenRepo.RevokeByHash(ctx, hash)
}

func (s *authService) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	return s.tokenRepo.RevokeAllByUserID(ctx, userID)
}

// ─── helpers ──────────────────────────────────────────────

func (s *authService) buildLoginRes(ctx context.Context, user *entity.User, userAgent, ip *string) (*dto.LoginRes, error) {
	tokens, err := s.issueTokens(ctx, user, userAgent, ip)
	if err != nil {
		return nil, err
	}

	return &dto.LoginRes{
		User:  mapUserToAuthRes(user),
		Token: *tokens,
	}, nil
}

func (s *authService) issueTokens(ctx context.Context, user *entity.User, userAgent, ip *string) (*dto.TokenRes, error) {
	accessToken, err := jwt.GenerateAccess(user.ID, string(user.Role), s.jwtCfg.AccessSecret, s.jwtCfg.AccessExpiry)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	rawRefresh, refreshToken, err := jwt.GenerateRefresh(user.ID, s.jwtCfg.RefreshSecret, s.jwtCfg.RefreshExpiry)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	rt := &entity.RefreshToken{
		UserID:    user.ID,
		TokenHash: jwt.HashToken(rawRefresh),
		UserAgent: userAgent,
		IPAddress: ip,
		ExpiresAt: refreshToken.ExpiresAt,
	}
	if err := s.tokenRepo.Create(ctx, rt); err != nil {
		return nil, err
	}

	return &dto.TokenRes{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtCfg.AccessExpiry.Seconds()),
	}, nil
}

func mapUserToAuthRes(u *entity.User) dto.AuthUserRes {
	return dto.AuthUserRes{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Role:        string(u.Role),
		StudentCode: u.StudentCode,
		IsPremium:   u.IsPremium,
	}
}
