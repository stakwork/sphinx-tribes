package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func PeopleRoutes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Get("/", handlers.GetListedPeople)
		r.Get("/search", handlers.GetPeopleBySearch)
		r.Get("/posts", handlers.GetListedPosts)
		r.Get("/wanteds/assigned/{uuid}", handlers.GetPersonAssignedBounties)
		r.Get("/wanteds/created/{uuid}", handlers.GetPersonCreatedBounties)
		r.Get("/wanteds/header", handlers.GetWantedsHeader)
		r.Get("/short", handlers.GetPeopleShortList)
		r.Get("/offers", handlers.GetListedOffers)
		r.Get("/bounty/leaderboard", handlers.GetBountiesLeaderboard)
	})
	return r
}
