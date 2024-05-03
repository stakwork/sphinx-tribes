package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func MetricsRoutes() chi.Router {
	r := chi.NewRouter()
	mh := handlers.NewMetricHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContextSuperAdmin)

		r.Get("/workspaces", handlers.GetAdminWorkspaces)

		r.Post("/payment", handlers.PaymentMetrics)
		r.Post("/people", handlers.PeopleMetrics)
		r.Post("/organization", handlers.WorkspacetMetrics)
		r.Post("/bounty_stats", mh.BountyMetrics)
		r.Post("/bounties", mh.MetricsBounties)
		r.Post("/bounties/count", mh.MetricsBountiesCount)
		r.Post("/bounties/providers", mh.MetricsBountiesProviders)
		r.Post("/csv", handlers.MetricsCsv)
	})
	return r
}
