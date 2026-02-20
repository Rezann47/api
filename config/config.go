package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config tüm uygulama konfigürasyonunu tutar
type Config struct {
	App      AppConfig
	DB       DBConfig
	JWT      JWTConfig
	RateLimit RateLimitConfig
}

type AppConfig struct {
	Env  string
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Timezone string
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type RateLimitConfig struct {
	Limit int
	Burst int
}

// DSN PostgreSQL bağlantı stringini oluşturur
func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode, d.Timezone,
	)
}

// Load .env dosyasını okur ve Config struct'ını doldurur
func Load() *Config {
	// .env dosyası yoksa sistem env değişkenlerini kullan
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env dosyası bulunamadı, sistem env değişkenleri kullanılıyor")
	}

	expireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	rateLimit, _   := strconv.Atoi(getEnv("RATE_LIMIT", "100"))
	rateBurst, _   := strconv.Atoi(getEnv("RATE_LIMIT_BURST", "20"))

	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "go_crud_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Timezone: getEnv("DB_TIMEZONE", "Europe/Istanbul"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "change-me"),
			ExpireHours: expireHours,
		},
		RateLimit: RateLimitConfig{
			Limit: rateLimit,
			Burst: rateBurst,
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
