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

	req := httptest.NewRequest("POST", "/workflows/request", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusNotFound, w.Code)
}
