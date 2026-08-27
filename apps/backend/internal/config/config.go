package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RazorpaySecret string
	DatabaseURL    string
	RedisAddr      string
	RedisURL       string
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

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisURL := "redis://" + redisAddr + "/0"

	return &Config{
		RazorpaySecret: secret,
		DatabaseURL:    dbUrl,
		RedisAddr:      redisAddr,
		RedisURL:       redisURL,
	}
}
