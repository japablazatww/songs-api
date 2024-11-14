package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/japablazatww/song-searcher/cmd/api/internal/adapters/inbound/http/handlers"
	"github.com/japablazatww/song-searcher/cmd/api/internal/adapters/inbound/http/routes"
	"github.com/japablazatww/song-searcher/cmd/api/internal/adapters/outbound/providers"
	"github.com/japablazatww/song-searcher/cmd/api/internal/adapters/outbound/repositories"
	"github.com/japablazatww/song-searcher/cmd/api/internal/application/services"
	"github.com/japablazatww/song-searcher/cmd/api/internal/core/ports"
	database "github.com/japablazatww/song-searcher/cmd/api/internal/infraestructure/repositories"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Inicializa repositorios
	postgresRepo := repositories.NewPostgresRepository(database.PostgresDB.DB)

	// Inicializa providers
	providers := []ports.MusicProvider{
		providers.NewiTunesProvider(),
		providers.NewChartLyricsProvider(),
	}

	// Inicializa configuracion del Scoring Service
	scoringWeightsconfig, err := services.LoadScoringWeights()
	if err != nil {
		log.Fatal(err)
	}

	// Inicializa Scoring Service
	scorer := services.SongScorer{
		Weights: scoringWeightsconfig,
	}

	// Inicializa Music Service
	musicService := &services.MusicService{
		Repo:      postgresRepo,
		Providers: providers,
		Scorer:    scorer,
	}

	// Inicializa servicios
	authService := services.NewAuthService(postgresRepo)

	// Inicializa handlers
	authHandler := handlers.NewAuthHandler(authService)
	searchHandler := handlers.NewSearchHandler(musicService)

	// Crear el HandlerContainer
	handlers := &routes.HandlerContainer{
		SearchHandler: searchHandler,
		AuthHandler:   authHandler,
		AuthService:   authService,
	}

	// Injectar el HandlerContainer en las rutas
	router := routes.Routes(handlers)

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	fmt.Printf("Starting server on port %s\n", port)
	err = svr.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
