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
	Auth *auth.Service
	DB   *database.Queries
}

func NewService(auth *auth.Service, db *database.Queries) *Service {
	return &Service{
		Auth: auth,
		DB:   db,
	}
}

func (s *Service) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", logger.Middleware(s.handleHomeGet, "Now Playing (GET) handler"))

	return mux
}

// HandleHome is the handler for the home route ("/")
func (s *Service) handleHomeGet(w http.ResponseWriter, r *http.Request) {
	dbUser, err := s.Auth.AuthJWTCookie(r)
	if err != nil && err != http.ErrNoCookie && err != sql.ErrNoRows {
		errors.Render(w, http.StatusInternalServerError)
		return
	}
	loggedIn := err == nil

	nowPlaying, err := tmdbapi.GetNowPlaying(5)
	if err != nil {
		fmt.Println(err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}

	reviews, err := s.DB.GetReviews(r.Context())
	if err != nil {
		fmt.Println(err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}

	templData := IndexPageData{
		LoggedIn:   loggedIn,
		User:       dbUser,
		NowPlaying: nowPlaying,
		Reviews:    tmdbapi.GetReviewMovieDetails(reviews),
	}

	err = IndexPage(templData).Render(r.Context(), w)
	if err != nil {
		fmt.Println(err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}
}
