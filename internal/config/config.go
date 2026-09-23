package config

import (
	"os"

	"github.com/joho/godotenv"
)

//If we want to export anything it's name should be capital

type Config struct {
	Port string
	Env  string
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

	return Config{
		Port: port,
		Env:  env,
	}
}
