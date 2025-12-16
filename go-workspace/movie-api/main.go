package main

import (
	"encoding/json"
	"log"

	"movie-api/routes"

	"github.com/rohankarmacharya/movie-lib/config"
	"github.com/rohankarmacharya/movie-lib/models"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize database
	config.ConnectDB() // Make sure this sets up the DB connection
	db, err := config.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer db.Close()

	// Auto migrate models
	if err := config.DB.AutoMigrate(&models.Movie{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Setup routes
	routes.MovieRoutes(app)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
