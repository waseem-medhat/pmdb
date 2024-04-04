package nowplaying

import (
	"net/http"

	"github.com/wipdev-tech/pmdb/internal/errors"
	"github.com/wipdev-tech/pmdb/internal/logger"
	"github.com/wipdev-tech/pmdb/internal/tmdbapi"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", logger.Middleware(s.handleNowPlayingGet, "Now Playing (GET) handler"))

	return mux
}

func (s *Service) handleNowPlayingGet(w http.ResponseWriter, r *http.Request) {
	nowPlaying, err := tmdbapi.GetNowPlaying(-1)
	if err != nil {
		errors.Render(w, http.StatusInternalServerError)
		return
	}

	templData := NowPlayingData{
		NowPlaying: nowPlaying,
	}

	err = NowPlaying(templData).Render(r.Context(), w)
	if err != nil {
		errors.Render(w, http.StatusInternalServerError)
		return
	}
}
