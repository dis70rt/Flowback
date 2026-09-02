package database

import (
	"log/slog"

	"github.com/hibiken/asynq"
)

// InitAsynqClient connects to Redis and returns an Asynq client for enqueuing tasks.
func InitAsynqClient(redisAddr string) *asynq.Client {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	slog.Info("asynq client initialized", "redis_addr", redisAddr)

	return client
}
