package tmdbapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
)

type NowPlayingRes struct {
	Page         int               `json:"page"`
	Results      []NowPlayingMovie `json:"results"`
	TotalPages   int               `json:"total_pages"`
	TotalResults int               `json:"total_results"`
}

type NowPlayingMovie struct {
	ID          int     `json:"id"`
	Popularity  float64 `json:"popularity"`
	PosterPath  string  `json:"poster_path"`
	ReleaseDate string  `json:"release_date"`
	Title       string  `json:"title"`
	// Adult        bool   `json:"adult"`
	// BackdropPath string `json:"backdrop_path"`
	// GenreIds     []int  `json:"genre_ids"`
	// OriginalLanguage string `json:"original_language"`
	// OriginalTitle    string `json:"original_title"`
	// Overview         string  `json:"overview"`
	// Video            bool    `json:"video"`
	// VoteAverage float64 `json:"vote_average"`
	// VoteCount   int     `json:"vote_count"`
}

type MovieDetails struct {
	// Adult               bool   `json:"adult"`
	// BackdropPath        string `json:"backdrop_path"`
	// BelongsToCollection any    `json:"belongs_to_collection"`
	// Budget              int    `json:"budget"`
	ID          int     `json:"id"`
	ImdbID      string  `json:"imdb_id"`
	Overview    string  `json:"overview"`
	Popularity  float64 `json:"popularity"`
	PosterPath  string  `json:"poster_path"`
	ReleaseDate string  `json:"release_date"`
	Revenue     int     `json:"revenue"`
	Runtime     int     `json:"runtime"`
	Tagline     string  `json:"tagline"`
	Title       string  `json:"title"`
	Video       bool    `json:"video"`
	Genres      []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	// Homepage            string `json:"homepage"`
	// OriginalLanguage    string `json:"original_language"`
	// OriginalTitle       string `json:"original_title"`
	// ProductionCompanies []struct {
	// 	ID            int    `json:"id"`
	// 	LogoPath      string `json:"logo_path"`
	// 	Name          string `json:"name"`
	// 	OriginCountry string `json:"origin_country"`
	// } `json:"production_companies"`
	// ProductionCountries []struct {
	// 	Iso31661 string `json:"iso_3166_1"`
	// 	Name     string `json:"name"`
	// } `json:"production_countries"`
	// SpokenLanguages []struct {
	// 	EnglishName string `json:"english_name"`
	// 	Iso6391     string `json:"iso_639_1"`
	// 	Name        string `json:"name"`
	// } `json:"spoken_languages"`
	// Status      string  `json:"status"`
	// VoteAverage float64 `json:"vote_average"`
	// VoteCount   int     `json:"vote_count"`
}

func GetNowPlaying() []NowPlayingMovie {
	url := "https://api.themoviedb.org/3/movie/now_playing?language=en-US&page=1"
	responseBody, err := callAPI(url)
	if err != nil {
		log.Fatal("error calling API - ", err)
	}

	nowPlaying := NowPlayingRes{}
	err = json.Unmarshal(responseBody, &nowPlaying)
	if err != nil {
		log.Fatal("couldn't unmarshal now playing - ", err)
	}

	results := nowPlaying.Results
	slices.SortFunc(results, sortByPopularity)
	return results
}

func GetMovieDetails(movieID string) MovieDetails {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?language=en-US", movieID)
	responseBody, err := callAPI(url)
	if err != nil {
		log.Fatal("error calling API - ", err)
	}

	movieDetails := MovieDetails{}
	err = json.Unmarshal(responseBody, &movieDetails)
	if err != nil {
		log.Fatal("couldn't unmarshal movie details - ", err)
	}
	return movieDetails
}

func callAPI(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("TMDB_TOKEN"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't fetch now playing - %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("couldn't read body - %v", err)
	}

	return body, err
}

func sortByPopularity(a, b NowPlayingMovie) int {
	if a.Popularity < b.Popularity {
		return 1
	}
	if a.Popularity > b.Popularity {
		return -1
	}
	return 0
}
