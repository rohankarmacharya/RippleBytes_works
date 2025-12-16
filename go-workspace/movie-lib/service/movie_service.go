package service

import (
	"fmt"
	"time"

	"github.com/rohankarmacharya/movie-lib/client"
	"github.com/rohankarmacharya/movie-lib/models"
	"github.com/rohankarmacharya/movie-lib/repository"
)

func SyncWithAPI() error {
	movies, err := client.FetchMovies()
	if err != nil {
		return fmt.Errorf("failed to fetch movies: %w", err)
	}

	for _, m := range movies {
		releaseDate := m.ReleaseDate

		movie := models.Movie{
			ExternalID:  m.ExternalID,
			Title:       m.Title,
			Description: m.Description,
			ReleaseDate: releaseDate,
			Rating:      m.Rating,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   time.Now(),
		}

		existingMovie, err := repository.GetMovieByTitleAndDate(m.Title, releaseDate)

		if err == nil && existingMovie != nil {
			// Movie exists, update it
			movie.ID = existingMovie.ID
			err = repository.UpdateMovie(&movie)
		} else {
			// Movie doesn't exist, create it
			err = repository.CreateMovie(&movie)
		}

		if err != nil {
			return fmt.Errorf("failed to sync movie %s: %v", m.Title, err)
		}
	}

	return nil
}

// GetAllMovies fetches all movies from the repository
func GetAllMovies() ([]models.Movie, error) {
	return repository.GetAllMovies()
}
