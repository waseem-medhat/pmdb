package movies

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/wipdev-tech/pmdb/internal/errors"
	"github.com/wipdev-tech/pmdb/internal/logger"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

type Service struct {
	// Auth *auth.Service
	// DB   *database.Queries
}

// func NewService(auth *auth.Service, db *database.Queries) *Service {
func NewService() *Service {
	return &Service{
		// Auth: auth,
		// DB:   db,
	}
}

func (s *Service) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{movieID}", logger.Middleware(s.handleMoviesGet, "Movies (GET) handler"))

	return mux
}

func (s *Service) handleMoviesGet(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	templData := MoviePageData{}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	var fetchErr error
	go func() {
		movieDetails, err := tmdbapi.GetMovieDetails(movieID)
		if err != nil {
			fetchErr = err
		}
		templData.Details = movieDetails
		wg.Done()
	}()

	go func() {
		movieCast, err := tmdbapi.GetMovieCast(movieID)
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
