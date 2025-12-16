package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rohankarmacharya/movie-lib/models"
)

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
