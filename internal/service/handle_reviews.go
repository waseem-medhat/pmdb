package service

import (
	"log"
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/templs"
)

func (s *Service) HandleReviewsNewGet(w http.ResponseWriter, r *http.Request) {
	_, err := s.authJWTCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	err = templs.NewReview().Render(r.Context(), w)
	if err != nil {
		log.Fatal(err)
	}
}
