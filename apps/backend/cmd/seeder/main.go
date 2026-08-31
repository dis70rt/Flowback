package main

import (
	"database/sql"
	"log"

	"github.com/dis70rt/flowback/internal/config"
	"github.com/dis70rt/flowback/internal/seed"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	if cfg.RazorpayKeyID == "" || cfg.RazorpayKeySecret == "" {
		log.Fatal("[FATAL] Missing Razorpay Keys in .env (RAZORPAY_KEY_ID, RAZORPAY_KEY_SECRET)")
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	seed.RunTUI(db, cfg)
}
