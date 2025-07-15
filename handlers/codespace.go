package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
)

type codeSpaceHandler struct {
	db db.Database
}

func NewCodeSpaceHandler(database db.Database) *codeSpaceHandler {
	return &codeSpaceHandler{
		db: database,
	}
}

type CodeSpaceQuery struct {
	WorkspaceID string `json:"workspaceID"`
	UserPubkey  string `json:"userPubkey"`
}

type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (ch *codeSpaceHandler) GetAllCodeSpaceMaps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[codespace] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	codespaces, err := ch.db.GetCodeSpaceMaps()
	if err != nil {
		logger.Log.Error("[codespace] error getting codespace mappings: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve codespace mappings"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(codespaces)
}

func (ch *codeSpaceHandler) GetCodeSpaceMapsByWorkspace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[codespace] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	workspaceID := chi.URLParam(r, "workspaceID")
	if workspaceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Workspace ID is required"})
		return
	}

	codespaces, err := ch.db.GetCodeSpaceMapByWorkspace(workspaceID)
	if err != nil {
		logger.Log.Error("[codespace] error getting codespace mappings: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve codespace mappings"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(codespaces)
}

func (ch *codeSpaceHandler) GetCodeSpaceMapByWorkspaceAndUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[codespace] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	workspaceID := chi.URLParam(r, "workspaceID")
	userPubkey := chi.URLParam(r, "userPubkey")

	if workspaceID == "" || userPubkey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Workspace ID and User public key are required"})
		return
	}

	codespace, err := ch.db.GetCodeSpaceMapByWorkspaceAndUser(workspaceID, userPubkey)
	if err != nil {
		if err.Error() == "codespace mapping not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Codespace mapping not found"})
			return
		}
		logger.Log.Error("[codespace] error getting codespace mapping: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve codespace mapping"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(codespace)
}

func (ch *codeSpaceHandler) GetCodeSpaceMapsByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[codespace] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	userPubkey := chi.URLParam(r, "userPubkey")
	if userPubkey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User public key is required"})
		return
	}

	codespaces, err := ch.db.GetCodeSpaceMapByUser(userPubkey)
	if err != nil {
		logger.Log.Error("[codespace] error getting codespace mappings: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve codespace mappings"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(codespaces)
}

func (ch *codeSpaceHandler) GetCodeSpaceMapsByURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[codespace] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	codeSpaceURL := r.URL.Query().Get("codeSpaceURL")
	if codeSpaceURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "CodeSpace URL is required"})
		return
	}

	codespaces, err := ch.db.GetCodeSpaceMapByURL(codeSpaceURL)
	if err != nil {
		logger.Log.Error("[codespace] error getting codespace mappings: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve codespace mappings"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(codespaces)
}

func (ch *codeSpaceHandler) QueryCodeSpaceMaps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[codespace] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	workspaceID := r.URL.Query().Get("workspaceID")
	userPubkey := r.URL.Query().Get("userPubkey")

	logger.Log.Info("[codespace] Query params - workspaceID: %s, userPubkey: %s", workspaceID, userPubkey)

	if workspaceID != "" && userPubkey != "" {
		logger.Log.Info("[codespace] Querying by workspace and user")
		codeSpace, err := ch.db.GetCodeSpaceMapByWorkspaceAndUser(workspaceID, userPubkey)
		if err != nil {
			if err.Error() == "codespace mapping not found" {
				logger.Log.Info("[codespace] No mapping found for workspace %s and user %s", workspaceID, userPubkey)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]db.CodeSpaceMap{})
				return
			}
			logger.Log.Error("[codespace] error querying codespace mapping: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to query codespace mapping"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]db.CodeSpaceMap{codeSpace})
		return
	}

	if workspaceID != "" {
		logger.Log.Info("[codespace] Querying by workspace")
		codespaces, err := ch.db.GetCodeSpaceMapByWorkspace(workspaceID)
		if err != nil {
			logger.Log.Error("[codespace] error querying codespace mappings: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to query codespace mappings"})
			return
		}
		logger.Log.Info("[codespace] Found %d mappings for workspace %s", len(codespaces), workspaceID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(codespaces)
		return
	}

	if userPubkey != "" {
		logger.Log.Info("[codespace] Querying by user")
		codespaces, err := ch.db.GetCodeSpaceMapByUser(userPubkey)
		if err != nil {
			logger.Log.Error("[codespace] error querying codespace mappings: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to query codespace mappings"})
			return
		}
		logger.Log.Info("[codespace] Found %d mappings for user %s", len(codespaces), userPubkey)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(codespaces)
		return
	}

	logger.Log.Info("[codespace] Querying all mappings")
	codespaces, err := ch.db.GetCodeSpaceMaps()
	if err != nil {
		logger.Log.Error("[codespace] error querying all codespace mappings: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to query codespace mappings"})
		return
	}
	logger.Log.Info("[codespace] Found %d total mappings", len(codespaces))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(codespaces)
}

func (ch *codeSpaceHandler) CreateCodeSpaceMap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[codespace] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	var codeSpace db.CodeSpaceMap
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		logger.Log.Error("[codespace] error reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read request body"})
		return
	}

	err = json.Unmarshal(body, &codeSpace)
	if err != nil {
		logger.Log.Error("[codespace] error unmarshaling request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request format"})
		return
	}

	if codeSpace.WorkspaceID == "" || codeSpace.CodeSpaceURL == "" || codeSpace.UserPubkey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "WorkspaceID, CodeSpaceURL, and UserPubkey are required"})
		return
	}

	createdCodeSpace, err := ch.db.CreateCodeSpaceMap(codeSpace)
	if err != nil {
		logger.Log.Error("[codespace] error creating codespace mapping: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create codespace mapping"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCodeSpace)
}

func (ch *codeSpaceHandler) UpdateCodeSpaceMap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[codespace] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID format"})
		return
	}

	var codeSpace db.CodeSpaceMap
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		logger.Log.Error("[codespace] error reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read request body"})
		return
	}

	err = json.Unmarshal(body, &codeSpace)
	if err != nil {
		logger.Log.Error("[codespace] error unmarshaling request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request format"})
		return
	}

	updates := make(map[string]interface{})
	if codeSpace.WorkspaceID != "" {
		updates["workspace_id"] = codeSpace.WorkspaceID
	}
	if codeSpace.CodeSpaceURL != "" {
		updates["code_space_url"] = codeSpace.CodeSpaceURL
	}
	if codeSpace.UserPubkey != "" {
		updates["user_pubkey"] = codeSpace.UserPubkey
	}
	// Also allow updating Username, GithubPat and BaseBranch, even if empty to clear them
	updates["username"] = codeSpace.Username
	updates["github_pat"] = codeSpace.GithubPat
	updates["base_branch"] = codeSpace.BaseBranch
	updates["pool_api_key"] = codeSpace.PoolAPIKey

	updatedCodeSpace, err := ch.db.UpdateCodeSpaceMap(id, updates)
	if err != nil {
		if err.Error() == "codespace mapping not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "CodeSpace mapping not found"})
			return
		}
		logger.Log.Error("[codespace] error updating codespace mapping: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update codespace mapping"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedCodeSpace)
}

func (ch *codeSpaceHandler) DeleteCodeSpaceMap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[codespace] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID format"})
		return
	}

	err = ch.db.DeleteCodeSpaceMap(id)
	if err != nil {
		if err.Error() == "codespace mapping not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "CodeSpace mapping not found"})
			return
		}
		logger.Log.Error("[codespace] error deleting codespace mapping: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete codespace mapping"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := DeleteResponse{
		Success: true,
		Message: "CodeSpace mapping deleted successfully",
	}

	json.NewEncoder(w).Encode(response)
}
