package routes

import (
	"net/http"

	"github.com/stakwork/sphinx-tribes/auth"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func BountyRoutes() chi.Router {
	r := chi.NewRouter()
	bountyHandler := handlers.NewBountyHandler(http.DefaultClient, db.DB)
	tribeHandlers := handlers.NewTribeHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Get("/all", bountyHandler.GetAllBounties)
		r.Get("/featured/all", bountyHandler.GetAllFeaturedBounties)

		r.Get("/code/{code}", bountyHandler.GetBountyByCode)
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
		r.Get("/workspace/timerange/{workspaceId}/{daysStart}/{daysEnd}", bountyHandler.GetBountiesByWorkspaceTime)

		r.Get("/stakes", bountyHandler.GetAllBountyStakes)
		r.Get("/stake/bounty/{bountyId}", bountyHandler.GetBountyStakesByBountyID)
		r.Get("/stake/{id}", bountyHandler.GetBountyStakeByID)
		r.Get("/stake/hunter/{hunterPubKey}", bountyHandler.GetBountyStakesByHunterPubKey)
		r.Post("/process_stake/{bountyId}", tribeHandlers.ProcessStake)
	})
	r.Group(func(r chi.Router) {
		r.Use(auth.CombinedAuthContext)
		r.Post("/featured/create", bountyHandler.CreateFeaturedBounty)
		r.Put("/featured/update", bountyHandler.UpdateFeaturedBounty)
		r.Delete("/featured/delete/{bountyId}", bountyHandler.DeleteFeaturedBounty)

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

		r.Get("/{id}/timing", bountyHandler.GetBountyTimingStats)
		r.Put("/{id}/timing/start", bountyHandler.StartBountyTiming)
		r.Put("/{id}/timing/close", bountyHandler.CloseBountyTiming)
		r.Delete("/{id}/timing", bountyHandler.DeleteBountyTiming)

		r.Post("/stake", bountyHandler.CreateBountyStake)
		r.Put("/stake/{id}", bountyHandler.UpdateBountyStake)
		r.Delete("/stake/{id}", bountyHandler.DeleteBountyStake)

		r.Post("/stake/stakeprocessing", bountyHandler.CreateBountyStakeProcess)
		r.Get("/stake/stakeprocessing", bountyHandler.GetAllBountyStakeProcesses)
		r.Get("/stake/stakeprocessing/{id}", bountyHandler.GetBountyStakeProcessByID)
		r.Get("/stake/stakeprocessing/bounty/{bountyId}", bountyHandler.GetBountyStakeProcessesByBountyID)
		r.Put("/stake/stakeprocessing/{id}", bountyHandler.UpdateBountyStakeProcess)
		r.Delete("/stake/stakeprocessing/{id}", bountyHandler.DeleteBountyStakeProcess)
	})
	return r
}
