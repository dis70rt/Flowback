package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/dis70rt/flowback/internal/api"
	"github.com/dis70rt/flowback/internal/config"
	"github.com/dis70rt/flowback/internal/database"
	"github.com/dis70rt/flowback/internal/events"
	"github.com/dis70rt/flowback/internal/pubsub"
	"github.com/dis70rt/flowback/internal/razorpay"
	"github.com/dis70rt/flowback/internal/repo"
	"github.com/dis70rt/flowback/internal/telemetry"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	shutdown, err := telemetry.Init(ctx, "flowback-api", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		slog.Error("failed to init telemetry", "error", err)
	} else {
		defer shutdown(ctx)
	}

	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	queries := repo.New(db)

	asynqClient := database.InitAsynqClient(cfg.RedisAddr)
	defer asynqClient.Close()
	enqueuer := events.NewEnqueuer(asynqClient)

	bus, err := pubsub.New(cfg.RedisURL)
	if err != nil {
		slog.Error("failed to init pubsub", "error", err)
		os.Exit(1)
	}
	defer bus.Close()

	rzpClient := razorpay.NewClient(cfg.RazorpayKeyID, cfg.RazorpayKeySecret)

	router := api.NewRouter(api.RouterDeps{
		Queries:        queries,
		Enqueuer:       enqueuer,
		Bus:            bus,
		RazorpaySecret: cfg.RazorpaySecret,
		RazorpayClient: rzpClient,
	})

	slog.Info("starting flowback api server", "port", 8080)
	if err := router.Run(":8080"); err != nil {
		slog.Error("api server exited", "error", err)
		os.Exit(1)
	}
}
