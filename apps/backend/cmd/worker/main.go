package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strings"

	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	
	flowagent "github.com/dis70rt/flowback/internal/agent"
	"github.com/dis70rt/flowback/internal/agent/core"
	"github.com/dis70rt/flowback/internal/agent/nodes"
	"github.com/dis70rt/flowback/internal/config"
	"github.com/dis70rt/flowback/internal/events"
	"github.com/dis70rt/flowback/internal/pubsub"
	"github.com/dis70rt/flowback/internal/repo"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/genai"
)

const (
	TypeProcessWebhook = "webhook:received"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()
	queries := repo.New(db)

	ctx := context.Background()
	registry, err := core.InitModels(ctx, cfg.OpenRouterAPIKey, cfg.OpenRouterModel)
	if err != nil {
		log.Fatalf("failed to init models: %v", err)
	}

	strategyAgent, err := nodes.NewStrategyAgent(registry, queries)
	if err != nil {
		log.Fatalf("failed to init strategy agent: %v", err)
	}
	
	copywriterAgent, err := nodes.NewCopywriterAgent(registry)
	if err != nil {
		log.Fatalf("failed to init copywriter agent: %v", err)
	}

	bus, _ := pubsub.New(cfg.RedisURL)
	orchestrator, err := flowagent.BuildOrchestrator(ctx, queries, strategyAgent, copywriterAgent, bus)
	if err != nil {
		log.Fatalf("failed to build orchestrator: %v", err)
	}

	r, err := runner.NewInMemory("recovery_worker", orchestrator)
	if err != nil {
		log.Fatalf("failed to create runner: %v", err)
	}

	var redisConnOpt asynq.RedisConnOpt
	if strings.HasPrefix(cfg.RedisURL, "redis://") || strings.HasPrefix(cfg.RedisURL, "rediss://") {
		opt, err := asynq.ParseRedisURI(cfg.RedisURL)
		if err != nil {
			log.Fatalf("failed to parse redis url: %v", err)
		}
		redisConnOpt = opt
	} else {
		redisConnOpt = asynq.RedisClientOpt{Addr: cfg.RedisAddr}
	}

	srv := asynq.NewServer(
		redisConnOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	
	mux.HandleFunc(TypeProcessWebhook, func(c context.Context, t *asynq.Task) error {
		var wp events.WebhookPayload
		if err := json.Unmarshal(t.Payload(), &wp); err != nil {
			log.Printf("[WORKER] Failed to parse task payload: %v", err)
			return err
		}
		
		log.Printf("[WORKER] Started processing webhook task. Event: %s", wp.Event)
		
		userID := "system"
		sessionID := t.ResultWriter().TaskID() 
		
		// Extract the true Razorpay JSON to feed to the ADK graph
		inputContent := genai.Text(string(wp.RawJSON))[0]

		iter := r.Run(c, userID, sessionID, inputContent, agent.RunConfig{})
		
		for ev, err := range iter {
			if err != nil {
				log.Printf("[WORKER] Graph Error: %v", err)
				return err 
			}
			if ev != nil {
				log.Printf("[WORKER] Finished processing node step.")
			}
		}

		log.Println("[WORKER] AI Graph execution finished successfully.")
		return nil
	})

	log.Printf("Worker booting up. Listening on Redis: %s", cfg.RedisURL)
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run worker server: %v", err)
	}
}
