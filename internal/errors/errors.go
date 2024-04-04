package errors

import (
	"context"
	"log"
	"net/http"
)

func Render(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	err := ErrorPage(statusCode).Render(context.Background(), w)
	if err != nil {
		log.Fatal(err)
	}
}
