package routes

import (
	"github.com/stakwork/sphinx-tribes/auth"
	"net/http"

	"github.com/go-chi/chi"
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
	})

	return r
}
