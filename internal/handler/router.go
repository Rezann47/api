package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Rezann47/YksKoc/internal/config"
	"github.com/Rezann47/YksKoc/internal/middleware"
	"github.com/Rezann47/YksKoc/internal/service"
)

type Handlers struct {
	Auth       *AuthHandler
	User       *UserHandler
	Subject    *SubjectHandler
	Pomodoro   *PomodoroHandler
	ExamResult *ExamResultHandler
	Instructor *InstructorHandler
	StudyPlan  *StudyPlanHandler
	Message    *MessageHandler
	Health     *HealthHandler
}

func NewRouter(h *Handlers, cfg config.JWTConfig) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// ─── Global Middleware ─────────────────────────────────
	r.Use(
		middleware.Recovery(nil),
		middleware.RequestID(),
		middleware.CORS(),
	)

	// ─── Health ───────────────────────────────────────────
	r.GET("/health", h.Health.Check)

	// ─── Swagger ──────────────────────────────────────────
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ─── API v1 ───────────────────────────────────────────
	v1 := r.Group("/api/v1")

	// Auth (public)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
		auth.POST("/refresh", h.Auth.Refresh)
		auth.POST("/logout", middleware.Auth(cfg), h.Auth.Logout)
		auth.POST("/logout-all", middleware.Auth(cfg), h.Auth.LogoutAll)
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.Auth(cfg))

	// ─── User ─────────────────────────────────────────────
	users := protected.Group("/users")
	{
		users.GET("/me", h.User.GetProfile)
		users.PATCH("/me", h.User.UpdateProfile)
		users.POST("/me/change-password", h.User.ChangePassword)
		users.GET("/me/premium", h.User.GetPremiumStatus)
		users.POST("/me/premium/activate", h.User.ActivatePremium)
		users.POST("/me/ping", h.User.Ping)
	}

	// ─── Students (öğrenci rolü) ──────────────────────────
	students := protected.Group("/students")
	students.Use(middleware.RequireRole("student"))
	{
		students.GET("/my-instructors", h.Instructor.ListMyInstructors)
	}

	// ─── Subjects & Topics ────────────────────────────────
	subjects := protected.Group("/subjects")
	subjects.Use(middleware.RequireRole("student", "instructor"))
	{
		subjects.GET("", h.Subject.ListSubjects)
		subjects.GET("/:subjectID/topics", h.Subject.ListTopics)
		subjects.GET("/:subjectID/progress", h.Subject.GetSubjectProgress)
		subjects.GET("/progress", h.Subject.GetAllProgress)
	}

	topics := protected.Group("/topics")
	topics.Use(middleware.RequireRole("student"))
	{
		topics.PATCH("/:topicID/mark", h.Subject.MarkTopic)
	}

	// ─── Pomodoro (öğrenci) ───────────────────────────────
	pomodoros := protected.Group("/pomodoros")
	pomodoros.Use(middleware.RequireRole("student"))
	{
		pomodoros.POST("", h.Pomodoro.Create)
		pomodoros.GET("", h.Pomodoro.List)
		pomodoros.GET("/stats", h.Pomodoro.GetStats)
		pomodoros.DELETE("/:id", h.Pomodoro.Delete)
	}

	// ─── Exam Results (öğrenci) ───────────────────────────
	exams := protected.Group("/exam-results")
	exams.Use(middleware.RequireRole("student"))
	{
		exams.POST("", h.ExamResult.Create)
		exams.GET("", h.ExamResult.List)
		exams.GET("/stats", h.ExamResult.GetStats)
		exams.DELETE("/:id", h.ExamResult.Delete)
	}

	// ─── Instructor ───────────────────────────────────────
	instructor := protected.Group("/instructor")
	instructor.Use(middleware.RequireRole("instructor"))
	{
		instructor.POST("/students", h.Instructor.AddStudent)
		instructor.GET("/students", h.Instructor.ListStudents)
		instructor.DELETE("/students/:studentID", h.Instructor.RemoveStudent)

		instructor.GET("/students/:studentID/pomodoros", h.Instructor.GetStudentPomodoros)
		instructor.GET("/students/:studentID/progress", h.Instructor.GetStudentProgress)
		instructor.GET("/students/:studentID/exam-results", h.Instructor.GetStudentExamResults)

		instructor.POST("/students/:studentID/study-plans", h.StudyPlan.CreateForStudent)
		instructor.GET("/students/:studentID/study-plans", h.StudyPlan.GetStudentPlans)
	}

	// ─── Messages (her iki rol) ───────────────────────────
	messages := protected.Group("/messages")
	{
		messages.POST("", h.Message.Send)
		messages.GET("/conversations", h.Message.ListConversations)
		messages.GET("/conversations/:peerID", h.Message.GetConversation)
		messages.POST("/conversations/:peerID/read", h.Message.MarkRead)
		messages.GET("/unread", h.Message.UnreadCount)
	}

	// ─── Study Plans (öğrenci) ────────────────────────────
	studyPlans := protected.Group("/study-plans")
	studyPlans.Use(middleware.RequireRole("student"))
	{
		studyPlans.POST("", h.StudyPlan.Create)
		studyPlans.GET("", h.StudyPlan.ListByDate)
		studyPlans.GET("/month", h.StudyPlan.ListByMonth)
		studyPlans.DELETE("/:id", h.StudyPlan.Delete)
		studyPlans.PATCH("/:id/items/:itemID/complete", h.StudyPlan.CompleteItem)
		studyPlans.PATCH("/:id/items/:itemID/uncomplete", h.StudyPlan.UncompleteItem)
	}

	return r
}

// NewHandlers tüm handler'ları service bağımlılıklarıyla oluşturur
func NewHandlers(
	authSvc service.AuthService,
	userSvc service.UserService,
	subjectSvc service.SubjectService,
	pomodoroSvc service.PomodoroService,
	examSvc service.ExamResultService,
	instructorSvc service.InstructorService,
	studyPlanSvc service.StudyPlanService,
	messageSvc service.MessageService,
) *Handlers {
	return &Handlers{
		Auth:       NewAuthHandler(authSvc),
		User:       NewUserHandler(userSvc),
		Subject:    NewSubjectHandler(subjectSvc),
		Pomodoro:   NewPomodoroHandler(pomodoroSvc),
		ExamResult: NewExamResultHandler(examSvc),
		Instructor: NewInstructorHandler(instructorSvc),
		StudyPlan:  NewStudyPlanHandler(studyPlanSvc),
		Message:    NewMessageHandler(messageSvc),
		Health:     NewHealthHandler(),
	}
}
