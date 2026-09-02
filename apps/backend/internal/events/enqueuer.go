package events

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/hibiken/asynq"
)

type Enqueuer struct {
	client *asynq.Client
}

func NewEnqueuer(client *asynq.Client) *Enqueuer {
	return &Enqueuer{client: client}
}

func (e *Enqueuer) EnqueueWebhook(event string, rawJSON []byte) error {
	payload := WebhookPayload{
		Event:   event,
		RawJSON: rawJSON,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(
		TopicWebhookReceived,
		payloadBytes,
		asynq.MaxRetry(5),
		asynq.Timeout(5*time.Minute),  // Each attempt must finish within 5mins``
		asynq.Retention(24*time.Hour), // Keep completed tasks for 24h for debugging
	)

	info, err := e.client.Enqueue(task)
	if err != nil {
		slog.Error("could not enqueue webhook task", "error", err)
		return err
	}

	slog.Info("task enqueued", "task_id", info.ID, "queue", info.Queue, "topic", TopicWebhookReceived, "event", event)
	return nil
}
