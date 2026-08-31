package database

import (
	"log"

	"github.com/hibiken/asynq"
)

// InitAsynqClient connects to Redis and returns an Asynq client for enqueuing tasks.
func InitAsynqClient(redisAddr string) *asynq.Client {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	log.Printf("REDIS: Asynq client initialized at %s\n", redisAddr)
	
	return client
}
