package services

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/japablazatww/song-searcher/cmd/api/internal/adapters/dto"
	"github.com/japablazatww/song-searcher/cmd/api/internal/core/domain"
	"github.com/japablazatww/song-searcher/cmd/api/internal/core/ports"
)

type MusicService struct {
	Repo      ports.SongRepository
	Providers []ports.MusicProvider
	Scorer    SongScorer
}

func (s *MusicService) Search(ctx context.Context, query dto.QueryParams) ([]domain.Song, error) {
	allSongs := make([]domain.Song, 0)
	var failedProviders []string
	songChan := make(chan []domain.Song)
	failedProviderChan := make(chan string, len(s.Providers))
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 1. Buscar en proveedores externos (con goroutines)
	for _, provider := range s.Providers {
		wg.Add(1)
		go func(p ports.MusicProvider) {
			defer wg.Done()
			songs, originName, err := p.Search(ctx, query)
			if err != nil {
				fmt.Println("dentro del error", originName)
				failedProviderChan <- originName
				fmt.Println(err)
				return
			}
			songChan <- songs

		}(provider)
	}

	// Recolectar resultados de los providers
	go func() {
		wg.Wait()
		close(songChan)
		close(failedProviderChan)
	}()

	for songs := range songChan {
		mu.Lock()
		allSongs = append(allSongs, songs...)
		mu.Unlock()
	}

	for originName := range failedProviderChan {
		mu.Lock()
		failedProviders = append(failedProviders, originName)
		mu.Unlock()
	}

	// 2. Buscar en la base de datos para los proveedores fallidos (con goroutine)
	for _, originName := range failedProviders {
		wg.Add(1)
		go func(origin string) {
			defer wg.Done()
			fmt.Println("Origin fallidos", origin)
			query.Origin = origin
			songs, err := s.Repo.Search(ctx, query)
			if err != nil {
				fmt.Println(err) // Usa un logger
				return
			}
			mu.Lock()
			allSongs = append(allSongs, songs...)
			mu.Unlock()
		}(originName)
	}
	wg.Wait()

	// 3. Ordenar resultados (sin goroutine, ya que es una operación rápida)
	sort.Slice(allSongs, func(i, j int) bool {
		scoreI := s.Scorer.CalculateScore(&allSongs[i], query)
		scoreJ := s.Scorer.CalculateScore(&allSongs[j], query)
		return scoreI > scoreJ
	})

	// 4. Guardar resultados ordenados (con goroutine)
	go func(songs []domain.Song) {
		saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Repo.Save(saveCtx, songs); err != nil {
			fmt.Println("Error en Save:", err)
		}
	}(allSongs)

	return allSongs, nil
}
