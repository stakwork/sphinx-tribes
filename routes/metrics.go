package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func MetricsRoutes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Post("/payment_metrics", handlers.PaymentMetrics)
		r.Post("/people_created", handlers.PeopleMetrics)
		r.Post("/organization_created", handlers.OrganizationtMetrics)
		r.Post("/bounty_metrics", handlers.BountyMetrics)
	})
	return r
}
