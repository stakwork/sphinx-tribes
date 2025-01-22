package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSnippet(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()

	sHandler := NewSnippetHandler(&http.Client{}, db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	tests := []struct {
		name         string
		requestBody  map[string]string
		auth         string
		workspaceID  string
		expectedCode int
		validate     func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			requestBody: map[string]string{
				"title":   "Test Snippet",
				"snippet": "Test Content",
			},
			auth:         person.OwnerPubKey,
			workspaceID:  workspace.Uuid,
			expectedCode: http.StatusCreated,
			validate: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var snippet db.TextSnippet
				err := json.NewDecoder(rr.Body).Decode(&snippet)
				require.NoError(t, err)
				assert.Equal(t, "Test Snippet", snippet.Title)
				assert.Equal(t, "Test Content", snippet.Snippet)
				assert.Equal(t, workspace.Uuid, snippet.WorkspaceUUID)
			},
		},
		{
			name: "unauthorized - no auth token",
			requestBody: map[string]string{
				"title":   "Test Snippet",
				"snippet": "Test Content",
			},
			workspaceID:  workspace.Uuid,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "bad request - missing workspace_uuid",
			requestBody: map[string]string{
				"title":   "Test Snippet",
				"snippet": "Test Content",
			},
			auth:         person.OwnerPubKey,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "bad request - missing title",
			requestBody: map[string]string{
				"snippet": "Test Content",
			},
			auth:         person.OwnerPubKey,
			workspaceID:  workspace.Uuid,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "bad request - missing snippet content",
			requestBody: map[string]string{
				"title": "Test Snippet",
			},
			auth:         person.OwnerPubKey,
			workspaceID:  workspace.Uuid,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest(http.MethodPost, "/snippet/create", bytes.NewReader(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			if tt.workspaceID != "" {
				req.URL.RawQuery = "workspace_uuid=" + tt.workspaceID
			}

			if tt.auth != "" {
				req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, tt.auth))
			}

			sHandler.CreateSnippet(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.validate != nil {
				tt.validate(t, rr)
			}
		})
	}
}

func TestGetSnippetsByWorkspace(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()

	sHandler := NewSnippetHandler(&http.Client{}, db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	snippet := db.TextSnippet{
		WorkspaceUUID: workspace.Uuid,
		Title:         "Test Snippet",
		Snippet:       "Test Content",
	}
	createdSnippet, _ := db.TestDB.CreateSnippet(&snippet)

	tests := []struct {
		name          string
		workspaceUuid string
		auth          string
		expectedCode  int
		validate      func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:          "success",
			workspaceUuid: workspace.Uuid,
			auth:          person.OwnerPubKey,
			expectedCode:  http.StatusOK,
			validate: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var snippets []db.TextSnippet
				err := json.NewDecoder(rr.Body).Decode(&snippets)
				require.NoError(t, err)
				assert.Len(t, snippets, 1)
				assert.Equal(t, createdSnippet.Title, snippets[0].Title)
			},
		},
		{
			name:          "unauthorized",
			workspaceUuid: workspace.Uuid,
			expectedCode:  http.StatusUnauthorized,
		},
		{
			name:         "empty workspace uuid",
			auth:         person.OwnerPubKey,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/snippet/workspace/"+tt.workspaceUuid, nil)

			if tt.workspaceUuid != "" {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("workspace_uuid", tt.workspaceUuid)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			}

			if tt.auth != "" {
				req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, tt.auth))
			}

			sHandler.GetSnippetsByWorkspace(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.validate != nil {
				tt.validate(t, rr)
			}
		})
	}
}

func TestGetSnippetByID(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()

	sHandler := NewSnippetHandler(&http.Client{}, db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}

	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}

	db.TestDB.CreateOrEditWorkspace(workspace)

	snippet := db.TextSnippet{
		WorkspaceUUID: workspace.Uuid,
		Title:         "Test Snippet",
		Snippet:       "Test Content",
	}

	createdSnippet, err := db.TestDB.CreateSnippet(&snippet)
	require.NoError(t, err)

	tests := []struct {
		name         string
		snippetID    string
		auth         string
		expectedCode int
		validate     func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:         "success",
			snippetID:    fmt.Sprintf("%d", createdSnippet.ID),
			auth:         person.OwnerPubKey,
			expectedCode: http.StatusOK,
			validate: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var snippet db.TextSnippet
				err := json.NewDecoder(rr.Body).Decode(&snippet)
				require.NoError(t, err)
				assert.Equal(t, createdSnippet.Title, snippet.Title)
				assert.Equal(t, createdSnippet.Snippet, snippet.Snippet)
				assert.Equal(t, createdSnippet.WorkspaceUUID, snippet.WorkspaceUUID)
			},
		},
		{
			name:         "unauthorized",
			snippetID:    "1",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid snippet id",
			snippetID:    "invalid",
			auth:         person.OwnerPubKey,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "snippet not found",
			snippetID:    "999",
			auth:         person.OwnerPubKey,
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/snippet/"+tt.snippetID, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.snippetID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.auth != "" {
				req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, tt.auth))
			}

			sHandler.GetSnippetByID(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.validate != nil {
				tt.validate(t, rr)
			}
		})
	}
}

func TestUpdateSnippet(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	sHandler := NewSnippetHandler(&http.Client{}, db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	snippet := db.TextSnippet{
		WorkspaceUUID: workspace.Uuid,
		Title:         "Test Snippet",
		Snippet:       "Test Content",
	}

	createdSnippet, err := db.TestDB.CreateSnippet(&snippet)
	require.NoError(t, err)

	tests := []struct {
		name         string
		snippetID    string
		requestBody  map[string]string
		auth         string
		expectedCode int
		validate     func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "success",
			snippetID: fmt.Sprintf("%d", createdSnippet.ID),
			requestBody: map[string]string{
				"title":   "Updated Snippet",
				"snippet": "Updated Content",
			},
			auth:         person.OwnerPubKey,
			expectedCode: http.StatusOK,
			validate: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var snippet db.TextSnippet
				err := json.NewDecoder(rr.Body).Decode(&snippet)
				require.NoError(t, err)
				assert.Equal(t, "Updated Snippet", snippet.Title)
				assert.Equal(t, "Updated Content", snippet.Snippet)
				assert.Equal(t, createdSnippet.WorkspaceUUID, snippet.WorkspaceUUID)
			},
		},
		{
			name:      "unauthorized",
			snippetID: fmt.Sprintf("%d", createdSnippet.ID),
			requestBody: map[string]string{
				"title":   "Updated Snippet",
				"snippet": "Updated Content",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid snippet id",
			snippetID:    "invalid",
			auth:         person.OwnerPubKey,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:      "snippet not found",
			snippetID: "999",
			requestBody: map[string]string{
				"title":   "Updated Snippet",
				"snippet": "Updated Content",
			},
			auth:         person.OwnerPubKey,
			expectedCode: http.StatusNotFound,
		},
		{
			name:      "invalid request body",
			snippetID: fmt.Sprintf("%d", createdSnippet.ID),
			auth:      person.OwnerPubKey,
			requestBody: map[string]string{
				"invalid": "data",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			requestBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/snippet/"+tt.snippetID, bytes.NewReader(requestBody))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.snippetID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.auth != "" {
				req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, tt.auth))
			}

			sHandler.UpdateSnippet(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.validate != nil {
				tt.validate(t, rr)
			}
		})
	}
}

func TestDeleteSnippet(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()

	sHandler := NewSnippetHandler(&http.Client{}, db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	snippet := db.TextSnippet{
		WorkspaceUUID: workspace.Uuid,
		Title:         "Test Snippet",
		Snippet:       "Test Content",
	}

	createdSnippet, err := db.TestDB.CreateSnippet(&snippet)
	require.NoError(t, err)

	tests := []struct {
		name         string
		snippetID    string
		auth         string
		expectedCode int
		validate     func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:         "success",
			snippetID:    fmt.Sprintf("%d", createdSnippet.ID),
			auth:         person.OwnerPubKey,
			expectedCode: http.StatusOK,
			validate: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.NewDecoder(rr.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, "Snippet deleted successfully", response["message"])

				// Verify snippet was deleted
				_, err = db.TestDB.GetSnippetByID(createdSnippet.ID)
				assert.Error(t, err)
				assert.Equal(t, "record not found", err.Error())
			},
		},
		{
			name:         "unauthorized",
			snippetID:    fmt.Sprintf("%d", createdSnippet.ID),
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid snippet id",
			snippetID:    "invalid",
			auth:         person.OwnerPubKey,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "snippet not found",
			snippetID:    "999",
			auth:         person.OwnerPubKey,
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/snippet/"+tt.snippetID, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.snippetID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.auth != "" {
				req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, tt.auth))
			}

			sHandler.DeleteSnippet(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.validate != nil {
				tt.validate(t, rr)
			}
		})
	}
}
