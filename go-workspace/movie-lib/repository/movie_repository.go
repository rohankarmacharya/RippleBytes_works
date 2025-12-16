package repository

import (
	"time"

	"github.com/rohankarmacharya/movie-lib/config"
	"github.com/rohankarmacharya/movie-lib/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateMovie creates a new movie
func CreateMovie(movie *models.Movie) error {
	// Upsert by external_id: insert or update core fields on conflict
	return config.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "external_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"title", "description", "poster_path", "release_date", "rating", "updated_at"}),
		}).
		Create(movie).Error
}

// SaveMovies saves multiple movies to the database and returns the count of saved movies
func SaveMovies(movies []models.Movie) (int64, error) {
	// Bulk upsert by external_id
	result := config.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "external_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"title", "description", "poster_path", "release_date", "rating", "updated_at"}),
		}).
		Create(&movies)
	return result.RowsAffected, result.Error
}

// GetMovieByID gets a movie by ID
func GetMovieByID(id uint) (*models.Movie, error) {
	var movie models.Movie
	err := config.DB.First(&movie, id).Error
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

// UpdateMovie updates an existing movie
func UpdateMovie(movie *models.Movie) error {
	result := config.DB.Save(movie)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// DeleteMovie deletes a movie by ID
func DeleteMovie(id uint) error {
	result := config.DB.Delete(&models.Movie{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// GetAllMovies retrieves all movies from the database
func GetAllMovies() ([]models.Movie, error) {
	var movies []models.Movie
	err := config.DB.Find(&movies).Error
	if err != nil {
		return nil, err
	}
	return movies, nil
}

// GetMovieByTitleAndDate finds a movie by its title and release date
func GetMovieByTitleAndDate(title string, releaseDate time.Time) (*models.Movie, error) {
	var movie models.Movie
	err := config.DB.Where("title = ? AND release_date = ?", title, releaseDate).First(&movie).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &movie, nil
}

// GetMoviesWithPagination retrieves movies with filtering, searching, and pagination
func GetMoviesWithPagination(params models.MovieQueryParams) (*models.PaginatedResponse, error) {
	var movies []models.Movie
	var total int64

	// Start building the query
	query := config.DB.Model(&models.Movie{})

	// Apply search (case-insensitive search in title and description)
	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query = query.Where("LOWER(title) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)",
			searchTerm, searchTerm)
	}

	// Apply rating filters
	if params.MinRating != nil {
		query = query.Where("rating >= ?", *params.MinRating)
	}
	if params.MaxRating != nil {
		query = query.Where("rating <= ?", *params.MaxRating)
	}

	// Apply release date filters (parse strings as YYYY-MM-DD, ignore invalid formats)
	if params.ReleaseFrom != "" {
		if from, err := time.Parse("2006-01-02", params.ReleaseFrom); err == nil {
			query = query.Where("release_date >= ?", from)
		}
	}
	if params.ReleaseTo != "" {
		if to, err := time.Parse("2006-01-02", params.ReleaseTo); err == nil {
			endOfDay := to.Add(24 * time.Hour)
			query = query.Where("release_date < ?", endOfDay)
		}
	}

	// Get total count for pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate pagination
	offset := (params.Page - 1) * params.Limit
	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	// Apply pagination and ordering
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(params.Limit).
		Find(&movies).Error; err != nil {
		return nil, err
	}

	// Build the response
	response := &models.PaginatedResponse{
		Data:       movies,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}

	return response, nil
}
