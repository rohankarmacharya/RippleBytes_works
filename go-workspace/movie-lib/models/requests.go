// MovieQueryParams represents the query parameters for filtering and pagination
package models

type MovieQueryParams struct {
	Page        int      `query:"page" default:"1"`
	Limit       int      `query:"limit" default:"10"`
	Search      string   `query:"search"`
	MinRating   *float64 `query:"min_rating"`
	MaxRating   *float64 `query:"max_rating"`
	ReleaseFrom string   `query:"release_from"`
	ReleaseTo   string   `query:"release_to"`
}

// PaginatedResponse represents the paginated response structure
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"totalPages"`
}
