package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func FeatureFlagMockHandler(t *testing.T, expectedStatus int, validateReq func(*http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isProtectedFeatureFlagRoute(r.URL.Path, r.Method) {
			token := r.Header.Get("Authorization")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		if validateReq != nil && !validateReq(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if strings.Contains(r.URL.Path, "invalid") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if strings.Contains(r.URL.Path, "notfound") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(expectedStatus)
	}
}

func TestFeatureFlagRoutes(t *testing.T) {
	r := chi.NewRouter()
	featureFlagRouter := chi.NewRouter()

	featureFlagRouter.Get("/", FeatureFlagMockHandler(t, http.StatusOK, nil))
	featureFlagRouter.Post("/", FeatureFlagMockHandler(t, http.StatusCreated, validateFeatureFlagCreate))
	featureFlagRouter.Put("/{id}", FeatureFlagMockHandler(t, http.StatusOK, validateFeatureFlagUpdate))
	featureFlagRouter.Delete("/{id}", FeatureFlagMockHandler(t, http.StatusNoContent, nil))

	featureFlagRouter.Post("/{feature_flag_id}/endpoints", FeatureFlagMockHandler(t, http.StatusCreated, validateEndpointCreate))
	featureFlagRouter.Put("/{feature_flag_id}/endpoints/{endpoint_id}", FeatureFlagMockHandler(t, http.StatusOK, validateEndpointUpdate))
	featureFlagRouter.Delete("/{feature_flag_id}/endpoints/{endpoint_id}", FeatureFlagMockHandler(t, http.StatusNoContent, nil))

	r.Mount("/feature-flags", featureFlagRouter)

	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		token          string
		expectedStatus int
	}{
		{"Create Feature Flag", "POST", "/feature-flags/", map[string]interface{}{"name": "test", "enabled": true}, "valid-token", http.StatusCreated},
		{"Get All Feature Flags", "GET", "/feature-flags/", nil, "", http.StatusOK},
		{"Create Feature Flag - Invalid", "POST", "/feature-flags/", map[string]interface{}{}, "valid-token", http.StatusBadRequest},
		{"Update Feature Flag", "PUT", "/feature-flags/1", map[string]interface{}{"name": "updated"}, "valid-token", http.StatusOK},
		{"Delete Feature Flag", "DELETE", "/feature-flags/1", nil, "valid-token", http.StatusNoContent},
		{"Add Endpoint to Feature Flag", "POST", "/feature-flags/1/endpoints", map[string]interface{}{"url": "/test"}, "valid-token", http.StatusCreated},
		{"Update Feature Flag Endpoint", "PUT", "/feature-flags/1/endpoints/1", map[string]interface{}{"url": "/updated"}, "valid-token", http.StatusOK},
		{"Delete Feature Flag Endpoint", "DELETE", "/feature-flags/1/endpoints/1", nil, "valid-token", http.StatusNoContent},
		{"Create Feature Flag - Invalid Data", "POST", "/feature-flags/", map[string]interface{}{}, "valid-token", http.StatusBadRequest},
		{"Update Non-Existent Feature Flag", "PUT", "/feature-flags/notfound", map[string]interface{}{"name": "updated"}, "valid-token", http.StatusNotFound},
		{"Delete Non-Existent Feature Flag", "DELETE", "/feature-flags/notfound", nil, "valid-token", http.StatusNotFound},
		{"Add Endpoint to Non-Existent Feature Flag", "POST", "/feature-flags/notfound/endpoints", map[string]interface{}{"url": "/test"}, "valid-token", http.StatusNotFound},
		{"Update Non-Existent Feature Flag Endpoint", "PUT", "/feature-flags/1/endpoints/notfound", map[string]interface{}{"url": "/updated"}, "valid-token", http.StatusNotFound},
		{"Delete Non-Existent Feature Flag Endpoint", "DELETE", "/feature-flags/1/endpoints/notfound", nil, "valid-token", http.StatusNotFound},
		{"Invalid Feature Flag ID Format", "PUT", "/feature-flags/invalid_id", map[string]interface{}{"name": "updated"}, "valid-token", http.StatusBadRequest},
		{"Invalid HTTP Method", "PATCH", "/feature-flags/", nil, "valid-token", http.StatusMethodNotAllowed},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.body != nil {
				bodyBytes, _ := json.Marshal(tc.body)
				req = httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tc.method, tc.path, nil)
			}

			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code, "Handler returned wrong status code for test: "+tc.name)
		})
	}
}

func isProtectedFeatureFlagRoute(path string, method string) bool {
	publicRoutes := map[string][]string{
		"GET": {"/feature-flags/"},
	}

	for publicMethod, publicPaths := range publicRoutes {
		if method == publicMethod {
			for _, publicPath := range publicPaths {
				if strings.HasPrefix(path, publicPath) {
					return false
				}
			}
		}
	}

	return true
}

func validateFeatureFlagCreate(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasName := body["name"]
	_, hasEnabled := body["enabled"]
	return hasName && hasEnabled
}

func validateFeatureFlagUpdate(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasName := body["name"]
	return hasName
}

func validateEndpointCreate(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasURL := body["url"]
	return hasURL
}

func validateEndpointUpdate(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasURL := body["url"]
	return hasURL
}
