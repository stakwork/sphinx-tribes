package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func WorkspaceRoutes() chi.Router {
	r := chi.NewRouter()
	workspaceHandlers := handlers.NewWorkspaceHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Get("/", handlers.GetWorkspaces)
		r.Get("/count", handlers.GetWorkspacesCount)
		r.Get("/{uuid}", handlers.GetWorkspaceByUuid)
		r.Get("/users/{uuid}", handlers.GetWorkspaceUsers)
		r.Get("/users/{uuid}/count", handlers.GetWorkspaceUsersCount)
		r.Get("/bounties/{uuid}", workspaceHandlers.GetWorkspaceBounties)
		r.Get("/bounties/{uuid}/count", workspaceHandlers.GetWorkspaceBountiesCount)
		r.Get("/user/{userId}", handlers.GetUserWorkspaces)
		r.Get("/user/dropdown/{userId}", workspaceHandlers.GetUserDropdownWorkspaces)
	})
	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Post("/", workspaceHandlers.CreateOrEditWorkspace)
		r.Post("/users/{uuid}", handlers.CreateWorkspaceUser)
		r.Delete("/users/{uuid}", handlers.DeleteWorkspaceUser)
		r.Post("/users/role/{uuid}/{user}", handlers.AddUserRoles)

		r.Get("/foruser/{uuid}", handlers.GetWorkspaceUser)
		r.Get("/bounty/roles", handlers.GetBountyRoles)
		r.Get("/users/role/{uuid}/{user}", handlers.GetUserRoles)
		r.Get("/budget/{uuid}", workspaceHandlers.GetWorkspaceBudget)
		r.Get("/budget/history/{uuid}", workspaceHandlers.GetWorkspaceBudgetHistory)
		r.Get("/payments/{uuid}", handlers.GetPaymentHistory)
		r.Get("/poll/invoices/{uuid}", workspaceHandlers.PollBudgetInvoices)
		r.Get("/invoices/count/{uuid}", handlers.GetInvoicesCount)
		r.Delete("/delete/{uuid}", workspaceHandlers.DeleteWorkspace)

		r.Post("/mission", workspaceHandlers.UpdateWorkspace)
		r.Post("/tactics", workspaceHandlers.UpdateWorkspace)
		r.Post("/schematicurl", workspaceHandlers.UpdateWorkspace)
<<<<<<< HEAD

		r.Post("/repositories", workspaceHandlers.CreateWorkspaceRepository)
		r.Get("/repositories/{uuid}", workspaceHandlers.GetWorkspaceRepositorByWorkspaceUuid)
=======
		r.Get("/{workspace_uuid}/features", workspaceHandlers.GetFeaturesByWorkspaceUuid)
>>>>>>> e1f721d5 (Modify Features endpoint and add delete feature)
	})
	return r
}
