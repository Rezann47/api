package model

import (
	"time"

	"gorm.io/gorm"
)

// User veritabanı modeli
// GORM tag'leri ile tablo ve kolon ayarları yapılır
type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"           json:"id"`
	Name      string         `gorm:"type:varchar(100);not null"          json:"name"`
	Email     string         `gorm:"type:varchar(150);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null"          json:"-"` // json:"-" → response'a dahil etme
	Role      string         `gorm:"type:varchar(20);default:'user'"    json:"role"`
	IsActive  bool           `gorm:"default:true"                        json:"is_active"`
	CreatedAt time.Time      `                                           json:"created_at"`
	UpdatedAt time.Time      `                                           json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                               json:"-"` // Soft delete
}

// TableName GORM'un kullanacağı tablo adını belirtir
func (User) TableName() string {
	return "users"
}

// Product ikinci örnek model — ilişki göstermek için
type Product struct {
	ID          uint           `gorm:"primaryKey;autoIncrement"        json:"id"`
	Name        string         `gorm:"type:varchar(200);not null"       json:"name"`
	Description string         `gorm:"type:text"                        json:"description"`
	Price       float64        `gorm:"type:decimal(10,2);not null"      json:"price"`
	Stock       int            `gorm:"default:0"                        json:"stock"`
	UserID      uint           `gorm:"not null;index"                   json:"user_id"`  // Foreign key
	User        User           `gorm:"foreignKey:UserID"                json:"user,omitempty"` // Belongs to
	CreatedAt   time.Time      `                                        json:"created_at"`
	UpdatedAt   time.Time      `                                        json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                            json:"-"`
}

func (Product) TableName() string {
	return "products"
}
