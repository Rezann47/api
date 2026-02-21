package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/Rezann47/YksKoc/internal/domain/dto"
	"github.com/Rezann47/YksKoc/internal/middleware"
	"github.com/Rezann47/YksKoc/internal/service"
	"github.com/Rezann47/YksKoc/pkg/response"
)

type AuthHandler struct{ svc service.AuthService }

func NewAuthHandler(svc service.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

// Register godoc
// @Summary      Kayıt ol
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.RegisterReq true "Kayıt bilgileri"
// @Success      201  {object} dto.LoginRes
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	res, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, res)
}

// Login godoc
// @Summary      Giriş yap
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.LoginReq true "Giriş bilgileri"
// @Success      200  {object} dto.LoginRes
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	res, err := h.svc.Login(
		c.Request.Context(), req,
		c.Request.UserAgent(), c.ClientIP(),
	)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

// Refresh godoc
// @Summary      Token yenile
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.RefreshReq true "Refresh token"
// @Success      200  {object} dto.TokenRes
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	res, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

// Logout godoc
// @Summary      Çıkış yap
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.RefreshReq true "Refresh token"
// @Success      204
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	if err := h.svc.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

// LogoutAll godoc
// @Summary      Tüm cihazlardan çıkış
// @Tags         auth
// @Security     BearerAuth
// @Success      204
// @Router       /auth/logout-all [post]
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.svc.LogoutAll(c.Request.Context(), userID); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}
