package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go-gin-crud/config"
	"go-gin-crud/internal/dto"
	"go-gin-crud/internal/model"
	"go-gin-crud/internal/repository"
)

// ==================== USER SERVICE ====================

type UserService interface {
	Register(req dto.CreateUserRequest) (*dto.UserResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	GetAll(query dto.PaginationQuery) (*dto.PaginatedResponse, error)
	GetByID(id uint) (*dto.UserResponse, error)
	Update(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(id uint) error
}

type userService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewUserService(repo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{repo: repo, cfg: cfg}
}

// Register yeni kullanıcı kaydeder
func (s *userService) Register(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Email daha önce kullanılmış mı?
	if s.repo.Exists(req.Email) {
		return nil, errors.New("bu email zaten kullanılıyor")
	}

	// Şifreyi hashle (bcrypt cost=12 production için ideal)
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("şifre işlenemedi: %w", err)
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPw),
		Role:     role,
		IsActive: true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("kullanıcı kaydedilemedi: %w", err)
	}

	return toUserResponse(user), nil
}

// Login kullanıcıyı doğrular ve JWT döner
func (s *userService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email veya şifre hatalı")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("hesap aktif değil")
	}

	// Şifre kontrolü
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("email veya şifre hatalı")
	}

	// JWT oluştur
	expiresAt := time.Now().Add(time.Duration(s.cfg.JWT.ExpireHours) * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     expiresAt.Unix(),
	})

	tokenStr, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, fmt.Errorf("token oluşturulamadı: %w", err)
	}

	return &dto.LoginResponse{
		Token:     tokenStr,
		ExpiresAt: expiresAt.Format(time.RFC3339),
		User:      *toUserResponse(user),
	}, nil
}

// GetAll sayfalı kullanıcı listesi
func (s *userService) GetAll(query dto.PaginationQuery) (*dto.PaginatedResponse, error) {
	users, total, err := s.repo.FindAll(query)
	if err != nil {
		return nil, err
	}

	// Model → DTO dönüşümü
	responses := make([]dto.UserResponse, len(users))
	for i, u := range users {
		responses[i] = *toUserResponse(&u)
	}

	limit := query.Limit
	if limit == 0 {
		limit = 10
	}

	return &dto.PaginatedResponse{
		Data:       responses,
		Total:      total,
		Page:       query.Page,
		Limit:      limit,
		TotalPages: repository.TotalPages(total, limit),
	}, nil
}

// GetByID tek kullanıcı getirir
func (s *userService) GetByID(id uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kullanıcı bulunamadı")
		}
		return nil, err
	}
	return toUserResponse(user), nil
}

// Update kullanıcıyı günceller
func (s *userService) Update(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kullanıcı bulunamadı")
		}
		return nil, err
	}

	// Sadece gönderilen alanları güncelle
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" && req.Email != user.Email {
		if s.repo.Exists(req.Email) {
			return nil, errors.New("bu email zaten kullanılıyor")
		}
		user.Email = req.Email
	}
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.repo.Update(user); err != nil {
		return nil, fmt.Errorf("güncelleme başarısız: %w", err)
	}

	return toUserResponse(user), nil
}

// Delete kullanıcıyı soft delete yapar
func (s *userService) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("kullanıcı bulunamadı")
		}
		return err
	}
	return s.repo.Delete(id)
}

// toUserResponse model → DTO dönüşümü (şifreyi gizler)
func toUserResponse(u *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

// ==================== PRODUCT SERVICE ====================

type ProductService interface {
	Create(userID uint, req dto.CreateProductRequest) (*model.Product, error)
	GetAll(userID uint, query dto.PaginationQuery) (*dto.PaginatedResponse, error)
	GetByID(id uint) (*model.Product, error)
	Update(id, userID uint, req dto.UpdateProductRequest) (*model.Product, error)
	Delete(id, userID uint) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(userID uint, req dto.CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		UserID:      userID,
	}
	if err := s.repo.Create(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) GetAll(userID uint, query dto.PaginationQuery) (*dto.PaginatedResponse, error) {
	products, total, err := s.repo.FindAll(userID, query)
	if err != nil {
		return nil, err
	}

	limit := query.Limit
	if limit == 0 {
		limit = 10
	}

	return &dto.PaginatedResponse{
		Data:       products,
		Total:      total,
		Page:       query.Page,
		Limit:      limit,
		TotalPages: repository.TotalPages(total, limit),
	}, nil
}

func (s *productService) GetByID(id uint) (*model.Product, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ürün bulunamadı")
		}
		return nil, err
	}
	return p, nil
}

func (s *productService) Update(id, userID uint, req dto.UpdateProductRequest) (*model.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("ürün bulunamadı")
	}
	if product.UserID != userID {
		return nil, errors.New("bu ürünü güncelleme yetkiniz yok")
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) Delete(id, userID uint) error {
	return s.repo.Delete(id, userID)
}
