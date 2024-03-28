package service

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/templs"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

func (s *Service) HandleReviewsNewGet(w http.ResponseWriter, r *http.Request) {
	_, err := s.authJWTCookie(r)
	if err != nil {
		cookie := createCookie("pmdb-requested-url", r.URL.String(), "/login", 3600)
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	movieId := r.URL.Query().Get("movieId")
	if movieId == "" {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	movieDetails, err := tmdbapi.GetMovieDetails(movieId)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	err = templs.NewReview(templs.NewReviewData{Movie: movieDetails}).Render(r.Context(), w)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleReviewsNewPost(w http.ResponseWriter, r *http.Request) {
	dbUser, err := s.authJWTCookie(r)
	if err != nil {
		log.Fatal(err)
	}

	rating, err := strconv.Atoi(r.FormValue("rating"))
	if err != nil {
		log.Fatal(err)
	}

	var publicReview int64
	if r.FormValue("public-review") == "on" {
		publicReview = 1
	}

	_, err = s.DB.CreateReview(r.Context(), database.CreateReviewParams{
		ID:           uuid.NewString(),
		UserID:       dbUser.ID,
		MovieTmdbID:  r.FormValue("movieID"),
		Rating:       int64(rating),
		Review:       r.FormValue("review"),
		PublicReview: publicReview,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Review saved!")
}
