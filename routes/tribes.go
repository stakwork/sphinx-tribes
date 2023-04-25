package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func TribeRoutes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get("/", handlers.GetListedTribes)
		r.Get("/{uuid}", handlers.GetTribe)
		r.Get("/total", handlers.GetTotalribes)
		r.Post("/", handlers.CreateOrEditTribe)
	})
	return r
}
