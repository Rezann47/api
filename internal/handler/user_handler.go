package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/Rezann47/YksKoc/internal/domain/dto"
	"github.com/Rezann47/YksKoc/internal/middleware"
	"github.com/Rezann47/YksKoc/internal/service"
	"github.com/Rezann47/YksKoc/pkg/response"
)

type UserHandler struct{ svc service.UserService }

func NewUserHandler(svc service.UserService) *UserHandler { return &UserHandler{svc: svc} }

func (h *UserHandler) GetProfile(c *gin.Context) {
	res, err := h.svc.GetProfile(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	res, err := h.svc.UpdateProfile(c.Request.Context(), middleware.GetUserID(c), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	if err := h.svc.ChangePassword(c.Request.Context(), middleware.GetUserID(c), req); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

func (h *UserHandler) GetPremiumStatus(c *gin.Context) {
	res, err := h.svc.GetPremiumStatus(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *UserHandler) ActivatePremium(c *gin.Context) {
	res, err := h.svc.ActivatePremium(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}
