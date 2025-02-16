package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"ttavito/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(cfg *config.Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		slog.Error("Failed to parse database configuration", "error", err)
		return nil, fmt.Errorf("failed to parse database configuration: %w", err)
	}

	config.MaxConns = 100
	config.MinConns = 10
	config.MaxConnLifetime = time.Hour
	config.HealthCheckPeriod = 20 * time.Second

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		slog.Error("Failed to connect to the database", "error", err)
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Ping(ctx); err != nil {
		db.Close()
		slog.Error("Database is unavailable", "error", err)
		return nil, fmt.Errorf("database is unavailable: %w", err)
	}

	slog.Info("Successfully connected to PostgreSQL (pgx)")
	return db, nil
}
