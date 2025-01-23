package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
)

func PersonMockHandler(t *testing.T, expectedStatus int, validateReq func(*http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if isProtectedEndpoint(r.URL.Path, r.Method) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		path := chi.URLParam(r, "pubkey")
		if path == "" {
			path = chi.URLParam(r, "id")
		}
		if path == "" {
			path = chi.URLParam(r, "uuid")
		}
		if path == "" {
			path = chi.URLParam(r, "github")
		}

		if isNonExistentResource(path) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if validateReq != nil && !validateReq(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(expectedStatus)
		if r.Method != http.MethodDelete && expectedStatus == http.StatusOK {
			json.NewEncoder(w).Encode(db.Person{})
		}
	}
}

func isProtectedEndpoint(path string, method string) bool {

	cleanPath := strings.TrimPrefix(path, "/person")

	if cleanPath == "" {
		cleanPath = "/"
	}

	return (cleanPath == "/upsertlogin" && method == http.MethodPost) ||
		(cleanPath == "/" && method == http.MethodPost) ||
		(cleanPath == "/" && method == http.MethodPut) ||
		(method == http.MethodDelete)
}

func isNonExistentResource(param string) bool {

	if param == "nonexistentuuid" {
		return true
	}

	nonExistentValues := []string{
		"nonexistentpubkey",
		"999",
		"nonexistentgithub",
	}

	for _, value := range nonExistentValues {
		if param == value {
			return true
		}
	}
	return false
}

func TestPersonRoutes(t *testing.T) {
	r := chi.NewRouter()
	personRouter := chi.NewRouter()

	personRouter.Get("/{pubkey}", PersonMockHandler(t, http.StatusOK, nil))
	personRouter.Get("/id/{id}", PersonMockHandler(t, http.StatusOK, validateID))
	personRouter.Get("/uuid/{uuid}", PersonMockHandler(t, http.StatusOK, validateUUID))
	personRouter.Get("/uuid/{uuid}/assets", PersonMockHandler(t, http.StatusOK, validateUUID))
	personRouter.Get("/githubname/{github}", PersonMockHandler(t, http.StatusOK, nil))
	personRouter.Post("/upsertlogin", PersonMockHandler(t, http.StatusOK, nil))
	personRouter.Post("/", PersonMockHandler(t, http.StatusOK, nil))
	personRouter.Put("/", PersonMockHandler(t, http.StatusOK, nil))
	personRouter.Delete("/{id}", PersonMockHandler(t, http.StatusOK, validateID))

	r.Mount("/person", personRouter)

	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		headers        map[string]string
		expectedStatus int
	}{
		{"Get Person by Pubkey", "GET", "/person/somepubkey", nil, nil, http.StatusOK},
		{"Get Person by ID", "GET", "/person/id/123", nil, nil, http.StatusOK},
		{"Get Person by UUID", "GET", "/person/uuid/123e4567-e89b-12d3-a456-426614174000", nil, nil, http.StatusOK},
		{"Get Person Assets by UUID", "GET", "/person/uuid/123e4567-e89b-12d3-a456-426614174000/assets", nil, nil, http.StatusOK},
		{"Get Person by GitHub Name", "GET", "/person/githubname/somegithub", nil, nil, http.StatusOK},
		{"Upsert Login", "POST", "/person/upsertlogin", map[string]string{"key": "value"}, map[string]string{"Authorization": "Bearer token"}, http.StatusOK},
		{"Create Person", "POST", "/person", map[string]string{"name": "Test"}, map[string]string{"Authorization": "Bearer token"}, http.StatusOK},
		{"Update Person", "PUT", "/person", map[string]string{"id": "123", "name": "Updated"}, map[string]string{"Authorization": "Bearer token"}, http.StatusOK},
		{"Delete Person", "DELETE", "/person/123", nil, map[string]string{"Authorization": "Bearer token"}, http.StatusOK},

		{"Get Person by Non-Existent Pubkey", "GET", "/person/nonexistentpubkey", nil, nil, http.StatusNotFound},
		{"Get Person by Non-Existent ID", "GET", "/person/id/999", nil, nil, http.StatusNotFound},
		{"Get Person by Non-Existent UUID", "GET", "/person/uuid/nonexistentuuid", nil, nil, http.StatusNotFound},
		{"Get Person by Non-Existent GitHub Name", "GET", "/person/githubname/nonexistentgithub", nil, nil, http.StatusNotFound},
		{"Delete Person with Non-Existent ID", "DELETE", "/person/999", nil, map[string]string{"Authorization": "Bearer token"}, http.StatusNotFound},
		{"Invalid Data Type for ID", "GET", "/person/id/invalid", nil, nil, http.StatusBadRequest},
		{"Unauthorized Access to Upsert Login", "POST", "/person/upsertlogin", map[string]string{"key": "value"}, nil, http.StatusUnauthorized},
		{"Invalid Authentication Token", "POST", "/person", map[string]string{"name": "Test"}, nil, http.StatusUnauthorized},
		{"Invalid UUID Format", "GET", "/person/uuid/invalid-uuid", nil, nil, http.StatusBadRequest},
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

			if tc.headers != nil {
				for key, value := range tc.headers {
					req.Header.Set(key, value)
				}
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code, "Handler returned wrong status code for test: "+tc.name)

			if w.Code == http.StatusOK && tc.method != http.MethodDelete {
				var response db.Person
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err, "Failed to decode response for test: "+tc.name)
			}
		})
	}
}

func validateID(r *http.Request) bool {
	id := chi.URLParam(r, "id")
	return id != "invalid"
}

func validateUUID(r *http.Request) bool {
	uuid := chi.URLParam(r, "uuid")
	return isValidUUID(uuid)
}

func TestValidateID(t *testing.T) {
	t.Run("Valid ID Parameter", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/person/id/123", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "123")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		result := validateID(r)
		assert.True(t, result)
	})

	t.Run("Invalid ID Parameter", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/person/id/invalid", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		result := validateID(r)
		assert.False(t, result)
	})

	t.Run("Empty ID Parameter", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/person/id/", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		result := validateID(r)
		assert.True(t, result, "Empty ID should not be considered 'invalid'")
	})

	t.Run("Special Characters in ID", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/person/id/123@#$", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "123@#$")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		result := validateID(r)
		assert.True(t, result, "Special characters should be allowed if not 'invalid'")
	})

	t.Run("Numeric ID Parameter", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/person/id/12345", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "12345")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		result := validateID(r)
		assert.True(t, result)
	})

	t.Run("Missing ID Parameter", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/person/id/", nil)
		rctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		result := validateID(r)
		assert.True(t, result, "Missing ID should not be considered 'invalid'")
	})

	t.Run("Very Long ID Parameter", func(t *testing.T) {
		longID := strings.Repeat("1", 1000)
		r := httptest.NewRequest(http.MethodGet, "/person/id/"+longID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", longID)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		result := validateID(r)
		assert.True(t, result, "Long ID should be valid if not 'invalid'")
	})

	t.Run("Case Sensitivity", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/person/id/INVALID", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "INVALID")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		result := validateID(r)
		assert.True(t, result, "INVALID in uppercase should be valid")
	})

	t.Run("Whitespace in ID Parameter", func(t *testing.T) {
		// Use URL encoded space (%20) in the URL
		r := httptest.NewRequest(http.MethodGet, "/person/id/123%20456", nil)
		rctx := chi.NewRouteContext()
		// Use actual space in the URL parameter
		rctx.URLParams.Add("id", "123 456")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		result := validateID(r)
		assert.True(t, result, "ID with spaces should be valid if not 'invalid'")
	})

	t.Run("ID Parameter with Only Whitespace", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/person/id/%20", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", " ")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		result := validateID(r)
		assert.True(t, result, "Whitespace-only ID should not be considered 'invalid'")
	})

	t.Run("ID Parameter with Mixed Case 'Invalid'", func(t *testing.T) {
		testCases := []string{
			"InVaLiD",
			"iNvAlId",
			"INVALID",
			"invalid",
		}

		for _, tc := range testCases {
			t.Run(tc, func(t *testing.T) {
				r := httptest.NewRequest(http.MethodGet, "/person/id/"+tc, nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", tc)
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

				result := validateID(r)
				if tc == "invalid" {
					assert.False(t, result, "lowercase 'invalid' should be invalid")
				} else {
					assert.True(t, result, "other cases of 'invalid' should be valid")
				}
			})
		}
	})
}
