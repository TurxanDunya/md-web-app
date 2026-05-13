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

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, email, password, created_at, updated_at`

	var createdUser models.User
	err := r.pool.QueryRow(ctx, query, user.Email, user.Password).Scan(
		&createdUser.ID,
		&createdUser.Email,
		&createdUser.Password,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("db: query timeout creating user", "email", user.Email, "error", err)
		} else {
			slog.Error("db: failed to create user", "email", user.Email, "error", err)
		}
		return nil, err
	}

	return &createdUser, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, email, is_active, created_at, updated_at
		FROM users
		WHERE email = $1`

	var user models.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("db: query timeout getting user by email", "error", err)
		} else if err != pgx.ErrNoRows {
			slog.Error("db: failed to get user by email", "error", err)
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, email, is_active, created_at, updated_at
		FROM users
		WHERE id = $1`

	var user models.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("db: query timeout getting user by id", "user_id", id, "error", err)
		} else if err != pgx.ErrNoRows {
			slog.Error("db: failed to get user by id", "user_id", id, "error", err)
		}
		return nil, err
	}

	return &user, nil
}
