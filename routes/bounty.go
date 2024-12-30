package routes

import (
	"github.com/stakwork/sphinx-tribes/auth"
	"net/http"

	"github.com/go-chi/chi"
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
		r.Get("/filter/count", bountyHandler.GetFilterCount)

	})
	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Get("/bounty-cards", bountyHandler.GetBountyCards)
		r.Post("/budget/withdraw", bountyHandler.BountyBudgetWithdraw)
		r.Post("/pay/{id}", bountyHandler.MakeBountyPayment)
		r.Get("/payment/status/{id}", bountyHandler.GetBountyPaymentStatus)
		r.Get("/payment/{bountyId}", handlers.GetPaymentByBountyId)
		r.Put("/payment/status/{id}", bountyHandler.UpdateBountyPaymentStatus)

		r.Post("/{id}/proof", bountyHandler.AddProofOfWork)
		r.Get("/{id}/proofs", bountyHandler.GetProofsByBounty)
		r.Delete("/{id}/proofs/{proofId}", bountyHandler.DeleteProof)
		r.Patch("/{id}/proofs/{proofId}/status", bountyHandler.UpdateProofStatus)

		r.Post("/", bountyHandler.CreateOrEditBounty)
		r.Delete("/assignee", bountyHandler.DeleteBountyAssignee)
		r.Delete("/{pubkey}/{created}", bountyHandler.DeleteBounty)
		r.Post("/paymentstatus/{created}", handlers.UpdatePaymentStatus)
		r.Post("/completedstatus/{created}", handlers.UpdateCompletedStatus)
	})
	return r
}
