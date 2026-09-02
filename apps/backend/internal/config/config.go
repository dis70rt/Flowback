package config

import (
	"log/slog"
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
	TreeLog           bool // Controlled by FLOWBACK_TREE_LOG=true
}

func Load() *Config {
	if err := godotenv.Load("../../.env"); err != nil {
		_ = godotenv.Load(".env") // Fallback
	}

	webhookSecret := os.Getenv("RAZORPAY_WEBHOOK_SECRET")
	if webhookSecret == "" {
		slog.Info("WARNING: RAZORPAY_WEBHOOK_SECRET is not set in .env")
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
		slog.Info("WARNING: OPENROUTER_API_KEY is not set in .env")
	}
	openRouterModel := os.Getenv("OPENROUTER_MODEL")
	if openRouterModel == "" {
		openRouterModel = "z-ai/glm-5.3-flash"
	}

	treeLog := os.Getenv("FLOWBACK_TREE_LOG") == "true"

	return &Config{
		RazorpaySecret:    webhookSecret,
		RazorpayKeyID:     keyID,
		RazorpayKeySecret: keySecret,
		DatabaseURL:       dbUrl,
		RedisAddr:         redisAddr,
		RedisURL:          redisURL,
		OpenRouterAPIKey:  openRouterAPIKey,
		OpenRouterModel:   openRouterModel,
		TreeLog:           treeLog,
	}
}
