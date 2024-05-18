package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func FeatureRoutes() chi.Router {
	r := chi.NewRouter()
	featureHandlers := handlers.NewFeatureHandler(&db.DB)
	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Post("/", featureHandlers.CreateOrEditFeatures)
		r.Get("/{uuid}", featureHandlers.GetFeatureByUuid)
		r.Get("/{workspace_uuid}/workspaces", featureHandlers.GetFeaturesByWorkspaceUuid)
		r.Get("/workspace/count/{uuid}", featureHandlers.GetWorkspaceFeaturesCount)
		r.Delete("/{uuid}", featureHandlers.DeleteFeature)

		r.Post("/phase", featureHandlers.CreateOrEditFeaturePhase)
		r.Get("/{feature_uuid}/phase", featureHandlers.GetFeaturePhases)
		r.Get("/{feature_uuid}/phase/{phase_uuid}", featureHandlers.GetFeaturePhaseByUUID)
		r.Delete("/{feature_uuid}/phase/{phase_uuid}", featureHandlers.DeleteFeaturePhase)

		r.Post("/story", featureHandlers.CreateOrEditStory)
		r.Get("/{feature_uuid}/story", featureHandlers.GetStoriesByFeatureUuid)
		r.Get("/{feature_uuid}/story/{story_uuid}", featureHandlers.GetStoryByUuid)
		r.Delete("/{feature_uuid}/story/{story_uuid}", featureHandlers.DeleteStory)
	})
	return r
}
