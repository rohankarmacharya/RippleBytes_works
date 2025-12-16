package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rohankarmacharya/movie-api/routes"
	"github.com/rohankarmacharya/movie-lib/config"
)

func main() {
	config.ConnectDB() // connect to Postgres
	fmt.Println("Database connected successfully!")

	app := fiber.New()

	routes.MovieRoutes(app)

	app.Listen(":3000")
}
