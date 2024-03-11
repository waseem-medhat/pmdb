package service

import (
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/wipdev-tech/pmdb/internal/database"
)

// HandleHome is the handler for the home route ("/")
func (s *Service) HandleHome(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := s.DB.ListUsers(r.Context())
	if err != nil {
		log.Fatal(err)
	}

	tmplData := struct {
		Users []database.User
	}{
		Users: dbUsers,
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/fragments.html"))
	err = tmpl.Execute(w, tmplData)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	newDisplayName := r.PostFormValue("display-name")
	newUserName := r.PostFormValue("user-name")
	_, err := s.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:          uuid.NewString(),
		UserName:    newUserName,
		DisplayName: newDisplayName,
	})
	if err != nil {
		log.Fatal(err)
	}

	errMsg := ""
	if err != nil {
		errMsg = "Could not add user :("
	}

	dbUsers, err := s.DB.ListUsers(r.Context())
	tmplData := struct {
		Users        []database.User
		ErrorMessage string
	}{
		Users:        dbUsers,
		ErrorMessage: errMsg,
	}

	tmpl := template.Must(template.ParseFiles("templates/users.html"))
	err = tmpl.Execute(w, tmplData)
	if err != nil {
		log.Fatal(err)
	}
}
