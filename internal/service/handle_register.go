package service

import (
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/wipdev-tech/pmdb/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) HandleRegisterGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/register.html",
		"templates/blocks/_top.html",
		"templates/blocks/_bottom.html",
	))
	err := tmpl.Execute(w, struct{ ErrorMsgs []string }{})
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleRegisterPost(w http.ResponseWriter, r *http.Request) {
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
		DisplayName string
	}{
		DisplayName: dbUser.DisplayName,
	}

	tmpl := template.Must(template.ParseFiles("templates/register-success.html"))
	err = tmpl.ExecuteTemplate(w, "register-success-main", tmplData)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleRegisterValidate(w http.ResponseWriter, r *http.Request) {
	errorMsgs := []string{}

	inDisplayName := r.PostFormValue("display-name")
	if inDisplayName == "" {
		errorMsgs = append(errorMsgs, "Full name shouldn't be empty.")
	}

	dbUsers, err := s.DB.ListUsers(r.Context())
	if err != nil {
		log.Fatal(err)
	}

	inUserName := r.PostFormValue("user-name")
	if inUserName == "" {
		errorMsgs = append(errorMsgs, "User name shouldn't be empty.")
	}

	inPassword := r.PostFormValue("password")
	if len(inPassword) < 8 {
		errorMsgs = append(errorMsgs, "Password must be at least 8 characters long.")
	}

	inConfirmPassword := r.PostFormValue("confirm-password")
	if inConfirmPassword != inPassword {
		errorMsgs = append(errorMsgs, "Passwords should match.")
	}

	for _, u := range dbUsers {
		if u.UserName == inUserName {
			errorMsgs = append(errorMsgs, "User name already exists.")
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	err = tmpl.ExecuteTemplate(w, "check", struct{ ErrorMsgs []string }{ErrorMsgs: errorMsgs})
	if err != nil {
		log.Fatal(err)
	}
}
