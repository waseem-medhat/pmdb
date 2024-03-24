package service

import (
	"log"
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/templs"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

func (s *Service) HandleReviewsNewGet(w http.ResponseWriter, r *http.Request) {
	_, err := s.authJWTCookie(r)
	if err != nil {
		cookie := createCookie("pmdb-requested-url", r.URL.Path, "/login", 3600)
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	movieId := r.URL.Query().Get("movieId")
	if movieId == "" {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}

	movieDetails, err := tmdbapi.GetMovieDetails(movieId)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}

	err = templs.NewReview(templs.NewReviewData{Movie: movieDetails}).Render(r.Context(), w)
	if err != nil {
		log.Fatal(err)
	}
}
