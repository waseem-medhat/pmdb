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
	"github.com/wipdev-tech/pmdb/internal/home"
	"github.com/wipdev-tech/pmdb/internal/movies"
	"github.com/wipdev-tech/pmdb/internal/nowplaying"
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
	nowPlayingService := nowplaying.NewService()
	homeService := home.NewService(authService, dbConn)
	movieService := movies.NewService()

	r := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	r.Handle("GET /static/", http.StripPrefix("/static/", fs))
	r.Handle("/", homeService.NewRouter())
	r.Handle("/users/", http.StripPrefix("/users", authService.NewRouter()))
	r.Handle("/now-playing/", http.StripPrefix("/now-playing", nowPlayingService.NewRouter()))
	r.Handle("/movies/", http.StripPrefix("/movies", movieService.NewRouter()))

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
