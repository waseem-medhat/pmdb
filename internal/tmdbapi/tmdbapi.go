package tmdbapi

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/wipdev-tech/pmdb/internal/database"
)

// Review is a type  that extends the DB review data by adding the external
// TMDB API data needed before rendering.
type Review struct {
	database.GetReviewsRow
	Title      string
	PosterPath string
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

// callAPI wraps around the boilerplate needed to make an HTTP call to the TMDB
// api, including the addition of auth headers and error handling
func callAPI(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("TMDB_TOKEN"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("error fetching now playing - %v", err)
	}
	defer res.Body.Close()

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
