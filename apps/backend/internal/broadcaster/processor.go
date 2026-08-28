package broadcaster

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"

	"github.com/dis70rt/flowback/internal/decline"
	"github.com/dis70rt/flowback/internal/events"
	"github.com/dis70rt/flowback/internal/pubsub"
)

type Processor struct {
	bus pubsub.Publisher // Dependency injected publisher
}

func NewProcessor(bus pubsub.Publisher) *Processor {
	return &Processor{bus: bus}
}

func (p *Processor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload events.WebhookPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Printf("ERROR: Failed to decode task payload: %v\n", err)
		return err
	}

	log.Printf("CONSUMED: Event=%s\n", payload.Event)

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
			p.bus.Publish(ctx, events.TopicDeclineHard, result)
			// -> Customer needs to provide new method
			// -> Log halted in DB

		case decline.Soft:
			log.Printf("-> BROADCAST [%s] PaymentID=%s\n", events.TopicDeclineSoft, result.PaymentID)
			p.bus.Publish(ctx, events.TopicDeclineSoft, result)
			// -> Schedule retry
			// -> Log retry attempt in DB

		case decline.Unknown:
			log.Printf("-> BROADCAST [%s] PaymentID=%s\n", events.TopicDeclineUnknown, result.PaymentID)
			p.bus.Publish(ctx, events.TopicDeclineUnknown, result)
			// -> AI agent decides
		}

	case events.EventSubscriptionHalted:
		log.Printf("-> Subscription halted event received (raw=%d bytes)\n", len(payload.RawJSON))
		// TODO: Trigger dunning flow

	case events.EventSubscriptionPending:
		log.Printf("-> Subscription pending event received (raw=%d bytes)\n", len(payload.RawJSON))
		// TODO: Razorpay is retrying

	case events.EventPaymentCaptured:
		log.Printf("-> Payment captured event received (raw=%d bytes)\n", len(payload.RawJSON))
		// TODO: Mark invoice paid

	default:
		log.Printf("-> Unhandled event type=%q (raw=%d bytes) — skipping\n", payload.Event, len(payload.RawJSON))
	}

	return nil
}
