package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func OrganizationRoutes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get("/", handlers.GetOrganizations)
		r.Get("/{uuid}", handlers.GetOrganizationByUuid)
		r.Get("/users/{uuid}", handlers.GetOrganizationUsers)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Post("/", handlers.CreateOrEditOrganization)
	})

	return r
}
