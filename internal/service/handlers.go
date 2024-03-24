package service

import (
	"log"
	"net/http"
	"sync"

	"github.com/wipdev-tech/pmdb/internal/templs"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

// HandleHome is the handler for the home route ("/")
func (s *Service) HandleHome(w http.ResponseWriter, r *http.Request) {
	tmplData := templs.IndexData{}

	dbUser, err := s.authJWTCookie(r)
	if err != nil && err != http.ErrNoCookie {
		log.Fatal(err)
	}
	tmplData.LoggedIn = err == nil
	tmplData.User = dbUser

	nowPlaying, err := tmdbapi.GetNowPlaying(5)
	if err != nil {
		log.Fatal(err)
	}
	tmplData.NowPlaying = nowPlaying

	err = templs.Index(tmplData).Render(r.Context(), w)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleProfilesGet(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("userName")
	dbUser, err := s.DB.GetUser(r.Context(), userName)
	if err != nil {
		log.Fatal("couldn't get user - ", err)
	}

	err = templs.Profile(templs.ProfileData{User: dbUser}).Render(r.Context(), w)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleMoviesGet(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	templData := templs.MovieData{}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		movieDetails, err := tmdbapi.GetMovieDetails(movieID)
		if err != nil {
			log.Fatal(err)
		}

		templData.Details = movieDetails
		wg.Done()
	}()

	go func() {
		movieCast, err := tmdbapi.GetMovieCast(movieID)
		if err != nil {
			log.Fatal(err)
		}

		templData.Cast = movieCast
		wg.Done()
	}()

	wg.Wait()
	err := templs.Movie(templData).Render(r.Context(), w)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleNowPlayingGet(w http.ResponseWriter, r *http.Request) {
	nowPlaying, err := tmdbapi.GetNowPlaying(-1)
	if err != nil {
		log.Fatal(err)
	}

	templData := templs.NowPlayingData{
		NowPlaying: nowPlaying,
	}

	err = templs.NowPlaying(templData).Render(r.Context(), w)
	if err != nil {
		log.Fatal(err)
	}
}
