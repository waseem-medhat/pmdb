package reviews

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/errors"
)

func (s *Service) handleReviewsGet(w http.ResponseWriter, r *http.Request) {
	reviews, err := s.db.GetReviews(r.Context())
	if err != nil {
		fmt.Println(err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}

	templData := ReviewsPageData{Reviews: s.tmdb.GetReviewMovieDetails(reviews)}

	err = ReviewsPage(templData).Render(r.Context(), w)
	if err != nil {
		fmt.Println(err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}
}

func (s *Service) handleReviewsNewGet(w http.ResponseWriter, r *http.Request, user database.GetUserRow) {
	if user.UserName == "guest" {
		cookie := s.auth.CreateCookie("pmdb-requested-url", r.RequestURI, "/users/login", 3600)
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	movieID := r.URL.Query().Get("movieID")
	if movieID == "" {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	movieDetails, err := s.tmdb.GetMovieDetails(movieID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	err = NewReviewPage(NewReviewPageData{Movie: movieDetails}).Render(r.Context(), w)
	if err != nil {
		errors.Render(w, http.StatusInternalServerError)
		return
	}

}

func (s *Service) handleReviewsNewPost(w http.ResponseWriter, r *http.Request, dbUser database.GetUserRow) {
	if dbUser.UserName == "guest" {
		cookie := s.auth.CreateCookie("pmdb-requested-url", r.RequestURI, "/users/login", 3600)
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	var dbTimeLayout = "2006-01-02 15:04"

	movieID := r.URL.Query().Get("movieID")
	if movieID == "" {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	rating, err := strconv.Atoi(r.FormValue("rating"))
	if err != nil {
		errors.Render(w, http.StatusBadRequest)
		return
	}

	var publicReview int64
	if r.FormValue("public-review") == "on" {
		publicReview = 1
	}

	_, err = s.db.CreateReview(r.Context(), database.CreateReviewParams{
		ID:           uuid.NewString(),
		CreatedAt:    time.Now().Format(dbTimeLayout),
		UpdatedAt:    time.Now().Format(dbTimeLayout),
		UserID:       dbUser.ID,
		MovieTmdbID:  movieID,
		Rating:       int64(rating),
		Review:       r.FormValue("review"),
		PublicReview: publicReview,
	})
	if err != nil {
		errors.Render(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusFound)
}
