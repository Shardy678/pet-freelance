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

	// Repositories & services
	userRepo := repository.NewUserRepository(db.DB)
	authSvc := service.NewAuthService(userRepo, cfg)
	offerRepo := repository.NewServiceOfferRepository(db.DB)
	serviceRepo := repository.NewServiceRepository(db.DB)
	slotRepo := repository.NewAvailabilitySlotRepository(db.DB)

	// Handlers
	authH := handlers.NewAuthHandler(authSvc)
	profH := handlers.NewProfileHandler(userRepo)
	offerH := handlers.NewServiceOfferHandler(offerRepo)
	serviceH := handlers.NewServiceHandler(service.NewServiceService(serviceRepo))
	slotH := handlers.NewAvailabilitySlotHandler(service.NewAvailabilitySlotService(slotRepo))

	api := r.Group("/api")
	{
		api.GET("/health", handlers.HealthCheck)
		api.GET("/hello", handlers.HelloWorld)

		// Auth & profile
		api.POST("/auth/register", authH.Register)
		api.POST("/auth/login", authH.Login)

		// Protected routes
		secure := api.Group("/")
		secure.Use(middleware.JWT(cfg))
		{
			secure.GET("/profile/me", profH.Me)
			secure.POST("/offers", offerH.Create)
			secure.POST("/services", serviceH.Create)
		}

		// Public services
		api.GET("/services", serviceH.List)
		api.GET("/services/:id", serviceH.Get)

		// Offers & nested slots
		offers := api.Group("/offers")
		{
			offers.GET("", offerH.List)

			specific := offers.Group("/:offer_id")
			{
				// GET  /api/offers/:offer_id
				specific.GET("", offerH.Get)

				// GET  /api/offers/:offer_id/slots
				specific.GET("/slots", slotH.List)

				// POST /api/offers/:offer_id/slots (protected)
				specificAuth := specific.Group("")
				specificAuth.Use(middleware.JWT(cfg))
				{
					specificAuth.POST("/slots", slotH.Create)
				}
			}
		}

		// Slots by ID (update/delete)
		slotsByID := api.Group("/slots")
		slotsByID.Use(middleware.JWT(cfg))
		{
			slotsByID.PUT("/:slot_id", slotH.Update)
			slotsByID.DELETE("/:slot_id", slotH.Delete)
		}
	}
}
