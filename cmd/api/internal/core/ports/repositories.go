package ports

import (
	"context"

	"github.com/japablazatww/song-searcher/cmd/api/internal/adapters/dto"
	"github.com/japablazatww/song-searcher/cmd/api/internal/core/domain"
)

type SongRepository interface {
	Save(ctx context.Context, songs []domain.Song) error
	Search(ctx context.Context, query dto.QueryParams) ([]domain.Song, error)
}

type AuthRepository interface {
	CreateApp(app *domain.AppAuth) error
	GetAppByClientCredentials(clientID, clientSecret string) (*domain.AppAuth, error)
	UpdateAppToken(clientID string) error
	ValidateToken(token string) error
}
