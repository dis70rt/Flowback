package broadcaster

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/hibiken/asynq"

	"github.com/dis70rt/flowback/internal/decline"
	"github.com/dis70rt/flowback/internal/events"
	"github.com/dis70rt/flowback/internal/pubsub"
)

type Processor struct {
	bus pubsub.Publisher
}

func NewProcessor(bus pubsub.Publisher) *Processor {
	return &Processor{
		bus: bus,
	}
}

func (p *Processor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload events.WebhookPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		slog.ErrorContext(ctx, "failed to decode task payload", "error", err)
		return err
	}

	slog.InfoContext(ctx, "consumed event", "event_type", payload.Event)

	switch payload.Event {

	case events.EventPaymentFailed:
		result, err := decline.ClassifyPayload(payload.RawJSON)
		if err != nil {
			slog.ErrorContext(ctx, "failed to classify payload", "error", err)
			return err
		}

		slog.InfoContext(ctx, "payment classified",
			"payment_id", result.PaymentID,
			"reason", result.ErrorReason,
			"classification", result.Type,
		)

		switch result.Type {

		case decline.Hard:
			slog.InfoContext(ctx, "broadcasting hard decline", "topic", events.TopicDeclineHard, "payment_id", result.PaymentID)
			p.bus.Publish(ctx, events.TopicDeclineHard, result)

		case decline.Soft:
			slog.InfoContext(ctx, "broadcasting soft decline", "topic", events.TopicDeclineSoft, "payment_id", result.PaymentID)
			p.bus.Publish(ctx, events.TopicDeclineSoft, result)

		case decline.Unknown:
			slog.InfoContext(ctx, "broadcasting unknown decline", "topic", events.TopicDeclineUnknown, "payment_id", result.PaymentID)
			p.bus.Publish(ctx, events.TopicDeclineUnknown, result)
		}

	case events.EventSubscriptionHalted:
		slog.InfoContext(ctx, "subscription halted event received", "bytes", len(payload.RawJSON))

	case events.EventSubscriptionPending:
		slog.InfoContext(ctx, "subscription pending event received", "bytes", len(payload.RawJSON))

	case events.EventPaymentCaptured:
		slog.InfoContext(ctx, "payment captured event received", "bytes", len(payload.RawJSON))

	default:
		slog.InfoContext(ctx, "unhandled event type skipped", "event_type", payload.Event, "bytes", len(payload.RawJSON))
	}

	return nil
}
