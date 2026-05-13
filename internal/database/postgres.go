package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseUrl string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		slog.Error("unable to parse DATABASE_URL", "error", err)
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error("unable to create connection pool", "error", err)
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		slog.Error("unable to ping database", "error", err)
		pool.Close()
		return nil, err
	}

	slog.Info("connected to database")
	return pool, nil
}
