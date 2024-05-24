// Package auth defines the service used for user authentication/authorization
// and management, including related routes, handlers, and templates.
package auth

import (
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/logger"
)

// Service holds the router, handlers, and functions related to auth and user
// management. Fields should be private to prevent access by other services.
type Service struct {
	db           *database.Queries
	jwtSecretEnv string
}

// NewService is the constructor function for creating the auth service.
func NewService(db *database.Queries, jwtSecret string) *Service {
	return &Service{
		db:           db,
		jwtSecretEnv: jwtSecret,
	}
}

// NewRouter creates a http.Handler with routes related to authentication and
// user profiles.
func (s *Service) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /login", logger.Middleware(s.handleLoginGet, "Login (GET) handler"))
	mux.HandleFunc("POST /login", logger.Middleware(s.handleLoginPost, "Login (POST) handler"))
	mux.HandleFunc("GET /logout", logger.Middleware(s.handleLogoutPost, "Logout handler"))

	mux.HandleFunc("GET /register", logger.Middleware(s.handleRegisterGet, "Register (GET) handler"))
	mux.HandleFunc("POST /register", logger.Middleware(s.handleRegisterPost, ""))
	mux.HandleFunc("POST /register/validate", logger.Middleware(s.handleRegisterValidate, "Register validator handler"))

	mux.HandleFunc("GET /{userName}", s.handleProfilesGet)

	return mux
}

// Middleware wraps around special handlers that have the database user as an
// extra parameter. If authentication failed, a "guest" user is passed into the
// inner handler function.
func (s *Service) Middleware(h func(http.ResponseWriter, *http.Request, database.GetUserRow)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbUser, err := s.AuthJWTCookie(r)
		if err != nil {
			dbUser.UserName = "guest"
			dbUser.DisplayName = "guest"
		}

		h(w, r, dbUser)
	}
}
