// Package tmdbapi defines the service used for interacting with the TMDB API
// for retrieving movie data.
package tmdbapi

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/database"
)

// Service holds the functions needed to consume the TMDB API. Fields should be
// private to prevent access by other services.
type Service struct {
	tmdbTokenEnv string
}

// Review is a type that extends the DB review data by adding the external TMDB
// API data needed before rendering.
type Review struct {
	database.GetReviewsRow
	Title      string
	PosterPath string
}

// NewService is the constructor function for creating the TMDB API service.
func NewService(tmdbToken string) *Service {
	return &Service{
		tmdbTokenEnv: tmdbToken,
	}
}

// GenreMap maps from genre ID to genre name. This was hard-coded based on a
// call to the TMDB movie genre list endpoint
// (https://developer.themoviedb.org/reference/genre-movie-list)
var GenreMap = map[int]string{
	28:    "Action",
	12:    "Adventure",
	16:    "Animation",
	35:    "Comedy",
	80:    "Crime",
	99:    "Documentary",
	18:    "Drama",
	10751: "Family",
	14:    "Fantasy",
	36:    "History",
	27:    "Horror",
	10402: "Music",
	9648:  "Mystery",
	10749: "Romance",
	878:   "Science Fiction",
	10770: "TV Movie",
	53:    "Thriller",
	10752: "War",
	37:    "Western",
}

// IsNotFound checks if the error returned form a TMDB API call is a "not
// found" error
func IsNotFound(err error) bool {
	return err != nil && err.Error() == "not found"
}

// callAPI wraps around the boilerplate needed to make an HTTP call to the TMDB
// api, including the addition of auth headers and error handling
func (s *Service) callAPI(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.tmdbTokenEnv)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("error fetching now playing - %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return []byte{}, fmt.Errorf("not found")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("error reading body - %v", err)
	}

	return body, err
}

// sortByPopularity is a sortFunc used to order movies by descending popularity
func sortByPopularity(a, b NowPlayingMovie) int {
	if a.Popularity < b.Popularity {
		return 1
	}
	if a.Popularity > b.Popularity {
		return -1
	}
	return 0
}

// GetReviewMovieDetails takes a slice of DB reviews and attached the TMDB data
// to them.
func (s *Service) GetReviewMovieDetails(reviews []database.GetReviewsRow) []Review {
	type newDetails struct {
		title      string
		posterPath string
	}

	tmdbData := map[string]newDetails{}
	for _, r := range reviews {
		if _, ok := tmdbData[r.MovieTmdbID]; ok {
			continue
		}

		details, err := s.GetMovieDetails(r.MovieTmdbID)
		if err != nil {
			log.Fatal(err)
		}

		tmdbData[r.MovieTmdbID] = newDetails{
			title:      details.Title,
			posterPath: details.PosterPath,
		}
	}

	newReviews := make([]Review, 0, len(reviews))
	for _, r := range reviews {
		details := tmdbData[r.MovieTmdbID]
		newReview := Review{
			GetReviewsRow: r,
			Title:         details.title,
			PosterPath:    details.posterPath,
		}
		newReviews = append(newReviews, newReview)
	}

	return newReviews
}
