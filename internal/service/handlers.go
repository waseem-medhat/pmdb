package service

import (
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/wipdev-tech/pmdb/internal/database"
	"golang.org/x/crypto/bcrypt"
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

	tmpl := template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/_top.html",
		"templates/_bottom.html",
	))
	err = tmpl.Execute(w, tmplData)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	displayName := r.PostFormValue("display-name")
	userName := r.PostFormValue("user-name")
	password := r.PostFormValue("password")
	confirmPassword := r.PostFormValue("confirm-password")

	if password != confirmPassword {
		log.Fatal("passwords don't match")
		return
	}

	hPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Fatal("couldn't hash password")
		return
	}

	dbUser, err := s.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:          uuid.NewString(),
		UserName:    userName,
		DisplayName: displayName,
		Password:    string(hPassword),
	})
	if err != nil {
		log.Fatal(err)
	}

	tmplData := struct {
		User database.User
	}{
		User: dbUser,
	}

	tmpl := template.Must(template.ParseFiles("templates/hx_register_success.html"))
	err = tmpl.Execute(w, tmplData)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleRegister(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/register.html",
		"templates/_top.html",
		"templates/_bottom.html",
	))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
