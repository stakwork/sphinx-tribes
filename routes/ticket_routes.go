package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func TicketRoutes() chi.Router {
	r := chi.NewRouter()
	ticketHandler := handlers.NewTicketHandler(http.DefaultClient, db.DB)

	r.Group(func(r chi.Router) {
		r.Get("/{uuid}", ticketHandler.GetTicket)
		r.Post("/review", ticketHandler.ProcessTicketReview)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Get("/feature/{feature_uuid}/phase/{phase_uuid}", ticketHandler.GetTicketsByPhaseUUID)
		r.Post("/review/send", ticketHandler.PostTicketDataToStakwork)
		r.Post("/{uuid}", ticketHandler.UpdateTicket)
		r.Post("/{ticket_group}/sequence", ticketHandler.UpdateTicketSequence)
		r.Post("/{ticket_uuid}/bounty", ticketHandler.TicketToBounty)
		r.Delete("/{uuid}", ticketHandler.DeleteTicket)
		r.Get("/group/{group_uuid}", ticketHandler.GetTicketsByGroup)

		r.Post("/workspace/{workspace_uuid}/draft", ticketHandler.CreateWorkspaceDraftTicket)
		r.Get("/workspace/{workspace_uuid}/draft/{uuid}", ticketHandler.GetWorkspaceDraftTicket)
		r.Post("/workspace/{workspace_uuid}/draft/{uuid}", ticketHandler.UpdateWorkspaceDraftTicket)
		r.Delete("/workspace/{workspace_uuid}/draft/{uuid}", ticketHandler.DeleteWorkspaceDraftTicket)

		r.Post("/plan", ticketHandler.CreateTicketPlan)
		r.Get("/plan/{uuid}", ticketHandler.GetTicketPlan)
		r.Delete("/plan/{uuid}", ticketHandler.DeleteTicketPlan)
		r.Get("/plan/feature/{feature_uuid}", ticketHandler.GetTicketPlansByFeature)
		r.Get("/plan/phase/{phase_uuid}", ticketHandler.GetTicketPlansByPhase)
		r.Get("/plan/workspace/{workspace_uuid}", ticketHandler.GetTicketPlansByWorkspace)
	})

	return r
}
