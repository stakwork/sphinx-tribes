package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
	"github.com/stakwork/sphinx-tribes/utils"
)

type snippetHandler struct {
	httpClient HttpClient
	db         db.Database
}

func NewSnippetHandler(httpClient HttpClient, db db.Database) *snippetHandler {
	return &snippetHandler{
		httpClient: httpClient,
		db:         db,
	}
}

type SnippetRequest struct {
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
}

func (sh *snippetHandler) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[snippet] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	workspaceUUID := r.URL.Query().Get("workspace_uuid")
	if workspaceUUID == "" {
		http.Error(w, "workspace_uuid is required", http.StatusBadRequest)
		return
	}

	var req SnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Snippet == "" {
		http.Error(w, "Title and snippet are required", http.StatusBadRequest)
		return
	}

	snippet := &db.TextSnippet{
		WorkspaceUUID: workspaceUUID,
		Title:         req.Title,
		Snippet:       req.Snippet,
	}

	createdSnippet, err := sh.db.CreateSnippet(snippet)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to create snippet: %v", err))
		http.Error(w, "Failed to create snippet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSnippet)
}

func (sh *snippetHandler) GetSnippetsByWorkspace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[snippet] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	workspaceUUID := chi.URLParam(r, "workspace_uuid")
	if workspaceUUID == "" {
		http.Error(w, "workspace_uuid is required", http.StatusBadRequest)
		return
	}

	snippets, err := sh.db.GetSnippetsByWorkspace(workspaceUUID)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to fetch snippets: %v", err))
		http.Error(w, "Failed to fetch snippets", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(snippets)
}

func (sh *snippetHandler) GetSnippetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[snippet] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := utils.ConvertStringToUint(idStr)
	if err != nil {
		http.Error(w, "Invalid snippet ID", http.StatusBadRequest)
		return
	}

	snippet, err := sh.db.GetSnippetByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			http.Error(w, "Snippet not found", http.StatusNotFound)
			return
		}
		logger.Log.Error(fmt.Sprintf("Failed to fetch snippet: %v", err))
		http.Error(w, "Failed to fetch snippet", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(snippet)
}

func (sh *snippetHandler) UpdateSnippet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[snippet] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := utils.ConvertStringToUint(idStr)
	if err != nil {
		http.Error(w, "Invalid snippet ID", http.StatusBadRequest)
		return
	}

	var req SnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Snippet == "" {
		http.Error(w, "Title and snippet are required", http.StatusBadRequest)
		return
	}

	snippet := &db.TextSnippet{
		ID:      id,
		Title:   req.Title,
		Snippet: req.Snippet,
	}

	updatedSnippet, err := sh.db.UpdateSnippet(snippet)
	if err != nil {
		if err.Error() == "record not found" {
			http.Error(w, "Snippet not found", http.StatusNotFound)
			return
		}
		logger.Log.Error(fmt.Sprintf("Failed to update snippet: %v", err))
		http.Error(w, "Failed to update snippet", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedSnippet)
}

func (sh *snippetHandler) DeleteSnippet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[snippet] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := utils.ConvertStringToUint(idStr)
	if err != nil {
		http.Error(w, "Invalid snippet ID", http.StatusBadRequest)
		return
	}

	err = sh.db.DeleteSnippet(id)
	if err != nil {
		if err.Error() == "snippet not found" {
			http.Error(w, "Snippet not found", http.StatusNotFound)
			return
		}
		logger.Log.Error(fmt.Sprintf("Failed to delete snippet: %v", err))
		http.Error(w, "Failed to delete snippet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Snippet deleted successfully"})
}
