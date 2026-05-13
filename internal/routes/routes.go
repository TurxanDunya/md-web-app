package routes

import (
	"md_api/internal/config"
	"md_api/internal/handlers"
	"md_api/internal/middleware"
	"md_api/internal/repository"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(mux *http.ServeMux, pool *pgxpool.Pool, cfg *config.Config) {
	gameRepo := repository.NewGameRepository(pool)
	userRepo := repository.NewUserRepository(pool)

	gameHandler := handlers.NewGameHandler(gameRepo)
	userHandler := handlers.NewUserHandler(userRepo, cfg)

	// Ping
	mux.HandleFunc("GET /ping", middleware.RequestLogger(func(w http.ResponseWriter, r *http.Request) {
		handlers.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "pong",
			"status":  "success",
		})
	}))

	// Games
	mux.HandleFunc("POST /api/v1/games",
		middleware.RequestLogger(
			middleware.AuthMiddleware(cfg, gameHandler.Create())))
	mux.HandleFunc("GET /api/v1/games",
		middleware.RequestLogger(gameHandler.GetAll()))
	mux.HandleFunc("GET /api/v1/games/{id}",
		middleware.RequestLogger(gameHandler.GetByID()))
	mux.HandleFunc("PUT /api/v1/games/{id}",
		middleware.RequestLogger(
			middleware.AuthMiddleware(cfg, gameHandler.Update())))
	mux.HandleFunc("DELETE /api/v1/games/{id}",
		middleware.RequestLogger(
			middleware.AuthMiddleware(cfg, gameHandler.Delete())))

	// Auth
	mux.HandleFunc("POST /api/v1/auth/register",
		middleware.RequestLogger(userHandler.Register()))
	mux.HandleFunc("POST /api/v1/auth/login",
		middleware.RequestLogger(userHandler.Login()))
}
