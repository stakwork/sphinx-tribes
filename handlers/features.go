package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"os"

	"github.com/go-chi/chi"
	"github.com/rs/xid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
)

type PostData struct {
	ProductBrief string   `json:"productBrief"`
	FeatureName  string   `json:"featureName"`
	Description  string   `json:"description"`
	Examples     []string `json:"examples"`
	WebhookURL   string   `json:"webhook_url"`
	FeatureUUID  string   `json:"featureUUID"`
}

type FeatureBriefRequest struct {
	Output struct {
		FeatureBrief string `json:"featureBrief"`
		AudioLink    string `json:"audioLink"`
		FeatureUUID  string `json:"featureUUID"`
	} `json:"output"`
}
type AudioBriefPostData struct {
	AudioLink   string   `json:"audioLink"`
	FeatureUUID string   `json:"featureUUID"`
	Source      string   `json:"source"`
	Examples    []string `json:"examples"`
}

type featureHandler struct {
	db                    db.Database
	generateBountyHandler func(bounties []db.NewBounty) []db.BountyResponse
}

func NewFeatureHandler(database db.Database) *featureHandler {
	bHandler := NewBountyHandler(http.DefaultClient, database)
	return &featureHandler{
		db:                    database,
		generateBountyHandler: bHandler.GenerateBountyResponse,
	}
}

func (oh *featureHandler) CreateOrEditFeatures(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	features := db.WorkspaceFeatures{}
	body, _ := io.ReadAll(r.Body)
	r.Body.Close()
	err := json.Unmarshal(body, &features)

	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	features.CreatedBy = pubKeyFromAuth

	if features.Uuid == "" {
		features.Uuid = xid.New().String()
		features.FeatStatus = db.ActiveFeature
	} else {
		features.UpdatedBy = pubKeyFromAuth
	}

	// Validate struct data
	err = db.Validate.Struct(features)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Error: did not pass validation test : %s", err)
		json.NewEncoder(w).Encode(msg)
		return
	}

	// Check if workspace exists
	workpace := oh.db.GetWorkspaceByUuid(features.WorkspaceUuid)
	if workpace.Uuid != features.WorkspaceUuid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Workspace does not exists")
		return
	}

	p, err := oh.db.CreateOrEditFeature(features)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

func (oh *featureHandler) DeleteFeature(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "uuid")
	err := oh.db.DeleteFeatureByUuid(uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Feature deleted successfully")
}

// Old Method for getting features for workspace uuid
func (oh *featureHandler) GetFeaturesByWorkspaceUuid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "workspace_uuid")
	workspaceFeatures := oh.db.GetFeaturesByWorkspaceUuid(uuid, r)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceFeatures)
}

func (oh *featureHandler) GetWorkspaceFeaturesCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "uuid")
	workspaceFeatures := oh.db.GetWorkspaceFeaturesCount(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceFeatures)
}

func (oh *featureHandler) GetFeatureByUuid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "uuid")
	workspaceFeature := oh.db.GetFeatureByUuid(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceFeature)
}

func (oh *featureHandler) UpdateFeatureBrief(w http.ResponseWriter, r *http.Request) {

	var req FeatureBriefRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request payload")
		return
	}

	featureUUID := req.Output.FeatureUUID
	newFeatureBrief := req.Output.FeatureBrief

	if featureUUID == "" || newFeatureBrief == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing required fields")
		return
	}

	prevFeatureBrief := oh.db.GetFeatureByUuid(featureUUID)

	if prevFeatureBrief.Uuid == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Feature not found")
		return
	}

	var updatedFeatureBrief string
	if prevFeatureBrief.Brief == "" {
		updatedFeatureBrief = newFeatureBrief
	} else {

		updatedFeatureBrief = prevFeatureBrief.Brief + "\n\n* Generated Feature Brief *\n\n" + newFeatureBrief
	}

	featureToUpdate := db.WorkspaceFeatures{
		Uuid:                   featureUUID,
		WorkspaceUuid:          prevFeatureBrief.WorkspaceUuid,
		Name:                   prevFeatureBrief.Name,
		Brief:                  updatedFeatureBrief,
		Requirements:           prevFeatureBrief.Requirements,
		Architecture:           prevFeatureBrief.Architecture,
		Url:                    prevFeatureBrief.Url,
		Priority:               prevFeatureBrief.Priority,
		BountiesCountCompleted: prevFeatureBrief.BountiesCountCompleted,
		BountiesCountAssigned:  prevFeatureBrief.BountiesCountAssigned,
		BountiesCountOpen:      prevFeatureBrief.BountiesCountOpen,
	}

	p, err := oh.db.CreateOrEditFeature(featureToUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

func (oh *featureHandler) CreateOrEditFeaturePhase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newPhase := db.FeaturePhase{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newPhase)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	if newPhase.Uuid == "" {
		newPhase.Uuid = xid.New().String()
	}

	existingPhase, _ := oh.db.GetFeaturePhaseByUuid(newPhase.FeatureUuid, newPhase.Uuid)

	if existingPhase.CreatedBy == "" {
		newPhase.CreatedBy = pubKeyFromAuth
	}

	newPhase.UpdatedBy = pubKeyFromAuth

	// Check if feature exists
	feature := oh.db.GetFeatureByUuid(newPhase.FeatureUuid)
	if feature.Uuid != newPhase.FeatureUuid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Feature does not exists")
		return
	}

	phase, err := oh.db.CreateOrEditFeaturePhase(newPhase)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating feature phase: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(phase)
}

func (oh *featureHandler) GetFeaturePhases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	featureUuid := chi.URLParam(r, "feature_uuid")
	phases := oh.db.GetPhasesByFeatureUuid(featureUuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(phases)
}

func (oh *featureHandler) GetFeaturePhaseByUUID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	featureUuid := chi.URLParam(r, "feature_uuid")
	phaseUuid := chi.URLParam(r, "phase_uuid")

	phase, err := oh.db.GetFeaturePhaseByUuid(featureUuid, phaseUuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(phase)
}

func (oh *featureHandler) DeleteFeaturePhase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	featureUuid := chi.URLParam(r, "feature_uuid")
	phaseUuid := chi.URLParam(r, "phase_uuid")

	err := oh.db.DeleteFeaturePhase(featureUuid, phaseUuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Phase deleted successfully"})
}

func (oh *featureHandler) CreateOrEditStory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newStory := db.FeatureStory{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newStory)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	if newStory.Uuid == "" {
		newStory.Uuid = xid.New().String()
	}

	existingStory, _ := oh.db.GetFeatureStoryByUuid(newStory.FeatureUuid, newStory.Uuid)

	if existingStory.CreatedBy == "" {
		newStory.CreatedBy = pubKeyFromAuth
	}

	newStory.UpdatedBy = pubKeyFromAuth

	story, err := oh.db.CreateOrEditFeatureStory(newStory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating feature story: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(story)
}

func (oh *featureHandler) GetStoriesByFeatureUuid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	featureUuid := chi.URLParam(r, "feature_uuid")
	stories, err := oh.db.GetFeatureStoriesByFeatureUuid(featureUuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stories)
}

func (oh *featureHandler) GetStoryByUuid(w http.ResponseWriter, r *http.Request) {
	featureUuid := chi.URLParam(r, "feature_uuid")
	storyUuid := chi.URLParam(r, "story_uuid")

	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	story, err := oh.db.GetFeatureStoryByUuid(featureUuid, storyUuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(story)
}
func (oh *featureHandler) DeleteStory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	featureUuid := chi.URLParam(r, "feature_uuid")
	storyUuid := chi.URLParam(r, "story_uuid")

	err := oh.db.DeleteFeatureStoryByUuid(featureUuid, storyUuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Story deleted successfully"})
}

func (oh *featureHandler) GetBountiesByFeatureAndPhaseUuid(w http.ResponseWriter, r *http.Request) {
	featureUuid := chi.URLParam(r, "feature_uuid")
	phaseUuid := chi.URLParam(r, "phase_uuid")

	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bounties, err := oh.db.GetBountiesByFeatureAndPhaseUuid(featureUuid, phaseUuid, r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bountyResponses := oh.generateBountyHandler(bounties)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyResponses)
}

func (oh *featureHandler) GetBountiesCountByFeatureAndPhaseUuid(w http.ResponseWriter, r *http.Request) {
	featureUuid := chi.URLParam(r, "feature_uuid")
	phaseUuid := chi.URLParam(r, "phase_uuid")

	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bountiesCount := oh.db.GetBountiesCountByFeatureAndPhaseUuid(featureUuid, phaseUuid, r)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountiesCount)
}

func (oh *featureHandler) GetFeatureStories(w http.ResponseWriter, r *http.Request) {
	featureStories := db.FeatureStoriesReponse{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&featureStories)

	featureUuid := featureStories.Output.FeatureUuid

	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	log.Println("Webhook Feature Uuid", featureUuid)

	log.Println("Webhook Feature Stories === ", featureStories.Output.Stories)

	// check if feature story exists
	feature := oh.db.GetFeatureByUuid(featureUuid)

	if feature.ID == 0 {
		msg := "Feature ID does not exists"
		log.Println(msg, featureUuid)
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(msg)
		return
	}

	for _, story := range featureStories.Output.Stories {

		now := time.Now()

		// Add story to database
		featureStory := db.FeatureStory{
			Uuid:        xid.New().String(),
			Description: story.UserStory,
			FeatureUuid: featureUuid,
			Created:     &now,
			Updated:     &now,
		}

		oh.db.CreateOrEditFeatureStory(featureStory)
		log.Println("Created user story for : ", featureStory.FeatureUuid)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User stories added successfully")
}

func (oh *featureHandler) StoriesSend(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		http.Error(w, "Failed to read requests body", http.StatusBadRequest)
		return
	}

	var postData PostData
	err = json.Unmarshal(body, &postData)
	if err != nil {
		logger.Log.Error("[StoriesSend] JSON Unmarshal error: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusNotAcceptable)
		return
	}

	apiKey := os.Getenv("SWWFKEY")
	if apiKey == "" {
		panic("API key not set in environment")
		return
	}

	stakworkPayload := map[string]interface{}{
		"name":        "string",
		"workflow_id": 35080,
		"workflow_params": map[string]interface{}{
			"set_var": map[string]interface{}{
				"attributes": map[string]interface{}{
					"vars": postData,
				},
			},
		},
	}

	stakworkPayloadJSON, err := json.Marshal(stakworkPayload)
	if err != nil {
		panic("Failed to encode payload")
		return
	}

	req, err := http.NewRequest("POST", "https://api.stakwork.com/api/v1/projects", bytes.NewBuffer(stakworkPayloadJSON))
	if err != nil {
		panic("Failed to create request to Stakwork API")
		return
	}
	req.Header.Set("Authorization", "Token token="+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic("Failed to send request to Stakwork API")
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("Failed to read response from Stakwork API")
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func (oh *featureHandler) BriefSend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		http.Error(w, "Failed to read requests body", http.StatusBadRequest)
		return
	}

	var postData AudioBriefPostData
	err = json.Unmarshal(body, &postData)
	if err != nil {
		logger.Log.Error("[BriefSend] JSON Unmarshal error: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusNotAcceptable)
		return
	}

	host := os.Getenv("HOST")
	if host == "" {
		panic("HOST environment variable not set")
		return
	}

	completePostData := struct {
		AudioBriefPostData
		WebhookURL string `json:"webhook_url"`
	}{
		AudioBriefPostData: postData,
		WebhookURL:         fmt.Sprintf("%s/feature/brief", host),
	}

	apiKey := os.Getenv("SWWFKEY")
	if apiKey == "" {
		panic("API key not set in environment")
		return
	}

	stakworkPayload := map[string]interface{}{
		"name":        "string",
		"workflow_id": 36928,
		"workflow_params": map[string]interface{}{
			"set_var": map[string]interface{}{
				"attributes": map[string]interface{}{
					"vars": completePostData,
				},
			},
		},
	}

	stakworkPayloadJSON, err := json.Marshal(stakworkPayload)
	if err != nil {
		panic("Failed to encode payload")
		return
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.stakwork.com/api/v1/projects", bytes.NewBuffer(stakworkPayloadJSON))
	if err != nil {
		panic("Failed to create request to Stakwork API")
		return
	}
	req.Header.Set("Authorization", "Token token="+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic("Failed to send request to Stakwork API")
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("Failed to read response from Stakwork API")
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func (oh *featureHandler) UpdateFeatureStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "uuid")
	var req struct {
		Status db.FeatureStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("invalid request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Status != db.ActiveFeature && req.Status != db.ArchivedFeature {
		logger.Log.Info("invalid feature status")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedFeature, err := oh.db.UpdateFeatureStatus(uuid, req.Status)
	if err != nil {
		logger.Log.Error("failed to update feature status", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedFeature)
}
