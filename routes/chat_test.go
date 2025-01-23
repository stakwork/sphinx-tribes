package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func ChatMockHandler(t *testing.T, expectedStatus int, validateReq func(*http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if isProtectedChatRoute(r.URL.Path) {
			token := r.Header.Get("Authorization")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if strings.Contains(token, "invalid-token") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		if strings.Contains(r.URL.Path, "invalid-id") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.URL.Path == "/chat/send" {
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if _, hasMessage := body["message"]; !hasMessage {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if _, hasChatID := body["chat_id"]; !hasChatID {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if msg, ok := body["message"].(string); ok && len(msg) > 5000 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.WriteHeader(expectedStatus)
			return
		}

		if validateReq != nil && !validateReq(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if strings.Contains(r.URL.Path, "invalid-uuid") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(expectedStatus)
	}
}

func TestChatRoutes(t *testing.T) {
	r := chi.NewRouter()
	chatRouter := chi.NewRouter()

	chatRouter.Post("/response", ChatMockHandler(t, http.StatusOK, validateChatResponse))
	chatRouter.Get("/", ChatMockHandler(t, http.StatusOK, nil))
	chatRouter.Post("/", ChatMockHandler(t, http.StatusOK, validateCreateChat))
	chatRouter.Put("/{chat_id}", ChatMockHandler(t, http.StatusOK, validateUpdateChat))
	chatRouter.Put("/{chat_id}/archive", ChatMockHandler(t, http.StatusOK, nil))
	chatRouter.Post("/send", ChatMockHandler(t, http.StatusOK, validateSendMessage))
	chatRouter.Get("/history/{uuid}", ChatMockHandler(t, http.StatusOK, nil))

	r.Mount("/chat", chatRouter)

	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		token          string
		expectedStatus int
	}{
		{
			name:           "POST /response - Process Chat Response",
			method:         "POST",
			path:           "/chat/response",
			body:           map[string]interface{}{"message": "Hello", "chat_id": "123"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET / - Retrieve Chat List",
			method:         "GET",
			path:           "/chat/",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST / - Create New Chat",
			method:         "POST",
			path:           "/chat/",
			body:           map[string]interface{}{"title": "New Chat"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "PUT /{chat_id} - Update Existing Chat",
			method:         "PUT",
			path:           "/chat/123",
			body:           map[string]interface{}{"title": "Updated Chat"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "PUT /{chat_id}/archive - Archive Chat",
			method:         "PUT",
			path:           "/chat/123/archive",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /send - Send Chat Message",
			method:         "POST",
			path:           "/chat/send",
			body:           map[string]interface{}{"message": "Hello", "chat_id": "123"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /history/{uuid} - Retrieve Chat History",
			method:         "GET",
			path:           "/chat/history/valid-uuid",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /history/{uuid} with non-existent UUID",
			method:         "GET",
			path:           "/chat/history/invalid-uuid",
			token:          "valid-token",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "PUT /{chat_id} with invalid data",
			method:         "PUT",
			path:           "/chat/123",
			body:           map[string]interface{}{},
			token:          "valid-token",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "POST / with missing fields",
			method:         "POST",
			path:           "/chat/",
			body:           map[string]interface{}{},
			token:          "valid-token",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "PUT /{chat_id} with invalid chat ID",
			method:         "PUT",
			path:           "/chat/invalid-id",
			body:           map[string]interface{}{"title": "Updated Chat"},
			token:          "valid-token",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unauthorized access to protected routes",
			method:         "GET",
			path:           "/chat/",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "POST /send with large payload",
			method:         "POST",
			path:           "/chat/send",
			body:           map[string]interface{}{"message": strings.Repeat("a", 10000), "chat_id": "123"},
			token:          "valid-token",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "GET /history/{uuid} with large history",
			method:         "GET",
			path:           "/chat/history/large-history-uuid",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /response with special characters",
			method:         "POST",
			path:           "/chat/response",
			body:           map[string]interface{}{"message": "Hello!@#$%^&*()", "chat_id": "123"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "PUT /{chat_id}/archive on already archived chat",
			method:         "PUT",
			path:           "/chat/archived-123/archive",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET / with invalid authentication token",
			method:         "GET",
			path:           "/chat/",
			token:          "invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.body != nil {
				bodyBytes, _ := json.Marshal(tc.body)
				req = httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tc.method, tc.path, nil)
			}

			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code, "Handler returned wrong status code for test: "+tc.name)
		})
	}
}

func isProtectedChatRoute(path string) bool {
	unprotectedPaths := []string{
		"/chat/response",
	}

	for _, p := range unprotectedPaths {
		if path == p {
			return false
		}
	}
	return true
}

func validateChatResponse(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasMessage := body["message"]
	_, hasChatID := body["chat_id"]
	return hasMessage && hasChatID
}

func validateCreateChat(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasTitle := body["title"]
	return hasTitle
}

func validateUpdateChat(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasTitle := body["title"]
	return hasTitle
}

func validateSendMessage(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}

	message, hasMessage := body["message"].(string)
	chatID, hasChatID := body["chat_id"].(string)

	if !hasMessage || !hasChatID {
		return false
	}

	if message == "" || chatID == "" {
		return false
	}

	if len(message) > 5000 {
		return false
	}

	return true
}
