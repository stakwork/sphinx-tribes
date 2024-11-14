package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func WorkflowRoutes() chi.Router {
	r := chi.NewRouter()
	workflowHandlers := handlers.NewWorkFlowHandler(db.DB)

	r.Group(func(r chi.Router) {
		r.Post("/request", workflowHandlers.HandleWorkflowRequest)

	})
	return r
}
