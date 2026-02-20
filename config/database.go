package config

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB GORM + PostgreSQL bağlantısını kurar
func NewDB(cfg *Config) *gorm.DB {
	logLevel := logger.Info
	if cfg.App.Env == "production" {
		logLevel = logger.Error // Production'da sadece hataları logla
	}

	db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		// Soft delete için DeletedAt alanını otomatik yönet
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		log.Fatalf("❌ Veritabanına bağlanılamadı: %v", err)
	}

	// Connection Pool ayarları
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ SQL DB alınamadı: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)           // Minimum açık bağlantı
	sqlDB.SetMaxOpenConns(100)          // Maksimum açık bağlantı
	sqlDB.SetConnMaxLifetime(time.Hour) // Bağlantı ömrü

	log.Println("✅ PostgreSQL bağlantısı kuruldu")
	return db
}
