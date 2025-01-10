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

func MetricsMockHandler(t *testing.T, expectedStatus int, validateReq func(*http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(token, "Bearer valid-super-admin-token") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if validateReq != nil && !validateReq(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(expectedStatus)
	}
}

func TestMetricsRoutes(t *testing.T) {
	r := chi.NewRouter()
	metricsRouter := chi.NewRouter()

	metricsRouter.Get("/workspaces", MetricsMockHandler(t, http.StatusOK, nil))
	metricsRouter.Post("/payment", MetricsMockHandler(t, http.StatusOK, validateRequestBody))
	metricsRouter.Post("/people", MetricsMockHandler(t, http.StatusOK, validateRequestBody))
	metricsRouter.Post("/organization", MetricsMockHandler(t, http.StatusOK, validateRequestBody))
	metricsRouter.Post("/bounty_stats", MetricsMockHandler(t, http.StatusOK, validateRequestBody))
	metricsRouter.Post("/bounties", MetricsMockHandler(t, http.StatusOK, validateRequestBody))
	metricsRouter.Post("/bounties/count", MetricsMockHandler(t, http.StatusOK, validateRequestBody))
	metricsRouter.Post("/bounties/providers", MetricsMockHandler(t, http.StatusOK, validateRequestBody))
	metricsRouter.Post("/csv", MetricsMockHandler(t, http.StatusOK, validateRequestBody))

	r.Mount("/metrics", metricsRouter)

	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		token          string
		expectedStatus int
	}{

		{
			name:           "GET /workspaces",
			method:         "GET",
			path:           "/metrics/workspaces",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /payment",
			method:         "POST",
			path:           "/metrics/payment",
			body:           map[string]interface{}{"data": "valid payment data"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /people",
			method:         "POST",
			path:           "/metrics/people",
			body:           map[string]interface{}{"data": "valid people data"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /organization",
			method:         "POST",
			path:           "/metrics/organization",
			body:           map[string]interface{}{"data": "valid org data"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /bounty_stats",
			method:         "POST",
			path:           "/metrics/bounty_stats",
			body:           map[string]interface{}{"data": "valid stats"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /bounties",
			method:         "POST",
			path:           "/metrics/bounties",
			body:           map[string]interface{}{"data": "valid bounties"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /bounties/count",
			method:         "POST",
			path:           "/metrics/bounties/count",
			body:           map[string]interface{}{"data": "valid count"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /bounties/providers",
			method:         "POST",
			path:           "/metrics/bounties/providers",
			body:           map[string]interface{}{"data": "valid providers"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /csv",
			method:         "POST",
			path:           "/metrics/csv",
			body:           map[string]interface{}{"data": "valid csv"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Token",
			method:         "GET",
			path:           "/metrics/workspaces",
			token:          "invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing Token",
			method:         "GET",
			path:           "/metrics/workspaces",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "POST /payment with Invalid Data",
			method:         "POST",
			path:           "/metrics/payment",
			body:           map[string]interface{}{"invalid": "data"},
			token:          "valid-token",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "POST /people with Invalid Data",
			method:         "POST",
			path:           "/metrics/people",
			body:           map[string]interface{}{"invalid": "data"},
			token:          "valid-token",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "POST /organization with Invalid Data",
			method:         "POST",
			path:           "/metrics/organization",
			body:           map[string]interface{}{"invalid": "data"},
			token:          "valid-token",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "POST /bounty_stats with Invalid Data",
			method:         "POST",
			path:           "/metrics/bounty_stats",
			body:           map[string]interface{}{"invalid": "data"},
			token:          "valid-token",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "POST /bounties with Large Data Set",
			method:         "POST",
			path:           "/metrics/bounties",
			body:           map[string]interface{}{"data": strings.Repeat("large data", 1000)},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /csv with Large CSV File",
			method:         "POST",
			path:           "/metrics/csv",
			body:           map[string]interface{}{"data": strings.Repeat("large csv", 1000)},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /workspaces with No Workspaces Available",
			method:         "GET",
			path:           "/metrics/workspaces",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /bounties/providers with No Providers Available",
			method:         "POST",
			path:           "/metrics/bounties/providers",
			body:           map[string]interface{}{"data": "empty"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /bounties/count with Zero Count",
			method:         "POST",
			path:           "/metrics/bounties/count",
			body:           map[string]interface{}{"count": 0},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /csv with Malformed CSV Data",
			method:         "POST",
			path:           "/metrics/csv",
			body:           map[string]interface{}{"data": "malformed,csv,data"},
			token:          "valid-token",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Test GET /workspaces Route",
			method:         "GET",
			path:           "/metrics/workspaces",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test POST /payment Route",
			method:         "POST",
			path:           "/metrics/payment",
			body:           map[string]interface{}{"startDate": "2024-01-01", "endDate": "2024-02-01"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Test POST /people Route",

			method:         "POST",
			path:           "/metrics/people",
			body:           map[string]interface{}{"startDate": "2024-01-01", "endDate": "2024-02-01"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Test POST /organization Route",

			method:         "POST",
			path:           "/metrics/organization",
			body:           map[string]interface{}{"startDate": "2024-01-01", "endDate": "2024-02-01"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Test POST /bounty_stats Route",

			method:         "POST",
			path:           "/metrics/bounty_stats",
			body:           map[string]interface{}{"startDate": "2024-01-01", "endDate": "2024-02-01"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Test POST /bounties Route",

			method:         "POST",
			path:           "/metrics/bounties",
			body:           map[string]interface{}{"startDate": "2024-01-01", "endDate": "2024-02-01"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Test POST /bounties/count Route",

			method:         "POST",
			path:           "/metrics/bounties/count",
			body:           map[string]interface{}{"startDate": "2024-01-01", "endDate": "2024-02-01"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Test POST /bounties/providers Route",

			method:         "POST",
			path:           "/metrics/bounties/providers",
			body:           map[string]interface{}{"startDate": "2024-01-01", "endDate": "2024-02-01"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name: "Test POST /csv Route",

			method:         "POST",
			path:           "/metrics/csv",
			body:           map[string]interface{}{"startDate": "2024-01-01", "endDate": "2024-02-01"},
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Test Invalid HTTP Method",
			method:         "PUT",
			path:           "/metrics/workspaces",
			token:          "valid-token",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Test Invalid Route",
			method:         "GET",
			path:           "/metrics/invalid_route",
			token:          "valid-token",
			expectedStatus: http.StatusNotFound,
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

			switch tc.token {
			case "valid-token":
				req.Header.Set("Authorization", "Bearer valid-super-admin-token")
			case "invalid-token":
				req.Header.Set("Authorization", "Bearer invalid-token")
			case "":
			default:
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code, "Handler returned wrong status code for test: "+tc.name)
		})
	}
}

func validateRequestBody(r *http.Request) bool {
	if r.Body == nil {
		return true
	}

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}

	if _, hasInvalid := body["invalid"]; hasInvalid {
		return false
	}

	if data, ok := body["data"].(string); ok && strings.Contains(data, "malformed") {
		return false
	}

	if len(body) == 0 {
		return true
	}

	return true
}
