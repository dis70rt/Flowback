package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RazorpaySecret string
	DatabaseURL    string
}

func Load() *Config {
	if err := godotenv.Load("../../.env"); err != nil {
		_ = godotenv.Load(".env") // Fallback
	}

	secret := os.Getenv("RAZORPAY_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("FATAL: RAZORPAY_WEBHOOK_SECRET is not set in .env")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		// Default to localhost if running outside docker
		dbUrl = "postgres://flowback:postgres@localhost:5432/flowback?sslmode=disable"
	}

	return &Config{
		RazorpaySecret: secret,
		DatabaseURL:    dbUrl,
	}
}
