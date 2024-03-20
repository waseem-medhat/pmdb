package service

import (
	"fmt"
	"log"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/templs"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

// HandleHome is the handler for the home route ("/")
func (s *Service) HandleHome(w http.ResponseWriter, r *http.Request) {
	dbUser, err := s.authJWTCookie(r)
	if err != nil {
		fmt.Println(err)
	}

	tmplData := templs.IndexData{
		LoggedIn:   err == nil,
		User:       dbUser,
		NowPlaying: tmdbapi.GetNowPlaying(),
	}

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

	templs.Profile(templs.ProfileData{User: dbUser}).Render(r.Context(), w)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleMoviesGet(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	movieDetails := tmdbapi.GetMovieDetails(movieID)

	err := templs.Movie(movieDetails).Render(r.Context(), w)
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
	if err != nil {
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
