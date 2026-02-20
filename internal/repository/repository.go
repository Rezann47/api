package repository

import (
	"go-gin-crud/internal/dto"
	"go-gin-crud/internal/model"
	"math"

	"gorm.io/gorm"
)

// UserRepository arayüzü — dependency injection için
type UserRepository interface {
	Create(user *model.User) error
	FindAll(query dto.PaginationQuery) ([]model.User, int64, error)
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	Exists(email string) bool
}

// userRepo somut implementasyon
type userRepo struct {
	db *gorm.DB
}

// NewUserRepository constructor (dependency injection)
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// Create yeni kullanıcı ekler
func (r *userRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindAll sayfalama + arama + sıralama destekli liste
func (r *userRepo) FindAll(query dto.PaginationQuery) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// Default değerler
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}
	if query.Sort == "" {
		query.Sort = "created_at"
	}
	if query.Order == "" {
		query.Order = "desc"
	}

	db := r.db.Model(&model.User{})

	// Arama filtresi
	if query.Search != "" {
		search := "%" + query.Search + "%"
		db = db.Where("name ILIKE ? OR email ILIKE ?", search, search)
	}

	// Toplam kayıt sayısı
	db.Count(&total)

	// Sayfalama ve sıralama
	offset := (query.Page - 1) * query.Limit
	err := db.
		Order(query.Sort + " " + query.Order).
		Offset(offset).
		Limit(query.Limit).
		Find(&users).Error

	return users, total, err
}

// FindByID tek kullanıcı getirir
func (r *userRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail email ile kullanıcı getirir (login için)
func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update kullanıcıyı günceller (sadece dolu alanları)
func (r *userRepo) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete soft delete — kayıt silinmez, deleted_at dolar
func (r *userRepo) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// Exists email var mı kontrol eder
func (r *userRepo) Exists(email string) bool {
	var count int64
	r.db.Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// ==================== PRODUCT REPOSITORY ====================

type ProductRepository interface {
	Create(product *model.Product) error
	FindAll(userID uint, query dto.PaginationQuery) ([]model.Product, int64, error)
	FindByID(id uint) (*model.Product, error)
	Update(product *model.Product) error
	Delete(id, userID uint) error
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepo) FindAll(userID uint, query dto.PaginationQuery) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	db := r.db.Model(&model.Product{}).Where("user_id = ?", userID)

	if query.Search != "" {
		db = db.Where("name ILIKE ?", "%"+query.Search+"%")
	}

	db.Count(&total)

	offset := (query.Page - 1) * query.Limit
	err := db.
		Preload("User"). // İlişkili User verisini de çek
		Offset(offset).
		Limit(query.Limit).
		Find(&products).Error

	return products, total, err
}

func (r *productRepo) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("User").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

// Delete sadece kendi ürününü silebilir (userID kontrolü)
func (r *productRepo) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Product{}).Error
}

// TotalPages hesaplamak için yardımcı fonksiyon
func TotalPages(total int64, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}
