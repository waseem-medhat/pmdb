package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/wipdev-tech/pmdb/internal/database"
)

type environment struct {
	dev       bool
	dbURL     string
	jwtSecret string
	tmdbToken string
	port      string
}

func loadEnv() (environment, error) {
	var env environment
	err := godotenv.Load()

	// In the Render prod deployment there is no .env by design
	// So will error out locally but skip for prod
	env.dev = os.Getenv("ENV") == "dev"
	if env.dev && os.IsNotExist(err) {
		return env, err
	}
	err = nil

	env.dbURL = os.Getenv("DBURL")
	if env.dbURL == "" {
		err = fmt.Errorf("DBURL environment variable is not set")
		return env, err
	}

	env.jwtSecret = os.Getenv("JWT_SECRET")
	if env.jwtSecret == "" {
		err = fmt.Errorf("JWT_SECRET environment variable is not set")
		return env, err
	}

	env.tmdbToken = os.Getenv("TMDB_TOKEN")
	if env.tmdbToken == "" {
		err = fmt.Errorf("TMDB_TOKEN environment variable is not set")
		return env, err
	}

	env.port = os.Getenv("PORT")
	if env.port == "" {
		env.port = "8080"
	}

	return env, err
}

func initDB(dbURL string) *database.Queries {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error opening db connection ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("error pinging db ", err)
	}

	return database.New(db)
}
