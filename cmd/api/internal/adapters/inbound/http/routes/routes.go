package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/japablazatww/song-searcher/cmd/api/internal/adapters/inbound/http/handlers"
	authMiddleware "github.com/japablazatww/song-searcher/cmd/api/internal/adapters/inbound/http/middleware"
	"github.com/japablazatww/song-searcher/cmd/api/internal/application/services"
)

type HandlerContainer struct {
	SearchHandler *handlers.SearchHandler
	AuthHandler   *handlers.AuthHandler
	AuthService   *services.AuthService
}

func Routes(h *HandlerContainer) http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.AuthHandler.CreateApp)
		r.Post("/token", h.AuthHandler.GenerateToken)
	})

	// Rutas protegidas
	mux.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthMiddleware(h.AuthService))
		r.Get("/search", h.SearchHandler.SearchSong)
	})

	return mux
}
