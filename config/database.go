package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(cfg *Config) *gorm.DB {
	logLevel := logger.Info
	if cfg.App.Env == "production" {
		logLevel = logger.Error
	}

	// Önce DATABASE_URL kontrol et
	dsn := os.Getenv("DATABASE_URL")

	// Eğer yoksa local DSN kullan
	if dsn == "" {
		dsn = cfg.DB.DSN()
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		log.Fatalf("❌ Veritabanına bağlanılamadı: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ SQL DB alınamadı: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ PostgreSQL bağlantısı kuruldu")
	return db
}
