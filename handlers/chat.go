package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"net/http"
	"time"

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

	chat := &db.Chat{
		ID:          xid.New().String(),
		WorkspaceID: request.WorkspaceID,
		Title:       request.Title,
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

func (ch *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ChatID      string `json:"chatId"`
		Message     string `json:"message"`
		ContextTags []struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"contextTags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: "Invalid request body",
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
	}

	createdMessage, err := ch.db.AddChatMessage(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to send message: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Message: "Message sent successfully",
		Data:    createdMessage,
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ChatResponse{
		Success: true,
		Message: "Stubbed out - process chat response",
	})
}
