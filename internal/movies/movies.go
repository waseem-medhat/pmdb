// Package movies defines the service used for the movies and reviews pages and
// logic, including related routes, handlers, and templates.
package movies

import (
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/auth"
	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/logger"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

// Service holds the router, handlers, and functions related to the movies and
// reviews pages. Fields should be private to prevent access by other services.
type Service struct {
	auth *auth.Service
	tmdb *tmdbapi.Service
	db   *database.Queries
}

// NewService is the constructor function for creating the movies service.
func NewService(auth *auth.Service, tmdb *tmdbapi.Service, db *database.Queries) *Service {
	return &Service{
		auth: auth,
		tmdb: tmdb,
		db:   db,
	}
}

// NewRouter creates a http.Handler with the routes for the movies and reviews pages.
func (s *Service) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{movieID}", logger.Middleware(s.handleMoviesGet, "Movies (GET) handler"))
	mux.HandleFunc("GET /{movieID}/reviews/new", s.auth.MiddlewareAuth(s.handleReviewsNewGet))
	mux.HandleFunc("POST /{movieID}/reviews/new", s.auth.MiddlewareAuth(s.handleReviewsNewPost))

	return mux
}
