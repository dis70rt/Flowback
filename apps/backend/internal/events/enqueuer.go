package events

import (
	"encoding/json"
	"log"

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

	task := asynq.NewTask(TopicWebhookReceived, payloadBytes, asynq.MaxRetry(5))
	
	info, err := e.client.Enqueue(task)
	if err != nil {
		log.Printf("ERROR: Could not enqueue webhook task: %v\n", err)
		return err
	}

	log.Printf("ENQUEUED: Task %s | Queue: %s | Topic: %s | Event: %s\n", info.ID, info.Queue, TopicWebhookReceived, event)
	return nil
}
