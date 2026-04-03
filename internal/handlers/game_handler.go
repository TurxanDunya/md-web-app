package handlers

import (
	"md_api/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateGameRequest struct {
	Title             string    `json:"title"`
	DevelopmentStatus string    `json:"development_status"`
	Description       string    `json:"description"`
	ReleaseDate       time.Time `json:"release_date"`
}

type UpdateGameRequest struct {
	Title             string    `json:"title"`
	DevelopmentStatus string    `json:"development_status"`
	Description       string    `json:"description"`
	ReleaseDate       time.Time `json:"release_date"`
}

func CreateGameHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateGameRequest
		if err := readJSON(r, &request); err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if request.Title == "" || request.DevelopmentStatus == "" {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "title and development_status are required"})
			return
		}

		game, err := repository.CreateGame(pool, request.Title, request.DevelopmentStatus, request.Description, request.ReleaseDate)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		WriteJSON(w, http.StatusCreated, game)
	}
}

func GetAllGamesHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		games, err := repository.GetAllGames(pool)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		WriteJSON(w, http.StatusOK, games)
	}
}

func GetGameByIDHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid game ID"})
			return
		}

		game, err := repository.GetGameByID(pool, int64(idInt))
		if err != nil {
			if err == pgx.ErrNoRows {
				WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Game not found"})
				return
			}
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		WriteJSON(w, http.StatusOK, game)
	}
}

func UpdateGameHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid game ID"})
			return
		}

		var request UpdateGameRequest
		if err := readJSON(r, &request); err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		game, err := repository.UpdateGame(pool, int64(idInt), request.Title, request.DevelopmentStatus, request.Description, request.ReleaseDate)
		if err != nil {
			if err == pgx.ErrNoRows {
				WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Game not found"})
				return
			}
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		WriteJSON(w, http.StatusOK, game)
	}
}

func DeleteGameHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid game ID"})
			return
		}

		err = repository.DeleteGame(pool, int64(idInt))
		if err != nil {
			if err == pgx.ErrNoRows {
				WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Game not found"})
				return
			}
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		WriteJSON(w, http.StatusOK, map[string]string{"message": "Game deleted successfully"})
	}
}
