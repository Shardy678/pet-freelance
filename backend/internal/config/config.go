package config

import "os"

type AppConfig struct {
	DSN       string
	JWTSecret string
}

func Load() *AppConfig {
	return &AppConfig{
		DSN:       getenv("DATABASE_URL", "host=localhost user=app dbname=app sslmode=disable"),
		JWTSecret: getenv("JWT_SECRET", "dev‑only‑secret"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
