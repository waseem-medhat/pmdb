package service

import (
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/database"
)

func (s *Service) MiddlewareAuth(h func(http.ResponseWriter, *http.Request, database.GetUserRow)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbUser, err := s.authJWTCookie(r)
		if err != nil {
			cookie := createCookie("pmdb-requested-url", r.URL.String(), "/login", 3600)
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		h(w, r, dbUser)
	}
}
