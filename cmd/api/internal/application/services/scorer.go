package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/japablazatww/song-searcher/cmd/api/internal/adapters/dto"
	"github.com/japablazatww/song-searcher/cmd/api/internal/core/domain"
)

type SongScorer struct {
	Weights ScoringWeights
}

type ScoringWeights struct {
	ExactMatchName     int
	PartialMatchName   int
	ExactMatchArtist   int
	PartialMatchArtist int
	ExactMatchAlbum    int
	PartialMatchAlbum  int
	HasArtwork         int
	HasPrice           int
	Origin             map[string]int
}

func loadOriginWeights(filePath string) (map[string]int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir archivo de configuración: %w", err)
	}
	defer file.Close()

	var weights map[string]int
	if err := json.NewDecoder(file).Decode(&weights); err != nil {
		return nil, fmt.Errorf("error al decodificar el JSON: %w", err)
	}

	return weights, nil
}

func LoadScoringWeights() (ScoringWeights, error) {

	originWeights, err := loadOriginWeights("./origin_weights.json")
	if err != nil {
		fmt.Println("Error cargando los pesos de origen:", err)
		originWeights = make(map[string]int)
	}

	return ScoringWeights{
		ExactMatchName:     10,
		PartialMatchName:   5,
		ExactMatchArtist:   8,
		PartialMatchArtist: 4,
		ExactMatchAlbum:    6,
		PartialMatchAlbum:  3,
		HasArtwork:         3,
		HasPrice:           2,
		Origin:             originWeights,
	}, nil
}

func (s *SongScorer) CalculateScore(song *domain.Song, query dto.QueryParams) int {
	score := 0

	// Coincidencia de la canción
	if query.Song != "" {
		if strings.EqualFold(song.Name, query.Song) {
			score += s.Weights.ExactMatchName
		} else if strings.Contains(strings.ToLower(song.Name), strings.ToLower(query.Song)) {
			score += s.Weights.PartialMatchName
		}
	}

	// Coincidencia del artista
	if query.Artist != "" {
		if strings.EqualFold(song.Artist, query.Artist) {
			score += s.Weights.ExactMatchArtist
		} else if strings.Contains(strings.ToLower(song.Artist), strings.ToLower(query.Artist)) {
			score += s.Weights.PartialMatchArtist
		}
	}

	// Coincidencia del álbum
	if query.Album != "" {
		if strings.EqualFold(song.Album, query.Album) {
			score += s.Weights.ExactMatchAlbum
		} else if strings.Contains(strings.ToLower(song.Album), strings.ToLower(query.Album)) {
			score += s.Weights.PartialMatchAlbum
		}
	}

	// Calidad de datos
	if song.Artwork != "" {
		score += s.Weights.HasArtwork
	}
	if song.Price != "" {
		score += s.Weights.HasPrice
	}

	// Bonus por origen
	if originScore, exists := s.Weights.Origin[song.Origin]; exists {
		score += originScore
	}

	return score
}
