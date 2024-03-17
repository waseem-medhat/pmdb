package service

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) HandleLoginGet(w http.ResponseWriter, r *http.Request) {
	err := template.Must(template.ParseFiles(
		"templates/login.html",
		"templates/blocks/_top.html",
		"templates/blocks/_bottom.html",
	)).Execute(w, struct{ LoginError bool }{LoginError: false})

	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleLoginPost(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("user-name")
	password := r.FormValue("password")

	dbUser, err := s.DB.GetUserForLogin(r.Context(), userName)
	if err == sql.ErrNoRows {
		err := template.Must(template.ParseFiles(
			"templates/login.html",
		)).ExecuteTemplate(w, "login-error", struct{ LoginError bool }{LoginError: true})

		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if err != nil {
		err := template.Must(template.ParseFiles(
			"templates/login.html",
		)).ExecuteTemplate(w, "login-error", struct{ LoginError bool }{LoginError: true})

		if err != nil {
			log.Fatal(err)
		}
		return
	}

	access, err := generateJWTAccess(userName)
	if err != nil {
		log.Fatal(err)
	}

	cookie := &http.Cookie{
		Name:     "jwt-access",
		Value:    access,
		Secure:   true,
		HttpOnly: true,
		MaxAge:   3600,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusFound)
}

func generateJWTAccess(userName string) (string, error) {
	access := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "pmdb-acess",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Subject:   userName,
		},
	)

	accessStr, err := access.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return accessStr, fmt.Errorf("couldn't sign access jwt - %v", err)
	}

	return accessStr, err
}

func (s *Service) HandleLogoutPost(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "jwt-access",
		Value:    "",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusFound)
}
