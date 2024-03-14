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
	claims := &jwt.RegisteredClaims{}
	keyfunc := func(toke *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}

	accessCookie, err := r.Cookie("jwt-access")

	dbUser := database.GetUserRow{}
	if err == nil {
		bearer := accessCookie.Value

		token, err := jwt.ParseWithClaims(bearer, claims, keyfunc)
		if err != nil {
			log.Fatal(err)
		}

		userName, err := token.Claims.GetSubject()
		if err != nil {
			log.Fatal(err)
		}

		dbUser, err = s.DB.GetUser(r.Context(), userName)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(dbUser.DisplayName)

	tmpl := template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/blocks/_top.html",
		"templates/blocks/_bottom.html",
	))
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
