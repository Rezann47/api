package dto

// ==================== USER DTOs ====================

// CreateUserRequest — POST /users body
type CreateUserRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"     binding:"omitempty,oneof=user admin"`
}

// UpdateUserRequest — PUT /users/:id body (tüm alanlar opsiyonel)
type UpdateUserRequest struct {
	Name     string `json:"name"      binding:"omitempty,min=2,max=100"`
	Email    string `json:"email"     binding:"omitempty,email"`
	Password string `json:"password"  binding:"omitempty,min=6"`
	Role     string `json:"role"      binding:"omitempty,oneof=user admin"`
	IsActive *bool  `json:"is_active"` // pointer: false değerini de alabilsin
}

// UserResponse — kullanıcıya dönecek temiz response (şifre yok)
type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

// LoginRequest — POST /auth/login
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse — JWT token response
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt string       `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// ==================== PRODUCT DTOs ====================

type CreateProductRequest struct {
	Name        string  `json:"name"        binding:"required,min=2,max=200"`
	Description string  `json:"description" binding:"omitempty"`
	Price       float64 `json:"price"       binding:"required,gt=0"`
	Stock       int     `json:"stock"       binding:"omitempty,gte=0"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"        binding:"omitempty,min=2,max=200"`
	Description string  `json:"description" binding:"omitempty"`
	Price       float64 `json:"price"       binding:"omitempty,gt=0"`
	Stock       int     `json:"stock"       binding:"omitempty,gte=0"`
}

// ==================== GENEL ====================

// PaginationQuery — GET /users?page=1&limit=10&sort=created_at&order=desc
type PaginationQuery struct {
	Page   int    `form:"page"   binding:"omitempty,min=1"`
	Limit  int    `form:"limit"  binding:"omitempty,min=1,max=100"`
	Sort   string `form:"sort"   binding:"omitempty"`
	Order  string `form:"order"  binding:"omitempty,oneof=asc desc"`
	Search string `form:"search" binding:"omitempty"`
}

// PaginatedResponse — sayfalı liste response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// APIResponse — standart API yanıtı
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Başarılı response oluşturucu
func Success(data interface{}, message string) APIResponse {
	return APIResponse{Success: true, Message: message, Data: data}
}

// Hata response oluşturucu
func Fail(err string) APIResponse {
	return APIResponse{Success: false, Error: err}
}
