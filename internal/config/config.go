package config

import (
	"os"

	"github.com/joho/godotenv"
)

//If we want to export anything it's name should be capital

type Config struct {
	Port                string
	Env                 string
	Database_url        string
	JWT_token           string
	StorageAccountID    string
	StorageAccessKey    string
	StorageAccessSecret string
	StorageBucket       string
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

	jwt_token := os.Getenv("JWT_TOKEN")
	if jwt_token == "" {
		panic("JWT_TOKEN is required")
	}

	storageAccountID := os.Getenv("STORAGE_ACCOUNT_ID")
	if storageAccountID == "" {
		panic("STORAGE_ACCOUNT_ID is required")
	}

	storageAccessKey := os.Getenv("STORAGE_ACCESS_KEY")
	if storageAccessKey == "" {
		panic("STORAGE_ACCESS_KEY is required")
	}

	storageAccessSecret := os.Getenv("STORAGE_ACCESS_SECRET")
	if storageAccessSecret == "" {
		panic("STORAGE_ACCESS_SECRET is required")
	}

	storageBucket := os.Getenv("STORAGE_BUCKET")
	if storageBucket == "" {
		panic("STORAGE_BUCKET is required")
	}

	return Config{
		Port:                port,
		Env:                 env,
		Database_url:        database_url,
		JWT_token:           jwt_token,
		StorageAccountID:    storageAccountID,
		StorageAccessKey:    storageAccessKey,
		StorageAccessSecret: storageAccessSecret,
		StorageBucket:       storageBucket,
	}
}
