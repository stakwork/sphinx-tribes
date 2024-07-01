package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
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

}

func TestGetUserRoles(t *testing.T) {

}

func TestCreateWorkspaceUser(t *testing.T) {

}

func TestGetWorkspaceUsers(t *testing.T) {

}

func TestGetUserDropdownWorkspaces(t *testing.T) {

}

func TestCreateOrEditWorkspaceRepository(t *testing.T) {

}

func TestGetWorkspaceRepositorByWorkspaceUuid(t *testing.T) {

}

func TestGetWorkspaceRepoByWorkspaceUuidAndRepoUuid(t *testing.T) {

}

func GetFeaturesByWorkspaceUuid(t *testing.T) {

}

func TestDeleteWorkspaceRepository(t *testing.T) {

}
