package movies

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wipdev-tech/pmdb/internal/database"
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
