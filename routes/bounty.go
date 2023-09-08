package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func BountyRoutes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/", handlers.CreateOrEditBounty)
		r.Delete("/assignee", handlers.DeleteBountyAssignee)
		r.Delete("/{pubkey}/{created}", handlers.DeleteBounty)
		r.Get("/all", handlers.GetAllBounties)
		r.Get("/id/{bountyId}", handlers.GetBountyById)
		r.Get("/count/{personKey}/{tabType}", handlers.GetBountyCount)
		r.Post("/paymentstatus/{created}", handlers.UpdatePaymentStatus)
	})

	return r
}
