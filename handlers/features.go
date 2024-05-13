package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
)

type featureHandler struct {
	db db.Database
}

func NewFeatureHandler(database db.Database) *featureHandler {
	return &featureHandler{
		db: database,
	}
}

func (oh *featureHandler) CreateOrEditFeatures(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	features := db.WorkspaceFeatures{}
	body, _ := io.ReadAll(r.Body)
	r.Body.Close()
	err := json.Unmarshal(body, &features)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	features.CreatedBy = pubKeyFromAuth

	// Validate struct data
	err = db.Validate.Struct(features)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Error: did not pass validation test : %s", err)
		json.NewEncoder(w).Encode(msg)
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

func (oh *featureHandler) GetFeaturesByWorkspaceUuid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "uuid")
	workspaceFeatures := oh.db.GetFeaturesByWorkspaceUuid(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceFeatures)
}

func (oh *featureHandler) GetFeatureByUuid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "uuid")
	workspaceFeature := oh.db.GetFeatureByUuid(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceFeature)
}

func (oh *featureHandler) CreateOrEditStory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newStory := db.FeatureStory{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newStory)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body: %v", err)
		return
	}

	newStory.CreatedBy = pubKeyFromAuth

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
	storyUuid := chi.URLParam(r, "story_uuid")
	story, err := oh.db.GetFeatureStoryByUuid(storyUuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(story)
}

func (oh *featureHandler) DeleteStory(w http.ResponseWriter, r *http.Request) {
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
