package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func ConnectionCodesRoutes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Post("/", handlers.CreateConnectionCode)
		r.Get("/", handlers.GetConnectionCode)
	})
	return r
}
