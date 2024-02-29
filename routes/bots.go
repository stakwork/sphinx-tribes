package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func BotsRoutes() chi.Router {
	r := chi.NewRouter()
	botHandler := handlers.NewBotHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Post("/", handlers.CreateOrEditBot)
		r.Get("/", handlers.GetListedBots)
		r.Get("/owner/{pubkey}", botHandler.GetBotsByOwner)
		r.Get("/{uuid}", handlers.GetBot)
	})
	return r
}
