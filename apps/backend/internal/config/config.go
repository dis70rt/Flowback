package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RazorpaySecret string
}

func Load() *Config {
	if err := godotenv.Load("../../.env"); err != nil {
		_ = godotenv.Load(".env") // Fallback
	}

	secret := os.Getenv("RAZORPAY_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("FATAL: RAZORPAY_WEBHOOK_SECRET is not set in .env")
	}

	return &Config{
		RazorpaySecret: secret,
	}
}
