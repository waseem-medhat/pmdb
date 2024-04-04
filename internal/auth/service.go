package auth

import (
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/logger"
)

type Service struct {
	db *database.Queries
}

func NewService(db *database.Queries) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /login", logger.Middleware(s.handleLoginGet, "Login (GET) handler"))
	mux.HandleFunc("POST /login", logger.Middleware(s.handleLoginPost, "Login (POST) handler"))
	mux.HandleFunc("GET /logout", logger.Middleware(s.handleLogoutPost, "Logout handler"))

	mux.HandleFunc("GET /register", logger.Middleware(s.HandleRegisterGet, "Register (GET) handler"))
	mux.HandleFunc("POST /register", logger.Middleware(s.HandleRegisterPost, ""))
	mux.HandleFunc("POST /register/validate", logger.Middleware(s.HandleRegisterValidate, "Register validator handler"))

	mux.HandleFunc("GET /{userName}", s.HandleProfilesGet)

	return mux
}

func (s *Service) MiddlewareAuth(h func(http.ResponseWriter, *http.Request, database.GetUserRow)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbUser, err := s.AuthJWTCookie(r)
		if err != nil {
			cookie := createCookie("pmdb-requested-url", r.RequestURI, "/users/login", 3600)
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
			return
		}

		h(w, r, dbUser)
	}
}
