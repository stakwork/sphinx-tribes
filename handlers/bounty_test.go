package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stakwork/sphinx-tribes/utils"

	"github.com/go-chi/chi"
	"github.com/lib/pq"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers/mocks"
	dbMocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrEditBounty(t *testing.T) {

	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := dbMocks.NewDatabase(t)
	mockClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockClient, mockDb)

	t.Run("should return error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)

		invalidJson := []byte(`{"key": "value"`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code, "invalid status received")
	})

	t.Run("missing required field, bounty type", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)

		invalidBody := []byte(`{"type": ""}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("missing required field, bounty title", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)

		invalidBody := []byte(`{"type": "bounty_type", "title": ""}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("missing required field, bounty description", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)

		invalidBody := []byte(`{"type": "bounty_type", "title": "first bounty", "description": ""}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("return error if trying to update other user bounty", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)

		existingBounty := db.Bounty{
			ID:          1,
			Type:        "coding",
			Title:       "first bounty",
			Description: "first bounty description",
			OrgUuid:     "org-1",
			Assignee:    "user1",
		}
		mockDb.On("UpdateBountyBoolColumn", mock.AnythingOfType("db.Bounty"), "show").Return(existingBounty)
		mockDb.On("GetBounty", uint(1)).Return(existingBounty).Once()

		body := []byte(`{"id": 1, "type": "bounty_type", "title": "first bounty", "description": "my first bounty", "tribe": "random-value", "assignee": "john-doe", "owner_id": "second-user"}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, strings.TrimRight(rr.Body.String(), "\n"), "Cannot edit another user's bounty")
		mockDb.AssertExpectations(t)
	})

	t.Run("return error if user does not have required roles", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)

		mockOrg := db.Organization{
			ID:          1,
			Uuid:        "org-1",
			Name:        "custom org",
			OwnerPubKey: "org-key",
		}
		existingBounty := db.Bounty{
			ID:          1,
			Type:        "coding",
			Title:       "first bounty",
			Description: "first bounty description",
			OrgUuid:     "org-1",
			OwnerID:     "second-user",
		}
		updatedBounty := existingBounty
		updatedBounty.Title = "first bounty updated"
		mockDb.On("UpdateBountyBoolColumn", mock.AnythingOfType("db.Bounty"), "show").Return(existingBounty)
		mockDb.On("UpdateBountyNullColumn", mock.AnythingOfType("db.Bounty"), "assignee").Return(existingBounty)
		mockDb.On("GetBounty", uint(1)).Return(existingBounty).Once()
		mockDb.On("UserHasManageBountyRoles", "test-key", mockOrg.Uuid).Return(false).Once()

		body, _ := json.Marshal(updatedBounty)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should allow to add or edit bounty if user has role", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)

		mockOrg := db.Organization{
			ID:          1,
			Uuid:        "org-1",
			Name:        "custom org",
			OwnerPubKey: "org-key",
		}
		existingBounty := db.Bounty{
			ID:          1,
			Type:        "coding",
			Title:       "first bounty",
			Description: "first bounty description",
			OrgUuid:     "org-1",
			OwnerID:     "second-user",
		}
		updatedBounty := existingBounty
		updatedBounty.Title = "first bounty updated"
		mockDb.On("UpdateBountyBoolColumn", mock.AnythingOfType("db.Bounty"), "show").Return(existingBounty)
		mockDb.On("UpdateBountyNullColumn", mock.AnythingOfType("db.Bounty"), "assignee").Return(existingBounty)
		mockDb.On("GetBounty", uint(1)).Return(existingBounty).Once()
		mockDb.On("UserHasManageBountyRoles", "test-key", mockOrg.Uuid).Return(true).Once()
		mockDb.On("CreateOrEditBounty", mock.AnythingOfType("db.Bounty")).Return(updatedBounty, nil).Once()

		body, _ := json.Marshal(updatedBounty)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("should not update created at when bounty is updated", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)
		now := time.Now().UnixMilli()
		mockOrg := db.Organization{
			ID:          1,
			Uuid:        "org-1",
			Name:        "custom org",
			OwnerPubKey: "org-key",
		}
		existingBounty := db.Bounty{
			ID:          1,
			Type:        "coding",
			Title:       "first bounty",
			Description: "first bounty description",
			OrgUuid:     "org-1",
			OwnerID:     "second-user",
			Created:     now,
		}
		updatedBounty := existingBounty
		updatedBounty.Title = "first bounty updated"
		mockDb.On("UpdateBountyBoolColumn", mock.AnythingOfType("db.Bounty"), "show").Return(existingBounty)
		mockDb.On("UpdateBountyNullColumn", mock.AnythingOfType("db.Bounty"), "assignee").Return(existingBounty)
		mockDb.On("GetBounty", uint(1)).Return(existingBounty).Once()
		mockDb.On("UserHasManageBountyRoles", "test-key", mockOrg.Uuid).Return(true).Once()
		mockDb.On("CreateOrEditBounty", mock.MatchedBy(func(b db.Bounty) bool {
			return b.Created == now
		})).Return(updatedBounty, nil).Once()

		body, _ := json.Marshal(updatedBounty)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("should return error if failed to add new bounty", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)
		newBounty := db.Bounty{
			Type:        "coding",
			Title:       "first bounty",
			Description: "first bounty description",
			OrgUuid:     "org-1",
			OwnerID:     "test-key",
		}
		mockDb.On("UpdateBountyNullColumn", mock.AnythingOfType("db.Bounty"), "assignee").Return(db.Bounty{Assignee: "test-key"})
		mockDb.On("CreateOrEditBounty", mock.AnythingOfType("db.Bounty")).Return(db.Bounty{}, errors.New("failed to add")).Once()

		body, _ := json.Marshal(newBounty)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("add bounty if not present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)
		newBounty := db.Bounty{
			Type:        "coding",
			Title:       "first bounty",
			Description: "first bounty description",
			OrgUuid:     "org-1",
			OwnerID:     "test-key",
		}
		mockDb.On("UpdateBountyNullColumn", mock.AnythingOfType("db.Bounty"), "assignee").Return(db.Bounty{Assignee: "test-key"})
		mockDb.On("CreateOrEditBounty", mock.AnythingOfType("db.Bounty")).Return(newBounty, nil).Once()

		body, _ := json.Marshal(newBounty)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
	})
}

func TestPayLightningInvoice(t *testing.T) {
	expectedUrl := fmt.Sprintf("%s/invoices", config.RelayUrl)
	expectedBody := `{"payment_request": "req-id"}`

	t.Run("validate request url, body and headers", func(t *testing.T) {
		mockHttpClient := &mocks.HttpClient{}
		mockDb := &dbMocks.Database{}
		handler := NewBountyHandler(mockHttpClient, mockDb)
		mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			bodyByt, _ := io.ReadAll(req.Body)
			return req.Method == http.MethodPut && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
		})).Return(nil, errors.New("some-error")).Once()

		success, invoicePayErr := handler.PayLightningInvoice("req-id")

		assert.Empty(t, invoicePayErr)
		assert.Empty(t, success)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("put on invoice request failed with error status and invalid json", func(t *testing.T) {
		mockHttpClient := &mocks.HttpClient{}
		mockDb := &dbMocks.Database{}
		handler := NewBountyHandler(mockHttpClient, mockDb)
		r := io.NopCloser(bytes.NewReader([]byte(`"internal server error"`)))
		mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			bodyByt, _ := io.ReadAll(req.Body)
			return req.Method == http.MethodPut && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
		})).Return(&http.Response{
			StatusCode: 500,
			Body:       r,
		}, nil)

		success, invoicePayErr := handler.PayLightningInvoice("req-id")

		assert.False(t, invoicePayErr.Success)
		assert.Empty(t, success)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("put on invoice request failed with error status", func(t *testing.T) {
		mockHttpClient := &mocks.HttpClient{}
		mockDb := &dbMocks.Database{}
		handler := NewBountyHandler(mockHttpClient, mockDb)
		r := io.NopCloser(bytes.NewReader([]byte(`{"error": "internal server error"}`)))
		mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			bodyByt, _ := io.ReadAll(req.Body)
			return req.Method == http.MethodPut && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
		})).Return(&http.Response{
			StatusCode: 500,
			Body:       r,
		}, nil).Once()

		success, invoicePayErr := handler.PayLightningInvoice("req-id")

		assert.Equal(t, invoicePayErr.Error, "internal server error")
		assert.Empty(t, success)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("put on invoice request succeed with invalid json", func(t *testing.T) {
		mockHttpClient := &mocks.HttpClient{}
		mockDb := &dbMocks.Database{}
		handler := NewBountyHandler(mockHttpClient, mockDb)
		r := io.NopCloser(bytes.NewReader([]byte(`"invalid json"`)))
		mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			bodyByt, _ := io.ReadAll(req.Body)
			return req.Method == http.MethodPut && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
		})).Return(&http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil).Once()

		success, invoicePayErr := handler.PayLightningInvoice("req-id")

		assert.False(t, success.Success)
		assert.Empty(t, invoicePayErr)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("should unmarshal the response properly after success", func(t *testing.T) {
		mockHttpClient := &mocks.HttpClient{}
		mockDb := &dbMocks.Database{}
		handler := NewBountyHandler(mockHttpClient, mockDb)
		r := io.NopCloser(bytes.NewReader([]byte(`{"success": true, "response": { "settled": true, "payment_request": "req", "payment_hash": "hash", "preimage": "random-string", "amount": "1000"}}`)))
		expectedSuccessMsg := db.InvoicePaySuccess{
			Success: true,
			Response: db.InvoiceCheckResponse{
				Settled:         true,
				Payment_request: "req",
				Payment_hash:    "hash",
				Preimage:        "random-string",
				Amount:          "1000",
			},
		}
		mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			bodyByt, _ := io.ReadAll(req.Body)
			return req.Method == http.MethodPut && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
		})).Return(&http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil).Once()

		success, invoicePayErr := handler.PayLightningInvoice("req-id")

		assert.Empty(t, invoicePayErr)
		assert.EqualValues(t, expectedSuccessMsg, success)
		mockHttpClient.AssertExpectations(t)
	})

}

func TestDeleteBounty(t *testing.T) {
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")

	t.Run("should return unauthorized error if users public key not present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.DeleteBounty)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return unauthorized error if public key not present in route", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.DeleteBounty)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("pubkey", "")
		rctx.URLParams.Add("created", "1111")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "//1111", nil)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return unauthorized error if created at key not present in route", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.DeleteBounty)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("pubkey", "pub-key")
		rctx.URLParams.Add("created", "")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/pub-key/", nil)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return error if failed to delete from db", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.DeleteBounty)
		mockDb.On("DeleteBounty", "pub-key", "1111").Return(db.Bounty{}, errors.New("some-error")).Once()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("pubkey", "pub-key")
		rctx.URLParams.Add("created", "1111")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/pub-key/createdAt", nil)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("should successfully delete bounty from db", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.DeleteBounty)
		existingBounty := db.Bounty{
			OwnerID: "pub-key",
			Created: 1111,
		}
		mockDb.On("DeleteBounty", "pub-key", "1111").Return(existingBounty, nil).Once()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("pubkey", "pub-key")
		rctx.URLParams.Add("created", "1111")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/pub-key/1111", nil)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)

		var returnedBounty db.Bounty
		_ = json.Unmarshal(rr.Body.Bytes(), &returnedBounty)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, existingBounty, returnedBounty)
		mockDb.AssertExpectations(t)
	})
}

func TestGetBountyByCreated(t *testing.T) {
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)

	t.Run("Should return bounty by its created value", func(t *testing.T) {
		mockGenerateBountyResponse := func(bounties []db.Bounty) []db.BountyResponse {
			var bountyResponses []db.BountyResponse

			for _, bounty := range bounties {
				owner := db.Person{
					ID: 1,
				}
				assignee := db.Person{
					ID: 1,
				}
				organization := db.OrganizationShort{
					Uuid: "uuid",
				}

				bountyResponse := db.BountyResponse{
					Bounty:       bounty,
					Assignee:     assignee,
					Owner:        owner,
					Organization: organization,
				}
				bountyResponses = append(bountyResponses, bountyResponse)
			}

			return bountyResponses
		}
		bHandler.generateBountyResponse = mockGenerateBountyResponse

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyByCreated)
		bounty := db.Bounty{
			ID:          1,
			Type:        "coding",
			Title:       "first bounty",
			Description: "first bounty description",
			OrgUuid:     "org-1",
			Assignee:    "user1",
			Created:     1707991475,
			OwnerID:     "owner-1",
		}
		createdStr := strconv.FormatInt(bounty.Created, 10)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("created", "1707991475")
		req, _ := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/created/1707991475", nil)
		mockDb.On("GetBountyDataByCreated", createdStr).Return([]db.Bounty{bounty}, nil).Once()
		mockDb.On("GetPersonByPubkey", "owner-1").Return(db.Person{}).Once()
		mockDb.On("GetPersonByPubkey", "user1").Return(db.Person{}).Once()
		mockDb.On("GetOrganizationByUuid", "org-1").Return(db.Organization{}).Once()
		handler.ServeHTTP(rr, req)

		var returnedBounty []db.BountyResponse
		err := json.Unmarshal(rr.Body.Bytes(), &returnedBounty)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotEmpty(t, returnedBounty)

	})
	t.Run("Should return 404 if bounty is not present in db", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyByCreated)
		createdStr := ""

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("created", createdStr)
		req, _ := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/created/"+createdStr, nil)

		mockDb.On("GetBountyDataByCreated", createdStr).Return([]db.Bounty{}, nil).Once()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code, "Expected 404 Not Found for nonexistent bounty")

		mockDb.AssertExpectations(t)
	})

}

func TestGetPersonAssignedBounties(t *testing.T) {
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)
	t.Run("Should successfull Get Person Assigned Bounties", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetPersonAssignedBounties)
		bounty := db.Bounty{
			ID:          1,
			Type:        "coding",
			Title:       "first bounty",
			Description: "first bounty description",
			OrgUuid:     "org-1",
			Assignee:    "user1",
			Created:     1707991475,
			OwnerID:     "owner-1",
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", "clu80datu2rjujsmim40")
		rctx.URLParams.Add("sortBy", "paid")
		rctx.URLParams.Add("page", "1")
		rctx.URLParams.Add("limit", "20")
		rctx.URLParams.Add("search", "")
		req, _ := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/people/wanteds/assigned/clu80datu2rjujsmim40?sortBy=paid&page=1&limit=20&search=", nil)

		mockDb.On("GetAssignedBounties", req).Return([]db.Bounty{bounty}, nil).Once()
		mockDb.On("GetPersonByPubkey", "owner-1").Return(db.Person{}, nil).Once()
		mockDb.On("GetPersonByPubkey", "user1").Return(db.Person{}, nil).Once()
		mockDb.On("GetOrganizationByUuid", "org-1").Return(db.Organization{}, nil).Once()
		handler.ServeHTTP(rr, req)

		var returnedBounty []db.BountyResponse
		err := json.Unmarshal(rr.Body.Bytes(), &returnedBounty)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotEmpty(t, returnedBounty)
	})
}

func TestGetPersonCreatedBounties(t *testing.T) {
	ctx := context.Background()
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)

	t.Run("should return bounties created by the user", func(t *testing.T) {
		mockGenerateBountyResponse := func(bounties []db.Bounty) []db.BountyResponse {
			var bountyResponses []db.BountyResponse

			for _, bounty := range bounties {
				owner := db.Person{
					ID: 1,
				}
				assignee := db.Person{
					ID: 1,
				}
				organization := db.OrganizationShort{
					Uuid: "uuid",
				}

				bountyResponse := db.BountyResponse{
					Bounty:       bounty,
					Assignee:     assignee,
					Owner:        owner,
					Organization: organization,
				}
				bountyResponses = append(bountyResponses, bountyResponse)
			}

			return bountyResponses
		}
		bHandler.generateBountyResponse = mockGenerateBountyResponse

		expectedBounties := []db.Bounty{
			{ID: 1, OwnerID: "user1"},
			{ID: 2, OwnerID: "user1"},
		}

		mockDb.On("GetCreatedBounties", mock.Anything).Return(expectedBounties, nil).Once()
		mockDb.On("GetPersonByPubkey", mock.Anything).Return(db.Person{}, nil)
		mockDb.On("GetOrganizationByUuid", mock.Anything).Return(db.Organization{}, nil)
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/people/wanteds/created/uuid", nil)
		req = req.WithContext(ctx)
		if err != nil {
			t.Fatal(err)
		}

		bHandler.GetPersonCreatedBounties(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseData []db.BountyResponse
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}

		assert.NotEmpty(t, responseData)
		assert.Len(t, responseData, 2)

		for i, expectedBounty := range expectedBounties {
			assert.Equal(t, expectedBounty.ID, responseData[i].Bounty.ID)
			assert.Equal(t, expectedBounty.Assignee, responseData[i].Bounty.Assignee)
		}
	})

	t.Run("should not return bounties created by other users", func(t *testing.T) {
		mockGenerateBountyResponse := func(bounties []db.Bounty) []db.BountyResponse {
			return []db.BountyResponse{}
		}
		bHandler.generateBountyResponse = mockGenerateBountyResponse

		mockDb.On("GetCreatedBounties", mock.Anything).Return([]db.Bounty{}, nil).Once()

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/people/wanteds/created/uuid", nil)
		req = req.WithContext(ctx)
		if err != nil {
			t.Fatal(err)
		}

		bHandler.GetPersonCreatedBounties(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseData []db.BountyResponse
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}

		assert.Empty(t, responseData)
		assert.Len(t, responseData, 0)
	})

	t.Run("should filter bounties by status and apply pagination", func(t *testing.T) {
		mockGenerateBountyResponse := func(bounties []db.Bounty) []db.BountyResponse {
			var bountyResponses []db.BountyResponse

			for _, bounty := range bounties {
				owner := db.Person{
					ID: 1,
				}
				assignee := db.Person{
					ID: 1,
				}
				organization := db.OrganizationShort{
					Uuid: "uuid",
				}

				bountyResponse := db.BountyResponse{
					Bounty:       bounty,
					Assignee:     assignee,
					Owner:        owner,
					Organization: organization,
				}
				bountyResponses = append(bountyResponses, bountyResponse)
			}

			return bountyResponses
		}
		bHandler.generateBountyResponse = mockGenerateBountyResponse

		expectedBounties := []db.Bounty{
			{ID: 1, OwnerID: "user1", Assignee: "assignee1"},
			{ID: 2, OwnerID: "user1", Assignee: "assignee2", Paid: true},
			{ID: 3, OwnerID: "user1", Assignee: "", Paid: true},
		}

		mockDb.On("GetCreatedBounties", mock.Anything).Return(expectedBounties, nil).Once()
		mockDb.On("GetPersonByPubkey", mock.Anything).Return(db.Person{}, nil)
		mockDb.On("GetOrganizationByUuid", mock.Anything).Return(db.Organization{}, nil)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/people/wanteds/created/uuid?Open=true&Assigned=true&Paid=true&offset=0&limit=2", nil)
		req = req.WithContext(ctx)
		if err != nil {
			t.Fatal(err)
		}

		bHandler.GetPersonCreatedBounties(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseData []db.BountyResponse
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}

		assert.Len(t, responseData, 3)

		// Assert that bounties are filtered correctly
		assert.Equal(t, expectedBounties[0].ID, responseData[0].Bounty.ID)
		assert.Equal(t, expectedBounties[1].ID, responseData[1].Bounty.ID)
		assert.Equal(t, expectedBounties[2].ID, responseData[2].Bounty.ID)
	})
}

func TestGetNextBountyByCreated(t *testing.T) {
	ctx := context.Background()

	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)

	t.Run("Should test that the next bounty on the bounties homepage can be gotten by its created value and the selected filters", func(t *testing.T) {
		mockDb.On("GetNextBountyByCreated", mock.Anything).Return(uint(1), nil).Once()

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/next/123456789", nil)

		bHandler.GetNextBountyByCreated(rr, req.WithContext(ctx))

		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
	})
}

func TestGetPreviousBountyByCreated(t *testing.T) {
	ctx := context.Background()

	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)

	t.Run("Should test that the previous bounty on the bounties homepage can be gotten by its created value and the selected filters", func(t *testing.T) {
		mockDb.On("GetPreviousBountyByCreated", mock.Anything).Return(uint(1), nil).Once()

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/previous/123456789", nil)

		bHandler.GetPreviousBountyByCreated(rr, req.WithContext(ctx))

		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
	})
}

func TestGetOrganizationNextBountyByCreated(t *testing.T) {
	ctx := context.Background()
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)

	t.Run("Should test that the next bounty on the organization bounties homepage can be gotten by its created value and the selected filters", func(t *testing.T) {
		mockDb.On("GetNextOrganizationBountyByCreated", mock.AnythingOfType("*http.Request")).Return(uint(1), nil).Once()

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/org/next/org-uuid/123456789", nil)

		bHandler.GetOrganizationNextBountyByCreated(rr, req.WithContext(ctx))

		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
	})
}

func TestGetOrganizationPreviousBountyByCreated(t *testing.T) {
	ctx := context.Background()
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)

	t.Run("Should test that the previous bounty on the organization bounties homepage can be gotten by its created value and the selected filters", func(t *testing.T) {
		mockDb.On("GetPreviousOrganizationBountyByCreated", mock.AnythingOfType("*http.Request")).Return(uint(1), nil).Once()

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/org/previous/org-uuid/123456789", nil)

		bHandler.GetOrganizationPreviousBountyByCreated(rr, req.WithContext(ctx))

		assert.Equal(t, http.StatusOK, rr.Code)

		mockDb.AssertExpectations(t)
	})
}

func TestGetBountyById(t *testing.T) {

	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)

	t.Run("successful retrieval of bounty by ID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyById)

		bounty := db.Bounty{
			ID:                      1,
			OwnerID:                 "owner123",
			Paid:                    false,
			Show:                    true,
			Type:                    "bug fix",
			Award:                   "500",
			AssignedHours:           10,
			BountyExpires:           "2023-12-31",
			CommitmentFee:           1000,
			Price:                   500,
			Title:                   "Fix critical bug in payment system",
			Tribe:                   "development",
			Assignee:                "user1",
			TicketUrl:               "http://example.com/issues/1",
			OrgUuid:                 "org-789",
			Description:             "This bounty is for fixing a critical bug in the payment system that causes transactions to fail under certain conditions.",
			WantedType:              "immediate",
			Deliverables:            "A pull request with a fix, including tests",
			GithubDescription:       true,
			OneSentenceSummary:      "Fix a critical payment system bug",
			EstimatedSessionLength:  "2 hours",
			EstimatedCompletionDate: "2023-10-01",
			Created:                 time.Now().Unix(),
			Updated:                 nil,
			AssignedDate:            nil,
			CompletionDate:          nil,
			MarkAsPaidDate:          nil,
			PaidDate:                nil,
			CodingLanguages:         pq.StringArray{"Go", "Python"},
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("bountyId", strconv.Itoa(int(bounty.ID)))
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/bounty/1", nil)
		assert.NoError(t, err)

		mockDb.On("GetBountyById", mock.Anything).Return([]db.Bounty{bounty}, nil).Once()
		mockDb.On("GetPersonByPubkey", "owner123").Return(db.Person{}).Once()
		mockDb.On("GetPersonByPubkey", "user1").Return(db.Person{}).Once()
		mockDb.On("GetOrganizationByUuid", "org-789").Return(db.Organization{}).Once()

		handler.ServeHTTP(rr, req)

		var returnedBounty []db.BountyResponse
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBounty)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("bounty not found", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyById)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("bountyId", "999")
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/bounty/999", nil)
		assert.NoError(t, err)

		mockDb.On("GetBountyById", "999").Return(nil, errors.New("not-found")).Once()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockDb.AssertExpectations(t)
	})
}

func GetPersonAssigned(t *testing.T) {
	ctx := context.Background()
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)
	t.Run("should return bounties assigned to the user", func(t *testing.T) {
		mockGenerateBountyResponse := func(bounties []db.Bounty) []db.BountyResponse {
			var bountyResponses []db.BountyResponse

			for _, bounty := range bounties {
				owner := db.Person{
					ID: 1,
				}
				assignee := db.Person{
					ID: 1,
				}
				organization := db.OrganizationShort{
					Uuid: "uuid",
				}

				bountyResponse := db.BountyResponse{
					Bounty:       bounty,
					Assignee:     assignee,
					Owner:        owner,
					Organization: organization,
				}
				bountyResponses = append(bountyResponses, bountyResponse)
			}

			return bountyResponses
		}
		bHandler.generateBountyResponse = mockGenerateBountyResponse

		expectedBounties := []db.Bounty{
			{ID: 1, Assignee: "user1"},
			{ID: 2, Assignee: "user1"},
			{ID: 3, OwnerID: "user2", Assignee: "user1"},
			{ID: 4, OwnerID: "user2", Assignee: "user1", Paid: true},
		}

		mockDb.On("GetAssignedBounties", mock.Anything).Return(expectedBounties, nil).Once()
		mockDb.On("GetPersonByPubkey", mock.Anything).Return(db.Person{}, nil)
		mockDb.On("GetOrganizationByUuid", mock.Anything).Return(db.Organization{}, nil)
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/wanteds/assigned/uuid?Assigned=true&Paid=true&offset=0&limit=4", nil)
		req = req.WithContext(ctx)
		if err != nil {
			t.Fatal(err)
		}

		bHandler.GetPersonAssignedBounties(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseData []db.BountyResponse
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}

		assert.NotEmpty(t, responseData)
		assert.Len(t, responseData, 2)

		for i, expectedBounty := range expectedBounties {
			assert.Equal(t, expectedBounty.ID, responseData[i].Bounty.ID)
			assert.Equal(t, expectedBounty.Assignee, responseData[i].Bounty.Assignee)
		}
	})

	t.Run("should not return bounties assigned to other users", func(t *testing.T) {
		mockGenerateBountyResponse := func(bounties []db.Bounty) []db.BountyResponse {
			return []db.BountyResponse{}
		}
		bHandler.generateBountyResponse = mockGenerateBountyResponse

		mockDb.On("GetAssignedBounties", mock.Anything).Return([]db.Bounty{}, nil).Once()

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/wanteds/assigned/uuid", nil)
		req = req.WithContext(ctx)
		if err != nil {
			t.Fatal(err)
		}

		bHandler.GetPersonAssignedBounties(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseData []db.BountyResponse
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}

		assert.Empty(t, responseData)
		assert.Len(t, responseData, 0)
	})
}
func TestGetBountyIndexById(t *testing.T) {
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)

	t.Run("successful retrieval of bounty by Index ID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyIndexById)

		bounty := db.Bounty{
			ID: 1,
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("bountyId", strconv.Itoa(int(bounty.ID)))
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/index/1", nil)
		assert.NoError(t, err)

		mockDb.On("GetBountyIndexById", "1").Return(int64(12), nil).Once()

		handler.ServeHTTP(rr, req)

		responseBody := rr.Body.Bytes()
		responseString := strings.TrimSpace(string(responseBody))
		returnedIndex, err := strconv.Atoi(responseString)
		assert.NoError(t, err)
		assert.Equal(t, 12, returnedIndex)

		assert.Equal(t, http.StatusOK, rr.Code)

		mockDb.AssertExpectations(t)
	})

	t.Run("bounty index by ID not found", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyIndexById)

		bountyID := ""
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("bountyId", bountyID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/index/"+bountyID, nil)
		assert.NoError(t, err)

		mockDb.On("GetBountyIndexById", bountyID).Return(int64(0), fmt.Errorf("bounty not found")).Once()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)

		mockDb.AssertExpectations(t)
	})
}

func TestGetAllBounties(t *testing.T) {
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)
	t.Run("Should successfull All Bounties", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetAllBounties)
		bounties := []db.Bounty{
			{ID: 1,
				Type:        "coding",
				Title:       "first bounty",
				Description: "first bounty description",
				OrgUuid:     "org-1",
				Assignee:    "user1",
				Created:     1707991475,
				OwnerID:     "owner-1",
			},
		}

		rctx := chi.NewRouteContext()
		req, _ := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/all", nil)

		mockDb.On("GetAllBounties", req).Return(bounties)
		mockDb.On("GetPersonByPubkey", mock.Anything).Return(db.Person{}, nil)
		mockDb.On("GetOrganizationByUuid", mock.Anything).Return(db.Organization{}, nil)
		handler.ServeHTTP(rr, req)

		var returnedBounty []db.BountyResponse
		err := json.Unmarshal(rr.Body.Bytes(), &returnedBounty)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotEmpty(t, returnedBounty)

	})
}

func MockNewWSServer(t *testing.T) (*httptest.Server, *websocket.Conn) {

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{}

		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("upgrade:", err)
			return
		}
		defer ws.Close()
	}))
	wsURL := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	return s, ws
}

func TestMakeBountyPayment(t *testing.T) {
	ctx := context.Background()
	mockDb := &dbMocks.Database{}
	mockHttpClient := &mocks.HttpClient{}
	mockGetSocketConnections := func(host string) (db.Client, error) {
		s, ws := MockNewWSServer(t)
		defer s.Close()
		defer ws.Close()

		mockClient := db.Client{
			Host: "mocked_host",
			Conn: ws,
		}

		return mockClient, nil
	}
	bHandler := NewBountyHandler(mockHttpClient, mockDb)

	unauthorizedCtx := context.WithValue(ctx, auth.ContextKey, "")
	authorizedCtx := context.WithValue(ctx, auth.ContextKey, "valid-key")

	var mutex sync.Mutex
	var processingTimes []time.Time

	bountyID := uint(1)
	bounty := db.Bounty{
		ID:       bountyID,
		OrgUuid:  "org-1",
		Assignee: "assignee-1",
		Price:    uint(1000),
	}

	t.Run("mutex lock ensures sequential access", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mutex.Lock()
			processingTimes = append(processingTimes, time.Now())
			time.Sleep(10 * time.Millisecond)
			mutex.Unlock()

			bHandler.MakeBountyPayment(w, r)
		}))
		defer server.Close()

		var wg sync.WaitGroup
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := http.Get(server.URL)
				if err != nil {
					t.Errorf("Failed to send request: %v", err)
				}
			}()
		}
		wg.Wait()

		for i := 1; i < len(processingTimes); i++ {
			assert.True(t, processingTimes[i].After(processingTimes[i-1]),
				"Expected processing times to be sequential, indicating mutex is locking effectively.")
		}
	})

	t.Run("401 unauthorized error when unauthorized user hits endpoint", func(t *testing.T) {

		r := chi.NewRouter()
		r.Post("/gobounties/pay/{id}", bHandler.MakeBountyPayment)

		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(unauthorizedCtx, http.MethodPost, "/gobounties/pay/1", nil)

		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected 401 Unauthorized for unauthorized access")
		mockDb.AssertExpectations(t)
	})

	t.Run("405 when trying to pay an already-paid bounty", func(t *testing.T) {
		mockDb.ExpectedCalls = nil
		mockDb.On("GetBounty", mock.AnythingOfType("uint")).Return(db.Bounty{
			ID:       1,
			Price:    1000,
			OrgUuid:  "org-1",
			Assignee: "assignee-1",
			Paid:     true,
		}, nil)

		r := chi.NewRouter()
		r.Post("/gobounties/pay/{id}", bHandler.MakeBountyPayment)

		requestBody := bytes.NewBuffer([]byte("{}"))
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/gobounties/pay/1", requestBody)
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Expected 405 Method Not Allowed for an already-paid bounty")
		mockDb.AssertExpectations(t)
	})

	t.Run("401 error if user not organization admin or does not have PAY BOUNTY role", func(t *testing.T) {
		mockDb.On("GetBounty", mock.AnythingOfType("uint")).Return(db.Bounty{
			ID:       1,
			Price:    1000,
			OrgUuid:  "org-1",
			Assignee: "assignee-1",
			Paid:     false,
		}, nil)
		mockDb.On("UserHasAccess", "valid-key", "org-1", db.PayBounty).Return(false)

		r := chi.NewRouter()
		r.Post("/gobounties/pay/{id}", bHandler.MakeBountyPayment)

		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(unauthorizedCtx, http.MethodPost, "/gobounties/pay/1", bytes.NewBufferString(`{}`))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected 401 Unauthorized when the user lacks the PAY BOUNTY role")

	})

	t.Run("403 error when amount exceeds organization's budget balance", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.ContextKey, "valid-key")

		mockDb := dbMocks.NewDatabase(t)
		mockHttpClient := mocks.NewHttpClient(t)
		bHandler := NewBountyHandler(mockHttpClient, mockDb)
		mockDb.On("GetBounty", mock.AnythingOfType("uint")).Return(db.Bounty{
			ID:       1,
			Price:    1000,
			OrgUuid:  "org-1",
			Assignee: "assignee-1",
			Paid:     false,
		}, nil)
		mockDb.On("UserHasAccess", "valid-key", "org-1", db.PayBounty).Return(true)
		mockDb.On("GetOrganizationBudget", "org-1").Return(db.BountyBudget{
			TotalBudget: 500,
		}, nil)

		r := chi.NewRouter()
		r.Post("/gobounties/pay/{id}", bHandler.MakeBountyPayment)

		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/gobounties/pay/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code, "Expected 403 Forbidden when the payment exceeds the organization's budget")

	})

	t.Run("Should test that a successful WebSocket message is sent if the payment is successful", func(t *testing.T) {
		mockDb.ExpectedCalls = nil
		bHandler.getSocketConnections = mockGetSocketConnections

		now := time.Now()
		expectedBounty := db.Bounty{
			ID:             bountyID,
			OrgUuid:        "org-1",
			Assignee:       "assignee-1",
			Price:          uint(1000),
			Paid:           true,
			PaidDate:       &now,
			CompletionDate: &now,
		}

		mockDb.On("GetBounty", bountyID).Return(bounty, nil)
		mockDb.On("UserHasAccess", "valid-key", bounty.OrgUuid, db.PayBounty).Return(true)
		mockDb.On("GetOrganizationBudget", bounty.OrgUuid).Return(db.BountyBudget{TotalBudget: 2000}, nil)
		mockDb.On("GetPersonByPubkey", bounty.Assignee).Return(db.Person{OwnerPubKey: "assignee-1", OwnerRouteHint: "OwnerRouteHint"}, nil)
		mockDb.On("AddPaymentHistory", mock.AnythingOfType("db.PaymentHistory")).Return(db.PaymentHistory{ID: 1})
		mockDb.On("UpdateBounty", mock.AnythingOfType("db.Bounty")).Run(func(args mock.Arguments) {
			updatedBounty := args.Get(0).(db.Bounty)
			assert.True(t, updatedBounty.Paid)
			assert.NotNil(t, updatedBounty.PaidDate)
			assert.NotNil(t, updatedBounty.CompletionDate)
		}).Return(expectedBounty, nil).Once()

		expectedUrl := fmt.Sprintf("%s/payment", config.RelayUrl)
		expectedBody := `{"amount": 1000, "destination_key": "assignee-1", "route_hint": "OwnerRouteHint", "text": "memotext added for notification"}`

		r := io.NopCloser(bytes.NewReader([]byte(`{"success": true, "response": { "sumAmount": "1"}}`)))
		mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			bodyByt, _ := io.ReadAll(req.Body)
			return req.Method == http.MethodPost && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
		})).Return(&http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil).Once()

		ro := chi.NewRouter()
		ro.Post("/gobounties/pay/{id}", bHandler.MakeBountyPayment)

		requestBody := bytes.NewBuffer([]byte("{}"))
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/gobounties/pay/1", requestBody)
		if err != nil {
			t.Fatal(err)
		}

		ro.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("Should test that an error WebSocket message is sent if the payment fails", func(t *testing.T) {
		mockDb2 := &dbMocks.Database{}
		mockHttpClient2 := &mocks.HttpClient{}
		mockDb2.ExpectedCalls = nil

		bHandler2 := NewBountyHandler(mockHttpClient2, mockDb2)
		bHandler2.getSocketConnections = mockGetSocketConnections

		mockDb2.On("GetBounty", bountyID).Return(bounty, nil)
		mockDb2.On("UserHasAccess", "valid-key", bounty.OrgUuid, db.PayBounty).Return(true)
		mockDb2.On("GetOrganizationBudget", bounty.OrgUuid).Return(db.BountyBudget{TotalBudget: 2000}, nil)
		mockDb2.On("GetPersonByPubkey", bounty.Assignee).Return(db.Person{OwnerPubKey: "assignee-1", OwnerRouteHint: "OwnerRouteHint"}, nil)

		expectedUrl := fmt.Sprintf("%s/payment", config.RelayUrl)
		expectedBody := `{"amount": 1000, "destination_key": "assignee-1", "route_hint": "OwnerRouteHint", "text": "memotext added for notification"}`

		r := io.NopCloser(bytes.NewReader([]byte(`"internal server error"`)))
		mockHttpClient2.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			bodyByt, _ := io.ReadAll(req.Body)
			return req.Method == http.MethodPost && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
		})).Return(&http.Response{
			StatusCode: 500,
			Body:       r,
		}, nil).Once()

		ro := chi.NewRouter()
		ro.Post("/gobounties/pay/{id}", bHandler2.MakeBountyPayment)

		requestBody := bytes.NewBuffer([]byte("{}"))
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/gobounties/pay/1", requestBody)
		if err != nil {
			t.Fatal(err)
		}

		ro.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb2.AssertExpectations(t)
		mockHttpClient2.AssertExpectations(t)
	})
}

func TestBountyBudgetWithdraw(t *testing.T) {
	ctx := context.Background()
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)
	unauthorizedCtx := context.WithValue(context.Background(), auth.ContextKey, "")
	authorizedCtx := context.WithValue(ctx, auth.ContextKey, "valid-key")

	t.Run("401 error if user is unauthorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.BountyBudgetWithdraw)

		req, err := http.NewRequestWithContext(unauthorizedCtx, http.MethodPost, "/budget/withdraw", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that a 406 error is returned if wrong data is passed", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.BountyBudgetWithdraw)

		invalidJson := []byte(`"key": "value"`)

		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/budget/withdraw", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("401 error if user is not the organization admin or does not have WithdrawBudget role", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.BountyBudgetWithdraw)
		mockDb.On("UserHasAccess", "valid-key", mock.AnythingOfType("string"), db.WithdrawBudget).Return(false)

		validData := []byte(`{"orgUuid": "org-1", "paymentRequest": "invoice"}`)
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/budget/withdraw", bytes.NewReader(validData))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "You don't have appropriate permissions to withdraw bounty budget")
	})

	t.Run("403 error when amount exceeds organization's budget", func(t *testing.T) {
		ctxs := context.WithValue(context.Background(), auth.ContextKey, "valid-key")
		mockDb := dbMocks.NewDatabase(t)
		mockHttpClient := mocks.NewHttpClient(t)
		bHandler := NewBountyHandler(mockHttpClient, mockDb)

		mockDb.On("UserHasAccess", "valid-key", "org-1", db.WithdrawBudget).Return(true)
		mockDb.On("GetOrganizationBudget", "org-1").Return(db.BountyBudget{
			TotalBudget: 500,
		}, nil)
		invoice := "lnbc15u1p3xnhl2pp5jptserfk3zk4qy42tlucycrfwxhydvlemu9pqr93tuzlv9cc7g3sdqsvfhkcap3xyhx7un8cqzpgxqzjcsp5f8c52y2stc300gl6s4xswtjpc37hrnnr3c9wvtgjfuvqmpm35evq9qyyssqy4lgd8tj637qcjp05rdpxxykjenthxftej7a2zzmwrmrl70fyj9hvj0rewhzj7jfyuwkwcg9g2jpwtk3wkjtwnkdks84hsnu8xps5vsq4gj5hs"

		amount := utils.GetInvoiceAmount(invoice)
		assert.Equal(t, uint(1500), amount)

		withdrawRequest := db.WithdrawBudgetRequest{
			PaymentRequest: invoice,
			OrgUuid:        "org-1",
		}
		requestBody, _ := json.Marshal(withdrawRequest)
		req, _ := http.NewRequestWithContext(ctxs, http.MethodPost, "/budget/withdraw", bytes.NewReader(requestBody))

		rr := httptest.NewRecorder()

		bHandler.BountyBudgetWithdraw(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code, "Expected 403 Forbidden when the payment exceeds the organization's budget")
		assert.Contains(t, rr.Body.String(), "Organization budget is not enough to withdraw the amount", "Expected specific error message")
	})

	t.Run("budget invoices get paid if amount is lesser than organization's budget", func(t *testing.T) {
		ctxs := context.WithValue(context.Background(), auth.ContextKey, "valid-key")
		mockDb := dbMocks.NewDatabase(t)
		mockHttpClient := mocks.NewHttpClient(t)
		bHandler := NewBountyHandler(mockHttpClient, mockDb)

		paymentAmount := uint(1500)

		mockDb.On("UserHasAccess", "valid-key", "org-1", db.WithdrawBudget).Return(true)
		mockDb.On("GetOrganizationBudget", "org-1").Return(db.BountyBudget{
			TotalBudget: 5000,
		}, nil)
		mockDb.On("WithdrawBudget", "valid-key", "org-1", paymentAmount).Return(nil)
		mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"success": true}`)),
		}, nil)

		invoice := "lnbc15u1p3xnhl2pp5jptserfk3zk4qy42tlucycrfwxhydvlemu9pqr93tuzlv9cc7g3sdqsvfhkcap3xyhx7un8cqzpgxqzjcsp5f8c52y2stc300gl6s4xswtjpc37hrnnr3c9wvtgjfuvqmpm35evq9qyyssqy4lgd8tj637qcjp05rdpxxykjenthxftej7a2zzmwrmrl70fyj9hvj0rewhzj7jfyuwkwcg9g2jpwtk3wkjtwnkdks84hsnu8xps5vsq4gj5hs"

		withdrawRequest := db.WithdrawBudgetRequest{
			PaymentRequest: invoice,
			OrgUuid:        "org-1",
		}
		requestBody, _ := json.Marshal(withdrawRequest)
		req, _ := http.NewRequestWithContext(ctxs, http.MethodPost, "/budget/withdraw", bytes.NewReader(requestBody))

		rr := httptest.NewRecorder()

		bHandler.BountyBudgetWithdraw(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var response db.InvoicePaySuccess
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success, "Expected invoice payment to succeed")

		mockDb.AssertCalled(t, "WithdrawBudget", "valid-key", "org-1", paymentAmount)
	})

	t.Run("400 BadRequest error if there is an error with invoice payment", func(t *testing.T) {
		ctxs := context.WithValue(context.Background(), auth.ContextKey, "valid-key")
		mockDb := dbMocks.NewDatabase(t)
		mockHttpClient := mocks.NewHttpClient(t)
		bHandler := NewBountyHandler(mockHttpClient, mockDb)

		mockDb.On("UserHasAccess", "valid-key", "org-1", db.WithdrawBudget).Return(true)
		mockDb.On("GetOrganizationBudget", "org-1").Return(db.BountyBudget{
			TotalBudget: 5000,
		}, nil)
		mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewBufferString(`{"success": false, "error": "Payment error"}`)),
		}, nil)

		invoice := "lnbc15u1p3xnhl2pp5jptserfk3zk4qy42tlucycrfwxhydvlemu9pqr93tuzlv9cc7g3sdqsvfhkcap3xyhx7un8cqzpgxqzjcsp5f8c52y2stc300gl6s4xswtjpc37hrnnr3c9wvtgjfuvqmpm35evq9qyyssqy4lgd8tj637qcjp05rdpxxykjenthxftej7a2zzmwrmrl70fyj9hvj0rewhzj7jfyuwkwcg9g2jpwtk3wkjtwnkdks84hsnu8xps5vsq4gj5hs"

		withdrawRequest := db.WithdrawBudgetRequest{
			PaymentRequest: invoice,
			OrgUuid:        "org-1",
		}
		requestBody, _ := json.Marshal(withdrawRequest)
		req, _ := http.NewRequestWithContext(ctxs, http.MethodPost, "/budget/withdraw", bytes.NewReader(requestBody))

		rr := httptest.NewRecorder()

		bHandler.BountyBudgetWithdraw(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["success"].(bool))
		assert.Equal(t, "Payment error", response["error"].(string))
		mockHttpClient.AssertCalled(t, "Do", mock.AnythingOfType("*http.Request"))
	})

	t.Run("Should test that an Organization's Budget Total Amount is accurate after three (3) successful 'Budget Withdrawal Requests'", func(t *testing.T) {
		ctxs := context.WithValue(context.Background(), auth.ContextKey, "valid-key")
		mockDb := dbMocks.NewDatabase(t)
		mockHttpClient := mocks.NewHttpClient(t)
		bHandler := NewBountyHandler(mockHttpClient, mockDb)

		paymentAmount := uint(1500)
		initialBudget := uint(5000)
		invoice := "lnbc15u1p3xnhl2pp5jptserfk3zk4qy42tlucycrfwxhydvlemu9pqr93tuzlv9cc7g3sdqsvfhkcap3xyhx7un8cqzpgxqzjcsp5f8c52y2stc300gl6s4xswtjpc37hrnnr3c9wvtgjfuvqmpm35evq9qyyssqy4lgd8tj637qcjp05rdpxxykjenthxftej7a2zzmwrmrl70fyj9hvj0rewhzj7jfyuwkwcg9g2jpwtk3wkjtwnkdks84hsnu8xps5vsq4gj5hs"

		for i := 0; i < 3; i++ {
			expectedFinalBudget := initialBudget - (paymentAmount * uint(i))

			mockDb.ExpectedCalls = nil
			mockDb.Calls = nil
			mockHttpClient.ExpectedCalls = nil
			mockHttpClient.Calls = nil

			mockDb.On("UserHasAccess", "valid-key", "org-1", db.WithdrawBudget).Return(true)
			mockDb.On("GetOrganizationBudget", "org-1").Return(db.BountyBudget{
				TotalBudget: expectedFinalBudget,
			}, nil)
			mockDb.On("WithdrawBudget", "valid-key", "org-1", paymentAmount).Return(nil)
			mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`{"success": true}`)),
			}, nil)

			withdrawRequest := db.WithdrawBudgetRequest{
				PaymentRequest: invoice,
				OrgUuid:        "org-1",
			}
			requestBody, _ := json.Marshal(withdrawRequest)
			req, _ := http.NewRequestWithContext(ctxs, http.MethodPost, "/budget/withdraw", bytes.NewReader(requestBody))

			rr := httptest.NewRecorder()

			bHandler.BountyBudgetWithdraw(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code)
			var response db.InvoicePaySuccess
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.True(t, response.Success, "Expected invoice payment to succeed")
			finalBudget := mockDb.GetOrganizationBudget("org-1")
			assert.Equal(t, expectedFinalBudget, finalBudget.TotalBudget, "The organization's final budget should reflect the deductions from the successful withdrawals")

		}
	})
}

func TestPollInvoice(t *testing.T) {
	ctx := context.Background()
	mockDb := &dbMocks.Database{}
	mockHttpClient := &mocks.HttpClient{}
	bHandler := NewBountyHandler(mockHttpClient, mockDb)

	unauthorizedCtx := context.WithValue(ctx, auth.ContextKey, "")
	authorizedCtx := context.WithValue(ctx, auth.ContextKey, "valid-key")

	t.Run("Should test that a 401 error is returned if a user is unauthorized", func(t *testing.T) {
		r := chi.NewRouter()
		r.Post("/poll/invoice/{paymentRequest}", bHandler.PollInvoice)

		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(unauthorizedCtx, http.MethodPost, "/poll/invoice/1", bytes.NewBufferString(`{}`))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected 401 error if a user is unauthorized")
	})

	t.Run("Should test that a 403 error is returned if there is an invoice error", func(t *testing.T) {
		expectedUrl := fmt.Sprintf("%s/invoice?payment_request=%s", config.RelayUrl, "1")

		r := io.NopCloser(bytes.NewReader([]byte(`{"success": false, "error": "Internel server error"}`)))
		mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.Method == http.MethodGet && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey
		})).Return(&http.Response{
			StatusCode: 500,
			Body:       r,
		}, nil).Once()

		ro := chi.NewRouter()
		ro.Post("/poll/invoice/{paymentRequest}", bHandler.PollInvoice)

		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/poll/invoice/1", bytes.NewBufferString(`{}`))
		if err != nil {
			t.Fatal(err)
		}

		ro.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code, "Expected 403 error if there is an invoice error")
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("Should mock relay payment is successful update the bounty associated with the invoice and set the paid as true", func(t *testing.T) {
		expectedUrl := fmt.Sprintf("%s/invoice?payment_request=%s", config.RelayUrl, "1")

		r := io.NopCloser(bytes.NewReader([]byte(`{"success": true, "response": { "settled": true, "payment_request": "1", "payment_hash": "payment_hash", "preimage": "preimage", "Amount": "1000"}}`)))
		mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.Method == http.MethodGet && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey
		})).Return(&http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil).Once()

		bountyID := uint(1)
		bounty := db.Bounty{
			ID:       bountyID,
			OrgUuid:  "org-1",
			Assignee: "assignee-1",
			Price:    uint(1000),
		}

		now := time.Now()
		expectedBounty := db.Bounty{
			ID:             bountyID,
			OrgUuid:        "org-1",
			Assignee:       "assignee-1",
			Price:          uint(1000),
			Paid:           true,
			PaidDate:       &now,
			CompletionDate: &now,
		}

		mockDb.On("GetInvoice", "1").Return(db.InvoiceList{Type: "KEYSEND"})
		mockDb.On("GetUserInvoiceData", "1").Return(db.UserInvoiceData{Amount: 1000, UserPubkey: "UserPubkey", RouteHint: "RouteHint", Created: 1234})
		mockDb.On("GetInvoice", "1").Return(db.InvoiceList{Status: false})
		mockDb.On("GetBountyByCreated", uint(1234)).Return(bounty, nil)
		mockDb.On("UpdateBounty", mock.AnythingOfType("db.Bounty")).Run(func(args mock.Arguments) {
			updatedBounty := args.Get(0).(db.Bounty)
			assert.True(t, updatedBounty.Paid)
		}).Return(expectedBounty, nil).Once()
		mockDb.On("UpdateInvoice", "1").Return(db.InvoiceList{}).Once()

		expectedPaymentUrl := fmt.Sprintf("%s/payment", config.RelayUrl)
		expectedPaymentBody := `{"amount": 1000, "destination_key": "UserPubkey", "route_hint": "RouteHint", "text": "memotext added for notification"}`

		r2 := io.NopCloser(bytes.NewReader([]byte(`{"success": true, "response": { "sumAmount": "1"}}`)))
		mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			bodyByt, _ := io.ReadAll(req.Body)
			return req.Method == http.MethodPost && expectedPaymentUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedPaymentBody == string(bodyByt)
		})).Return(&http.Response{
			StatusCode: 200,
			Body:       r2,
		}, nil).Once()

		ro := chi.NewRouter()
		ro.Post("/poll/invoice/{paymentRequest}", bHandler.PollInvoice)

		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/poll/invoice/1", bytes.NewBufferString(`{}`))
		if err != nil {
			t.Fatal(err)
		}

		ro.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("If the invoice is settled and the invoice.Type is equal to BUDGET the invoice amount should be added to the organization budget and the payment status of the related invoice should be sent to true on the payment history table", func(t *testing.T) {
		ctx := context.Background()
		mockDb := &dbMocks.Database{}
		mockHttpClient := &mocks.HttpClient{}
		bHandler := NewBountyHandler(mockHttpClient, mockDb)
		authorizedCtx := context.WithValue(ctx, auth.ContextKey, "valid-key")
		expectedUrl := fmt.Sprintf("%s/invoice?payment_request=%s", config.RelayUrl, "1")

		r := io.NopCloser(bytes.NewReader([]byte(`{"success": true, "response": { "settled": true, "payment_request": "1", "payment_hash": "payment_hash", "preimage": "preimage", "Amount": "1000"}}`)))
		mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.Method == http.MethodGet && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey
		})).Return(&http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil).Once()

		mockDb.On("GetInvoice", "1").Return(db.InvoiceList{Type: "BUDGET"})
		mockDb.On("GetUserInvoiceData", "1").Return(db.UserInvoiceData{Amount: 1000, UserPubkey: "UserPubkey", RouteHint: "RouteHint", Created: 1234})
		mockDb.On("GetInvoice", "1").Return(db.InvoiceList{Status: false})
		mockDb.On("AddAndUpdateBudget", mock.Anything).Return(db.PaymentHistory{})
		mockDb.On("UpdateInvoice", "1").Return(db.InvoiceList{}).Once()

		ro := chi.NewRouter()
		ro.Post("/poll/invoice/{paymentRequest}", bHandler.PollInvoice)

		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/poll/invoice/1", bytes.NewBufferString(`{}`))
		if err != nil {
			t.Fatal(err)
		}

		ro.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockHttpClient.AssertExpectations(t)
	})
}
