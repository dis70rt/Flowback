package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RazorpaySecret    string
	RazorpayKeyID     string
	RazorpayKeySecret string
	DatabaseURL       string
	RedisAddr         string
	RedisURL          string
	OpenRouterAPIKey  string
	OpenRouterModel   string
}

func Load() *Config {
	if err := godotenv.Load("../../.env"); err != nil {
		_ = godotenv.Load(".env") // Fallback
	}

	webhookSecret := os.Getenv("RAZORPAY_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Println("WARNING: RAZORPAY_WEBHOOK_SECRET is not set in .env")
	}

	keyID := os.Getenv("RAZORPAY_KEY_ID")
	keySecret := os.Getenv("RAZORPAY_KEY_SECRET")

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "postgres://flowback:postgres@localhost:5432/flowback?sslmode=disable"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisURL := "redis://" + redisAddr + "/0"

	openRouterAPIKey := os.Getenv("OPENROUTER_API_KEY")
	if openRouterAPIKey == "" {
		log.Println("WARNING: OPENROUTER_API_KEY is not set in .env")
	}
	openRouterModel := os.Getenv("OPENROUTER_MODEL")
	if openRouterModel == "" {
		openRouterModel = "z-ai/glm-5.3-flash"
	}

	return &Config{
		RazorpaySecret:    webhookSecret,
		RazorpayKeyID:     keyID,
		RazorpayKeySecret: keySecret,
		DatabaseURL:       dbUrl,
		RedisAddr:         redisAddr,
		RedisURL:          redisURL,
		OpenRouterAPIKey:  openRouterAPIKey,
		OpenRouterModel:   openRouterModel,
	}
}
