package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUnitCreateOrEditWorkspace(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	oHandler := NewWorkspaceHandler(db.TestDB)

	t.Run("should return error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspace)

		invalidJson := []byte(`{"key": "value"`)

		// Include a dummy public key in the context
		ctx := context.WithValue(context.Background(), auth.ContextKey, "dummy-pub-key")

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("should return error if public key not present", func(t *testing.T) { //passed
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspace)

		invalidJson := []byte(`{"key": "value"}`)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return error org name is empty", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspace)

		invalidJson := []byte(`{"name": ""}`)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return error org name is more than 20", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspace)

		invalidJson := []byte(`{"name": "DemoTestingNewWorkspace"}`)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return error if org name contains only spaces", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspace)

		invalidJson := []byte(`{"name": "   "}`)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should trim spaces from workspace name", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspace)

		const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		rand.Seed(int64(time.Now().UnixNano()))

		b := make([]byte, 10)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		name := string(b)

		spacedName := "  " + name + "  "

		jsonInput := []byte(fmt.Sprintf(`{"name": "%s", "owner_pubkey": "test-key", "description": "Workspace Bounties Description"}`, spacedName))

		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(jsonInput))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseOrg db.Workspace
		err = json.Unmarshal(rr.Body.Bytes(), &responseOrg)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, name, responseOrg.Name)
	})

	t.Run("should successfully add workspace if request is valid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspace)

		const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		rand.Seed(int64(time.Now().UnixNano()))

		b := make([]byte, 10)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		name := string(b)

		workspace := db.Workspace{
			Uuid:        uuid.New().String(),
			Name:        name,
			OwnerPubKey: uuid.New().String(),
			Github:      "https://github.com/bounties",
			Website:     "https://www.bountieswebsite.com",
			Description: "Workspace Bounties Description",
		}
		db.TestDB.CreateOrEditWorkspace(workspace)

		Workspace := db.TestDB.GetWorkspaceByUuid(workspace.Uuid)
		workspace.ID = Workspace.ID

		requestBody, _ := json.Marshal(workspace)
		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, workspace, Workspace)
	})
	t.Run("should return error if org description is empty or too long", func(t *testing.T) {
		tests := []struct {
			name        string
			description string
			wantStatus  int
		}{
			{"long description", strings.Repeat("a", 121), http.StatusBadRequest},
		}

		for _, tc := range tests {
			t.Run(tc.description, func(t *testing.T) {
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(oHandler.CreateOrEditWorkspace)
				invalidJson := []byte(fmt.Sprintf(`{"name": "TestWorkspace", "owner_pubkey": "test-key", "description": "%s"}`, tc.description))
				ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
				req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
				if err != nil {
					t.Fatal(err)
				}

				handler.ServeHTTP(rr, req)

				assert.Equal(t, tc.wantStatus, rr.Code)
			})
		}
	})
}

func TestDeleteWorkspace(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	oHandler := NewWorkspaceHandler(db.TestDB)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        fmt.Sprintf("Workspace %s", uuid.New().String()),
		OwnerPubKey: "test-key",
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "Workspace Description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)
	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)
	ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

	t.Run("should return error if not authorized", func(t *testing.T) {
		workspaceUUID := workspace.Uuid
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteWorkspace)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodDelete, "/delete/"+workspaceUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should set workspace fields to null and delete users on successful delete", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteWorkspace)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/delete/"+workspaceUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		updatedOrg := db.TestDB.GetWorkspaceByUuid(workspaceUUID)
		assert.Equal(t, true, updatedOrg.Deleted)
		assert.Equal(t, "", updatedOrg.Website)
		assert.Equal(t, "", updatedOrg.Github)
		assert.Equal(t, "", updatedOrg.Description)
	})

	t.Run("should handle failures in database updates", func(t *testing.T) {
		workspaceUUID := workspace.Uuid
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if chi.URLParam(r, "uuid") == workspaceUUID {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			oHandler.DeleteWorkspace(w, r)
		})

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/delete/"+workspaceUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("should set workspace's deleted column to true", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteWorkspace)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/delete/"+workspaceUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		updatedOrg := db.TestDB.GetWorkspaceByUuid(workspaceUUID)
		assert.Equal(t, true, updatedOrg.Deleted)
	})

	t.Run("should set Website, Github, and Description to empty strings", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteWorkspace)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/delete/"+workspaceUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		updatedOrg := db.TestDB.GetWorkspaceByUuid(workspaceUUID)
		assert.Equal(t, "", updatedOrg.Website)
		assert.Equal(t, "", updatedOrg.Github)
		assert.Equal(t, "", updatedOrg.Description)
	})

	t.Run("should delete all users from the workspace", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteWorkspace)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/delete/"+workspaceUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		updatedOrg := db.TestDB.GetWorkspaceByUuid(workspaceUUID)
		assert.Equal(t, true, updatedOrg.Deleted)
	})
}

func TestGetWorkspaceBounties(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	mockGenerateBountyHandler := func(bounties []db.NewBounty) []db.BountyResponse {
		return []db.BountyResponse{} // Mocked response
	}
	oHandler := NewWorkspaceHandler(db.TestDB)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        uuid.New().String(),
		OwnerPubKey: "workspace_owner_bounties_pubkey",
		Github:      "https://github.com/bounties",
		Website:     "https://www.bountieswebsite.com",
		Description: "Workspace Bounties Description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	bounty := db.NewBounty{
		Type:          "coding",
		Title:         "existing bounty",
		Description:   "existing bounty description",
		WorkspaceUuid: workspace.Uuid,
		OwnerID:       "workspace-user",
		Price:         2000,
	}
	db.TestDB.CreateOrEditBounty(bounty)

	t.Run("Should test that a workspace's bounties can be listed without authentication", func(t *testing.T) {

		oHandler.generateBountyHandler = mockGenerateBountyHandler
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceBounties)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/bounties/"+workspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		fetchedWorkspace := db.TestDB.GetWorkspaceByUuid(workspace.Uuid)
		workspace.ID = fetchedWorkspace.ID

		fetchedBounty := db.TestDB.GetWorkspaceBounties(req, bounty.WorkspaceUuid)
		bounty.ID = fetchedBounty[0].ID
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, workspace, fetchedWorkspace)
		assert.Equal(t, bounty, fetchedBounty[0])
	})

	t.Run("should return empty array when wrong workspace UUID is passed", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceBounties)
		workspaceUUID := "wrong-uuid"

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/bounties/"+workspaceUUID+"?limit=10&sortBy=created&search=test&page=1&resetPage=true", nil)
		if err != nil {
			t.Fatal(err)
		}

		fetchedWorkspaceWrong := db.TestDB.GetWorkspaceByUuid(workspaceUUID)

		handler.ServeHTTP(rr, req)

		// Assert that the response status code is as expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Assert that the response body is an empty array
		assert.Equal(t, "[]\n", rr.Body.String())
		assert.NotEqual(t, workspace, fetchedWorkspaceWrong)
	})
}

func TestGetWorkspaceBudget(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	oHandler := NewWorkspaceHandler(db.TestDB)
	handlerUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
		return true
	}
	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "Workspace Budget Name " + uuid.New().String(),
		OwnerPubKey: "workspace_owner_budget_pubkey",
		Github:      "https://github.com/budget",
		Website:     "https://www.budgetwebsite.com",
		Description: "Workspace Budget Description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	budgetAmount := uint(5000)
	bounty := db.NewBountyBudget{
		WorkspaceUuid: workspace.Uuid,
		TotalBudget:   budgetAmount,
	}
	db.TestDB.CreateWorkspaceBudget(bounty)

	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

	t.Run("Should test that a 401 is returned when trying to view an workspace's budget without a token", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/budget/"+workspaceUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetWorkspaceBudget).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that the right workspace budget is returned, if the user is the workspace admin or has the ViewReport role", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		oHandler.userHasAccess = handlerUserHasAccess

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/budget/"+workspaceUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetWorkspaceBudget).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseBudget db.StatusBudget
		err = json.Unmarshal(rr.Body.Bytes(), &responseBudget)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, budgetAmount, responseBudget.CurrentBudget)
	})

	t.Run("Valid Request with Admin Access", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceBudget)

		oHandler.userHasAccess = func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/budget/"+workspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseBudget db.StatusBudget
		err = json.Unmarshal(rr.Body.Bytes(), &responseBudget)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, budgetAmount, responseBudget.CurrentBudget)
	})

	t.Run("Missing Public Key", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceBudget)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodGet, "/budget/"+workspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceBudget)

		oHandler.userHasAccess = func(pubKeyFromAuth string, uuid string, role string) bool {
			return false
		}

		ctx := context.WithValue(context.Background(), auth.ContextKey, "unauthorized_user")
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/budget/"+workspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Non-Existent UUID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceBudget)

		oHandler.userHasAccess = func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}

		nonExistentUUID := uuid.New().String()
		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", nonExistentUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/budget/"+nonExistentUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var responseBudget db.StatusBudget
		err = json.Unmarshal(rr.Body.Bytes(), &responseBudget)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, uint(0), responseBudget.CurrentBudget)
	})

	t.Run("Workspace Budget Not Set", func(t *testing.T) {

		workspaceNoBudget := db.Workspace{
			Uuid:        uuid.New().String(),
			Name:        "Test No Budget Workspace " + uuid.New().String(),
			OwnerPubKey: "test_no_budget_owner",
		}
		db.TestDB.CreateOrEditWorkspace(workspaceNoBudget)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceBudget)

		oHandler.userHasAccess = func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspaceNoBudget.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceNoBudget.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/budget/"+workspaceNoBudget.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var responseBudget db.StatusBudget
		err = json.Unmarshal(rr.Body.Bytes(), &responseBudget)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, uint(0), responseBudget.CurrentBudget)
	})

	t.Run("Empty UUID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceBudget)

		oHandler.userHasAccess = func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", "")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/budget/", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var responseBudget db.StatusBudget
		err = json.Unmarshal(rr.Body.Bytes(), &responseBudget)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, uint(0), responseBudget.CurrentBudget)
	})

	t.Run("Invalid UUID Format", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceBudget)

		oHandler.userHasAccess = func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}

		invalidUUID := "invalid-uuid-format"
		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", invalidUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/budget/"+invalidUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var responseBudget db.StatusBudget
		err = json.Unmarshal(rr.Body.Bytes(), &responseBudget)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, uint(0), responseBudget.CurrentBudget)
	})

	t.Run("Large Number of Workspaces", func(t *testing.T) {

		numWorkspaces := 100
		workspaces := make([]db.Workspace, numWorkspaces)
		for i := 0; i < numWorkspaces; i++ {
			workspace := db.Workspace{
				Uuid:        uuid.New().String(),
				Name:        fmt.Sprintf("Test Budget Workspace %d", i),
				OwnerPubKey: "test_budget_owner",
			}
			db.TestDB.CreateOrEditWorkspace(workspace)
			workspaces[i] = workspace

			budget := db.NewBountyBudget{
				WorkspaceUuid: workspace.Uuid,
				TotalBudget:   uint(1000 * (i + 1)),
			}
			db.TestDB.CreateWorkspaceBudget(budget)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceBudget)

		oHandler.userHasAccess = func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}

		lastWorkspace := workspaces[numWorkspaces-1]
		ctx := context.WithValue(context.Background(), auth.ContextKey, lastWorkspace.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", lastWorkspace.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/budget/"+lastWorkspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var responseBudget db.StatusBudget
		err = json.Unmarshal(rr.Body.Bytes(), &responseBudget)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, uint(1000*numWorkspaces), responseBudget.CurrentBudget)
	})
}

func TestGetWorkspaceBudgetHistory(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	oHandler := NewWorkspaceHandler(db.TestDB)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "Workspace History Name" + uuid.New().String(),
		OwnerPubKey: "test-key",
		Github:      "https://github.com/history",
		Website:     "https://www.historywebsite.com",
		Description: "Workspace History Description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)
	ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

	budgetAmount := uint(5000)
	bounty := db.NewBountyBudget{
		WorkspaceUuid: workspace.Uuid,
		TotalBudget:   budgetAmount,
	}
	db.TestDB.CreateWorkspaceBudget(bounty)

	now := time.Now()
	paymentHistory := db.NewPaymentHistory{
		WorkspaceUuid:  workspace.Uuid,
		Amount:         budgetAmount,
		Status:         true,
		PaymentType:    "budget",
		Created:        &now,
		Updated:        &now,
		SenderPubKey:   workspace.OwnerPubKey,
		ReceiverPubKey: "",
		BountyId:       0,
	}
	db.TestDB.AddPaymentHistory(paymentHistory)

	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

	t.Run("Should test that a 401 is returned when trying to view an workspace's budget history without a token", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		handlerUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
			return false
		}
		oHandler.userHasAccess = handlerUserHasAccess

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/budget/history/"+workspaceUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetWorkspaceBudgetHistory).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that the right budget history is returned, if the user is the workspace admin or has the ViewReport role", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		handlerUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}
		oHandler.userHasAccess = handlerUserHasAccess

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/budget/history/"+workspaceUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetWorkspaceBudgetHistory).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseBudgetHistory []db.BudgetHistoryData
		err = json.Unmarshal(rr.Body.Bytes(), &responseBudgetHistory)
		if err != nil {
			t.Fatal(err)
		}

		expectedBudgetHistory := db.TestDB.GetWorkspaceBudgetHistory(workspaceUUID)

		assert.Equal(t, expectedBudgetHistory, responseBudgetHistory)
	})
}

func TestGetWorkspaceBountiesCount(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	oHandler := NewWorkspaceHandler(db.TestDB)

	t.Run("should return the count of workspace bounties", func(t *testing.T) {

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceBountiesCount)

		expectedCount := int(1)

		workspace := db.Workspace{
			Uuid:        uuid.New().String(),
			Name:        uuid.New().String(),
			OwnerPubKey: uuid.New().String(),
			Github:      "https://github.com/bounties",
			Website:     "https://www.bountieswebsite.com",
			Description: "Workspace Bounties Description",
		}
		db.TestDB.CreateOrEditWorkspace(workspace)
		bounty := db.NewBounty{
			Type:          "coding",
			Title:         "existing bounty",
			Description:   "existing bounty description",
			WorkspaceUuid: workspace.Uuid,
			OwnerID:       "workspace-user",
			Price:         2000,
		}

		db.TestDB.CreateOrEditBounty(bounty)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)
		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/bounties/"+workspace.Uuid+"/count/", nil)
		if err != nil {
			t.Fatal(err)
		}

		fetchedWorkspace := db.TestDB.GetWorkspaceByUuid(workspace.Uuid)
		workspace.ID = fetchedWorkspace.ID

		fetchedBounty := db.TestDB.GetWorkspaceBounties(req, bounty.WorkspaceUuid)
		bounty.ID = fetchedBounty[0].ID

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		assert.Equal(t, expectedCount, len(fetchedBounty))
		assert.Equal(t, workspace, fetchedWorkspace)
		assert.Equal(t, bounty, fetchedBounty[0])
	})
}

func TestAddUserRoles(t *testing.T) {

	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewWorkspaceHandler(db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "alias",
		UniqueName:  "unique_name",
		OwnerPubKey: "pubkey",
		PriceToMeet: 0,
		Description: "description",
	}

	person2 := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "alias2",
		UniqueName:  "unique_name2",
		OwnerPubKey: "pubkey2",
		PriceToMeet: 0,
		Description: "description2",
	}
	db.TestDB.CreateOrEditPerson(person)
	db.TestDB.CreateOrEditPerson(person2)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "workspace_name",
		OwnerPubKey: person2.OwnerPubKey,
		Github:      "gtihub",
		Website:     "website",
		Description: "description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	userRoles := []db.WorkspaceUserRoles{
		db.WorkspaceUserRoles{
			WorkspaceUuid: workspace.Uuid,
			OwnerPubKey:   person2.OwnerPubKey,
			Role:          "ADD BOUNTY",
		},
	}

	workspaceUser := db.WorkspaceUsers{
		OwnerPubKey:   person2.OwnerPubKey,
		OrgUuid:       workspace.Uuid,
		WorkspaceUuid: workspace.Uuid,
	}

	db.TestDB.CreateWorkspaceUser(workspaceUser)

	t.Run("Should test that when the right conditions are met a user can be added to a workspace", func(t *testing.T) {
		handlerUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}
		oHandler.userHasAccess = handlerUserHasAccess

		ctx := context.WithValue(context.Background(), auth.ContextKey, "pub-key")

		requestBody, _ := json.Marshal(userRoles)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)
		rctx.URLParams.Add("user", person2.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/users/role/"+workspace.Uuid+"/"+person2.OwnerPubKey, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		fetchedWorkspaceUser := db.TestDB.GetWorkspaceUser(person2.OwnerPubKey, workspace.Uuid)

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.AddUserRoles).ServeHTTP(rr, req)

		fetchedUserRole := db.TestDB.GetUserRoles(workspace.Uuid, person2.OwnerPubKey)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, person2.OwnerPubKey, fetchedWorkspaceUser.OwnerPubKey)
		assert.Equal(t, userRoles[0].Role, fetchedUserRole[0].Role)

	})

	t.Run("Should test that when an unauthorized user hits the endpoint it returns a 401 error", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		requestBody, _ := json.Marshal(userRoles)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		rctx.URLParams.Add("user", person2.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodPost, "/users/role/"+workspaceUUID+"/"+person2.OwnerPubKey, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.AddUserRoles).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that if a user or uuid parameters are not passed it returns a 401 error", func(t *testing.T) {

		requestBody, _ := json.Marshal(userRoles)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", "")
		rctx.URLParams.Add("user", "")
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodPost, "/users/role/"+""+"/"+"", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.AddUserRoles).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that when a wrong body data is sent to the endpoint it returns a 406 error", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		InvalidJson := []byte(`{"key": "value"`)
		requestBody, _ := json.Marshal(InvalidJson)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		rctx.URLParams.Add("user", person2.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodPost, "/users/role/"+workspaceUUID+"/"+person2.OwnerPubKey, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.AddUserRoles).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("Should test that if a user is not the creator of the workspace or does not have an ADD USER ROLE it returns a 401 error", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		handlerUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
			return false
		}
		oHandler.userHasAccess = handlerUserHasAccess
		userRoles[0].OwnerPubKey = person.OwnerPubKey
		requestBody, _ := json.Marshal(userRoles)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		rctx.URLParams.Add("user", person.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodPost, "/users/role/"+workspaceUUID+"/"+person.OwnerPubKey, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.AddUserRoles).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that when the pubkey from URL param does not match the pubkey from JWT AUTH claims it returns a 401 error", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		ctx := context.WithValue(context.Background(), auth.ContextKey, "mismatching_pubkey")

		requestBody, _ := json.Marshal(userRoles)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		rctx.URLParams.Add("user", person2.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/users/role/"+workspaceUUID+"/"+person2.OwnerPubKey, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.AddUserRoles).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that if user doesn't exists in workspace it returns a 401 error", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		handlerUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}
		oHandler.userHasAccess = handlerUserHasAccess
		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

		userRoles[0].OwnerPubKey = person.OwnerPubKey
		requestBody, _ := json.Marshal(userRoles)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		rctx.URLParams.Add("user", person.OwnerPubKey)

		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/users/role/"+workspaceUUID+"/"+person.OwnerPubKey, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.AddUserRoles).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

}

func TestGetUserRoles(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	oHandler := NewWorkspaceHandler(db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "alias",
		UniqueName:  "unique_name",
		OwnerPubKey: "pubkey",
		PriceToMeet: 0,
		Description: "description",
	}

	person2 := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "alias2",
		UniqueName:  "unique_name2",
		OwnerPubKey: "pubkey2",
		PriceToMeet: 0,
		Description: "description2",
	}
	db.TestDB.CreateOrEditPerson(person)
	db.TestDB.CreateOrEditPerson(person2)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        uuid.New().String(),
		OwnerPubKey: person2.OwnerPubKey,
		Github:      "gtihub",
		Website:     "website",
		Description: "description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	userRoles := []db.WorkspaceUserRoles{
		db.WorkspaceUserRoles{
			WorkspaceUuid: workspace.Uuid,
			OwnerPubKey:   person2.OwnerPubKey,
			Role:          "ADD BOUNTY",
		},
	}

	db.TestDB.CreateUserRoles(userRoles, workspace.Uuid, person2.OwnerPubKey)

	t.Run("Should test that the ADD BOUNTY role is returned for person2 from the API call response and the API response array length is 1", func(t *testing.T) {

		ctx := context.WithValue(context.Background(), auth.ContextKey, "pub-key")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)
		rctx.URLParams.Add("user", person2.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/users/role/"+workspace.Uuid+"/"+person2.OwnerPubKey, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetUserRoles).ServeHTTP(rr, req)

		var returnedUserRole []db.WorkspaceUserRoles
		err = json.Unmarshal(rr.Body.Bytes(), &returnedUserRole)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, userRoles[0].Role, returnedUserRole[0].Role)
		assert.Equal(t, 1, len(returnedUserRole))

	})
}

func TestCreateWorkspaceUser(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewWorkspaceHandler(db.TestDB)

	person := db.Person{
		Uuid:        "uuid",
		OwnerAlias:  "alias",
		UniqueName:  "unique_name",
		OwnerPubKey: "pubkey",
		PriceToMeet: 0,
		Description: "description",
	}

	person2 := db.Person{
		Uuid:        "uuid2",
		OwnerAlias:  "alias2",
		UniqueName:  "unique_name2",
		OwnerPubKey: "pubkey2",
		PriceToMeet: 0,
		Description: "description2",
	}
	db.TestDB.CreateOrEditPerson(person)
	db.TestDB.CreateOrEditPerson(person2)

	workspace := db.Workspace{
		Uuid:        "workspace_uuid",
		Name:        "workspace_name",
		OwnerPubKey: "person.OwnerPubkey",
		Github:      "gtihub",
		Website:     "website",
		Description: "description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

	workspaceUser := db.WorkspaceUsers{
		OwnerPubKey:   person.OwnerPubKey,
		OrgUuid:       workspace.Uuid,
		WorkspaceUuid: workspace.Uuid,
	}

	workspaceUserData := db.WorkspaceUsersData{
		OrgUuid:       workspace.Uuid,
		WorkspaceUuid: workspace.Uuid,
		Person:        person,
	}
	db.TestDB.DeleteWorkspaceUser(workspaceUserData, workspace.Uuid)

	workspaceUserData.Person = person2
	db.TestDB.DeleteWorkspaceUser(workspaceUserData, workspace.Uuid)

	t.Run("Should test that when an unauthorized user hits the endpoint it returns a 401 error", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		requestBody, _ := json.Marshal(workspaceUser)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodPost, "/users/"+workspaceUUID, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.CreateWorkspaceUser).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that when a wrong body data is sent to the endpoint it returns a 406 error", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		InvalidJson := []byte(`{"key": "value"`)
		requestBody, _ := json.Marshal(InvalidJson)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/users/"+workspaceUUID, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.CreateWorkspaceUser).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("Should test that if a user is not the creator of the workspace or does not have an ADD USER ROLE it returns a 401 error", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		handlerUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
			return false
		}
		oHandler.userHasAccess = handlerUserHasAccess

		requestBody, _ := json.Marshal(workspaceUser)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/users/"+workspaceUUID, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.CreateWorkspaceUser).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that when the pubkey from URL param does not match the pubkey from JWT AUTH claims it returns a 401 error", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		ctx := context.WithValue(context.Background(), auth.ContextKey, "mismatching_pubkey")

		requestBody, _ := json.Marshal(workspaceUser)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/users/"+workspaceUUID, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.CreateWorkspaceUser).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that a user cannot add themselves it should return a 401 error", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		requestBody, _ := json.Marshal(workspaceUser)
		ctx := context.WithValue(context.Background(), auth.ContextKey, workspaceUser.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/users/"+workspaceUUID, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.CreateWorkspaceUser).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that Cannot add workspace admin as a user it should return a 401 error", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		workspaceUser.OwnerPubKey = workspace.OwnerPubKey
		requestBody, _ := json.Marshal(workspaceUser)
		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/users/"+workspaceUUID, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.CreateWorkspaceUser).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that if user doesn't exists in people it returns a 401 error", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		handlerUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}
		oHandler.userHasAccess = handlerUserHasAccess

		workspaceUser.OwnerPubKey = "OwnerPubKey"
		requestBody, _ := json.Marshal(workspaceUser)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/users/"+workspaceUUID, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.CreateWorkspaceUser).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that when the right conditions are met a user can be added to a workspace", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		handlerUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}
		oHandler.userHasAccess = handlerUserHasAccess

		workspaceUser.OwnerPubKey = person.OwnerPubKey
		requestBody, _ := json.Marshal(workspaceUser)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/users/"+workspaceUUID, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.CreateWorkspaceUser).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should test that when the right conditions are met another user can be added to a workspace", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		handlerUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}
		oHandler.userHasAccess = handlerUserHasAccess

		workspaceUser.OwnerPubKey = person2.OwnerPubKey
		requestBody, _ := json.Marshal(workspaceUser)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/users/"+workspaceUUID, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.CreateWorkspaceUser).ServeHTTP(rr, req)

		updatedWorkspaceUsers, err := db.TestDB.GetWorkspaceUsers(workspaceUUID)
		if err != nil {
			t.Fatal(err)
		}

		updatedWorkspaceUser := db.TestDB.GetWorkspaceUser(person2.OwnerPubKey, workspaceUUID)

		assert.Equal(t, 2, len(updatedWorkspaceUsers))
		assert.Equal(t, person2.OwnerPubKey, updatedWorkspaceUser.OwnerPubKey)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should test that an existing user cannot be added to the workspace it returns a 401 error", func(t *testing.T) {
		workspaceUUID := workspace.Uuid

		handlerUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}
		oHandler.userHasAccess = handlerUserHasAccess

		workspaceUser.OwnerPubKey = person.OwnerPubKey
		requestBody, _ := json.Marshal(workspaceUser)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspaceUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/users/"+workspaceUUID, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.CreateWorkspaceUser).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestGetWorkspaceUsers(t *testing.T) {

}

func TestGetUserDropdownWorkspaces(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewWorkspaceHandler(db.TestDB)

	db.TestDB.DeleteWorkspace()

	person := db.Person{
		Uuid:        "uuid",
		OwnerAlias:  "alias",
		UniqueName:  "unique_name",
		OwnerPubKey: "pubkey",
		PriceToMeet: 0,
		Description: "description",
	}

	person2 := db.Person{
		Uuid:        "uuid2",
		OwnerAlias:  "alias2",
		UniqueName:  "unique_name2",
		OwnerPubKey: "pubkey2",
		PriceToMeet: 0,
		Description: "description2",
	}
	db.TestDB.CreateOrEditPerson(person)
	db.TestDB.CreateOrEditPerson(person2)

	workspace := db.Workspace{
		Uuid:        "workspace_uuid",
		Name:        "workspace_name",
		OwnerPubKey: "person.OwnerPubkey",
		Github:      "gtihub",
		Website:     "website",
		Description: "description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	roles := []db.WorkspaceUserRoles{
		db.WorkspaceUserRoles{
			WorkspaceUuid: workspace.Uuid,
			OwnerPubKey:   person2.OwnerPubKey,
			Role:          "ADD BOUNTY",
		},
		db.WorkspaceUserRoles{
			WorkspaceUuid: workspace.Uuid,
			OwnerPubKey:   person2.OwnerPubKey,
			Role:          "UPDATE BOUNTY",
		},
		db.WorkspaceUserRoles{
			WorkspaceUuid: workspace.Uuid,
			OwnerPubKey:   person2.OwnerPubKey,
			Role:          "DELETE BOUNTY",
		},
		db.WorkspaceUserRoles{
			WorkspaceUuid: workspace.Uuid,
			OwnerPubKey:   person2.OwnerPubKey,
			Role:          "PAY BOUNTY",
		},
	}

	db.TestDB.CreateUserRoles(roles, workspace.Uuid, person2.OwnerPubKey)

	dbPerson := db.TestDB.GetPersonByUuid(person2.Uuid)

	ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

	t.Run("should return user dropdown workspaces", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("userId", strconv.Itoa(int(dbPerson.ID)))
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/user/dropdown/"+strconv.Itoa(int(dbPerson.ID)), nil)
		if err != nil {
			t.Fatal(err)
		}

		handler := http.HandlerFunc(oHandler.GetUserDropdownWorkspaces)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseWorkspaces []db.Workspace
		err = json.Unmarshal(rr.Body.Bytes(), &responseWorkspaces)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotEmpty(t, responseWorkspaces)
		assert.Equal(t, workspace.Uuid, responseWorkspaces[0].Uuid)
		assert.Equal(t, workspace.Name, responseWorkspaces[0].Name)
		assert.Equal(t, workspace.OwnerPubKey, responseWorkspaces[0].OwnerPubKey)
		assert.Equal(t, workspace.Github, responseWorkspaces[0].Github)
		assert.Equal(t, workspace.Website, responseWorkspaces[0].Website)
		assert.Equal(t, workspace.Description, responseWorkspaces[0].Description)
	})
}

func TestCreateOrEditWorkspaceRepository(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	oHandler := NewWorkspaceHandler(db.TestDB)

	t.Run("should return error if a user is not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspaceRepository)

		bodyJson := []byte(`{"key": "value"}`)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/repositories", bytes.NewReader(bodyJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspaceRepository)

		invalidJson := []byte(`{"key": "value"`)

		ctx := context.WithValue(context.Background(), auth.ContextKey, "pub-key")

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/repositories", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("should return error if a Workspace UUID that does not exist Is passed to the API body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspaceRepository)

		workspace := db.Workspace{
			Uuid:        uuid.New().String(),
			Name:        uuid.New().String(),
			OwnerPubKey: "workspace_owner_bounties_pubkey",
			Github:      "https://github.com/bounties",
			Website:     "https://www.bountieswebsite.com",
			Description: "Workspace Bounties Description",
		}
		db.TestDB.CreateOrEditWorkspace(workspace)

		repository := db.WorkspaceRepositories{
			Uuid:          uuid.New().String(),
			WorkspaceUuid: "wrongid",
			Name:          "workspacerepo",
			Url:           "https://github.com/bounties",
		}

		db.TestDB.CreateOrEditWorkspaceRepository(repository)
		requestBody, _ := json.Marshal(repository)

		ctx := context.WithValue(context.Background(), auth.ContextKey, "pub-key")

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/repositories", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("user should be able to add a workspace repository when the right conditions are met", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspaceRepository)

		workspace := db.Workspace{
			Uuid:        uuid.New().String(),
			Name:        uuid.New().String(),
			OwnerPubKey: "workspace_owner_bounties_pubkey",
			Github:      "https://github.com/bounties",
			Website:     "https://www.bountieswebsite.com",
			Description: "Workspace Bounties Description",
		}
		db.TestDB.CreateOrEditWorkspace(workspace)

		repository := db.WorkspaceRepositories{
			Uuid:          uuid.New().String(),
			WorkspaceUuid: workspace.Uuid,
			Name:          "workspacerepo",
			Url:           "https://github.com/bounties",
		}

		db.TestDB.CreateOrEditWorkspaceRepository(repository)
		requestBody, _ := json.Marshal(repository)

		ctx := context.WithValue(context.Background(), auth.ContextKey, "pub-key")

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/repositories", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		getWorkspaceRepo := db.TestDB.GetWorkspaceRepositorByWorkspaceUuid(workspace.Uuid)

		handler.ServeHTTP(rr, req)

		var returnedWorkspaceRepo db.WorkspaceRepositories
		err = json.Unmarshal(rr.Body.Bytes(), &returnedWorkspaceRepo)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		// Assert that the workspace repository is created by using the GetWorkspaceRepositorByWorkspaceUuid function
		assert.Equal(t, repository.Name, getWorkspaceRepo[0].Name)
		assert.Equal(t, repository.Url, getWorkspaceRepo[0].Url)
		// Assert that the Name and Url  of the repository returned matches what was sent in the API body.
		assert.Equal(t, repository.Name, returnedWorkspaceRepo.Name)
		assert.Equal(t, repository.Url, returnedWorkspaceRepo.Url)
	})

}

func TestGetWorkspaceRepositorByWorkspaceUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewWorkspaceHandler(db.TestDB)

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

	// Create a workspace repository
	repository := db.WorkspaceRepositories{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-repo",
		Url:           "https://github.com/test-repo",
	}
	db.TestDB.CreateOrEditWorkspaceRepository(repository)

	t.Run("should return error if user is not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceRepositorByWorkspaceUuid)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/repositories/"+workspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return workspace repositories if user is authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceRepositorByWorkspaceUuid)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)
		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/repositories/"+workspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var returnedRepos []db.WorkspaceRepositories
		err = json.Unmarshal(rr.Body.Bytes(), &returnedRepos)
		assert.NoError(t, err)
		assert.Len(t, returnedRepos, 1)
		assert.Equal(t, repository.Name, returnedRepos[0].Name)
		assert.Equal(t, repository.Url, returnedRepos[0].Url)
	})
}

func TestGetWorkspaceRepoByWorkspaceUuidAndRepoUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewWorkspaceHandler(db.TestDB)

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

	// Create a workspace repository
	repository := db.WorkspaceRepositories{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-repo",
		Url:           "https://github.com/test-repo",
	}
	db.TestDB.CreateOrEditWorkspaceRepository(repository)

	t.Run("should return error if user is not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceRepoByWorkspaceUuidAndRepoUuid)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		rctx.URLParams.Add("uuid", repository.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/repositories/"+workspace.Uuid+"/repository/"+repository.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return workspace repository if user is authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceRepoByWorkspaceUuidAndRepoUuid)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		rctx.URLParams.Add("uuid", repository.Uuid)
		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/repositories/"+workspace.Uuid+"/repository/"+repository.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var returnedRepos db.WorkspaceRepositories
		err = json.Unmarshal(rr.Body.Bytes(), &returnedRepos)
		assert.NoError(t, err)
		assert.Equal(t, repository.Name, returnedRepos.Name)
		assert.Equal(t, repository.Url, returnedRepos.Url)
	})
}

func TestDeleteWorkspaceRepository(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewWorkspaceHandler(db.TestDB)

	person := db.Person{
		Uuid:        "uuid",
		OwnerAlias:  "alias",
		UniqueName:  "unique_name",
		OwnerPubKey: "pubkey",
		PriceToMeet: 0,
		Description: "description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        "workspace_uuid",
		Name:        "workspace_name",
		OwnerPubKey: "person.OwnerPubkey",
		Github:      "gtihub",
		Website:     "website",
		Description: "description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	repository := db.WorkspaceRepositories{
		Uuid:          "repo_uuid",
		WorkspaceUuid: workspace.Uuid,
		Name:          "repo_name",
		Url:           "repo_url",
	}
	db.TestDB.CreateOrEditWorkspaceRepository(repository)

	ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

	t.Run("Should test that it throws a 401 error if a user is not authorized", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", repository.Uuid)
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodDelete, "/"+workspace.Uuid+"/repository/"+repository.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.DeleteWorkspaceRepository).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that the repository is deleted after the Delete API request is successful", func(t *testing.T) {
		workspaceRepo, err := db.TestDB.GetWorkspaceRepoByWorkspaceUuidAndRepoUuid(workspace.Uuid, repository.Uuid)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotEmpty(t, workspaceRepo)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", repository.Uuid)
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/"+workspace.Uuid+"/repository/"+repository.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.DeleteWorkspaceRepository).ServeHTTP(rr, req)

		_, err = db.TestDB.GetWorkspaceRepoByWorkspaceUuidAndRepoUuid(workspace.Uuid, repository.Uuid)
		assert.Error(t, err)
		assert.Equal(t, "workspace repository not found", err.Error())
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestCreateOrEditWorkspaceCodeGraph(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	oHandler := NewWorkspaceHandler(db.TestDB)

	t.Run("should return error if a user is not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspaceCodeGraph)

		bodyJson := []byte(`{"key": "value"}`)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/codegraph", bytes.NewReader(bodyJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspaceCodeGraph)

		invalidJson := []byte(`{"key": "value"`)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "pub-key")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/codegraph", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("should return error if a Workspace UUID that does not exist is passed to the API body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspaceCodeGraph)

		workspace := db.Workspace{
			Uuid:        uuid.New().String(),
			Name:        uuid.New().String(),
			OwnerPubKey: "workspace_owner_pubkey",
			Github:      "https://github.com/test",
			Website:     "https://www.test.com",
			Description: "Test Description",
		}
		db.TestDB.CreateOrEditWorkspace(workspace)

		codeGraph := db.WorkspaceCodeGraph{
			Uuid:          uuid.New().String(),
			WorkspaceUuid: "wrongid",
			Name:          "testgraph",
			Url:           "https://github.com/test/graph",
		}

		requestBody, _ := json.Marshal(codeGraph)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "pub-key")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/codegraph", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("user should be able to add a workspace code graph when the right conditions are met", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspaceCodeGraph)

		workspace := db.Workspace{
			Uuid:        uuid.New().String(),
			Name:        uuid.New().String(),
			OwnerPubKey: "workspace_owner_pubkey",
			Github:      "https://github.com/test",
			Website:     "https://www.test.com",
			Description: "Test Description",
		}
		db.TestDB.CreateOrEditWorkspace(workspace)

		codeGraph := db.WorkspaceCodeGraph{
			Uuid:          uuid.New().String(),
			WorkspaceUuid: workspace.Uuid,
			Name:          "testgraph",
			Url:           "https://github.com/test/graph",
		}

		requestBody, _ := json.Marshal(codeGraph)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "pub-key")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/codegraph", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var returnedCodeGraph db.WorkspaceCodeGraph
		err = json.Unmarshal(rr.Body.Bytes(), &returnedCodeGraph)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, codeGraph.Name, returnedCodeGraph.Name)
		assert.Equal(t, codeGraph.Url, returnedCodeGraph.Url)
	})
}

func TestGetWorkspaceCodeGraphByUUID(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewWorkspaceHandler(db.TestDB)

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
		Website:     "https://www.test.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	codeGraph := db.WorkspaceCodeGraph{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-graph",
		Url:           "https://github.com/test/graph",
	}
	db.TestDB.CreateOrEditCodeGraph(codeGraph)

	t.Run("should return error if user is not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceCodeGraphByUUID)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", codeGraph.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodGet, "/codegraph/"+codeGraph.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return code graph if user is authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceCodeGraphByUUID)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", codeGraph.Uuid)
		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/codegraph/"+codeGraph.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var returnedCodeGraph db.WorkspaceCodeGraph
		err = json.Unmarshal(rr.Body.Bytes(), &returnedCodeGraph)
		assert.NoError(t, err)
		assert.Equal(t, codeGraph.Name, returnedCodeGraph.Name)
		assert.Equal(t, codeGraph.Url, returnedCodeGraph.Url)
	})
}

func TestGetCodeGraphsByWorkspaceUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewWorkspaceHandler(db.TestDB)

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
		Website:     "https://www.test.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	codeGraph := db.WorkspaceCodeGraph{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-graph",
		Url:           "https://github.com/test/graph",
	}
	db.TestDB.CreateOrEditCodeGraph(codeGraph)

	t.Run("should return error if user is not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetCodeGraphsByWorkspaceUuid)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodGet, "/"+workspace.Uuid+"/codegraph", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return workspace code graphs if user is authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetCodeGraphsByWorkspaceUuid)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/"+workspace.Uuid+"/codegraph", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var returnedCodeGraphs []db.WorkspaceCodeGraph
		err = json.Unmarshal(rr.Body.Bytes(), &returnedCodeGraphs)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(returnedCodeGraphs))
		assert.Equal(t, codeGraph.Name, returnedCodeGraphs[0].Name)
		assert.Equal(t, codeGraph.Url, returnedCodeGraphs[0].Url)
	})
}

func TestDeleteWorkspaceCodeGraph(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewWorkspaceHandler(db.TestDB)

	person := db.Person{
		Uuid:        "uuid",
		OwnerAlias:  "alias",
		UniqueName:  "unique_name",
		OwnerPubKey: "pubkey",
		PriceToMeet: 0,
		Description: "description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        "workspace_uuid",
		Name:        "workspace_name",
		OwnerPubKey: "person.OwnerPubkey",
		Github:      "github",
		Website:     "website",
		Description: "description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	codeGraph := db.WorkspaceCodeGraph{
		Uuid:          "graph_uuid",
		WorkspaceUuid: workspace.Uuid,
		Name:          "graph_name",
		Url:           "graph_url",
	}
	db.TestDB.CreateOrEditCodeGraph(codeGraph)

	ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

	t.Run("Should test that it throws a 401 error if a user is not authorized", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", codeGraph.Uuid)
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodDelete, "/"+workspace.Uuid+"/codegraph/"+codeGraph.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.DeleteWorkspaceCodeGraph).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that the code graph is deleted after the Delete API request is successful", func(t *testing.T) {
		existingGraph, err := db.TestDB.GetCodeGraphByUUID(codeGraph.Uuid)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotEmpty(t, existingGraph)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", codeGraph.Uuid)
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodDelete, "/"+workspace.Uuid+"/codegraph/"+codeGraph.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.DeleteWorkspaceCodeGraph).ServeHTTP(rr, req)

		_, err = db.TestDB.GetCodeGraphByUUID(codeGraph.Uuid)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error())
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestPollBudgetInvoices(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	config.IsV2Payment = true
	config.V2BotUrl = "http://v2-bot-url.com"
	config.V2BotToken = "v2-bot-token"

	// create a user
	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerPubKey: "test_user_poll_budget",
		OwnerAlias:  "test_user_poll_update_budget",
		Description: "test_user_poll_update_budget_description",
	}

	db.TestDB.CreateOrEditPerson(person)

	t.Run("Should return 401 if the user is not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		oHandler := NewWorkspaceHandler(db.TestDB)
		handler := http.HandlerFunc(oHandler.PollBudgetInvoices)

		req, err := http.NewRequest(http.MethodGet, "/poll/invoices/"+workspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should update the budget if the invoice if a single invoice is paid invoice is paid", func(t *testing.T) {
		now := time.Now()

		amount := 50000

		// Create a new workspace
		uuid := uuid.New()

		randomWorkspaceName := fmt.Sprintf("Test Poll Workspace Budget %d", rand.Intn(1000))
		workspace := db.Workspace{
			OwnerPubKey: person.OwnerPubKey,
			Uuid:        uuid.String(),
			Name:        randomWorkspaceName,
		}

		db.TestDB.CreateOrEditWorkspace(workspace)

		randomPaymentRequest := fmt.Sprintf("test_update_budget_payment_request_%d", rand.Intn(1000))
		// create invoice
		invoice := db.NewInvoiceList{
			WorkspaceUuid:  workspace.Uuid,
			PaymentRequest: randomPaymentRequest,
			Status:         false,
			OwnerPubkey:    person.OwnerPubKey,
			Created:        &now,
			Type:           "BUDGET",
		}

		db.TestDB.AddInvoice(invoice)

		mockGetLightningInvoice := func(payment_request string) (db.InvoiceResult, db.InvoiceError) {
			return db.InvoiceResult{
				Success: true,
				Response: db.InvoiceCheckResponse{
					Settled:         true,
					Amount:          strconv.Itoa(amount),
					Payment_request: invoice.PaymentRequest,
				},
			}, db.InvoiceError{}
		}

		// create paymentHistory
		paymentHistory := db.NewPaymentHistory{
			WorkspaceUuid: workspace.Uuid,
			Amount:        uint(amount),
			PaymentStatus: db.PaymentComplete,
			PaymentType:   db.Deposit,
			SenderPubKey:  person.OwnerPubKey,
			Created:       &now,
		}

		db.TestDB.AddPaymentHistory(paymentHistory)

		ctx := context.Background()
		oHandler := NewWorkspaceHandler(db.TestDB)

		oHandler.getLightningInvoice = mockGetLightningInvoice
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.PollBudgetInvoices)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)

		authorizedCtx := context.WithValue(ctx, auth.ContextKey, invoice.OwnerPubkey)
		req, err := http.NewRequestWithContext(context.WithValue(authorizedCtx, chi.RouteCtxKey, rctx), http.MethodGet, "/poll/invoices/"+workspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		workspaceBudget := db.TestDB.GetWorkspaceBudget(workspace.Uuid)
		assert.Equal(t, workspaceBudget.TotalBudget, uint(amount))
	})

	t.Run("Should update the budget if the invoice if a multiple invoice is paid invoice is paid ", func(t *testing.T) {
		now := time.Now()

		amount := 10000

		randomWorkspaceName := fmt.Sprintf("Test Poll Multiple Workspace Budget %d", rand.Intn(1000))
		workspace := db.Workspace{
			OwnerPubKey: person.OwnerPubKey,
			Uuid:        uuid.New().String(),
			Name:        randomWorkspaceName,
		}

		db.TestDB.CreateOrEditWorkspace(workspace)

		for i := 0; i < 3; i++ {
			randomPaymentRequest := fmt.Sprintf("test_update_budget_payment_request_%d%d", rand.Intn(amount), i)
			// create invoice
			invoice := db.NewInvoiceList{
				WorkspaceUuid:  workspace.Uuid,
				PaymentRequest: randomPaymentRequest,
				Status:         false,
				OwnerPubkey:    person.OwnerPubKey,
				Created:        &now,
				Type:           "BUDGET",
			}
			db.TestDB.AddInvoice(invoice)
		}

		for i := 0; i < 3; i++ {
			// create paymentHistory
			paymentHistory := db.NewPaymentHistory{
				WorkspaceUuid: workspace.Uuid,
				Amount:        uint(amount),
				PaymentStatus: db.PaymentComplete,
				PaymentType:   db.Deposit,
				SenderPubKey:  person.OwnerPubKey,
				Created:       &now,
			}
			db.TestDB.AddPaymentHistory(paymentHistory)
		}

		mockGetLightningInvoice := func(payment_request string) (db.InvoiceResult, db.InvoiceError) {
			return db.InvoiceResult{
				Success: true,
				Response: db.InvoiceCheckResponse{
					Settled:         true,
					Amount:          strconv.Itoa(amount),
					Payment_request: "",
				},
			}, db.InvoiceError{}
		}

		ctx := context.Background()
		oHandler := NewWorkspaceHandler(db.TestDB)

		oHandler.getLightningInvoice = mockGetLightningInvoice

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.PollBudgetInvoices)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)

		authorizedCtx := context.WithValue(ctx, auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(authorizedCtx, chi.RouteCtxKey, rctx), http.MethodGet, "/poll/invoices/"+workspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		workspaceBudget := db.TestDB.GetWorkspaceBudget(workspace.Uuid)
		assert.Equal(t, workspaceBudget.TotalBudget, uint(amount*3))
	})
}

func TestPollUserWorkspacesBudget(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	config.IsV2Payment = true
	config.V2BotUrl = "http://v2-bot-url.com"
	config.V2BotToken = "v2-bot-token"

	// create a user
	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerPubKey: "test_user_poll_budget",
		OwnerAlias:  "test_user_poll_update_budget",
		Description: "test_user_poll_update_budget_description",
	}

	db.TestDB.CreateOrEditPerson(person)

	// create second user
	person2 := db.Person{
		Uuid:        uuid.New().String(),
		OwnerPubKey: "test_user_poll_budget_2",
		OwnerAlias:  "test_user_poll_update_budget_2",
		Description: "test_user_poll_update_budget_description_2",
	}
	db.TestDB.CreateOrEditPerson(person2)

	t.Run("Should return 401 if the user is not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		oHandler := NewWorkspaceHandler(db.TestDB)
		handler := http.HandlerFunc(oHandler.PollUserWorkspacesBudget)

		req, err := http.NewRequest(http.MethodGet, "/poll/user/invoices", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should update the user workspaces budget if the invoice if a single invoice is paid invoice is paid", func(t *testing.T) {
		now := time.Now()

		amount := 50000

		// Create a new workspace
		uuid := uuid.New()

		randomWorkspaceName := fmt.Sprintf("Test Poll User Workspace Budget %d", rand.Intn(1000))
		workspace := db.Workspace{
			OwnerPubKey: person.OwnerPubKey,
			Uuid:        uuid.String(),
			Name:        randomWorkspaceName,
		}

		db.TestDB.CreateOrEditWorkspace(workspace)

		// add user 2 to workspace users
		db.TestDB.CreateWorkspaceUser(db.WorkspaceUsers{
			WorkspaceUuid: workspace.Uuid,
			OwnerPubKey:   person2.OwnerPubKey,
		})

		randomPaymentRequest := fmt.Sprintf("test_update_budget_payment_request_%d", rand.Intn(1000))
		// create invoice
		invoice := db.NewInvoiceList{
			WorkspaceUuid:  workspace.Uuid,
			PaymentRequest: randomPaymentRequest,
			Status:         false,
			OwnerPubkey:    person2.OwnerPubKey,
			Created:        &now,
			Type:           "BUDGET",
		}

		db.TestDB.AddInvoice(invoice)

		mockGetLightningInvoice := func(payment_request string) (db.InvoiceResult, db.InvoiceError) {
			return db.InvoiceResult{
				Success: true,
				Response: db.InvoiceCheckResponse{
					Settled:         true,
					Amount:          strconv.Itoa(amount),
					Payment_request: invoice.PaymentRequest,
				},
			}, db.InvoiceError{}
		}

		mockGetAllUserWorkspaces := func(pubKeyFromAuth string) []db.Workspace {
			return []db.Workspace{workspace}
		}

		// create paymentHistory
		paymentHistory := db.NewPaymentHistory{
			WorkspaceUuid: workspace.Uuid,
			Amount:        uint(amount),
			PaymentStatus: db.PaymentComplete,
			PaymentType:   db.Deposit,
			SenderPubKey:  person.OwnerPubKey,
			Created:       &now,
		}

		db.TestDB.AddPaymentHistory(paymentHistory)

		ctx := context.Background()
		oHandler := NewWorkspaceHandler(db.TestDB)

		oHandler.getLightningInvoice = mockGetLightningInvoice
		oHandler.getAllUserWorkspaces = mockGetAllUserWorkspaces

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.PollUserWorkspacesBudget)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)

		authorizedCtx := context.WithValue(ctx, auth.ContextKey, person2.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(authorizedCtx, chi.RouteCtxKey, rctx), http.MethodGet, "/poll/user/invoices", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		workspaceBudget := db.TestDB.GetWorkspaceBudget(workspace.Uuid)
		assert.Equal(t, workspaceBudget.TotalBudget, uint(amount))
	})

	t.Run("Should update the user workspaces budget if the invoice if a multiple invoice is paid invoice is paid ", func(t *testing.T) {
		now := time.Now()

		amount := 10000

		randomWorkspaceName := fmt.Sprintf("Test Poll Multiple Workspace Budget %d", rand.Intn(1000))
		workspace := db.Workspace{
			OwnerPubKey: person.OwnerPubKey,
			Uuid:        uuid.New().String(),
			Name:        randomWorkspaceName,
		}

		db.TestDB.CreateOrEditWorkspace(workspace)

		// add user 2 to workspace users
		db.TestDB.CreateWorkspaceUser(db.WorkspaceUsers{
			WorkspaceUuid: workspace.Uuid,
			OwnerPubKey:   person2.OwnerPubKey,
		})

		for i := 0; i < 3; i++ {
			randomPaymentRequest := fmt.Sprintf("test_user_workspace_update_budget_payment_request_%d%d", rand.Intn(amount), i)
			// create invoice
			invoice := db.NewInvoiceList{
				WorkspaceUuid:  workspace.Uuid,
				PaymentRequest: randomPaymentRequest,
				Status:         false,
				OwnerPubkey:    person2.OwnerPubKey,
				Created:        &now,
				Type:           "BUDGET",
			}
			db.TestDB.AddInvoice(invoice)
		}

		for i := 0; i < 3; i++ {
			// create paymentHistory
			paymentHistory := db.NewPaymentHistory{
				WorkspaceUuid: workspace.Uuid,
				Amount:        uint(amount),
				PaymentStatus: db.PaymentComplete,
				PaymentType:   db.Deposit,
				SenderPubKey:  person.OwnerPubKey,
				Created:       &now,
			}
			db.TestDB.AddPaymentHistory(paymentHistory)
		}

		mockGetLightningInvoice := func(payment_request string) (db.InvoiceResult, db.InvoiceError) {
			return db.InvoiceResult{
				Success: true,
				Response: db.InvoiceCheckResponse{
					Settled:         true,
					Amount:          strconv.Itoa(amount),
					Payment_request: "",
				},
			}, db.InvoiceError{}
		}

		mockGetAllUserWorkspaces := func(pubKeyFromAuth string) []db.Workspace {
			return []db.Workspace{workspace}
		}

		ctx := context.Background()
		oHandler := NewWorkspaceHandler(db.TestDB)

		oHandler.getLightningInvoice = mockGetLightningInvoice
		oHandler.getAllUserWorkspaces = mockGetAllUserWorkspaces

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.PollUserWorkspacesBudget)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)

		authorizedCtx := context.WithValue(ctx, auth.ContextKey, person2.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(authorizedCtx, chi.RouteCtxKey, rctx), http.MethodGet, "/poll/user/invoices", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		workspaceBudget := db.TestDB.GetWorkspaceBudget(workspace.Uuid)
		assert.Equal(t, workspaceBudget.TotalBudget, uint(amount*3))
	})
}

func TestHandlerGetFeaturesByWorkspaceUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	oHandler := NewWorkspaceHandler(db.TestDB)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace-" + uuid.New().String(),
		OwnerPubKey: "test-pubkey",
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	t.Run("should return unauthorized when no pubkey provided", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetFeaturesByWorkspaceUuid)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/workspaces/"+workspace.Uuid+"/features", nil)
		if err != nil {
			t.Fatal(err)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return bad request when workspace_uuid is missing", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetFeaturesByWorkspaceUuid)

		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-pubkey")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/workspaces//features", nil)
		if err != nil {
			t.Fatal(err)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", "")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return bad request when workspace_uuid contains invalid characters", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetFeaturesByWorkspaceUuid)

		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-pubkey")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/workspaces/invalid!@#$/features", nil)
		if err != nil {
			t.Fatal(err)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", "invalid!@#$")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return empty array when no features exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetFeaturesByWorkspaceUuid)

		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-pubkey")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/workspaces/"+workspace.Uuid+"/features", nil)
		if err != nil {
			t.Fatal(err)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var features []db.WorkspaceFeatures
		err = json.Unmarshal(rr.Body.Bytes(), &features)
		assert.NoError(t, err)
		assert.Empty(t, features)
	})

	t.Run("should return features with correct bounty counts", func(t *testing.T) {

		feature := db.WorkspaceFeatures{
			Uuid:          uuid.New().String(),
			WorkspaceUuid: workspace.Uuid,
			Name:          "test-feature",
			Brief:         "test brief",
			FeatStatus:    db.ActiveFeature,
		}
		db.TestDB.CreateOrEditFeature(feature)

		phase := db.FeaturePhase{
			Uuid:        uuid.New().String(),
			FeatureUuid: feature.Uuid,
			Name:        "test-phase",
		}
		db.TestDB.CreateOrEditFeaturePhase(phase)

		bounties := []db.NewBounty{
			{
				ID:            1,
				Type:          "coding",
				Title:         "open bounty",
				Description:   "open bounty description",
				WorkspaceUuid: workspace.Uuid,
				PhaseUuid:     phase.Uuid,
				Assignee:      "", // open
				Created:       time.Now().Unix(),
				OwnerID:       "test-owner",
				Completed:     false,
				Paid:          false,
			},
			{
				ID:            2,
				Type:          "coding",
				Title:         "assigned bounty",
				Description:   "assigned bounty description",
				WorkspaceUuid: workspace.Uuid,
				PhaseUuid:     phase.Uuid,
				Assignee:      "test-assignee", // assigned
				Created:       time.Now().Unix(),
				OwnerID:       "test-owner",
				Completed:     false,
				Paid:          false,
			},
			{
				ID:            3,
				Type:          "coding",
				Title:         "completed bounty",
				Description:   "completed bounty description",
				WorkspaceUuid: workspace.Uuid,
				PhaseUuid:     phase.Uuid,
				Assignee:      "test-assignee", // completed
				Created:       time.Now().Unix(),
				OwnerID:       "test-owner",
				Completed:     true,
				Paid:          false,
			},
		}

		for _, bounty := range bounties {
			db.TestDB.CreateOrEditBounty(bounty)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetFeaturesByWorkspaceUuid)

		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-pubkey")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/workspaces/"+workspace.Uuid+"/features", nil)
		if err != nil {
			t.Fatal(err)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var features []db.WorkspaceFeatures
		err = json.Unmarshal(rr.Body.Bytes(), &features)
		assert.NoError(t, err)
		assert.Len(t, features, 1)

		assert.Equal(t, 1, features[0].BountiesCountOpen)
		assert.Equal(t, 1, features[0].BountiesCountAssigned)
		assert.Equal(t, 1, features[0].BountiesCountCompleted)
	})

	t.Run("should return features without bounty counts when no phases exist", func(t *testing.T) {

		feature := db.WorkspaceFeatures{
			Uuid:          uuid.New().String(),
			WorkspaceUuid: workspace.Uuid,
			Name:          "test-feature-no-phases",
			Brief:         "test brief",
			FeatStatus:    db.ActiveFeature,
		}
		db.TestDB.CreateOrEditFeature(feature)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetFeaturesByWorkspaceUuid)

		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-pubkey")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/workspaces/"+workspace.Uuid+"/features", nil)
		if err != nil {
			t.Fatal(err)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var features []db.WorkspaceFeatures
		err = json.Unmarshal(rr.Body.Bytes(), &features)
		assert.NoError(t, err)

		var testFeature *db.WorkspaceFeatures
		for _, f := range features {
			if f.Uuid == feature.Uuid {
				testFeature = &f
				break
			}
		}

		assert.NotNil(t, testFeature)
		assert.Equal(t, 0, testFeature.BountiesCountOpen)
		assert.Equal(t, 0, testFeature.BountiesCountAssigned)
		assert.Equal(t, 0, testFeature.BountiesCountCompleted)
	})

	t.Run("Valid Request with Phases but No Bounties", func(t *testing.T) {
		feature := db.WorkspaceFeatures{
			Uuid:          uuid.New().String(),
			WorkspaceUuid: workspace.Uuid,
			Name:          "feature-with-phases-no-bounties",
			Brief:         "test brief",
			FeatStatus:    db.ActiveFeature,
		}
		db.TestDB.CreateOrEditFeature(feature)

		phase := db.FeaturePhase{
			Uuid:        uuid.New().String(),
			FeatureUuid: feature.Uuid,
			Name:        "test-phase",
		}
		db.TestDB.CreateOrEditFeaturePhase(phase)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetFeaturesByWorkspaceUuid)

		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-pubkey")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/workspaces/"+workspace.Uuid+"/features", nil)
		if err != nil {
			t.Fatal(err)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var features []db.WorkspaceFeatures
		err = json.Unmarshal(rr.Body.Bytes(), &features)
		assert.NoError(t, err)

		var testFeature *db.WorkspaceFeatures
		for _, f := range features {
			if f.Uuid == feature.Uuid {
				testFeature = &f
				break
			}
		}

		assert.NotNil(t, testFeature)
		assert.Equal(t, 0, testFeature.BountiesCountOpen)
		assert.Equal(t, 0, testFeature.BountiesCountAssigned)
		assert.Equal(t, 0, testFeature.BountiesCountCompleted)
	})

	t.Run("Large Number of Features and Phases", func(t *testing.T) {

		for i := 0; i < 10; i++ {
			feature := db.WorkspaceFeatures{
				Uuid:          uuid.New().String(),
				WorkspaceUuid: workspace.Uuid,
				Name:          fmt.Sprintf("bulk-feature-%d", i),
				Brief:         "test brief",
				FeatStatus:    db.ActiveFeature,
			}
			db.TestDB.CreateOrEditFeature(feature)

			for j := 0; j < 5; j++ {
				phase := db.FeaturePhase{
					Uuid:        uuid.New().String(),
					FeatureUuid: feature.Uuid,
					Name:        fmt.Sprintf("phase-%d", j),
				}
				db.TestDB.CreateOrEditFeaturePhase(phase)
			}
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetFeaturesByWorkspaceUuid)

		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-pubkey")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/workspaces/"+workspace.Uuid+"/features", nil)
		if err != nil {
			t.Fatal(err)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var features []db.WorkspaceFeatures
		err = json.Unmarshal(rr.Body.Bytes(), &features)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(features), 10)
	})

	t.Run("Features with Mixed Bounty States", func(t *testing.T) {
		feature := db.WorkspaceFeatures{
			Uuid:          uuid.New().String(),
			WorkspaceUuid: workspace.Uuid,
			Name:          "mixed-states-feature",
			Brief:         "test brief",
			FeatStatus:    db.ActiveFeature,
		}
		db.TestDB.CreateOrEditFeature(feature)

		phase := db.FeaturePhase{
			Uuid:        uuid.New().String(),
			FeatureUuid: feature.Uuid,
			Name:        "test-phase",
		}
		db.TestDB.CreateOrEditFeaturePhase(phase)

		bounties := []db.NewBounty{
			{
				ID:            10,
				Type:          "coding",
				Title:         "mixed state bounty 1",
				Description:   "description",
				WorkspaceUuid: workspace.Uuid,
				PhaseUuid:     phase.Uuid,
				Assignee:      "",
				Created:       time.Now().Unix(),
				OwnerID:       "test-owner",
				Completed:     false,
				Paid:          false,
			},
			{
				ID:            11,
				Type:          "coding",
				Title:         "mixed state bounty 2",
				Description:   "description",
				WorkspaceUuid: workspace.Uuid,
				PhaseUuid:     phase.Uuid,
				Assignee:      "test-assignee",
				Created:       time.Now().Unix(),
				OwnerID:       "test-owner",
				Completed:     true,
				Paid:          true,
			},
		}

		for _, bounty := range bounties {
			db.TestDB.CreateOrEditBounty(bounty)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetFeaturesByWorkspaceUuid)

		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-pubkey")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/workspaces/"+workspace.Uuid+"/features", nil)
		if err != nil {
			t.Fatal(err)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var features []db.WorkspaceFeatures
		err = json.Unmarshal(rr.Body.Bytes(), &features)
		assert.NoError(t, err)

		var testFeature *db.WorkspaceFeatures
		for _, f := range features {
			if f.Uuid == feature.Uuid {
				testFeature = &f
				break
			}
		}

		assert.NotNil(t, testFeature)
		assert.Equal(t, 1, testFeature.BountiesCountOpen)
		assert.Equal(t, 0, testFeature.BountiesCountAssigned)
		assert.Equal(t, 1, testFeature.BountiesCountCompleted)
	})

}
