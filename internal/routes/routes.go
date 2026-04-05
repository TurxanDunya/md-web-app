package routes

import (
	"log/slog"
	"md_api/internal/config"
	"md_api/internal/handlers"
	"md_api/internal/middleware"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(mux *http.ServeMux, pool *pgxpool.Pool, cfg *config.Config, logger *slog.Logger) {
	// Ping
	mux.HandleFunc("GET /ping", middleware.RequestLogger(logger, func(w http.ResponseWriter, r *http.Request) {
		handlers.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "pong",
			"status":  "success",
		})
	}))

	// Games
	mux.HandleFunc("POST /api/v1/games",
		middleware.RequestLogger(logger,
			middleware.AuthMiddleware(cfg, logger, handlers.CreateGameHandler(pool, logger))))
	mux.HandleFunc("GET /api/v1/games",
		middleware.RequestLogger(logger, handlers.GetAllGamesHandler(pool, logger)))
	mux.HandleFunc("GET /api/v1/games/{id}",
		middleware.RequestLogger(logger, handlers.GetGameByIDHandler(pool, logger)))
	mux.HandleFunc("PUT /api/v1/games/{id}",
		middleware.RequestLogger(logger,
			middleware.AuthMiddleware(cfg, logger, handlers.UpdateGameHandler(pool, logger))))
	mux.HandleFunc("DELETE /api/v1/games/{id}",
		middleware.RequestLogger(logger,
			middleware.AuthMiddleware(cfg, logger, handlers.DeleteGameHandler(pool, logger))))

	// Auth
	mux.HandleFunc("POST /api/v1/auth/register",
		middleware.RequestLogger(logger, handlers.RegisterUserHandler(pool, logger)))
	mux.HandleFunc("POST /api/v1/auth/login",
		middleware.RequestLogger(logger, handlers.LoginUserHandler(pool, cfg, logger)))
}
