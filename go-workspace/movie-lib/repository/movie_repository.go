package repository

import (
	"context"

	"github.com/rohankarmacharya/movie-lib/config"
	"github.com/rohankarmacharya/movie-lib/models"
)

// GetAllMovies retrieves all movies from the database
func GetAllMovies() ([]models.Movie, error) {
	db := config.DB
	rows, err := db.Query(context.Background(),
		`SELECT id, external_id, title, description, release_date, rating, created_at, updated_at 
        FROM movies ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		err := rows.Scan(
			&m.ID,
			&m.ExternalID,
			&m.Title,
			&m.Description,
			&m.ReleaseDate,
			&m.Rating,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	if err = rows.Err(); err != nil {
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
	db := config.DB
	var movie models.Movie
	err := db.QueryRow(context.Background(),
		`SELECT id, external_id, title, description, release_date, rating, created_at, updated_at 
        FROM movies WHERE id = $1`, id).Scan(
		&movie.ID,
		&movie.ExternalID,
		&movie.Title,
		&movie.Description,
		&movie.ReleaseDate,
		&movie.Rating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &movie, nil
}

// CreateOrUpdateMovie creates a new movie or updates it if it already exists
func CreateOrUpdateMovie(movie models.Movie) error {
	db := config.DB
	_, err := db.Exec(context.Background(),
		`INSERT INTO movies (external_id, title, description, release_date, rating, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (external_id) 
		DO UPDATE SET 
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			release_date = EXCLUDED.release_date,
			rating = EXCLUDED.rating,
			updated_at = EXCLUDED.updated_at`,
		movie.ExternalID, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating, movie.CreatedAt, movie.UpdatedAt,
	)
	return err
}

// GetTotalMoviesCount returns the total number of movies in the database
func GetTotalMoviesCount() (int, error) {
	db := config.DB
	var count int
	err := db.QueryRow(context.Background(), "SELECT COUNT(*) FROM movies").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
