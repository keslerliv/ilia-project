package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/keslerliv/ilia-project/user/internal/handlers"
	"github.com/keslerliv/ilia-project/user/internal/middlewares"
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

		r.Route("/users", func(r chi.Router) {
			r.Post("/", handlers.PostUser)

			r.Group(func(r chi.Router) {
				r.Use(middlewares.Auth)

				r.Get("/", handlers.GetUsers)
				r.Get("/{id}", handlers.GetUser)
				r.Patch("/{id}", handlers.PatchUser)
				r.Delete("/{id}", handlers.DeleteUser)
			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/", handlers.Auth)
		})
	})

	return router
}
