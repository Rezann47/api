package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/Rezann47/YksKoc/internal/middleware"
	"github.com/Rezann47/YksKoc/internal/service"
	"github.com/Rezann47/YksKoc/pkg/response"
)

type StreakHandler struct {
	svc *service.StreakService
}

func NewStreakHandler(svc *service.StreakService) *StreakHandler {
	return &StreakHandler{svc: svc}
}

// GetMyStreak godoc
// @Summary      Streak bilgisi
// @Tags         streak
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /streak/me [get]
func (h *StreakHandler) GetMyStreak(c *gin.Context) {
	userID := middleware.GetUserID(c)
	info, err := h.svc.GetMyStreak(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, info)
}

// GetMyBadges godoc
// @Summary      Kazanılan rozetler
// @Tags         streak
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} entity.Badge
// @Router       /badges/me [get]
func (h *StreakHandler) GetMyBadges(c *gin.Context) {
	userID := middleware.GetUserID(c)
	badges, err := h.svc.GetMyBadges(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, badges)
}

// GetLeaderboard godoc
// @Summary      Haftalık liderlik tablosu
// @Tags         streak
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /leaderboard [get]
func (h *StreakHandler) GetLeaderboard(c *gin.Context) {
	userID := middleware.GetUserID(c)
	entries, myRank, err := h.svc.GetLeaderboard(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{
		"entries": entries,
		"my_rank": myRank,
	})
}
