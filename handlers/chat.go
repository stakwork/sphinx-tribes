package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
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
	ChatID            string `json:"chat_id"`
	Message           string `json:"message"`
	PDFURL            string `json:"pdf_url,omitempty"`
	ModelSelection    string `json:"modelSelection,omitempty"`
	ContextTags       []struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"contextTags"`
	SourceWebsocketID string `json:"sourceWebsocketId"`
	WorkspaceUUID     string `json:"workspaceUUID"`
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

func buildVarsPayload(request SendMessageRequest, createdMessage *db.ChatMessage, messageHistory []map[string]string, context interface{}, user *db.Person, codeGraph *db.WorkspaceCodeGraph) map[string]interface{} {
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
		vars["codeGraph"] = url
		vars["codeGraphAlias"] = codeGraph.SecretAlias
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
		codeGraphs, err := ch.db.GetCodeGraphsByWorkspaceUuid(workspaceID)
		if err == nil && len(codeGraphs) > 0 {
			codeGraph = &codeGraphs[0]
		}
	}

	vars := buildVarsPayload(request, &createdMessage, messageHistory, context, &user, codeGraph)

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

	projectID, err := ch.sendToStakwork(stakworkPayload)
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

func (ch *ChatHandler) sendToStakwork(payload StakworkChatPayload) (int64, error) {

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

	apiKey := os.Getenv("SWWFKEY")
	if apiKey == "" {
		return 0, fmt.Errorf("SWWFKEY environment variable not set")
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
	var request struct {
		Value struct {
			ChatID            string `json:"chatId"`
			MessageID         string `json:"messageId"`
			Response          string `json:"response"`
			SourceWebsocketID string `json:"sourceWebsocketId"`
		} `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	chatID := request.Value.ChatID
	response := request.Value.Response
	sourceWebsocketID := request.Value.SourceWebsocketID

	if chatID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "ChatID is required for message creation",
		})
		return
	}

	message := &db.ChatMessage{
		ID:        xid.New().String(),
		ChatID:    chatID,
		Message:   response,
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

	wsMessage := websocket.TicketMessage{
		BroadcastType:   "direct",
		SourceSessionID: sourceWebsocketID,
		Message:         "Response received",
		Action:          "message",
		ChatMessage:     createdMessage,
	}

	if err := websocket.WebsocketPool.SendTicketMessage(wsMessage); err != nil {
		log.Printf("Failed to send websocket message: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Message: "Response processed successfully",
		Data:    createdMessage,
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
