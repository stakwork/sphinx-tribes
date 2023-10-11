package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func BotsRoutes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Post("/", handlers.CreateOrEditBot)
		r.Get("/", handlers.GetListedBots)
		r.Get("/owner/{pubkey}", handlers.GetBotsByOwner)
		r.Get("/{uuid}", handlers.GetBot)
	})
	return r
}
