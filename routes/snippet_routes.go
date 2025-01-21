package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func SnippetRoutes() chi.Router {
	r := chi.NewRouter()
	snippetHandler := handlers.NewSnippetHandler(http.DefaultClient, db.DB)

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Post("/create", snippetHandler.CreateSnippet)
		r.Get("/workspace/{workspace_uuid}", snippetHandler.GetSnippetsByWorkspace)
		r.Get("/{id}", snippetHandler.GetSnippetByID)
		r.Put("/{id}", snippetHandler.UpdateSnippet)
		r.Delete("/{id}", snippetHandler.DeleteSnippet)
	})

	return r
}
