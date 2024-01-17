package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
  "github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	mocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
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
			{"empty description", "", http.StatusBadRequest},
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

	t.Run("should set Website, Github, and Description to null", func(t *testing.T) {
		orgUUID := "org-uuid"

		updatedOrg := db.Organization{
			Uuid:        orgUUID,
			OwnerPubKey: "test-key",
			Website:     nil,
			Github:      nil,
			Description: nil,
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

		assert.Nil(t, returnedOrg.Website)
		assert.Nil(t, returnedOrg.Github)
		assert.Nil(t, returnedOrg.Description)
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
