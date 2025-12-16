package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/rohankarmacharya/movie-lib/models"
)

// FetchMovieDetails gets detailed information about a specific movie
func FetchMovieDetails(movieID string) (*models.Movie, error) {
	apiKey := os.Getenv("TMDB_API_KEY")
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?api_key=%s", movieID, apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var movie models.Movie
	if err := json.Unmarshal(body, &movie); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &movie, nil
}

// FetchMovies gets a list of popular movies
func FetchMovies() ([]models.Movie, error) {
	apiKey := os.Getenv("TMDB_API_KEY")
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tmdbResp TMDBResponse
	if err := json.Unmarshal(body, &tmdbResp); err != nil {
		return nil, err
	}

	var movies []models.Movie
	for _, m := range tmdbResp.Results {
		releaseDate, _ := time.Parse("2006-01-02", m.ReleaseDate)
		movies = append(movies, models.Movie{
			ExternalID:  fmt.Sprint(m.ID),
			Title:       m.Title,
			Description: m.Overview,
			ReleaseDate: releaseDate,
			Rating:      float64(m.VoteAverage),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	return movies, nil
}

// SearchMovies searches for movies by query using the TMDB API
func SearchMovies(query string) ([]models.Movie, error) {
	apiKey := os.Getenv("TMDB_API_KEY")
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s", apiKey, url.QueryEscape(query))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making search request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var tmdbResp TMDBResponse
	if err := json.Unmarshal(body, &tmdbResp); err != nil {
		return nil, fmt.Errorf("error decoding search response: %v", err)
	}

	var movies []models.Movie
	for _, m := range tmdbResp.Results {
		releaseDate, _ := time.Parse("2006-01-02", m.ReleaseDate)
		movies = append(movies, models.Movie{
			ExternalID:  fmt.Sprint(m.ID),
			Title:       m.Title,
			Description: m.Overview,
			ReleaseDate: releaseDate,
			Rating:      float64(m.VoteAverage),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	return movies, nil
}
