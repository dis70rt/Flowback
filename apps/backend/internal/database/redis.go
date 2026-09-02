package database

import (
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
)

// InitAsynqClient connects to Redis and returns an Asynq client for enqueuing tasks.
func InitAsynqClient(redisAddr string) *asynq.Client {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	slog.Info(fmt.Sprintf("REDIS: Asynq client initialized at %s\n", redisAddr))
	
	return client
}
