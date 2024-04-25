package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func WorkspaceRoutes() chi.Router {
	r := chi.NewRouter()
	organizationHandlers := handlers.NewWorkspaceHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Get("/", handlers.GetWorkspaces)
		r.Get("/count", handlers.GetWorkspacesCount)
		r.Get("/{uuid}", handlers.GetWorkspaceByUuid)
		r.Get("/users/{uuid}", handlers.GetWorkspaceUsers)
		r.Get("/users/{uuid}/count", handlers.GetWorkspaceUsersCount)
		r.Get("/bounties/{uuid}", organizationHandlers.GetWorkspaceBounties)
		r.Get("/bounties/{uuid}/count", organizationHandlers.GetWorkspaceBountiesCount)
		r.Get("/user/{userId}", handlers.GetUserWorkspaces)
		r.Get("/user/dropdown/{userId}", organizationHandlers.GetUserDropdownWorkspaces)
	})
	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Post("/", organizationHandlers.CreateOrEditWorkspace)
		r.Post("/users/{uuid}", handlers.CreateWorkspaceUser)
		r.Delete("/users/{uuid}", handlers.DeleteWorkspaceUser)
		r.Post("/users/role/{uuid}/{user}", handlers.AddUserRoles)

		r.Get("/foruser/{uuid}", handlers.GetWorkspaceUser)
		r.Get("/bounty/roles", handlers.GetBountyRoles)
		r.Get("/users/role/{uuid}/{user}", handlers.GetUserRoles)
		r.Get("/budget/{uuid}", organizationHandlers.GetWorkspaceBudget)
		r.Get("/budget/history/{uuid}", organizationHandlers.GetWorkspaceBudgetHistory)
		r.Get("/payments/{uuid}", handlers.GetPaymentHistory)
		r.Get("/poll/invoices/{uuid}", organizationHandlers.PollBudgetInvoices)
		r.Get("/invoices/count/{uuid}", handlers.GetInvoicesCount)
		r.Delete("/delete/{uuid}", organizationHandlers.DeleteWorkspace)

		r.Post("/mission", organizationHandlers.UpdateWorkspaceMission)
		r.Post("/tactics", organizationHandlers.UpdateWorkspaceTactics)
		r.Post("/schematicurl", organizationHandlers.UpdateWorkspaceSchematicUrl)
	})
	return r
}
