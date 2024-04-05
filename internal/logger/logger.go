// Package logger holds tooling for logging
package logger

import (
	"log"
	"net/http"
)

// Middleware wraps around a http.HandlerFunc to log some basic request info to
// stdout along with a custom message.
func Middleware(h http.HandlerFunc, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v '%v' | %v", r.Method, r.URL.String(), message)
		h(w, r)
	}
}
