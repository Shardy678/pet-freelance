package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/shardy678/pet-freelance/backend/internal/routes"
	"github.com/shardy678/pet-freelance/backend/internal/db"
)

func main() {
	_ = godotenv.Load()

	db.Init() 

	r := gin.Default()

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
