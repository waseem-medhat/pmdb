package service

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	"github.com/wipdev-tech/pmdb/internal/templs"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

// HandleHome is the handler for the home route ("/")
func (s *Service) HandleHome(w http.ResponseWriter, r *http.Request) {
	dbUser, err := s.authJWTCookie(r)
	if err != nil && err != http.ErrNoCookie && err != sql.ErrNoRows {
		renderError(w, http.StatusInternalServerError)
		return
	}
	loggedIn := err == nil

	nowPlaying, err := tmdbapi.GetNowPlaying(5)
	if err != nil {
		fmt.Println(err)
		renderError(w, http.StatusInternalServerError)
		return
	}

	reviews, err := s.DB.GetReviews(r.Context())
	if err != nil {
		fmt.Println(err)
		renderError(w, http.StatusInternalServerError)
		return
	}

	templData := templs.IndexData{
		LoggedIn:   loggedIn,
		User:       dbUser,
		NowPlaying: nowPlaying,
		Reviews:    getReviewData(reviews),
	}

	err = templs.Index(templData).Render(r.Context(), w)
	if err != nil {
		fmt.Println(err)
		renderError(w, http.StatusInternalServerError)
		return
	}
}

func (s *Service) HandleProfilesGet(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("userName")
	dbUser, err := s.DB.GetUser(r.Context(), userName)
	if err == sql.ErrNoRows {
		renderError(w, http.StatusNotFound)
		return
	}

	if err != nil {
		fmt.Println("couldn't get user - ", err)
		renderError(w, http.StatusInternalServerError)
		return
	}

	err = templs.Profile(templs.ProfileData{User: dbUser}).Render(r.Context(), w)
	if err != nil {
		fmt.Println(err)
		renderError(w, http.StatusInternalServerError)
		return
	}
}

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
