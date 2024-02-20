package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func getDBEnv() (string, string, error) {
	dbURL := ""
	dbToken := ""
	err := godotenv.Load()

	// In the Render prod deployment there is no .env by design
	// So will error out locally but skip for prod
	if os.Getenv("ENV") == "dev" && os.IsNotExist(err) {
		return dbURL, dbToken, err
	}
	err = nil

	dbURL = os.Getenv("DBURL")
	if dbURL == "" {
		err = fmt.Errorf("DBURL environment variable is not set")
		return dbURL, dbToken, err
	}

	dbToken = os.Getenv("DBTOKEN")
	if dbURL == "" {
		err = fmt.Errorf("DBTOKEN environment variable is not set")
		return dbURL, dbToken, err
	}

	return dbURL, dbToken, err
}
