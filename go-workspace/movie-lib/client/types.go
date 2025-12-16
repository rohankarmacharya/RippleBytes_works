package client

// TMDBMovie represents a movie from the TMDB API
type TMDBMovie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	ReleaseDate string  `json:"release_date"`
	VoteAverage float64 `json:"vote_average"`
}

// TMDBResponse represents the response from TMDB API
type TMDBResponse struct {
	Results []TMDBMovie `json:"results"`
}
