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
		http.Error(w, "Failed to process workflow request", http.StatusInternalServerError)
		return
	}

	request.Status = db.StatusNew

	if err := wh.db.CreateWorkflowRequest(&request); err != nil {
		http.Error(w, "Failed to create workflow request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"request_id": processedRequestID,
		"status":     "success",
	})
}
