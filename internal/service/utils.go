package service

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/wipdev-tech/pmdb/internal/database"
)

// createCookie is a wrapper that makes it easier and more concise to create a
// *http.Cookie. Some cookie attributes are pre-set to make it secure,
// HTTP-only, with a "Strict" same-site mode.
func createCookie(name, val, path string, maxAgeSec int) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    val,
		Secure:   true,
		HttpOnly: true,
		MaxAge:   maxAgeSec,
		SameSite: http.SameSiteStrictMode,
		Path:     path,
	}
}

// authJWTCookie checks the request's cookies for an access token and tries to
// authenticate the user who owns the JWT.
func (s *Service) authJWTCookie(r *http.Request) (database.GetUserRow, error) {
	dbUser := database.GetUserRow{}
	claims := &jwt.RegisteredClaims{}
	keyfunc := func(toke *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}

	accessCookie, err := r.Cookie("pmdb-jwt-access")
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
	if err == sql.ErrNoRows {
		return dbUser, err
	}

	if err != nil {
		return dbUser, fmt.Errorf("couldn't query user - %v", err)
	}

	return dbUser, err
}
