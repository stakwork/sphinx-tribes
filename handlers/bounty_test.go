package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/handlers/mocks"
	dbMocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var bountyOwner = db.Person{
	Uuid:        "user_3_uuid",
	OwnerAlias:  "user3",
	UniqueName:  "user3",
	OwnerPubKey: "user_3_pubkey",
	PriceToMeet: 0,
	Description: "this is test user 3",
}

var bountyAssignee = db.Person{
	Uuid:        "user_4_uuid",
	OwnerAlias:  "user4",
	UniqueName:  "user4",
	OwnerPubKey: "user_4_pubkey",
	PriceToMeet: 0,
	Description: "this is user 4",
}

var bountyPrev = db.NewBounty{
	Type:          "coding",
	Title:         "Previous bounty",
	Description:   "Previous bounty description",
	OrgUuid:       "org-4",
	WorkspaceUuid: "work-4",
	Assignee:      bountyAssignee.OwnerPubKey,
	OwnerID:       bountyOwner.OwnerPubKey,
	Show:          true,
	Created:       111111111,
}

var bountyNext = db.NewBounty{
	Type:          "coding",
	Title:         "Next bounty",
	Description:   "Next bounty description",
	WorkspaceUuid: "work-4",
	Assignee:      "",
	OwnerID:       bountyOwner.OwnerPubKey,
	Show:          true,
	Created:       111111112,
}

var workspace = db.Workspace{
	Uuid:        "workspace_uuid13",
	Name:        "TestWorkspace",
	Description: "This is a test workspace",
	OwnerPubKey: bountyOwner.OwnerPubKey,
	Img:         "",
	Website:     "",
}

var workBountyPrev = db.NewBounty{
	Type:          "coding",
	Title:         "Workspace Previous bounty",
	Description:   "Workspace Previous bounty description",
	WorkspaceUuid: workspace.Uuid,
	Assignee:      bountyAssignee.OwnerPubKey,
	OwnerID:       bountyOwner.OwnerPubKey,
	Show:          true,
	Created:       111111113,
}

var workBountyNext = db.NewBounty{
	Type:          "coding",
	Title:         "Workpace Next bounty",
	Description:   "Workspace Next bounty description",
	WorkspaceUuid: workspace.Uuid,
	Assignee:      "",
	OwnerID:       bountyOwner.OwnerPubKey,
	Show:          true,
	Created:       111111114,
}

func SetupSuite(_ *testing.T) func(tb testing.TB) {
	db.InitTestDB()

	return func(_ testing.TB) {
		log.Println("Teardown test")
	}
}

func AddExisitingDB(existingBounty db.NewBounty) {
	bounty := db.TestDB.GetBounty(1)
	if bounty.ID == 0 {
		// add existing bounty to db
		db.TestDB.CreateOrEditBounty(existingBounty)
	}
}

func TestCreateOrEditBounty(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	existingBounty := db.NewBounty{
		Type:          "coding",
		Title:         "existing bounty",
		Description:   "existing bounty description",
		WorkspaceUuid: "work-1",
		OwnerID:       "first-user",
		Price:         2000,
	}

	// Add initial Bounty
	AddExisitingDB(existingBounty)

	newBounty := db.NewBounty{
		Type:          "coding",
		Title:         "new bounty",
		Description:   "new bounty description",
		WorkspaceUuid: "work-1",
		OwnerID:       "test-key",
		Price:         1500,
	}

	failedBounty := db.NewBounty{
		Type:          "coding",
		Title:         "new bounty",
		Description:   "failed bounty description",
		WorkspaceUuid: "work-1",
		Price:         1500,
	}

	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockClient := mocks.NewHttpClient(t)
	mockUserHasManageBountyRolesTrue := func(pubKeyFromAuth string, uuid string) bool {
		return true
	}
	mockUserHasManageBountyRolesFalse := func(pubKeyFromAuth string, uuid string) bool {
		return false
	}
	bHandler := NewBountyHandler(mockClient, db.TestDB)

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

	t.Run("return error if trying to update other user's bounty", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)
		bHandler.userHasManageBountyRoles = mockUserHasManageBountyRolesFalse

		updatedBounty := existingBounty
		updatedBounty.ID = 1
		updatedBounty.Show = true
		updatedBounty.WorkspaceUuid = ""

		json, err := json.Marshal(updatedBounty)
		if err != nil {
			fmt.Println("Could not marshal json data")
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(json))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, strings.TrimRight(rr.Body.String(), "\n"), "Cannot edit another user's bounty")
	})

	t.Run("return error if user does not have required roles", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)
		bHandler.userHasManageBountyRoles = mockUserHasManageBountyRolesFalse

		updatedBounty := existingBounty
		updatedBounty.Title = "Existing bounty updated"
		updatedBounty.ID = 1

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
		bHandler.userHasManageBountyRoles = mockUserHasManageBountyRolesTrue

		updatedBounty := existingBounty
		updatedBounty.Title = "first bounty updated"
		updatedBounty.ID = 1

		body, _ := json.Marshal(updatedBounty)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		bounty := db.TestDB.GetBounty(1)
		assert.Equal(t, bounty.Title, updatedBounty.Title)
	})

	t.Run("should not update created at when bounty is updated", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)
		bHandler.userHasManageBountyRoles = mockUserHasManageBountyRolesTrue

		updatedBounty := existingBounty
		updatedBounty.Title = "second bounty updated"

		body, _ := json.Marshal(updatedBounty)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var returnedBounty db.Bounty
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBounty)
		assert.NoError(t, err)
		assert.NotEqual(t, returnedBounty.Created, returnedBounty.Updated)
		// Check the response body or any other expected behavior
	})

	t.Run("should return error if failed to add new bounty", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)
		body, _ := json.Marshal(failedBounty)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("add bounty if error not present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBounty)

		body, _ := json.Marshal(newBounty)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
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
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	existingBounty := db.NewBounty{
		Type:          "coding",
		Title:         "existing bounty",
		Description:   "existing bounty description",
		WorkspaceUuid: "work-1",
		OwnerID:       "first-user",
		Price:         2000,
	}

	// Add initial Bounty
	AddExisitingDB(existingBounty)

	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)
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

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("pubkey", "pub-key")
		rctx.URLParams.Add("created", "1111")

		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/pub-key/createdAt", nil)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("should successfully delete bounty from db", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.DeleteBounty)
		existingBounty := db.TestDB.GetBounty(1)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("pubkey", existingBounty.OwnerID)

		created := fmt.Sprintf("%d", existingBounty.Created)
		rctx.URLParams.Add("created", created)

		route := fmt.Sprintf("/%s/%d", existingBounty.OwnerID, existingBounty.Created)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, route, nil)

		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rr, req)

		// get Bounty from DB
		checkBounty := db.TestDB.GetBounty(1)
		// chcek that the bounty's ID is now zero
		assert.Equal(t, 0, int(checkBounty.ID))
	})
}

func TestGetBountyByCreated(t *testing.T) {
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, mockDb)

	t.Run("Should return bounty by its created value", func(t *testing.T) {
		mockGenerateBountyResponse := func(bounties []db.NewBounty) []db.BountyResponse {
			var bountyResponses []db.BountyResponse

			for _, bounty := range bounties {
				owner := db.Person{
					ID: 1,
				}
				assignee := db.Person{
					ID: 1,
				}
				workspace := db.WorkspaceShort{
					Uuid: "uuid",
				}

				bountyResponse := db.BountyResponse{
					Bounty:       bounty,
					Assignee:     assignee,
					Owner:        owner,
					Organization: workspace,
					Workspace:    workspace,
				}
				bountyResponses = append(bountyResponses, bountyResponse)
			}

			return bountyResponses
		}
		bHandler.generateBountyResponse = mockGenerateBountyResponse

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyByCreated)
		bounty := db.NewBounty{
			ID:            1,
			Type:          "coding",
			Title:         "first bounty",
			Description:   "first bounty description",
			OrgUuid:       "org-1",
			WorkspaceUuid: "work-1",
			Assignee:      "user1",
			Created:       1707991475,
			OwnerID:       "owner-1",
		}
		createdStr := strconv.FormatInt(bounty.Created, 10)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("created", "1707991475")
		req, _ := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/created/1707991475", nil)
		mockDb.On("GetBountyDataByCreated", createdStr).Return([]db.NewBounty{bounty}, nil).Once()
		mockDb.On("GetPersonByPubkey", "owner-1").Return(db.Person{}).Once()
		mockDb.On("GetPersonByPubkey", "user1").Return(db.Person{}).Once()
		mockDb.On("GetWorkspaceByUuid", "work-1").Return(db.Workspace{}).Once()
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

		mockDb.On("GetBountyDataByCreated", createdStr).Return([]db.NewBounty{}, nil).Once()

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code, "Expected 404 Not Found for nonexistent bounty")

		mockDb.AssertExpectations(t)
	})

}

func TestGetPersonAssignedBounties(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	bountyOwner := db.Person{
		Uuid:        "user_1_uuid",
		OwnerAlias:  "user1",
		UniqueName:  "user1",
		OwnerPubKey: "user_1_pubkey",
		PriceToMeet: 0,
		Description: "this is test user 1",
	}

	bountyAssignee := db.Person{
		Uuid:        "user_2_uuid",
		OwnerAlias:  "user2",
		UniqueName:  "user2",
		OwnerPubKey: "user_2_pubkey",
		PriceToMeet: 0,
		Description: "this is user 2",
	}

	bounty := db.NewBounty{
		Type:          "coding",
		Title:         "first bounty",
		Description:   "first bounty description",
		OrgUuid:       "org-1",
		WorkspaceUuid: "work-1",
		Assignee:      bountyAssignee.OwnerPubKey,
		OwnerID:       bountyOwner.OwnerPubKey,
		Show:          true,
	}

	t.Run("Should successfull Get Person Assigned Bounties", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetPersonAssignedBounties)

		// create users
		db.TestDB.CreateOrEditPerson(bountyOwner)
		db.TestDB.CreateOrEditPerson(bountyAssignee)

		// create bounty
		db.TestDB.CreateOrEditBounty(bounty)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", bountyAssignee.Uuid)
		rctx.URLParams.Add("sortBy", "paid")
		rctx.URLParams.Add("page", "0")
		rctx.URLParams.Add("limit", "20")
		rctx.URLParams.Add("search", "")

		route := fmt.Sprintf("/people/wanteds/assigned/%s?sortBy=paid&page=0&limit=20&search=''", bountyAssignee.Uuid)
		req, _ := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, route, nil)

		handler.ServeHTTP(rr, req)

		// bounty from db
		expectedBounty, _ := db.TestDB.GetAssignedBounties(req)

		var returnedBounty []db.BountyResponse
		err := json.Unmarshal(rr.Body.Bytes(), &returnedBounty)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotEmpty(t, returnedBounty)
		assert.Equal(t, len(expectedBounty), len(returnedBounty))
	})
}

func TestGetPersonCreatedBounties(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	ctx := context.Background()
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	bounty := db.NewBounty{
		Type:          "coding",
		Title:         "first bounty 3",
		Description:   "first bounty description",
		WorkspaceUuid: "work-4",
		Assignee:      bountyAssignee.OwnerPubKey,
		OwnerID:       bountyOwner.OwnerPubKey,
		Show:          true,
	}

	bounty2 := db.NewBounty{
		Type:          "coding 2",
		Title:         "second bounty 3",
		Description:   "second bounty description 2",
		OrgUuid:       "org-4",
		WorkspaceUuid: "work-4",
		Assignee:      bountyAssignee.OwnerPubKey,
		OwnerID:       bountyOwner.OwnerPubKey,
		Show:          true,
		Created:       11111111,
	}

	bounty3 := db.NewBounty{
		Type:          "coding 2",
		Title:         "second bounty 4",
		Description:   "second bounty description 2",
		WorkspaceUuid: "work-4",
		Assignee:      "",
		OwnerID:       bountyOwner.OwnerPubKey,
		Show:          true,
		Created:       2222222,
	}

	// create users
	db.TestDB.CreateOrEditPerson(bountyOwner)
	db.TestDB.CreateOrEditPerson(bountyAssignee)

	// create bounty
	db.TestDB.CreateOrEditBounty(bounty)
	db.TestDB.CreateOrEditBounty(bounty2)
	db.TestDB.CreateOrEditBounty(bounty3)

	t.Run("should return bounties created by the user", func(t *testing.T) {
		rr := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", bountyOwner.Uuid)

		route := fmt.Sprintf("/people/wanteds/created/%s?sortBy=paid&page=1&limit=20&search=''", bountyOwner.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, route, nil)
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

		// bounty from db
		expectedBounty, _ := db.TestDB.GetCreatedBounties(req)

		assert.NotEmpty(t, responseData)
		assert.Equal(t, len(expectedBounty), len(responseData))
	})

	t.Run("should not return bounties created by other users", func(t *testing.T) {
		rr := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", bountyAssignee.Uuid)

		route := fmt.Sprintf("/people/wanteds/created/%s?sortBy=paid&page=1&limit=20&search=''", bountyAssignee.Uuid)
		req, err := http.NewRequest("GET", route, nil)
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
		rr := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", bountyOwner.Uuid)

		route := fmt.Sprintf("/people/wanteds/created/%s?Assigned=true&page=1&limit=2", bountyOwner.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, route, nil)
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

		assert.Len(t, responseData, 2)

		// Assert that bounties are filtered correctly
		// bounty from db
		expectedBounty, _ := db.TestDB.GetCreatedBounties(req)
		assert.Equal(t, len(expectedBounty), len(responseData))
	})
}

func TestGetNextBountyByCreated(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.TestDB.CreateOrEditBounty(bountyPrev)
	db.TestDB.CreateOrEditBounty(bountyNext)

	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	t.Run("Should test that the next bounty on the bounties homepage can be gotten by its created value and the selected filters", func(t *testing.T) {
		rr := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		created := fmt.Sprintf("%d", bountyPrev.Created)
		rctx.URLParams.Add("created", created)

		route := fmt.Sprintf("/next/%d", bountyPrev.Created)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, route, nil)
		if err != nil {
			t.Fatal(err)
		}
		bHandler.GetNextBountyByCreated(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseData uint
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.Greater(t, responseData, uint(1))
	})
}

func TestGetPreviousBountyByCreated(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	t.Run("Should test that the previous bounty on the bounties homepage can be gotten by its created value and the selected filters", func(t *testing.T) {
		rr := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		created := fmt.Sprintf("%d", bountyPrev.Created)
		rctx.URLParams.Add("created", created)

		route := fmt.Sprintf("/previous/%d", bountyNext.Created)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, route, nil)
		if err != nil {
			t.Fatal(err)
		}

		bHandler.GetPreviousBountyByCreated(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseData uint
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.Greater(t, responseData, uint(1))
	})
}

func TestGetWorkspaceNextBountyByCreated(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.TestDB.CreateOrEditWorkspace(workspace)
	db.TestDB.CreateOrEditBounty(workBountyPrev)
	db.TestDB.CreateOrEditBounty(workBountyNext)

	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	t.Run("Should test that the next bounty on the workspace bounties homepage can be gotten by its created value and the selected filters", func(t *testing.T) {
		rr := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		created := fmt.Sprintf("%d", workBountyPrev.Created)
		rctx.URLParams.Add("created", created)
		rctx.URLParams.Add("uuid", workspace.Uuid)

		route := fmt.Sprintf("/org/next/%s/%d", workspace.Uuid, workBountyPrev.Created)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, route, nil)
		if err != nil {
			t.Fatal(err)
		}

		bHandler.GetWorkspaceNextBountyByCreated(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseData uint
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.Greater(t, responseData, uint(2))
	})
}

func TestGetWorkspacePreviousBountyByCreated(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	t.Run("Should test that the previous bounty on the workspace bounties homepage can be gotten by its created value and the selected filters", func(t *testing.T) {
		rr := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		created := fmt.Sprintf("%d", workBountyNext.Created)
		rctx.URLParams.Add("created", created)
		rctx.URLParams.Add("uuid", workspace.Uuid)

		route := fmt.Sprintf("/org/previous/%s/%d", workspace.Uuid, workBountyNext.Created)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, route, nil)
		if err != nil {
			t.Fatal(err)
		}

		bHandler.GetWorkspacePreviousBountyByCreated(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseData uint
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.Greater(t, responseData, uint(2))
	})
}

func TestGetBountyById(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	t.Run("successful retrieval of bounty by ID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyById)

		now := time.Now().Unix()
		bounty := db.NewBounty{
			Type:          "coding",
			Title:         "Bounty With ID",
			Description:   "Bounty ID description",
			WorkspaceUuid: "",
			Assignee:      "",
			OwnerID:       bountyOwner.OwnerPubKey,
			Show:          true,
			Created:       now,
		}

		db.TestDB.CreateOrEditBounty(bounty)

		bountyInDb, err := db.TestDB.GetBountyByCreated(uint(bounty.Created))
		assert.NoError(t, err)
		assert.NotNil(t, bountyInDb)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("bountyId", strconv.Itoa(int(bountyInDb.ID)))
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/bounty/"+strconv.Itoa(int(bountyInDb.ID)), nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var returnedBounty []db.BountyResponse
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBounty)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotEmpty(t, returnedBounty)
	})

	t.Run("bounty not found", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyById)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("bountyId", "Invalid-id")
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/bounty/Invalid-id", nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestGetBountyIndexById(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	t.Run("successful retrieval of bounty by Index ID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyIndexById)

		now := time.Now().Unix()
		bounty := db.NewBounty{
			ID:            1,
			Type:          "coding",
			Title:         "Bounty With ID",
			Description:   "Bounty description",
			WorkspaceUuid: "",
			Assignee:      "",
			OwnerID:       bountyOwner.OwnerPubKey,
			Show:          true,
			Created:       now,
		}

		db.TestDB.CreateOrEditBounty(bounty)

		bountyInDb, err := db.TestDB.GetBountyByCreated(uint(bounty.Created))
		assert.Equal(t, bounty, bountyInDb)
		assert.NoError(t, err)

		bountyIndex := db.TestDB.GetBountyIndexById(strconv.Itoa(int(bountyInDb.ID)))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("bountyId", strconv.Itoa(int(bountyInDb.ID)))
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/index/"+strconv.Itoa(int(bountyInDb.ID)), nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		responseBody := rr.Body.Bytes()
		responseString := strings.TrimSpace(string(responseBody))
		returnedIndex, err := strconv.ParseInt(responseString, 10, 64)
		assert.NoError(t, err)
		assert.Equal(t, bountyIndex, returnedIndex)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("bounty index by ID not found", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyIndexById)

		bountyID := ""
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("bountyId", bountyID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/index/"+bountyID, nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestGetAllBounties(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	t.Run("Should successfully return all bounties", func(t *testing.T) {
		now := time.Now().Unix()
		bounty := db.NewBounty{
			Type:          "coding",
			Title:         "Bounty With ID",
			Description:   "Bounty ID description",
			WorkspaceUuid: "",
			Assignee:      "",
			OwnerID:       "test-owner",
			Show:          true,
			Created:       now,
		}
		db.TestDB.CreateOrEditBounty(bounty)

		bountyInDb, err := db.TestDB.GetBountyByCreated(uint(bounty.Created))
		assert.NoError(t, err)
		assert.NotNil(t, bountyInDb)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetAllBounties)

		rctx := chi.NewRouteContext()
		req, _ := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/all", nil)

		handler.ServeHTTP(rr, req)

		var returnedBounty []db.BountyResponse
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBounty)
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

	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mockHttpClient := &mocks.HttpClient{}
	mockUserHasAccessTrue := func(pubKeyFromAuth string, uuid string, role string) bool {
		return true
	}
	mockUserHasAccessFalse := func(pubKeyFromAuth string, uuid string, role string) bool {
		return false
	}
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
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	var mutex sync.Mutex
	var processingTimes []time.Time

	now := time.Now().Unix()
	bountyOwnerId := "owner_pubkey"

	db.TestDB.DeleteWorkspaceBudget()
	db.TestDB.DeleteBounty(bountyOwnerId, strconv.FormatInt(now, 10))

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
		OwnerPubKey: person.OwnerPubKey,
		Github:      "gtihub",
		Website:     "website",
		Description: "description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	budgetAmount := uint(5000)
	bountyBudget := db.NewBountyBudget{
		WorkspaceUuid: workspace.Uuid,
		TotalBudget:   budgetAmount,
	}
	db.TestDB.CreateWorkspaceBudget(bountyBudget)

	bountyAmount := uint(3000)
	bounty := db.NewBounty{
		OwnerID:       bountyOwnerId,
		Price:         bountyAmount,
		Created:       now,
		Type:          "coding",
		Title:         "bountyTitle",
		Description:   "bountyDescription",
		Assignee:      person.OwnerPubKey,
		Show:          true,
		WorkspaceUuid: workspace.Uuid,
		Paid:          false,
	}
	db.TestDB.CreateOrEditBounty(bounty)

	dbBounty, err := db.TestDB.GetBountyDataByCreated(strconv.FormatInt(bounty.Created, 10))
	if err != nil {
		t.Fatal(err)
	}

	bountyId := dbBounty[0].ID
	bountyIdStr := strconv.FormatInt(int64(bountyId), 10)

	unauthorizedCtx := context.WithValue(ctx, auth.ContextKey, "")
	authorizedCtx := context.WithValue(ctx, auth.ContextKey, person.OwnerPubKey)

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
		req, err := http.NewRequestWithContext(unauthorizedCtx, http.MethodPost, "/gobounties/pay/"+bountyIdStr, nil)

		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected 401 Unauthorized for unauthorized access")
	})

	t.Run("401 error if user not workspace admin or does not have PAY BOUNTY role", func(t *testing.T) {
		bHandler.userHasAccess = mockUserHasAccessFalse

		r := chi.NewRouter()
		r.Post("/gobounties/pay/{id}", bHandler.MakeBountyPayment)

		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(unauthorizedCtx, http.MethodPost, "/gobounties/pay/"+bountyIdStr, bytes.NewBufferString(`{}`))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected 401 Unauthorized when the user lacks the PAY BOUNTY role")

	})

	t.Run("Should test that an error WebSocket message is sent if the payment fails", func(t *testing.T) {
		mockHttpClient2 := &mocks.HttpClient{}

		bHandler2 := NewBountyHandler(mockHttpClient2, db.TestDB)
		bHandler2.getSocketConnections = mockGetSocketConnections
		bHandler2.userHasAccess = mockUserHasAccessTrue

		expectedUrl := fmt.Sprintf("%s/payment", config.RelayUrl)
		expectedBody := fmt.Sprintf(`{"amount": %d, "destination_key": "%s", "text": "memotext added for notification"}`, bountyAmount, person.OwnerPubKey)

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
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/gobounties/pay/"+bountyIdStr, requestBody)
		if err != nil {
			t.Fatal(err)
		}

		ro.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockHttpClient2.AssertExpectations(t)
	})

	t.Run("Should test that a successful WebSocket message is sent if the payment is successful", func(t *testing.T) {
		bHandler.getSocketConnections = mockGetSocketConnections
		bHandler.userHasAccess = mockUserHasAccessTrue

		expectedUrl := fmt.Sprintf("%s/payment", config.RelayUrl)
		expectedBody := fmt.Sprintf(`{"amount": %d, "destination_key": "%s", "text": "memotext added for notification"}`, bountyAmount, person.OwnerPubKey)

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
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/gobounties/pay/"+bountyIdStr, requestBody)
		if err != nil {
			t.Fatal(err)
		}

		ro.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("405 when trying to pay an already-paid bounty", func(t *testing.T) {
		r := chi.NewRouter()
		r.Post("/gobounties/pay/{id}", bHandler.MakeBountyPayment)

		requestBody := bytes.NewBuffer([]byte("{}"))
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/gobounties/pay/"+bountyIdStr, requestBody)
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Expected 405 Method Not Allowed for an already-paid bounty")
	})

	t.Run("403 error when amount exceeds workspace's budget balance", func(t *testing.T) {
		db.TestDB.DeleteBounty(bountyOwnerId, strconv.FormatInt(now, 10))
		bounty.Paid = false
		db.TestDB.CreateOrEditBounty(bounty)

		dbBounty, err := db.TestDB.GetBountyDataByCreated(strconv.FormatInt(bounty.Created, 10))
		if err != nil {
			t.Fatal(err)
		}

		bountyId := dbBounty[0].ID
		bountyIdStr := strconv.FormatInt(int64(bountyId), 10)

		mockHttpClient := mocks.NewHttpClient(t)
		bHandler := NewBountyHandler(mockHttpClient, db.TestDB)
		bHandler.userHasAccess = mockUserHasAccessTrue

		r := chi.NewRouter()
		r.Post("/gobounties/pay/{id}", bHandler.MakeBountyPayment)

		requestBody := bytes.NewBuffer([]byte("{}"))
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/gobounties/pay/"+bountyIdStr, requestBody)
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code, "Expected 403 Forbidden when the payment exceeds the workspace's budget")
	})
}

func TestBountyBudgetWithdraw(t *testing.T) {
	ctx := context.Background()
	mockDb := dbMocks.NewDatabase(t)
	mockHttpClient := mocks.NewHttpClient(t)
	mockUserHasAccessTrue := func(pubKeyFromAuth string, uuid string, role string) bool {
		return true
	}
	mockUserHasAccessFalse := func(pubKeyFromAuth string, uuid string, role string) bool {
		return false
	}
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

	t.Run("401 error if user is not the workspace admin or does not have WithdrawBudget role", func(t *testing.T) {
		bHandler.userHasAccess = mockUserHasAccessFalse

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.BountyBudgetWithdraw)

		validData := []byte(`{"orgUuid": "org-1", "paymentRequest": "invoice"}`)
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/budget/withdraw", bytes.NewReader(validData))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "You don't have appropriate permissions to withdraw bounty budget")
	})

	t.Run("403 error when amount exceeds workspace's budget", func(t *testing.T) {
		ctxs := context.WithValue(context.Background(), auth.ContextKey, "valid-key")
		mockDb := dbMocks.NewDatabase(t)
		mockHttpClient := mocks.NewHttpClient(t)
		bHandler := NewBountyHandler(mockHttpClient, mockDb)
		bHandler.userHasAccess = mockUserHasAccessTrue

		mockDb.On("GetWorkspaceBudget", "org-1").Return(db.NewBountyBudget{
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

		assert.Equal(t, http.StatusForbidden, rr.Code, "Expected 403 Forbidden when the payment exceeds the workspace's budget")
		assert.Contains(t, rr.Body.String(), "Workspace budget is not enough to withdraw the amount", "Expected specific error message")
	})

	t.Run("budget invoices get paid if amount is lesser than workspace's budget", func(t *testing.T) {
		ctxs := context.WithValue(context.Background(), auth.ContextKey, "valid-key")
		mockDb := dbMocks.NewDatabase(t)
		mockHttpClient := mocks.NewHttpClient(t)
		bHandler := NewBountyHandler(mockHttpClient, mockDb)
		bHandler.userHasAccess = mockUserHasAccessTrue

		paymentAmount := uint(1500)

		mockDb.On("GetWorkspaceBudget", "org-1").Return(db.NewBountyBudget{
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
		bHandler.userHasAccess = mockUserHasAccessTrue

		mockDb.On("GetWorkspaceBudget", "org-1").Return(db.NewBountyBudget{
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

	t.Run("Should test that an Workspace's Budget Total Amount is accurate after three (3) successful 'Budget Withdrawal Requests'", func(t *testing.T) {
		ctxs := context.WithValue(context.Background(), auth.ContextKey, "valid-key")
		mockDb := dbMocks.NewDatabase(t)
		mockHttpClient := mocks.NewHttpClient(t)
		bHandler := NewBountyHandler(mockHttpClient, mockDb)
		bHandler.userHasAccess = mockUserHasAccessTrue

		paymentAmount := uint(1500)
		initialBudget := uint(5000)
		invoice := "lnbc15u1p3xnhl2pp5jptserfk3zk4qy42tlucycrfwxhydvlemu9pqr93tuzlv9cc7g3sdqsvfhkcap3xyhx7un8cqzpgxqzjcsp5f8c52y2stc300gl6s4xswtjpc37hrnnr3c9wvtgjfuvqmpm35evq9qyyssqy4lgd8tj637qcjp05rdpxxykjenthxftej7a2zzmwrmrl70fyj9hvj0rewhzj7jfyuwkwcg9g2jpwtk3wkjtwnkdks84hsnu8xps5vsq4gj5hs"

		for i := 0; i < 3; i++ {
			expectedFinalBudget := initialBudget - (paymentAmount * uint(i))

			mockDb.ExpectedCalls = nil
			mockDb.Calls = nil
			mockHttpClient.ExpectedCalls = nil
			mockHttpClient.Calls = nil

			mockDb.On("GetWorkspaceBudget", "org-1").Return(db.NewBountyBudget{
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
			finalBudget := mockDb.GetWorkspaceBudget("org-1")
			assert.Equal(t, expectedFinalBudget, finalBudget.TotalBudget, "The workspace's final budget should reflect the deductions from the successful withdrawals")

		}
	})

	t.Run("Should test that the BountyBudgetWithdraw handler gets locked by go mutex when it is called i.e. the handler has to be fully executed before it processes another request.", func(t *testing.T) {
		mockDb := dbMocks.NewDatabase(t)
		mockHttpClient := mocks.NewHttpClient(t)
		bHandler := NewBountyHandler(mockHttpClient, mockDb)

		var processingTimes []time.Time
		var mutex sync.Mutex

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mutex.Lock()
			processingTimes = append(processingTimes, time.Now())
			time.Sleep(10 * time.Millisecond)
			mutex.Unlock()

			bHandler.BountyBudgetWithdraw(w, r)
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

}

func TestPollInvoice(t *testing.T) {
	ctx := context.Background()

	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mockHttpClient := &mocks.HttpClient{}
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	paymentRequest := "DummyPayReq"

	now := time.Now()
	bountyAmount := uint(5000)
	invoice := db.NewInvoiceList{
		PaymentRequest: paymentRequest,
		Status:         false,
		Type:           "KEYSEND",
		OwnerPubkey:    "owner_pubkey",
		WorkspaceUuid:  "workspace_uuid",
		Created:        &now,
	}
	db.TestDB.AddInvoice(invoice)

	invoiceData := db.UserInvoiceData{
		PaymentRequest: invoice.PaymentRequest,
		Amount:         bountyAmount,
		UserPubkey:     invoice.OwnerPubkey,
		Created:        int(now.Unix()),
	}
	db.TestDB.AddUserInvoiceData(invoiceData)

	bounty := db.NewBounty{
		OwnerID:     "owner_pubkey",
		Price:       bountyAmount,
		Created:     now.Unix(),
		Type:        "coding",
		Title:       "bountyTitle",
		Description: "bountyDescription",
		Assignee:    "assigneePubkey",
		Show:        true,
		Paid:        false,
	}
	db.TestDB.CreateOrEditBounty(bounty)

	unauthorizedCtx := context.WithValue(ctx, auth.ContextKey, "")
	authorizedCtx := context.WithValue(ctx, auth.ContextKey, invoice.OwnerPubkey)

	t.Run("Should test that a 401 error is returned if a user is unauthorized", func(t *testing.T) {
		r := chi.NewRouter()
		r.Post("/poll/invoice/{paymentRequest}", bHandler.PollInvoice)

		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(unauthorizedCtx, http.MethodPost, "/poll/invoice/"+invoice.PaymentRequest, bytes.NewBufferString(`{}`))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected 401 error if a user is unauthorized")
	})

	t.Run("Should test that a 403 error is returned if there is an invoice error", func(t *testing.T) {
		expectedUrl := fmt.Sprintf("%s/invoice?payment_request=%s", config.RelayUrl, invoice.PaymentRequest)

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
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/poll/invoice/"+invoice.PaymentRequest, bytes.NewBufferString(`{}`))
		if err != nil {
			t.Fatal(err)
		}

		ro.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code, "Expected 403 error if there is an invoice error")
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("Should mock relay payment is successful update the bounty associated with the invoice and set the paid as true", func(t *testing.T) {
		expectedUrl := fmt.Sprintf("%s/invoice?payment_request=%s", config.RelayUrl, invoice.PaymentRequest)
		expectedBody := fmt.Sprintf(`{"success": true, "response": { "settled": true, "payment_request": "%s", "payment_hash": "payment_hash", "preimage": "preimage", "Amount": %d}}`, invoice.OwnerPubkey, bountyAmount)

		r := io.NopCloser(bytes.NewReader([]byte(expectedBody)))
		mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.Method == http.MethodGet && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey
		})).Return(&http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil).Once()

		expectedPaymentUrl := fmt.Sprintf("%s/payment", config.RelayUrl)
		expectedPaymentBody := fmt.Sprintf(`{"amount": %d, "destination_key": "%s", "text": "memotext added for notification"}`, bountyAmount, invoice.OwnerPubkey)

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
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/poll/invoice/"+invoice.PaymentRequest, bytes.NewBufferString(`{}`))
		if err != nil {
			t.Fatal(err)
		}

		ro.ServeHTTP(rr, req)

		invData := db.TestDB.GetUserInvoiceData(invoice.PaymentRequest)
		updatedBounty, err := db.TestDB.GetBountyByCreated(uint(invData.Created))
		if err != nil {
			t.Fatal(err)
		}
		updatedInvoice := db.TestDB.GetInvoice(invoice.PaymentRequest)

		assert.True(t, updatedBounty.Paid, "Expected bounty to be marked as paid")
		assert.True(t, updatedInvoice.Status, "Expected invoice status to be true")
		assert.Equal(t, http.StatusOK, rr.Code)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("If the invoice is settled and the invoice.Type is equal to BUDGET the invoice amount should be added to the workspace budget and the payment status of the related invoice should be sent to true on the payment history table", func(t *testing.T) {
		db.TestDB.DeleteInvoice(paymentRequest)

		invoice := db.NewInvoiceList{
			PaymentRequest: paymentRequest,
			Status:         false,
			OwnerPubkey:    "owner_pubkey",
			WorkspaceUuid:  "workspace_uuid",
			Created:        &now,
		}

		db.TestDB.AddInvoice(invoice)

		ctx := context.Background()
		mockHttpClient := &mocks.HttpClient{}
		bHandler := NewBountyHandler(mockHttpClient, db.TestDB)
		authorizedCtx := context.WithValue(ctx, auth.ContextKey, invoice.OwnerPubkey)
		expectedUrl := fmt.Sprintf("%s/invoice?payment_request=%s", config.RelayUrl, invoice.PaymentRequest)
		expectedBody := fmt.Sprintf(`{"success": true, "response": { "settled": true, "payment_request": "%s", "payment_hash": "payment_hash", "preimage": "preimage", "Amount": %d}}`, invoice.OwnerPubkey, bountyAmount)

		r := io.NopCloser(bytes.NewReader([]byte(expectedBody)))
		mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.Method == http.MethodGet && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey
		})).Return(&http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil).Once()

		ro := chi.NewRouter()
		ro.Post("/poll/invoice/{paymentRequest}", bHandler.PollInvoice)

		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/poll/invoice/"+invoice.PaymentRequest, bytes.NewBufferString(`{}`))
		if err != nil {
			t.Fatal(err)
		}

		ro.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockHttpClient.AssertExpectations(t)
	})
}
