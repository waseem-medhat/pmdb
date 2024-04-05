// Package errors holds a helper and a template for rendering HTTP response
// errors
package errors

import (
	"context"
	"log"
	"net/http"
)

// Render returns an HTML error page with the response status code with the
// text (e.g., "not found")
func Render(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	err := ErrorPage(statusCode).Render(context.Background(), w)
	if err != nil {
		log.Fatal(err)
	}
}
