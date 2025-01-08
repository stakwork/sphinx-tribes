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
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stakwork/sphinx-tribes/logger"
	"github.com/stakwork/sphinx-tribes/utils"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
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
		defer db.CloseTestDB()
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

	// create user
	db.TestDB.CreateOrEditPerson(bountyOwner)

	existingBounty := db.NewBounty{
		Type:          "coding",
		Title:         "existing bounty",
		Description:   "existing bounty description",
		WorkspaceUuid: "work-1",
		OwnerID:       bountyOwner.OwnerPubKey,
		Price:         2000,
	}

	// Add initial Bounty
	AddExisitingDB(existingBounty)

	newBounty := db.NewBounty{
		Type:          "coding",
		Title:         "new bounty",
		Description:   "new bounty description",
		WorkspaceUuid: "work-1",
		OwnerID:       bountyOwner.OwnerPubKey,
		Price:         1500,
	}

	failedBounty := db.NewBounty{
		Title:         "new bounty",
		Description:   "failed bounty description",
		WorkspaceUuid: "work-1",
		OwnerID:       bountyOwner.OwnerPubKey,
		Price:         1500,
	}

	ctx := context.WithValue(context.Background(), auth.ContextKey, bountyOwner.OwnerPubKey)
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
			logger.Log.Error("Could not marshal json data")
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
	botURL := os.Getenv("V2_BOT_URL")
	botToken := os.Getenv("V2_BOT_TOKEN")

	expectedUrl := fmt.Sprintf("%s/invoices", config.RelayUrl)
	expectedBody := `{"payment_request": "req-id"}`

	expectedV2Url := fmt.Sprintf("%s/pay_invoice", botURL)
	expectedV2Body := `{"bolt11": "req-id", "wait": true}`

	t.Run("validate request url, body and headers", func(t *testing.T) {
		mockHttpClient := &mocks.HttpClient{}
		mockDb := &dbMocks.Database{}
		handler := NewBountyHandler(mockHttpClient, mockDb)

		if botURL != "" && botToken != "" {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPost && expectedV2Url == req.URL.String() && req.Header.Get("x-admin-token") == botToken && expectedV2Body == string(bodyByt)
			})).Return(nil, errors.New("some-error")).Once()
		} else {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPut && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
			})).Return(nil, errors.New("some-error")).Once()
		}

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

		if botURL != "" && botToken != "" {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPost && expectedV2Url == req.URL.String() && req.Header.Get("x-admin-token") == botToken && expectedV2Body == string(bodyByt)
			})).Return(&http.Response{
				StatusCode: 500,
				Body:       r,
			}, nil)
		} else {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPut && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
			})).Return(&http.Response{
				StatusCode: 500,
				Body:       r,
			}, nil)
		}

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

		if botURL != "" && botToken != "" {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPost && expectedV2Url == req.URL.String() && req.Header.Get("x-admin-token") == botToken && expectedV2Body == string(bodyByt)
			})).Return(&http.Response{
				StatusCode: 500,
				Body:       r,
			}, nil)
		} else {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPut && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
			})).Return(&http.Response{
				StatusCode: 500,
				Body:       r,
			}, nil).Once()
		}
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

		if botURL != "" && botToken != "" {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPost && expectedV2Url == req.URL.String() && req.Header.Get("x-admin-token") == botToken && expectedV2Body == string(bodyByt)
			})).Return(&http.Response{
				StatusCode: 500,
				Body:       r,
			}, nil)
		} else {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPut && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
			})).Return(&http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil).Once()
		}

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

		rv3 := io.NopCloser(bytes.NewReader([]byte(`{"status": "COMPLETE", "amt_msat": "1000", "timestamp": "" }`)))

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

		expectedV2SuccessMsg := db.InvoicePaySuccess{
			Success: true,
			Response: db.InvoiceCheckResponse{
				Settled:         true,
				Payment_request: "req-id",
				Payment_hash:    "",
				Preimage:        "",
				Amount:          "",
			},
		}

		if botURL != "" && botToken != "" {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPost && expectedV2Url == req.URL.String() && req.Header.Get("x-admin-token") == botToken && expectedV2Body == string(bodyByt)
			})).Return(&http.Response{
				StatusCode: 200,
				Body:       rv3,
			}, nil)

			success, invoicePayErr := handler.PayLightningInvoice("req-id")

			assert.Empty(t, invoicePayErr)
			assert.EqualValues(t, expectedV2SuccessMsg, success)
			mockHttpClient.AssertExpectations(t)
		} else {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPut && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
			})).Return(&http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil).Once()

			success, invoicePayErr := handler.PayLightningInvoice("req")

			assert.Empty(t, invoicePayErr)
			assert.EqualValues(t, expectedSuccessMsg, success)
			mockHttpClient.AssertExpectations(t)
		}
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
		mockDb.On("GetProofsByBountyID", bounty.ID).Return([]db.ProofOfWork{}).Once()
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

		now := time.Now().UnixMilli()
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
			logger.Log.Error("upgrade error: %v", err)
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

	now := time.Now().UnixMilli()
	bountyOwnerId := "owner_pubkey"

	botURL := os.Getenv("V2_BOT_URL")
	botToken := os.Getenv("V2_BOT_TOKEN")

	person := db.Person{
		Uuid:           "uuid",
		OwnerAlias:     "alias",
		UniqueName:     "unique_name",
		OwnerPubKey:    "03b2205df68d90f8f9913650bc3161761b61d743e615a9faa7ffecea3380a93fc1",
		OwnerRouteHint: "02162c52716637fb8120ab0261e410b185d268d768cc6f6227c58102d194ad0bc2_1099607703554",
		PriceToMeet:    0,
		Description:    "description",
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
		mockHttpClient := &mocks.HttpClient{}

		bHandler2 := NewBountyHandler(mockHttpClient, db.TestDB)
		bHandler2.getSocketConnections = mockGetSocketConnections
		bHandler2.userHasAccess = mockUserHasAccessTrue

		memoData := fmt.Sprintf("Payment For: %ss", bounty.Title)
		memoText := url.QueryEscape(memoData)

		expectedUrl := fmt.Sprintf("%s/payment", config.RelayUrl)
		expectedBody := fmt.Sprintf(`{"amount": %d, "destination_key": "%s", "text": "memotext added for notification", "data": "%s"}`, bountyAmount, person.OwnerPubKey, memoText)

		expectedV2Url := fmt.Sprintf("%s/pay", botURL)
		expectedV2Body :=
			fmt.Sprintf(`{"amt_msat": %d, "dest": "%s", "route_hint": "%s", "data": "%s", "wait": true}`, bountyAmount*1000, person.OwnerPubKey, person.OwnerRouteHint, memoText)

		r := io.NopCloser(bytes.NewReader([]byte(`"internal server error"`)))
		if botURL != "" && botToken != "" {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPost && expectedV2Url == req.URL.String() && req.Header.Get("x-admin-token") == botToken && expectedV2Body == string(bodyByt)
			})).Return(&http.Response{
				StatusCode: 406,
				Body:       r,
			}, nil).Once()
		} else {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPost && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
			})).Return(&http.Response{
				StatusCode: 500,
				Body:       r,
			}, nil).Once()
		}

		ro := chi.NewRouter()
		ro.Post("/gobounties/pay/{id}", bHandler2.MakeBountyPayment)

		requestBody := bytes.NewBuffer([]byte("{}"))
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/gobounties/pay/"+bountyIdStr, requestBody)
		if err != nil {
			t.Fatal(err)
		}

		ro.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("Should test that a successful WebSocket message is sent if the payment is successful", func(t *testing.T) {

		bHandler.getSocketConnections = mockGetSocketConnections
		bHandler.userHasAccess = mockUserHasAccessTrue

		memoData := fmt.Sprintf("Payment For: %ss", bounty.Title)
		memoText := url.QueryEscape(memoData)

		expectedUrl := fmt.Sprintf("%s/payment", config.RelayUrl)
		expectedBody := fmt.Sprintf(`{"amount": %d, "destination_key": "%s", "text": "memotext added for notification", "data": "%s"}`, bountyAmount, person.OwnerPubKey, memoText)

		expectedV2Url := fmt.Sprintf("%s/pay", botURL)
		expectedV2Body :=
			fmt.Sprintf(`{"amt_msat": %d, "dest": "%s", "route_hint": "%s", "data": "%s", "wait": true}`, bountyAmount*1000, person.OwnerPubKey, person.OwnerRouteHint, memoText)

		if botURL != "" && botToken != "" {
			rv2 := io.NopCloser(bytes.NewReader([]byte(`{"status": "COMPLETE", "tag": "", "preimage": "", "payment_hash": "" }`)))
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPost && expectedV2Url == req.URL.String() && req.Header.Get("x-admin-token") == botToken && expectedV2Body == string(bodyByt)
			})).Return(&http.Response{
				StatusCode: 200,
				Body:       rv2,
			}, nil).Once()
		} else {
			r := io.NopCloser(bytes.NewReader([]byte(`{"success": true, "response": { "sumAmount": "1"}}`)))
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				bodyByt, _ := io.ReadAll(req.Body)
				return req.Method == http.MethodPost && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey && expectedBody == string(bodyByt)
			})).Return(&http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil).Once()
		}

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

		updatedBounty := db.TestDB.GetBounty(bountyId)
		assert.True(t, updatedBounty.Paid, "Expected bounty to be marked as paid")

		updatedWorkspaceBudget := db.TestDB.GetWorkspaceBudget(bounty.WorkspaceUuid)
		assert.Equal(t, budgetAmount-bountyAmount, updatedWorkspaceBudget.TotalBudget, "Expected workspace budget to be reduced by bounty amount")
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

func TestUpdateBountyPaymentStatus(t *testing.T) {
	ctx := context.Background()

	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mockHttpClient := &mocks.HttpClient{}

	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	paymentTag := "update_tag"

	mockPendingGetInvoiceStatusByTag := func(tag string) db.V2TagRes {
		return db.V2TagRes{
			Status: db.PaymentPending,
			Tag:    paymentTag,
			Error:  "",
		}

	}
	mockCompleteGetInvoiceStatusByTag := func(tag string) db.V2TagRes {
		return db.V2TagRes{
			Status: db.PaymentComplete,
			Tag:    paymentTag,
			Error:  "",
		}
	}

	now := time.Now().UnixMilli()
	bountyOwnerId := "owner_pubkey"

	person := db.Person{
		Uuid:           "update_payment_uuid",
		OwnerAlias:     "update_alias",
		UniqueName:     "update_unique_name",
		OwnerPubKey:    "03b2205df68d90f8f9913650bc3161761b61d743e615a9faa7ffecea3380a99fg1",
		OwnerRouteHint: "02162c52716637fb8120ab0261e410b185d268d768cc6f6227c58102d194ad0bc2_1088607703554",
		PriceToMeet:    0,
		Description:    "update_description",
	}

	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        "update_workspace_uuid",
		Name:        "update_workspace_name",
		OwnerPubKey: person.OwnerPubKey,
		Github:      "gtihub",
		Website:     "website",
		Description: "update_description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	budgetAmount := uint(10000)
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
		Title:         "updateBountyTitle",
		Description:   "updateBountyDescription",
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

	paymentTime := time.Now()

	payment := db.NewPaymentHistory{
		BountyId:       bountyId,
		PaymentStatus:  db.PaymentPending,
		WorkspaceUuid:  workspace.Uuid,
		PaymentType:    db.Payment,
		SenderPubKey:   person.OwnerPubKey,
		ReceiverPubKey: person.OwnerPubKey,
		Tag:            paymentTag,
		Status:         true,
		Created:        &paymentTime,
		Updated:        &paymentTime,
	}

	db.TestDB.AddPaymentHistory(payment)

	unauthorizedCtx := context.WithValue(ctx, auth.ContextKey, "")
	authorizedCtx := context.WithValue(ctx, auth.ContextKey, person.OwnerPubKey)

	t.Run("401 unauthorized error when unauthorized user hits endpoint", func(t *testing.T) {

		r := chi.NewRouter()
		r.Post("/gobounties/payment/status/{id}", bHandler.UpdateBountyPaymentStatus)

		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(unauthorizedCtx, http.MethodPost, "/gobounties/payment/status/"+bountyIdStr, nil)

		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected 401 Unauthorized for unauthorized access")
	})

	t.Run("Should test that a PENDING payment_status is sent if the payment is not successful", func(t *testing.T) {
		mockHttpClient := &mocks.HttpClient{}

		bHandler := NewBountyHandler(mockHttpClient, db.TestDB)
		bHandler.getInvoiceStatusByTag = mockPendingGetInvoiceStatusByTag

		ro := chi.NewRouter()
		ro.Put("/gobounties/payment/status/{id}", bHandler.UpdateBountyPaymentStatus)

		rr := httptest.NewRecorder()
		requestBody := bytes.NewBuffer([]byte("{}"))
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPut, "/gobounties/payment/status/"+bountyIdStr, requestBody)
		if err != nil {
			t.Fatal(err)
		}

		ro.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("Should test that a COMPLETE payment_status is sent if the payment is successful", func(t *testing.T) {
		mockHttpClient := &mocks.HttpClient{}

		bHandler := NewBountyHandler(mockHttpClient, db.TestDB)
		bHandler.getInvoiceStatusByTag = mockCompleteGetInvoiceStatusByTag

		ro := chi.NewRouter()
		ro.Put("/gobounties/payment/status/{id}", bHandler.UpdateBountyPaymentStatus)

		requestBody := bytes.NewBuffer([]byte("{}"))
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPut, "/gobounties/payment/status/"+bountyIdStr, requestBody)
		if err != nil {
			t.Fatal(err)
		}

		ro.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockHttpClient.AssertExpectations(t)

		payment := db.TestDB.GetPaymentByBountyId(payment.BountyId)

		updatedBounty := db.TestDB.GetBounty(bountyId)
		assert.True(t, updatedBounty.Paid, "Expected bounty to be marked as paid")
		assert.Equal(t, payment.PaymentStatus, db.PaymentComplete, "Expected Payment Status To be Complete")
	})

	t.Run("405 when trying to update an already-paid bounty", func(t *testing.T) {
		r := chi.NewRouter()
		r.Put("/gobounties/payment/status/{id}", bHandler.UpdateBountyPaymentStatus)

		requestBody := bytes.NewBuffer([]byte("{}"))
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPut, "/gobounties/payment/status/"+bountyIdStr, requestBody)
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Expected 405 Method Not Allowed for an already-paid bounty")
	})
}

func TestBountyBudgetWithdraw(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	ctx := context.Background()
	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	handlerUserHasAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
		return true
	}

	handlerUserNotAccess := func(pubKeyFromAuth string, uuid string, role string) bool {
		return false
	}

	getHoursDifference := func(createdDate int64, endDate *time.Time) int64 {
		return 2
	}

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

	budgetAmount := uint(5000)

	paymentTime := time.Now()

	payment := db.NewPaymentHistory{
		Amount:         budgetAmount,
		WorkspaceUuid:  workspace.Uuid,
		PaymentType:    db.Deposit,
		SenderPubKey:   person.OwnerPubKey,
		ReceiverPubKey: person.OwnerPubKey,
		Tag:            "test_deposit",
		Status:         true,
		Created:        &paymentTime,
		Updated:        &paymentTime,
	}

	db.TestDB.AddPaymentHistory(payment)

	budget := db.NewBountyBudget{
		WorkspaceUuid: workspace.Uuid,
		TotalBudget:   budgetAmount,
	}
	db.TestDB.CreateWorkspaceBudget(budget)

	unauthorizedCtx := context.WithValue(context.Background(), auth.ContextKey, "")
	authorizedCtx := context.WithValue(ctx, auth.ContextKey, person.OwnerPubKey)

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
		bHandler.userHasAccess = handlerUserNotAccess

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.BountyBudgetWithdraw)

		validData := []byte(`{"workspace_uuid": "workspace-uuid", "paymentRequest": "invoice"}`)
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/budget/withdraw", bytes.NewReader(validData))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "You don't have appropriate permissions to withdraw bounty budget")
	})

	t.Run("403 error when amount exceeds workspace's budget", func(t *testing.T) {

		bHandler.userHasAccess = handlerUserHasAccess

		invoice := "lnbc100u1png0l8ypp5hna5vnd2hcskpf69rt5y9dly2p202lejcacj53md32wx87vc2mnqdqzvscqzpgxqyz5vqrzjqwnw5tv745sjpvft6e3f9w62xqk826vrm3zaev4nvj6xr3n065aukqqqqyqqpmgqqyqqqqqqqqqqqqqqqqsp5cdg0c2qhuewz4j8680pf5va0l9a382qa5sakg4uga4nv4wnuf5qs9qrssqpdddmqtflxz3553gm5xq8ptdpl2t3ew49hgjnta0v0eyz747drkkhmnk5yxg676kvmgyugm35cts9dmrnt9mcgejg64kwk9nwxqg43cqcvxm44"

		amount := utils.GetInvoiceAmount(invoice)
		assert.Equal(t, uint(10000), amount)

		withdrawRequest := db.NewWithdrawBudgetRequest{
			PaymentRequest: invoice,
			WorkspaceUuid:  workspace.Uuid,
		}
		requestBody, _ := json.Marshal(withdrawRequest)
		req, _ := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/budget/withdraw", bytes.NewReader(requestBody))

		rr := httptest.NewRecorder()

		bHandler.BountyBudgetWithdraw(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code, "Expected 403 Forbidden when the payment exceeds the workspace's budget")
		assert.Contains(t, rr.Body.String(), "Workspace budget is not enough to withdraw the amount", "Expected specific error message")
	})

	t.Run("budget invoices get paid if amount is lesser than workspace's budget", func(t *testing.T) {
		mockHttpClient := mocks.NewHttpClient(t)
		bHandler := NewBountyHandler(mockHttpClient, db.TestDB)
		bHandler.userHasAccess = handlerUserHasAccess

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.BountyBudgetWithdraw)
		paymentAmount := uint(300)
		initialBudget := budget.TotalBudget
		expectedFinalBudget := initialBudget - paymentAmount
		budget.TotalBudget = expectedFinalBudget

		mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"status": "COMPLETE", "amt_msat": "1000", "timestamp": "" }`)),
		}, nil)

		invoice := "lnbc3u1pngsqv8pp5vl6ep8llmg3f9sfu8j7ctcnphylpnjduuyljqf3sc30z6ejmrunqdqzvscqzpgxqyz5vqrzjqwnw5tv745sjpvft6e3f9w62xqk826vrm3zaev4nvj6xr3n065aukqqqqyqqz9gqqyqqqqqqqqqqqqqqqqsp5n9hrrw6pr89qn3c82vvhy697wp45zdsyhm7tnu536ga77ytvxxaq9qrssqqqhenjtquz8wz5tym8v830h9gjezynjsazystzj6muhw4rd9ccc40p8sazjuk77hhcj0xn72lfyee3tsfl7lucxkx5xgtfaqya9qldcqr3072z"

		withdrawRequest := db.NewWithdrawBudgetRequest{
			PaymentRequest: invoice,
			WorkspaceUuid:  workspace.Uuid,
		}

		requestBody, _ := json.Marshal(withdrawRequest)
		req, _ := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/budget/withdraw", bytes.NewReader(requestBody))

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response db.InvoicePaySuccess
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success, "Expected invoice payment to succeed")

		finalBudget := db.TestDB.GetWorkspaceBudget(workspace.Uuid)
		assert.Equal(t, expectedFinalBudget, finalBudget.TotalBudget, "The workspace's final budget should reflect the deductions from the successful withdrawals")
	})

	t.Run("400 BadRequest error if there is an error with invoice payment", func(t *testing.T) {
		bHandler.getHoursDifference = getHoursDifference

		mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewBufferString(`{"success": false, "error": "Payment error"}`)),
		}, nil)

		invoice := "lnbcrt1u1pnv5ejzdqad9h8vmmfvdjjqen0wgsrzvpsxqcrqpp58xyhvymlhc8q05z930fknk2vdl8wnpm5zlx5lgp4ev9u8h7yd4kssp5nu652c5y0epuxeawn8szcgdrjxwk7pfkdh9tsu44r7hacg52nfgq9qrsgqcqpjxqrrssrzjqgtzc5n3vcmlhqfq4vpxreqskxzay6xhdrxx7c38ckqs95v5459uyqqqqyqq9ggqqsqqqqqqqqqqqqqq9gwyffzjpnrwt6yswwd4znt2xqnwjwxgq63qxudru95a8pqeer2r7sduurtstz5x60y4e7m4y9nx6rqy5sr9k08vtwv6s37xh0z5pdwpgqxeqdtv"

		withdrawRequest := db.NewWithdrawBudgetRequest{
			PaymentRequest: invoice,
			WorkspaceUuid:  workspace.Uuid,
		}
		requestBody, _ := json.Marshal(withdrawRequest)
		req, _ := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/budget/withdraw", bytes.NewReader(requestBody))

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
		paymentAmount := uint(1000)
		initialBudget := budget.TotalBudget
		invoice := "lnbcrt10u1pnv7nz6dqld9h8vmmfvdjjqen0wgsrzvpsxqcrqvqpp54v0synj4q3j2usthzt8g5umteky6d2apvgtaxd7wkepkygxgqdyssp5lhv2878qjas3azv3nnu8r6g3tlgejl7mu7cjzc9q5haygrpapd4s9qrsgqcqpjxqrrssrzjqgtzc5n3vcmlhqfq4vpxreqskxzay6xhdrxx7c38ckqs95v5459uyqqqqyqqtwsqqgqqqqqqqqqqqqqq9gea2fjj7q302ncprk2pawk4zdtayycvm0wtjpprml96h9vujvmqdp0n5z8v7lqk44mq9620jszwaevj0mws7rwd2cegxvlmfszwgpgfqp2xafjf"

		bHandler.userHasAccess = handlerUserHasAccess
		bHandler.getHoursDifference = getHoursDifference

		for i := 0; i < 3; i++ {
			expectedFinalBudget := initialBudget - (paymentAmount * uint(i+1))
			mockHttpClient.ExpectedCalls = nil
			mockHttpClient.Calls = nil

			// add a zero amount withdrawal with a time lesser than 2 + loop index hours to beat the 1 hour withdrawal timer
			dur := int(time.Hour.Hours())*2 + i + 1
			paymentTime = time.Now().Add(-time.Hour * time.Duration(dur))

			mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`{"status": "COMPLETE", "amt_msat": "1000", "timestamp": "" }`)),
			}, nil)

			withdrawRequest := db.NewWithdrawBudgetRequest{
				PaymentRequest: invoice,
				WorkspaceUuid:  workspace.Uuid,
			}
			requestBody, _ := json.Marshal(withdrawRequest)
			req, _ := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/budget/withdraw", bytes.NewReader(requestBody))

			rr := httptest.NewRecorder()

			bHandler.BountyBudgetWithdraw(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code)
			var response db.InvoicePaySuccess
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.True(t, response.Success, "Expected invoice payment to succeed")

			finalBudget := db.TestDB.GetWorkspaceBudget(workspace.Uuid)
			assert.Equal(t, expectedFinalBudget, finalBudget.TotalBudget, "The workspace's final budget should reflect the deductions from the successful withdrawals")

		}
	})

	t.Run("Should test that the BountyBudgetWithdraw handler gets locked by go mutex when it is called i.e. the handler has to be fully executed before it processes another request.", func(t *testing.T) {

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

	paymentRequest := "lnbcrt10u1pnv7nz6dqld9h8vmmfvdjjqen0wgsrzvpsxqcrqvqpp54v0synj4q3j2usthzt8g5umteky6d2apvgtaxd7wkepkygxgqdyssp5lhv2878qjas3azv3nnu8r6g3tlgejl7mu7cjzc9q5haygrpapd4s9qrsgqcqpjxqrrssrzjqgtzc5n3vcmlhqfq4vpxreqskxzay6xhdrxx7c38ckqs95v5459uyqqqqyqqtwsqqgqqqqqqqqqqqqqq9gea2fjj7q302ncprk2pawk4zdtayycvm0wtjpprml96h9vujvmqdp0n5z8v7lqk44mq9620jszwaevj0mws7rwd2cegxvlmfszwgpgfqp2xafj"

	botURL := os.Getenv("V2_BOT_URL")
	botToken := os.Getenv("V2_BOT_TOKEN")

	now := time.Now()
	bountyAmount := uint(5000)
	invoice := db.NewInvoiceList{
		PaymentRequest: paymentRequest,
		Status:         false,
		Type:           "KEYSEND",
		OwnerPubkey:    "03b2205df68d90f8f9913650bc3161761b61d743e615a9faa7ffecea3380a93fc1",
		WorkspaceUuid:  "workspace_uuid",
		Created:        &now,
	}
	db.TestDB.AddInvoice(invoice)

	invoiceData := db.UserInvoiceData{
		PaymentRequest: invoice.PaymentRequest,
		Amount:         bountyAmount,
		UserPubkey:     invoice.OwnerPubkey,
		Created:        int(now.Unix()),
		RouteHint:      "02162c52716637fb8120ab0261e410b185d268d768cc6f6227c58102d194ad0bc2_1099607703554",
	}
	db.TestDB.AddUserInvoiceData(invoiceData)

	bounty := db.NewBounty{
		OwnerID:     "owner_pubkey",
		Price:       bountyAmount,
		Created:     now.Unix(),
		Type:        "coding",
		Title:       "bountyTitle",
		Description: "bountyDescription",
		Assignee:    "03b2205df68d90f8f9913650bc3161761b61d743e615a9faa7ffecea3380a93fc1",
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

		expectedV2Url := fmt.Sprintf("%s/check_invoice", botURL)

		r := io.NopCloser(bytes.NewReader([]byte(`{"success": false, "error": "Internel server error"}`)))

		if botURL != "" && botToken != "" {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				return req.Method == http.MethodPost && expectedV2Url == req.URL.String() && req.Header.Get("x-admin-token") == botToken
			})).Return(&http.Response{
				StatusCode: 500,
				Body:       r,
			}, nil).Once()
		} else {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				return req.Method == http.MethodGet && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey
			})).Return(&http.Response{
				StatusCode: 500,
				Body:       r,
			}, nil).Once()
		}

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

		expectedV2Url := fmt.Sprintf("%s/check_invoice", botURL)
		expectedV2InvoiceBody := `{"status": "paid", "amt_msat": "", "timestamp": ""}`

		r := io.NopCloser(bytes.NewReader([]byte(expectedBody)))
		rv2 := io.NopCloser(bytes.NewReader([]byte(expectedV2InvoiceBody)))

		if botURL != "" && botToken != "" {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				return req.Method == http.MethodPost && expectedV2Url == req.URL.String() && req.Header.Get("x-admin-token") == botToken
			})).Return(&http.Response{
				StatusCode: 200,
				Body:       rv2,
			}, nil).Once()
		} else {
			mockHttpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				return req.Method == http.MethodGet && expectedUrl == req.URL.String() && req.Header.Get("x-user-token") == config.RelayAuthKey
			})).Return(&http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil).Once()
		}

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

func TestGetBountyCards(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	db.CleanTestData()

	workspace := db.Workspace{
		ID:          1,
		Uuid:        "test-workspace-uuid",
		Name:        "Test Workspace",
		Description: "Test Workspace Description",
		OwnerPubKey: "test-owner",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	phase := db.FeaturePhase{
		Uuid:        "test-phase-uuid",
		Name:        "Test Phase",
		FeatureUuid: "test-feature-uuid",
	}
	db.TestDB.CreateOrEditFeaturePhase(phase)

	feature := db.WorkspaceFeatures{
		Uuid:          "test-feature-uuid",
		Name:          "Test Feature",
		WorkspaceUuid: workspace.Uuid,
	}
	db.TestDB.CreateOrEditFeature(feature)

	assignee := db.Person{
		OwnerPubKey: "test-assignee",
		Img:         "test-image-url",
	}
	db.TestDB.CreateOrEditPerson(assignee)

	now := time.Now()
	bounty := db.NewBounty{
		ID:            1,
		Type:          "coding",
		Title:         "Test Bounty",
		Description:   "Test Description",
		WorkspaceUuid: workspace.Uuid,
		PhaseUuid:     phase.Uuid,
		Assignee:      assignee.OwnerPubKey,
		Show:          true,
		Created:       now.Unix(),
		OwnerID:       "test-owner",
		Price:         1000,
		Paid:          false,
	}
	db.TestDB.CreateOrEditBounty(bounty)

	t.Run("should successfully return bounty cards", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyCards)

		req, err := http.NewRequest(http.MethodGet, "/gobounties/bounty-cards?workspace_uuid="+workspace.Uuid, nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var response []db.BountyCard
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotEmpty(t, response, "Response should not be empty")

		firstCard := response[0]
		assert.Equal(t, bounty.ID, firstCard.BountyID)
		assert.Equal(t, bounty.Title, firstCard.Title)
		assert.Equal(t, assignee.Img, firstCard.AssigneePic)

		assert.Equal(t, feature.Uuid, firstCard.Features.Uuid)
		assert.Equal(t, feature.Name, firstCard.Features.Name)
		assert.Equal(t, feature.WorkspaceUuid, firstCard.Features.WorkspaceUuid)

		assert.Equal(t, phase.Uuid, firstCard.Phase.Uuid)
		assert.Equal(t, phase.Name, firstCard.Phase.Name)
		assert.Equal(t, phase.FeatureUuid, firstCard.Phase.FeatureUuid)

		assert.Equal(t, workspace, firstCard.Workspace)
	})

	t.Run("should return empty array when no bounties exist", func(t *testing.T) {

		db.TestDB.DeleteAllBounties()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyCards)

		req, err := http.NewRequest(http.MethodGet, "/gobounties/bounty-cards", nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var response []db.BountyCard
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, response)
	})

	t.Run("should handle bounties without phase and feature", func(t *testing.T) {
		bountyWithoutPhase := db.NewBounty{
			ID:            2,
			Type:          "coding",
			Title:         "Test Bounty Without Phase",
			Description:   "Test Description",
			WorkspaceUuid: workspace.Uuid,
			Assignee:      assignee.OwnerPubKey,
			Show:          true,
			Created:       now.Unix(),
			OwnerID:       "test-owner",
			Price:         1000,
			Paid:          false,
		}
		db.TestDB.CreateOrEditBounty(bountyWithoutPhase)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyCards)

		req, err := http.NewRequest(http.MethodGet, "/gobounties/bounty-cards?workspace_uuid="+workspace.Uuid, nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var response []db.BountyCard
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotEmpty(t, response)

		var cardWithoutPhase db.BountyCard
		for _, card := range response {
			if card.BountyID == bountyWithoutPhase.ID {
				cardWithoutPhase = card
				break
			}
		}

		assert.Equal(t, bountyWithoutPhase.ID, cardWithoutPhase.BountyID)
		assert.Equal(t, bountyWithoutPhase.Title, cardWithoutPhase.Title)
		assert.Equal(t, assignee.Img, cardWithoutPhase.AssigneePic)
		assert.Empty(t, cardWithoutPhase.Phase.Uuid)
		assert.Empty(t, cardWithoutPhase.Features.Uuid)
	})

	t.Run("should handle bounties without assignee", func(t *testing.T) {

		db.TestDB.DeleteAllBounties()

		bountyWithoutAssignee := db.NewBounty{
			ID:            1,
			Type:          "coding",
			Title:         "Test Bounty Without Assignee",
			Description:   "Test Description",
			WorkspaceUuid: workspace.Uuid,
			PhaseUuid:     phase.Uuid,
			Show:          true,
			Created:       now.Unix(),
			OwnerID:       "test-owner",
			Price:         1000,
			Paid:          false,
		}
		db.TestDB.CreateOrEditBounty(bountyWithoutAssignee)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyCards)

		req, err := http.NewRequest(http.MethodGet, "/gobounties/bounty-cards?workspace_uuid="+workspace.Uuid, nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var response []db.BountyCard
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotEmpty(t, response, "Response should not be empty")

		cardWithoutAssignee := response[0]
		assert.Equal(t, bountyWithoutAssignee.ID, cardWithoutAssignee.BountyID)
		assert.Equal(t, bountyWithoutAssignee.Title, cardWithoutAssignee.Title)
		assert.Empty(t, cardWithoutAssignee.AssigneePic)
	})
}

func TestDeleteBountyAssignee(t *testing.T) {

	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mockHttpClient := mocks.NewHttpClient(t)

	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	db.CleanTestData()

	db.TestDB.CreateOrEditBounty(db.NewBounty{
		Type:          "coding",
		Title:         "Bounty 1",
		Description:   "Description for Bounty 1",
		WorkspaceUuid: "work-1",
		OwnerID:       "validOwner",
		Price:         1500,
		Created:       1234567890,
	})

	db.TestDB.CreateOrEditBounty(db.NewBounty{
		Type:          "design",
		Title:         "Bounty 2",
		Description:   "Description for Bounty 2",
		WorkspaceUuid: "work-2",
		OwnerID:       "nonExistentOwner",
		Price:         2000,
		Created:       1234567891,
	})

	db.TestDB.CreateOrEditBounty(db.NewBounty{
		Type:          "design",
		Title:         "Bounty 2",
		Description:   "Description for Bounty 2",
		WorkspaceUuid: "work-2",
		OwnerID:       "validOwner",
		Price:         2000,
		Created:       0,
	})

	tests := []struct {
		name           string
		input          interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   bool
	}{
		{
			name: "Valid Input - Successful Deletion",
			input: db.DeleteBountyAssignee{
				Owner_pubkey: "validOwner",
				Created:      "1234567890",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   true,
		},
		{
			name:           "Empty JSON Body",
			input:          nil,
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   false,
		},
		{
			name:           "Invalid JSON Format",
			input:          `{"Owner_pubkey": "abc", "Created": }`,
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   false,
		},
		{
			name: "Non-Existent Bounty",
			input: db.DeleteBountyAssignee{
				Owner_pubkey: "nonExistentOwner",
				Created:      "1234567890",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   false,
		},
		{
			name: "Mismatched Owner Key",
			input: db.DeleteBountyAssignee{
				Owner_pubkey: "wrongOwner",
				Created:      "1234567890",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   false,
		},
		{
			name: "Invalid Data Types",
			input: db.DeleteBountyAssignee{
				Owner_pubkey: "validOwners",
				Created:      "invalidDate",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   false,
		},
		{
			name: "Null Values",
			input: db.DeleteBountyAssignee{
				Owner_pubkey: "",
				Created:      "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   false,
		},
		{
			name: "Large JSON Body",
			input: map[string]interface{}{
				"Owner_pubkey": "validOwner",
				"Created":      "1234567890",
				"Extra":        make([]byte, 10000),
			},
			expectedStatus: http.StatusOK,
			expectedBody:   true,
		},
		{
			name: "Boundary Date Value",
			input: db.DeleteBountyAssignee{
				Owner_pubkey: "validOwner",
				Created:      "0",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var body []byte
			if tt.input != nil {
				switch v := tt.input.(type) {
				case string:
					body = []byte(v)
				default:
					var err error
					body, err = json.Marshal(tt.input)
					if err != nil {
						t.Fatalf("Failed to marshal input: %v", err)
					}
				}
			}

			req := httptest.NewRequest(http.MethodDelete, "/gobounties/assignee", bytes.NewReader(body))

			w := httptest.NewRecorder()

			bHandler.DeleteBountyAssignee(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest {

				var result bool
				err := json.NewDecoder(resp.Body).Decode(&result)
				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				assert.Equal(t, tt.expectedBody, result)
			}
		})
	}

}

func TestBountyGetFilterCount(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	tests := []struct {
		name          string
		setupBounties []db.NewBounty
		expected      db.FilterStattuCount
	}{
		{
			name:          "Empty Database",
			setupBounties: []db.NewBounty{},
			expected: db.FilterStattuCount{
				Open: 0, Assigned: 0, Completed: 0,
				Paid: 0, Pending: 0, Failed: 0,
			},
		},
		{
			name: "Only Open Bounties",
			setupBounties: []db.NewBounty{
				{
					Show:     true,
					Assignee: "",
					Paid:     false,
					OwnerID:  "test-owner-1",
					Type:     "coding",
					Title:    "Test Bounty 1",
				},
				{
					Show:     true,
					Assignee: "",
					Paid:     false,
					OwnerID:  "test-owner-2",
					Type:     "coding",
					Title:    "Test Bounty 2",
				},
			},
			expected: db.FilterStattuCount{
				Open:      2,
				Assigned:  0,
				Completed: 0,
				Paid:      0,
				Pending:   0,
				Failed:    0,
			},
		},
		{
			name: "Only Assigned Bounties",
			setupBounties: []db.NewBounty{
				{
					Show:      true,
					Assignee:  "user1",
					Paid:      false,
					Completed: false,
					OwnerID:   "test-owner-1",
					Type:      "coding",
					Title:     "Test Bounty 1",
					Created:   time.Now().Unix(),
				},
				{
					Show:      true,
					Assignee:  "user2",
					Paid:      false,
					Completed: false,
					OwnerID:   "test-owner-2",
					Type:      "coding",
					Title:     "Test Bounty 2",
					Created:   time.Now().Unix(),
				},
			},
			expected: db.FilterStattuCount{
				Open:      0,
				Assigned:  2,
				Completed: 0,
				Paid:      0,
				Pending:   0,
				Failed:    0,
			},
		},
		{
			name: "Only Completed Bounties",
			setupBounties: []db.NewBounty{
				{
					Show:      true,
					Assignee:  "user1",
					Completed: true,
					Paid:      false,
					OwnerID:   "test-owner-1",
					Type:      "coding",
					Title:     "Test Bounty 1",
					Created:   time.Now().Unix(),
				},
				{
					Show:      true,
					Assignee:  "user2",
					Completed: true,
					Paid:      false,
					OwnerID:   "test-owner-2",
					Type:      "coding",
					Title:     "Test Bounty 2",
					Created:   time.Now().Unix(),
				},
			},
			expected: db.FilterStattuCount{
				Open:      0,
				Assigned:  2,
				Completed: 2,
				Paid:      0,
				Pending:   0,
				Failed:    0,
			},
		},
		{
			name: "Only Paid Bounties",
			setupBounties: []db.NewBounty{
				{
					Show:     true,
					Assignee: "user1",
					Paid:     true,
					OwnerID:  "test-owner-1",
					Type:     "coding",
					Title:    "Test Bounty 1",
					Created:  time.Now().Unix(),
				},
				{
					Show:     true,
					Assignee: "user2",
					Paid:     true,
					OwnerID:  "test-owner-2",
					Type:     "coding",
					Title:    "Test Bounty 2",
					Created:  time.Now().Unix(),
				},
			},
			expected: db.FilterStattuCount{
				Open: 0, Assigned: 0, Completed: 0,
				Paid: 2, Pending: 0, Failed: 0,
			},
		},
		{
			name: "Only Pending Payment Bounties",
			setupBounties: []db.NewBounty{
				{
					Show:           true,
					Assignee:       "user1",
					PaymentPending: true,
					OwnerID:        "test-owner-1",
					Type:           "coding",
					Title:          "Test Bounty 1",
					Created:        time.Now().Unix(),
				},
				{
					Show:           true,
					Assignee:       "user2",
					PaymentPending: true,
					OwnerID:        "test-owner-2",
					Type:           "coding",
					Title:          "Test Bounty 2",
					Created:        time.Now().Unix(),
				},
			},
			expected: db.FilterStattuCount{
				Open: 0, Assigned: 2, Completed: 0,
				Paid: 0, Pending: 2, Failed: 0,
			},
		},
		{
			name: "Only Failed Payment Bounties",
			setupBounties: []db.NewBounty{
				{
					Show:          true,
					Assignee:      "user1",
					PaymentFailed: true,
					OwnerID:       "test-owner-1",
					Type:          "coding",
					Title:         "Test Bounty 1",
					Created:       time.Now().Unix(),
				},
				{
					Show:          true,
					Assignee:      "user2",
					PaymentFailed: true,
					OwnerID:       "test-owner-2",
					Type:          "coding",
					Title:         "Test Bounty 2",
					Created:       time.Now().Unix(),
				},
			},
			expected: db.FilterStattuCount{
				Open: 0, Assigned: 2, Completed: 0,
				Paid: 0, Pending: 0, Failed: 2,
			},
		},
		{
			name: "Hidden Bounties Should Not Count",
			setupBounties: []db.NewBounty{
				{
					Show:     false,
					Assignee: "",
					Paid:     false,
					OwnerID:  "test-owner-1",
					Type:     "coding",
					Title:    "Test Bounty 1",
					Created:  time.Now().Unix(),
				},
				{
					Show:      false,
					Assignee:  "user1",
					Completed: true,
					OwnerID:   "test-owner-2",
					Type:      "coding",
					Title:     "Test Bounty 2",
					Created:   time.Now().Unix(),
				},
			},
			expected: db.FilterStattuCount{
				Open: 0, Assigned: 0, Completed: 0,
				Paid: 0, Pending: 0, Failed: 0,
			},
		},
		{
			name: "Mixed Status Bounties",
			setupBounties: []db.NewBounty{
				{
					Show: true, Assignee: "", Paid: false,
					OwnerID: "test-owner-1", Type: "coding", Title: "Open Bounty",
					Created: time.Now().Unix(),
				},
				{
					Show: true, Assignee: "user1", Paid: false,
					OwnerID: "test-owner-2", Type: "coding", Title: "Assigned Bounty",
					Created: time.Now().Unix(),
				},
				{
					Show: true, Assignee: "user2", Completed: true, Paid: false,
					OwnerID: "test-owner-3", Type: "coding", Title: "Completed Bounty",
					Created: time.Now().Unix(),
				},
				{
					Show: true, Assignee: "user3", Paid: true,
					OwnerID: "test-owner-4", Type: "coding", Title: "Paid Bounty",
					Created: time.Now().Unix(),
				},
				{
					Show: true, Assignee: "user4", PaymentPending: true,
					OwnerID: "test-owner-5", Type: "coding", Title: "Pending Bounty",
					Created: time.Now().Unix(),
				},
				{
					Show: true, Assignee: "user5", PaymentFailed: true,
					OwnerID: "test-owner-6", Type: "coding", Title: "Failed Bounty",
					Created: time.Now().Unix(),
				},
				{
					Show: false, Assignee: "user6", Paid: true,
					OwnerID: "test-owner-7", Type: "coding", Title: "Hidden Bounty",
					Created: time.Now().Unix(),
				},
			},
			expected: db.FilterStattuCount{
				Open: 1, Assigned: 4, Completed: 1,
				Paid: 1, Pending: 1, Failed: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db.TestDB.DeleteAllBounties()

			for _, bounty := range tt.setupBounties {
				_, err := db.TestDB.CreateOrEditBounty(bounty)
				if err != nil {
					t.Fatalf("Failed to create test bounty: %v", err)
				}
			}

			rr := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/filter/count", nil)
			if err != nil {
				t.Fatal(err)
			}

			bHandler.GetFilterCount(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			var result db.FilterStattuCount
			err = json.NewDecoder(rr.Body).Decode(&result)
			if err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateBountyCardResponse(t *testing.T) {

	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	db.CleanTestData()

	workspace := db.Workspace{
		ID:          1,
		Uuid:        "test-workspace-uuid",
		Name:        "Test Workspace",
		Description: "Test Workspace Description",
		OwnerPubKey: "test-owner",
	}
	_, err := db.TestDB.CreateOrEditWorkspace(workspace)
	assert.NoError(t, err)

	phase := db.FeaturePhase{
		Uuid:        "test-phase-uuid",
		Name:        "Test Phase",
		FeatureUuid: "test-feature-uuid",
	}
	db.TestDB.CreateOrEditFeaturePhase(phase)

	feature := db.WorkspaceFeatures{
		Uuid:          "test-feature-uuid",
		Name:          "Test Feature",
		WorkspaceUuid: workspace.Uuid,
	}
	db.TestDB.CreateOrEditFeature(feature)

	assignee := db.Person{
		OwnerPubKey: "test-assignee",
		Img:         "test-image-url",
	}
	db.TestDB.CreateOrEditPerson(assignee)

	now := time.Now()

	publicBounty := db.NewBounty{
		ID:            1,
		Type:          "coding",
		Title:         "Public Bounty",
		Description:   "Test Description",
		WorkspaceUuid: workspace.Uuid,
		PhaseUuid:     phase.Uuid,
		Assignee:      assignee.OwnerPubKey,
		Show:          true,
		Created:       now.Unix(),
		OwnerID:       "test-owner",
		Price:         1000,
	}
	_, err = db.TestDB.CreateOrEditBounty(publicBounty)
	assert.NoError(t, err)

	privateBounty := db.NewBounty{
		ID:            2,
		Type:          "coding",
		Title:         "Private Bounty",
		Description:   "Test Description",
		WorkspaceUuid: workspace.Uuid,
		PhaseUuid:     phase.Uuid,
		Assignee:      assignee.OwnerPubKey,
		Show:          false,
		Created:       now.Unix(),
		OwnerID:       "test-owner",
		Price:         2000,
	}
	_, err = db.TestDB.CreateOrEditBounty(privateBounty)
	assert.NoError(t, err)

	inputBounties := []db.NewBounty{publicBounty, privateBounty}

	response := bHandler.GenerateBountyCardResponse(inputBounties)

	assert.Equal(t, 2, len(response), "Should return cards for both bounties")

	titles := make(map[string]bool)
	for _, card := range response {
		titles[card.Title] = true

		assert.Equal(t, workspace.Uuid, card.Workspace.Uuid)
		assert.Equal(t, assignee.Img, card.AssigneePic)
		assert.Equal(t, phase.Uuid, card.Phase.Uuid)
		assert.Equal(t, feature.Uuid, card.Features.Uuid)
	}

	assert.True(t, titles["Public Bounty"], "Public bounty should be present")
	assert.True(t, titles["Private Bounty"], "Private bounty should be present")
}

func TestGetWorkspaceBountyCards(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mockHttpClient := mocks.NewHttpClient(t)
	bHandler := NewBountyHandler(mockHttpClient, db.TestDB)

	db.CleanTestData()

	workspace := db.Workspace{
		ID:          1,
		Uuid:        "test-workspace-uuid",
		Name:        "Test Workspace",
		Description: "Test Workspace Description",
		OwnerPubKey: "test-owner",
	}
	_, err := db.TestDB.CreateOrEditWorkspace(workspace)
	assert.NoError(t, err)

	phase := db.FeaturePhase{
		Uuid:        "test-phase-uuid",
		Name:        "Test Phase",
		FeatureUuid: "test-feature-uuid",
	}
	db.TestDB.CreateOrEditFeaturePhase(phase)

	feature := db.WorkspaceFeatures{
		Uuid:          "test-feature-uuid",
		Name:          "Test Feature",
		WorkspaceUuid: workspace.Uuid,
	}
	db.TestDB.CreateOrEditFeature(feature)

	assignee := db.Person{
		OwnerPubKey: "test-assignee",
		Img:         "test-image-url",
	}
	db.TestDB.CreateOrEditPerson(assignee)

	now := time.Now()

	publicBounty := db.NewBounty{
		ID:            1,
		Type:          "coding",
		Title:         "Public Bounty",
		Description:   "Test Description",
		WorkspaceUuid: workspace.Uuid,
		PhaseUuid:     phase.Uuid,
		Assignee:      assignee.OwnerPubKey,
		Show:          true,
		Created:       now.Unix(),
		OwnerID:       "test-owner",
		Price:         1000,
	}

	privateBounty := db.NewBounty{
		ID:            2,
		Type:          "coding",
		Title:         "Private Bounty",
		Description:   "Test Description",
		WorkspaceUuid: workspace.Uuid,
		PhaseUuid:     phase.Uuid,
		Assignee:      assignee.OwnerPubKey,
		Show:          false,
		Created:       now.Add(time.Hour).Unix(),
		OwnerID:       "test-owner",
		Price:         2000,
	}

	fiveWeeksAgo := now.Add(-5 * 7 * 24 * time.Hour)
	threeWeeksAgo := now.Add(-3 * 7 * 24 * time.Hour)

	t.Run("should only get public bounty", func(t *testing.T) {
		db.TestDB.DeleteAllBounties()
		_, err := db.TestDB.CreateOrEditBounty(publicBounty)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyCards)

		req, err := http.NewRequest(http.MethodGet, "/gobounties/bounty-cards", nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var response []db.BountyCard
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 1, len(response))
		assert.Equal(t, "Public Bounty", response[0].Title)
	})

	t.Run("should get private bounty in workspace context", func(t *testing.T) {
		db.TestDB.DeleteAllBounties()
		_, err := db.TestDB.CreateOrEditBounty(privateBounty)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyCards)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/gobounties/bounty-cards?workspace_uuid=%s", workspace.Uuid), nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var response []db.BountyCard
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 1, len(response))
		assert.Equal(t, "Private Bounty", response[0].Title)
	})

	t.Run("should include recent unpaid bounty", func(t *testing.T) {
		db.TestDB.DeleteAllBounties()

		recentUnpaidBounty := db.NewBounty{
			ID:            1,
			Type:          "coding",
			Title:         "Recent Unpaid",
			Description:   "Test Description",
			WorkspaceUuid: workspace.Uuid,
			PhaseUuid:     phase.Uuid,
			Assignee:      assignee.OwnerPubKey,
			Show:          true,
			Created:       now.Unix(),
			OwnerID:       "test-owner",
			Price:         1000,
			Updated:       &now,
			Paid:          false,
		}
		_, err := db.TestDB.CreateOrEditBounty(recentUnpaidBounty)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyCards)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/gobounties/bounty-cards?workspace_uuid=%s", workspace.Uuid), nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var response []db.BountyCard
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 1, len(response))
		assert.Equal(t, "Recent Unpaid", response[0].Title)
	})

	t.Run("should include recent paid bounty", func(t *testing.T) {
		db.TestDB.DeleteAllBounties()

		recentPaidBounty := db.NewBounty{
			ID:            1,
			Type:          "coding",
			Title:         "Recent Paid",
			Description:   "Test Description",
			WorkspaceUuid: workspace.Uuid,
			PhaseUuid:     phase.Uuid,
			Assignee:      assignee.OwnerPubKey,
			Show:          true,
			Created:       now.Unix(),
			OwnerID:       "test-owner",
			Price:         1000,
			Updated:       &now,
			Paid:          true,
		}
		_, err := db.TestDB.CreateOrEditBounty(recentPaidBounty)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyCards)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/gobounties/bounty-cards?workspace_uuid=%s", workspace.Uuid), nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var response []db.BountyCard
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 1, len(response))
		assert.Equal(t, "Recent Paid", response[0].Title)
	})

	t.Run("should exclude old unpaid bounty", func(t *testing.T) {
		db.TestDB.DeleteAllBounties()

		oldUnpaidBounty := db.NewBounty{
			ID:            1,
			Type:          "coding",
			Title:         "Old Unpaid",
			Description:   "Test Description",
			WorkspaceUuid: workspace.Uuid,
			PhaseUuid:     phase.Uuid,
			Assignee:      assignee.OwnerPubKey,
			Show:          true,
			Created:       fiveWeeksAgo.Unix(),
			OwnerID:       "test-owner",
			Price:         1000,
			Updated:       &fiveWeeksAgo,
			Paid:          false,
		}
		_, err := db.TestDB.CreateOrEditBounty(oldUnpaidBounty)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyCards)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/gobounties/bounty-cards?workspace_uuid=%s", workspace.Uuid), nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var response []db.BountyCard
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 0, len(response))
	})

	t.Run("should exclude old paid bounty", func(t *testing.T) {
		db.TestDB.DeleteAllBounties()

		oldPaidBounty := db.NewBounty{
			ID:            1,
			Type:          "coding",
			Title:         "Old Paid",
			Description:   "Test Description",
			WorkspaceUuid: workspace.Uuid,
			PhaseUuid:     phase.Uuid,
			Assignee:      assignee.OwnerPubKey,
			Show:          true,
			Created:       threeWeeksAgo.Unix(),
			OwnerID:       "test-owner",
			Price:         1000,
			Updated:       &threeWeeksAgo,
			Paid:          true,
		}
		_, err := db.TestDB.CreateOrEditBounty(oldPaidBounty)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBountyCards)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/gobounties/bounty-cards?workspace_uuid=%s", workspace.Uuid), nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var response []db.BountyCard
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 0, len(response))
	})
}
