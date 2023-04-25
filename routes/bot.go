package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func BotRoutes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Put("/", handlers.CreateOrEditBot)
		r.Get("/{name}", handlers.GetBotByUniqueName)
		r.Delete("/{uuid}", handlers.DeleteBot)
	})

	return r
}
