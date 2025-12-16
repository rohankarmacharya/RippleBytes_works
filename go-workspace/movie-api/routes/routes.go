package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rohankarmacharya/movie-lib/client"
	"github.com/rohankarmacharya/movie-lib/models"
	"github.com/rohankarmacharya/movie-lib/repository"
)

func MovieRoutes(app *fiber.App) {
	app.Get("/movies", func(c *fiber.Ctx) error {
		movies, err := repository.GetAllMovies()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(movies)
	})

	app.Get("movies/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		movie, err := repository.GetMovieByID(id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(movie)
	})

	app.Post("/movies/sync", func(c *fiber.Ctx) error {
		movies, err := client.FetchMovies()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		syncedCount, err := repository.SaveMovies(movies)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		totalCount, err := repository.GetTotalMoviesCount()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"message":      "Movies synced successfully",
			"synced_count": syncedCount,
			"total_movies": totalCount,
		})
	})

	app.Post("/movies", func(c *fiber.Ctx) error {
		movie := new(models.Movie)
		if err := c.BodyParser(movie); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		err := repository.CreateOrUpdateMovie(*movie)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(201).JSON(movie)

	})
	app.Get("/movies/sync", func(c *fiber.Ctx) error {
		syncedCount, err := client.SyncWithAPI()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// Get total count from DB
		totalCount, err := repository.GetTotalMoviesCount()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"synced_count": syncedCount, // movies synced from TMDB
			"total_count":  totalCount,  // total movies in DB now
			"message":      "Movies synced successfully!",
		})
	})

}
