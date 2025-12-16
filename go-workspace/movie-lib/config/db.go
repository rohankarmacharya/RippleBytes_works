package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rohankarmacharya/movie-lib/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}
}

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database.", err)
	}

	fmt.Println("Database connected successfully!")

	// Drop the table if it exists to avoid constraint issues
	err = DB.Migrator().DropTable(&models.Movie{})
	if err != nil {
		log.Println("Warning: Could not drop movies table:", err)
	}

	// Auto-migrate the model with the new schema
	err = DB.AutoMigrate(&models.Movie{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Ensure the unique index on external_id exists
	if !DB.Migrator().HasIndex(&models.Movie{}, "external_id") {
		err = DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_movies_external_id ON movies(external_id)`).Error
		if err != nil {
			log.Fatal("Failed to create external_id index:", err)
		}
	}

	fmt.Println("Database migrated successfully!")
}
