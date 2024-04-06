// Package movies defines the service used for the movies and reviews pages and
// logic, including related routes, handlers, and templates.
package movies

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wipdev-tech/pmdb/internal/auth"
	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/errors"
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

func (s *Service) handleMoviesGet(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	templData := MoviePageData{}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	var fetchErr error
	go func() {
		movieDetails, err := s.tmdb.GetMovieDetails(movieID)
		if err != nil {
			fetchErr = err
		}
		templData.Details = movieDetails
		wg.Done()
	}()

	go func() {
		movieCast, err := s.tmdb.GetMovieCast(movieID)
		if err != nil {
			fetchErr = err
		}
		templData.Cast = movieCast
		wg.Done()
	}()

	wg.Wait()
	if tmdbapi.IsNotFound(fetchErr) {
		errors.Render(w, http.StatusNotFound)
		return
	}
	if fetchErr != nil {
		fmt.Println(fetchErr)
		errors.Render(w, http.StatusInternalServerError)
		return
	}

	err := MoviePage(templData).Render(r.Context(), w)
	if err != nil {
		errors.Render(w, http.StatusInternalServerError)
		return
	}
}

func (s *Service) handleReviewsNewGet(w http.ResponseWriter, r *http.Request, _ database.GetUserRow) {
	movieID := r.PathValue("movieID")
	if movieID == "" {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	movieDetails, err := s.tmdb.GetMovieDetails(movieID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	err = NewReviewPage(NewReviewPageData{Movie: movieDetails}).Render(r.Context(), w)
	if err != nil {
		errors.Render(w, http.StatusInternalServerError)
		return
	}

}

func (s *Service) handleReviewsNewPost(w http.ResponseWriter, r *http.Request, dbUser database.GetUserRow) {
	var dbTimeLayout = "2006-01-02 15:04"

	movieID := r.PathValue("movieID")
	if movieID == "" {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	rating, err := strconv.Atoi(r.FormValue("rating"))
	if err != nil {
		errors.Render(w, http.StatusBadRequest)
		return
	}

	var publicReview int64
	if r.FormValue("public-review") == "on" {
		publicReview = 1
	}

	_, err = s.db.CreateReview(r.Context(), database.CreateReviewParams{
		ID:           uuid.NewString(),
		CreatedAt:    time.Now().Format(dbTimeLayout),
		UpdatedAt:    time.Now().Format(dbTimeLayout),
		UserID:       dbUser.ID,
		MovieTmdbID:  movieID,
		Rating:       int64(rating),
		Review:       r.FormValue("review"),
		PublicReview: publicReview,
	})
	if err != nil {
		errors.Render(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusFound)
}
