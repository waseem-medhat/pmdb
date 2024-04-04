package home

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/auth"
	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/errors"
	"github.com/wipdev-tech/pmdb/internal/logger"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

type Service struct {
	auth *auth.Service
	tmdb *tmdbapi.Service
	db   *database.Queries
}

func NewService(auth *auth.Service, tmdb *tmdbapi.Service, db *database.Queries) *Service {
	return &Service{
		auth: auth,
		tmdb: tmdb,
		db:   db,
	}
}

func (s *Service) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", logger.Middleware(s.handleHomeGet, "Home (GET) handler"))

	return mux
}

// HandleHome is the handler for the home route ("/")
func (s *Service) handleHomeGet(w http.ResponseWriter, r *http.Request) {
	dbUser, err := s.auth.AuthJWTCookie(r)
	if err != nil && err != http.ErrNoCookie && err != sql.ErrNoRows {
		errors.Render(w, http.StatusInternalServerError)
		return
	}
	loggedIn := err == nil

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

	templData := IndexPageData{
		LoggedIn:   loggedIn,
		User:       dbUser,
		NowPlaying: nowPlaying,
		Reviews:    s.tmdb.GetReviewMovieDetails(reviews),
	}

	err = IndexPage(templData).Render(r.Context(), w)
	if err != nil {
		fmt.Println(err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}
}
