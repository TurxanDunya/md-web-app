package main

import (
	"md_api/internal/config"
	"md_api/internal/database"
	"md_api/internal/handlers"
	"md_api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	if err != nil {
		panic(err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		panic(err)
	}

	defer pool.Close()

	var router *gin.Engine =  gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status": "success",
		})
	})

	router.POST("/api/v1/games", middleware.AuthMiddleware(cfg), handlers.CreateGameHandler(pool))
	router.GET("/api/v1/games", handlers.GetAllGamesHandler(pool))
	router.GET("/api/v1/games/:id", handlers.GetGameByIDHandler(pool))
	router.PUT("/api/v1/games/:id", middleware.AuthMiddleware(cfg), handlers.UpdateGameHandler(pool))
	router.DELETE("/api/v1/games/:id", middleware.AuthMiddleware(cfg), handlers.DeleteGameHandler(pool))

	router.POST("/api/v1/auth/register", handlers.RegisterUserHandler(pool))
	router.POST("/api/v1/auth/login", handlers.LoginUserHandler(pool, cfg))

	router.Run(":" + cfg.Port)
}