package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func PeopleRoutes() chi.Router {
	r := chi.NewRouter()
	peopleHandler := handlers.NewPeopleHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Get("/", peopleHandler.GetListedPeople)
		r.Get("/search", peopleHandler.GetPeopleBySearch)
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
