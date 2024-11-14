package dto

import "github.com/japablazatww/song-searcher/cmd/api/internal/core/domain"

type Response struct {
	TotalSongs int           `json:"totalSongs"`
	Songs      []domain.Song `json:"songs"`
}
