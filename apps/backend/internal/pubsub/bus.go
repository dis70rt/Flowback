package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Publisher interface {
	Publish(ctx context.Context, channel string, message any) error
}

type MessageHandler func(ctx context.Context, payload []byte)

type Subscriber interface {
	Subscribe(ctx context.Context, channel string, handler MessageHandler)
}

// Bus combines Publisher and Subscriber.
type Bus interface {
	Publisher
	Subscriber
	Close() error
}


type redisBus struct {
	client *redis.Client
}

// New creates a new Redis-backed pubsub Bus.
func New(redisURL string) (Bus, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("pubsub: invalid redis URL: %w", err)
	}

	rdb := redis.NewClient(opt)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("pubsub: redis ping failed: %w", err)
	}

	return &redisBus{client: rdb}, nil
}

func (b *redisBus) Publish(ctx context.Context, channel string, message any) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("pubsub: failed to marshal message for channel %q: %w", channel, err)
	}

	if err := b.client.Publish(ctx, channel, payload).Err(); err != nil {
		return fmt.Errorf("pubsub: failed to publish to channel %q: %w", channel, err)
	}

	return nil
}

func (b *redisBus) Subscribe(ctx context.Context, channel string, handler MessageHandler) {
	go func() {
		sub := b.client.Subscribe(ctx, channel)
		defer func() {
			if err := sub.Close(); err != nil {
				log.Printf("[PUBSUB] Error closing subscription for channel %q: %v\n", channel, err)
			}
		}()

		log.Printf("[PUBSUB] Subscribed -> channel=%q\n", channel)

		ch := sub.Channel()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					log.Printf("[PUBSUB] Channel closed -> channel=%q\n", channel)
					return
				}
				handler(ctx, []byte(msg.Payload))

			case <-ctx.Done():
				log.Printf("[PUBSUB] Context cancelled -> unsubscribing from channel=%q\n", channel)
				return
			}
		}
	}()
}

func (b *redisBus) Close() error {
	return b.client.Close()
}
