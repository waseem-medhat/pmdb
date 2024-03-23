package tmdbapi

import (
	"encoding/json"
	"fmt"
	"slices"
)

// NowPlayingRes matches the JSON response coming from a call to the Now
// Playing API
type NowPlayingRes struct {
	Page         int               `json:"page"`
	Results      []NowPlayingMovie `json:"results"`
	TotalPages   int               `json:"total_pages"`
	TotalResults int               `json:"total_results"`
}

// NowPlayingMovie matches the JSON structure of a single movie in a
// NowPlayingRes struct
type NowPlayingMovie struct {
	ID          int     `json:"id"`
	Popularity  float64 `json:"popularity"`
	PosterPath  string  `json:"poster_path"`
	ReleaseDate string  `json:"release_date"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	GenreIds    []int   `json:"genre_ids"`
	// Adult        bool   `json:"adult"`
	// BackdropPath string `json:"backdrop_path"`
	// OriginalLanguage string `json:"original_language"`
	// OriginalTitle    string `json:"original_title"`
	// Video            bool    `json:"video"`
	// VoteAverage float64 `json:"vote_average"`
	// VoteCount   int     `json:"vote_count"`
}

// GetNowPlaying makes the call to the Now Playing API endpoint and sorts them
// by descending popularity
func GetNowPlaying(n int) ([]NowPlayingMovie, error) {
	url := "https://api.themoviedb.org/3/movie/now_playing?language=en-US&page=1"
	responseBody, err := callAPI(url)
	if err != nil {
		return []NowPlayingMovie{}, fmt.Errorf("error calling API - %v", err)
	}

	nowPlaying := NowPlayingRes{}
	err = json.Unmarshal(responseBody, &nowPlaying)
	if err != nil {
		return nowPlaying.Results, fmt.Errorf("error unmarshaling now playing - %v", err)
	}

	results := nowPlaying.Results
	slices.SortFunc(results, sortByPopularity)

	if n > 0 {
		results = results[:n]
	}

	return results, err
}
