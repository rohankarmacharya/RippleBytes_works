package models

import "time"

type Movie struct {
	ID          int       `json:"id" db:"id"`
	ExternalID  string    `json:"external_id" db:"external_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	ReleaseDate time.Time `json:"release_date" db:"release_date"`
	Rating      float64   `json:"rating" db:"rating"`
	CreatedAt   time.Time `json:"CreatedAt" db:"created_at"`
	UpdatedAt   time.Time `json:"UpdatedAt" db:"updated_at"`
}
