package movies

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/wipdev-tech/pmdb/internal/errors"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

func (s *Service) handleMoviesGet(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	templData := MoviePageData{}

	wg := &sync.WaitGroup{}
	wg.Add(3)

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

	go func() {
		reviews, err := s.db.GetReviewsForMovie(r.Context(), movieID)
		if err != nil {
			fetchErr = err
		}
		templData.Reviews = reviews
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
