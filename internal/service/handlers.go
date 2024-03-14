package service

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/wipdev-tech/pmdb/internal/database"
)

// HandleHome is the handler for the home route ("/")
func (s *Service) HandleHome(w http.ResponseWriter, r *http.Request) {
	tmplData := struct {
		LoggedIn bool
		User     database.GetUserRow
	}{}

	dbUser, err := s.authJWTCookie(r)
	tmplData.User = dbUser
	tmplData.LoggedIn = err == nil
	if err != nil {
		fmt.Println(err)
	}

	err = template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/blocks/_top.html",
		"templates/blocks/_bottom.html",
	)).Execute(w, tmplData)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleProfileGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.PathValue("userName"))
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
