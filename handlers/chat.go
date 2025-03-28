package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/logger"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stakwork/sphinx-tribes/websocket"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/sse"
)

// ChatHandler handles chat-related requests
type ChatHandler struct {
	httpClient *http.Client
	db         db.Database
}

// ChatResponse is the response format for chat requests
type ChatResponse struct {
	Success   bool          `json:"success"`
	Message   string        `json:"message"`
	Data      interface{}   `json:"data,omitempty"`
	Artifacts []db.Artifact `json:"artifacts,omitempty"`
}

type HistoryChatResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type ChatHistoryResponse struct {
	Messages []db.ChatMessage `json:"messages"`
}

type StakworkChatPayload struct {
	Name           string                 `json:"name"`
	WorkflowID     int                    `json:"workflow_id"`
	WorkflowParams map[string]interface{} `json:"workflow_params"`
}

type ChatWebhookResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	ChatID  string           `json:"chat_id"`
	History []db.ChatMessage `json:"history"`
}

type FileResponse struct {
	Success    bool         `json:"success"`
	URL        string       `json:"url"`
	IsExisting bool         `json:"isExisting"`
	Asset      db.FileAsset `json:"asset"`
	UploadTime time.Time    `json:"uploadTime"`
}

type SendMessageRequest struct {
	ChatID         string `json:"chat_id"`
	Message        string `json:"message"`
	PDFURL         string `json:"pdf_url,omitempty"`
	ModelSelection string `json:"modelSelection,omitempty"`
	ContextTags    []struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"contextTags"`
	SourceWebsocketID string `json:"sourceWebsocketId"`
	WorkspaceUUID     string `json:"workspaceUUID"`
	Mode              string `json:"mode,omitempty"`
}

type BuildMessageRequest struct {
	Question string `json:"question"`
}

type ChatResponseRequest struct {
	Value struct {
		ChatID            string                `json:"chatId"`
		MessageID         string                `json:"messageId"`
		Response          string                `json:"response"`
		SourceWebsocketID string                `json:"sourceWebsocketId"`
		Artifacts         []ChatMessageArtifact `json:"artifacts,omitempty"`
	} `json:"value"`
}

type ChatMessageArtifact struct {
	ID      string          `json:"id"`
	Type    db.ArtifactType `json:"type"`
	Content interface{}     `json:"content"`
}

type ActionMessageRequest struct {
	ActionWebhook     string `json:"action_webhook"`
	ChatID            string `json:"chatId"`
	MessageID         string `json:"messageId"`
	Message           string `json:"message"`
	SourceWebsocketID string `json:"sourceWebsocketId"`
}

type ActionPayload struct {
	ChatID            string              `json:"chatId"`
	MessageID         string              `json:"messageId"`
	Message           string              `json:"message"`
	History           []map[string]string `json:"history"`
	CodeGraph         string              `json:"codeGraph,omitempty"`
	CodeGraphAlias    string              `json:"codeGraphAlias,omitempty"`
	SourceWebsocketID string              `json:"sourceWebsocketId"`
	WebhookURL        string              `json:"webhook_url"`
	CodeSpaceURL      string              `json:"codeSpaceURL"`
}

type CreateOrEditChatRequest struct {
	WorkspaceID string `json:"workspaceId"`
	Title       string `json:"title"`
}

type PaginationResponse struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalItems  int `json:"totalItems"`
	TotalPages  int `json:"totalPages"`
}

type ListFilesResponse struct {
	Success    bool               `json:"success" example:"true"`
	Data       []db.FileAsset     `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type ChatWorkflowRequest struct {
	WorkspaceID string `json:"workspaceId"`
	URL         string `json:"url"`
	StackworkID string `json:"stackworkId,omitempty"`
}

type ChatWorkflowResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message,omitempty"`
	Data    *db.ChatWorkflow `json:"data,omitempty"`
}

type StopSSERequest struct {
	SSEURL string `json:"sse_url"`
	ChatID string `json:"chatID"`
}

type ChatStatusRequest struct {
	ChatID  string `json:"chat_id"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type ChatStatusResponse struct {
	Success   bool                  `json:"success"`
	Message   string                `json:"message,omitempty"`
	Data      *db.ChatWorkflowStatus `json:"data,omitempty"`
	DataArray []db.ChatWorkflowStatus `json:"data_array,omitempty"`
}

type SSEMaintenanceRequest struct {
	StopAllClients bool  `json:"stop_all_clients"`
	CleanupLogs    bool  `json:"cleanup_logs"`
	LogMaxAgeHours int64 `json:"log_max_age_hours"`
}

type SSEMaintenanceResponse struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	ClientsStopped int    `json:"clients_stopped"`
	LogsRemoved    int64  `json:"logs_removed"`
}

type WebhookPayload struct {
	ProjectStatus string `json:"project_status"`
	Error         *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type ChatStatusWebhookResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewChatHandler(httpClient *http.Client, database db.Database) *ChatHandler {
	return &ChatHandler{
		httpClient: httpClient,
		db:         database,
	}
}

// CreateChat creates a new chat
//
//	@Summary		Create a new chat
//	@Description	Create a new chat with the given workspace ID and title
//	@Tags			Hive Chat
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			request	body		CreateOrEditChatRequest	true	"Chat creation request"
//	@Success		200		{object}	ChatResponse
//	@Failure		400		{object}	ChatResponse
//	@Failure		500		{object}	ChatResponse
//	@Router			/hivechat [post]
func (ch *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var request CreateOrEditChatRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if request.WorkspaceID == "" || request.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	chat := &db.Chat{
		ID:          xid.New().String(),
		WorkspaceID: request.WorkspaceID,
		Title:       request.Title,
		Status:      "active",
	}

	createdChat, err := ch.db.AddChat(chat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create chat: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Message: "Chat created successfully",
		Data:    createdChat,
	})
}

// UpdateChat updates an existing chat
//
//	@Summary		Update an existing chat
//	@Description	Update the title of an existing chat
//	@Tags			Hive Chat
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			chat_id	path		string					true	"Chat ID"
//	@Param			request	body		CreateOrEditChatRequest	true	"Chat update request"
//	@Success		200		{object}	ChatResponse
//	@Failure		400		{object}	ChatResponse
//	@Failure		404		{object}	ChatResponse
//	@Failure		500		{object}	ChatResponse
//	@Router			/hivechat/{chat_id} [put]
func (ch *ChatHandler) UpdateChat(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "chat_id")
	if chatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Chat ID is required",
		})
		return
	}

	var request CreateOrEditChatRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	existingChat, err := ch.db.GetChatByChatID(chatID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Chat not found",
		})
		return
	}

	updatedChat := existingChat
	updatedChat.Title = request.Title

	updatedChat, err = ch.db.UpdateChat(&updatedChat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to update chat: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Message: "Chat updated successfully",
		Data:    updatedChat,
	})
}

// ArchiveChat archives an existing chat
//
//	@Summary		Archive an existing chat
//	@Description	Archive a chat by changing its status to archived
//	@Tags			Hive Chat
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			chat_id	path		string	true	"Chat ID"
//	@Success		200		{object}	ChatResponse
//	@Failure		400		{object}	ChatResponse
//	@Failure		404		{object}	ChatResponse
//	@Failure		500		{object}	ChatResponse
//	@Router			/hivechat/{chat_id}/archive [put]
func (ch *ChatHandler) ArchiveChat(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "chat_id")
	if chatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Chat ID is required",
		})
		return
	}

	existingChat, err := ch.db.GetChatByChatID(chatID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Chat not found",
		})
		return
	}

	updatedChat := existingChat
	updatedChat.Status = db.ArchiveStatus

	updatedChat, err = ch.db.UpdateChat(&updatedChat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to archive chat: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Message: "Chat archived successfully",
		Data:    updatedChat,
	})
}

func buildArtifactList(artefacts []db.Artifact) []map[string]string {
	formattedArtefacts := []map[string]string{}
	for _, art := range artefacts {
		contentJSON, err := json.Marshal(art.Content)
		if err != nil {
			contentJSON = []byte("{}")
		}
		formattedArtefacts = append(formattedArtefacts, map[string]string{
			"artifactId": art.ID.String(),
			"content":    string(contentJSON),
		})
	}
	return formattedArtefacts
}

func buildVarsPayload(request SendMessageRequest, createdMessage *db.ChatMessage, messageHistory []map[string]string, context interface{}, user *db.Person, codeGraph *db.WorkspaceCodeGraph, codeSpace db.CodeSpaceMap, mode string) map[string]interface{} {
	vars := map[string]interface{}{
		"chatId":            request.ChatID,
		"messageId":         createdMessage.ID,
		"message":           request.Message,
		"history":           messageHistory,
		"contextTags":       context,
		"sourceWebsocketId": request.SourceWebsocketID,
		"webhook_url":       fmt.Sprintf("%s/hivechat/response", os.Getenv("HOST")),
		"alias":             user.OwnerAlias,
		"pdf_url":           request.PDFURL,
		"modelSelection":    request.ModelSelection,
		"workspaceId":       request.WorkspaceUUID,
	}

	if codeGraph != nil && codeGraph.Url != "" {
		url := strings.TrimSuffix(codeGraph.Url, "/")
		if !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}
		if mode == "Build" {
			vars["2b_base_url"] = url
			vars["secret"] = codeGraph.SecretAlias
		} else {
			vars["codeGraph"] = url
			vars["codeGraphAlias"] = codeGraph.SecretAlias
			vars["2b_base_url"] = url
			vars["secret"] = codeGraph.SecretAlias
		}
	}

	if codeSpace.CodeSpaceURL != "" {
		vars["codeSpaceURL"] = codeSpace.CodeSpaceURL
	}

	vars["query"] = request.Message

	return vars
}

// SendMessage sends a message in a chat
//
//	@Summary		Send a message in a chat
//	@Description	Send a message in a chat with the given details
//	@Tags			Hive Chat
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			request	body		SendMessageRequest	true	"Send message request"
//	@Success		200		{object}	ChatResponse
//	@Failure		400		{object}	ChatResponse
//	@Failure		401		{object}	ChatResponse
//	@Failure		500		{object}	ChatResponse
//	@Router			/hivechat/send [post]
func (ch *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := ch.db.GetPersonByPubkey(pubKeyFromAuth)

	if user.OwnerPubKey != pubKeyFromAuth {
		logger.Log.Info("Person not exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if request.WorkspaceUUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "workspaceUUID is required",
		})
		return
	}

	context, err := ch.db.GetProductBrief(request.WorkspaceUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Error retrieving product brief",
		})
		return
	}

	history, err := ch.db.GetChatMessagesForChatID(request.ChatID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to fetch chat history: %v", err),
		})
		return
	}

	start := 0
	if len(history) > 20 {
		start = len(history) - 20
	}
	recentHistory := history[start:]

	messageHistory := make([]map[string]string, len(recentHistory))
	for i, msg := range recentHistory {
		artefacts, err := ch.db.GetArtifactsByMessageID(msg.ID)
		if err != nil {
			artefacts = []db.Artifact{}
		}

		artifactJSON, err := json.Marshal(map[string]interface{}{"artifacts": buildArtifactList(artefacts)})
		if err != nil {
			artifactJSON = []byte(`{"artifacts": []}`)
		}

		messageHistory[i] = map[string]string{
			"role":    string(msg.Role),
			"content": msg.Message + "\nArtifacts: " + string(artifactJSON),
		}
	}

	message := &db.ChatMessage{
		ID:        xid.New().String(),
		ChatID:    request.ChatID,
		Message:   request.Message,
		PDFURL:    request.PDFURL,
		Role:      "user",
		Timestamp: time.Now(),
		Status:    "sending",
		Source:    "user",
	}

	createdMessage, err := ch.db.AddChatMessage(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to save message: %v", err),
		})
		return
	}

	var codeGraph *db.WorkspaceCodeGraph
	if workspaceID := request.WorkspaceUUID; workspaceID != "" {
		codeGraphResult, err := ch.db.GetCodeGraphByWorkspaceUuid(workspaceID)
		if err == nil {
			codeGraph = &codeGraphResult
		}
	}

	mode := "Chat"
	if request.Mode != "" {
		mode = request.Mode
	}

	var codeSpace db.CodeSpaceMap
	if workspaceID := request.WorkspaceUUID; workspaceID != "" {
		codeSpaceResult, err := ch.db.GetCodeSpaceMapByWorkspaceAndUser(workspaceID, pubKeyFromAuth)
		if err == nil {
			codeSpace = codeSpaceResult
		}
	}

	vars := buildVarsPayload(request, &createdMessage, messageHistory, context, &user, codeGraph, codeSpace, mode)

	stakworkPayload := StakworkChatPayload{
		Name:       "Hive Chat Processor",
		WorkflowID: 38842,
		WorkflowParams: map[string]interface{}{
			"set_var": map[string]interface{}{
				"attributes": map[string]interface{}{
					"vars": vars,
				},
			},
		},
	}

	if mode == "Build" {
		stakworkPayload.Name = "hive_autogen"
		stakworkPayload.WorkflowID = 43859
	}

	apiKeyEnv := "SWWFKEY"
	if mode == "Build" {
		//apiKeyEnv = "SWWFSWKEY"
		apiKeyEnv = "SWPR"
	}

	apiKey := os.Getenv(apiKeyEnv)
	if apiKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "environment variable is not set",
		})
		return
	}

	projectID, err := ch.sendToStakwork(stakworkPayload, apiKey)
	if err != nil {
		createdMessage.Status = "error"
		ch.db.UpdateChatMessage(&createdMessage)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to process message: %v", err),
		})
		return
	}

	projectMsg := websocket.TicketMessage{
		BroadcastType:   "direct",
		SourceSessionID: request.SourceWebsocketID,
		Message:         fmt.Sprintf("https://jobs.stakwork.com/admin/projects/%d", projectID),
		Action:          "swrun",
	}

	if err := websocket.WebsocketPool.SendTicketMessage(projectMsg); err != nil {
		log.Printf("Failed to send Stakwork project WebSocket message: %v", err)
	}

	wsMessage := websocket.TicketMessage{
		BroadcastType:   "direct",
		SourceSessionID: request.SourceWebsocketID,
		Message:         "Message sent",
		Action:          "process",
		ChatMessage:     createdMessage,
	}

	if err := websocket.WebsocketPool.SendTicketMessage(wsMessage); err != nil {
		log.Printf("Failed to send websocket message: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Message: "Message sent successfully",
		Data:    createdMessage,
	})
}

func (ch *ChatHandler) sendToStakwork(payload StakworkChatPayload, apiKey string) (int64, error) {

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("error marshaling payload: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.stakwork.com/api/v1/projects",
		bytes.NewBuffer(payloadJSON),
	)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Token token="+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ch.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("stakwork API error: %s", string(body))
	}

	var stakworkResp StakworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&stakworkResp); err != nil {
		return 0, fmt.Errorf("error decoding response: %v", err)
	}

	return stakworkResp.Data.ProjectID, nil
}

// GetChat retrieves chats for a workspace
//
//	@Summary		Retrieve chats for a workspace
//	@Description	Retrieve chats for a workspace with the given ID and status
//	@Tags			Hive Chat
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			workspace_id	query		string	true	"Workspace ID"
//	@Param			status			query		string	false	"Chat status"
//	@Success		200				{object}	ChatResponse
//	@Failure		400				{object}	ChatResponse
//	@Failure		500				{object}	ChatResponse
//	@Router			/hivechat [get]
func (ch *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.URL.Query().Get("workspace_id")
	chatStatus := r.URL.Query().Get("status")

	if workspaceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "workspace_id query parameter is required",
		})
		return
	}

	chats, err := ch.db.GetChatsForWorkspace(workspaceID, chatStatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to fetch chats: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Data:    chats,
	})
}

// GetChatHistory retrieves the history of a chat
//
//	@Summary		Retrieve chat history
//	@Description	Retrieve the history of a chat with the given ID
//	@Tags			Hive Chat
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path		string	true	"Chat ID"
//	@Success		200		{object}	HistoryChatResponse
//	@Failure		400		{object}	ChatResponse
//	@Failure		500		{object}	ChatResponse
//	@Router			/hivechat/history/{uuid} [get]
func (ch *ChatHandler) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "uuid")
	if chatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Chat ID is required",
		})
		return
	}

	messages, err := ch.db.GetChatMessagesForChatID(chatID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to fetch chat history: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HistoryChatResponse{
		Success: true,
		Data:    messages,
	})
}

// ProcessChatResponse processes a chat response
//
//	@Summary		Process a chat response
//	@Description	Process a chat response with the given details
//	@Tags			Hive Chat
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ChatResponseRequest	true	"Chat response request"
//	@Success		200		{object}	ChatResponse
//	@Failure		400		{object}	ChatResponse
//	@Failure		500		{object}	ChatResponse
//	@Router			/hivechat/response [post]
func (ch *ChatHandler) ProcessChatResponse(w http.ResponseWriter, r *http.Request) {
	var request ChatResponseRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if request.Value.ChatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "ChatID is required for message creation",
		})
		return
	}

	existingMessages, err := ch.db.GetChatMessagesForChatID(request.Value.ChatID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to check existing messages: %v", err),
		})
		return
	}

	for _, msg := range existingMessages {
		if msg.Role == "assistant" &&
			msg.Message == request.Value.Response &&
			time.Since(msg.Timestamp) < 5*time.Second {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ChatResponse{
				Success: true,
				Message: "Similar message already processed recently",
			})
			return
		}
	}

	message := &db.ChatMessage{
		ID:        xid.New().String(),
		ChatID:    request.Value.ChatID,
		Message:   request.Value.Response,
		Role:      "assistant",
		Timestamp: time.Now(),
		Status:    "sent",
		Source:    "agent",
	}

	createdMessage, err := ch.db.AddChatMessage(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to save response message: %v", err),
		})
		return
	}

	var artifacts []db.Artifact
	if len(request.Value.Artifacts) > 0 {
		for _, artifact := range request.Value.Artifacts {
			content := db.PropertyMap{}
			if contentMap, ok := artifact.Content.(map[string]interface{}); ok {
				content = db.PropertyMap(contentMap)
			}

			newArtifact := &db.Artifact{
				ID:        uuid.New(),
				MessageID: createdMessage.ID,
				Type:      artifact.Type,
				Content:   content,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			processedArtifact, err := ch.db.CreateArtifact(newArtifact)
			if err != nil {
				log.Printf("Error processing artifact: %v", err)
				continue
			}
			artifacts = append(artifacts, *processedArtifact)

			if artifact.Type == db.SSEArtifact {
				go HandleSSEConnectionArtifact(ch.db, artifact, request.Value.ChatID)
			}
		}
	}

	wsMessage := websocket.TicketMessage{
		BroadcastType:   "direct",
		SourceSessionID: request.Value.SourceWebsocketID,
		Message:         "Response received",
		Action:          "message",
		ChatMessage:     createdMessage,
		Artifacts:       artifacts,
	}

	if err := websocket.WebsocketPool.SendTicketMessage(wsMessage); err != nil {
		log.Printf("Failed to send websocket message: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success:   true,
		Message:   "Response processed successfully",
		Data:      createdMessage,
		Artifacts: artifacts,
	})
}

// UploadFile uploads a file to a chat
//
//	@Summary		Upload a file to a chat
//	@Description	Upload a file to a chat with the given details
//	@Tags			Hive Chat
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			file	formData	file	true	"File to upload"
//	@Success		200		{object}	FileResponse
//	@Failure		400		{object}	ChatResponse
//	@Failure		500		{object}	ChatResponse
//	@Router			/hivechat/upload [post]
func (ch *ChatHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "No file provided",
		})
		return
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	if !isAllowedFileType(mimeType) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "File type not allowed",
		})
		return
	}

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Failed to process file",
		})
		return
	}
	fileHash := hex.EncodeToString(h.Sum(nil))

	if existing, err := ch.db.GetFileAssetByHash(fileHash); err == nil {
		if err := ch.db.UpdateFileAssetReference(existing.ID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ChatResponse{
				Success: false,
				Message: "Failed to update reference time",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(FileResponse{
			Success:    true,
			URL:        existing.StoragePath,
			IsExisting: true,
			Asset:      *existing,
			UploadTime: existing.UploadTime,
		})
		return
	}

	file.Seek(0, 0)

	uploadFilename := uuid.New().String() + filepath.Ext(header.Filename)

	challenge := GetMemeChallenge()
	signer := SignChallenge(challenge.Challenge)
	mErr, mToken := GetMemeToken(challenge.Id, signer.Response.Sig)

	if mErr != "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Failed to get meme token",
		})
		return
	}

	err, uploadURL := UploadMemeImage(file, mToken.Token, header.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Failed to upload file",
		})
		return
	}

	asset := &db.FileAsset{
		OriginFilename: header.Filename,
		FileHash:       fileHash,
		UploadFilename: uploadFilename,
		FileSize:       header.Size,
		MimeType:       mimeType,
		Status:         db.ActiveFileStatus,
		UploadedBy:     r.Context().Value("pubkey").(string),
		StoragePath:    uploadURL,
		WorkspaceID:    r.URL.Query().Get("workspaceId"),
	}

	asset, err = ch.db.CreateFileAsset(asset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Failed to create asset record",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FileResponse{
		Success:    true,
		URL:        uploadURL,
		IsExisting: false,
		Asset:      *asset,
		UploadTime: asset.UploadTime,
	})
}

// GetFile retrieves a file from a chat
//
//	@Summary		Retrieve a file from a chat
//	@Description	Retrieve a file from a chat with the given ID
//	@Tags			Hive Chat
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path		string	true	"File ID"
//	@Success		200	{object}	FileResponse
//	@Failure		400	{object}	ChatResponse
//	@Failure		404	{object}	ChatResponse
//	@Router			/hivechat/file/{id} [get]
func (ch *ChatHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "File ID is required",
		})
		return
	}

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Invalid file ID",
		})
		return
	}

	asset, err := ch.db.GetFileAssetByID(uint(idUint))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "File not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FileResponse{
		Success:    true,
		URL:        asset.StoragePath,
		IsExisting: true,
		Asset:      *asset,
		UploadTime: asset.UploadTime,
	})
}

// ListFiles lists all files in a chat
//
//	@Summary		List all files in a chat
//	@Description	List all files in a chat with the given parameters
//	@Tags			Hive Chat
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			status		query		string	false	"File status"
//	@Param			mimeType	query		string	false	"File MIME type"
//	@Param			workspaceId	query		string	false	"Workspace ID"
//	@Param			page		query		int		false	"Page number"
//	@Param			pageSize	query		int		false	"Page size"
//	@Success		200			{object}	ListFilesResponse
//	@Failure		400			{object}	ChatResponse
//	@Failure		500			{object}	ChatResponse
//	@Router			/hivechat/file/all [get]
func (ch *ChatHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	var params db.ListFileAssetsParams
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Invalid parameters",
		})
		return
	}

	if status := r.URL.Query().Get("status"); status != "" {
		fileStatus := db.FileStatus(status)
		params.Status = &fileStatus
	}
	if mimeType := r.URL.Query().Get("mimeType"); mimeType != "" {
		params.MimeType = &mimeType
	}
	if workspaceID := r.URL.Query().Get("workspaceId"); workspaceID != "" {
		params.WorkspaceID = &workspaceID
	}

	params.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if params.Page <= 0 {
		params.Page = 1
	}
	params.PageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))
	if params.PageSize <= 0 {
		params.PageSize = 50
	}

	assets, total, err := ch.db.ListFileAssets(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Failed to retrieve files",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    assets,
		"pagination": map[string]interface{}{
			"currentPage": params.Page,
			"pageSize":    params.PageSize,
			"totalItems":  total,
			"totalPages":  int(math.Ceil(float64(total) / float64(params.PageSize))),
		},
	})
}

// DeleteFile deletes a file from a chat
//
//	@Summary		Delete a file from a chat
//	@Description	Delete a file from a chat with the given ID
//	@Tags			Hive Chat
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path		string	true	"File ID"
//	@Success		200	{object}	ChatResponse
//	@Failure		400	{object}	ChatResponse
//	@Failure		404	{object}	ChatResponse
//	@Failure		500	{object}	ChatResponse
//	@Router			/hivechat/file/{id} [delete]
func (ch *ChatHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "File ID is required",
		})
		return
	}

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Invalid file ID",
		})
		return
	}

	err = ch.db.DeleteFileAsset(uint(idUint))
	if err != nil {
		if strings.Contains(err.Error(), "file not found") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ChatResponse{
				Success: false,
				Message: "File not found",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Failed to delete file",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Message: "File deleted successfully",
	})
}

// SendBuildMessage sends a build message in a chat
//
//	@Summary		Send a build message in a chat
//	@Description	Send a build message in a chat with the given details
//	@Tags			Hive Chat
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			request	body		BuildMessageRequest	true	"Build message request"
//	@Success		200		{object}	ChatResponse
//	@Failure		400		{object}	ChatResponse
//	@Failure		500		{object}	ChatResponse
//	@Router			/hivechat/send/build [post]
func (ch *ChatHandler) SendBuildMessage(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req BuildMessageRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("SWWFSWKEY")
	if apiKey == "" {
		http.Error(w, "API key not set in environment", http.StatusInternalServerError)
		return
	}

	stakworkPayload := map[string]interface{}{
		"name":        "hive_autogen",
		"workflow_id": 43198,
		"workflow_params": map[string]interface{}{
			"set_var": map[string]interface{}{
				"attributes": map[string]interface{}{
					"vars": map[string]interface{}{
						"query": req.Question,
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(stakworkPayload)
	if err != nil {
		http.Error(w, "Error encoding payload", http.StatusInternalServerError)
		return
	}

	stakworkURL := "https://api.stakwork.com/api/v1/projects"
	reqStakwork, err := http.NewRequest("POST", stakworkURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Error creating Stakwork request", http.StatusInternalServerError)
		return
	}

	reqStakwork.Header.Set("Authorization", "Token token="+apiKey)
	reqStakwork.Header.Set("Content-Type", "application/json")

	resp, err := ch.httpClient.Do(reqStakwork)
	if err != nil {
		http.Error(w, "Error sending request to Stakwork", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading Stakwork response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func isAllowedFileType(mimeType string) bool {
	allowedTypes := map[string]bool{
		"application/pdf":  true,
		"image/jpeg":       true,
		"image/png":        true,
		"image/gif":        true,
		"text/plain":       true,
		"application/json": true,
	}
	return allowedTypes[mimeType]
}

func jsonErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (ch *ChatHandler) CreateArtefact(w http.ResponseWriter, r *http.Request) {
	var artifact db.Artifact
	if err := json.NewDecoder(r.Body).Decode(&artifact); err != nil {
		jsonErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdArtifact, err := ch.db.CreateArtifact(&artifact)
	if err != nil {
		jsonErrorResponse(w, fmt.Sprintf("Failed to create artifact: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdArtifact)
}

func (ch *ChatHandler) GetArtefactsByChatID(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "chatId")
	if chatID == "" {
		jsonErrorResponse(w, "Chat ID is required", http.StatusBadRequest)
		return
	}

	artifacts, err := ch.db.GetAllArtifactsByChatID(chatID)
	if err != nil {
		jsonErrorResponse(w, fmt.Sprintf("Failed to fetch artifacts: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(artifacts)
}

func (ch *ChatHandler) GetArtefactByID(w http.ResponseWriter, r *http.Request) {
	artifactIDStr := chi.URLParam(r, "artifactId")
	if artifactIDStr == "" {
		jsonErrorResponse(w, "Artifact ID is required", http.StatusBadRequest)
		return
	}

	artifactID, err := uuid.Parse(artifactIDStr)
	if err != nil {
		jsonErrorResponse(w, "Invalid artifact ID format", http.StatusBadRequest)
		return
	}

	artifact, err := ch.db.GetArtifactByID(artifactID)
	if err != nil {
		jsonErrorResponse(w, fmt.Sprintf("Failed to fetch artifact: %v", err), http.StatusInternalServerError)
		return
	}

	if artifact == nil {
		jsonErrorResponse(w, "Artifact not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(artifact)
}

func (ch *ChatHandler) GetArtefactsByMessageID(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "messageId")
	if messageID == "" {
		jsonErrorResponse(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	artifacts, err := ch.db.GetArtifactsByMessageID(messageID)
	if err != nil {
		jsonErrorResponse(w, fmt.Sprintf("Failed to fetch artifacts: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(artifacts)
}

func (ch *ChatHandler) UpdateArtefact(w http.ResponseWriter, r *http.Request) {
	var artifact db.Artifact
	if err := json.NewDecoder(r.Body).Decode(&artifact); err != nil {
		jsonErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedArtifact, err := ch.db.UpdateArtifact(&artifact)
	if err != nil {
		jsonErrorResponse(w, fmt.Sprintf("Failed to update artifact: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedArtifact)
}

func (ch *ChatHandler) DeleteArtefactByID(w http.ResponseWriter, r *http.Request) {
	artifactIDStr := chi.URLParam(r, "artifactId")
	if artifactIDStr == "" {
		jsonErrorResponse(w, "Artifact ID is required", http.StatusBadRequest)
		return
	}

	artifactID, err := uuid.Parse(artifactIDStr)
	if err != nil {
		jsonErrorResponse(w, "Invalid artifact ID format", http.StatusBadRequest)
		return
	}

	if err := ch.db.DeleteArtifactByID(artifactID); err != nil {
		jsonErrorResponse(w, fmt.Sprintf("Failed to delete artifact: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ch *ChatHandler) DeleteAllArtefactsByChatID(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "chatId")
	if chatID == "" {
		jsonErrorResponse(w, "Chat ID is required", http.StatusBadRequest)
		return
	}

	if err := ch.db.DeleteAllArtifactsByChatID(chatID); err != nil {
		jsonErrorResponse(w, fmt.Sprintf("Failed to delete artifacts: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ch *ChatHandler) SendActionMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	var request ActionMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if request.ActionWebhook == "" || request.ChatID == "" || request.MessageID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "action_webhook, chatId, and messageId are required",
		})
		return
	}

	history, err := ch.db.GetChatMessagesForChatID(request.ChatID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to get chat history: %v", err),
		})
		return
	}

	start := 0
	if len(history) > 20 {
		start = len(history) - 20
	}
	recentHistory := history[start:]

	messageHistory := make([]map[string]string, len(recentHistory))
	for i, msg := range recentHistory {
		artefacts, err := ch.db.GetArtifactsByMessageID(msg.ID)
		if err != nil {
			artefacts = []db.Artifact{}
		}

		artifactJSON, err := json.Marshal(map[string]interface{}{"artifacts": buildArtifactList(artefacts)})
		if err != nil {
			artifactJSON = []byte(`{"artifacts": []}`)
		}

		messageHistory[i] = map[string]string{
			"role":    string(msg.Role),
			"content": msg.Message + "\nArtifacts: " + string(artifactJSON),
		}
	}

	chat, err := ch.db.GetChatByChatID(request.ChatID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to get chat: %v", err),
		})
		return
	}

	codeGraph, err := ch.db.GetCodeGraphByWorkspaceUuid(chat.WorkspaceID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to get code graph: %v", err),
		})
		return
	}

	message := &db.ChatMessage{
		ID:        xid.New().String(),
		ChatID:    request.ChatID,
		Message:   request.Message,
		Role:      "user",
		Timestamp: time.Now(),
		Status:    "sending",
		Source:    "user",
	}

	createdMessage, err := ch.db.AddChatMessage(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to save message: %v", err),
		})
		return
	}

	var codeSpace db.CodeSpaceMap
	if workspaceID := chat.WorkspaceID; workspaceID != "" {
		codeSpaceResult, err := ch.db.GetCodeSpaceMapByWorkspaceAndUser(workspaceID, pubKeyFromAuth)
		if err == nil {
			codeSpace = codeSpaceResult
		}
	}

	payload := ActionPayload{
		ChatID:            request.ChatID,
		MessageID:         request.MessageID,
		Message:           request.Message,
		History:           messageHistory,
		SourceWebsocketID: request.SourceWebsocketID,
		WebhookURL:        fmt.Sprintf("%s/hivechat/response", os.Getenv("HOST")),
	}

	if codeGraph.Url != "" {
		url := strings.TrimSuffix(codeGraph.Url, "/")
		if !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}
		payload.CodeGraph = url
		payload.CodeGraphAlias = codeGraph.SecretAlias
	}

	if codeSpace.CodeSpaceURL != "" {
		payload.CodeSpaceURL = codeSpace.CodeSpaceURL
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to marshal payload: %v", err),
		})
		return
	}

	req, err := http.NewRequest("POST", request.ActionWebhook, bytes.NewBuffer(payloadBytes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create request: %v", err),
		})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := ch.httpClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to send message to action webhook: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Action webhook returned error: %s", string(body)),
		})
		return
	}

	wsMessage := websocket.TicketMessage{
		BroadcastType:   "direct",
		SourceSessionID: request.SourceWebsocketID,
		Message:         "Message sent",
		Action:          "process",
		ChatMessage:     createdMessage,
	}

	if err := websocket.WebsocketPool.SendTicketMessage(wsMessage); err != nil {
		log.Printf("Failed to send websocket message: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Message: "Message sent to action webhook successfully",
	})
}

func (ch *ChatHandler) CreateOrEditChatWorkflow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var request ChatWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatWorkflowResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if request.WorkspaceID == "" || request.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatWorkflowResponse{
			Success: false,
			Message: "WorkspaceID and URL are required",
		})
		return
	}

	workspace := ch.db.GetWorkspaceByUuid(request.WorkspaceID)
	if workspace.Uuid == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ChatWorkflowResponse{
			Success: false,
			Message: "Workspace not found",
		})
		return
	}

	workflow := &db.ChatWorkflow{
		WorkspaceID: request.WorkspaceID,
		URL:         request.URL,
		StackworkID: request.StackworkID,
	}

	result, err := ch.db.CreateOrEditChatWorkflow(workflow)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatWorkflowResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create/update chat workflow: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatWorkflowResponse{
		Success: true,
		Message: "Chat workflow created/updated successfully",
		Data:    result,
	})
}

func (ch *ChatHandler) GetChatWorkflow(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	if workspaceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatWorkflowResponse{
			Success: false,
			Message: "WorkspaceID is required",
		})
		return
	}

	workflow, err := ch.db.GetChatWorkflowByWorkspaceID(workspaceID)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(ChatWorkflowResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatWorkflowResponse{
		Success: true,
		Data:    workflow,
	})
}

func (ch *ChatHandler) DeleteChatWorkflow(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	if workspaceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatWorkflowResponse{
			Success: false,
			Message: "WorkspaceID is required",
		})
		return
	}

	if err := ch.db.DeleteChatWorkflow(workspaceID); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(ChatWorkflowResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatWorkflowResponse{
		Success: true,
		Message: "Chat workflow deleted successfully",
	})
}

func (ch *ChatHandler) StopSSEClient(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SSEURL string `json:"sse_url"`
		ChatID string `json:"chatID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if request.SSEURL == "" || request.ChatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "sse_url and chatID are required",
		})
		return
	}

	if sse.ClientRegistry.Unregister(request.SSEURL, request.ChatID) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: true,
			Message: "SSE client stopped successfully",
		})
		return
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: false,
		Message: "No active SSE client found for the provided sse_url and chatID",
	})
}

func (ch *ChatHandler) GetSSEMessagesByChatID(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "chat_id")
	if chatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Chat ID is required",
		})
		return
	}

	messages, err := ch.db.GetNewSSEMessageLogsByChatID(chatID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to retrieve SSE messages: %v", err),
		})
		return
	}

	if len(messages) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: true,
			Message: "No unsent SSE messages found",
			Data:    []interface{}{},
		})
		return
	}

	var messageIDs []uuid.UUID
	for _, msg := range messages {
		messageIDs = append(messageIDs, msg.ID)
	}

	err = ch.db.UpdateSSEMessageLogStatusBatch(messageIDs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to update message status: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Message: fmt.Sprintf("Retrieved %d SSE messages", len(messages)),
		Data:    messages,
	})
}

func (ch *ChatHandler) StartSSEClient(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SSEURL     string `json:"sse_url"`
		ChatID     string `json:"chatID"`
		WebhookURL string `json:"webhook_url,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if request.SSEURL == "" || request.ChatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Both sse_url and chatID are required",
		})
		return
	}

	if sse.ClientRegistry.HasClient(request.SSEURL, request.ChatID) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "An SSE client is already running for this URL and chat ID",
			Data: map[string]string{
				"chatID":  request.ChatID,
				"sse_url": request.SSEURL,
			},
		})
		return
	}

	webhookURL := request.WebhookURL
	if webhookURL == "" {
		webhookURL = fmt.Sprintf("%s/hivechat/response", os.Getenv("HOST"))
	}

	client := sse.NewClient(request.SSEURL, request.ChatID, webhookURL, ch.db)
	sse.ClientRegistry.Register(client)

	go client.Start()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Message: "SSE client started successfully",
		Data: map[string]string{
			"chatID":      request.ChatID,
			"sse_url":     request.SSEURL,
			"webhook_url": webhookURL,
		},
	})
}

func HandleSSEConnectionArtifact(database db.Database, artifact ChatMessageArtifact, chatID string) {

	content, ok := artifact.Content.(map[string]interface{})
	if !ok {
		log.Printf("Invalid SSE connection artifact content format")
		return
	}

	sseURL, ok := content["sse_url"].(string)
	if !ok || sseURL == "" {
		log.Printf("Missing or invalid sse_url in SSE connection artifact")
		return
	}

	webhookURL, ok := content["webhook_url"].(string)
	if !ok || webhookURL == "" {
		log.Printf("Missing or invalid webhook_url in SSE connection artifact")
		return
	}

	var delayMs int64 = 0
	switch v := content["delay"].(type) {
	case string:
		if parsedDelay, err := strconv.ParseInt(v, 10, 64); err == nil {
			delayMs = parsedDelay
		}
	case float64:
		delayMs = int64(v)
	}

	client := sse.NewClient(sseURL, chatID, webhookURL, database)
	sse.ClientRegistry.Register(client)

	go client.Start()

	log.Printf("Started SSE client for chatID %s connecting to %s", chatID, sseURL)

	if delayMs > 0 {
		log.Printf("Triggering webhook payload with delay: %dms", delayMs)
		go func() {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
			SendEventPayloadToWebhook(database, chatID, webhookURL, delayMs)
		}()
	}
}

func SendEventPayloadToWebhook(database db.Database, chatID string, webhookURL string, delayMs int64) {
	if delayMs > 0 {
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}

	unsentEvents, err := database.GetNewSSEMessageLogsByChatID(chatID)
	if err != nil {
		log.Printf("Error retrieving unsent events for chatID %s: %v", chatID, err)
		return
	}

	sseURL := ""
	if len(unsentEvents) > 0 {
		sseURL = unsentEvents[0].From
	}

	eventsList := make([]map[string]interface{}, len(unsentEvents))
	for i, event := range unsentEvents {
		eventsList[i] = map[string]interface{}{
			"event": event.Event,
		}
	}

	payload := map[string]interface{}{
		"chatID":  chatID,
		"events":  eventsList,
		"sse_url": sseURL,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling webhook payload: %v", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		log.Printf("Error sending events to webhook %s: %v", webhookURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Webhook returned non-success status %d: %s", resp.StatusCode, string(body))
		return
	}

	var eventIDs []uuid.UUID
	for _, event := range unsentEvents {
		eventIDs = append(eventIDs, event.ID)
	}

	err = database.UpdateSSEMessageLogStatusBatch(eventIDs)
	if err != nil {
		log.Printf("Error updating event status: %v", err)
		return
	}

	log.Printf("Successfully sent %d events for chatID %s to webhook %s", len(unsentEvents), chatID, webhookURL)
}

func (ch *ChatHandler) GetAllSSEMessagesByChatID(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "chat_id")
	if chatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Chat ID is required",
		})
		return
	}

	limit := 200
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}
	status := r.URL.Query().Get("status")

	messages, total, err := ch.db.GetSSEMessagesByChatID(chatID, limit, offset, status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to retrieve SSE messages: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Message: fmt.Sprintf("Retrieved %d SSE messages", len(messages)),
		Data: map[string]interface{}{
			"messages": messages,
			"total":    total,
			"limit":    limit,
			"offset":   offset,
		},
	})
}

// GetChatStatus retrieves all status entries for a specific chat
//
//	@Summary		Get all chat statuses
//	@Description	Retrieve all status entries for a specific chat
//	@Tags			Hive Chat
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			chat_id	path		string	true	"Chat ID"
//	@Success		200		{object}	ChatStatusResponse
//	@Failure		400		{object}	ChatStatusResponse
//	@Failure		500		{object}	ChatStatusResponse
//	@Router			/hivechat/status/{chat_id} [get]
func (ch *ChatHandler) GetAllChatStatus(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "chat_id")
	if chatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: "Chat ID is required",
		})
		return
	}

	statuses, err := ch.db.GetChatStatusByChatID(chatID)
	if err != nil {
		logger.Log.Error("Failed to get chat statuses: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to get chat statuses: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatStatusResponse{
		Success:   true,
		Message:   "Chat statuses retrieved successfully",
		DataArray: statuses,
	})
}

// GetLatestChatStatus retrieves the most recent status for a specific chat
//
//	@Summary		Get latest chat status
//	@Description	Retrieve the most recent status for a specific chat
//	@Tags			Hive Chat
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			chat_id	path		string	true	"Chat ID"
//	@Success		200		{object}	ChatStatusResponse
//	@Failure		400		{object}	ChatStatusResponse
//	@Failure		404		{object}	ChatStatusResponse
//	@Failure		500		{object}	ChatStatusResponse
//	@Router			/hivechat/status/{chat_id}/latest [get]
func (ch *ChatHandler) GetLatestChatStatus(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "chat_id")
	if chatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: "Chat ID is required",
		})
		return
	}

	status, err := ch.db.GetLatestChatStatusByChatID(chatID)
	if err != nil {
		if strings.Contains(err.Error(), "no chat status found") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ChatStatusResponse{
				Success: false,
				Message: "No status found for this chat",
			})
			return
		}
		logger.Log.Error("Failed to get latest chat status: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to get latest chat status: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatStatusResponse{
		Success: true,
		Message: "Latest chat status retrieved successfully",
		Data:    &status,
	})
}

// CreateChatStatus creates a new status entry for a chat
//
//	@Summary		Create chat status
//	@Description	Create a new status entry for a chat
//	@Tags			Hive Chat
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			request	body		ChatStatusRequest	true	"Chat status creation request"
//	@Success		201		{object}	ChatStatusResponse
//	@Failure		400		{object}	ChatStatusResponse
//	@Failure		404		{object}	ChatStatusResponse
//	@Failure		500		{object}	ChatStatusResponse
//	@Router			/hivechat/status [post]
func (ch *ChatHandler) CreateChatStatus(w http.ResponseWriter, r *http.Request) {
	var request ChatStatusRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if request.ChatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: "Chat ID is required",
		})
		return
	}

	if request.Status == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: "Status is required",
		})
		return
	}

	_, err = ch.db.GetChatByChatID(request.ChatID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: "Chat not found",
		})
		return
	}

	chatStatus := &db.ChatWorkflowStatus{
		ChatID:  request.ChatID,
		Status:  request.Status,
		Message: request.Message,
	}

	createdStatus, err := ch.db.AddChatStatus(chatStatus)
	if err != nil {
		logger.Log.Error("Failed to create chat status: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create chat status: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ChatStatusResponse{
		Success: true,
		Message: "Chat status created successfully",
		Data:    &createdStatus,
	})
}

// UpdateChatStatus updates an existing status entry
//
//	@Summary		Update chat status
//	@Description	Update an existing status entry
//	@Tags			Hive Chat
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path		string				true	"Status UUID"
//	@Param			request	body		ChatStatusRequest	true	"Chat status update request"
//	@Success		200		{object}	ChatStatusResponse
//	@Failure		400		{object}	ChatStatusResponse
//	@Failure		404		{object}	ChatStatusResponse
//	@Failure		500		{object}	ChatStatusResponse
//	@Router			/hivechat/status/{uuid} [put]
func (ch *ChatHandler) UpdateChatStatus(w http.ResponseWriter, r *http.Request) {
	statusUUID := chi.URLParam(r, "uuid")
	if statusUUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: "Status UUID is required",
		})
		return
	}

	var request ChatStatusRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	parsedUUID, err := uuid.Parse(statusUUID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: "Invalid UUID format",
		})
		return
	}

	chatStatus := &db.ChatWorkflowStatus{
		UUID:    parsedUUID,
		Status:  request.Status,
		Message: request.Message,
	}

	updatedStatus, err := ch.db.UpdateChatStatus(chatStatus)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ChatStatusResponse{
				Success: false,
				Message: "Chat status not found",
			})
			return
		}
		logger.Log.Error("Failed to update chat status: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to update chat status: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatStatusResponse{
		Success: true,
		Message: "Chat status updated successfully",
		Data:    &updatedStatus,
	})
}

// DeleteChatStatus deletes an existing status entry
//
//	@Summary		Delete chat status
//	@Description	Delete an existing status entry
//	@Tags			Hive Chat
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path		string	true	"Status UUID"
//	@Success		200		{object}	ChatStatusResponse
//	@Failure		400		{object}	ChatStatusResponse
//	@Failure		404		{object}	ChatStatusResponse
//	@Failure		500		{object}	ChatStatusResponse
//	@Router			/hivechat/status/{uuid} [delete]
func (ch *ChatHandler) DeleteChatStatus(w http.ResponseWriter, r *http.Request) {
	statusUUID := chi.URLParam(r, "uuid")
	if statusUUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: "Status UUID is required",
		})
		return
	}

	parsedUUID, err := uuid.Parse(statusUUID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: "Invalid UUID format",
		})
		return
	}

	err = ch.db.DeleteChatStatus(parsedUUID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ChatStatusResponse{
				Success: false,
				Message: "Chat status not found",
			})
			return
		}
		logger.Log.Error("Failed to delete chat status: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatStatusResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to delete chat status: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatStatusResponse{
		Success: true,
		Message: "Chat status deleted successfully",
	})
}

// HandleChatWebhook processes webhook requests from Stakwork
//
//	@Summary		Process chat webhook
//	@Description	Receives status updates from Stakwork workflow system
//	@Tags			Hive Chat
//	@Accept			json
//	@Produce		json
//	@Param			chat_id	path		string	true	"Chat ID"
//	@Param			payload	body		WebhookPayload	true	"Webhook payload"
//	@Success		200		{object}	ChatStatusWebhookResponse
//	@Failure		400		{object}	ChatStatusWebhookResponse
//	@Failure		404		{object}	ChatStatusWebhookResponse
//	@Failure		500		{object}	ChatStatusWebhookResponse
//	@Router			/hivechat/{chat_id}/update [post]
func (ch *ChatHandler) HandleChatWebhook(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "chat_id")
	if chatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatStatusWebhookResponse{
			Status:  "error",
			Message: "Chat ID is required",
		})
		return
	}

	_, err := ch.db.GetChatByChatID(chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Error("Chat not found for webhook: %s", chatID)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ChatStatusWebhookResponse{
				Status:  "error",
				Message: "Chat not found",
			})
			return
		}
		logger.Log.Error("Error fetching chat for webhook: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatStatusWebhookResponse{
			Status:  "error",
			Message: "Failed to verify chat",
		})
		return
	}

	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logger.Log.Error("Error parsing webhook payload: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatStatusWebhookResponse{
			Status:  "error",
			Message: "Invalid payload format",
		})
		return
	}

	payloadBytes, _ := json.Marshal(payload)
	logger.Log.Info("Received webhook for chat %s: %s", chatID, string(payloadBytes))

	status := ""
	message := ""

	if payload.ProjectStatus == "completed" {
		status = "success"
	} else if payload.ProjectStatus == "error" {
		status = "error"
		if payload.Error != nil {
			message = payload.Error.Message
		} else {
			message = "An error occurred during workflow execution"
		}
	} else {
		status = payload.ProjectStatus
	}

	chatStatus := &db.ChatWorkflowStatus{
		ChatID:  chatID,
		Status:  status,
		Message: message,
	}

	createdStatus, err := ch.db.AddChatStatus(chatStatus)
	if err != nil {
		logger.Log.Error("Failed to create chat status: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatStatusWebhookResponse{
			Status:  "error",
			Message: "Failed to process webhook",
		})
		return
	}

	logger.Log.Info("Created chat status for chat %s: %s - %s", 
		chatID, createdStatus.Status, createdStatus.Message)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatStatusWebhookResponse{
		Status:  "success",
		Message: "Webhook processed successfully",
	})
}

// SSEMaintenance performs maintenance on SSE connections and logs
//
//	@Summary		Perform SSE maintenance
//	@Description	Stop all SSE client connections and clean up old logs
//	@Tags			Hive Chat
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			request	body		SSEMaintenanceRequest	true	"Maintenance options"
//	@Success		200		{object}	SSEMaintenanceResponse
//	@Failure		400		{object}	ChatResponse
//	@Failure		500		{object}	ChatResponse
//	@Router			/hivechat/sse/maintenance [post]
func (ch *ChatHandler) SSEMaintenance(w http.ResponseWriter, r *http.Request) {
	var request SSEMaintenanceRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {

		request.StopAllClients = true
		request.CleanupLogs = true
		request.LogMaxAgeHours = 2
	}

	response := SSEMaintenanceResponse{
		Success: true,
	}

	if request.StopAllClients {
		response.ClientsStopped = sse.ClientRegistry.StopAllClients()
		logger.Log.Info("Stopped %d SSE clients during maintenance", response.ClientsStopped)
	}

	if request.CleanupLogs {
		if request.LogMaxAgeHours <= 0 {
			request.LogMaxAgeHours = 2
		}

		maxAge := time.Duration(request.LogMaxAgeHours) * time.Hour
		logsRemoved, err := ch.db.DeleteOldSSEMessageLogs(maxAge)
		if err != nil {
			logger.Log.Error("Error cleaning up SSE logs: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ChatResponse{
				Success: false,
				Message: fmt.Sprintf("Error cleaning up logs: %v", err),
			})
			return
		}

		response.LogsRemoved = logsRemoved
		logger.Log.Info("Removed %d SSE logs older than %v during maintenance", logsRemoved, maxAge)
	}

	response.Message = fmt.Sprintf("Maintenance completed: stopped %d clients, removed %d logs",
		response.ClientsStopped, response.LogsRemoved)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}