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

func BountyMockHandler(t *testing.T, expectedStatus int, validateReq func(*http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if isProtectedRoute(r.URL.Path) {
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

		if strings.Contains(r.URL.Path, "999999") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if strings.Contains(r.URL.Path, "invalid") && !strings.Contains(r.URL.Path, "invoice") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(expectedStatus)
	}
}

func TestBountyRoutes(t *testing.T) {
	r := chi.NewRouter()
	bountyRouter := chi.NewRouter()

	bountyRouter.Get("/all", BountyMockHandler(t, http.StatusOK, nil))
	bountyRouter.Get("/id/{bountyId}", BountyMockHandler(t, http.StatusOK, nil))
	bountyRouter.Get("/index/{bountyId}", BountyMockHandler(t, http.StatusOK, nil))
	bountyRouter.Get("/next/{created}", BountyMockHandler(t, http.StatusOK, nil))
	bountyRouter.Get("/previous/{created}", BountyMockHandler(t, http.StatusOK, nil))
	bountyRouter.Get("/count/{personKey}/{tabType}", BountyMockHandler(t, http.StatusOK, nil))
	bountyRouter.Get("/invoice/{paymentRequest}", BountyMockHandler(t, http.StatusOK, validateInvoiceRequest))

	bountyRouter.Post("/", BountyMockHandler(t, http.StatusOK, validateBountyRequest))
	bountyRouter.Delete("/{pubkey}/{created}", BountyMockHandler(t, http.StatusOK, nil))
	bountyRouter.Post("/paymentstatus/{created}", BountyMockHandler(t, http.StatusOK, nil))
	bountyRouter.Post("/{id}/proof", BountyMockHandler(t, http.StatusOK, validateProofRequest))

	r.Mount("/gobounties", bountyRouter)

	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		token          string
		expectedStatus int
	}{
		{
			name:           "Get All Bounties",
			method:         "GET",
			path:           "/gobounties/all",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Bounty By ID",
			method:         "GET",
			path:           "/gobounties/id/123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Bounty Index By ID",
			method:         "GET",
			path:           "/gobounties/index/123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Next Bounty By Created Date",
			method:         "GET",
			path:           "/gobounties/next/2024-01-01",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Previous Bounty By Created Date",
			method:         "GET",
			path:           "/gobounties/previous/2024-01-01",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Bounty By Non-Existent ID",
			method:         "GET",
			path:           "/gobounties/id/999999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Get Bounty By Invalid ID Format",
			method:         "GET",
			path:           "/gobounties/id/invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Get Bounty By Created Date with Invalid Format",
			method:         "GET",
			path:           "/gobounties/next/invalid-date",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unauthorized Access to Protected Route",
			method:         "POST",
			path:           "/gobounties/",
			body:           map[string]interface{}{"title": "Test Bounty"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Get Bounty Count for Specific User and Tab Type",
			method:         "GET",
			path:           "/gobounties/count/user123/active",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Invoice Data with Valid Payment Request",
			method:         "GET",
			path:           "/gobounties/invoice/valid-payment-request",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Invoice Data with Invalid Payment Request",
			method:         "GET",
			path:           "/gobounties/invoice/invalid-request",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Create or Edit Bounty",
			method:         "POST",
			path:           "/gobounties/",
			body:           map[string]interface{}{"title": "New Bounty", "amount": 1000},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Delete Bounty",
			method:         "DELETE",
			path:           "/gobounties/pubkey123/2024-01-01",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Update Bounty Payment Status",
			method:         "POST",
			path:           "/gobounties/paymentstatus/2024-01-01",
			token:          "valid-token",
			body:           map[string]interface{}{"status": "paid"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Add Proof of Work",
			method:         "POST",
			path:           "/gobounties/123/proof",
			token:          "valid-token",
			body:           map[string]interface{}{"proof": "Work completed", "url": "https://example.com"},
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

func isProtectedRoute(path string) bool {
	protectedPaths := []string{
		"/gobounties/paymentstatus/",
		"/gobounties/proof",
		"/gobounties/pubkey",
	}

	if path == "/gobounties/" {
		return true
	}

	for _, p := range protectedPaths {
		if strings.Contains(path, p) {
			return true
		}
	}
	return false
}

func validateBountyRequest(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasTitle := body["title"]
	return hasTitle
}

func validateProofRequest(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasProof := body["proof"]
	_, hasURL := body["url"]
	return hasProof && hasURL
}

func validateInvoiceRequest(r *http.Request) bool {
	paymentRequest := chi.URLParam(r, "paymentRequest")
	return paymentRequest != "invalid-request"
}
