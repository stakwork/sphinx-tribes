package handlers

import (
	"encoding/json"

	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/utils"

	"io"
	"net/http"
)

type workflowHandler struct {
	db db.Database
}

func NewWorkFlowHandler(database db.Database) *workflowHandler {
	return &workflowHandler{
		db: database,
	}
}

func (wh *workflowHandler) HandleWorkflowRequest(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var request db.WfRequest
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if request.WorkflowID == "" || request.Source == "" {
		http.Error(w, "Missing required fields: workflow_id or source", http.StatusBadRequest)
		return
	}

	processedRequestID, err := utils.ProcessWorkflowRequest(request.RequestID, request.Source)
	if err != nil {
		panic("Failed to process workflow request")
		return
	}

	request.RequestID = processedRequestID
	request.Status = db.StatusNew

	if err := wh.db.CreateWorkflowRequest(&request); err != nil {
		panic("Failed to create workflow request")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"request_id": processedRequestID,
		"status":     "success",
	})
}

func (wh *workflowHandler) HandleWorkflowResponse(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusBadRequest)
		return
	}

	var response struct {
		RequestID    string         `json:"request_id"`
		ResponseData db.PropertyMap `json:"response_data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		http.Error(w, "Invalid response format", http.StatusBadRequest)
		return
	}

	if response.RequestID == "" {
		http.Error(w, "Request ID is required", http.StatusBadRequest)
		return
	}

	request, err := wh.db.GetWorkflowRequest(response.RequestID)
	if err != nil {
		panic("Failed to retrieve original request")
		return
	}
	if request == nil {
		http.Error(w, "Original request not found", http.StatusNotFound)
		return
	}

	processingMap, err := wh.db.GetProcessingMapByKey(request.Source, request.Action)
	if err != nil {
		panic("Failed to check processing requirements")
		return
	}

	status := db.StatusCompleted
	if processingMap != nil && processingMap.RequiresProcessing {
		status = db.StatusPending
	}

	err = wh.db.UpdateWorkflowRequestStatusAndResponse(
		response.RequestID,
		status,
		response.ResponseData,
	)
	if err != nil {
		panic("Failed to update workflow request")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":     "success",
		"request_id": response.RequestID,
	})
}
