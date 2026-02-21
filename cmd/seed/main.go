// cmd/seed — geliştirme / staging için örnek veri yükler
// Production'da çalıştırılmaz.

package main

import (
	"log"
	"os"

	"github.com/Rezann47/YksKoc/internal/config"
	"github.com/Rezann47/YksKoc/internal/config/database"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "production" {
		log.Fatal("seed komutu production'da çalıştırılamaz")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.Connect(&cfg.DB, cfg.App.Env)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// TODO: seed fonksiyonları buraya eklenecek
	// Örn: seedTestUsers(db), seedExtraTopics(db)
	log.Println("✓ seed tamamlandı (subjects/topics migration'dan gelir)")
}
