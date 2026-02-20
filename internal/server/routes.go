package server

import (
	"github.com/gin-gonic/gin"

	"go-gin-crud/config"
	"go-gin-crud/internal/handler"
	"go-gin-crud/internal/middleware"
)

func registerRoutes(
	r *gin.Engine,
	cfg *config.Config,
	userHandler *handler.UserHandler,
	productHandler *handler.ProductHandler,
) {
	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/me", userHandler.Me)

		users := protected.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		products := protected.Group("/products")
		{
			products.POST("", productHandler.Create)
			products.GET("", productHandler.GetAll)
			products.GET("/:id", productHandler.GetOne)
			products.PUT("/:id", productHandler.Update)
			products.DELETE("/:id", productHandler.Delete)
		}
	}
}
