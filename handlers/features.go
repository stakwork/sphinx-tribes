package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/rs/xid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"io"
	"net/http"
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

	if features.Uuid == "" {
		features.Uuid = xid.New().String()
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

	p, err := oh.db.CreateOrEditFeature(features)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

func (oh *featureHandler) GetWorkspaceFeaturesCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
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
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "uuid")
	workspaceFeature := oh.db.GetFeatureByUuid(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaceFeature)
}

func (oh *featureHandler) DeleteFeature(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
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
