// package main is the entry point of the PMDb app
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/router"
	"github.com/wipdev-tech/pmdb/internal/service"
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

	s := service.New()
	s.DB = database.New(db)
	r := router.New(s)
	server := http.Server{Handler: r}

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
