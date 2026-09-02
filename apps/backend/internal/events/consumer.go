package events

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
)

type TaskHandler func(ctx context.Context, t *asynq.Task) error

type Consumer struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewConsumer(redisAddr string, concurrency int) *Consumer {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				"default": 10,
			},
		},
	)

	return &Consumer{
		server: srv,
		mux:    asynq.NewServeMux(),
	}
}

func (c *Consumer) Register(topic string, handler TaskHandler) {
	c.mux.HandleFunc(topic, handler)
	slog.Info(fmt.Sprintf("[CONSUMER] Registered handler for topic: %s\n", topic))
}

func (c *Consumer) Start() error {
	slog.Info("[CONSUMER] Starting background task consumption...")
	return c.server.Run(c.mux)
}
