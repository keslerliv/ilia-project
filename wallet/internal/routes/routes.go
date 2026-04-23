package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/keslerliv/wallet/internal/handlers"
	"github.com/keslerliv/wallet/internal/middlewares"
)

func LoadRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Route("/wallet", func(r chi.Router) {
			r.Get("/", handlers.GetBalance)
			r.Post("/", handlers.PostValue)
		})
	})

	return router
}
