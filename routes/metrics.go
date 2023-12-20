package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func MetricsRoutes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContextSuperAdmin)

		r.Post("/payment", handlers.PaymentMetrics)
		r.Post("/people", handlers.PeopleMetrics)
		r.Post("/organization", handlers.OrganizationtMetrics)
		r.Post("/bounty_stats", handlers.BountyMetrics)
		r.Post("/bounties", handlers.MetricsBounties)
		r.Post("/bounties/count", handlers.MetricsBountiesCount)
	})
	return r
}
