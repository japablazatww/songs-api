package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/japablazatww/song-searcher/cmd/api/internal/adapters/dto"
	"github.com/japablazatww/song-searcher/cmd/api/internal/core/domain"
)

const itunesOriginName = "apple"

type iTunesProvider struct {
	baseURL string
}

func NewiTunesProvider() *iTunesProvider {
	return &iTunesProvider{
		baseURL: "https://itunes.apple.com/search",
	}
}

func (p *iTunesProvider) Search(ctx context.Context, query dto.QueryParams) ([]domain.Song, string, error) {
	queryString := buildSearchTerm(query)

	url := fmt.Sprintf("%s?term=%s", p.baseURL, queryString)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, itunesOriginName, fmt.Errorf("error creando la solicitud a iTunes: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, itunesOriginName, fmt.Errorf("error al realizar la solicitud a iTunes: %w", err)
	}
	defer resp.Body.Close()

	var iTunesResp struct {
		Results []iTunesTrack `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&iTunesResp); err != nil {
		return nil, itunesOriginName, fmt.Errorf("error al decodificar la respuesta de iTunes: %w", err)
	}

	var songs []domain.Song
	for _, track := range iTunesResp.Results {
		songs = append(songs, transformToSong(track))
	}

	return songs, itunesOriginName, nil
}

func transformToSong(track iTunesTrack) domain.Song {
	return domain.Song{
		ID:       track.TrackID,
		Name:     track.TrackName,
		Artist:   track.ArtistName,
		Duration: formatDuration(track.TrackTimeMillis),
		Album:    track.CollectionName,
		Artwork:  track.ArtworkUrl100,
		Price:    fmt.Sprintf("%s %.2f", track.Currency, track.TrackPrice),
		Origin:   itunesOriginName,
	}
}

func formatDuration(ms int) string {
	minutes := ms / (1000 * 60)
	seconds := (ms / 1000) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func buildSearchTerm(params dto.QueryParams) string {
	var words []string

	addWords := func(s string) {
		if s != "" {
			words = append(words, strings.Fields(s)...)
		}
	}

	addWords(params.Song)
	addWords(params.Album)
	addWords(params.Artist)

	return strings.Join(words, "+")
}

type iTunesTrack struct {
	TrackID         int     `json:"trackId"`
	TrackName       string  `json:"trackName"`
	ArtistName      string  `json:"artistName"`
	TrackTimeMillis int     `json:"trackTimeMillis"`
	CollectionName  string  `json:"collectionName"`
	ArtworkUrl100   string  `json:"artworkUrl100"`
	TrackPrice      float64 `json:"trackPrice"`
	Currency        string  `json:"currency"`
}
