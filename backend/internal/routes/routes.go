package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shardy678/pet-freelance/backend/internal/handlers"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/health", handlers.HealthCheck)
		api.GET("/hello", handlers.HelloWorld)
	}
}
