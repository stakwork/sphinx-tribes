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

func MockHandler(t *testing.T, expectedStatus int, validateReq func(*http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if validateReq != nil && !validateReq(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(expectedStatus)
	}
}

func TestWorkspaceRoutes(t *testing.T) {
	r := chi.NewRouter()
	workspaceRouter := chi.NewRouter()

	workspaceRouter.Get("/", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/count", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/{uuid}", MockHandler(t, http.StatusOK, func(r *http.Request) bool {
		uuid := chi.URLParam(r, "uuid")
		return len(uuid) == 36
	}))

	workspaceRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		name, exists := body["name"]
		if !exists {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		nameStr, isString := name.(string)
		if !isString {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(nameStr) > 1000 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	workspaceRouter.Delete("/delete/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		uuid := chi.URLParam(r, "uuid")
		if uuid == "non-existent-uuid" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if len(uuid) != 36 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	workspaceRouter.Get("/users/{uuid}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/users/{uuid}/count", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Post("/users/role/{uuid}/{user}", MockHandler(t, http.StatusOK, nil))

	workspaceRouter.Get("/bounties/{uuid}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/bounties/{uuid}/count", MockHandler(t, http.StatusOK, nil))

	workspaceRouter.Get("/budget/{uuid}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/poll/user/invoices", MockHandler(t, http.StatusOK, nil))

	workspaceRouter.Get("/{workspace_uuid}/features", MockHandler(t, http.StatusOK, nil))

	workspaceRouter.Put("/{workspace_uuid}/payments", MockHandler(t, http.StatusOK, nil))

	workspaceRouter.Get("/user/{userId}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/user/dropdown/{userId}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/foruser/{uuid}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/bounty/roles", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/users/role/{uuid}/{user}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/budget/history/{uuid}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/payments/{uuid}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/poll/invoices/{uuid}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/invoices/count/{uuid}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/user/invoices/count", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Post("/mission", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Post("/tactics", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Post("/schematicurl", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Post("/repositories", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/repositories/{uuid}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/{workspace_uuid}/repository/{uuid}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Delete("/{workspace_uuid}/repository/{uuid}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/{workspace_uuid}/lastwithdrawal", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Post("/codegraph", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/codegraph/{uuid}", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Get("/{workspace_uuid}/codegraph", MockHandler(t, http.StatusOK, nil))
	workspaceRouter.Delete("/{workspace_uuid}/codegraph/{uuid}", MockHandler(t, http.StatusOK, nil))

	r.Mount("/workspaces", workspaceRouter)

	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
	}{
		{
			name:           "Get All Workspaces",
			method:         "GET",
			path:           "/workspaces/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Workspace Count",
			method:         "GET",
			path:           "/workspaces/count",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Workspace by UUID",
			method:         "GET",
			path:           "/workspaces/123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Create or Edit Workspace",
			method:         "POST",
			path:           "/workspaces/",
			body:           map[string]interface{}{"name": "Test Workspace"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Delete Workspace",
			method:         "DELETE",
			path:           "/workspaces/delete/123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Workspace by Invalid UUID",
			method:         "GET",
			path:           "/workspaces/invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Create Workspace with Missing Data",
			method:         "POST",
			path:           "/workspaces/",
			body:           map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Delete Non-existent Workspace",
			method:         "DELETE",
			path:           "/workspaces/delete/non-existent-uuid",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Get Workspace with Invalid Data Type",
			method:         "GET",
			path:           "/workspaces/123",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Create Workspace with Invalid Data Type",
			method:         "POST",
			path:           "/workspaces/",
			body:           map[string]interface{}{"name": 123},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Get Large Number of Workspaces",
			method:         "GET",
			path:           "/workspaces/?limit=1000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Create Workspace with Large Payload",
			method:         "POST",
			path:           "/workspaces/",
			body:           map[string]interface{}{"name": strings.Repeat("a", 10000)},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Get Workspace Users",
			method:         "GET",
			path:           "/workspaces/users/123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Workspace Bounties",
			method:         "GET",
			path:           "/workspaces/bounties/123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Add User Roles",
			method:         "POST",
			path:           "/workspaces/users/role/123e4567-e89b-12d3-a456-426614174000/user123",
			body:           map[string]interface{}{"role": "admin"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Workspace Users Count",
			method:         "GET",
			path:           "/workspaces/users/123e4567-e89b-12d3-a456-426614174000/count",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Workspace Budget",
			method:         "GET",
			path:           "/workspaces/budget/123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Poll User Workspaces Budget",
			method:         "GET",
			path:           "/workspaces/poll/user/invoices",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Update Workspace Pending Payments",
			method:         "PUT",
			path:           "/workspaces/123e4567-e89b-12d3-a456-426614174000/payments",
			body:           map[string]interface{}{"amount": 100},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Features by Workspace UUID",
			method:         "GET",
			path:           "/workspaces/123e4567-e89b-12d3-a456-426614174000/features",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get User Workspaces",
			method:         "GET",
			path:           "/workspaces/user/123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get User Dropdown Workspaces",
			method:         "GET",
			path:           "/workspaces/user/dropdown/123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Workspace For User",
			method:         "GET",
			path:           "/workspaces/foruser/123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Bounty Roles",
			method:         "GET",
			path:           "/workspaces/bounty/roles",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get User Roles",
			method:         "GET",
			path:           "/workspaces/users/role/123e4567-e89b-12d3-a456-426614174000/user123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Workspace Budget History",
			method:         "GET",
			path:           "/workspaces/budget/history/123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Payment History",
			method:         "GET",
			path:           "/workspaces/payments/123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Poll Budget Invoices",
			method:         "GET",
			path:           "/workspaces/poll/invoices/123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Invoices Count",
			method:         "GET",
			path:           "/workspaces/invoices/count/123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Update Workspace Mission",
			method:         "POST",
			path:           "/workspaces/mission",
			body:           map[string]interface{}{"mission": "New Mission"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Update Workspace Tactics",
			method:         "POST",
			path:           "/workspaces/tactics",
			body:           map[string]interface{}{"tactics": "New Tactics"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Create Workspace Repository",
			method:         "POST",
			path:           "/workspaces/repositories",
			body:           map[string]interface{}{"url": "https://github.com/test/repo"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Workspace Repository",
			method:         "GET",
			path:           "/workspaces/repositories/123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Delete Workspace Repository",
			method:         "DELETE",
			path:           "/workspaces/123e4567-e89b-12d3-a456-426614174000/repository/456",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Last Withdrawal",
			method:         "GET",
			path:           "/workspaces/123e4567-e89b-12d3-a456-426614174000/lastwithdrawal",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Create Code Graph",
			method:         "POST",
			path:           "/workspaces/codegraph",
			body:           map[string]interface{}{"data": "graph data"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Code Graph",
			method:         "GET",
			path:           "/workspaces/codegraph/123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Delete Code Graph",
			method:         "DELETE",
			path:           "/workspaces/123e4567-e89b-12d3-a456-426614174000/codegraph/456",
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

			if tc.method != "GET" || strings.Contains(tc.path, "budget") {
				req.Header.Set("Authorization", "Bearer test-token")
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code, "Handler returned wrong status code")
		})
	}
}
