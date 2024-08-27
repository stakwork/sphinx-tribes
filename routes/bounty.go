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
		r.Get("/next/{created}", bountyHandler.GetNextBountyByCreated)
		r.Get("/previous/{created}", bountyHandler.GetPreviousBountyByCreated)
		r.Get("/org/next/{uuid}/{created}", bountyHandler.GetWorkspaceNextBountyByCreated)
		r.Get("/org/previous/{uuid}/{created}", bountyHandler.GetWorkspacePreviousBountyByCreated)
		r.Get("/workspace/next/{uuid}/{created}", bountyHandler.GetWorkspaceNextBountyByCreated)
		r.Get("/workspace/previous/{uuid}/{created}", bountyHandler.GetWorkspacePreviousBountyByCreated)

		r.Get("/created/{created}", bountyHandler.GetBountyByCreated)
		r.Get("/count/{personKey}/{tabType}", handlers.GetUserBountyCount)
		r.Get("/count", handlers.GetBountyCount)
		r.Get("/invoice/{paymentRequest}", bountyHandler.GetInvoiceData)
		r.Get("/filter/count", handlers.GetFilterCount)

	})
	r.Group(func(r chi.Router) {
		r.Post("/budget/withdraw", bountyHandler.BountyBudgetWithdraw)

		r.Use(auth.PubKeyContext)
		r.Post("/pay/{id}", bountyHandler.MakeBountyPayment)
		r.Post("/budget_workspace/withdraw", bountyHandler.NewBountyBudgetWithdraw)

		r.Post("/", bountyHandler.CreateOrEditBounty)
		r.Delete("/assignee", handlers.DeleteBountyAssignee)
		r.Delete("/{pubkey}/{created}", bountyHandler.DeleteBounty)
		r.Post("/paymentstatus/{created}", handlers.UpdatePaymentStatus)
		r.Post("/completedstatus/{created}", handlers.UpdateCompletedStatus)
	})
	return r
}
