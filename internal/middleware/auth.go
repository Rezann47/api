package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Rezann47/YksKoc/internal/config"
	"github.com/Rezann47/YksKoc/internal/domain/apperror"
	"github.com/Rezann47/YksKoc/pkg/jwt"
	"github.com/Rezann47/YksKoc/pkg/response"
)

const (
	CtxUserID = "userID"
	CtxRole   = "role"
)

func Auth(jwtCfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, apperror.NewUnauthorized("token gerekli"))
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwt.ValidateAccess(tokenStr, jwtCfg.AccessSecret)
		if err != nil {
			response.Error(c, apperror.NewUnauthorized("geçersiz veya süresi dolmuş token"))
			c.Abort()
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role, _ := c.Get(CtxRole)
		if !allowed[role.(string)] {
			response.Error(c, apperror.NewForbidden("bu işlem için yetkiniz yok"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetUserID(c *gin.Context) uuid.UUID {
	v, exists := c.Get(CtxUserID)
	if !exists {
		return uuid.Nil
	}

	switch id := v.(type) {
	case uuid.UUID:
		return id
	case string:
		u, err := uuid.Parse(id)
		if err != nil {
			return uuid.Nil
		}
		return u
	default:
		return uuid.Nil
	}
}

func GetRole(c *gin.Context) string {
	v, _ := c.Get(CtxRole)
	role, _ := v.(string)
	return role
}
