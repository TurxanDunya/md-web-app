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
	configureLogger()
	cfg := configureConfig()

	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	mux := http.NewServeMux()
	routes.Setup(mux, pool, cfg)

	slog.Info("server starting", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func configureLogger() {
	var handler slog.Handler
	if os.Getenv("APP_ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}
	slog.SetDefault(slog.New(handler))
}

func configureConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	return cfg
}
