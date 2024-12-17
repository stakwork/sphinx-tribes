package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
	"github.com/stakwork/sphinx-tribes/utils"
)

func PersonRoutes() chi.Router {
	r := chi.NewRouter()
	peopleHandler := handlers.NewPeopleHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Get("/{pubkey}", peopleHandler.GetPersonByPubkey)
		r.Get("/id/{id}", peopleHandler.GetPersonById)
		r.Get("/uuid/{uuid}", utils.TraceWithLogging(peopleHandler.GetPersonByUuid))
		r.Get("/uuid/{uuid}/assets", handlers.GetPersonAssetsByUuid)
		r.Get("/githubname/{github}", handlers.GetPersonByGithubName)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.CypressContext)
		r.Post("/upsertlogin", peopleHandler.UpsertLogin)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Post("/", peopleHandler.CreatePerson)
		r.Put("/", peopleHandler.UpdatePerson)
		r.Delete("/{id}", peopleHandler.DeletePerson)
	})
	return r
}
