package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TicketMockHandler(t *testing.T, expectedStatus int, validateReq func(*http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if isProtectedTicketRoute(r.URL.Path, r.Method) {
			token := r.Header.Get("Authorization")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		if validateReq != nil && !validateReq(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if strings.Contains(r.URL.Path, "invalid") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if strings.Contains(r.URL.Path, "notfound") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(expectedStatus)
	}
}

func TestTicketRoutes(t *testing.T) {
	r := chi.NewRouter()
	ticketRouter := chi.NewRouter()

	ticketRouter.Get("/{uuid}", TicketMockHandler(t, http.StatusOK, nil))
	ticketRouter.Post("/review", TicketMockHandler(t, http.StatusOK, validateTicketReview))
	ticketRouter.Get("/feature/{feature_uuid}/phase/{phase_uuid}", TicketMockHandler(t, http.StatusOK, nil))
	ticketRouter.Post("/review/send", TicketMockHandler(t, http.StatusOK, validateReviewSend))
	ticketRouter.Post("/{uuid}", TicketMockHandler(t, http.StatusOK, validateTicketUpdate))
	ticketRouter.Post("/{ticket_group}/sequence", TicketMockHandler(t, http.StatusOK, validateSequence))
	ticketRouter.Post("/{ticket_uuid}/bounty", TicketMockHandler(t, http.StatusOK, validateTicketToBounty))
	ticketRouter.Delete("/{uuid}", TicketMockHandler(t, http.StatusOK, nil))

	r.Mount("/bounties/ticket", ticketRouter)

	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		token          string
		expectedStatus int
	}{
		{
			name:           "Get Ticket By UUID",
			method:         "GET",
			path:           "/bounties/ticket/123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Ticket By Invalid UUID",
			method:         "GET",
			path:           "/bounties/ticket/invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Get Non-existent Ticket",
			method:         "GET",
			path:           "/bounties/ticket/notfound",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Process Ticket Review",
			method:         "POST",
			path:           "/bounties/ticket/review",
			body:           map[string]interface{}{"review_data": "test review"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Process Invalid Ticket Review",
			method:         "POST",
			path:           "/bounties/ticket/review",
			body:           map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Get Tickets By Phase UUID - Unauthorized",
			method:         "GET",
			path:           "/bounties/ticket/feature/feature123/phase/phase123",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Get Tickets By Phase UUID - Authorized",
			method:         "GET",
			path:           "/bounties/ticket/feature/feature123/phase/phase123",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Post Ticket Review to Stakwork - Unauthorized",
			method:         "POST",
			path:           "/bounties/ticket/review/send",
			body:           map[string]interface{}{"data": "test review"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Post Ticket Review to Stakwork - Authorized",
			method:         "POST",
			path:           "/bounties/ticket/review/send",
			body:           map[string]interface{}{"data": "test review"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Update Ticket - Unauthorized",
			method:         "POST",
			path:           "/bounties/ticket/123",
			body:           map[string]interface{}{"status": "completed"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Update Ticket - Authorized",
			method:         "POST",
			path:           "/bounties/ticket/123",
			body:           map[string]interface{}{"status": "completed"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Update Ticket Sequence - Unauthorized",
			method:         "POST",
			path:           "/bounties/ticket/group123/sequence",
			body:           map[string]interface{}{"sequence": []int{1, 2, 3}},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Update Ticket Sequence - Authorized",
			method:         "POST",
			path:           "/bounties/ticket/group123/sequence",
			body:           map[string]interface{}{"sequence": []int{1, 2, 3}},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Convert Ticket to Bounty - Unauthorized",
			method:         "POST",
			path:           "/bounties/ticket/ticket123/bounty",
			body:           map[string]interface{}{"bounty_data": "test data"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Convert Ticket to Bounty - Authorized",
			method:         "POST",
			path:           "/bounties/ticket/ticket123/bounty",
			body:           map[string]interface{}{"bounty_data": "test data"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Delete Ticket - Unauthorized",
			method:         "DELETE",
			path:           "/bounties/ticket/123",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Delete Ticket - Authorized",
			method:         "DELETE",
			path:           "/bounties/ticket/123",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
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

func isProtectedTicketRoute(path string, method string) bool {

	if matched, _ := regexp.MatchString(`^/bounties/ticket/[a-zA-Z0-9-]+$`, path); matched && method == "GET" {
		return false
	}
	if path == "/bounties/ticket/review" && method == "POST" {
		return false
	}

	protectedPaths := map[string]bool{
		"/bounties/ticket/feature/":    true,
		"/bounties/ticket/review/send": true,
		"/bounties/ticket/sequence":    true,
		"/bounties/ticket/bounty":      true,
	}

	for protectedPath := range protectedPaths {
		if strings.Contains(path, protectedPath) {
			return true
		}
	}

	if matched, _ := regexp.MatchString(`^/bounties/ticket/[a-zA-Z0-9-]+$`, path); matched {
		return method == "DELETE" || method == "POST" || method == "PUT" || method == "PATCH"
	}

	return true
}

func validateTicketReview(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasReviewData := body["review_data"]
	return hasReviewData
}

func validateReviewSend(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasData := body["data"]
	return hasData
}

func validateTicketUpdate(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasStatus := body["status"]
	return hasStatus
}

func validateSequence(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasSequence := body["sequence"]
	return hasSequence
}

func validateTicketToBounty(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasBountyData := body["bounty_data"]
	return hasBountyData
}
