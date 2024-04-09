// Package movies defines the service used for movies pages logic, including
// related routes, handlers, and templates.
package movies

import (
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/auth"
	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/logger"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

// Service holds the router, handlers, and functions related to the movies
// pages. Fields should be private to prevent access by other services.
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

// NewRouter creates a http.Handler with the routes for the movies pages.
func (s *Service) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{movieID}", logger.Middleware(s.auth.Middleware(s.handleMoviesGet), "Movies (GET) handler"))
	return mux
}
