package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TribeMockHandler(t *testing.T, expectedStatus int, validateReq func(*http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if validateReq != nil && !validateReq(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.URL.Path == "/nonExistentUUID" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(expectedStatus)
	}
}

func TestTribeRoutes(t *testing.T) {
	r := chi.NewRouter()
	tribeRouter := chi.NewRouter()

	tribeRouter.Get("/", TribeMockHandler(t, http.StatusOK, nil))
	tribeRouter.Get("/app_url/{app_url}", TribeMockHandler(t, http.StatusOK, validateAppURL))
	tribeRouter.Get("/app_urls/{app_urls}", TribeMockHandler(t, http.StatusOK, validateAppURLs))
	tribeRouter.Get("/{uuid}", TribeMockHandler(t, http.StatusOK, tribesValidateUUID))
	tribeRouter.Get("/total", TribeMockHandler(t, http.StatusOK, nil))
	tribeRouter.Post("/", TribeMockHandler(t, http.StatusCreated, validateCreateOrEditTribe))

	r.Mount("/tribes", tribeRouter)

	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
	}{
		{"Get Listed Tribes", "GET", "/tribes/", nil, http.StatusOK},
		{"Get Tribes by App URL", "GET", "/tribes/app_url/sampleAppUrl", nil, http.StatusOK},
		{"Get Tribes by App URLs", "GET", "/tribes/app_urls/sampleAppUrl1,sampleAppUrl2", nil, http.StatusOK},
		{"Get Tribe by UUID", "GET", "/tribes/123e4567-e89b-12d3-a456-426614174000", nil, http.StatusOK},
		{"Get Total Tribes", "GET", "/tribes/total", nil, http.StatusOK},
		{"Create Tribe", "POST", "/tribes/", map[string]interface{}{"name": "New Tribe"}, http.StatusCreated},
		{"Get Tribe by Non-Existent UUID", "GET", "/tribes/nonExistentUUID", nil, http.StatusBadRequest},
		{"Create Tribe with Missing Fields", "POST", "/tribes/", map[string]interface{}{}, http.StatusBadRequest},
		{"Invalid App URL Format", "GET", "/tribes/app_url/invalidFormat!@#", nil, http.StatusBadRequest},
		{"Invalid JSON in Create Tribe", "POST", "/tribes/", "Invalid JSON", http.StatusBadRequest},
		{"Empty App URL Parameter", "GET", "/tribes/app_url/", nil, http.StatusNotFound},
		{"Empty UUID Parameter", "GET", "/tribes/", nil, http.StatusOK},
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

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code, "Handler returned wrong status code for test: "+tc.name)
		})
	}
}

func validateAppURL(r *http.Request) bool {
	appURL := chi.URLParam(r, "app_url")
	return appURL != "" && !containsInvalidCharacters(appURL)
}

func validateAppURLs(r *http.Request) bool {
	appURLs := chi.URLParam(r, "app_urls")
	return appURLs != "" && !containsInvalidCharacters(appURLs)
}

func tribesValidateUUID(r *http.Request) bool {
	uuid := chi.URLParam(r, "uuid")
	return uuid != "" && isValidUUID(uuid)
}

func validateCreateOrEditTribe(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasName := body["name"]
	return hasName
}

func containsInvalidCharacters(input string) bool {
	for _, c := range input {
		if c == '!' || c == '@' || c == '#' {
			return true
		}
	}
	return false
}

func isValidUUID(uuid string) bool {
	return len(uuid) == 36
}
