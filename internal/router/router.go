package router

import (
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/service"
)

func New(s *service.Service) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /reviews/new", s.MiddlewareAuth(s.HandleReviewsNewGet))
	mux.HandleFunc("POST /reviews/new", s.MiddlewareAuth(s.HandleReviewsNewPost))

	return mux
}
