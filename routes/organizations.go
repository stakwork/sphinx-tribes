package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func OrganizationRoutes() chi.Router {
	r := chi.NewRouter()
	organizationHandlers := handlers.NewOrganizationHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Get("/", handlers.GetOrganizations)
		r.Get("/count", handlers.GetOrganizationsCount)
		r.Get("/{uuid}", handlers.GetOrganizationByUuid)
		r.Get("/users/{uuid}", handlers.GetOrganizationUsers)
		r.Get("/users/{uuid}/count", handlers.GetOrganizationUsersCount)
		r.Get("/bounties/{uuid}", organizationHandlers.GetOrganizationBounties)
		r.Get("/bounties/{uuid}/count", organizationHandlers.GetOrganizationBountiesCount)
		r.Get("/user/{userId}", handlers.GetUserOrganizations)
		r.Get("/user/dropdown/{userId}", organizationHandlers.GetUserDropdownOrganizations)
	})
	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Post("/", organizationHandlers.CreateOrEditOrganization)
		r.Post("/users/{uuid}", handlers.CreateOrganizationUser)
		r.Delete("/users/{uuid}", handlers.DeleteOrganizationUser)
		r.Post("/users/role/{uuid}/{user}", handlers.AddUserRoles)

		r.Get("/foruser/{uuid}", handlers.GetOrganizationUser)
		r.Get("/bounty/roles", handlers.GetBountyRoles)
		r.Get("/users/role/{uuid}/{user}", handlers.GetUserRoles)
		r.Get("/budget/{uuid}", organizationHandlers.GetOrganizationBudget)
		r.Get("/budget/history/{uuid}", organizationHandlers.GetOrganizationBudgetHistory)
		r.Get("/payments/{uuid}", handlers.GetPaymentHistory)
		r.Get("/poll/invoices/{uuid}", handlers.PollBudgetInvoices)
		r.Get("/invoices/count/{uuid}", handlers.GetInvoicesCount)
		r.Delete("/delete/{uuid}", organizationHandlers.DeleteOrganization)
	})
	return r
}
