package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
)

type activityHandler struct {
	httpClient HttpClient
	db         db.Database
}

func NewActivityHandler(httpClient HttpClient, database db.Database) *activityHandler {
	return &activityHandler{
		httpClient: httpClient,
		db:         database,
	}
}

type CreateActivityRequest struct {
	ContentType string         `json:"content_type"`
	Content     string         `json:"content"`
	Workspace   string         `json:"workspace"`
	FeatureUUID string         `json:"feature_uuid"`
	PhaseUUID   string         `json:"phase_uuid"`
	Actions     []string       `json:"actions,omitempty"`
	Questions   []string       `json:"questions,omitempty"`
	Author      db.AuthorType  `json:"author"`
	AuthorRef   string         `json:"author_ref"`
}

type ActivityResponse struct {
	Success bool         `json:"success"`
	Data    *db.Activity `json:"data,omitempty"`
	Error   string       `json:"error,omitempty"`
}

func (ah *activityHandler) GetActivity(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	activity, err := ah.db.GetActivity(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "activity not found" {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	json.NewEncoder(w).Encode(ActivityResponse{
		Success: true,
		Data:    activity,
	})
}

func (ah *activityHandler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	activity := &db.Activity{
		ID:          uuid.New(),
		ContentType: db.ContentType(req.ContentType),
		Content:     req.Content,
		Workspace:   req.Workspace,
		FeatureUUID: req.FeatureUUID,
		PhaseUUID:   req.PhaseUUID,
		Actions:     req.Actions,
		Questions:   req.Questions,
		Author:      req.Author,
		AuthorRef:   req.AuthorRef,
		TimeCreated: time.Now(),
		TimeUpdated: time.Now(),
		Status:      "active",
	}

	createdActivity, err := ah.db.CreateActivity(activity)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create activity: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ActivityResponse{
		Success: true,
		Data:    createdActivity,
	})
}

func (ah *activityHandler) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var activity db.Activity
	if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	activity.TimeUpdated = time.Now()
	updatedActivity, err := ah.db.UpdateActivity(&activity)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update activity: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ActivityResponse{
		Success: true,
		Data:    updatedActivity,
	})
}

func (ah *activityHandler) CreateActivityThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sourceID := r.URL.Query().Get("source_id")
	if sourceID == "" {
		http.Error(w, "source_id query parameter is required", http.StatusBadRequest)
		return
	}

	activity := &db.Activity{
		ContentType: db.ContentType(req.ContentType),
		Content:     req.Content,
		Workspace:   req.Workspace,
		FeatureUUID: req.FeatureUUID,
		PhaseUUID:   req.PhaseUUID,
		Actions:     req.Actions,
		Questions:   req.Questions,
		Author:      req.Author,
		AuthorRef:   req.AuthorRef,
	}

	createdActivity, err := ah.db.CreateActivityThread(sourceID, activity)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create activity thread: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ActivityResponse{
		Success: true,
		Data:    createdActivity,
	})
}

func (ah *activityHandler) GetActivitiesByThread(w http.ResponseWriter, r *http.Request) {
	threadID := chi.URLParam(r, "thread_id")
	if threadID == "" {
		http.Error(w, "thread_id is required", http.StatusBadRequest)
		return
	}

	activities, err := ah.db.GetActivitiesByThread(threadID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get activities: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    activities,
	})
}

func (ah *activityHandler) GetLatestActivityByThread(w http.ResponseWriter, r *http.Request) {
	threadID := chi.URLParam(r, "thread_id")
	if threadID == "" {
		http.Error(w, "thread_id is required", http.StatusBadRequest)
		return
	}

	activity, err := ah.db.GetLatestActivityByThread(threadID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get latest activity: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ActivityResponse{
		Success: true,
		Data:    activity,
	})
}

func (ah *activityHandler) GetActivitiesByFeature(w http.ResponseWriter, r *http.Request) {
	featureUUID := chi.URLParam(r, "feature_uuid")
	if featureUUID == "" {
		http.Error(w, "feature_uuid is required", http.StatusBadRequest)
		return
	}

	activities, err := ah.db.GetActivitiesByFeature(featureUUID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get activities: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    activities,
	})
}

func (ah *activityHandler) GetActivitiesByPhase(w http.ResponseWriter, r *http.Request) {
	phaseUUID := chi.URLParam(r, "phase_uuid")
	if phaseUUID == "" {
		http.Error(w, "phase_uuid is required", http.StatusBadRequest)
		return
	}

	activities, err := ah.db.GetActivitiesByPhase(phaseUUID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get activities: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    activities,
	})
}

func (ah *activityHandler) GetActivitiesByWorkspace(w http.ResponseWriter, r *http.Request) {
	workspace := chi.URLParam(r, "workspace")
	if workspace == "" {
		http.Error(w, "workspace is required", http.StatusBadRequest)
		return
	}

	activities, err := ah.db.GetActivitiesByWorkspace(workspace)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get activities: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    activities,
	})
}

func (ah *activityHandler) DeleteActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
    if id == "" {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }

    err := ah.db.DeleteActivity(id)
    if err != nil {
        if err.Error() == "activity not found" {
            http.Error(w, "Activity not found", http.StatusNotFound)
            return
        }
        http.Error(w, fmt.Sprintf("Failed to delete activity: %v", err), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "message": "Activity deleted successfully",
    })
}

type ActivityContentRequest struct {
	Content string `json:"content"`
}

func (ah *activityHandler) AddActivityActions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var req ActivityContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	activity, err := ah.db.GetActivity(id)
	if err != nil {
		if err.Error() == "activity not found" {
			http.Error(w, "Activity not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	activity.Actions = append(activity.Actions, req.Content)
	activity.TimeUpdated = time.Now()

	updatedActivity, err := ah.db.UpdateActivity(activity)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add action: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ActivityResponse{
		Success: true,
		Data:    updatedActivity,
	})
}

func (ah *activityHandler) AddActivityQuestions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var req ActivityContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	activity, err := ah.db.GetActivity(id)
	if err != nil {
		if err.Error() == "activity not found" {
			http.Error(w, "Activity not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	activity.Questions = append(activity.Questions, req.Content)
	activity.TimeUpdated = time.Now()

	updatedActivity, err := ah.db.UpdateActivity(activity)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add question: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ActivityResponse{
		Success: true,
		Data:    updatedActivity,
	})
}

func (ah *activityHandler) RemoveActivityAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	actionID := chi.URLParam(r, "action_id")
	if id == "" || actionID == "" {
		http.Error(w, "ID and action_id are required", http.StatusBadRequest)
		return
	}

	activity, err := ah.db.GetActivity(id)
	if err != nil {
		if err.Error() == "activity not found" {
			http.Error(w, "Activity not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	actionIndex := -1
	for i, action := range activity.Actions {
		if action == actionID {
			actionIndex = i
			break
		}
	}

	if actionIndex == -1 {
		http.Error(w, "Action not found", http.StatusNotFound)
		return
	}

	activity.Actions = append(activity.Actions[:actionIndex], activity.Actions[actionIndex+1:]...)
	activity.TimeUpdated = time.Now()

	updatedActivity, err := ah.db.UpdateActivity(activity)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove action: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ActivityResponse{
		Success: true,
		Data:    updatedActivity,
	})
}

func (ah *activityHandler) RemoveActivityQuestion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	questionID := chi.URLParam(r, "question_id")
	if id == "" || questionID == "" {
		http.Error(w, "ID and question_id are required", http.StatusBadRequest)
		return
	}

	activity, err := ah.db.GetActivity(id)
	if err != nil {
		if err.Error() == "activity not found" {
			http.Error(w, "Activity not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	questionIndex := -1
	for i, question := range activity.Questions {
		if question == questionID {
			questionIndex = i
			break
		}
	}

	if questionIndex == -1 {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	activity.Questions = append(activity.Questions[:questionIndex], activity.Questions[questionIndex+1:]...)
	activity.TimeUpdated = time.Now()

	updatedActivity, err := ah.db.UpdateActivity(activity)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove question: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ActivityResponse{
		Success: true,
		Data:    updatedActivity,
	})
} 