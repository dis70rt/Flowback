package main

import (
	"log"

	"github.com/dis70rt/flowback/internal/broadcaster"
	"github.com/dis70rt/flowback/internal/config"
	"github.com/dis70rt/flowback/internal/events"
	"github.com/dis70rt/flowback/internal/pubsub"
)

func main() {
	cfg := config.Load()

	ps, err := pubsub.New(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis Pub/Sub: %v", err)
	}
	log.Println("[WORKER] Redis Pub/Sub engine connected.")

	processor := broadcaster.NewProcessor(ps)
	consumer := events.NewConsumer(cfg.RedisAddr, 10)
	consumer.Register(events.TopicWebhookReceived, processor.ProcessTask)

	if err := consumer.Start(); err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
}
