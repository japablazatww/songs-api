package handlers

import (
	"encoding/json"
	"net/http"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/japablazatww/song-searcher/cmd/api/internal/adapters/dto"
	"github.com/japablazatww/song-searcher/cmd/api/internal/application/services"
)

type SearchHandler struct {
	musicService *services.MusicService
}

func NewSearchHandler(musicService *services.MusicService) *SearchHandler {
	return &SearchHandler{musicService: musicService}
}

func (s *SearchHandler) SearchSong(w http.ResponseWriter, r *http.Request) {
	queryParams := dto.QueryParams{
		Song:   r.URL.Query().Get("song"),
		Album:  r.URL.Query().Get("album"),
		Artist: r.URL.Query().Get("artist"),
	}

	if IsEmptyQueryParams(queryParams) {
		http.Error(w, "Debes proporcionar al menos un parámetro de búsqueda (song, album o artist)", http.StatusBadRequest)
		return
	}

	songsResponse, err := s.musicService.Search(r.Context(), queryParams)
	if err != nil {
		http.Error(w, "Error al buscar canciones", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.Response{
		TotalSongs: len(songsResponse),
		Songs:      songsResponse,
	})
}

func IsEmptyQueryParams(params dto.QueryParams) bool {
	return params.Song == "" && params.Album == "" && params.Artist == ""
}
