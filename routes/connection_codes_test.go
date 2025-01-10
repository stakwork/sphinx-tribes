package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func ConnectionCodeMockHandler(t *testing.T, expectedStatus int, validateReq func(*http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.ContentLength > 10*1024*1024 {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}

		if isProtectedConnectionCodeRoute(r.URL.Path, r.Method) {
			token := r.Header.Get("Authorization")
			if token != "Bearer valid-super-admin-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if len(bodyBytes) > 10*1024*1024 {
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		if validateReq != nil && !validateReq(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(expectedStatus)
	}
}

func TestConnectionCodesRoutes(t *testing.T) {
	r := chi.NewRouter()
	connectionRouter := chi.NewRouter()

	connectionRouter.Get("/", ConnectionCodeMockHandler(t, http.StatusOK, nil))

	connectionRouter.Post("/", ConnectionCodeMockHandler(t, http.StatusCreated, validateConnectionCodeRequest))

	r.Mount("/connectioncodes", connectionRouter)

	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		token          string
		expectedStatus int
	}{
		{
			name:           "Get Connection Code",
			method:         "GET",
			path:           "/connectioncodes/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Create Connection Code - Success",
			method:         "POST",
			path:           "/connectioncodes/",
			body:           map[string]interface{}{"code": "test-code"},
			token:          "valid-super-admin-token",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Create Connection Code - Unauthorized",
			method:         "POST",
			path:           "/connectioncodes/",
			body:           map[string]interface{}{"code": "test-code"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Create Connection Code - Invalid Token",
			method:         "POST",
			path:           "/connectioncodes/",
			body:           map[string]interface{}{"code": "test-code"},
			token:          "invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Create Connection Code - Invalid Request",
			method:         "POST",
			path:           "/connectioncodes/",
			body:           map[string]interface{}{},
			token:          "valid-super-admin-token",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Method",
			method:         "PUT",
			path:           "/connectioncodes/",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Create Connection Code - Empty Body",
			method:         "POST",
			path:           "/connectioncodes/",
			token:          "valid-super-admin-token",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Create Connection Code - Large Payload",
			method:         "POST",
			path:           "/connectioncodes/",
			body:           generateLargePayload(),
			token:          "valid-super-admin-token",
			expectedStatus: http.StatusRequestEntityTooLarge,
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

func isProtectedConnectionCodeRoute(path string, method string) bool {
	return method == "POST" && path == "/connectioncodes/"
}

func validateConnectionCodeRequest(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	code, hasCode := body["code"]
	if !hasCode {
		return false
	}
	codeStr, ok := code.(string)
	return ok && codeStr != ""
}

func generateLargePayload() map[string]interface{} {
	largeString := make([]byte, 11*1024*1024)
	for i := range largeString {
		largeString[i] = 'a'
	}
	return map[string]interface{}{
		"code": string(largeString),
	}
}
