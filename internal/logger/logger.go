package logger

import (
	"log"
	"net/http"
)

func Middleware(h http.HandlerFunc, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v '%v' | %v", r.Method, r.URL.String(), message)
		h(w, r)
	}
}
