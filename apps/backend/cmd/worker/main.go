package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"

	"github.com/dis70rt/flowback/internal/config"
	"github.com/dis70rt/flowback/internal/decline"
	"github.com/dis70rt/flowback/internal/events"
)

func HandleWebhook(ctx context.Context, t *asynq.Task) error {
	var payload events.WebhookPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Printf("ERROR: Failed to decode task payload: %v\n", err)
		return err
	}

	log.Printf("CONSUMED: Event=%s\n", payload.Event)

	// Decide what to do based on the event type
	switch payload.Event {

	case events.EventPaymentFailed:
		result, err := decline.ClassifyPayload(payload.RawJSON)
		if err != nil {
			log.Printf("ERROR: Failed to classify payload: %v\n", err)
			return err
		}

		log.Printf("-> Payment %s | reason=%q | classification=%s\n",
			result.PaymentID, result.ErrorReason, result.Type)

		switch result.Type {

		case decline.Hard:
			log.Printf("-> BROADCAST [%s] PaymentID=%s\n", events.TopicDeclineHard, result.PaymentID)
			// TODO: broadcaster.Publish(ctx, events.TopicDeclineHard, result)
			// -> Notify customer to provide a new payment method
			// -> Mark subscription as halted in audit log

		case decline.Soft:
			log.Printf("-> BROADCAST [%s] PaymentID=%s\n", events.TopicDeclineSoft, result.PaymentID)
			// TODO: broadcaster.Publish(ctx, events.TopicDeclineSoft, result)
			// -> Schedule a retry with exponential backoff
			// -> Log soft decline attempt count in audit log

		case decline.Unknown:
			log.Printf("-> BROADCAST [%s] PaymentID=%s\n", events.TopicDeclineUnknown, result.PaymentID)
			// TODO: broadcaster.Publish(ctx, events.TopicDeclineUnknown, result)
			// -> Forward to AI agent to decide Hard or Soft
			// -> AI agent re-publishes to TopicDeclineHard or TopicDeclineSoft
		}

	case events.EventSubscriptionHalted:
		log.Printf("-> Subscription halted event received (raw=%d bytes)\n", len(payload.RawJSON))
		// TODO: Handle subscription halted — trigger customer dunning flow

	case events.EventSubscriptionPending:
		log.Printf("-> Subscription pending event received (raw=%d bytes)\n", len(payload.RawJSON))
		// TODO: Handle subscription pending — Razorpay is already retrying

	case events.EventPaymentCaptured:
		log.Printf("-> Payment captured event received (raw=%d bytes)\n", len(payload.RawJSON))
		// TODO: Mark invoice as paid in audit log

	default:
		log.Printf("-> Unhandled event type=%q (raw=%d bytes) — skipping\n", payload.Event, len(payload.RawJSON))
	}

	// Return nil to tell Asynq the task was processed successfully
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
