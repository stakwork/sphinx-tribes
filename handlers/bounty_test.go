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
	"testing"
	"time"

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

	t.Run("Should return if bounty is present in db", func(t *testing.T) {
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
		req, _ := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "gobounties/created/1707991475", nil)
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
		req, _ := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "people/wanteds/assigned/clu80datu2rjujsmim40?sortBy=paid&page=1&limit=20&search=", nil)

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
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)
	t.Run("Should successfull Get Person Created Bounties", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetPersonCreatedBounties)
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
		req, _ := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "people/wanteds/created/clu80datu2rjujsmim40?sortBy=paid&page=1&limit=20&search=", nil)

		mockDb.On("GetCreatedBounties", req).Return([]db.Bounty{bounty}, nil).Once()
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

func TestGetAllBounties(t *testing.T) {
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)

	t.Run("successful retrieval of bounties", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetAllBounties)
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
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/all", nil)
		assert.NoError(t, err)
		mockDb.On("GetAllBounties", mock.Anything).Return([]db.Bounty{bounty}, nil).Once()
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

}
