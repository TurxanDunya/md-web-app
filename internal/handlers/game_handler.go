package handlers

import (
	"log/slog"
	"md_api/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
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

type GameHandler struct {
	repo *repository.GameRepository
}

func NewGameHandler(repo *repository.GameRepository) *GameHandler {
	return &GameHandler{repo: repo}
}

func (h *GameHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateGameRequest
		if err := readJSON(r, &request); err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Internal server error"})
			return
		}

		if request.Title == "" || request.DevelopmentStatus == "" {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "title and development_status are required"})
			return
		}

		game, err := h.repo.Create(request.Title, request.DevelopmentStatus, request.Description, request.ReleaseDate)
		if err != nil {
			slog.Error("failed to create game", "error", err)
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			return
		}

		WriteJSON(w, http.StatusCreated, game)
	}
}

func (h *GameHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		games, err := h.repo.GetAll()
		if err != nil {
			slog.Error("failed to get all games", "error", err)
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			return
		}

		WriteJSON(w, http.StatusOK, games)
	}
}

func (h *GameHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid game ID"})
			return
		}

		game, err := h.repo.GetByID(int64(idInt))
		if err != nil {
			if err == pgx.ErrNoRows {
				WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Game not found"})
				return
			}
			slog.Error("failed to get game", "game_id", idInt, "error", err)
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			return
		}

		WriteJSON(w, http.StatusOK, game)
	}
}

func (h *GameHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid game ID"})
			return
		}

		var request UpdateGameRequest
		if err := readJSON(r, &request); err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Internal server error"})
			return
		}

		game, err := h.repo.Update(int64(idInt), request.Title, request.DevelopmentStatus, request.Description, request.ReleaseDate)
		if err != nil {
			if err == pgx.ErrNoRows {
				WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Game not found"})
				return
			}
			slog.Error("failed to update game", "game_id", idInt, "error", err)
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			return
		}

		WriteJSON(w, http.StatusOK, game)
	}
}

func (h *GameHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid game ID"})
			return
		}

		err = h.repo.Delete(int64(idInt))
		if err != nil {
			if err == pgx.ErrNoRows {
				WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Game not found"})
				return
			}
			slog.Error("failed to delete game", "game_id", idInt, "error", err)
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			return
		}

		WriteJSON(w, http.StatusOK, map[string]string{"message": "Game deleted successfully"})
	}
}
