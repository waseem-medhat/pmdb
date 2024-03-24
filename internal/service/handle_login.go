package service

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/wipdev-tech/pmdb/internal/templs"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) HandleLoginGet(w http.ResponseWriter, r *http.Request) {
	err := templs.Login(templs.LoginData{LoginError: false}).Render(r.Context(), w)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleLoginPost(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("user-name")
	password := r.FormValue("password")

	dbUser, err := s.DB.GetUserForLogin(r.Context(), userName)
	if err == sql.ErrNoRows {
		err := templs.ErrorAlert(templs.LoginData{LoginError: true}).Render(r.Context(), w)
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
		err := templs.ErrorAlert(templs.LoginData{LoginError: true}).Render(r.Context(), w)
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
		Name:     "pmdb-jwt-access",
		Value:    access,
		Secure:   true,
		HttpOnly: true,
		MaxAge:   3600,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	redirectCookie, err := r.Cookie("pmdb-requested-url")
	var redirectTo string
	if err != nil || redirectCookie.Value == "" {
		redirectTo = "/"
	} else {
		redirectTo = redirectCookie.Value
	}

	newRedirectCookie := &http.Cookie{
		Name:     "pmdb-requested-url",
		Value:    "",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Path:     "/login",
	}
	http.SetCookie(w, newRedirectCookie)
	w.Header().Set("HX-Redirect", redirectTo)
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
		Name:     "pmdb-jwt-access",
		Value:    "",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusFound)
}
