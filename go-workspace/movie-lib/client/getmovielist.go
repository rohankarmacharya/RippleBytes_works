package client

import (
	"fmt"

	// "github.com/rohankarmacharya/movie-lib/models"
	"github.com/rohankarmacharya/movie-lib/repository"
)

// SyncWithAPI fetches movies from external API and stores in DB
// Returns the number of movies synced and any error that occurred
func SyncWithAPI() (int, error) {
	// Fetch movies from external API
	movies, err := FetchMovies()
	if err != nil {
		return 0, fmt.Errorf("error fetching movies: %v", err)
	}

	// Save movies to repository
	syncedCount, err := repository.SaveMovies(movies)
	if err != nil {
		return 0, fmt.Errorf("error saving movies: %v", err)
	}

	fmt.Printf("Successfully synced %d movies!\n", syncedCount)
	return int(syncedCount), nil
}

// CreateMovie adds a new movie manually
// func CreateMovie(m models.Movie) {
// 	err := repository.CreateOrUpdateMovie(m)
// 	if err != nil {
// 		fmt.Println("Error creating movie:", err)
// 	} else {
// 		fmt.Println("Movie created successfully:", m.Title)
// 	}
// }
