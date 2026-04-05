package repository

import (
	"context"
	"errors"
	"log/slog"
	"md_api/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateGame(pool *pgxpool.Pool, title string, developmentStatus string, description string, releaseDate time.Time, logger *slog.Logger) (*models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO games (title, development_status, description, release_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, development_status, description, release_date, created_at, updated_at`

	var game models.Game
	err := pool.QueryRow(ctx, query, title, developmentStatus, description, releaseDate).Scan(
		&game.ID,
		&game.Title,
		&game.DevelopmentStatus,
		&game.Description,
		&game.ReleaseDate,
		&game.CreatedAt,
		&game.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Error("db: query timeout creating game", "error", err)
		} else {
			logger.Error("db: failed to create game", "error", err)
		}
		return nil, err
	}

	return &game, nil
}

func GetAllGames(pool *pgxpool.Pool, logger *slog.Logger) ([]*models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, title, development_status, description, release_date, created_at, updated_at
	            FROM games
		    ORDER BY created_at DESC`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Error("db: query timeout getting all games", "error", err)
		} else {
			logger.Error("db: failed to query games", "error", err)
		}
		return nil, err
	}
	defer rows.Close()

	var games []*models.Game
	for rows.Next() {
		var game models.Game
		err := rows.Scan(
			&game.ID,
			&game.Title,
			&game.DevelopmentStatus,
			&game.Description,
			&game.ReleaseDate,
			&game.CreatedAt,
			&game.UpdatedAt,
		)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Error("db: query timeout scanning game row", "error", err)
			} else {
				logger.Error("db: failed to scan game row", "error", err)
			}
			return nil, err
		}
		games = append(games, &game)
	}

	return games, nil
}

func GetGameByID(pool *pgxpool.Pool, id int64, logger *slog.Logger) (*models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, title, development_status, description, release_date, created_at, updated_at FROM games
		WHERE id = $1`

	var game models.Game
	err := pool.QueryRow(ctx, query, id).Scan(
		&game.ID,
		&game.Title,
		&game.DevelopmentStatus,
		&game.Description,
		&game.ReleaseDate,
		&game.CreatedAt,
		&game.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Error("db: query timeout getting game", "game_id", id, "error", err)
		} else if err != pgx.ErrNoRows {
			logger.Error("db: failed to get game", "game_id", id, "error", err)
		}
		return nil, err
	}

	return &game, nil
}

func UpdateGame(pool *pgxpool.Pool, id int64, title string, developmentStatus string, description string, releaseDate time.Time, logger *slog.Logger) (*models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE games SET title = $1, development_status = $2, description = $3, release_date = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, title, development_status, description, release_date, created_at, updated_at`

	var game models.Game
	err := pool.QueryRow(ctx, query, title, developmentStatus, description, releaseDate, id).Scan(
		&game.ID,
		&game.Title,
		&game.DevelopmentStatus,
		&game.Description,
		&game.ReleaseDate,
		&game.CreatedAt,
		&game.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Error("db: query timeout updating game", "game_id", id, "error", err)
		} else if err != pgx.ErrNoRows {
			logger.Error("db: failed to update game", "game_id", id, "error", err)
		}
		return nil, err
	}

	return &game, nil
}

func DeleteGame(pool *pgxpool.Pool, id int64, logger *slog.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM games WHERE id = $1`
	commandTag, err := pool.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Error("db: query timeout deleting game", "game_id", id, "error", err)
		} else {
			logger.Error("db: failed to delete game", "game_id", id, "error", err)
		}
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
