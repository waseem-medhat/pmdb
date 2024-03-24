package service

import (
	"log"
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/templs"
)

func (s *Service) HandleReviewsNewGet(w http.ResponseWriter, r *http.Request) {
	_, err := s.authJWTCookie(r)
	if err != nil {
		cookie := createCookie("pmdb-requested-url", r.URL.Path, "/login", 3600)
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	err = templs.NewReview().Render(r.Context(), w)
	if err != nil {
		log.Fatal(err)
	}
}
