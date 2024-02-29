package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func BotRoutes() chi.Router {
	r := chi.NewRouter()
	botHandler := handlers.NewBotHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Get("/{name}", botHandler.GetBotByUniqueName)
	})
	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Delete("/{uuid}", botHandler.DeleteBot)
		r.Put("/", botHandler.CreateOrEditBot)
	})
	return r
}
