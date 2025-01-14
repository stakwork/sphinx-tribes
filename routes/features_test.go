package routes

import (
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FeatureMockHandler(t *testing.T, expectedStatus int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/1234" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/1234/story" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == "/stories" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodPut && r.URL.Path == "/1234/status" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodDelete && r.URL.Path == "/1234" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == "/brief/send" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/workspace/count/1234" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

func TestFeatureRoutes(t *testing.T) {
	r := chi.NewRouter()

	r.Post("/stories", FeatureMockHandler(t, http.StatusOK))
	r.Post("/", FeatureMockHandler(t, http.StatusOK))
	r.Get("/{uuid}", FeatureMockHandler(t, http.StatusOK))
	r.Put("/{uuid}/status", FeatureMockHandler(t, http.StatusOK))
	r.Delete("/{uuid}", FeatureMockHandler(t, http.StatusOK))
	r.Post("/brief/send", FeatureMockHandler(t, http.StatusOK))
	r.Get("/workspace/count/{uuid}", FeatureMockHandler(t, http.StatusOK))
	r.Post("/phase", FeatureMockHandler(t, http.StatusOK))
	r.Get("/{feature_uuid}/story", FeatureMockHandler(t, http.StatusOK))

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Test POST /stories Route",
			method:         http.MethodPost,
			path:           "/stories",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test POST / Route",
			method:         http.MethodPost,
			path:           "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test GET /{uuid} Route",
			method:         http.MethodGet,
			path:           "/1234",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test PUT /{uuid}/status Route",
			method:         http.MethodPut,
			path:           "/1234/status",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test DELETE /{uuid} Route",
			method:         http.MethodDelete,
			path:           "/1234",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test POST /brief/send Route",
			method:         http.MethodPost,
			path:           "/brief/send",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test GET /workspace/count/{uuid} Routes",
			method:         http.MethodGet,
			path:           "/workspace/count/1234",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test GET /{feature_uuid}/story Route",
			method:         http.MethodGet,
			path:           "/1234/story",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code, "Handler returned wrong status code")
		})
	}
}
