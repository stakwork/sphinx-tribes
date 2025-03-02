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
	ContentType string        `json:"content_type"`
	Title       string        `json:"title,omitempty"`
	Content     string        `json:"content"`
	Workspace   string        `json:"workspace"`
	FeatureUUID string        `json:"feature_uuid"`
	PhaseUUID   string        `json:"phase_uuid"`
	Actions     []string      `json:"actions,omitempty"`
	Questions   []string      `json:"questions,omitempty"`
	Author      db.AuthorType `json:"author"`
	AuthorRef   string        `json:"author_ref"`
}

type ActivityResponse struct {
	Success bool         `json:"success"`
	Data    *db.Activity `json:"data,omitempty"`
	Error   string       `json:"error,omitempty"`
}

type ActivitiesResponse struct {
	Success bool          `json:"success"`
	Data    []db.Activity `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
}

type ActivityThreadResponse struct {
	Success bool          `json:"success"`
	Data    []db.Activity `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
}

type WebhookActivityRequest struct {
	ContentType string        `json:"content_type"`
	Title       string        `json:"title,omitempty"`
	Content     string        `json:"content"`
	Workspace   string        `json:"workspace"`
	ThreadID    string        `json:"thread_id,omitempty"`
	FeatureUUID string        `json:"feature_uuid,omitempty"`
	PhaseUUID   string        `json:"phase_uuid,omitempty"`
	Actions     []string      `json:"actions,omitempty"`
	Questions   []string      `json:"questions,omitempty"`
	Author      db.AuthorType `json:"author"`
	AuthorRef   string        `json:"author_ref"`
}

type WebhookResponse struct {
	Success    bool   `json:"success"`
	ActivityID string `json:"activity_id,omitempty"`
	Error      string `json:"error,omitempty"`
}

// GetActivity godoc
//
//	@Summary		Retrieve activity details
//	@Description	Fetch a specific activity by its unique identifier
//	@Tags			Activities
//	@Param			id	path	string	true	"Activity ID"
//	@Produce		json
//	@Success		200	{object}	ActivityResponse
//	@Failure		400	{string}	string	"ID is required"
//	@Failure		404	{string}	string	"activity not found"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/activities/{id} [get]
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

// CreateActivity godoc
//
//	@Summary		Create an activity
//	@Description	Create a new activity
//	@Tags			Activities
//	@Param			activity	body	CreateActivityRequest	true	"Activity object"
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		201	{object}	ActivityResponse
//	@Router			/activities [post]
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
		Title:       req.Title,
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

// UpdateActivity godoc
//
//	@Summary		Update an activity
//	@Description	Update an existing activity
//	@Tags			Activities
//	@Param			id			path	string		true	"Activity ID"
//	@Param			activity	body	db.Activity	true	"Activity object"
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{object}	ActivityResponse
//	@Router			/activities/{id} [put]
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

	existing, err := ah.db.GetActivity(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get activity: %v", err), http.StatusInternalServerError)
		return
	}

	var updateReq db.Activity
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateReq.ID = existing.ID
	updateReq.ThreadID = existing.ThreadID
	updateReq.Sequence = existing.Sequence
	updateReq.TimeUpdated = time.Now()

	updatedActivity, err := ah.db.UpdateActivity(&updateReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update activity: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ActivityResponse{
		Success: true,
		Data:    updatedActivity,
	})
}

// CreateActivityThread godoc
//
//	@Summary		Create an activity thread
//	@Description	Create a new activity thread
//	@Tags			Activities
//	@Param			activity	body	CreateActivityRequest	true	"Activity object"
//	@Param			source_id	query	string					true	"Source ID"
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		201	{object}	ActivityResponse
//	@Router			/activities/thread [post]
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

// GetActivitiesByThread godoc
//
//	@Summary		Get activities by thread
//	@Description	Get activities by thread ID
//	@Tags			Activities
//	@Param			thread_id	path	string	true	"Thread ID"
//	@Produce		json
//	@Success		200	{object}	ActivityThreadResponse
//	@Failure		400	{string}	string	"thread_id is required"
//	@Failure		500	{string}	string	"Failed to get activities"
//	@Router			/activities/thread/{thread_id} [get]
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

	json.NewEncoder(w).Encode(ActivityThreadResponse{
		Success: true,
		Data:    activities,
	})
}

// GetLatestActivityByThread godoc
//
//	@Summary		Get the latest activity by thread
//	@Description	Get the latest activity by thread ID
//	@Tags			Activities
//	@Param			thread_id	path	string	true	"Thread ID"
//	@Produce		json
//	@Success		200	{object}	ActivityResponse
//	@Failure		400	{string}	string	"thread_id is required"
//	@Failure		500	{string}	string	"Failed to get latest activity"
//	@Router			/activities/thread/{thread_id}/latest [get]
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

// GetActivitiesByFeature godoc
//
//	@Summary		Get activities by feature
//	@Description	Get activities by feature UUID
//	@Tags			Activities
//	@Param			feature_uuid	path	string	true	"Feature UUID"
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{object}	ActivitiesResponse
//	@Router			/activities/feature/{feature_uuid} [get]
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

	json.NewEncoder(w).Encode(ActivitiesResponse{
		Success: true,
		Data:    activities,
	})
}

// GetActivitiesByPhase godoc
//
//	@Summary		Get activities by phase
//	@Description	Get activities by phase UUID
//	@Tags			Activities
//	@Param			phase_uuid	path	string	true	"Phase UUID"
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{object}	ActivitiesResponse
//	@Router			/activities/phase/{phase_uuid} [get]
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

	json.NewEncoder(w).Encode(ActivitiesResponse{
		Success: true,
		Data:    activities,
	})
}

// GetActivitiesByWorkspace godoc
//
//	@Summary		Get activities by workspace
//	@Description	Get activities by workspace
//	@Tags			Activities
//	@Param			workspace	path	string	true	"Workspace"
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{object}	ActivitiesResponse
//	@Router			/activities/workspace/{workspace} [get]
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

	json.NewEncoder(w).Encode(ActivitiesResponse{
		Success: true,
		Data:    activities,
	})
}

// DeleteActivity godoc
//
//	@Summary		Delete an activity
//	@Description	Delete an activity by ID
//	@Tags			Activities
//	@Param			id	path	string	true	"Activity ID"
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{object}	map[string]interface{}
//	@Router			/activities/{id} [delete]
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

// AddActivityActions godoc
//
//	@Summary		Add actions to an activity
//	@Description	Add actions to an activity by ID
//	@Tags			Activities
//	@Param			id		path	string					true	"Activity ID"
//	@Param			action	body	ActivityContentRequest	true	"Action content"
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{object}	ActivityResponse
//	@Router			/activities/{id}/actions [post]
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

// AddActivityQuestions godoc
//
//	@Summary		Add questions to an activity
//	@Description	Add questions to an activity by ID
//	@Tags			Activities
//	@Param			id			path	string					true	"Activity ID"
//	@Param			question	body	ActivityContentRequest	true	"Question content"
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{object}	ActivityResponse
//	@Router			/activities/{id}/questions [post]
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

// RemoveActivityAction godoc
//
//	@Summary		Remove an action from an activity
//	@Description	Remove an action from an activity by ID
//	@Tags			Activities
//	@Param			id			path	string	true	"Activity ID"
//	@Param			action_id	path	string	true	"Action ID"
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{object}	ActivityResponse
//	@Router			/activities/{id}/actions/{action_id} [delete]
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

// RemoveActivityQuestion godoc
//
//	@Summary		Remove a question from an activity
//	@Description	Remove a question from an activity by ID
//	@Tags			Activities
//	@Param			id			path	string	true	"Activity ID"
//	@Param			question_id	path	string	true	"Question ID"
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{object}	ActivityResponse
//	@Router			/activities/{id}/questions/{question_id} [delete]
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

// ReceiveActivity godoc
//
//	@Summary		Receive and process a new activity
//	@Description	Receives activity data from a webhook, validates it, and creates a new activity record
//	@Tags			Activities
//	@Accept			json
//	@Produce		json
//	@Param			request	body		WebhookActivityRequest	true	"Activity information"
//	@Success		201		{object}	WebhookResponse
//	@Failure		400		{object}	WebhookResponse	"Invalid request payload, invalid public key format, invalid source ID format, or other validation errors"
//	@Failure		500		{object}	WebhookResponse	"Internal server error"
//	@Router			/activities/receive [post]
func (ah *activityHandler) ReceiveActivity(w http.ResponseWriter, r *http.Request) {
	var req WebhookActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(WebhookResponse{
			Success: false,
			Error:   "Invalid request payload",
		})
		return
	}

	if req.Author == db.HumansAuthor {
		if len(req.AuthorRef) < 32 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(WebhookResponse{
				Success: false,
				Error:   "invalid public key format for human author",
			})
			return
		}
	}

	if req.ThreadID != "" {
		if _, err := uuid.Parse(req.ThreadID); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(WebhookResponse{
				Success: false,
				Error:   "invalid source ID format",
			})
			return
		}
	}

	activity := &db.Activity{
		ID:          uuid.New(),
		Title:       req.Title,
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

	var createdActivity *db.Activity
	var err error

	if req.ThreadID != "" {
		createdActivity, err = ah.db.CreateActivityThread(req.ThreadID, activity)
	} else {
		createdActivity, err = ah.db.CreateActivity(activity)
	}

	if err != nil {
		status := http.StatusInternalServerError
		if err == db.ErrInvalidContent || err == db.ErrInvalidAuthorRef ||
			err == db.ErrInvalidContentType || err == db.ErrInvalidAuthorType ||
			err == db.ErrInvalidWorkspace {
			status = http.StatusBadRequest
		}

		w.WriteHeader(status)
		json.NewEncoder(w).Encode(WebhookResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(WebhookResponse{
		Success:    true,
		ActivityID: createdActivity.ID.String(),
	})
}
