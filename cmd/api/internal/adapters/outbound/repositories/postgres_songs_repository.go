package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/japablazatww/song-searcher/cmd/api/internal/adapters/dto"
	"github.com/japablazatww/song-searcher/cmd/api/internal/core/domain"
)

func (r *PostgresRepository) Save(ctx context.Context, songs []domain.Song) error {
	query := `
        INSERT INTO songs (id, name, artist, duration, album, artwork, price, origin)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			artist = EXCLUDED.artist,
			duration = EXCLUDED.duration,
			album = EXCLUDED.album,
			artwork = EXCLUDED.artwork,
			price = EXCLUDED.price,
			origin = EXCLUDED.origin;
    `
	for _, song := range songs {
		_, err := r.DB.ExecContext(ctx, query,
			song.ID, song.Name, song.Artist, song.Duration, song.Album, song.Artwork, song.Price, song.Origin,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresRepository) Search(ctx context.Context, params dto.QueryParams) ([]domain.Song, error) {
	query := `
        SELECT 
            id, name, artist, duration, album, artwork, price, origin
        FROM songs
        WHERE 1=1
    `
	var args []interface{}
	var filterClauses []string

	if params.Artist != "" {
		filterClauses = append(filterClauses, "artist ILIKE $"+fmt.Sprint(len(args)+1))
		args = append(args, "%"+params.Artist+"%")
	}

	if params.Album != "" {
		filterClauses = append(filterClauses, "album ILIKE $"+fmt.Sprint(len(args)+1))
		args = append(args, "%"+params.Album+"%")
	}

	if params.Song != "" {
		filterClauses = append(filterClauses, "name ILIKE $"+fmt.Sprint(len(args)+1))
		args = append(args, "%"+params.Song+"%")
	}

	// Concatena las condiciones OR entre parentesis
	if len(filterClauses) > 0 {
		query += " AND (" + strings.Join(filterClauses, " OR ") + ")"
	}

	// Agrega la condicion AND al final para garantizar el origen
	if params.Origin != "" {
		query += " AND origin = $" + fmt.Sprint(len(args)+1)
		args = append(args, params.Origin)
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []domain.Song
	for rows.Next() {
		var song domain.Song
		if err := rows.Scan(
			&song.ID, &song.Name, &song.Artist, &song.Duration, &song.Album, &song.Artwork, &song.Price, &song.Origin,
		); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	fmt.Printf("Found %+v songs on db\n", songs)

	return songs, nil
}
