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

func PeopleMockHandler(t *testing.T, expectedStatus int, validateReq func(*http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if validateReq != nil && !validateReq(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.URL.Path == "/non-existent" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(expectedStatus)
	}
}

func TestPeopleRoutes(t *testing.T) {
	r := chi.NewRouter()
	peopleRouter := chi.NewRouter()

	peopleRouter.Get("/", PeopleMockHandler(t, http.StatusOK, nil))
	peopleRouter.Get("/search", PeopleMockHandler(t, http.StatusOK, nil))
	peopleRouter.Get("/posts", PeopleMockHandler(t, http.StatusOK, nil))
	peopleRouter.Get("/wanteds/assigned/{uuid}", PeopleMockHandler(t, http.StatusOK, peopleValidateUUID))
	peopleRouter.Get("/wanteds/created/{uuid}", PeopleMockHandler(t, http.StatusOK, peopleValidateUUID))
	peopleRouter.Get("/wanteds/header", PeopleMockHandler(t, http.StatusOK, nil))
	peopleRouter.Get("/short", PeopleMockHandler(t, http.StatusOK, nil))
	peopleRouter.Get("/offers", PeopleMockHandler(t, http.StatusOK, nil))
	peopleRouter.Get("/bounty/leaderboard", PeopleMockHandler(t, http.StatusOK, nil))

	r.Mount("/people", peopleRouter)

	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
	}{
		{"Root Endpoint", "GET", "/people/", nil, http.StatusOK},
		{"Search Endpoint", "GET", "/people/search", nil, http.StatusOK},
		{"Posts Endpoint", "GET", "/people/posts", nil, http.StatusOK},
		{"Assigned Bounties Endpoint", "GET", "/people/wanteds/assigned/123e4567-e89b-12d3-a456-426614174000", nil, http.StatusOK},
		{"Created Bounties Endpoint", "GET", "/people/wanteds/created/123e4567-e89b-12d3-a456-426614174000", nil, http.StatusOK},
		{"Wanteds Header Endpoint", "GET", "/people/wanteds/header", nil, http.StatusOK},
		{"Short List Endpoint", "GET", "/people/short", nil, http.StatusOK},
		{"Offers Endpoint", "GET", "/people/offers", nil, http.StatusOK},
		{"Bounty Leaderboard Endpoint", "GET", "/people/bounty/leaderboard", nil, http.StatusOK},
		{"Invalid UUID Format", "GET", "/people/wanteds/assigned/invalid-uuid", nil, http.StatusBadRequest},
		{"Non-Existent Endpoint", "GET", "/people/non-existent", nil, http.StatusNotFound},
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

func peopleValidateUUID(r *http.Request) bool {
	uuid := chi.URLParam(r, "uuid")
	return uuid != "" && isValidUUID(uuid)
}
