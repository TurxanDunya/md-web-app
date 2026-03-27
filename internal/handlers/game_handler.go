package handlers

import (
	"md_api/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateGameRequest struct {
	Title             string    `json:"title" binding:"required"`
	DevelopmentStatus string    `json:"development_status" binding:"required"`
	Description       string    `json:"description"`
	ReleaseDate       time.Time `json:"release_date"`
}

type UpdateGameRequest struct {
	Title             string    `json:"title"`
	DevelopmentStatus string    `json:"development_status"`
	Description       string    `json:"description"`
	ReleaseDate       time.Time `json:"release_date"`
}

func CreateGameHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		var request CreateGameRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		game, err := repository.CreateGame(pool, request.Title, request.DevelopmentStatus, request.Description, request.ReleaseDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, game)
	}
}

func GetAllGamesHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		games, err := repository.GetAllGames(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, games)
	}
}

func GetGameByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
			return
		}

		game, err := repository.GetGameByID(pool, int64(idInt))
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, game)
	}
}

func UpdateGameHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
			return
		}

		var request UpdateGameRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		game, err := repository.UpdateGame(pool, int64(idInt), request.Title, request.DevelopmentStatus, request.Description, request.ReleaseDate)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, game)
	}	
}

func DeleteGameHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func (c *gin.Context) {
		id := c.Param("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
			return
		}

		err = repository.DeleteGame(pool, int64(idInt))
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Game deleted successfully"})
	}	
}