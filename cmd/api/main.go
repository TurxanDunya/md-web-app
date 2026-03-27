package main

import (
	"md_api/internal/config"
	"md_api/internal/database"
	"md_api/internal/handlers"

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

	router.POST("/games", handlers.CreateGameHandler(pool))
	router.GET("/games", handlers.GetAllGamesHandler(pool))
	router.GET("/games/:id", handlers.GetGameByIDHandler(pool))
	router.PUT("/games/:id", handlers.UpdateGameHandler(pool))
	router.DELETE("/games/:id", handlers.DeleteGameHandler(pool))

	router.Run(":" + cfg.Port)
}