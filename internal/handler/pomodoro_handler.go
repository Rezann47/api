package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Rezann47/YksKoc/internal/domain/dto"
	"github.com/Rezann47/YksKoc/internal/middleware"
	"github.com/Rezann47/YksKoc/internal/service"
	"github.com/Rezann47/YksKoc/pkg/response"
)

type PomodoroHandler struct {
	svc       service.PomodoroService
	streakSvc *service.StreakService // ← YENİ
}

func NewPomodoroHandler(svc service.PomodoroService, streakSvc *service.StreakService) *PomodoroHandler {
	return &PomodoroHandler{svc: svc, streakSvc: streakSvc}
}

func (h *PomodoroHandler) Create(c *gin.Context) {
	var req dto.CreatePomodoroReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	res, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	// ── Streak güncelle ──────────────────────────────────────
	// Pomodoro kaydedildikten sonra günlük seriyi güncelle.
	// Hata olsa bile pomodoro yanıtını döndür (streak kritik değil).
	streakResult, streakErr := h.streakSvc.UpdateStreak(c.Request.Context(), userID)
	if streakErr == nil && len(streakResult.NewBadges) > 0 {
		// Yeni rozet kazanıldı → frontend popup gösterebilir
		response.Created(c, gin.H{
			"pomodoro": res,
			"streak":   streakResult,
		})
		return
	}
	// ── /Streak ─────────────────────────────────────────────

	response.Created(c, res)
}

func (h *PomodoroHandler) List(c *gin.Context) {
	var filter dto.PomodoroListFilter
	c.ShouldBindQuery(&filter) //nolint
	filter.PaginationReq = dto.ParsePagination(c)

	res, err := h.svc.List(c.Request.Context(), middleware.GetUserID(c), filter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *PomodoroHandler) GetStats(c *gin.Context) {
	fromStr := c.DefaultQuery("from", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	toStr := c.DefaultQuery("to", time.Now().Format("2006-01-02"))

	from, _ := time.Parse("2006-01-02", fromStr)
	to, _ := time.Parse("2006-01-02", toStr)
	to = to.Add(24*time.Hour - time.Second) // günün sonuna kadar

	res, err := h.svc.GetStats(c.Request.Context(), middleware.GetUserID(c), from, to)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *PomodoroHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "geçersiz ID")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id, middleware.GetUserID(c)); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}
