package providers

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/japablazatww/song-searcher/cmd/api/internal/adapters/dto"
	"github.com/japablazatww/song-searcher/cmd/api/internal/core/domain"
)

const chartlyricsOriginName = "chartlyrics"

type chartLyricsProvider struct {
	baseUrl string
}

type ChartLyricsResponse struct {
	XMLName xml.Name            `xml:"ArrayOfSearchLyricResult"`
	Results []ChartLyricsResult `xml:"SearchLyricResult"`
}

type ChartLyricsResult struct {
	TrackId   int    `xml:"TrackId"`
	LyricId   int    `xml:"LyricId"`
	SongUrl   string `xml:"SongUrl"`
	ArtistUrl string `xml:"ArtistUrl"`
	Artist    string `xml:"Artist"`
	Song      string `xml:"Song"`
	SongRank  int    `xml:"SongRank"`
}

func NewChartLyricsProvider() *chartLyricsProvider {
	return &chartLyricsProvider{
		baseUrl: "http://api.chartlyrics.com",
	}
}

func (c *chartLyricsProvider) Search(ctx context.Context, query dto.QueryParams) ([]domain.Song, string, error) {
	if queryContainsNameAndArtist(query) {
		return nil, chartlyricsOriginName, nil
	}

	// consulta a chartlyrics pasandole artist y song
	url := fmt.Sprintf("%s/apiv1.asmx/SearchLyric?artist=%s&song=%s", c.baseUrl, query.Artist, query.Song)
	fmt.Println("url", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, chartlyricsOriginName, fmt.Errorf("error creando la solicitud a chartlyrics: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, chartlyricsOriginName, fmt.Errorf("error al enviar la solicitud a chartlyrics: %w", err)
	}
	defer res.Body.Close()

	// Leer el cuerpo de la respuesta SOAP
	soapResponse, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, chartlyricsOriginName, fmt.Errorf("error leyendo la respuesta SOAP de chartlyrics: %w", err)
	}

	// Decodificar la respuesta SOAP en la estructura
	var chartLyricsResponse ChartLyricsResponse
	err = xml.Unmarshal(soapResponse, &chartLyricsResponse)
	if err != nil {
		return nil, chartlyricsOriginName, fmt.Errorf("error unmarshaling SOAP response: %w", err)
	}

	// Mapear la respuesta a la estructura de domain.Song
	song, err := c.mapSongFromChartLyricsResponse(chartLyricsResponse)
	if err != nil {
		return nil, chartlyricsOriginName, fmt.Errorf("error al pasar la información a domain.Song: %w", err)
	}

	return song[:len(song)-1], chartlyricsOriginName, nil
}

func queryContainsNameAndArtist(query dto.QueryParams) bool {
	return query.Song == "" && query.Artist == ""
}

func (c *chartLyricsProvider) mapSongFromChartLyricsResponse(chartLyricsResponse ChartLyricsResponse) ([]domain.Song, error) {
	var songs []domain.Song
	for _, result := range chartLyricsResponse.Results {
		fmt.Println("ID", result.LyricId+result.TrackId)
		song := domain.Song{
			ID:      result.LyricId + result.TrackId,
			Name:    result.Song,
			Artist:  result.Artist,
			Artwork: result.ArtistUrl,
			Origin:  chartlyricsOriginName,
		}
		songs = append(songs, song)
	}
	return songs, nil
}
