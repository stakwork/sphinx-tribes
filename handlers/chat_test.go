package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
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
