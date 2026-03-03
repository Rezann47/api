package server

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/Rezann47/YksKoc/internal/config"
	"github.com/Rezann47/YksKoc/internal/config/database"
	"github.com/Rezann47/YksKoc/internal/handler"
	"github.com/Rezann47/YksKoc/internal/middleware"
	"github.com/Rezann47/YksKoc/internal/repository"
	"github.com/Rezann47/YksKoc/internal/service"
)

type Server struct {
	http *http.Server
	log  *zap.Logger
}

func New(cfg *config.Config, log *zap.Logger) (*Server, error) {
	// 1. DB bağlantısı
	db, err := database.Connect(&cfg.DB, cfg.App.Env)
	if err != nil {
		return nil, fmt.Errorf("db connect: %w", err)
	}

	// 2. Repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewRefreshTokenRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)
	topicRepo := repository.NewTopicRepository(db)
	progressRepo := repository.NewProgressRepository(db)
	pomodoroRepo := repository.NewPomodoroRepository(db)
	examResultRepo := repository.NewExamResultRepository(db)
	instructorRepo := repository.NewInstructorStudentRepository(db)
	studyPlanRepo := repository.NewStudyPlanRepository(db)
	msgRepo := repository.NewMessageRepository(db)

	// 3. Services
	authSvc := service.NewAuthService(userRepo, tokenRepo, cfg.JWT, log)

	// ✅ AppleSharedSecret eklendi
	userSvc := service.NewUserService(userRepo, log, cfg.AppleSharedSecret)

	subjectSvc := service.NewSubjectService(subjectRepo, topicRepo, progressRepo, log)
	pomodoroSvc := service.NewPomodoroService(pomodoroRepo, log)
	examSvc := service.NewExamResultService(examResultRepo, log)
	instructorSvc := service.NewInstructorService(
		instructorRepo, userRepo, pomodoroRepo,
		progressRepo, subjectRepo, examResultRepo, log,
	)
	studyPlanSvc := service.NewStudyPlanService(studyPlanRepo, log)
	streakSvc := service.NewStreakService(db)
	messageSvc := service.NewMessageService(msgRepo, instructorRepo, log)

	// 4. Handlers & Router
	handlers := handler.NewHandlers(
		authSvc,
		userSvc,
		subjectSvc,
		pomodoroSvc,
		examSvc,
		instructorSvc,
		studyPlanSvc,
		messageSvc,
		streakSvc,
	)

	router := handler.NewRouter(handlers, cfg.JWT)
	router.Use(middleware.Logger(log))

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &Server{http: srv, log: log}, nil
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
