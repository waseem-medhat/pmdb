// package main is the entry point of the PMDb app
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wipdev-tech/pmdb/internal/auth"
	"github.com/wipdev-tech/pmdb/internal/home"
	"github.com/wipdev-tech/pmdb/internal/movies"
	"github.com/wipdev-tech/pmdb/internal/nowplaying"
	"github.com/wipdev-tech/pmdb/internal/reviews"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

func main() {
	env, err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	dbConn := initDB(env.dbURL, env.dbToken)
	authService := auth.NewService(dbConn, env.jwtSecret)
	tmdbService := tmdbapi.NewService(env.tmdbToken)
	nowPlayingService := nowplaying.NewService(authService, tmdbService)
	homeService := home.NewService(authService, tmdbService, dbConn)
	movieService := movies.NewService(authService, tmdbService, dbConn)
	reviewService := reviews.NewService(authService, tmdbService, dbConn)

	r := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	r.Handle("GET /static/", http.StripPrefix("/static/", fs))
	r.Handle("/", homeService.NewRouter())
	r.Handle("/users/", http.StripPrefix("/users", authService.NewRouter()))
	r.Handle("/now-playing/", http.StripPrefix("/now-playing", nowPlayingService.NewRouter()))
	r.Handle("/movies/", http.StripPrefix("/movies", movieService.NewRouter()))
	r.Handle("/reviews/", http.StripPrefix("/reviews", reviewService.NewRouter()))

	server := &http.Server{
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Println("PMDb server let's Go! ")
	if env.dev {
		fmt.Println("Dev server started and running at http://localhost:" + env.port)
		server.Addr = "localhost:" + env.port
	} else {
		fmt.Println("Server started and running")
		server.Addr = "0.0.0.0:" + env.port
	}
	log.Fatal(server.ListenAndServe())
}
