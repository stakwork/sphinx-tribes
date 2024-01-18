package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func PersonRoutes() chi.Router {
	r := chi.NewRouter()
	personHandler := handlers.NewPersonyHandler(http.DefaultClient, db.DB)
	r.Group(func(r chi.Router) {
		r.Get("/{pubkey}", handlers.GetPersonByPubkey)
		r.Get("/id/{id}", personHandler.GetPersonById)
		r.Get("/uuid/{uuid}", handlers.GetPersonByUuid)
		r.Get("/uuid/{uuid}/assets", handlers.GetPersonAssetsByUuid)
		r.Get("/githubname/{github}", handlers.GetPersonByGithubName)
	})
	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Post("/", handlers.CreateOrEditPerson)
		r.Delete("/{id}", handlers.DeletePerson)
	})
	return r
}
