package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func TribeRoutes() chi.Router {
	r := chi.NewRouter()
	tribeHandlers := handlers.NewTribeHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Get("/", tribeHandlers.GetListedTribes)
		r.Get("/app_url/{app_url}", tribeHandlers.GetTribesByAppUrl)
		r.Get("/app_urls/{app_urls}", handlers.GetTribesByAppUrls)
		r.Get("/{uuid}", tribeHandlers.GetTribe)
		r.Get("/total", tribeHandlers.GetTotalribes)
		r.Post("/", tribeHandlers.CreateOrEditTribe)
	})
	return r
}
