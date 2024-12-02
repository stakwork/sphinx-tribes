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

	r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Group(func(r chi.Router) {
		r.Get("/{uuid}", ticketHandler.GetTicket)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Get("/feature/{feature_uuid}/phase/{phase_uuid}", ticketHandler.GetTicketsByPhaseUUID)
		r.Post("/review/send", ticketHandler.PostTicketDataToStakwork)
		r.Post("/review", ticketHandler.ProcessTicketReview)

		r.Post("/{uuid}", ticketHandler.UpdateTicket)
		r.Delete("/{uuid}", ticketHandler.DeleteTicket)
	})

	return r
}
