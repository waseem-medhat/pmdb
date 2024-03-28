package router

import (
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/service"
)

func New(s *service.Service) *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("GET /{$}", s.HandleHome)

	mux.HandleFunc("GET /register", s.HandleRegisterGet)
	mux.HandleFunc("POST /register", s.HandleRegisterPost)
	mux.HandleFunc("POST /register/validate", s.HandleRegisterValidate)

	mux.HandleFunc("GET /login", s.HandleLoginGet)
	mux.HandleFunc("POST /login", s.HandleLoginPost)
	mux.HandleFunc("GET /logout", s.HandleLogoutPost)

	mux.HandleFunc("GET /profiles/{userName}", s.HandleProfilesGet)
	mux.HandleFunc("GET /movies/{movieID}", s.HandleMoviesGet)
	mux.HandleFunc("GET /now-playing", s.HandleNowPlayingGet)

	mux.HandleFunc("GET /reviews/new", s.HandleReviewsNewGet)
	mux.HandleFunc("POST /reviews/new", s.HandleReviewsNewPost)

	return mux
}
