package fixture

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Rezann47/YksKoc/internal/domain/entity"
	"github.com/Rezann47/YksKoc/pkg/password"
)

// CreateUser test kullanıcısı oluşturur
func CreateUser(db *gorm.DB, role entity.Role) *entity.User {
	hash, _ := password.Hash("Test1234!")
	code := "YKS99999"
	user := &entity.User{
		ID:           uuid.New(),
		Name:         "Test User",
		Email:        uuid.NewString() + "@test.com",
		PasswordHash: hash,
		Role:         role,
	}
	if role == entity.RoleStudent {
		user.StudentCode = &code
	}
	db.WithContext(context.Background()).Create(user)
	return user
}
