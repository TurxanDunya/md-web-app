package routes

import (
	"md_api/internal/config"
	"md_api/internal/handlers"
	"md_api/internal/middleware"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(mux *http.ServeMux, pool *pgxpool.Pool, cfg *config.Config) {
	// Ping
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		handlers.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "pong",
			"status":  "success",
		})
	})

	// Games
	mux.HandleFunc("POST /api/v1/games", middleware.AuthMiddleware(cfg, handlers.CreateGameHandler(pool)))
	mux.HandleFunc("GET /api/v1/games", handlers.GetAllGamesHandler(pool))
	mux.HandleFunc("GET /api/v1/games/{id}", handlers.GetGameByIDHandler(pool))
	mux.HandleFunc("PUT /api/v1/games/{id}", middleware.AuthMiddleware(cfg, handlers.UpdateGameHandler(pool)))
	mux.HandleFunc("DELETE /api/v1/games/{id}", middleware.AuthMiddleware(cfg, handlers.DeleteGameHandler(pool)))

	// Auth
	mux.HandleFunc("POST /api/v1/auth/register", handlers.RegisterUserHandler(pool))
	mux.HandleFunc("POST /api/v1/auth/login", handlers.LoginUserHandler(pool, cfg))
}
