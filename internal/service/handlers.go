package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/wipdev-tech/pmdb/internal/database"
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

func (s *Service) authJWTCookie(r *http.Request) (database.GetUserRow, error) {
	dbUser := database.GetUserRow{}
	claims := &jwt.RegisteredClaims{}
	keyfunc := func(toke *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}

	accessCookie, err := r.Cookie("jwt-access")
	if err == http.ErrNoCookie {
		return dbUser, err
	} else if err != nil {
		return dbUser, fmt.Errorf("couldn't get cookie - %v", err)
	}

	bearer := accessCookie.Value
	token, err := jwt.ParseWithClaims(bearer, claims, keyfunc)
	if err != nil {
		return dbUser, fmt.Errorf("couldn't parse jwt - %v", err)
	}

	userName, err := token.Claims.GetSubject()
	if err != nil {
		return dbUser, fmt.Errorf("couldn't get jwt subject - %v", err)
	}

	dbUser, err = s.DB.GetUser(r.Context(), userName)
	if err != nil {
		return dbUser, fmt.Errorf("couldn't query user - %v", err)
	}

	return dbUser, err
}
