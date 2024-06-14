package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stakwork/sphinx-tribes/config"

	"github.com/go-chi/chi"
	"github.com/lib/pq"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	mocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTribesByOwner(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	tHandler := NewTribeHandler(mockDb)

	t.Run("Should test that all tribes that an owner did not delete are returned if all=true is added to the request query", func(t *testing.T) {
		// Mock data
		mockPubkey := "mock_pubkey"
		mockTribes := []db.Tribe{
			{UUID: "uuid", OwnerPubKey: mockPubkey, Deleted: false},
			{UUID: "uuid", OwnerPubKey: mockPubkey, Deleted: false},
		}
		mockDb.On("GetAllTribesByOwner", mock.Anything).Return(mockTribes).Once()

		// Create request with "all=true" query parameter
		req, err := http.NewRequest("GET", "/tribes_by_owner/"+mockPubkey+"?all=true", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Serve request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GetTribesByOwner)
		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		var responseData []db.Tribe
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.ElementsMatch(t, mockTribes, responseData)
	})

	t.Run("Should test that all tribes that are not unlisted by an owner are returned", func(t *testing.T) {
		// Mock data
		mockPubkey := "mock_pubkey"
		mockTribes := []db.Tribe{
			{UUID: "uuid", OwnerPubKey: mockPubkey, Unlisted: false},
			{UUID: "uuid", OwnerPubKey: mockPubkey, Unlisted: false},
		}
		mockDb.On("GetTribesByOwner", mock.Anything).Return(mockTribes)

		// Create request without "all=true" query parameter
		req, err := http.NewRequest("GET", "/tribes/"+mockPubkey, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Serve request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GetTribesByOwner)
		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		var responseData []db.Tribe
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.ElementsMatch(t, mockTribes, responseData)
	})
}

func TestGetTribe(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	tHandler := NewTribeHandler(mockDb)

	t.Run("Should test that a tribe can be returned when the right UUID is passed to the request parameter", func(t *testing.T) {
		// Mock data
		mockUUID := "valid_uuid"
		mockTribe := db.Tribe{
			UUID: mockUUID,
		}
		mockChannels := []db.Channel{
			{ID: 1, TribeUUID: mockUUID},
			{ID: 2, TribeUUID: mockUUID},
		}
		mockDb.On("GetTribe", mock.Anything).Return(mockTribe).Once()
		mockDb.On("GetChannelsByTribe", mock.Anything).Return(mockChannels).Once()

		// Serve request
		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", mockUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/"+mockUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler := http.HandlerFunc(tHandler.GetTribe)
		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		var responseData map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.Equal(t, mockTribe.UUID, responseData["uuid"])
	})

	t.Run("Should test that no tribe is returned when a nonexistent UUID is passed", func(t *testing.T) {
		// Mock data
		mockDb.ExpectedCalls = nil
		nonexistentUUID := "nonexistent_uuid"
		mockDb.On("GetTribe", nonexistentUUID).Return(db.Tribe{}).Once()
		mockDb.On("GetChannelsByTribe", mock.Anything).Return([]db.Channel{}).Once()

		// Serve request
		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", nonexistentUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/"+nonexistentUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler := http.HandlerFunc(tHandler.GetTribe)
		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		var responseData map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.Equal(t, "", responseData["uuid"])
	})
}

func TestGetTribesByAppUrl(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	tHandler := NewTribeHandler(mockDb)

	t.Run("Should test that a tribe is returned when the right app URL is passed", func(t *testing.T) {
		// Mock data
		mockAppURL := "valid_app_url"
		mockTribes := []db.Tribe{
			{UUID: "uuid", AppURL: mockAppURL},
			{UUID: "uuid", AppURL: mockAppURL},
		}
		mockDb.On("GetTribesByAppUrl", mockAppURL).Return(mockTribes).Once()

		// Serve request
		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("app_url", mockAppURL)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/app_url/"+mockAppURL, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler := http.HandlerFunc(tHandler.GetTribesByAppUrl)
		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		var responseData []db.Tribe
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.ElementsMatch(t, mockTribes, responseData)
	})
}

func TestDeleteTribe(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "owner_pubkey")
	mockDb := mocks.NewDatabase(t)
	tHandler := NewTribeHandler(mockDb)

	t.Run("Should test that the owner of a tribe can delete a tribe", func(t *testing.T) {
		// Mock data
		mockUUID := "valid_uuid"
		mockOwnerPubKey := "owner_pubkey"

		mockVerifyTribeUUID := func(uuid string, checkTimestamp bool) (string, error) {
			return mockOwnerPubKey, nil
		}
		mockDb.On("UpdateTribe", mock.Anything, map[string]interface{}{"deleted": true}).Return(true)

		tHandler.verifyTribeUUID = mockVerifyTribeUUID

		// Create and serve request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.DeleteTribe)

		req, err := http.NewRequestWithContext(ctx, "DELETE", "/tribe/"+mockUUID, nil)
		if err != nil {
			t.Fatal(err)
		}
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("uuid", "mockUUID")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		var responseData bool
		errors := json.Unmarshal(rr.Body.Bytes(), &responseData)
		assert.NoError(t, errors)
		assert.True(t, responseData)
	})

	t.Run("Should test that a 401 error is returned when a tribe is attempted to be deleted by someone other than the owner", func(t *testing.T) {
		// Mock data
		ctx := context.WithValue(context.Background(), auth.ContextKey, "pubkey")
		mockUUID := "valid_uuid"
		mockOwnerPubKey := "owner_pubkey"

		mockVerifyTribeUUID := func(uuid string, checkTimestamp bool) (string, error) {
			return mockOwnerPubKey, nil
		}

		tHandler.verifyTribeUUID = mockVerifyTribeUUID

		// Create and serve request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.DeleteTribe)

		req, err := http.NewRequestWithContext(ctx, "DELETE", "/tribe/"+mockUUID, nil)
		if err != nil {
			t.Fatal(err)
		}
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("uuid", "mockUUID")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestGetFirstTribeByFeed(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	tHandler := NewTribeHandler(mockDb)

	t.Run("Should test that a tribe can be gotten by passing the feed URL", func(t *testing.T) {
		// Mock data
		mockFeedURL := "valid_feed_url"
		mockTribe := db.Tribe{
			UUID: "valid_uuid",
		}
		mockChannels := []db.Channel{
			{ID: 1, TribeUUID: mockTribe.UUID},
		}

		mockDb.On("GetFirstTribeByFeedURL", mockFeedURL).Return(mockTribe).Once()
		mockDb.On("GetChannelsByTribe", mockTribe.UUID).Return(mockChannels).Once()

		// Create request with valid feed URL
		req, err := http.NewRequest("GET", "/tribe_by_feed?url="+mockFeedURL, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Serve request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GetFirstTribeByFeed)
		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		var responseData map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.Equal(t, mockTribe.UUID, responseData["uuid"])
	})
}

func TestSetTribePreview(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "owner_pubkey")
	mockDb := mocks.NewDatabase(t)
	tHandler := NewTribeHandler(mockDb)

	t.Run("Should test that the owner of a tribe can set tribe preview", func(t *testing.T) {
		// Mock data
		mockUUID := "valid_uuid"
		mockOwnerPubKey := "owner_pubkey"

		mockVerifyTribeUUID := func(uuid string, checkTimestamp bool) (string, error) {
			return mockOwnerPubKey, nil
		}
		mockDb.On("UpdateTribe", mock.Anything, map[string]interface{}{"preview": "preview"}).Return(true)

		tHandler.verifyTribeUUID = mockVerifyTribeUUID

		// Create and serve request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.SetTribePreview)

		req, err := http.NewRequestWithContext(ctx, "PUT", "/tribepreview/"+mockUUID+"?preview=preview", nil)
		if err != nil {
			t.Fatal(err)
		}
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("uuid", "mockUUID")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		var responseData bool
		errors := json.Unmarshal(rr.Body.Bytes(), &responseData)
		assert.NoError(t, errors)
		assert.True(t, responseData)
	})

	t.Run("Should test that a 401 error is returned when setting a tribe preview action by someone other than the owner", func(t *testing.T) {
		// Mock data
		ctx := context.WithValue(context.Background(), auth.ContextKey, "pubkey")
		mockUUID := "valid_uuid"
		mockOwnerPubKey := "owner_pubkey"

		mockVerifyTribeUUID := func(uuid string, checkTimestamp bool) (string, error) {
			return mockOwnerPubKey, nil
		}

		tHandler.verifyTribeUUID = mockVerifyTribeUUID

		// Create and serve request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.SetTribePreview)

		req, err := http.NewRequestWithContext(ctx, "PUT", "/tribepreview/"+mockUUID+"?preview=preview", nil)
		if err != nil {
			t.Fatal(err)
		}
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("uuid", "mockUUID")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestCreateOrEditTribe(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	tHandler := NewTribeHandler(mockDb)

	t.Run("Should test that a tribe can be created when the right data is passed", func(t *testing.T) {
		// Mock data
		mockPubKey := "valid_pubkey"
		mockUUID := "valid_uuid"
		mockName := "Test Tribe"
		mockDescription := "This is a test tribe."
		mockTags := []string{"tag1", "tag2"}

		mockVerifyTribeUUID := func(uuid string, checkTimestamp bool) (string, error) {
			return mockPubKey, nil
		}

		tHandler.verifyTribeUUID = mockVerifyTribeUUID

		// Mock request body
		requestBody := map[string]interface{}{
			"UUID":        mockUUID,
			"Name":        mockName,
			"Description": mockDescription,
			"Tags":        mockTags,
		}
		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatal(err)
		}

		// Mock database calls
		mockDb.On("GetTribe", mock.Anything).Return(db.Tribe{
			UUID:        mockUUID,
			OwnerPubKey: mockPubKey,
		}).Once()
		mockDb.On("CreateOrEditTribe", mock.Anything).Return(db.Tribe{
			UUID: mockUUID,
		}, nil)

		// Create request with mock body
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(requestBodyBytes))
		if err != nil {
			t.Fatal(err)
		}

		// Set context with mock pub key
		ctx := context.WithValue(req.Context(), auth.ContextKey, mockPubKey)
		req = req.WithContext(ctx)

		// Serve request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.CreateOrEditTribe)
		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		var responseData map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.Equal(t, mockUUID, responseData["uuid"])
	})
}

func TestGetTribeByUniqueName(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	tHandler := NewTribeHandler(mockDb)

	t.Run("Should test that a tribe can be fetched by its unique name", func(t *testing.T) {
		// Mock data
		mockUniqueName := "test_tribe"
		mockTribe := db.Tribe{
			UniqueName: mockUniqueName,
			UUID:       "valid_uuid",
		}
		mockChannels := []db.Channel{
			{ID: 1, TribeUUID: "UUID"},
			{ID: 2, TribeUUID: "UUID"},
		}

		// Mock database calls
		mockDb.On("GetTribeByUniqueName", mock.Anything).Return(mockTribe)
		mockDb.On("GetChannelsByTribe", mock.Anything).Return(mockChannels).Once()

		// Create request with mock unique name
		req, err := http.NewRequest("GET", "/tribe_by_un/"+mockUniqueName, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Serve request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GetTribeByUniqueName)
		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		var responseData map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.Equal(t, mockUniqueName, responseData["unique_name"])
	})
}

func TestGetAllTribes(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	tHandler := NewTribeHandler(mockDb)
	t.Run("should return all tribes", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GetAllTribes)

		expectedTribes := []db.Tribe{
			{UUID: "uuid", Name: "Tribe1"},
			{UUID: "uuid", Name: "Tribe2"},
			{UUID: "uuid", Name: "Tribe3"},
		}

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/", nil)
		assert.NoError(t, err)

		mockDb.On("GetAllTribes", mock.Anything).Return(expectedTribes)
		handler.ServeHTTP(rr, req)
		var returnedTribes []db.Tribe
		err = json.Unmarshal(rr.Body.Bytes(), &returnedTribes)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, expectedTribes, returnedTribes)
		mockDb.AssertExpectations(t)

	})
}

func TestGetTotalTribes(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	tHandler := NewTribeHandler(mockDb)
	t.Run("should return the total number of tribes", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GetTotalribes)

		expectedTribes := []db.Tribe{
			{UUID: "uuid", Name: "Tribe1"},
			{UUID: "uuid", Name: "Tribe2"},
			{UUID: "uuid", Name: "Tribe3"},
		}

		expectedTribesCount := int64(len(expectedTribes))

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/total", nil)
		assert.NoError(t, err)

		mockDb.On("GetTribesTotal", mock.Anything).Return(expectedTribesCount)

		handler.ServeHTTP(rr, req)
		var returnedTribesCount int64
		err = json.Unmarshal(rr.Body.Bytes(), &returnedTribesCount)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, expectedTribesCount, returnedTribesCount)
		mockDb.AssertExpectations(t)

	})
}

func TestGetListedTribes(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	tHandler := NewTribeHandler(mockDb)

	t.Run("should only return tribes associated with a passed tag query", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GetListedTribes)
		expectedTribes := []db.Tribe{
			{UUID: "1", Name: "Tribe 1", Tags: pq.StringArray{"tag1", "tag2", "tag3"}},
			{UUID: "2", Name: "Tribe 2", Tags: pq.StringArray{"tag4", "tag5"}},
			{UUID: "3", Name: "Tribe 3", Tags: pq.StringArray{"tag6", "tag7", "tag8"}},
		}
		req, err := http.NewRequest("GET", "/tribes", nil)
		if err != nil {
			t.Fatal(err)
		}
		query := req.URL.Query()
		tagVals := pq.StringArray{"tag1", "tag4", "tag7"}
		tags := strings.Join(tagVals, ",")
		query.Set("tags", tags)
		req.URL.RawQuery = query.Encode()
		if err != nil {
			t.Fatal(err)
		}

		mockDb.On("GetListedTribes", req).Return(expectedTribes)
		handler.ServeHTTP(rr, req)
		var returnedTribes []db.Tribe
		err = json.Unmarshal(rr.Body.Bytes(), &returnedTribes)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, expectedTribes, returnedTribes)

	})

	t.Run("should return all tribes when no tag queries are passed", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GetListedTribes)
		expectedTribes := []db.Tribe{
			{UUID: "1", Name: "Tribe 1", Tags: pq.StringArray{"tag1", "tag2", "tag3"}},
			{UUID: "2", Name: "Tribe 2", Tags: pq.StringArray{"tag4", "tag5"}},
			{UUID: "3", Name: "Tribe 3", Tags: pq.StringArray{"tag6", "tag7", "tag8"}},
		}

		req, err := http.NewRequest("GET", "/tribes", nil)
		if err != nil {
			t.Fatal(err)
		}
		query := req.URL.Query()
		tagVals := pq.StringArray{"tag1", "tag4", "tag7"}
		tags := strings.Join(tagVals, ",")
		query.Set("tags", tags)
		req.URL.RawQuery = query.Encode()
		if err != nil {
			t.Fatal(err)
		}

		mockDb.On("GetListedTribes", req).Return(expectedTribes)
		handler.ServeHTTP(rr, req)

		var returnedTribes []db.Tribe
		err = json.Unmarshal(rr.Body.Bytes(), &returnedTribes)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, expectedTribes, returnedTribes)

	})

}

func TestGenerateBudgetInvoice(t *testing.T) {
	ctx := context.Background()
	mockDb := mocks.NewDatabase(t)
	tHandler := NewTribeHandler(mockDb)
	authorizedCtx := context.WithValue(ctx, auth.ContextKey, "valid-key")

	userAmount := uint(1000)
	invoiceResponse := db.InvoiceResponse{
		Succcess: true,
		Response: db.Invoice{
			Invoice: "example_invoice",
		},
	}

	t.Run("Should test that a wrong Post body returns a 406 error", func(t *testing.T) {
		invalidBody := []byte(`"key": "value"`)
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/budgetinvoices", bytes.NewBuffer(invalidBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GenerateBudgetInvoice)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("Should mock a call to relay /invoices with the correct body", func(t *testing.T) {
		mockDb.On("ProcessBudgetInvoice", mock.AnythingOfType("db.NewPaymentHistory"), mock.AnythingOfType("db.NewInvoiceList")).Return(nil)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			expectedBody := map[string]interface{}{"amount": float64(0), "memo": "Budget Invoice"}
			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			assert.NoError(t, err)

			assert.Equal(t, expectedBody, body)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "success"})
		}))
		defer ts.Close()

		config.RelayUrl = ts.URL

		reqBody := map[string]interface{}{"amount": 0}
		bodyBytes, _ := json.Marshal(reqBody)

		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/budgetinvoices", bytes.NewBuffer(bodyBytes))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GenerateBudgetInvoice)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should test that the amount passed by the user is equal to the amount sent for invoice generation", func(t *testing.T) {

		userAmount := float64(1000)

		mockDb.On("ProcessBudgetInvoice", mock.AnythingOfType("db.NewPaymentHistory"), mock.AnythingOfType("db.NewInvoiceList")).Return(nil)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			assert.NoError(t, err)

			assert.Equal(t, userAmount, body["amount"])

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "success"})
		}))
		defer ts.Close()

		config.RelayUrl = ts.URL

		reqBody := map[string]interface{}{"amount": userAmount}
		bodyBytes, _ := json.Marshal(reqBody)

		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/budgetinvoices", bytes.NewBuffer(bodyBytes))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GenerateBudgetInvoice)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should add payments to the payment history and invoice to the invoice list upon successful relay call", func(t *testing.T) {
		mockDb.On("ProcessBudgetInvoice", mock.AnythingOfType("db.NewPaymentHistory"), mock.AnythingOfType("db.NewInvoiceList")).Return(nil)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(invoiceResponse)
		}))
		defer ts.Close()

		config.RelayUrl = ts.URL

		reqBody := map[string]interface{}{"amount": userAmount}
		bodyBytes, _ := json.Marshal(reqBody)
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/budgetinvoices", bytes.NewBuffer(bodyBytes))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GenerateBudgetInvoice)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response db.InvoiceResponse
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Succcess, "Invoice generation should be successful")
		assert.Equal(t, "example_invoice", response.Response.Invoice, "The invoice in the response should match the mock")

		mockDb.AssertCalled(t, "ProcessBudgetInvoice", mock.AnythingOfType("db.NewPaymentHistory"), mock.AnythingOfType("db.NewInvoiceList"))
	})
}
