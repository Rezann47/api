package response

import (
	"net/http"

	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/gin-gonic/gin"
)

type successRes struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type errorRes struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// OK 200
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, successRes{Success: true, Data: data})
}

// Created 201
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, successRes{Success: true, Data: data})
}

// NoContent 204
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error apperror veya generic error'u uygun HTTP kodu ile döner
func Error(c *gin.Context, err error) {
	status := apperror.HTTPStatus(err)
	code := "INTERNAL_ERROR"
	msg := "Beklenmedik bir hata oluştu"

	if appErr, ok := err.(*apperror.AppError); ok {
		code = appErr.Code
		msg = appErr.Message
	}

	c.JSON(status, errorRes{
		Success: false,
		Code:    code,
		Message: msg,
	})
}

// ValidationError 422
func ValidationError(c *gin.Context, msg string) {
	c.JSON(http.StatusUnprocessableEntity, errorRes{
		Success: false,
		Code:    "VALIDATION_ERROR",
		Message: msg,
	})
}
