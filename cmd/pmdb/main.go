// package main is the entry point of the PMDb app
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"github.com/wipdev-tech/pmdb/internal/auth"
	"github.com/wipdev-tech/pmdb/internal/database"
)

func main() {
	dbURL, dbToken, err := getDBEnv()
	if err != nil {
		log.Fatal(err)
	}

	connURL := fmt.Sprintf("%s?authToken=%s", dbURL, dbToken)
	db, err := sql.Open("libsql", connURL)
	if err != nil {
		log.Fatal(err)
	}

	dbConn := database.New(db)

	authService := auth.NewService(dbConn)
	authMux := authService.NewRouter()

	r := http.NewServeMux()
	r.Handle("/users/", http.StripPrefix("/users", authMux))
	server := &http.Server{
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Println("PMDb server let's Go! ")
	if os.Getenv("ENV") == "dev" {
		fmt.Println("Dev server started and running at http://localhost:8080")
		server.Addr = "localhost:8080"
	} else {
		fmt.Println("Server started and running")
		server.Addr = "0.0.0.0:" + os.Getenv("PORT")
	}
	log.Fatal(server.ListenAndServe())
}
