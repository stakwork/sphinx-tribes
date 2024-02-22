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

func TestUnitCreateOrEditOrganization(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	oHandler := NewOrganizationHandler(mockDb)

	t.Run("should return error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditOrganization)

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
		handler := http.HandlerFunc(oHandler.CreateOrEditOrganization)

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
		handler := http.HandlerFunc(oHandler.CreateOrEditOrganization)

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
		handler := http.HandlerFunc(oHandler.CreateOrEditOrganization)

		invalidJson := []byte(`{"name": "DemoTestingOrganization"}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return error if org name contains only spaces", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditOrganization)

		invalidJson := []byte(`{"name": "   "}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should successfully add organization if request is valid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditOrganization)

		mockDb.On("GetOrganizationByUuid", mock.AnythingOfType("string")).Return(db.Organization{}).Once()
		mockDb.On("GetOrganizationByName", "TestOrganization").Return(db.Organization{}).Once()
		mockDb.On("CreateOrEditOrganization", mock.MatchedBy(func(org db.Organization) bool {
			return org.Name == "TestOrganization" && org.Uuid != "" && org.Updated != nil && org.Created != nil
		})).Return(db.Organization{}, nil).Once()

		invalidJson := []byte(`{"name": "TestOrganization", "owner_pubkey": "test-key" ,"description": "Test"}`)
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
				handler := http.HandlerFunc(oHandler.CreateOrEditOrganization)
				invalidJson := []byte(fmt.Sprintf(`{"name": "TestOrganization", "owner_pubkey": "test-key", "description": "%s"}`, tc.description))

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

func TestDeleteOrganization(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	oHandler := NewOrganizationHandler(mockDb)

	t.Run("should return error if not authorized", func(t *testing.T) {
		orgUUID := "org-uuid"
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteOrganization)

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
		mockDb.On("GetOrganizationByUuid", orgUUID).Return(db.Organization{OwnerPubKey: "test-key"}).Once()
		mockDb.On("UpdateOrganizationForDeletion", orgUUID).Return(nil).Once()
		mockDb.On("DeleteAllUsersFromOrganization", orgUUID).Return(nil).Once()
		mockDb.On("ChangeOrganizationDeleteStatus", orgUUID, true).Return(db.Organization{Uuid: orgUUID, Deleted: true}).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteOrganization)

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
		mockDb.On("GetOrganizationByUuid", orgUUID).Return(db.Organization{OwnerPubKey: "test-key"}).Once()
		mockDb.On("UpdateOrganizationForDeletion", orgUUID).Return(errors.New("update error")).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteOrganization)

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
		mockDb.On("GetOrganizationByUuid", orgUUID).Return(db.Organization{OwnerPubKey: "test-key"}).Once()
		mockDb.On("UpdateOrganizationForDeletion", orgUUID).Return(nil).Once()
		mockDb.On("DeleteAllUsersFromOrganization", orgUUID).Return(nil).Once()
		mockDb.On("ChangeOrganizationDeleteStatus", orgUUID, true).Return(db.Organization{Uuid: orgUUID, Deleted: true}).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteOrganization)

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
		var updatedOrg db.Organization
		err = json.Unmarshal(rr.Body.Bytes(), &updatedOrg)
		if err != nil {
			t.Fatal(err)
		}
		assert.True(t, updatedOrg.Deleted)

		mockDb.AssertExpectations(t)
	})

	t.Run("should set Website, Github, and Description to empty strings", func(t *testing.T) {
		orgUUID := "org-uuid"

		updatedOrg := db.Organization{
			Uuid:        orgUUID,
			OwnerPubKey: "test-key",
			Website:     "",
			Github:      "",
			Description: "",
		}

		mockDb.On("GetOrganizationByUuid", orgUUID).Return(db.Organization{OwnerPubKey: "test-key"}).Once()
		mockDb.On("UpdateOrganizationForDeletion", orgUUID).Return(nil).Once()
		mockDb.On("DeleteAllUsersFromOrganization", orgUUID).Return(nil).Once()
		mockDb.On("ChangeOrganizationDeleteStatus", orgUUID, true).Return(updatedOrg).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteOrganization)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/delete/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var returnedOrg db.Organization
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
		mockDb.On("GetOrganizationByUuid", orgUUID).Return(db.Organization{OwnerPubKey: "test-key"}).Once()
		mockDb.On("UpdateOrganizationForDeletion", orgUUID).Return(nil).Once()
		mockDb.On("DeleteAllUsersFromOrganization", orgUUID).Return(nil).Run(func(args mock.Arguments) {}).Once()
		mockDb.On("ChangeOrganizationDeleteStatus", orgUUID, true).Return(db.Organization{Uuid: orgUUID, Deleted: true}).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteOrganization)

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

func TestGetOrganizationBounties(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	mockGenerateBountyHandler := func(bounties []db.Bounty) []db.BountyResponse {
		return []db.BountyResponse{} // Mocked response
	}
	oHandler := NewOrganizationHandler(mockDb)

	t.Run("Should test that an organization's bounties can be listed without authentication", func(t *testing.T) {
		orgUUID := "valid-uuid"
		oHandler.generateBountyHandler = mockGenerateBountyHandler

		expectedBounties := []db.Bounty{{}, {}} // Mocked response
		mockDb.On("GetOrganizationBounties", mock.AnythingOfType("*http.Request"), orgUUID).Return(expectedBounties).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetOrganizationBounties)

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

		mockDb.On("GetOrganizationBounties", mock.AnythingOfType("*http.Request"), orgUUID).Return([]db.Bounty{}).Once()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetOrganizationBounties)

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

func TestGetOrganizationBudget(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	oHandler := NewOrganizationHandler(mockDb)

	t.Run("Should test that a 401 is returned when trying to view an organization's budget without a token", func(t *testing.T) {
		orgUUID := "valid-uuid"

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/budget/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetOrganizationBudget).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that the right organization budget is returned, if the user is the organization admin or has the ViewReport role", func(t *testing.T) {
		orgUUID := "valid-uuid"
		expectedBudget := db.BountyBudget{
			ID:          1,
			OrgUuid:     orgUUID,
			TotalBudget: 1000,
			Created:     nil,
			Updated:     nil,
		}

		mockDb.On("UserHasAccess", "test-key", orgUUID, "VIEW REPORT").Return(true).Once()
		mockDb.On("GetOrganizationBudget", orgUUID).Return(expectedBudget).Once()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/budget/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetOrganizationBudget).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseBudget db.BountyBudget
		err = json.Unmarshal(rr.Body.Bytes(), &responseBudget)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedBudget, responseBudget)
	})
}

func TestGetOrganizationBudgetHistory(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	oHandler := NewOrganizationHandler(mockDb)

	t.Run("Should test that a 401 is returned when trying to view an organization's budget history without a token", func(t *testing.T) {
		orgUUID := "valid-uuid"

		mockDb.On("UserHasAccess", "", orgUUID, "VIEW REPORT").Return(false).Once()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/budget/history/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetOrganizationBudgetHistory).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that the right budget history is returned, if the user is the organization admin or has the ViewReport role", func(t *testing.T) {
		orgUUID := "valid-uuid"
		expectedBudgetHistory := []db.BudgetHistoryData{
			{BudgetHistory: db.BudgetHistory{ID: 1, OrgUuid: orgUUID, Created: nil, Updated: nil}, SenderName: "Sender1"},
			{BudgetHistory: db.BudgetHistory{ID: 2, OrgUuid: orgUUID, Created: nil, Updated: nil}, SenderName: "Sender2"},
		}

		mockDb.On("UserHasAccess", "test-key", orgUUID, "VIEW REPORT").Return(true).Once()
		mockDb.On("GetOrganizationBudgetHistory", orgUUID).Return(expectedBudgetHistory).Once()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/budget/history/"+orgUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetOrganizationBudgetHistory).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseBudgetHistory []db.BudgetHistoryData
		err = json.Unmarshal(rr.Body.Bytes(), &responseBudgetHistory)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedBudgetHistory, responseBudgetHistory)
	})
}

func TestGetOrganizationBountiesCount(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	oHandler := NewOrganizationHandler(mockDb)

	t.Run("should return the count of organization bounties", func(t *testing.T) {
		orgUUID := "valid-uuid"
		expectedCount := int64(5)

		mockDb.On("GetOrganizationBountiesCount", mock.AnythingOfType("*http.Request"), orgUUID).Return(expectedCount).Once()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", orgUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/bounties/"+orgUUID+"/count/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetOrganizationBountiesCount).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var count int64
		err = json.Unmarshal(rr.Body.Bytes(), &count)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedCount, count)
	})
}
