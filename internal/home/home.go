// Package home defines the service used for the home page, including related
// routes, handlers, and templates.
package home

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/wipdev-tech/pmdb/internal/auth"
	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/errors"
	"github.com/wipdev-tech/pmdb/internal/logger"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

// Service holds the router, handlers, and functions related to the home page.
// Fields should be private to prevent access by other services.
type Service struct {
	auth *auth.Service
	tmdb *tmdbapi.Service
	db   *database.Queries
}

// NewService is the constructor function for creating the home page service.
func NewService(auth *auth.Service, tmdb *tmdbapi.Service, db *database.Queries) *Service {
	return &Service{
		auth: auth,
		tmdb: tmdb,
		db:   db,
	}
}

// NewRouter creates a http.Handler with the route for the home page.
func (s *Service) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", logger.Middleware(s.auth.Middleware(s.handleHomeGet), "Home (GET) handler"))

	return mux
}

// HandleHome is the handler for the home route ("/")
func (s *Service) handleHomeGet(w http.ResponseWriter, r *http.Request, user database.GetUserRow) {
	loggedIn := user.UserName != "guest"

	nowPlaying, err := s.tmdb.GetNowPlaying(5)
	if err != nil {
		fmt.Println(err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}

	reviews, err := s.db.GetReviews(r.Context())
	if err != nil {
		fmt.Println(err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}

	reviewsCleaned := make([]database.GetReviewsRow, 0, len(reviews))
	for _, r := range reviews {
		r.Review = strings.ReplaceAll(r.Review, "\\n", "\n")
		reviewsCleaned = append(reviewsCleaned, r)
	}

	templData := IndexPageData{
		LoggedIn:   loggedIn,
		User:       user,
		NowPlaying: nowPlaying,
		Reviews:    s.tmdb.GetReviewMovieDetails(reviewsCleaned),
	}

	err = IndexPage(templData).Render(r.Context(), w)
	if err != nil {
		fmt.Println(err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}
}
