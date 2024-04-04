package service

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/wipdev-tech/pmdb/internal/templs"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

func (s *Service) HandleMoviesGet(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	templData := templs.MovieData{}

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
		renderError(w, http.StatusNotFound)
		return
	}
	if fetchErr != nil {
		fmt.Println(fetchErr)
		renderError(w, http.StatusInternalServerError)
		return
	}

	err := templs.Movie(templData).Render(r.Context(), w)
	if err != nil {
		renderError(w, http.StatusInternalServerError)
		return
	}
}
