package routes

import (
	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers"
)

func FeatureFlagRoutes() chi.Router {
	r := chi.NewRouter()
	featureFlagHandler := handlers.NewFeatureFlagHandler(db.DB)

	r.Group(func(r chi.Router) {
		r.Use(auth.PubKeyContext)

		r.Get("/", featureFlagHandler.GetFeatureFlags)
		r.Post("/", featureFlagHandler.CreateFeatureFlag)
		r.Put("/{id}", featureFlagHandler.UpdateFeatureFlag)
		r.Delete("/{id}", featureFlagHandler.DeleteFeatureFlag)

		r.Post("/{feature_flag_id}/endpoints", featureFlagHandler.AddFeatureFlagEndpoint)
		r.Put("/{feature_flag_id}/endpoints/{endpoint_id}", featureFlagHandler.UpdateFeatureFlagEndpoint)
		r.Delete("/{feature_flag_id}/endpoints/{endpoint_id}", featureFlagHandler.DeleteFeatureFlagEndpoint)
	})

	return r
}
