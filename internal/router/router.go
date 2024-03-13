package router

import (
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/service"
)

func New(s *service.Service) *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("GET /{$}", s.HandleHome)
	mux.HandleFunc("GET /register", s.HandleRegister)
	mux.HandleFunc("GET /login", s.HandleLoginGet)

	mux.HandleFunc("POST /create-user", s.HandleCreateUser)
	mux.HandleFunc("POST /register/validate", s.HandleValidateRegisterForm)
	mux.HandleFunc("POST /login", s.HandleLoginPost)

	return mux
}
