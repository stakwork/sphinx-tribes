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
	"mime/multipart"
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
)

type ChatHandler struct {
	httpClient *http.Client
	db         db.Database
}

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
	ChatID           string `json:"chatId"`
	MessageID        string `json:"messageId"`
	Message          string `json:"message"`
	SourceWebsocketID string `json:"sourceWebsocketId"`
}

type ActionPayload struct {
	ChatID            string              `json:"chatId"`
	MessageID         string              `json:"messageId"`
	Message           string              `json:"message"`
	History           []db.ChatMessage    `json:"history"`
	CodeGraph         string              `json:"codeGraph,omitempty"`
	CodeGraphAlias    string              `json:"codeGraphAlias,omitempty"`
	SourceWebsocketID string              `json:"sourceWebsocketId"`
	WebhookURL        string              `json:"webhook_url"`
}

func NewChatHandler(httpClient *http.Client, database db.Database) *ChatHandler {
	return &ChatHandler{
		httpClient: httpClient,
		db:         database,
	}
}

func (ch *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var request struct {
		WorkspaceID string `json:"workspaceId"`
		Title       string `json:"title"`
	}

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

	var request struct {
		WorkspaceID string `json:"workspaceId"`
		Title       string `json:"title"`
	}

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

func buildVarsPayload(request SendMessageRequest, createdMessage *db.ChatMessage, messageHistory []map[string]string, context interface{}, user *db.Person, codeGraph *db.WorkspaceCodeGraph, mode string) map[string]interface{} {
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
		}
	}

	if mode == "Build" {
		vars["query"] = request.Message
	}

	return vars
}

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
		messageHistory[i] = map[string]string{
			"role":    string(msg.Role),
			"content": msg.Message,
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

	vars := buildVarsPayload(request, &createdMessage, messageHistory, context, &user, codeGraph, mode)

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
		stakworkPayload.WorkflowID = 43198
	}

	apiKeyEnv := "SWWFKEY"
	if mode == "Build" {
		apiKeyEnv = "SWWFSWKEY"
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
	uploadURL, err := uploadToStorage(file, uploadFilename)
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

func uploadToStorage(file multipart.File, filename string) (string, error) {

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, file); err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, &buf); err != nil {
		return "", fmt.Errorf("failed to copy to form file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	memeServerURL := os.Getenv("MEME_SERVER_URL")
	if memeServerURL == "" {
		memeServerURL = "https://meme.sphinx.chat" // TODO: CHANGE TO PROD
	}

	req, err := http.NewRequest("POST", memeServerURL+"/public", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("meme server returned status %d", resp.StatusCode)
	}

	var response struct {
		Success bool   `json:"success"`
		URL     string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return "", fmt.Errorf("meme server upload failed")
	}

	return response.URL, nil
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

	payload := ActionPayload{
		ChatID:            request.ChatID,
		MessageID:         request.MessageID,
		Message:           request.Message,
		History:           history,
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Message: "Message sent to action webhook successfully",
	})
}
