package ports

import (
	"context"

	"github.com/japablazatww/song-searcher/cmd/api/internal/adapters/dto"
	"github.com/japablazatww/song-searcher/cmd/api/internal/core/domain"
)

type MusicProvider interface {
	Search(ctx context.Context, query dto.QueryParams) ([]domain.Song, string, error)
}
