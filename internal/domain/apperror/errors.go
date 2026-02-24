package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError uygulamanın standart hata tipi.
// Handler katmanı bunu HTTP response'a dönüştürür.
type AppError struct {
	Code       string // makine kodu: USER_NOT_FOUND
	Message    string // kullanıcıya gösterilecek mesaj
	HTTPStatus int
	Err        error // iç hata (log için)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Is(target error) bool {
	var t *AppError
	if errors.As(target, &t) {
		return e.Code == t.Code
	}
	return false
}

// ─── Constructor'lar ──────────────────────────────────────

func New(code, message string, status int, err error) *AppError {
	return &AppError{Code: code, Message: message, HTTPStatus: status, Err: err}
}

func NewNotFound(resource string, err error) *AppError {
	return New("NOT_FOUND", fmt.Sprintf("%s bulunamadı", resource), http.StatusNotFound, err)
}

func NewUnauthorized(msg string) *AppError {
	return New("UNAUTHORIZED", msg, http.StatusUnauthorized, nil)
}

func NewForbidden(msg string) *AppError {
	return New("FORBIDDEN", msg, http.StatusForbidden, nil)
}

func NewConflict(msg string, err error) *AppError {
	return New("CONFLICT", msg, http.StatusConflict, err)
}

func NewValidation(msg string) *AppError {
	return New("VALIDATION_ERROR", msg, http.StatusUnprocessableEntity, nil)
}

func NewInternal(err error) *AppError {
	return New("INTERNAL_ERROR", "Beklenmedik bir hata oluştu", http.StatusInternalServerError, err)
}
func NewBadRequest(msg string) *AppError {
	return New("BAD_REQUEST", msg, http.StatusBadRequest, nil)
}

// ─── Sentinel hatalar ─────────────────────────────────────

var (
	ErrNotFound     = New("NOT_FOUND", "kaynak bulunamadı", http.StatusNotFound, nil)
	ErrUnauthorized = New("UNAUTHORIZED", "kimlik doğrulama gerekli", http.StatusUnauthorized, nil)
	ErrForbidden    = New("FORBIDDEN", "bu işlem için yetkiniz yok", http.StatusForbidden, nil)
	ErrConflict     = New("CONFLICT", "kaynak zaten mevcut", http.StatusConflict, nil)
)

// Unwrap — errors.Is/As ile çalışmak için
func IsAppError(err error) bool {
	var e *AppError
	return errors.As(err, &e)
}

func HTTPStatus(err error) int {
	var e *AppError
	if errors.As(err, &e) {
		return e.HTTPStatus
	}
	return http.StatusInternalServerError
}
