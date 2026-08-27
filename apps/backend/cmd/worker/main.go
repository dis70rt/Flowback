package main

import (
	"context"
	"log"

	"github.com/hibiken/asynq"

	"github.com/dis70rt/flowback/internal/config"
	"github.com/dis70rt/flowback/internal/events"
)

func HandleWebhook(ctx context.Context, t *asynq.Task) error {
	rawBytes := t.Payload()
	log.Printf("CONSUMED TASK | Topic: %s | Payload: %s\n", t.Type(), string(rawBytes))
	return nil
}

func main() {
	cfg := config.Load()

	consumer := events.NewConsumer(cfg.RedisAddr, 10)
	consumer.Register(events.TopicWebhookReceived, HandleWebhook)

	if err := consumer.Start(); err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
}
