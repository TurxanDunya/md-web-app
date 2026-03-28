package models

import "time"

type User struct {
	ID                string       `json:"id" db:"id"`
	Email             string    `json:"email" db:"email"`
	Password          string    `json:"password" db:"password"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}