package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func FeatureRoutes() chi.Router {
	r := chi.NewRouter()
	featureHandlers := handlers.NewFeatureHandler(db.DB)
	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Post("/", featureHandlers.CreateOrEditFeatures)
		r.Get("/forworkspace/{uuid}", featureHandlers.GetFeaturesByWorkspaceUuid)
		r.Get("/{uuid}", featureHandlers.GetFeatureByUuid)
		r.Post("/story", featureHandlers.CreateOrEditStory)
		r.Get("/{feature_uuid}/story", featureHandlers.GetStoriesByFeatureUuid)
		r.Get("/{feature_uuid}/story/{story_uuid}", featureHandlers.GetStoryByUuid)
		r.Delete("/{feature_uuid}/story/{story_uuid}", featureHandlers.DeleteStory)
	})
	return r
}
