package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Rezann47/YksKoc/internal/config"
	"github.com/Rezann47/YksKoc/pkg/jwt"
	"github.com/Rezann47/YksKoc/pkg/response"
)

const (
	CtxUserID = "userID"
	CtxRole   = "role"
)

// Auth JWT access token doğrular ve userID + role'u context'e enjekte eder
func Auth(jwtCfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, errUnauthorized())
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwt.ValidateAccess(tokenStr, jwtCfg.AccessSecret)
		if err != nil {
			response.Error(c, errUnauthorized())
			c.Abort()
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

// RequireRole belirtilen rollerden birine sahip kullanıcı geçirir
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role, _ := c.Get(CtxRole)
		if !allowed[role.(string)] {
			response.Error(c, errForbidden())
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetUserID context'ten userID çeker (handler'larda kullanılır)
func GetUserID(c *gin.Context) uuid.UUID {
	v, _ := c.Get(CtxUserID)
	id, _ := v.(uuid.UUID)
	return id
}

// GetRole context'ten role çeker
func GetRole(c *gin.Context) string {
	v, _ := c.Get(CtxRole)
	role, _ := v.(string)
	return role
}

func errUnauthorized() error {
	return &unauthorizedErr{}
}

func errForbidden() error {
	return &forbiddenErr{}
}

// basit hata tipleri response paketini tetiklemek için
type unauthorizedErr struct{}

func (e *unauthorizedErr) Error() string { return "UNAUTHORIZED" }

type forbiddenErr struct{}

func (e *forbiddenErr) Error() string { return "FORBIDDEN" }
