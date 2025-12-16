package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rohankarmacharya/movie-lib/client"
	"github.com/rohankarmacharya/movie-lib/models"
	"github.com/rohankarmacharya/movie-lib/repository"
	"github.com/rohankarmacharya/movie-lib/service"
	"gorm.io/gorm"
)

type movieRequest struct {
	ExternalID  string   `json:"external_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	PosterPath  string   `json:"poster_path"`
	ReleaseDate string   `json:"release_date"`
	Rating      *float64 `json:"rating"`
}

// GetMovies handles GET /api/movies
func GetMovies(c *fiber.Ctx) error {
	var queryParams models.MovieQueryParams
	if err := c.QueryParser(&queryParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	// Set default values if not provided
	if queryParams.Page < 1 {
		queryParams.Page = 1
	}
	if queryParams.Limit < 1 || queryParams.Limit > 100 {
		queryParams.Limit = 10
	}

	result, err := repository.GetMoviesWithPagination(queryParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch movies",
		})
	}

	return c.JSON(result)
}

// GetMovie handles GET /api/movies/:id
func GetMovie(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid movie ID",
		})
	}

	movie, err := repository.GetMovieByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Movie not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch movie",
		})
	}

	return c.JSON(movie)
}

// CreateMovie handles POST /api/movies
func CreateMovie(c *fiber.Ctx) error {
	var req movieRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	var releaseDate time.Time
	if req.ReleaseDate != "" {
		parsed, err := time.Parse("2006-01-02", req.ReleaseDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid release_date format, expected YYYY-MM-DD",
			})
		}
		releaseDate = parsed
	}

	movie := models.Movie{
		ExternalID:  req.ExternalID,
		Title:       req.Title,
		Description: req.Description,
		PosterPath:  req.PosterPath,
		ReleaseDate: releaseDate,
	}
	if req.Rating != nil {
		movie.Rating = *req.Rating
	}

	if err := repository.CreateMovie(&movie); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create movie",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(movie)
}

// UpdateMovie handles PUT /api/movies/:id
func UpdateMovie(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid movie ID",
		})
	}

	var req movieRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var releaseDate time.Time
	if req.ReleaseDate != "" {
		parsed, err := time.Parse("2006-01-02", req.ReleaseDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid release_date format, expected YYYY-MM-DD",
			})
		}
		releaseDate = parsed
	}

	movie := models.Movie{
		ID:          uint(id),
		ExternalID:  req.ExternalID,
		Title:       req.Title,
		Description: req.Description,
		PosterPath:  req.PosterPath,
		ReleaseDate: releaseDate,
	}
	if req.Rating != nil {
		movie.Rating = *req.Rating
	}

	if err := repository.UpdateMovie(&movie); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Movie not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update movie",
		})
	}

	return c.JSON(movie)
}

// DeleteMovie handles DELETE /api/movies/:id
func DeleteMovie(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid movie ID",
		})
	}

	if err := repository.DeleteMovie(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Movie deleted",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete movie",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// SyncMovies handles POST /api/movies/sync
func SyncMovies(c *fiber.Ctx) error {
	if err := service.SyncWithAPI(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to sync movies from TMDB",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Movies synced successfully",
	})
}

// SearchTMDBMovies handles GET /api/tmdb/movies/search
func SearchTMDBMovies(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query parameter 'query' is required",
		})
	}

	movies, err := client.SearchMovies(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search TMDB",
		})
	}

	return c.JSON(movies)
}

// GetTMDBMovieDetails handles GET /api/tmdb/movies/:id
func GetTMDBMovieDetails(c *fiber.Ctx) error {
	id := c.Params("id")
	movie, err := client.FetchMovieDetails(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch movie details from TMDB",
		})
	}

	return c.JSON(movie)
}
