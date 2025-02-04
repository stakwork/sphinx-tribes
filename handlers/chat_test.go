package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Chat struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspaceId"`
	Title       string `json:"title"`
}

type ChatRes struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func TestUpdateChat(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	chatHandler := NewChatHandler(&http.Client{}, db.TestDB)

	t.Run("should successfully update chat when valid data is provided", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.UpdateChat)

		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: uuid.New().String(),
			Title:       "Old Title",
		}
		db.TestDB.AddChat(chat)

		requestBody := map[string]string{
			"workspaceId": chat.WorkspaceID,
			"title":       "New Title",
		}
		bodyBytes, _ := json.Marshal(requestBody)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("chat_id", chat.ID)
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodPut,
			"/hivechat/"+chat.ID,
			bytes.NewReader(bodyBytes),
		)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		assert.Equal(t, "Chat updated successfully", response.Message)
		responseData, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "Response data should be a map")
		assert.Equal(t, chat.ID, responseData["id"])
		assert.Equal(t, "New Title", responseData["title"])
	})

	t.Run("should return bad request when chat_id is missing", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.UpdateChat)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodPut,
			"/hivechat/",
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Chat ID is required", response.Message)
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.UpdateChat)

		invalidJson := []byte(`{"title": "New Title"`)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("chat_id", uuid.New().String())
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodPut,
			"/hivechat/"+uuid.New().String(),
			bytes.NewReader(invalidJson),
		)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request body", response.Message)
	})

	t.Run("should return not found when chat doesn't exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.UpdateChat)

		requestBody := map[string]string{
			"workspaceId": uuid.New().String(),
			"title":       "New Title",
		}
		bodyBytes, _ := json.Marshal(requestBody)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("chat_id", uuid.New().String())
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodPut,
			"/hivechat/"+uuid.New().String(),
			bytes.NewReader(bodyBytes),
		)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Chat not found", response.Message)
	})

	t.Run("should handle empty request body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.UpdateChat)

		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: uuid.New().String(),
			Title:       "Original Title",
		}
		db.TestDB.AddChat(chat)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("chat_id", chat.ID)
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodPut,
			"/hivechat/"+chat.ID,
			bytes.NewReader([]byte{}),
		)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request body", response.Message)
	})

	t.Run("should handle title with special characters", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.UpdateChat)

		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: uuid.New().String(),
			Title:       "Original Title",
		}
		db.TestDB.AddChat(chat)

		specialTitle := "!@#$%^&*()_+-=[]{}|;:'\",.<>?/"
		requestBody := map[string]string{
			"workspaceId": chat.WorkspaceID,
			"title":       specialTitle,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("chat_id", chat.ID)
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodPut,
			"/hivechat/"+chat.ID,
			bytes.NewReader(bodyBytes),
		)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		responseData, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, specialTitle, responseData["title"])
	})

	t.Run("should handle title with unicode characters", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.UpdateChat)

		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: uuid.New().String(),
			Title:       "Original Title",
		}
		db.TestDB.AddChat(chat)

		unicodeTitle := "测试标题 テストタイトル"
		requestBody := map[string]string{
			"workspaceId": chat.WorkspaceID,
			"title":       unicodeTitle,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("chat_id", chat.ID)
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodPut,
			"/hivechat/"+chat.ID,
			bytes.NewReader(bodyBytes),
		)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		responseData, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, unicodeTitle, responseData["title"])
	})

	t.Run("should handle concurrent update requests", func(t *testing.T) {
		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: uuid.New().String(),
			Title:       "Original Title",
		}
		db.TestDB.AddChat(chat)

		var wg sync.WaitGroup
		responses := make([]string, 5)

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(chatHandler.UpdateChat)

				requestBody := map[string]string{
					"workspaceId": chat.WorkspaceID,
					"title":       fmt.Sprintf("New Title %d", index),
				}
				bodyBytes, _ := json.Marshal(requestBody)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("chat_id", chat.ID)
				req, _ := http.NewRequestWithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
					http.MethodPut,
					"/hivechat/"+chat.ID,
					bytes.NewReader(bodyBytes),
				)

				handler.ServeHTTP(rr, req)

				var response ChatResponse
				json.NewDecoder(rr.Body).Decode(&response)
				if response.Success {
					responseData := response.Data.(map[string]interface{})
					responses[index] = responseData["title"].(string)
				}
			}(i)
		}
		wg.Wait()

		updatedChat, err := db.TestDB.GetChatByChatID(chat.ID)
		assert.NoError(t, err)
		assert.Contains(t, updatedChat.Title, "New Title")
	})

	t.Run("should handle very long title", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.UpdateChat)

		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: uuid.New().String(),
			Title:       "Original Title",
		}
		db.TestDB.AddChat(chat)

		longTitle := strings.Repeat("a", 1000)
		requestBody := map[string]string{
			"workspaceId": chat.WorkspaceID,
			"title":       longTitle,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("chat_id", chat.ID)
		req := httptest.NewRequest(
			http.MethodPut,
			"/hivechat/"+chat.ID,
			bytes.NewReader(bodyBytes),
		)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		responseData, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, longTitle, responseData["title"])
	})

	t.Run("should handle special characters in title", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.UpdateChat)

		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: uuid.New().String(),
			Title:       "Original Title",
		}
		db.TestDB.AddChat(chat)

		specialTitle := "!@#$%^&*()_+-=[]{}|;:'\",.<>?/\\`~"
		requestBody := map[string]string{
			"workspaceId": chat.WorkspaceID,
			"title":       specialTitle,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("chat_id", chat.ID)
		req := httptest.NewRequest(
			http.MethodPut,
			"/hivechat/"+chat.ID,
			bytes.NewReader(bodyBytes),
		)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		responseData, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, specialTitle, responseData["title"])
	})

	t.Run("should handle unicode characters in title", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.UpdateChat)

		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: uuid.New().String(),
			Title:       "Original Title",
		}
		db.TestDB.AddChat(chat)

		unicodeTitle := "测试标题 🌟 привет มาลอง"
		requestBody := map[string]string{
			"workspaceId": chat.WorkspaceID,
			"title":       unicodeTitle,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("chat_id", chat.ID)
		req := httptest.NewRequest(
			http.MethodPut,
			"/hivechat/"+chat.ID,
			bytes.NewReader(bodyBytes),
		)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		responseData, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, unicodeTitle, responseData["title"])
	})

	t.Run("should handle malformed JSON request", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.UpdateChat)

		malformedJSON := `{"workspaceId": "123", "title": "Test"`

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("chat_id", uuid.New().String())
		req := httptest.NewRequest(
			http.MethodPut,
			"/hivechat/"+uuid.New().String(),
			strings.NewReader(malformedJSON),
		)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request body", response.Message)
	})

	t.Run("should handle invalid content type", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.UpdateChat)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("chat_id", uuid.New().String())
		req := httptest.NewRequest(
			http.MethodPut,
			"/hivechat/"+uuid.New().String(),
			strings.NewReader("plain text body"),
		)
		req.Header.Set("Content-Type", "text/plain")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request body", response.Message)
	})
}

func TestCreateChat(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	chatHandler := NewChatHandler(&http.Client{}, db.TestDB)

	t.Run("should successfully get chats when valid workspace_id is provided", func(t *testing.T) {

		db.DeleteAllChats()

		chats := []*db.Chat{
			{
				ID:          uuid.New().String(),
				WorkspaceID: "workspace1",
				Title:       "Chat 1",
			},
			{
				ID:          uuid.New().String(),
				WorkspaceID: "workspace1",
				Title:       "Chat 2",
			},
		}
		for _, chat := range chats {
			db.TestDB.AddChat(chat)
		}

		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/hivechat?workspace_id=workspace1", nil)
		assert.NoError(t, err)

		handler := http.HandlerFunc(chatHandler.GetChat)
		handler.ServeHTTP(rr, req)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)

		responseChats, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Equal(t, 2, len(responseChats))

		firstChat := responseChats[0].(map[string]interface{})
		assert.NotEmpty(t, firstChat["id"])
		assert.Equal(t, "workspace1", firstChat["workspaceId"])
		assert.Contains(t, []string{"Chat 1", "Chat 2"}, firstChat["title"])
	})

	t.Run("should return empty array when no chats exist for workspace", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/hivechat?workspace_id=nonexistent", nil)
		assert.NoError(t, err)

		handler := http.HandlerFunc(chatHandler.GetChat)
		handler.ServeHTTP(rr, req)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)

		responseChats, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Empty(t, responseChats)
	})

	t.Run("should return bad request when workspace_id is missing", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/hivechat", nil)
		assert.NoError(t, err)

		handler := http.HandlerFunc(chatHandler.GetChat)
		handler.ServeHTTP(rr, req)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "workspace_id query parameter is required", response.Message)
	})

	t.Run("should return bad request when workspace_id is empty", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/hivechat?workspace_id=", nil)
		assert.NoError(t, err)

		handler := http.HandlerFunc(chatHandler.GetChat)
		handler.ServeHTTP(rr, req)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "workspace_id query parameter is required", response.Message)
	})

	t.Run("should handle special characters in workspace_id", func(t *testing.T) {
		workspaceID := "workspace-123-special"
		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: workspaceID,
			Title:       "Special Chat",
		}
		db.TestDB.AddChat(chat)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/hivechat?workspace_id=%s", url.QueryEscape(workspaceID)), nil)
		assert.NoError(t, err)

		handler := http.HandlerFunc(chatHandler.GetChat)
		handler.ServeHTTP(rr, req)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)

		responseChats, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Equal(t, 1, len(responseChats))
	})

	t.Run("should handle large number of chats", func(t *testing.T) {
		workspaceID := "workspace-large"

		for i := 0; i < 100; i++ {
			chat := &db.Chat{
				ID:          uuid.New().String(),
				WorkspaceID: workspaceID,
				Title:       fmt.Sprintf("Chat %d", i),
			}
			db.TestDB.AddChat(chat)
		}

		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/hivechat?workspace_id=%s", workspaceID), nil)
		assert.NoError(t, err)

		handler := http.HandlerFunc(chatHandler.GetChat)
		handler.ServeHTTP(rr, req)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)

		responseChats, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Equal(t, 100, len(responseChats))
	})

	t.Run("should successfully create chat when valid data is provided", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.CreateChat)

		requestBody := map[string]string{
			"workspaceId": "workspace123",
			"title":       "Test Chat",
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req, err := http.NewRequest(
			http.MethodPost,
			"/hivechat",
			bytes.NewReader(bodyBytes),
		)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		assert.Equal(t, "Chat created successfully", response.Message)

		responseData, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "Response data should be a map")
		assert.NotEmpty(t, responseData["id"])
		assert.Equal(t, "workspace123", responseData["workspaceId"])
		assert.Equal(t, "Test Chat", responseData["title"])
	})

	t.Run("should return bad request when request body is invalid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.CreateChat)

		invalidJson := []byte(`{"title": "Test Chat"`)
		req, err := http.NewRequest(
			http.MethodPost,
			"/hivechat",
			bytes.NewReader(invalidJson),
		)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request body", response.Message)
	})

	t.Run("should return bad request when required field workspaceId is missing", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.CreateChat)

		requestBody := map[string]string{
			"title": "Test Chat",
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req, err := http.NewRequest(
			http.MethodPost,
			"/hivechat",
			bytes.NewReader(bodyBytes),
		)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request body", response.Message)
	})

	t.Run("should return bad request when required field title is missing", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.CreateChat)

		requestBody := map[string]string{
			"workspaceId": "workspace123",
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req, err := http.NewRequest(
			http.MethodPost,
			"/hivechat",
			bytes.NewReader(bodyBytes),
		)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request body", response.Message)
	})

	t.Run("should handle empty strings in required fields", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.CreateChat)

		requestBody := map[string]string{
			"workspaceId": "",
			"title":       "",
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req, err := http.NewRequest(
			http.MethodPost,
			"/hivechat",
			bytes.NewReader(bodyBytes),
		)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request body", response.Message)
	})

	t.Run("should return bad request when non-string title is provided", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.CreateChat)

		rawJSON := []byte(`{
			"workspaceId": "workspace123",
			"title": 12345
		}`)

		req, err := http.NewRequest(
			http.MethodPost,
			"/hivechat",
			bytes.NewReader(rawJSON),
		)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request body", response.Message)
	})

	t.Run("should return bad request when non-string workspaceId is provided", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(chatHandler.CreateChat)

		rawJSON := []byte(`{
			"workspaceId": 12345,
			"title": "Test Chat"
		}`)

		req, err := http.NewRequest(
			http.MethodPost,
			"/hivechat",
			bytes.NewReader(rawJSON),
		)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var response ChatResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request body", response.Message)
	})
}

func TestProcessChatResponse(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	handler := NewChatHandler(&http.Client{}, db.TestDB)

	tests := []struct {
		name           string
		input          string
		mockDBResponse *db.ChatMessage
		mockDBError    error
		mockWSResponse error
		expectedStatus int
		expectedBody   ChatResponse
	}{
		{
			name: "Valid Input",
			input: `{
  			"value": {
  				"chatId": "validChatId",
  				"messageId": "validMessageId",
  				"response": "This is a response",
  				"sourceWebsocketId": "validWebsocketId"
  			}
  		}`,
			mockDBResponse: &db.ChatMessage{
				ID:        "generatedID",
				ChatID:    "validChatId",
				Message:   "This is a response",
				Role:      "assistant",
				Timestamp: time.Now(),
				Status:    "sent",
				Source:    "agent",
			},
			mockDBError:    nil,
			mockWSResponse: nil,
			expectedStatus: http.StatusOK,
			expectedBody: ChatResponse{
				Success: true,
				Message: "Response processed successfully",
				Data: &db.ChatMessage{
					ID:        "generatedID",
					ChatID:    "validChatId",
					Message:   "This is a response",
					Role:      "assistant",
					Timestamp: time.Now(),
					Status:    "sent",
					Source:    "agent",
				},
			},
		},
		{
			name: "Empty ChatID",
			input: `{
  			"value": {
  				"chatId": "",
  				"messageId": "validMessageId",
  				"response": "This is a response",
  				"sourceWebsocketId": "validWebsocketId"
  			}
  		}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody: ChatResponse{
				Success: false,
				Message: "ChatID is required for message creation",
			},
		},
		{
			name: "Empty Response",
			input: `{
  			"value": {
  				"chatId": "validChatId",
  				"messageId": "validMessageId",
  				"response": "",
  				"sourceWebsocketId": "validWebsocketId"
  			}
  		}`,
			mockDBResponse: &db.ChatMessage{
				ID:        "generatedID",
				ChatID:    "validChatId",
				Message:   "",
				Role:      "assistant",
				Timestamp: time.Now(),
				Status:    "sent",
				Source:    "agent",
			},
			mockDBError:    nil,
			mockWSResponse: nil,
			expectedStatus: http.StatusOK,
			expectedBody: ChatResponse{
				Success: true,
				Message: "Response processed successfully",
				Data: &db.ChatMessage{
					ID:        "generatedID",
					ChatID:    "validChatId",
					Message:   "",
					Role:      "assistant",
					Timestamp: time.Now(),
					Status:    "sent",
					Source:    "agent",
				},
			},
		},
		{
			name: "Invalid JSON Format",
			input: `{
  			"value": "invalidJson"
  		}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody: ChatResponse{
				Success: false,
				Message: "Invalid request body",
			},
		},
		{
			name: "Database Error",
			input: `{
  			"value": {
  				"messageId": "validMessageId",
  				"response": "This is a response",
  				"sourceWebsocketId": "validWebsocketId"
  			}
  		}`,
			mockDBResponse: nil,
			mockDBError:    errors.New("database error"),
			expectedStatus: http.StatusBadRequest,
			expectedBody: ChatResponse{
				Success: false,
				Message: "ChatID is required for message creation",
			},
		},
		{
			name: "WebSocket Error",
			input: `{
  			"value": {
  				"chatId": "validChatId",
  				"messageId": "validMessageId",
  				"response": "This is a response",
  				"sourceWebsocketId": "validWebsocketId"
  			}
  		}`,
			mockDBResponse: &db.ChatMessage{
				ID:        "generatedID",
				ChatID:    "validChatId",
				Message:   "This is a response",
				Role:      "assistant",
				Timestamp: time.Now(),
				Status:    "sent",
				Source:    "agent",
			},
			mockDBError:    nil,
			mockWSResponse: errors.New("websocket error"),
			expectedStatus: http.StatusOK,
			expectedBody: ChatResponse{
				Success: true,
				Message: "Response processed successfully",
				Data: &db.ChatMessage{
					ID:        "generatedID",
					ChatID:    "validChatId",
					Message:   "This is a response",
					Role:      "assistant",
					Timestamp: time.Now(),
					Status:    "sent",
					Source:    "agent",
				},
			},
		},
		{
			name: "Missing SourceWebsocketID",
			input: `{
  			"value": {
  				"chatId": "validChatId",
  				"messageId": "validMessageId",
  				"response": "This is a response"
  			}
  		}`,
			mockDBResponse: &db.ChatMessage{
				ID:        "generatedID",
				ChatID:    "validChatId",
				Message:   "This is a response",
				Role:      "assistant",
				Timestamp: time.Now(),
				Status:    "sent",
				Source:    "agent",
			},
			mockDBError:    nil,
			mockWSResponse: nil,
			expectedStatus: http.StatusOK,
			expectedBody: ChatResponse{
				Success: true,
				Message: "Response processed successfully",
				Data: &db.ChatMessage{
					ID:        "generatedID",
					ChatID:    "validChatId",
					Message:   "This is a response",
					Role:      "assistant",
					Timestamp: time.Now(),
					Status:    "sent",
					Source:    "agent",
				},
			},
		},
		{
			name: "All Fields Empty",
			input: `{
  			"value": {
  				"chatId": "",
  				"messageId": "",
  				"response": "",
  				"sourceWebsocketId": ""
  			}
  		}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody: ChatResponse{
				Success: false,
				Message: "ChatID is required for message creation",
			},
		},
		{
			name: "Missing MessageID",
			input: `{
  			"value": {
  				"chatId": "validChatId",
  				"response": "This is a response",
  				"sourceWebsocketId": "validWebsocketId"
  			}
  		}`,
			mockDBResponse: &db.ChatMessage{
				ID:        "generatedID",
				ChatID:    "validChatId",
				Message:   "This is a response",
				Role:      "assistant",
				Timestamp: time.Now(),
				Status:    "sent",
				Source:    "agent",
			},
			mockDBError:    nil,
			mockWSResponse: nil,
			expectedStatus: http.StatusOK,
			expectedBody: ChatResponse{
				Success: true,
				Message: "Response processed successfully",
				Data: &db.ChatMessage{
					ID:        "generatedID",
					ChatID:    "validChatId",
					Message:   "This is a response",
					Role:      "assistant",
					Timestamp: time.Now(),
					Status:    "sent",
					Source:    "agent",
				},
			},
		},
		{
			name: "Invalid request body",
			input: `{
  			"value": {
  				"chatId": 1,
  				"messageId": "validMessageId",
  				"response": "This is a response",
  				"sourceWebsocketId": "validWebsocketId"
  			}
  		}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody: ChatResponse{
				Success: false,
				Message: "Invalid request body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/response", bytes.NewBufferString(tt.input))
			w := httptest.NewRecorder()

			handler.ProcessChatResponse(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody ChatRes
			err := json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody.Success, responseBody.Success)
			assert.Equal(t, tt.expectedBody.Message, responseBody.Message)

			if tt.expectedBody.Data != nil {
				assert.NotNil(t, responseBody.Data)
				actualDataMap, ok := responseBody.Data.(map[string]interface{})
				assert.True(t, ok, "Response Data should be a map[string]interface{}")

				actualChatID, ok := actualDataMap["chatId"].(string)
				assert.True(t, ok, "ChatID in response should be a string")
				actualMessage, ok := actualDataMap["message"].(string)
				assert.True(t, ok, "Message in response should be a string")

				expectedData := tt.expectedBody.Data.(*db.ChatMessage)

				// Compare ChatID and Message
				assert.Equal(t, expectedData.ChatID, actualChatID)
				assert.Equal(t, expectedData.Message, actualMessage)
			} else {
				assert.Nil(t, responseBody.Data)
			}
		})
	}
}

func TestGetChatHistory(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	chatHandler := NewChatHandler(&http.Client{}, db.TestDB)

	t.Run("should successfully get chat history when valid chat_id is provided", func(t *testing.T) {

		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: "workspace1",
			Title:       "Test Chat",
		}
		db.TestDB.AddChat(chat)

		messages := []db.ChatMessage{
			{
				ID:        uuid.New().String(),
				ChatID:    chat.ID,
				Message:   "Message 1",
				Role:      "user",
				Timestamp: time.Now(),
			},
			{
				ID:        uuid.New().String(),
				ChatID:    chat.ID,
				Message:   "Message 2",
				Role:      "assistant",
				Timestamp: time.Now(),
			},
		}
		for _, msg := range messages {
			db.TestDB.AddChatMessage(&msg)
		}

		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", chat.ID)
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodGet,
			"/hivechat/history/"+chat.ID,
			nil,
		)
		assert.NoError(t, err)

		handler := http.HandlerFunc(chatHandler.GetChatHistory)
		handler.ServeHTTP(rr, req)

		var response HistoryChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		responseMessages, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Equal(t, 2, len(responseMessages))
	})

	t.Run("should return bad request when chat_id is missing", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodGet,
			"/hivechat/history/",
			nil,
		)
		assert.NoError(t, err)

		handler := http.HandlerFunc(chatHandler.GetChatHistory)
		handler.ServeHTTP(rr, req)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Chat ID is required", response.Message)
	})

	t.Run("should return empty array when chat has no messages", func(t *testing.T) {
		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: "workspace1",
			Title:       "Empty Chat",
		}
		db.TestDB.AddChat(chat)

		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", chat.ID)
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodGet,
			"/hivechat/history/"+chat.ID,
			nil,
		)
		assert.NoError(t, err)

		handler := http.HandlerFunc(chatHandler.GetChatHistory)
		handler.ServeHTTP(rr, req)

		var response HistoryChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		responseMessages, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Empty(t, responseMessages)
	})

	t.Run("should handle chat with large number of messages", func(t *testing.T) {
		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: "workspace1",
			Title:       "Large Chat",
		}
		db.TestDB.AddChat(chat)

		for i := 0; i < 100; i++ {
			message := &db.ChatMessage{
				ID:        uuid.New().String(),
				ChatID:    chat.ID,
				Message:   fmt.Sprintf("Message %d", i),
				Role:      "user",
				Timestamp: time.Now(),
			}
			db.TestDB.AddChatMessage(message)
		}

		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", chat.ID)
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodGet,
			"/hivechat/history/"+chat.ID,
			nil,
		)
		assert.NoError(t, err)

		handler := http.HandlerFunc(chatHandler.GetChatHistory)
		handler.ServeHTTP(rr, req)

		var response HistoryChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		responseMessages, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Equal(t, 100, len(responseMessages))
	})

	t.Run("Valid Chat ID with No Messages", func(t *testing.T) {
		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: "workspace1",
			Title:       "Empty Chat",
		}
		db.TestDB.AddChat(chat)

		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", chat.ID)
		req := httptest.NewRequest(http.MethodGet, "/hivechat/history/"+chat.ID, nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		chatHandler.GetChatHistory(rr, req)

		var response HistoryChatResponse
		json.NewDecoder(rr.Body).Decode(&response)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		messages, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Empty(t, messages)
	})

	t.Run("Chat ID with Special Characters", func(t *testing.T) {

		db.DeleteAllChats()
		db.DeleteAllChatMessages()

		specialChatID := "special!@#$%^&*()_+-=[]{}|;:,.<>?"
		chat := &db.Chat{
			ID:          specialChatID,
			WorkspaceID: "workspace1",
			Title:       "Special Chat",
		}
		db.TestDB.AddChat(chat)

		message := &db.ChatMessage{
			ID:        uuid.New().String(),
			ChatID:    specialChatID,
			Message:   "Special message",
			Role:      "user",
			Timestamp: time.Now(),
		}
		db.TestDB.AddChatMessage(message)

		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", specialChatID)
		req := httptest.NewRequest(
			http.MethodGet,
			"/hivechat/history/"+url.PathEscape(specialChatID),
			nil,
		)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		chatHandler.GetChatHistory(rr, req)

		var response HistoryChatResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)

		messages, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Equal(t, 1, len(messages))

		firstMessage := messages[0].(map[string]interface{})
		assert.Equal(t, specialChatID, firstMessage["chatId"])
		assert.Equal(t, "Special message", firstMessage["message"])
		assert.Equal(t, "user", firstMessage["role"])
	})

	t.Run("Non-Existent Chat ID", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", nonExistentID)
		req := httptest.NewRequest(http.MethodGet, "/hivechat/history/"+nonExistentID, nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		chatHandler.GetChatHistory(rr, req)

		var response HistoryChatResponse
		json.NewDecoder(rr.Body).Decode(&response)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		messages, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Empty(t, messages)
	})

	t.Run("Malformed Chat ID", func(t *testing.T) {
		malformedID := "malformed-chat-id-###"
		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", malformedID)
		req := httptest.NewRequest(http.MethodGet, "/hivechat/history/"+malformedID, nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		chatHandler.GetChatHistory(rr, req)

		var response HistoryChatResponse
		json.NewDecoder(rr.Body).Decode(&response)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		messages, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Empty(t, messages)
	})

}

func TestArchiveChat(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	chatHandler := NewChatHandler(&http.Client{}, db.TestDB)

	t.Run("should successfully archive chat when valid chat_id is provided", func(t *testing.T) {
		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: "workspace1",
			Title:       "Test Chat",
			Status:      "active",
		}
		db.TestDB.AddChat(chat)

		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("chat_id", chat.ID)
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodPut,
			"/hivechat/"+chat.ID+"/archive",
			nil,
		)
		assert.NoError(t, err)

		handler := http.HandlerFunc(chatHandler.ArchiveChat)
		handler.ServeHTTP(rr, req)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.True(t, response.Success)
		assert.Equal(t, "Chat archived successfully", response.Message)

		archivedChat, err := db.TestDB.GetChatByChatID(chat.ID)
		assert.NoError(t, err)
		assert.Equal(t, db.ArchiveStatus, archivedChat.Status)
	})

	t.Run("should return not found when chat doesn't exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("chat_id", uuid.New().String())
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodPut,
			"/hivechat/"+uuid.New().String()+"/archive",
			nil,
		)
		assert.NoError(t, err)

		handler := http.HandlerFunc(chatHandler.ArchiveChat)
		handler.ServeHTTP(rr, req)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Chat not found", response.Message)
	})

	t.Run("should return bad request when chat_id is missing", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodPut,
			"/hivechat//archive",
			nil,
		)
		assert.NoError(t, err)

		handler := http.HandlerFunc(chatHandler.ArchiveChat)
		handler.ServeHTTP(rr, req)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.False(t, response.Success)
		assert.Equal(t, "Chat ID is required", response.Message)
	})
}

func TestUploadFile(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()

	mockStorage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Error("Expected POST request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			t.Error("Failed to parse multipart form")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"success": true,
			"url":     "https://meme.sphinx.chat/public/test123.txt",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer mockStorage.Close()

	originalURL := os.Getenv("MEME_SERVER_URL")
	os.Setenv("MEME_SERVER_URL", mockStorage.URL)
	defer os.Setenv("MEME_SERVER_URL", originalURL)

	chatHandler := NewChatHandler(&http.Client{}, db.TestDB)

	createUploadRequest := func(filename, contentType string, content []byte, workspaceID string) (*http.Request, *httptest.ResponseRecorder) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
		h.Set("Content-Type", contentType)

		part, err := writer.CreatePart(h)
		require.NoError(t, err)

		_, err = part.Write(content)
		require.NoError(t, err)

		err = writer.Close()
		require.NoError(t, err)

		url := "/chat/upload"
		if workspaceID != "" {
			url = fmt.Sprintf("%s?workspaceId=%s", url, workspaceID)
		}

		req := httptest.NewRequest(http.MethodPost, url, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		ctx := context.WithValue(req.Context(), "pubkey", "test-pubkey-123")
		req = req.WithContext(ctx)

		return req, httptest.NewRecorder()
	}

	t.Run("should successfully upload new file", func(t *testing.T) {
		fileContent := []byte("test file content")

		req, rr := createUploadRequest(
			"test.txt",
			"text/plain",
			fileContent,
			"test-workspace-123",
		)

		chatHandler.UploadFile(rr, req)

		require.Equal(t, http.StatusOK, rr.Code, "Response body: %s", rr.Body.String())

		var response FileResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.False(t, response.IsExisting)
		assert.Equal(t, "https://meme.sphinx.chat/public/test123.txt", response.URL)

		assert.NotZero(t, response.Asset.ID)
		assert.Equal(t, "test.txt", response.Asset.OriginFilename)
		assert.Equal(t, "text/plain", response.Asset.MimeType)
		assert.Equal(t, db.ActiveFileStatus, response.Asset.Status)
		assert.Equal(t, "test-pubkey-123", response.Asset.UploadedBy)
		assert.Equal(t, "test-workspace-123", response.Asset.WorkspaceID)
		assert.NotEmpty(t, response.Asset.FileHash)
		assert.NotEmpty(t, response.Asset.UploadFilename)
		assert.Equal(t, int64(len(fileContent)), response.Asset.FileSize)
		assert.NotZero(t, response.Asset.UploadTime)
		assert.NotZero(t, response.Asset.LastReferenced)
		assert.Equal(t, "https://meme.sphinx.chat/public/test123.txt", response.Asset.StoragePath)

		storedAsset, err := db.TestDB.GetFileAssetByID(response.Asset.ID)
		require.NoError(t, err)
		assert.Equal(t, response.Asset.FileHash, storedAsset.FileHash)
	})

	t.Run("should reject file with unsupported mime type", func(t *testing.T) {
		req, rr := createUploadRequest(
			"test.xyz",
			"application/xyz",
			[]byte("test content"),
			"test-workspace-123",
		)

		chatHandler.UploadFile(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		var response ChatResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "File type not allowed", response.Message)
	})

	t.Run("should handle missing file in request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/chat/upload", nil)
		req.Header.Set("Content-Type", "multipart/form-data")
		ctx := context.WithValue(req.Context(), "pubkey", "test-pubkey-123")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		chatHandler.UploadFile(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		var response ChatResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "No file provided", response.Message)
	})

	t.Run("should handle storage service failure", func(t *testing.T) {
		failingStorage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Storage service error",
			})
		}))
		defer failingStorage.Close()

		failingHandler := NewChatHandler(&http.Client{}, db.TestDB)
		os.Setenv("MEME_SERVER_URL", failingStorage.URL)

		req, rr := createUploadRequest(
			"test.txt",
			"text/plain",
			[]byte("test content"),
			"test-workspace-123",
		)

		failingHandler.UploadFile(rr, req)

		require.Equal(t, http.StatusInternalServerError, rr.Code, "Response body: %s", rr.Body.String())
		var response ChatResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.False(t, response.Success)
		assert.Contains(t, response.Message, "Failed to upload file")

		os.Setenv("MEME_SERVER_URL", mockStorage.URL)
	})

	t.Run("should handle supported image types", func(t *testing.T) {
		imageTypes := []struct {
			ext         string
			contentType string
			content     []byte
		}{
			{"jpg", "image/jpeg", []byte("fake jpeg content")},
			{"png", "image/png", []byte("fake png content")},
			{"gif", "image/gif", []byte("fake gif content")},
		}

		for _, img := range imageTypes {
			t.Run(img.ext, func(t *testing.T) {

				imgStorage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"success": true,
						"url":     fmt.Sprintf("https://meme.sphinx.chat/public/test.%s", img.ext),
					})
				}))
				defer imgStorage.Close()

				os.Setenv("MEME_SERVER_URL", imgStorage.URL)
				imgHandler := NewChatHandler(&http.Client{}, db.TestDB)

				req, rr := createUploadRequest(
					fmt.Sprintf("test.%s", img.ext),
					img.contentType,
					img.content,
					"test-workspace-123",
				)

				imgHandler.UploadFile(rr, req)

				require.Equal(t, http.StatusOK, rr.Code, "Response body: %s", rr.Body.String())
				var response FileResponse
				err := json.NewDecoder(rr.Body).Decode(&response)
				require.NoError(t, err)
				assert.True(t, response.Success)
				assert.Equal(t, img.contentType, response.Asset.MimeType)
			})
		}

		os.Setenv("MEME_SERVER_URL", mockStorage.URL)
	})
}

func TestGetFile(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()
	chatHandler := NewChatHandler(&http.Client{}, db.TestDB)

	createTestFileAsset := func(t *testing.T, uploadFilename string) *db.FileAsset {
		if uploadFilename == "" {
			uploadFilename = fmt.Sprintf("test-upload-%d", time.Now().UnixNano())
		}

		asset := &db.FileAsset{
			OriginFilename: "test.txt",
			FileHash:       fmt.Sprintf("test-hash-%s", uploadFilename),
			UploadFilename: uploadFilename,
			FileSize:       100,
			MimeType:       "text/plain",
			StoragePath:    "https://meme.sphinx.chat/public/test123.txt",
			WorkspaceID:    "test-workspace-123",
			UploadedBy:     "test-pubkey-123",
			Status:         db.ActiveFileStatus,
		}

		createdAsset, err := db.TestDB.CreateFileAsset(asset)
		require.NoError(t, err)
		return createdAsset
	}

	createGetRequest := func(fileID string) (*http.Request, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, "/file/"+fileID, nil)
		ctx := context.WithValue(req.Context(), "pubkey", "test-pubkey-123")
		req = req.WithContext(ctx)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", fileID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		return req, httptest.NewRecorder()
	}

	t.Run("should successfully get file", func(t *testing.T) {
		asset := createTestFileAsset(t, "")

		req, rr := createGetRequest(fmt.Sprintf("%d", asset.ID))
		chatHandler.GetFile(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response FileResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.True(t, response.IsExisting)
		assert.Equal(t, asset.StoragePath, response.URL)
		assert.Equal(t, asset.ID, response.Asset.ID)
		assert.Equal(t, asset.OriginFilename, response.Asset.OriginFilename)
		assert.Equal(t, asset.MimeType, response.Asset.MimeType)
		assert.Equal(t, asset.Status, response.Asset.Status)
		assert.Equal(t, asset.WorkspaceID, response.Asset.WorkspaceID)
	})

	t.Run("should handle missing file ID", func(t *testing.T) {
		req, rr := createGetRequest("")
		chatHandler.GetFile(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)

		var response ChatResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "File ID is required", response.Message)
	})

	t.Run("should handle invalid file ID format", func(t *testing.T) {
		req, rr := createGetRequest("invalid-id")
		chatHandler.GetFile(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)

		var response ChatResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid file ID", response.Message)
	})

	t.Run("should handle non-existent file ID", func(t *testing.T) {
		req, rr := createGetRequest("999999")
		chatHandler.GetFile(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)

		var response ChatResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "File not found", response.Message)
	})

	t.Run("should handle deleted file", func(t *testing.T) {

		asset := createTestFileAsset(t, "test-upload-deleted")
		err := db.TestDB.DeleteFileAsset(asset.ID)
		require.NoError(t, err)

		req, rr := createGetRequest(fmt.Sprintf("%d", asset.ID))
		chatHandler.GetFile(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response FileResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.True(t, response.IsExisting)
		assert.Equal(t, db.DeletedFileStatus, response.Asset.Status)
	})

	t.Run("should handle file with different workspace", func(t *testing.T) {

		asset := &db.FileAsset{
			OriginFilename: "test.txt",
			FileHash:       "test-hash-diff-workspace",
			UploadFilename: fmt.Sprintf("test-upload-diff-workspace-%d", time.Now().UnixNano()),
			FileSize:       100,
			MimeType:       "text/plain",
			StoragePath:    "https://meme.sphinx.chat/public/test123.txt",
			WorkspaceID:    "different-workspace-123",
			UploadedBy:     "test-pubkey-123",
			Status:         db.ActiveFileStatus,
		}
		createdAsset, err := db.TestDB.CreateFileAsset(asset)
		require.NoError(t, err)

		req, rr := createGetRequest(fmt.Sprintf("%d", createdAsset.ID))
		chatHandler.GetFile(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response FileResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.True(t, response.IsExisting)
		assert.Equal(t, "different-workspace-123", response.Asset.WorkspaceID)
	})

	t.Run("should handle zero ID", func(t *testing.T) {
		req, rr := createGetRequest("0")
		chatHandler.GetFile(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code)

		var response ChatResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "File not found", response.Message)
	})
}

func TestListFiles(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()
	chatHandler := NewChatHandler(&http.Client{}, db.TestDB)

	createTestFileAsset := func(t *testing.T, opts map[string]string) *db.FileAsset {
		uploadFilename := fmt.Sprintf("test-upload-%d", time.Now().UnixNano())

		asset := &db.FileAsset{
			OriginFilename: opts["originFilename"],
			FileHash:       fmt.Sprintf("test-hash-%s", uploadFilename),
			UploadFilename: uploadFilename,
			FileSize:       100,
			MimeType:       opts["mimeType"],
			StoragePath:    "https://meme.sphinx.chat/public/test123.txt",
			WorkspaceID:    opts["workspaceId"],
			UploadedBy:     "test-pubkey-123",
			Status:         db.FileStatus(opts["status"]),
		}

		createdAsset, err := db.TestDB.CreateFileAsset(asset)
		require.NoError(t, err)
		return createdAsset
	}

	createListRequest := func(queryParams map[string]string) (*http.Request, *httptest.ResponseRecorder) {
		url := "/chat/files?"
		for key, value := range queryParams {
			url += fmt.Sprintf("%s=%s&", key, value)
		}
		req := httptest.NewRequest(http.MethodGet, url, nil)
		return req, httptest.NewRecorder()
	}

	t.Run("should list files with default pagination", func(t *testing.T) {

		for i := 0; i < 3; i++ {
			createTestFileAsset(t, map[string]string{
				"originFilename": fmt.Sprintf("test%d.txt", i),
				"mimeType":       "text/plain",
				"workspaceId":    "test-workspace-123",
				"status":         string(db.ActiveFileStatus),
			})
		}

		req, rr := createListRequest(map[string]string{})
		chatHandler.ListFiles(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.NotNil(t, response["data"])
		assert.NotNil(t, response["pagination"])

		pagination := response["pagination"].(map[string]interface{})
		assert.Equal(t, float64(1), pagination["currentPage"])
		assert.Equal(t, float64(50), pagination["pageSize"])
		assert.GreaterOrEqual(t, pagination["totalItems"].(float64), float64(3))
	})

	t.Run("should filter by status", func(t *testing.T) {

		createTestFileAsset(t, map[string]string{
			"originFilename": "active.txt",
			"mimeType":       "text/plain",
			"workspaceId":    "test-workspace-123",
			"status":         string(db.ActiveFileStatus),
		})
		deletedAsset := createTestFileAsset(t, map[string]string{
			"originFilename": "deleted.txt",
			"mimeType":       "text/plain",
			"workspaceId":    "test-workspace-123",
			"status":         string(db.ActiveFileStatus),
		})
		err := db.TestDB.DeleteFileAsset(deletedAsset.ID)
		require.NoError(t, err)

		req, rr := createListRequest(map[string]string{
			"status": string(db.DeletedFileStatus),
		})
		chatHandler.ListFiles(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		data := response["data"].([]interface{})
		for _, item := range data {
			asset := item.(map[string]interface{})
			assert.Equal(t, string(db.DeletedFileStatus), asset["status"])
		}
	})

	t.Run("should handle custom pagination", func(t *testing.T) {

		for i := 0; i < 15; i++ {
			createTestFileAsset(t, map[string]string{
				"originFilename": fmt.Sprintf("page%d.txt", i),
				"mimeType":       "text/plain",
				"workspaceId":    "test-workspace-123",
				"status":         string(db.ActiveFileStatus),
			})
		}

		req, rr := createListRequest(map[string]string{
			"page":     "2",
			"pageSize": "5",
		})
		chatHandler.ListFiles(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		pagination := response["pagination"].(map[string]interface{})
		assert.Equal(t, float64(2), pagination["currentPage"])
		assert.Equal(t, float64(5), pagination["pageSize"])
		data := response["data"].([]interface{})
		assert.LessOrEqual(t, len(data), 5)
	})

	t.Run("should handle invalid pagination parameters", func(t *testing.T) {
		req, rr := createListRequest(map[string]string{
			"page":     "-1",
			"pageSize": "0",
		})
		chatHandler.ListFiles(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		pagination := response["pagination"].(map[string]interface{})
		assert.Equal(t, float64(1), pagination["currentPage"])
		assert.Equal(t, float64(50), pagination["pageSize"])
	})

	t.Run("should handle empty result set", func(t *testing.T) {
		req, rr := createListRequest(map[string]string{
			"mimeType": "application/nonexistent",
		})
		chatHandler.ListFiles(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		data := response["data"].([]interface{})
		assert.Empty(t, data)
		pagination := response["pagination"].(map[string]interface{})
		assert.Equal(t, float64(0), pagination["totalItems"])
	})
}

func TestDeleteFile(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()
	chatHandler := NewChatHandler(&http.Client{}, db.TestDB)

	createTestFileAsset := func(t *testing.T) *db.FileAsset {
		uploadFilename := fmt.Sprintf("test-upload-%d", time.Now().UnixNano())
		asset := &db.FileAsset{
			OriginFilename: "test.txt",
			FileHash:       fmt.Sprintf("test-hash-%s", uploadFilename),
			UploadFilename: uploadFilename,
			FileSize:       100,
			MimeType:       "text/plain",
			StoragePath:    "https://meme.sphinx.chat/public/test123.txt",
			WorkspaceID:    "test-workspace-123",
			UploadedBy:     "test-pubkey-123",
			Status:         db.ActiveFileStatus,
		}

		createdAsset, err := db.TestDB.CreateFileAsset(asset)
		require.NoError(t, err)
		return createdAsset
	}

	createDeleteRequest := func(fileID string) (*http.Request, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodDelete, "/file/"+fileID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", fileID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		return req, httptest.NewRecorder()
	}

	t.Run("should successfully delete file", func(t *testing.T) {

		asset := createTestFileAsset(t)

		req, rr := createDeleteRequest(fmt.Sprintf("%d", asset.ID))
		chatHandler.DeleteFile(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response ChatResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "File deleted successfully", response.Message)

		deletedAsset, err := db.TestDB.GetFileAssetByID(asset.ID)
		require.NoError(t, err)
		assert.Equal(t, db.DeletedFileStatus, deletedAsset.Status)
		assert.NotNil(t, deletedAsset.DeletedAt)
	})

	t.Run("should handle missing file ID", func(t *testing.T) {
		req, rr := createDeleteRequest("")
		chatHandler.DeleteFile(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)

		var response ChatResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "File ID is required", response.Message)
	})

	t.Run("should handle invalid file ID format", func(t *testing.T) {
		req, rr := createDeleteRequest("invalid-id")
		chatHandler.DeleteFile(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)

		var response ChatResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid file ID", response.Message)
	})

	t.Run("should handle non-existent file ID", func(t *testing.T) {

		db.CleanTestData()

		nonExistentID := uint(999999)

		_, err := db.TestDB.GetFileAssetByID(nonExistentID)
		require.Error(t, err)

		req, rr := createDeleteRequest(fmt.Sprintf("%d", nonExistentID))
		chatHandler.DeleteFile(rr, req)

		require.Equal(t, http.StatusNotFound, rr.Code, "Should return 404 for non-existent file")

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "File not found", response.Message)
	})

	t.Run("should handle already deleted file", func(t *testing.T) {

		asset := createTestFileAsset(t)
		err := db.TestDB.DeleteFileAsset(asset.ID)
		require.NoError(t, err)

		req, rr := createDeleteRequest(fmt.Sprintf("%d", asset.ID))
		chatHandler.DeleteFile(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "File deleted successfully", response.Message)

		deletedAsset, err := db.TestDB.GetFileAssetByID(asset.ID)
		require.NoError(t, err)
		assert.Equal(t, db.DeletedFileStatus, deletedAsset.Status)
	})
}

func TestSendMessage(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()
	db.DeleteAllChatMessages()

	originalKey := os.Getenv("SWWFKEY")
	os.Setenv("SWWFKEY", "test-key")
	defer os.Setenv("SWWFKEY", originalKey)

	stakworkServer := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`{
                    "success": true,
                    "data": {
                        "project_id": 12345
                    }
                }`)),
				Header: make(http.Header),
			}, nil
		}),
	}

	websocket.WebsocketPool = &websocket.Pool{
		Clients: make(map[string]*websocket.ClientData),
	}

	chatHandler := NewChatHandler(stakworkServer, db.TestDB)

	t.Run("should successfully send message with PDF URL", func(t *testing.T) {

		person := db.Person{
			Uuid:        uuid.New().String(),
			OwnerAlias:  "test-alias",
			UniqueName:  "test-unique-name",
			OwnerPubKey: "test-pubkey",
			PriceToMeet: 0,
			Description: "test-description",
		}
		db.TestDB.CreateOrEditPerson(person)

		workspace := db.Workspace{
			Uuid:        uuid.New().String(),
			Name:        "test-workspace" + uuid.New().String(),
			OwnerPubKey: person.OwnerPubKey,
			Github:      "https://github.com/test",
			Website:     "https://www.testwebsite.com",
			Description: "test-description",
		}
		db.TestDB.CreateOrEditWorkspace(workspace)

		chatCreateReq := map[string]string{
			"workspaceId": workspace.Uuid,
			"title":       "Test Chat",
		}
		chatBodyBytes, _ := json.Marshal(chatCreateReq)
		chatReq := httptest.NewRequest(http.MethodPost, "/hivechat", bytes.NewReader(chatBodyBytes))
		chatRR := httptest.NewRecorder()
		chatHandler.CreateChat(chatRR, chatReq)

		var chatResponse ChatResponse
		err := json.NewDecoder(chatRR.Body).Decode(&chatResponse)
		require.NoError(t, err)
		require.True(t, chatResponse.Success)

		chatData := chatResponse.Data.(map[string]interface{})
		chatID := chatData["id"].(string)

		requestBody := SendMessageRequest{
			ChatID:  chatID,
			Message: "Test message with PDF",
			PDFURL:  "https://example.com/test.pdf",
			ContextTags: []struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			}{},
			SourceWebsocketID: "test-websocket-id",
			WorkspaceUUID:     workspace.Uuid,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(bodyBytes))
		rr := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), auth.ContextKey, "test-pubkey")
		req = req.WithContext(ctx)

		chatHandler.SendMessage(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Message sent successfully", response.Message)

		messages, err := db.TestDB.GetChatMessagesForChatID(chatID)
		require.NoError(t, err)
		assert.Equal(t, 1, len(messages))
		assert.Equal(t, "Test message with PDF", messages[0].Message)
		assert.Equal(t, db.UserRole, messages[0].Role)
		assert.Equal(t, db.SendingStatus, messages[0].Status)
	})

	t.Run("should handle unauthorized request", func(t *testing.T) {
		requestBody := SendMessageRequest{
			ChatID:            uuid.New().String(),
			Message:           "Test message",
			WorkspaceUUID:     uuid.New().String(),
			SourceWebsocketID: "test-websocket-id",
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(bodyBytes))
		rr := httptest.NewRecorder()

		chatHandler.SendMessage(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should handle invalid user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/send", nil)
		rr := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), auth.ContextKey, "non-existent-pubkey")
		req = req.WithContext(ctx)

		chatHandler.SendMessage(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should handle missing workspaceUUID", func(t *testing.T) {

		testUser := &db.Person{
			OwnerPubKey: "test-pubkey-2",
			OwnerAlias:  "test-user-2",
		}
		db.TestDB.CreateOrEditPerson(*testUser)

		requestBody := SendMessageRequest{
			ChatID:            uuid.New().String(),
			Message:           "Test message",
			SourceWebsocketID: "test-websocket-id",
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(bodyBytes))
		rr := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), auth.ContextKey, "test-pubkey-2")
		req = req.WithContext(ctx)

		chatHandler.SendMessage(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)

		var response ChatResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "workspaceUUID is required", response.Message)
	})

	t.Run("should handle invalid request body", func(t *testing.T) {

		testUser := &db.Person{
			OwnerPubKey: "test-pubkey-3",
			OwnerAlias:  "test-user-3",
		}
		db.TestDB.CreateOrEditPerson(*testUser)

		invalidJSON := []byte(`{"chat_id": "123", "message":`)
		req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(invalidJSON))
		rr := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), auth.ContextKey, "test-pubkey-3")
		req = req.WithContext(ctx)

		chatHandler.SendMessage(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)

		var response ChatResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid request body", response.Message)
	})

	t.Run("should successfully send message without PDF URL", func(t *testing.T) {

		testUser := &db.Person{
			Uuid:        uuid.New().String(),
			OwnerPubKey: "test-pubkey-4",
			OwnerAlias:  "test-user-4",
		}
		_, err := db.TestDB.CreateOrEditPerson(*testUser)
		require.NoError(t, err)

		workspace := db.Workspace{
			Uuid:        uuid.New().String(),
			Name:        "test-workspace" + uuid.New().String(),
			OwnerPubKey: testUser.OwnerPubKey,
			Github:      "https://github.com/test",
			Website:     "https://www.testwebsite.com",
			Description: "test-description",
		}
		_, err = db.TestDB.CreateOrEditWorkspace(workspace)
		require.NoError(t, err)

		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: workspace.Uuid,
			Title:       "Test Chat",
			Status:      db.ActiveStatus,
		}
		_, err = db.TestDB.AddChat(chat)
		require.NoError(t, err)

		requestBody := SendMessageRequest{
			ChatID:  chat.ID,
			Message: "Test message without PDF",
			ContextTags: []struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			}{},
			SourceWebsocketID: "test-websocket-id",
			WorkspaceUUID:     workspace.Uuid,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(bodyBytes))
		rr := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), auth.ContextKey, testUser.OwnerPubKey)
		req = req.WithContext(ctx)

		chatHandler.SendMessage(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Message sent successfully", response.Message)

		messages, err := db.TestDB.GetChatMessagesForChatID(chat.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, len(messages))
		assert.Equal(t, "Test message without PDF", messages[0].Message)
		assert.Equal(t, db.UserRole, messages[0].Role)
		assert.Equal(t, db.SendingStatus, messages[0].Status)
	})

	t.Run("should successfully send message with the model selection", func(t *testing.T) {
		person := db.Person{
			Uuid:        uuid.New().String(),
			OwnerAlias:  "test-alias-model",
			UniqueName:  "test-unique-name-model",
			OwnerPubKey: "test-pubkey-model",
			PriceToMeet: 0,
			Description: "test-description",
		}
		db.TestDB.CreateOrEditPerson(person)

		workspace := db.Workspace{
			Uuid:        uuid.New().String(),
			Name:        "test-workspace" + uuid.New().String(),
			OwnerPubKey: person.OwnerPubKey,
			Github:      "https://github.com/test",
			Website:     "https://www.testwebsite.com",
			Description: "test-description",
		}
		db.TestDB.CreateOrEditWorkspace(workspace)

		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: workspace.Uuid,
			Title:       "Test Chat with Model",
			Status:      db.ActiveStatus,
		}
		_, err := db.TestDB.AddChat(chat)
		require.NoError(t, err)

		requestBody := SendMessageRequest{
			ChatID:         chat.ID,
			Message:        "Test message with model selection",
			ModelSelection: "claude-3-sonnet",
			ContextTags: []struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			}{},
			SourceWebsocketID: "test-websocket-id",
			WorkspaceUUID:     workspace.Uuid,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(bodyBytes))
		rr := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
		req = req.WithContext(ctx)

		chatHandler.SendMessage(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Message sent successfully", response.Message)

		messages, err := db.TestDB.GetChatMessagesForChatID(chat.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, len(messages))
		assert.Equal(t, "Test message with model selection", messages[0].Message)
		assert.Equal(t, db.UserRole, messages[0].Role)
		assert.Equal(t, db.SendingStatus, messages[0].Status)
	})

	t.Run("should successfully send message with PDF URL and model selection", func(t *testing.T) {
		person := db.Person{
			Uuid:        uuid.New().String(),
			OwnerAlias:  "test-alias-pdf-model",
			UniqueName:  "test-unique-name-pdf-model",
			OwnerPubKey: "test-pubkey-pdf-model",
			PriceToMeet: 0,
			Description: "test-description",
		}
		db.TestDB.CreateOrEditPerson(person)

		workspace := db.Workspace{
			Uuid:        uuid.New().String(),
			Name:        "test-workspace" + uuid.New().String(),
			OwnerPubKey: person.OwnerPubKey,
			Github:      "https://github.com/test",
			Website:     "https://www.testwebsite.com",
			Description: "test-description",
		}
		db.TestDB.CreateOrEditWorkspace(workspace)

		chat := &db.Chat{
			ID:          uuid.New().String(),
			WorkspaceID: workspace.Uuid,
			Title:       "Test Chat PDF Model",
			Status:      db.ActiveStatus,
		}
		_, err := db.TestDB.AddChat(chat)
		require.NoError(t, err)

		requestBody := SendMessageRequest{
			ChatID:         chat.ID,
			Message:        "Test message with PDF and model",
			PDFURL:        "https://example.com/test.pdf",
			ModelSelection: "claude-3-opus",
			ContextTags: []struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			}{},
			SourceWebsocketID: "test-websocket-id",
			WorkspaceUUID:     workspace.Uuid,
		}
		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(bodyBytes))
		rr := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
		req = req.WithContext(ctx)

		chatHandler.SendMessage(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)

		var response ChatResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Message sent successfully", response.Message)

		messages, err := db.TestDB.GetChatMessagesForChatID(chat.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, len(messages))
		assert.Equal(t, "Test message with PDF and model", messages[0].Message)
		assert.Equal(t, db.UserRole, messages[0].Role)
		assert.Equal(t, db.SendingStatus, messages[0].Status)
	})

	t.Run("should successfully send message with code graph", func(t *testing.T) {
        person := db.Person{
            Uuid:        uuid.New().String(),
            OwnerAlias:  "test-alias-code-graph",
            UniqueName:  "test-unique-name-code-graph",
            OwnerPubKey: "test-pubkey-code-graph",
            PriceToMeet: 0,
            Description: "test-description",
        }
        db.TestDB.CreateOrEditPerson(person)

        workspace := db.Workspace{
            Uuid:        uuid.New().String(),
            Name:        "test-workspace" + uuid.New().String(),
            OwnerPubKey: person.OwnerPubKey,
            Github:      "https://github.com/test",
            Website:     "https://www.testwebsite.com",
            Description: "test-description",
        }
        db.TestDB.CreateOrEditWorkspace(workspace)

        codeGraph := db.WorkspaceCodeGraph{
            Uuid:           uuid.New().String(),
            WorkspaceUuid:  workspace.Uuid,
            Url:            "boltwall.swarm38.sphinx.chat/",
            CreatedBy:      person.OwnerPubKey,
            UpdatedBy:      person.OwnerPubKey,
        }
        _, err := db.TestDB.CreateOrEditCodeGraph(codeGraph)
        require.NoError(t, err)

        chat := &db.Chat{
            ID:          uuid.New().String(),
            WorkspaceID: workspace.Uuid,
            Title:       "Test Chat with Code Graph",
            Status:      db.ActiveStatus,
        }
        _, err = db.TestDB.AddChat(chat)
        require.NoError(t, err)

        requestBody := SendMessageRequest{
            ChatID:            chat.ID,
            Message:           "Test message with code graph",
            ContextTags:       []struct {
                Type string `json:"type"`
                ID   string `json:"id"`
            }{},
            SourceWebsocketID: "test-websocket-id",
            WorkspaceUUID:     workspace.Uuid,
        }
        bodyBytes, _ := json.Marshal(requestBody)

        req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(bodyBytes))
        rr := httptest.NewRecorder()

        ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
        req = req.WithContext(ctx)

        var capturedPayload StakworkChatPayload
        stakworkServer := &http.Client{
            Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
                body, _ := io.ReadAll(req.Body)
                json.Unmarshal(body, &capturedPayload)
                return &http.Response{
                    StatusCode: http.StatusOK,
                    Body: io.NopCloser(bytes.NewBufferString(`{
                        "success": true,
                        "data": {
                            "project_id": 12345
                        }
                    }`)),
                    Header: make(http.Header),
                }, nil
            }),
        }

        chatHandler := NewChatHandler(stakworkServer, db.TestDB)
        chatHandler.SendMessage(rr, req)

        require.Equal(t, http.StatusOK, rr.Code)

        vars, ok := capturedPayload.WorkflowParams["set_var"].(map[string]interface{})["attributes"].(map[string]interface{})["vars"].(map[string]interface{})
        require.True(t, ok, "Workflow params should contain vars")
        
        codeGraphUrl, exists := vars["codeGraph"].(string)
        require.True(t, exists, "codeGraph should exist in vars")
        assert.Equal(t, "https://boltwall.swarm38.sphinx.chat", codeGraphUrl)
    })

    t.Run("should send message without code graph", func(t *testing.T) {
        person := db.Person{
            Uuid:        uuid.New().String(),
            OwnerAlias:  "test-alias-no-code-graph",
            UniqueName:  "test-unique-name-no-code-graph",
            OwnerPubKey: "test-pubkey-no-code-graph",
            PriceToMeet: 0,
            Description: "test-description",
        }
        db.TestDB.CreateOrEditPerson(person)

        workspace := db.Workspace{
            Uuid:        uuid.New().String(),
            Name:        "test-workspace" + uuid.New().String(),
            OwnerPubKey: person.OwnerPubKey,
            Github:      "https://github.com/test",
            Website:     "https://www.testwebsite.com",
            Description: "test-description",
        }
        db.TestDB.CreateOrEditWorkspace(workspace)

        chat := &db.Chat{
            ID:          uuid.New().String(),
            WorkspaceID: workspace.Uuid,
            Title:       "Test Chat without Code Graph",
            Status:      db.ActiveStatus,
        }
        _, err := db.TestDB.AddChat(chat)
        require.NoError(t, err)

        requestBody := SendMessageRequest{
            ChatID:            chat.ID,
            Message:           "Test message without code graph",
            ContextTags:       []struct {
                Type string `json:"type"`
                ID   string `json:"id"`
            }{},
            SourceWebsocketID: "test-websocket-id",
            WorkspaceUUID:     workspace.Uuid,
        }
        bodyBytes, _ := json.Marshal(requestBody)

        req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(bodyBytes))
        rr := httptest.NewRecorder()

        ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
        req = req.WithContext(ctx)

        var capturedPayload StakworkChatPayload
        stakworkServer := &http.Client{
            Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
                body, _ := io.ReadAll(req.Body)
                json.Unmarshal(body, &capturedPayload)
                return &http.Response{
                    StatusCode: http.StatusOK,
                    Body: io.NopCloser(bytes.NewBufferString(`{
                        "success": true,
                        "data": {
                            "project_id": 12345
                        }
                    }`)),
                    Header: make(http.Header),
                }, nil
            }),
        }

        chatHandler := NewChatHandler(stakworkServer, db.TestDB)
        chatHandler.SendMessage(rr, req)

        require.Equal(t, http.StatusOK, rr.Code)

        vars, ok := capturedPayload.WorkflowParams["set_var"].(map[string]interface{})["attributes"].(map[string]interface{})["vars"].(map[string]interface{})
        require.True(t, ok, "Workflow params should contain vars")
        
        _, exists := vars["codeGraph"]
        assert.False(t, exists, "codeGraph should not exist in vars")
    })
}


type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
