// Package nowplaying defines the service used for the now playing page,
// including related routes, handlers, and templates.
package nowplaying

import (
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/auth"
	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/errors"
	"github.com/wipdev-tech/pmdb/internal/logger"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

// Service holds the router, handlers, and functions related to the now playing
// page. Fields should be private to prevent access by other services.
type Service struct {
	auth *auth.Service
	tmdb *tmdbapi.Service
}

// NewService is the constructor function for creating the now playing service.
func NewService(auth *auth.Service, tmdb *tmdbapi.Service) *Service {
	return &Service{
		auth: auth,
		tmdb: tmdb,
	}
}

// NewRouter creates a http.Handler with the route for the now playing page.
func (s *Service) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", logger.Middleware(s.auth.Middleware(s.handleNowPlayingGet), "Now Playing (GET) handler"))

	return mux
}

func (s *Service) handleNowPlayingGet(w http.ResponseWriter, r *http.Request, user database.GetUserRow) {
	nowPlaying, err := s.tmdb.GetNowPlaying(-1)
	if err != nil {
		errors.Render(w, http.StatusInternalServerError)
		return
	}

	templData := NowPlayingData{
		User:       user,
		NowPlaying: nowPlaying,
	}

	err = NowPlaying(templData).Render(r.Context(), w)
	if err != nil {
		errors.Render(w, http.StatusInternalServerError)
		return
	}
}
