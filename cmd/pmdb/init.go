package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	dev       bool
	dbURL     string
	dbToken   string
	jwtSecret string
	tmdbToken string
	port      string
}

func loadEnv() (Env, error) {
	var env Env
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

	env.dbToken = os.Getenv("DBTOKEN")
	if env.dbToken == "" {
		err = fmt.Errorf("DBTOKEN environment variable is not set")
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
