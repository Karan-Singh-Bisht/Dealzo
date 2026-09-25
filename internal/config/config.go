package config

import (
	"os"

	"github.com/joho/godotenv"
)

//If we want to export anything it's name should be capital

type Config struct {
	Port         string
	Env          string
	Database_url string
}

// Follows Fail Fast Pattern
func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is required")
	}

	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is required")
	}

	database_url := os.Getenv("DATABASE_URL")
	if database_url == "" {
		panic("DATABASE_URL is required")
	}

	return Config{
		Port:         port,
		Env:          env,
		Database_url: database_url,
	}
}
