package jwt

import (
	"crypto/sha256"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwtlib.RegisteredClaims
}

type RefreshClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwtlib.RegisteredClaims
}

type ParsedRefreshToken struct {
	ExpiresAt time.Time
}

// GenerateAccess kısa ömürlü access token üretir
func GenerateAccess(userID uuid.UUID, role, secret string, expiry time.Duration) (string, error) {
	claims := AccessClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}
	return jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// GenerateRefresh uzun ömürlü refresh token üretir
// raw (kullanıcıya gönderilecek), parsed (DB'ye yazılacak meta) döner
func GenerateRefresh(userID uuid.UUID, secret string, expiry time.Duration) (raw string, parsed *ParsedRefreshToken, err error) {
	expiresAt := time.Now().Add(expiry)
	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(expiresAt),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}
	raw, err = jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", nil, err
	}
	return raw, &ParsedRefreshToken{ExpiresAt: expiresAt}, nil
}

// ValidateAccess access token'ı doğrular ve claims döner
func ValidateAccess(tokenStr, secret string) (*AccessClaims, error) {
	token, err := jwtlib.ParseWithClaims(tokenStr, &AccessClaims{}, keyFunc(secret))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

// ValidateRefresh refresh token'ı doğrular
func ValidateRefresh(tokenStr, secret string) (*RefreshClaims, error) {
	token, err := jwtlib.ParseWithClaims(tokenStr, &RefreshClaims{}, keyFunc(secret))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh claims")
	}
	return claims, nil
}

// HashToken token'ın SHA-256 hash'ini döner (DB'de saklanır)
func HashToken(raw string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(raw)))
}

func keyFunc(secret string) jwtlib.Keyfunc {
	return func(token *jwtlib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}
}
