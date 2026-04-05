package main

import (
	"log/slog"
	"md_api/internal/config"
	"md_api/internal/database"
	"md_api/internal/routes"
	"net/http"
	"os"
)

func main() {
	var handler slog.Handler
	if os.Getenv("APP_ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}
	logger := slog.New(handler)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	pool, err := database.Connect(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	mux := http.NewServeMux()
	routes.Setup(mux, pool, cfg, logger)

	logger.Info("server starting", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
