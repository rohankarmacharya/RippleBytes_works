package repository

import (
	"errors"
	"time"

	"github.com/rohankarmacharya/movie-lib/config"
	"github.com/rohankarmacharya/movie-lib/models"
	"gorm.io/gorm"
)

// GetAllMovies retrieves all movies from the database
func GetAllMovies() ([]models.Movie, error) {
	var movies []models.Movie
	err := config.DB.
		Order("created_at DESC").
		Find(&movies).Error
	if err != nil {
		return nil, err
	}

	return movies, nil
}

// SaveMovies saves multiple movies to the database
// Returns the number of successfully saved movies and any error that occurred
func SaveMovies(movies []models.Movie) (int, error) {
	savedCount := 0
	for _, m := range movies {
		err := CreateOrUpdateMovie(m)
		if err != nil {
			return savedCount, err
		}
		savedCount++
	}
	return savedCount, nil
}

// GetMovieByID retrieves a single movie by its ID
func GetMovieByID(id string) (*models.Movie, error) {
	var movie models.Movie
	err := config.DB.Model(&models.Movie{}).
		Where("id = ?", id).
		First(&movie).Error

	if err != nil {
		return nil, err
	}

	return &movie, nil
}

// CreateOrUpdateMovie creates a new movie or updates it if it already exists
func CreateOrUpdateMovie(movie models.Movie) error {
	// Check if movie with external_id already exists
	var existingMovie models.Movie
	result := config.DB.Where("external_id = ?", movie.ExternalID).First(&existingMovie)

	if result.Error == nil {
		// Movie exists, update it
		movie.ID = existingMovie.ID
		movie.CreatedAt = existingMovie.CreatedAt
		movie.UpdatedAt = time.Now()
		return config.DB.Save(&movie).Error
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Movie doesn't exist, create new
		movie.CreatedAt = time.Now()
		movie.UpdatedAt = time.Now()
		return config.DB.Create(&movie).Error
	} else {
		// Some other error occurred
		return result.Error
	}
}

// GetTotalMoviesCount returns the total number of movies in the database
func GetTotalMoviesCount() (int64, error) {
	var count int64
	if err := config.DB.Model(&models.Movie{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
