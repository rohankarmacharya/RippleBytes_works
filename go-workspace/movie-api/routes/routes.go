package routes

import (
	"movie-api/handlers"

	"github.com/gofiber/fiber/v2"
)

// MovieRoutes configures all the movie-related routes
func MovieRoutes(app *fiber.App) {
	// API v1 routes
	api := app.Group("/api")
	{
		// Movie routes
		movies := api.Group("/movies")
		{
			// Get all movies with filtering and pagination
			movies.Get("/", handlers.GetMovies)

			// Get single movie by ID
			movies.Get("/:id", handlers.GetMovie)

			// Create new movie
			movies.Post("/", handlers.CreateMovie)

			// Sync movies from TMDB and store in DB
			movies.Post("/sync", handlers.SyncMovies)

			// Update existing movie
			movies.Put("/:id", handlers.UpdateMovie)

			// Delete movie
			movies.Delete("/:id", handlers.DeleteMovie)

			// TMDB integration routes
			tmdb := api.Group("/tmdb")
			{
				tmdb.Get("/movies/search", handlers.SearchTMDBMovies)
				tmdb.Get("/movies/:id", handlers.GetTMDBMovieDetails)
			}
		}
	}
}
