package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/wipdev-tech/pmdb/internal/database"
)

func (s *Service) generateJWTAccess(userName string) (string, error) {
	access := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "pmdb-acess",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Subject:   userName,
		},
	)

	accessStr, err := access.SignedString([]byte(s.jwtSecretEnv))
	if err != nil {
		return accessStr, fmt.Errorf("couldn't sign access jwt - %v", err)
	}

	return accessStr, err
}

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

// AuthJWTCookie checks the request's cookies for an access token and tries to
// authenticate the user who owns the JWT.
func (s *Service) AuthJWTCookie(r *http.Request) (database.GetUserRow, error) {
	dbUser := database.GetUserRow{}
	claims := &jwt.RegisteredClaims{}
	keyfunc := func(toke *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecretEnv), nil
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

	dbUser, err = s.db.GetUser(r.Context(), userName)
	if err == sql.ErrNoRows {
		return dbUser, err
	}

	if err != nil {
		return dbUser, fmt.Errorf("couldn't query user - %v", err)
	}

	return dbUser, err
}
