package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/dis70rt/flowback/migrations"
)

func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Tell goose to use our embedded SQL files from the migrations package
	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Run migrations (using "." because the FS is anchored in the migrations folder)
	if err := goose.Up(db, "."); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	slog.Info("DATABASE: Connected and Migrations Applied!")
	return db, nil
}
