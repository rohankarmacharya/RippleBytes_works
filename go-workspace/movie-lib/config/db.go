package config

import (
	"fmt"
	"os"

	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}
}

func ConnectDB() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@localhost:5432/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(fmt.Errorf("unable to parse config: %v", err))
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(fmt.Errorf("unable to connect to database: %v", err))
	}

	// fmt.Println("Database connected successfully!")
}
