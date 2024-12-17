package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
	"github.com/stakwork/sphinx-tribes/utils"
)

func PeopleRoutes() chi.Router {
	r := chi.NewRouter()
	bountyHandler := handlers.NewBountyHandler(http.DefaultClient, db.DB)

	peopleHandler := handlers.NewPeopleHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Get("/", utils.TraceWithLogging(peopleHandler.GetListedPeople))
		r.Get("/search", utils.TraceWithLogging(peopleHandler.GetPeopleBySearch))
		r.Get("/posts", utils.TraceWithLogging(handlers.GetListedPosts))
		r.Get("/wanteds/assigned/{uuid}", bountyHandler.GetPersonAssignedBounties)
		r.Get("/wanteds/created/{uuid}", bountyHandler.GetPersonCreatedBounties)
		r.Get("/wanteds/header", handlers.GetWantedsHeader)
		r.Get("/short", handlers.GetPeopleShortList)
		r.Get("/offers", handlers.GetListedOffers)
		r.Get("/bounty/leaderboard", handlers.GetBountiesLeaderboard)
	})
	return r
}
