package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func BountyRoutes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/", handlers.CreateOrEditBounty)
		r.Delete("/assignee", handlers.DeleteBountyAssignee)
		r.Get("/all", handlers.GetAllBounties)
		r.Get("/id/{bountyId}", handlers.GetBountyById)
		r.Get("/count/{personKey}/{tabType}", handlers.GetBountyCount)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Delete("/{pubKey}/{created}", handlers.DeleteBounty)
	})

	return r
}
