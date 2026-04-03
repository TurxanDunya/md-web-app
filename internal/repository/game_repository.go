package repository

import (
	"context"
	"md_api/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateGame(pool *pgxpool.Pool, title string, developmentStatus string, description string, releaseDate time.Time) (*models.Game, error) {
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
		return nil, err
	}

	return &game, nil
}

func GetAllGames(pool *pgxpool.Pool) ([]*models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, title, development_status, description, release_date, created_at, updated_at FROM games
		ORDER BY created_at DESC`

	rows, err := pool.Query(ctx, query)
	if err != nil {
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
			return nil, err
		}
		games = append(games, &game)
	}

	return games, nil
}

func GetGameByID(pool *pgxpool.Pool, id int64) (*models.Game, error) {
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
		return nil, err
	}

	return &game, nil
}

func UpdateGame(pool *pgxpool.Pool, id int64, title string, developmentStatus string, description string, releaseDate time.Time) (*models.Game, error) {
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
		return nil, err
	}

	return &game, nil
}

func DeleteGame(pool *pgxpool.Pool, id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM games WHERE id = $1`
	commandTag, err := pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
