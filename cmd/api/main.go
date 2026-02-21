package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Rezann47/YksKoc/internal/config"
	"github.com/Rezann47/YksKoc/internal/server"
	"github.com/Rezann47/YksKoc/pkg/logger"
	pkgmigrate "github.com/Rezann47/YksKoc/pkg/migrate"
)

// @title           YKS Tracker API
// @version         1.0
// @description     YKS (TYT & AYT) öğrenci ve eğitmen takip uygulaması API'si
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@ykstracker.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Bearer {token}" formatında giriniz

func main() {
	// 1. Config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config load: %v\n", err)
		os.Exit(1)
	}

	// 2. Logger
	log := logger.New(cfg.App.Env)
	defer log.Sync() //nolint:errcheck

	// 3. Migration
	if err := runMigrations(cfg, log); err != nil {
		log.Fatal("migration failed", zap.Error(err))
	}

	// 4. Server
	srv, err := server.New(cfg, log)
	if err != nil {
		log.Fatal("server init failed", zap.Error(err))
	}

	// 5. Graceful shutdown
	go func() {
		log.Info("server starting", zap.String("port", cfg.Server.Port))
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("forced shutdown", zap.Error(err))
	}
	log.Info("server stopped")
}

func runMigrations(cfg *config.Config, log *zap.Logger) error {
	runner, err := pkgmigrate.New(&cfg.DB, "file://migrations")
	if err != nil {
		return fmt.Errorf("create runner: %w", err)
	}
	defer runner.Close() //nolint:errcheck

	version, dirty, _ := runner.Version()
	if dirty {
		return fmt.Errorf("db dirty at version %d — fix and re-run", version)
	}

	if err := runner.Up(); err != nil {
		return err
	}

	newVer, _, _ := runner.Version()
	if newVer != version {
		log.Info("migrations applied", zap.Uint("version", newVer))
	} else {
		log.Info("migrations up to date", zap.Uint("version", newVer))
	}
	return nil
}
