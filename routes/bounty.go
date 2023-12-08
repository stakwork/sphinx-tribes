package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func BountyRoutes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Get("/all", handlers.GetAllBounties)
		r.Get("/id/{bountyId}", handlers.GetBountyById)
		r.Get("/created/{created}", handlers.GetBountyByCreated)
		r.Get("/count/{personKey}/{tabType}", handlers.GetUserBountyCount)
		r.Get("/count", handlers.GetBountyCount)
		r.Get("/invoice/{paymentRequest}", handlers.GetInvoiceData)

	})
	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)
		r.Post("/pay/{id}", handlers.MakeBountyPayment)
		r.Post("/budget/withdraw", handlers.BountyBudgetWithdraw)

		r.Post("/", handlers.CreateOrEditBounty)
		r.Delete("/assignee", handlers.DeleteBountyAssignee)
		r.Delete("/{pubkey}/{created}", handlers.DeleteBounty)
		r.Post("/paymentstatus/{created}", handlers.UpdatePaymentStatus)
	})
	return r
}
