package handlers

import (
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
)

func TestHandleWorkflowRequest(t *testing.T) {

	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	wh := NewWorkFlowHandler(db.TestDB)

	t.Run("successful workflow request", func(t *testing.T) {

		request := &db.WfRequest{
			Source:     "test_source_1",
			Action:     "test_action_1",
			WorkflowID: "test_workflow",
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
func TestHandleWorkflowResponse(t *testing.T) {

	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	wh := NewWorkFlowHandler(db.TestDB)
	t.Run("should process workflow response successfully", func(t *testing.T) {

		requestData := db.PropertyMap{
			"test_key": "test_value",
		}

		// Create a workflow request with unique source/action
		workflowRequest := &db.WfRequest{
			RequestID:   uuid.New().String(),
			Source:      "test_source_1",
			Action:      "test_action_1",
			Status:      db.StatusNew,
			WorkflowID:  "test_workflow",
			RequestData: requestData,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := db.TestDB.CreateWorkflowRequest(workflowRequest)
		assert.NoError(t, err)

		responseData := db.PropertyMap{
			"result": "success",
		}

		response := struct {
			RequestID    string             `json:"request_id"`
			Status       db.WfRequestStatus `json:"status"`
			ResponseData db.PropertyMap     `json:"response_data"`
		}{
			RequestID:    workflowRequest.RequestID,
			Status:       db.StatusCompleted,
			ResponseData: responseData,
		}

		payload, err := json.Marshal(response)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/workflows/response", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		wh.HandleWorkflowResponse(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var respBody map[string]string
		err = json.NewDecoder(w.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Equal(t, "success", respBody["status"])
		assert.Equal(t, workflowRequest.RequestID, respBody["request_id"])

		updatedReq, err := db.TestDB.GetWorkflowRequest(workflowRequest.RequestID)
		assert.NoError(t, err)
		assert.NotNil(t, updatedReq)
		assert.Equal(t, db.StatusCompleted, updatedReq.Status)
		assert.Equal(t, responseData, updatedReq.ResponseData)
	})

	t.Run("should process with processing map when required", func(t *testing.T) {

		requestData := db.PropertyMap{
			"test_key1": "test_value1",
		}

		// Create a workflow request
		workflowRequest := &db.WfRequest{
			RequestID:   uuid.New().String(),
			Source:      "test_source1",
			Action:      "test_action1",
			Status:      db.StatusNew,
			WorkflowID:  "test_workflow1",
			RequestData: requestData,
		}
		err := db.TestDB.CreateWorkflowRequest(workflowRequest)
		assert.NoError(t, err)

		// Create a processing map
		processingMap := &db.WfProcessingMap{
			Type:               workflowRequest.Source,
			ProcessKey:         workflowRequest.Action,
			RequiresProcessing: true,
			HandlerFunc:        "test_handler",
		}
		err = db.TestDB.CreateProcessingMap(processingMap)
		assert.NoError(t, err)

		response := struct {
			RequestID    string             `json:"request_id"`
			Status       db.WfRequestStatus `json:"status"`
			ResponseData db.PropertyMap     `json:"response_data"`
		}{
			RequestID:    workflowRequest.RequestID,
			Status:       db.StatusCompleted,
			ResponseData: db.PropertyMap{"result": "success"},
		}
		payload, _ := json.Marshal(response)

		req := httptest.NewRequest(http.MethodPost, "/workflows/response", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		wh.HandleWorkflowResponse(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify request was set to pending
		updatedReq, err := db.TestDB.GetWorkflowRequest(workflowRequest.RequestID)
		assert.NoError(t, err)
		assert.Equal(t, db.StatusPending, updatedReq.Status)
		assert.Equal(t, response.ResponseData, updatedReq.ResponseData)
	})

	t.Run("should return 400 error if request_id is missing", func(t *testing.T) {
		response := map[string]interface{}{
			"status":        string(db.StatusCompleted),
			"response_data": map[string]interface{}{"result": "success"},
		}
		payload, _ := json.Marshal(response)

		req := httptest.NewRequest(http.MethodPost, "/workflows/response", bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		wh.HandleWorkflowResponse(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 400 error if JSON format is invalid", func(t *testing.T) {
		invalidJSON := []byte(`{"request_id": "123", status: invalid}`)

		req := httptest.NewRequest(http.MethodPost, "/workflows/response", bytes.NewBuffer(invalidJSON))
		w := httptest.NewRecorder()

		wh.HandleWorkflowResponse(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 404 error if workflow request is not found", func(t *testing.T) {
		response := map[string]interface{}{
			"request_id":    uuid.New().String(),
			"status":        string(db.StatusCompleted),
			"response_data": map[string]interface{}{"result": "success"},
		}
		payload, _ := json.Marshal(response)

		req := httptest.NewRequest(http.MethodPost, "/workflows/response", bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		wh.HandleWorkflowResponse(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should handle empty response data", func(t *testing.T) {

		requestData := db.PropertyMap{
			"test_key": "test_value",
		}

		workflowRequest := &db.WfRequest{
			RequestID:   uuid.New().String(),
			Source:      "test_source",
			Action:      "test_action",
			Status:      db.StatusNew,
			WorkflowID:  "test_workflow",
			RequestData: requestData,
		}
		err := db.TestDB.CreateWorkflowRequest(workflowRequest)
		assert.NoError(t, err)

		response := map[string]interface{}{
			"request_id": workflowRequest.RequestID,
			"status":     string(db.StatusCompleted),
		}
		payload, _ := json.Marshal(response)

		req := httptest.NewRequest(http.MethodPost, "/workflows/response", bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		wh.HandleWorkflowResponse(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should maintain NEW status when processing map exists", func(t *testing.T) {
		requestData := db.PropertyMap{
			"test_key": "test_value",
		}

		workflowRequest := &db.WfRequest{
			RequestID:   uuid.New().String(),
			Source:      "test_source",
			Action:      "test_action",
			Status:      db.StatusNew,
			WorkflowID:  "test_workflow",
			RequestData: requestData,
		}
		err := db.TestDB.CreateWorkflowRequest(workflowRequest)
		assert.NoError(t, err)

		processingMap := &db.WfProcessingMap{
			Type:               workflowRequest.Source,
			ProcessKey:         workflowRequest.Action,
			RequiresProcessing: true,
			HandlerFunc:        "test_handler",
		}
		err = db.TestDB.CreateProcessingMap(processingMap)
		assert.NoError(t, err)

		response := struct {
			RequestID    string         `json:"request_id"`
			ResponseData db.PropertyMap `json:"response_data"`
		}{
			RequestID:    workflowRequest.RequestID,
			ResponseData: db.PropertyMap{"result": "success"},
		}
		payload, _ := json.Marshal(response)

		req := httptest.NewRequest(http.MethodPost, "/workflows/response", bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		wh.HandleWorkflowResponse(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify the request status is set to PENDING due to processing map
		updatedReq, err := db.TestDB.GetWorkflowRequest(workflowRequest.RequestID)
		assert.NoError(t, err)
		assert.Equal(t, db.StatusPending, updatedReq.Status)
		assert.Equal(t, response.ResponseData, updatedReq.ResponseData)
	})
}