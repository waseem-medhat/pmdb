package tmdbapi

import (
	"encoding/json"
	"fmt"
	"log"
)

// MovieDetails matches the JSON response coming from a call to the Movie
// Details API endpoint
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

type MovieCredits struct {
	ID   int         `json:"id"`
	Cast []MovieCast `json:"cast"`
	// Crew []struct {
	// 	Adult              bool    `json:"adult"`
	// 	Gender             int     `json:"gender"`
	// 	ID                 int     `json:"id"`
	// 	KnownForDepartment string  `json:"known_for_department"`
	// 	Name               string  `json:"name"`
	// 	OriginalName       string  `json:"original_name"`
	// 	Popularity         float64 `json:"popularity"`
	// 	ProfilePath        string  `json:"profile_path"`
	// 	CreditID           string  `json:"credit_id"`
	// 	Department         string  `json:"department"`
	// 	Job                string  `json:"job"`
	// } `json:"crew"`
}

type MovieCast struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	OriginalName string  `json:"original_name"`
	Popularity   float64 `json:"popularity"`
	ProfilePath  string  `json:"profile_path"`
	CastID       int     `json:"cast_id"`
	Character    string  `json:"character"`
	CreditID     string  `json:"credit_id"`
	Order        int     `json:"order"`
	// Adult              bool    `json:"adult"`
	// Gender             int     `json:"gender"`
	// KnownForDepartment string  `json:"known_for_department"`
}

// GetMovieDetails makes the call to the Movie Details API
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

// GetMovieCast makes the call to the Movie Credits API
func GetMovieCast(movieID string) []MovieCast {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s/credits", movieID)
	responseBody, err := callAPI(url)
	if err != nil {
		log.Fatal("error calling API - ", err)
	}

	movieCredits := MovieCredits{}
	err = json.Unmarshal(responseBody, &movieCredits)
	if err != nil {
		log.Fatal("couldn't unmarshal movie details - ", err)
	}
	return movieCredits.Cast[:10]
}
