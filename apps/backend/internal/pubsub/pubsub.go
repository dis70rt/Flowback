package pubsub

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

type Engine struct {
	client *redis.Client
}

func New(redisURL string) (*Engine, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opt)
	
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Engine{client: rdb}, nil
}

func (e *Engine) Publish(ctx context.Context, channel string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return e.client.Publish(ctx, channel, payload).Err()
}

func (e *Engine) Subscribe(ctx context.Context, channel string, handler func(payload []byte)) {
	go func() {
		sub := e.client.Subscribe(ctx, channel)
		defer sub.Close()

		log.Printf("[PUBSUB] Successfully subscribed to channel: %s", channel)

		ch := sub.Channel()
		for {
			select {
			case msg := <-ch:
				handler([]byte(msg.Payload))
			
			case <-ctx.Done():
				log.Printf("[PUBSUB] Context cancelled. Unsubscribing from channel: %s", channel)
				return
			}
		}
	}()
}
