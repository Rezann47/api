package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-gin-crud/internal/dto"
	"go-gin-crud/internal/service"
)

// ==================== USER HANDLER ====================

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Register godoc
// @Summary     Yeni kullanıcı kaydı
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body  body      dto.CreateUserRequest  true  "Kullanıcı bilgileri"
// @Success     201   {object}  dto.APIResponse
// @Failure     400   {object}  dto.APIResponse
// @Router      /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(err.Error()))
		return
	}

	user, err := h.svc.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(user, "Kayıt başarılı"))
}

// Login godoc
// @Summary     Giriş yap ve JWT al
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body  body      dto.LoginRequest  true  "Giriş bilgileri"
// @Success     200   {object}  dto.APIResponse
// @Failure     401   {object}  dto.APIResponse
// @Router      /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(err.Error()))
		return
	}

	resp, err := h.svc.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(resp, "Giriş başarılı"))
}

// GetUsers godoc
// @Summary     Tüm kullanıcıları listele (sayfalı)
// @Tags        users
// @Security    BearerAuth
// @Produce     json
// @Param       page    query     int     false  "Sayfa no"
// @Param       limit   query     int     false  "Sayfa başı kayıt"
// @Param       search  query     string  false  "Ad/email arama"
// @Param       sort    query     string  false  "Sıralama alanı"
// @Param       order   query     string  false  "asc veya desc"
// @Success     200     {object}  dto.APIResponse
// @Router      /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(err.Error()))
		return
	}

	result, err := h.svc.GetAll(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(result, ""))
}

// GetUser godoc
// @Summary  Tek kullanıcı getir
// @Tags     users
// @Security BearerAuth
// @Produce  json
// @Param    id   path      int  true  "Kullanıcı ID"
// @Success  200  {object}  dto.APIResponse
// @Failure  404  {object}  dto.APIResponse
// @Router   /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail("Geçersiz ID"))
		return
	}

	user, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(user, ""))
}

// UpdateUser godoc
// @Summary  Kullanıcı güncelle
// @Tags     users
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id    path      int                    true  "Kullanıcı ID"
// @Param    body  body      dto.UpdateUserRequest  true  "Güncellenecek alanlar"
// @Success  200   {object}  dto.APIResponse
// @Router   /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail("Geçersiz ID"))
		return
	}

	// Kullanıcı sadece kendi profilini güncelleyebilir (admin hariç)
	currentUserID, _ := c.Get("user_id")
	currentRole, _ := c.Get("role")
	if currentRole != "admin" && currentUserID.(uint) != id {
		c.JSON(http.StatusForbidden, dto.Fail("Sadece kendi profilinizi güncelleyebilirsiniz"))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(err.Error()))
		return
	}

	user, err := h.svc.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(user, "Kullanıcı güncellendi"))
}

// DeleteUser godoc
// @Summary  Kullanıcı sil (soft delete)
// @Tags     users
// @Security BearerAuth
// @Produce  json
// @Param    id   path      int  true  "Kullanıcı ID"
// @Success  200  {object}  dto.APIResponse
// @Router   /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail("Geçersiz ID"))
		return
	}

	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil, "Kullanıcı silindi"))
}

// Me — giriş yapmış kullanıcının profilini döner
func (h *UserHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.svc.GetByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(user, ""))
}

// ==================== PRODUCT HANDLER ====================

type ProductHandler struct {
	svc service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")
	product, err := h.svc.Create(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(product, "Ürün oluşturuldu"))
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	var query dto.PaginationQuery
	c.ShouldBindQuery(&query)

	userID, _ := c.Get("user_id")
	result, err := h.svc.GetAll(userID.(uint), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(result, ""))
}

func (h *ProductHandler) GetOne(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail("Geçersiz ID"))
		return
	}

	product, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(product, ""))
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail("Geçersiz ID"))
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")
	product, err := h.svc.Update(id, userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(product, "Ürün güncellendi"))
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail("Geçersiz ID"))
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.svc.Delete(id, userID.(uint)); err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil, "Ürün silindi"))
}

// ==================== YARDIMCI ====================

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	return uint(id), err
}
