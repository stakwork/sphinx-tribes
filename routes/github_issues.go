package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func GithubIssuesRoutes() chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Get("/{owner}/{repo}/{issue}", handlers.GetGithubIssue)
		r.Get("/status/open", handlers.GetOpenGithubIssues)
	})
	return r
}
