package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	mocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUnitCreateOrEditWorkspace(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	oHandler := NewWorkspaceHandler(mockDb)

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
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should trim spaces from organization name", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspace)

		mockDb.On("GetWorkspaceByUuid", mock.AnythingOfType("string")).Return(db.Workspace{}).Once()
		mockDb.On("GetWorkspaceByName", "Abdul").Return(db.Workspace{}).Once()
		mockDb.On("CreateOrEditWorkspace", mock.MatchedBy(func(org db.Workspace) bool {
			return org.Name == "Abdul" && org.Uuid != "" && org.Updated != nil && org.Created != nil
		})).Return(db.Workspace{Name: "Abdul"}, nil).Once()

		jsonInput := []byte(`{"name": " Abdul ", "owner_pubkey": "test-key" ,"description": "Test"}`)
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

		assert.Equal(t, "Abdul", responseOrg.Name)
	})

	t.Run("should successfully add organization if request is valid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditWorkspace)

		mockDb.On("GetWorkspaceByUuid", mock.AnythingOfType("string")).Return(db.Workspace{}).Once()
		mockDb.On("GetWorkspaceByName", "TestWorkspace").Return(db.Workspace{}).Once()
		mockDb.On("CreateOrEditWorkspace", mock.MatchedBy(func(org db.Workspace) bool {
			return org.Name == "TestWorkspace" && org.Uuid != "" && org.Updated != nil && org.Created != nil
		})).Return(db.Workspace{}, nil).Once()

		invalidJson := []byte(`{"name": "TestWorkspace", "owner_pubkey": "test-key" ,"description": "Test"}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
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
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	oHandler := NewWorkspaceHandler(mockDb)

	t.Run("should return error if not authorized", func(t *testing.T) {
		orgUUID := "org-uuid"
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteWorkspace)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodDelete, "/delete/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should set organization fields to null and delete users on successful delete", func(t *testing.T) {
		orgUUID := "org-uuid"

		// Mock expected database interactions
		mockDb.On("GetWorkspaceByUuid", orgUUID).Return(db.Workspace{OwnerPubKey: "test-key"}).Once()
		mockDb.On("UpdateWorkspaceForDeletion", orgUUID).Return(nil).Once()
		mockDb.On("DeleteAllUsersFromWorkspace", orgUUID).Return(nil).Once()
		mockDb.On("ChangeWorkspaceDeleteStatus", orgUUID, true).Return(db.Workspace{Uuid: orgUUID, Deleted: true}).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteWorkspace)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/delete/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("should handle failures in database updates", func(t *testing.T) {
		orgUUID := "org-uuid"

		// Mock database interactions with error
		mockDb.On("GetWorkspaceByUuid", orgUUID).Return(db.Workspace{OwnerPubKey: "test-key"}).Once()
		mockDb.On("UpdateWorkspaceForDeletion", orgUUID).Return(errors.New("update error")).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteWorkspace)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/delete/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("should set organization's deleted column to true", func(t *testing.T) {
		orgUUID := "org-uuid"

		// Mock the database interactions
		mockDb.On("GetWorkspaceByUuid", orgUUID).Return(db.Workspace{OwnerPubKey: "test-key"}).Once()
		mockDb.On("UpdateWorkspaceForDeletion", orgUUID).Return(nil).Once()
		mockDb.On("DeleteAllUsersFromWorkspace", orgUUID).Return(nil).Once()
		mockDb.On("ChangeWorkspaceDeleteStatus", orgUUID, true).Return(db.Workspace{Uuid: orgUUID, Deleted: true}).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteWorkspace)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/delete/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		// Asserting that the response status code is OK
		assert.Equal(t, http.StatusOK, rr.Code)

		// Decoding the response to check if Deleted field is true
		var updatedOrg db.Workspace
		err = json.Unmarshal(rr.Body.Bytes(), &updatedOrg)
		if err != nil {
			t.Fatal(err)
		}
		assert.True(t, updatedOrg.Deleted)

		mockDb.AssertExpectations(t)
	})

	t.Run("should set Website, Github, and Description to empty strings", func(t *testing.T) {
		orgUUID := "org-uuid"

		updatedOrg := db.Workspace{
			Uuid:        orgUUID,
			OwnerPubKey: "test-key",
			Website:     "",
			Github:      "",
			Description: "",
		}

		mockDb.On("GetWorkspaceByUuid", orgUUID).Return(db.Workspace{OwnerPubKey: "test-key"}).Once()
		mockDb.On("UpdateWorkspaceForDeletion", orgUUID).Return(nil).Once()
		mockDb.On("DeleteAllUsersFromWorkspace", orgUUID).Return(nil).Once()
		mockDb.On("ChangeWorkspaceDeleteStatus", orgUUID, true).Return(updatedOrg).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteWorkspace)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/delete/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var returnedOrg db.Workspace
		err = json.Unmarshal(rr.Body.Bytes(), &returnedOrg)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "", returnedOrg.Website)
		assert.Equal(t, "", returnedOrg.Github)
		assert.Equal(t, "", returnedOrg.Description)
		mockDb.AssertExpectations(t)
	})

	t.Run("should delete all users from the organization", func(t *testing.T) {
		orgUUID := "org-uuid"

		// Setting up the expected behavior of the mock database
		mockDb.On("GetWorkspaceByUuid", orgUUID).Return(db.Workspace{OwnerPubKey: "test-key"}).Once()
		mockDb.On("UpdateWorkspaceForDeletion", orgUUID).Return(nil).Once()
		mockDb.On("DeleteAllUsersFromWorkspace", orgUUID).Return(nil).Run(func(args mock.Arguments) {}).Once()
		mockDb.On("ChangeWorkspaceDeleteStatus", orgUUID, true).Return(db.Workspace{Uuid: orgUUID, Deleted: true}).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteWorkspace)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/delete/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		// Asserting that the response status code is as expected
		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
	})
}

func TestGetWorkspaceBounties(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	mockGenerateBountyHandler := func(bounties []db.NewBounty) []db.BountyResponse {
		return []db.BountyResponse{} // Mocked response
	}
	oHandler := NewWorkspaceHandler(mockDb)

	t.Run("Should test that an organization's bounties can be listed without authentication", func(t *testing.T) {
		orgUUID := "valid-uuid"
		oHandler.generateBountyHandler = mockGenerateBountyHandler

		expectedBounties := []db.Bounty{{}, {}} // Mocked response
		mockDb.On("GetWorkspaceBounties", mock.AnythingOfType("*http.Request"), orgUUID).Return(expectedBounties).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceBounties)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/bounties/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return empty array when wrong organization UUID is passed", func(t *testing.T) {
		orgUUID := "wrong-uuid"

		mockDb.On("GetWorkspaceBounties", mock.AnythingOfType("*http.Request"), orgUUID).Return([]db.Bounty{}).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetWorkspaceBounties)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/bounties/"+orgUUID+"?limit=10&sortBy=created&search=test&page=1&resetPage=true", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		// Assert that the response status code is as expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Assert that the response body is an empty array
		assert.Equal(t, "[]\n", rr.Body.String())
	})
}

func TestGetWorkspaceBudget(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	mockUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
		return true
	}
	oHandler := NewWorkspaceHandler(mockDb)

	t.Run("Should test that a 401 is returned when trying to view an organization's budget without a token", func(t *testing.T) {
		orgUUID := "valid-uuid"

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/budget/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetWorkspaceBudget).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that the right workspace budget is returned, if the user is the organization admin or has the ViewReport role", func(t *testing.T) {
		orgUUID := "valid-uuid"
		statusBudget := db.StatusBudget{
			OrgUuid:         orgUUID,
			CurrentBudget:   10000,
			OpenBudget:      1000,
			OpenCount:       10,
			AssignedBudget:  2000,
			AssignedCount:   15,
			CompletedBudget: 3000,
			CompletedCount:  5,
		}

		oHandler.userHasAccess = mockUserHasAccess
		mockDb.On("GetWorkspaceStatusBudget", orgUUID).Return(statusBudget).Once()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/budget/"+orgUUID, nil)
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

		assert.Equal(t, statusBudget, responseBudget)
	})
}

func TestGetWorkspaceBudgetHistory(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	oHandler := NewWorkspaceHandler(mockDb)

	t.Run("Should test that a 401 is returned when trying to view an organization's budget history without a token", func(t *testing.T) {
		orgUUID := "valid-uuid"

		mockUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
			return false
		}
		oHandler.userHasAccess = mockUserHasAccess

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/budget/history/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetWorkspaceBudgetHistory).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that the right budget history is returned, if the user is the organization admin or has the ViewReport role", func(t *testing.T) {
		orgUUID := "valid-uuid"
		expectedBudgetHistory := []db.BudgetHistoryData{
			{BudgetHistory: db.BudgetHistory{ID: 1, OrgUuid: orgUUID, Created: nil, Updated: nil}, SenderName: "Sender1"},
			{BudgetHistory: db.BudgetHistory{ID: 2, OrgUuid: orgUUID, Created: nil, Updated: nil}, SenderName: "Sender2"},
		}

		mockUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
			return true
		}
		oHandler.userHasAccess = mockUserHasAccess

		mockDb.On("GetWorkspaceBudgetHistory", orgUUID).Return(expectedBudgetHistory).Once()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/budget/history/"+orgUUID, nil)
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

		assert.Equal(t, expectedBudgetHistory, responseBudgetHistory)
	})
}

func TestGetWorkspaceBountiesCount(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	oHandler := NewWorkspaceHandler(mockDb)

	t.Run("should return the count of organization bounties", func(t *testing.T) {
		orgUUID := "valid-uuid"
		expectedCount := int64(5)

		mockDb.On("GetWorkspaceBountiesCount", mock.AnythingOfType("*http.Request"), orgUUID).Return(expectedCount).Once()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/bounties/"+orgUUID+"/count/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetWorkspaceBountiesCount).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var count int64
		err = json.Unmarshal(rr.Body.Bytes(), &count)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedCount, count)
	})
}
