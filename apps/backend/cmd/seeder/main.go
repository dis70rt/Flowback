package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/dis70rt/flowback/internal/seed"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	taskType := flag.String("task", "db", "Which seeder task to run: 'db', 'vip', 'genz', 'outage', 'ignorer', 'fraud', 'night'")
	flag.Parse()

	_ = godotenv.Load("../../.env")

	// 1. Database Seeding Task
	if *taskType == "db" {
		dsn := "postgres://flowback:postgres@localhost:5432/flowback?sslmode=disable"
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal("Failed to open DB:", err)
		}
		defer db.Close()
		
		seed.RunDBSeeder(db)
		return
	}

	// 2. Razorpay API Tasks
	keyID := os.Getenv("RAZORPAY_KEY_ID")
	secret := os.Getenv("RAZORPAY_KEY_SECRET")
	if keyID == "" || secret == "" {
		log.Fatal("[FATAL] Missing Razorpay Keys in .env")
	}

	engine := seed.New(keyID, secret)
	var task seed.Task

	switch *taskType {
	case "vip":
		task = &seed.FailedCardTask{
			CustomerID:  "cust_vip",
			Email:       "rajesh@vip.com",
			AmountPaise: 500000,
			CardNumber:  seed.CardErrorInsufficientFund,
		}
	case "genz":
		task = &seed.FailedCardTask{
			CustomerID:  "cust_genz",
			Email:       "aditi@genz.edu",
			AmountPaise: 15000,
			CardNumber:  seed.CardErrorInsufficientFund,
		}
	case "outage":
		task = &seed.FailedCardTask{
			CustomerID:  "cust_outage",
			Email:       "vikram@nepal.com",
			AmountPaise: 90000,
			CardNumber:  seed.CardErrorInsufficientFund,
		}
	case "ignorer":
		task = &seed.FailedCardTask{
			CustomerID:  "cust_ignorer",
			Email:       "priya@ignorer.com",
			AmountPaise: 25000,
			CardNumber:  seed.CardErrorInsufficientFund,
		}
	case "fraud":
		task = &seed.FailedCardTask{
			CustomerID:  "cust_fraud",
			Email:       "scammer@fraud.com",
			AmountPaise: 999000,
			CardNumber:  seed.CardErrorInsufficientFund,
		}
	case "night":
		task = &seed.FailedCardTask{
			CustomerID:  "cust_night",
			Email:       "rahul@night.com",
			AmountPaise: 30000,
			CardNumber:  seed.CardErrorInsufficientFund,
		}
	default:
		log.Fatalf("Unknown task: %s. Use 'db', 'vip', 'genz', 'outage', 'ignorer', 'fraud', or 'night'", *taskType)
	}

	engine.Execute(task)
}
