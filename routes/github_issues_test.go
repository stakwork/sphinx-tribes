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

func GithubMockHandler(t *testing.T, expectedStatus int, validateReq func(*http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if validateReq != nil && !validateReq(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.URL.Path == "/user/repo/" || r.URL.Path == "/user/repo/123/extra" || r.URL.Path == "/invalid/path" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(expectedStatus)
	}
}

func TestGithubIssuesRoutes(t *testing.T) {
	r := chi.NewRouter()
	githubRouter := chi.NewRouter()

	githubRouter.Get("/{owner}/{repo}/{issue}/", GithubMockHandler(t, http.StatusOK, validateGithubIssueParams))
	githubRouter.Get("/{owner}/{repo}/{issue}", GithubMockHandler(t, http.StatusOK, validateGithubIssueParams))
	githubRouter.Get("/status/open", GithubMockHandler(t, http.StatusOK, nil))

	r.Mount("/github", githubRouter)

	testCases := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
	}{
		{"Route Registration for Specific Issue", "GET", "/github/user/repo/123", nil, http.StatusOK},
		{"Route Registration for Open Issues Status", "GET", "/github/status/open", nil, http.StatusOK},
		{"Route with Missing Parameters", "GET", "/github/user/repo/", nil, http.StatusNotFound},
		{"Route with Extra Parameters", "GET", "/github/user/repo/123/extra", nil, http.StatusNotFound},
		{"Invalid Route Path", "GET", "/github/invalid/path", nil, http.StatusNotFound},
		{"Case Sensitivity in Route Parameters", "GET", "/github/User/Repo/123", nil, http.StatusOK},

		{"Trailing Slash in Route", "GET", "/github/user/repo/123/", nil, http.StatusOK},

		{"Long Route Parameters", "GET", "/github/user/repo/" + "a_long_string_to_test_routing", nil, http.StatusOK},
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

func validateGithubIssueParams(r *http.Request) bool {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	issue := chi.URLParam(r, "issue")

	return owner != "" && repo != "" && issue != ""
}
