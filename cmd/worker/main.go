// cmd/worker — arka plan işçisi (bildirimler, periyodik raporlar, vb.)
// Şu an placeholder; ilerleyen sprint'lerde doldurulacak.

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Rezann47/YksKoc/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	_ = cfg

	log.Println("worker başlatıldı — Ctrl+C ile dur")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("worker durduruldu")
}
