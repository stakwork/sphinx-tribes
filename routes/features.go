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
		r.Post("/stories", featureHandlers.GetFeatureStories)

	})

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Post("/", featureHandlers.CreateOrEditFeatures)
		r.Post("/brief", featureHandlers.UpdateFeatureBrief)
		r.Get("/{uuid}", featureHandlers.GetFeatureByUuid)
		r.Put("/{uuid}/status", featureHandlers.UpdateFeatureStatus)
		r.Post("/brief/send", featureHandlers.BriefSend)
		// Old route for to getting features for workspace uuid
		r.Get("/forworkspace/{workspace_uuid}", featureHandlers.GetFeaturesByWorkspaceUuid)
		r.Get("/workspace/count/{uuid}", featureHandlers.GetWorkspaceFeaturesCount)
		r.Delete("/{uuid}", featureHandlers.DeleteFeature)

		r.Post("/phase", featureHandlers.CreateOrEditFeaturePhase)
		r.Get("/{feature_uuid}/phase", featureHandlers.GetFeaturePhases)
		r.Get("/{feature_uuid}/phase/{phase_uuid}", featureHandlers.GetFeaturePhaseByUUID)
		r.Delete("/{feature_uuid}/phase/{phase_uuid}", featureHandlers.DeleteFeaturePhase)

		r.Post("/story", featureHandlers.CreateOrEditStory)
		r.Post("/stories/send", featureHandlers.StoriesSend)
		r.Get("/{feature_uuid}/story", featureHandlers.GetStoriesByFeatureUuid)
		r.Get("/{feature_uuid}/story/{story_uuid}", featureHandlers.GetStoryByUuid)
		r.Delete("/{feature_uuid}/story/{story_uuid}", featureHandlers.DeleteStory)
		r.Get("/{feature_uuid}/phase/{phase_uuid}/bounty", featureHandlers.GetBountiesByFeatureAndPhaseUuid)
		r.Get("/{feature_uuid}/phase/{phase_uuid}/bounty/count", featureHandlers.GetBountiesCountByFeatureAndPhaseUuid)

	})
	return r
}
