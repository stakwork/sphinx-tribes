package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func CodeSpaceRoutes() chi.Router {
	r := chi.NewRouter()
	codeSpaceHandler := handlers.NewCodeSpaceHandler(db.DB)

	r.Group(func(r chi.Router) {
		r.Use(auth.CombinedAuthContext)

		r.Get("/workspaces", codeSpaceHandler.GetAllCodeSpaceMaps)
		r.Get("/workspaces/workspace/{workspaceID}", codeSpaceHandler.GetCodeSpaceMapsByWorkspace)
		r.Get("/workspaces/user/{userPubkey}", codeSpaceHandler.GetCodeSpaceMapsByUser)
		r.Get("/workspaces/codespace", codeSpaceHandler.GetCodeSpaceMapsByURL)
		r.Get("/workspaces/query", codeSpaceHandler.QueryCodeSpaceMaps)
		r.Post("/workspaces", codeSpaceHandler.CreateCodeSpaceMap)
		r.Put("/workspaces/{id}", codeSpaceHandler.UpdateCodeSpaceMap)
		r.Delete("/workspaces/{id}", codeSpaceHandler.DeleteCodeSpaceMap)
	})

	return r
} 