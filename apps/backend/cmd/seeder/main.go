package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/dis70rt/flowback/internal/seed"
)

func main() {
	// Changed default from "upi" to "card" since UPI requires manual dashboard activation
	taskType := flag.String("task", "card", "Which seeder task to run: 'upi' or 'card'")
	flag.Parse()

	// 1. Load Environment Variables
	_ = godotenv.Load("../../.env")
	keyID := os.Getenv("RAZORPAY_KEY_ID")
	secret := os.Getenv("RAZORPAY_KEY_SECRET")

	if keyID == "" || secret == "" {
		log.Fatal("[FATAL] Missing Razorpay Keys in .env")
	}

	// 2. Initialize the Seeder Engine
	engine := seed.New(keyID, secret)

	// 3. Select and Execute the Interface Task
	var task seed.Task

	switch *taskType {
	case "card":
		task = &seed.FailedCardTask{AmountPaise: 50000, CardNumber: seed.CardErrorInsufficientFund}
	/*
	case "upi":
		task = &seed.FailedUPICollectTask{AmountPaise: seed.UPIErrorIncorrectPIN}
	*/
	default:
		log.Fatalf("Unknown task: %s. Use 'card'", *taskType)
	}

	engine.Execute(task)
}
