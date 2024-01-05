package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func ConnectionCodesRoutes() chi.Router {
	r := chi.NewRouter()
	authHandler := handlers.NewAuthHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Post("/", authHandler.CreateConnectionCode)
		r.Get("/", authHandler.GetConnectionCode)
	})
	return r
}
