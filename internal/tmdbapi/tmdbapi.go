package tmdbapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type TopRatedRes struct {
	Page         int             `json:"page"`
	Results      []TopRatedMovie `json:"results"`
	TotalPages   int             `json:"total_pages"`
	TotalResults int             `json:"total_results"`
}

type TopRatedMovie struct {
	Adult            bool   `json:"adult"`
	BackdropPath     string `json:"backdrop_path"`
	GenreIds         []int  `json:"genre_ids"`
	ID               int    `json:"id"`
	OriginalLanguage string `json:"original_language"`
	OriginalTitle    string `json:"original_title"`
	// Overview         string  `json:"overview"`
	// Popularity  float64 `json:"popularity"`
	PosterPath  string `json:"poster_path"`
	ReleaseDate string `json:"release_date"`
	Title       string `json:"title"`
	// Video            bool    `json:"video"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

func GetTopRated() TopRatedRes {
	url := "https://api.themoviedb.org/3/movie/now_playing?language=en-US&page=1"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("TMDB_TOKEN"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("couldn't fetch top rated - ", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	topRated := TopRatedRes{}
	json.Unmarshal(body, &topRated)

	return topRated
}
