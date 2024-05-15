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
		r.Get("/forworkspace/{uuid}", featureHandlers.GetFeaturesByWorkspaceUuid)
		r.Get("/{uuid}", featureHandlers.GetFeatureByUuid)
		r.Get("/workspace/count/{uuid}", featureHandlers.GetWorkspaceFeaturesCount)

		r.Post("/phase", featureHandlers.CreateOrEditFeaturePhase)
		r.Get("/{feature_uuid}/phase", featureHandlers.GetFeaturePhases)
		r.Get("/{feature_uuid}/phase/{phase_uuid}", featureHandlers.GetFeaturePhaseByUUID)
		r.Delete("/{feature_uuid}/phase/{phase_uuid}", featureHandlers.DeleteFeaturePhase)

	})
	return r
}
