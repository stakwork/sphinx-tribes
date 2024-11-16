package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestWorkflowRoutes(t *testing.T) {
	r := chi.NewRouter()
	r.Mount("/workflows", WorkflowRoutes())

	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "workflow request endpoint",
			method:         "POST",
			path:           "/workflows/request",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "workflow response endpoint",
			method:         "POST",
			path:           "/workflows/response",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusNotFound, w.Code, "Route should exist")
		})
	}
}
