package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"os"
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
	"github.com/dis70rt/flowback/internal/telemetry"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/genai"
)

const TypeProcessWebhook = "webhook:received"

func main() {
	cfg := config.Load()

	ctx := context.Background()
	shutdown, err := telemetry.Init(ctx, "flowback-worker", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		slog.Error("failed to init telemetry", "error", err)
	} else {
		defer shutdown(ctx)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	queries := repo.New(db)

	registry, err := core.InitModels(ctx, cfg.OpenRouterAPIKey, cfg.OpenRouterModel)
	if err != nil {
		slog.Error("failed to init models", "error", err)
		os.Exit(1)
	}

	strategyAgent, err := nodes.NewStrategyAgent(registry, queries)
	if err != nil {
		slog.Error("failed to init strategy agent", "error", err)
		os.Exit(1)
	}

	copywriterAgent, err := nodes.NewCopywriterAgent(registry)
	if err != nil {
		slog.Error("failed to init copywriter agent", "error", err)
		os.Exit(1)
	}

	voiceAgent, err := nodes.NewVoiceAgent(registry)
	if err != nil {
		slog.Error("failed to init voice agent", "error", err)
		os.Exit(1)
	}

	bus, _ := pubsub.New(cfg.RedisURL)
	orchestrator, err := flowagent.BuildOrchestrator(ctx, queries, strategyAgent, copywriterAgent, voiceAgent, bus, cfg.OpenRouterAPIKey)
	if err != nil {
		slog.Error("failed to build orchestrator", "error", err)
		os.Exit(1)
	}

	r, err := runner.NewInMemory("recovery_worker", orchestrator)
	if err != nil {
		slog.Error("failed to create runner", "error", err)
		os.Exit(1)
	}

	var redisConnOpt asynq.RedisConnOpt
	if strings.HasPrefix(cfg.RedisURL, "redis://") || strings.HasPrefix(cfg.RedisURL, "rediss://") {
		opt, err := asynq.ParseRedisURI(cfg.RedisURL)
		if err != nil {
			slog.Error("failed to parse redis url", "error", err)
			os.Exit(1)
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
			slog.ErrorContext(c, "failed to parse task payload", "error", err)
			return err
		}

		taskID := t.ResultWriter().TaskID()
		slog.InfoContext(c, "processing webhook", "event", wp.Event, "task_id", taskID)

		userID := "system"
		sessionID := taskID

		// Build a tree root if tree logging is enabled.
		var root *telemetry.TreeNode
		nodeMap := map[string]*telemetry.TreeNode{}
		if cfg.TreeLog {
			shortID := sessionID
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}
			root = telemetry.NewRootNode("RecoveryCase: " + shortID)
		}

		inputContent := genai.Text(string(wp.RawJSON))[0]
		iter := r.Run(c, userID, sessionID, inputContent, agent.RunConfig{})

		var graphErr error
		for ev, err := range iter {
			if err != nil {
				slog.ErrorContext(c, "graph execution error", "error", err)
				graphErr = err
				break
			}

			// Accumulate tree node for this event.
			if ev != nil && root != nil {
				// ev.Author is the agent/node name for LLM agents.
				// ev.NodeInfo.Path is set for function nodes.
				nodeName := ev.Author
				if ev.NodeInfo != nil && ev.NodeInfo.Path != "" {
					// Take just the first segment (top-level node name).
					seg := ev.NodeInfo.Path
					if idx := strings.Index(seg, "/"); idx >= 0 {
						seg = seg[:idx]
					}
					if at := strings.Index(seg, "@"); at >= 0 {
						seg = seg[:at]
					}
					if seg != "" {
						nodeName = seg
					}
				}

				if nodeName == "" || nodeName == "user" {
					continue
				}

				treeNode, exists := nodeMap[nodeName]
				if !exists {
					treeNode = root.AddChild(nodeName)
					nodeMap[nodeName] = treeNode
				}

				// Mark the node finished when it emits a final response or output.
				if ev.Output != nil || ev.IsFinalResponse() {
					treeNode.Finish(nil)
				}
			}
		}

		// Finish any unfinished tree nodes and the root, then print.
		if root != nil {
			for _, n := range nodeMap {
				n.FinishIfRunning(graphErr)
			}
			root.Finish(graphErr)
			root.Print(os.Stdout)
		}

		if graphErr != nil {
			return graphErr
		}

		slog.InfoContext(c, "recovery workflow completed", "task_id", taskID)
		return nil
	})

	slog.Info("worker starting", "redis", cfg.RedisURL, "tree_log", cfg.TreeLog)
	if err := srv.Run(mux); err != nil {
		slog.Error("worker server exited", "error", err)
		os.Exit(1)
	}
}
