package main

import (
	"log"

	"github.com/dis70rt/flowback/internal/api"
	"github.com/dis70rt/flowback/internal/config"
	"github.com/dis70rt/flowback/internal/database"
	"github.com/dis70rt/flowback/internal/events"
	"github.com/dis70rt/flowback/internal/pubsub"
	"github.com/dis70rt/flowback/internal/repo"
)

func main() {
	cfg := config.Load()
	
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("FATAL: %v", err)
	}
	defer db.Close()

	queries := repo.New(db)

	asynqClient := database.InitAsynqClient(cfg.RedisAddr)
	defer asynqClient.Close()
	enqueuer := events.NewEnqueuer(asynqClient)

	bus, err := pubsub.New(cfg.RedisURL)
	if err != nil {
		log.Fatalf("FATAL: failed to init pubsub: %v", err)
	}
	defer bus.Close()

	router := api.NewRouter(api.RouterDeps{
		Queries:        queries,
		Enqueuer:       enqueuer,
		Bus:            bus,
		RazorpaySecret: cfg.RazorpaySecret,
	})

	log.Println("STARTING: Flowback Backend listening on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
