package models

import "time"

type Movie struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ExternalID  string    `json:"external_id" gorm:"uniqueIndex:idx_external_id"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description,omitempty"`
	ReleaseDate time.Time `json:"release_date" time_format:"2006-01-02"`
	Rating      float64   `json:"rating,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
