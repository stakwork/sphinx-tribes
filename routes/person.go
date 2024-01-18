package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func PersonRoutes() chi.Router {
	r := chi.NewRouter()
	peopleHandler := handlers.NewPeopleHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Get("/{pubkey}", peopleHandler.GetPersonByPubkey)
		r.Get("/id/{id}", handlers.GetPersonById)
		r.Get("/uuid/{uuid}", handlers.GetPersonByUuid)
		r.Get("/uuid/{uuid}/assets", handlers.GetPersonAssetsByUuid)
		r.Get("/githubname/{github}", handlers.GetPersonByGithubName)
	})
	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Post("/", peopleHandler.CreateOrEditPerson)
		r.Delete("/{id}", handlers.DeletePerson)
	})
	return r
}
