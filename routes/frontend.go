package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/frontend"
)

func IndexRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/ping", frontend.PingRoute)
	r.Group(func(r chi.Router) {
		r.Get("/", frontend.IndexRoute)
		r.Get("/static/*", frontend.StaticRoute)
		r.Get("/manifest.json", frontend.ManifestRoute)
		r.Get("/favicon.ico", frontend.FaviconRoute)
	})
	r.Group(func(r chi.Router) {
		r.Get("/t/static/*", frontend.StaticRoute)
		r.Get("/t/manifest.json", frontend.ManifestRoute)
		r.Get("/t/favicon.ico", frontend.FaviconRoute)
		r.Get("/t/{unique_name}", frontend.IndexRoute)
		r.Get("/t", frontend.IndexRoute)
	})
	r.Group(func(r chi.Router) {
		r.Get("/p/static/*", frontend.StaticRoute)
		r.Get("/p/manifest.json", frontend.ManifestRoute)
		r.Get("/p/favicon.ico", frontend.FaviconRoute)
		r.Get("/p/{pubkey}", frontend.IndexRoute)
		r.Get("/p", frontend.IndexRoute)
		r.Get("/p/{pubkey}/offer", frontend.IndexRoute)
		r.Get("/p/{pubkey}/badges", frontend.IndexRoute)
		r.Get("/p/{pubkey}/wanted", frontend.IndexRoute)
		r.Get("/p/{pubkey}/wanted/{page}/{index}", frontend.IndexRoute)
		r.Get("/p/{pubkey}/usertickets", frontend.IndexRoute)
		r.Get("/p/{pubkey}/usertickets/{ticket_id}/{index}", frontend.IndexRoute)
		r.Get("/p/{pubkey}/organizations", frontend.IndexRoute)
		r.Get("/b", frontend.IndexRoute)
		r.Get("/tickets", frontend.IndexRoute)
		r.Get("/bounties", frontend.IndexRoute)
		r.Get("/bounty/*", frontend.IndexRoute)
		r.Get("/leaderboard", frontend.IndexRoute)
		r.Get("/org/bounties/*", frontend.IndexRoute)
		r.Get("/admin", frontend.IndexRoute)
	})
	return r
}
