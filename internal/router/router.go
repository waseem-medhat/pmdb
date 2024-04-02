package router

import (
	"log"
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/service"
)

func New(s *service.Service) *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("GET /{$}", mwLogger(s.HandleHome, "Home handler"))

	mux.HandleFunc("GET /register", mwLogger(s.HandleRegisterGet, "Register (GET) handler"))
	mux.HandleFunc("POST /register", mwLogger(s.HandleRegisterPost, ""))
	mux.HandleFunc("POST /register/validate", mwLogger(s.HandleRegisterValidate, ""))

	mux.HandleFunc("GET /login", mwLogger(s.HandleLoginGet, "Login (GET) handler"))
	mux.HandleFunc("POST /login", mwLogger(s.HandleLoginPost, "Login (POST) handler"))
	mux.HandleFunc("GET /logout", mwLogger(s.HandleLogoutPost, "Logout handler"))

	mux.HandleFunc("GET /profiles/{userName}", mwLogger(s.HandleProfilesGet, "Profiles handler"))
	mux.HandleFunc("GET /movies/{movieID}", mwLogger(s.HandleMoviesGet, "Movies handler"))
	mux.HandleFunc("GET /now-playing", mwLogger(s.HandleNowPlayingGet, "Now Playing handler"))

	mux.HandleFunc("GET /reviews/new", mwLogger(s.MiddlewareAuth(s.HandleReviewsNewGet), "New Review (GET) handler"))
	mux.HandleFunc("POST /reviews/new", mwLogger(s.MiddlewareAuth(s.HandleReviewsNewPost), "New Review (POST) handler"))

	return mux
}

func mwLogger(h http.HandlerFunc, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v '%v' | %v", r.Method, r.URL.String(), message)
		h(w, r)
	}
}
