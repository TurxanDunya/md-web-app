package models

import "time"

type Game struct {
	ID                int       `json:"id" db:"id"`
	Title             string    `json:"title" db:"title"`
	DevelopmentStatus string    `json:"development_status" db:"development_status"`
	Description       string    `json:"description" db:"description"`
	ReleaseDate       time.Time `json:"release_date" db:"release_date"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}
