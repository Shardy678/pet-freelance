package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shardy678/pet-freelance/backend/internal/config"
	"github.com/shardy678/pet-freelance/backend/internal/db"
	"github.com/shardy678/pet-freelance/backend/internal/handlers"
	"github.com/shardy678/pet-freelance/backend/internal/middleware"
	"github.com/shardy678/pet-freelance/backend/internal/repository"
	"github.com/shardy678/pet-freelance/backend/internal/service"
)

func SetupRoutes(r *gin.Engine) {
	cfg := config.Load()

	userRepo := repository.NewUserRepository(db.DB)
	authSvc := service.NewAuthService(userRepo, cfg)

	authH := handlers.NewAuthHandler(authSvc)
	profH := handlers.NewProfileHandler(userRepo)

	offerRepo := repository.NewServiceOfferRepository(db.DB)
	offerH := handlers.NewServiceOfferHandler(offerRepo)

	serviceRepo := repository.NewServiceRepository(db.DB)
	serviceSvc := service.NewServiceService(serviceRepo)
	serviceH := handlers.NewServiceHandler(serviceSvc)

	api := r.Group("/api")
	{
		api.GET("/health", handlers.HealthCheck)
		api.GET("/hello", handlers.HelloWorld)
		api.POST("/auth/register", authH.Register)
		api.POST("/auth/login", authH.Login)
		api.GET("/offers", offerH.List)
		api.GET("/offers/:id", offerH.Get)
		api.GET("/services", serviceH.List)
		api.GET("/services/:id", serviceH.Get)

		secure := api.Group("/")
		secure.Use(middleware.JWT(cfg))
		{
			secure.GET("/profile/me", profH.Me)
			secure.POST("/offers", offerH.Create)
			secure.POST("/services", serviceH.Create)
		}
	}
}
