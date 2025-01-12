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

func BotMockHandler(t *testing.T, expectedStatus int, validateReq func(*http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if validateReq != nil && !validateReq(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if chi.URLParam(r, "uuid") == "unknownuuid" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(expectedStatus)
	}
}

func TestBotsRoutes(t *testing.T) {
	r := chi.NewRouter()
	botsRouter := chi.NewRouter()

	botsRouter.Post("/", BotMockHandler(t, http.StatusCreated, validateCreateOrEditBot))
	botsRouter.Get("/", BotMockHandler(t, http.StatusOK, nil))
	botsRouter.Get("/owner/{pubkey}", BotMockHandler(t, http.StatusOK, validatePubkey))
	botsRouter.Get("/{uuid}", BotMockHandler(t, http.StatusOK, botsValidateUUID))

	r.Mount("/bots", botsRouter)

	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
	}{
		{"Create or Edit Bot", "POST", "/bots/", map[string]interface{}{"name": "New Bot"}, http.StatusCreated},
		{"Get Listed Bots", "GET", "/bots/", nil, http.StatusOK},
		{"Get Bots by Owner", "GET", "/bots/owner/testpubkey", nil, http.StatusOK},
		{"Get Bot by Unknown UUID", "GET", "/bots/unknownuuid", nil, http.StatusBadRequest},
		{"Create Bot with Missing Fields", "POST", "/bots/", map[string]interface{}{}, http.StatusBadRequest},
		{"Invalid Pubkey Format", "GET", "/bots/owner/invalid!@#", nil, http.StatusBadRequest},
		{"Empty Pubkey Parameter", "GET", "/bots/owner/", nil, http.StatusNotFound},
		{"Invalid JSON in Create Bot", "POST", "/bots/", "Invalid JSON", http.StatusBadRequest},
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

func validatePubkey(r *http.Request) bool {
	pubkey := chi.URLParam(r, "pubkey")
	return pubkey != "" && !containsInvalidCharacters(pubkey)
}

func botsValidateUUID(r *http.Request) bool {
	uuid := chi.URLParam(r, "uuid")
	return uuid != "" && isValidUUID(uuid)
}

func validateCreateOrEditBot(r *http.Request) bool {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return false
	}
	_, hasName := body["name"]
	return hasName
}
