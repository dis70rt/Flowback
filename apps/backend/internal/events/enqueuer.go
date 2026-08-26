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

func (e *Enqueuer) EnqueuePaymentFailed(payload PaymentFailedPayload) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TopicPaymentFailed, payloadBytes, asynq.MaxRetry(5))
	
	info, err := e.client.Enqueue(task)
	if err != nil {
		log.Printf("ERROR: Could not enqueue PaymentFailed task: %v\n", err)
		return err
	}

	log.Printf("ENQUEUED: Task %s | Queue: %s | Topic: %s\n", info.ID, info.Queue, TopicPaymentFailed)
	return nil
}

func (e *Enqueuer) EnqueueSubscriptionHalted(payload SubscriptionHaltedPayload) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TopicSubscriptionHalted, payloadBytes, asynq.MaxRetry(5))
	
	info, err := e.client.Enqueue(task)
	if err != nil {
		log.Printf("ERROR: Could not enqueue SubscriptionHalted task: %v\n", err)
		return err
	}

	log.Printf("ENQUEUED: Task %s | Queue: %s | Topic: %s\n", info.ID, info.Queue, TopicSubscriptionHalted)
	return nil
}
