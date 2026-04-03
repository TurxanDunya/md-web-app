package models

import "time"

type Game struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	DevelopmentStatus string    `json:"development_status"`
	Description       string    `json:"description"`
	ReleaseDate       time.Time `json:"release_date"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
