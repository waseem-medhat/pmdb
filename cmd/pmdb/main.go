package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("PMDb Server Let's Go! ")
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handleHome)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("ENV") == "dev" {
		fmt.Println("Dev server started and running at http://localhost:8080")
		log.Fatal(http.ListenAndServe("localhost:8080", nil))
	} else {
		fmt.Println("Server started and running")
		log.Fatal(http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), nil))
	}
}

// handleHome is the handler for the home route ("/")
func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/fragments.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
