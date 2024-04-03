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

	mux.HandleFunc("GET /profiles/{userName}", s.HandleProfilesGet)
	mux.HandleFunc("GET /movies/{movieID}", s.HandleMoviesGet)

	mux.HandleFunc("GET /reviews/new", s.MiddlewareAuth(s.HandleReviewsNewGet))
	mux.HandleFunc("POST /reviews/new", s.MiddlewareAuth(s.HandleReviewsNewPost))

	return mux
}
