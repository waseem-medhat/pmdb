package router

import (
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/service"
)

func New(s *service.Service) *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", s.HandleHome)
	mux.HandleFunc("/register", s.HandleRegister)

	mux.HandleFunc("/create-user", s.HandleCreateUser)

	return mux
}
