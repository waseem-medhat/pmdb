package reviews

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/errors"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

func (s *Service) handleReviewsGet(w http.ResponseWriter, r *http.Request, user database.GetUserRow) {
	reviews, err := s.db.GetReviews(r.Context())
	if err != nil {
		fmt.Println(err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}

	templData := ReviewsPageData{
		Reviews: s.tmdb.GetReviewMovieDetails(reviews),
		User:    user,
	}

	err = ReviewsPage(templData).Render(r.Context(), w)
	if err != nil {
		fmt.Println(err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}
}

func (s *Service) handleReviewsGetByID(w http.ResponseWriter, r *http.Request, user database.GetUserRow) {
	reviewIDStr := r.PathValue("reviewID")
	if reviewIDStr == "" {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}

	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}

	review, err := s.db.GetReviewByID(r.Context(), reviewID)
	if err == sql.ErrNoRows {
		fmt.Println(err)
		errors.Render(w, http.StatusNotFound)
		return
	}

	review.Review = strings.ReplaceAll(review.Review, "\\n", "\n")

	if err != nil {
		fmt.Println(err)
		errors.Render(w, http.StatusInternalServerError)
		return
	}

	movieDetails, err := s.tmdb.GetMovieDetails(review.MovieTmdbID)
	if tmdbapi.IsNotFound(err) {
		errors.Render(w, http.StatusNotFound)
		return
	}

	if err != nil {
		errors.Render(w, http.StatusInternalServerError)
		return
	}

	templData := ReviewPageData{
		User:   user,
		Review: review,
		Movie:  movieDetails,
	}

	err = ReviewPage(templData).Render(r.Context(), w)
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

	movieID := r.URL.Query().Get("movieID")
	if movieID == "" {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	rating, err := strconv.ParseInt(r.FormValue("rating"), 0, 32)
	if err != nil || rating < 0 || rating > 10 {
		errors.Render(w, http.StatusBadRequest)
		return
	}

	publicReview := r.FormValue("public-review") == "on"

	review := strings.ReplaceAll(r.FormValue("review"), "\n", "\\n")

	_, err = s.db.CreateReview(r.Context(), database.CreateReviewParams{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		UserID:       dbUser.ID,
		MovieTmdbID:  movieID,
		Rating:       int32(rating), // #nosec G109
		Review:       review,
		PublicReview: publicReview,
	})
	if err != nil {
		errors.Render(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusFound)
}
