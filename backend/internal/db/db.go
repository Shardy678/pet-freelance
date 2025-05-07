package db

import (
	"log"

	"github.com/shardy678/pet-freelance/backend/internal/config"
	"github.com/shardy678/pet-freelance/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.AppConfig) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
}

func Init() {
	// 1. Load config (DSN, JWT secret, etc)
	cfg := config.Load()

	// 2. Open the database
	conn, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("db.Init: failed to connect to database: %v", err)
	}

	// 3. Auto-migrate your core models
	if err := conn.AutoMigrate(
		&models.User{},
		&models.Service{},
		&models.ServiceOffer{},
	); err != nil {
		log.Fatalf("db.Init: auto-migrate failed: %v", err)
	}

	// 4. Assign to global
	DB = conn
}
