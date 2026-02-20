package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"go-gin-crud/config"
	"go-gin-crud/internal/handler"
	"go-gin-crud/internal/middleware"
	"go-gin-crud/internal/model"
	"go-gin-crud/internal/repository"
	"go-gin-crud/internal/service"
)

type HTTPServer struct {
	cfg *config.Config
}

func NewHTTPServer(cfg *config.Config) *HTTPServer {
	return &HTTPServer{cfg: cfg}
}

func (s *HTTPServer) Run() error {
	cfg := s.cfg

	// ── Database ─────────────────────────
	db := config.NewDB(cfg)

	if err := db.AutoMigrate(&model.User{}, &model.Product{}); err != nil {
		return err
	}
	log.Println("✅ Migration tamamlandı")

	// ── Dependency Injection ─────────────
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)

	userSvc := service.NewUserService(userRepo, cfg)
	productSvc := service.NewProductService(productRepo)

	userHandler := handler.NewUserHandler(userSvc)
	productHandler := handler.NewProductHandler(productSvc)

	// ── Gin ──────────────────────────────
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit(cfg.RateLimit.Limit))

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"env":     cfg.App.Env,
			"version": "1.0.0",
			"time":    time.Now().UTC(),
		})
	})

	// Routes
	registerRoutes(r, cfg, userHandler, productHandler)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint bulunamadı"})
	})

	// ── Server ───────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server running on :%s\n", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏳ Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
