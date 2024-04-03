package auth

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) handleLoginGet(w http.ResponseWriter, r *http.Request) {
	err := LoginPage(LoginPageData{LoginError: false}).Render(r.Context(), w)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) handleLoginPost(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("user-name")
	password := r.FormValue("password")

	dbUser, err := s.DB.GetUserForLogin(r.Context(), userName)
	if err == sql.ErrNoRows {
		err := ErrorAlert(LoginPageData{LoginError: true}).Render(r.Context(), w)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if err != nil {
		err := ErrorAlert(LoginPageData{LoginError: true}).Render(r.Context(), w)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	access, err := generateJWTAccess(userName)
	if err != nil {
		log.Fatal(err)
	}

	cookie := createCookie("pmdb-jwt-access", access, "/", 3600)
	http.SetCookie(w, cookie)

	redirectCookie, err := r.Cookie("pmdb-requested-url")
	var redirectTo string
	if err != nil || redirectCookie.Value == "" {
		redirectTo = "/"
	} else {
		redirectTo = redirectCookie.Value
	}

	newRedirectCookie := createCookie("pmdb-requested-url", "", "/login", -1)
	http.SetCookie(w, newRedirectCookie)
	w.Header().Set("HX-Redirect", redirectTo)
	w.WriteHeader(http.StatusFound)
}

func (s *Service) handleLogoutPost(w http.ResponseWriter, _ *http.Request) {
	cookie := createCookie("pmdb-jwt-access", "", "/", -1)
	http.SetCookie(w, cookie)

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusFound)
}

func (s *Service) HandleRegisterGet(w http.ResponseWriter, r *http.Request) {
	err := Register(RegisterData{ErrorMsgs: []string{}}).Render(r.Context(), w)
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

	err = RegisterSuccessHX(RegisterSuccessData{
		DisplayName: dbUser.DisplayName,
	}).Render(r.Context(), w)
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

	err = RegisterErrors(RegisterData{ErrorMsgs: errorMsgs}).Render(r.Context(), w)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) HandleProfilesGet(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("userName")
	dbUser, err := s.DB.GetUser(r.Context(), userName)
	if err == sql.ErrNoRows {
		errors.Render(w, http.StatusNotFound)
		return
	}

	if err != nil {
		fmt.Println("couldn't get user - ", err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}

	err = Profile(ProfileData{User: dbUser}).Render(r.Context(), w)
	if err != nil {
		fmt.Println(err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}
}
