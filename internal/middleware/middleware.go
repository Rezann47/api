package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"go-gin-crud/config"
	"go-gin-crud/internal/dto"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail("Authorization header eksik"))
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail("Geçersiz token formatı"))
			return
		}
		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail("Geçersiz veya süresi dolmuş token"))
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail("Token okunamadı"))
			return
		}
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("email", claims["email"].(string))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.Fail("Admin yetkisi gerekiyor"))
			return
		}
		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		statusColor := colorForStatus(param.StatusCode)
		methodColor := colorForMethod(param.Method)
		reset := "\033[0m"
		return fmt.Sprintf("[GIN] %v | %s%3d%s | %13v | %15s | %s%-7s%s %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, reset,
			param.Latency, param.ClientIP,
			methodColor, param.Method, reset,
			param.Path,
		)
	})
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[97;42m"
	case code >= 300 && code < 400:
		return "\033[90;47m"
	case code >= 400 && code < 500:
		return "\033[90;43m"
	default:
		return "\033[97;41m"
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return "\033[97;44m"
	case "POST":
		return "\033[97;46m"
	case "PUT", "PATCH":
		return "\033[97;43m"
	case "DELETE":
		return "\033[97;41m"
	default:
		return "\033[0m"
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

type rateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
}

var limiter = &rateLimiter{requests: make(map[string][]time.Time)}

func RateLimit(maxPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		window := now.Add(-1 * time.Minute)
		limiter.mu.Lock()
		defer limiter.mu.Unlock()
		var recent []time.Time
		for _, t := range limiter.requests[ip] {
			if t.After(window) {
				recent = append(recent, t)
			}
		}
		if len(recent) >= maxPerMinute {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.Fail("Çok fazla istek, lütfen bekleyin"))
			return
		}
		limiter.requests[ip] = append(recent, now)
		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.Fail("Sunucu hatası oluştu"))
	})
}
