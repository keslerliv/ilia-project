package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/keslerliv/ilia-project/wallet/internal/handlers"
	"github.com/keslerliv/ilia-project/wallet/internal/middlewares"
)

func LoadRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	router.Route("/", func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Route("/transactions", func(r chi.Router) {
			r.Get("/", handlers.GetTransactions)
			r.Post("/", handlers.PostTransaction)
		})

		r.Route("/balance", func(r chi.Router) {
			r.Get("/", handlers.GetBalance)
		})
	})

	return router
}
