package service

import (
	"log"
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/database"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

// getReviewData takes a slice of DB reviews and attached the TMDB data to
// them.
func getReviewData(reviews []database.GetReviewsRow) []tmdbapi.Review {
	type newDetails struct {
		title      string
		posterPath string
	}

	tmdbData := map[string]newDetails{}
	for _, r := range reviews {
		if _, ok := tmdbData[r.MovieTmdbID]; ok {
			continue
		}

		details, err := tmdbapi.GetMovieDetails(r.MovieTmdbID)
		if err != nil {
			log.Fatal(err)
		}

		tmdbData[r.MovieTmdbID] = newDetails{
			title:      details.Title,
			posterPath: details.PosterPath,
		}
	}

	newReviews := make([]tmdbapi.Review, 0, len(reviews))
	for _, r := range reviews {
		details := tmdbData[r.MovieTmdbID]
		newReview := tmdbapi.Review{
			GetReviewsRow: r,
			Title:         details.title,
			PosterPath:    details.posterPath,
		}
		newReviews = append(newReviews, newReview)
	}

	return newReviews
}

func (s *Service) MiddlewareAuth(h func(http.ResponseWriter, *http.Request, database.GetUserRow)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbUser, err := s.authJWTCookie(r)
		if err != nil {
			cookie := createCookie("pmdb-requested-url", r.URL.String(), "/login", 3600)
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		h(w, r, dbUser)
	}
}
