package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
)

func TestHandleWorkflowRequest(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	wh := NewWorkFlowHandler(db.TestDB)

	t.Run("successful workflow request", func(t *testing.T) {
		request := db.WfRequest{
			WorkflowID: uuid.New().String(),
			Source:     "test_source",
			RequestID:  uuid.New().String(),
		}
		body, _ := json.Marshal(request)

		req := httptest.NewRequest(http.MethodPost, "/workflows/request", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		wh.HandleWorkflowRequest(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var respBody map[string]string
		err := json.NewDecoder(w.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Equal(t, "success", respBody["status"])
		assert.NotEmpty(t, respBody["request_id"])
	})

	t.Run("invalid JSON format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/workflows/request", bytes.NewBuffer([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		wh.HandleWorkflowRequest(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request format")
	})

	t.Run("missing required fields", func(t *testing.T) {
		request := db.WfRequest{
			Source: "test_source",
		}
		body, _ := json.Marshal(request)

		req := httptest.NewRequest(http.MethodPost, "/workflows/request", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		wh.HandleWorkflowRequest(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Missing required fields: workflow_id or source")
	})

}
