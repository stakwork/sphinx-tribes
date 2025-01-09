package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/stakwork/sphinx-tribes/db"
)

type FeatureFlagResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func FeatureFlag(database db.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestPath := r.URL.Path

			endpoints, err := database.GetAllEndpoints()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			var found bool
			var endpoint db.Endpoint
			for _, e := range endpoints {
				if matchPath(e.Path, requestPath) {
					endpoint = e
					found = true
					break
				}
			}

			if !found {
				next.ServeHTTP(w, r)
				return
			}

			featureFlag, err := database.GetFeatureFlagByUUID(endpoint.FeatureFlagUUID)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if !featureFlag.Enabled {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(FeatureFlagResponse{
					Success: false,
					Message: "This feature is currently unavailable.",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func matchPath(pattern, requestPath string) bool {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	requestParts := strings.Split(strings.Trim(requestPath, "/"), "/")

	if len(patternParts) != len(requestParts) {
		return false
	}

	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			continue
		}
		if part != requestParts[i] {
			return false
		}
	}

	return true
}
