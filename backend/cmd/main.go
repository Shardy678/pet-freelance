package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/shardy678/pet-freelance/backend/internal/db"
	"github.com/shardy678/pet-freelance/backend/internal/routes"
)

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Println("No .env file found, continuing with environment variables…")
	}

	db.Init()

	router := gin.Default()
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
