package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func BountyRoutes() chi.Router {
	r := chi.NewRouter()
	bountyHandler := handlers.NewBountyHandler(http.DefaultClient, db.DB)
	r.Group(func(r chi.Router) {
		r.Get("/all", bountyHandler.GetAllBounties)
		r.Get("/id/{bountyId}", bountyHandler.GetBountyById)
		r.Get("/index/{bountyId}", bountyHandler.GetBountyIndexById)
		r.Get("/next/{created}", handlers.GetNextBountyByCreated)
		r.Get("/previous/{created}", handlers.GetPreviousBountyByCreated)
		r.Get("/org/next/{uuid}/{created}", handlers.GetOrganizationNextBountyByCreated)
		r.Get("/org/previous/{uuid}/{created}", handlers.GetOrganizationPreviousBountyByCreated)
		r.Get("/created/{created}", bountyHandler.GetBountyByCreated)
		r.Get("/count/{personKey}/{tabType}", handlers.GetUserBountyCount)
		r.Get("/count", handlers.GetBountyCount)
		r.Get("/invoice/{paymentRequest}", handlers.GetInvoiceData)
		r.Get("/filter/count", handlers.GetFilterCount)

	})
	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)
		r.Post("/pay/{id}", handlers.MakeBountyPayment)
		r.Post("/budget/withdraw", bountyHandler.BountyBudgetWithdraw)

		r.Post("/", bountyHandler.CreateOrEditBounty)
		r.Delete("/assignee", handlers.DeleteBountyAssignee)
		r.Delete("/{pubkey}/{created}", bountyHandler.DeleteBounty)
		r.Post("/paymentstatus/{created}", handlers.UpdatePaymentStatus)
	})
	return r
}
