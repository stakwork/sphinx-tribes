package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/lib/pq"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	dbMocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetBotByUniqueName(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	btHandler := NewBotHandler(db.TestDB)

	t.Run("successful retrieval of bots by uniqueName", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(btHandler.GetBotByUniqueName)

		bot := db.Bot{
			UUID:           "uuid-123",
			OwnerPubKey:    "owner-pubkey-123",
			OwnerAlias:     "owner-alias",
			Name:           "bot-name",
			UniqueName:     "unique-bot-name",
			Description:    "bot-description",
			Tags:           pq.StringArray{"tag1", "tag2"},
			Img:            "bot-img-url",
			PricePerUse:    100,
			Created:        nil,
			Updated:        nil,
			Unlisted:       false,
			Deleted:        false,
			MemberCount:    10,
			OwnerRouteHint: "route-hint",
		}

		db.TestDB.CreateOrEditBot(bot)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("name", bot.UniqueName)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/bot/"+bot.UniqueName, nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var returnedBot db.BotRes
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBot)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, bot.UUID, returnedBot.UUID)
		assert.Equal(t, bot.Name, returnedBot.Name)
		assert.Equal(t, bot.UniqueName, returnedBot.UniqueName)
	})
}

func TestGetBotsByOwner(t *testing.T) {

	mockDb := dbMocks.NewDatabase(t)
	btHandler := NewBotHandler(mockDb)

	t.Run("empty list is returned when a user has no bots", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(btHandler.GetBotsByOwner)

		bot := db.Bot{
			UUID:           "uuid-123",
			OwnerPubKey:    "owner-pubkey-123",
			OwnerAlias:     "owner-alias",
			Name:           "bot-name",
			UniqueName:     "unique-bot-name",
			Description:    "bot-description",
			Tags:           pq.StringArray{"tag1", "tag2"},
			Img:            "bot-img-url",
			PricePerUse:    100,
			Created:        nil,
			Updated:        nil,
			Unlisted:       false,
			Deleted:        false,
			MemberCount:    10,
			OwnerRouteHint: "route-hint",
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("pubkey", bot.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/bots/owner/"+bot.OwnerPubKey, nil)
		assert.NoError(t, err)

		mockDb.On("GetBotsByOwner", bot.OwnerPubKey).Return([]db.Bot{}, nil).Once()

		handler.ServeHTTP(rr, req)

		var returnedBot []db.BotRes
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBot)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("retrieval all bots by an owner", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(btHandler.GetBotsByOwner)

		bot := db.Bot{
			UUID:           "uuid-123",
			OwnerPubKey:    "owner-pubkey-123",
			OwnerAlias:     "owner-alias",
			Name:           "bot-name",
			UniqueName:     "unique-bot-name",
			Description:    "bot-description",
			Tags:           pq.StringArray{"tag1", "tag2"},
			Img:            "bot-img-url",
			PricePerUse:    100,
			Created:        nil,
			Updated:        nil,
			Unlisted:       false,
			Deleted:        false,
			MemberCount:    10,
			OwnerRouteHint: "route-hint",
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("pubkey", bot.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/bots/owner/"+bot.OwnerPubKey, nil)
		assert.NoError(t, err)

		mockDb.On("GetBotsByOwner", bot.OwnerPubKey).Return([]db.Bot{bot}, nil)

		handler.ServeHTTP(rr, req)

		var returnedBot []db.BotRes
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBot)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
	})
}

func TestSearchBots(t *testing.T) {
	mockDb := dbMocks.NewDatabase(t)
	btHandler := NewBotHandler(mockDb)

	t.Run("successful search query returns data", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(btHandler.SearchBots)

		query := "bot"

		bots := []db.BotRes{
			{
				UUID:        "uuid-1",
				OwnerPubKey: "owner-pubkey-1",
				Name:        "Bot 1",
				UniqueName:  "unique-bot-1",
				Description: "Description for Bot 1",
				Tags:        pq.StringArray{"tag1", "tag2"},
				Img:         "bot-img-url-1",
				PricePerUse: 100,
			},
		}

		mockDb.On("SearchBots", query, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(bots)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("query", query)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/search/bots/"+query, nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var returnedBots []db.BotRes
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBots)
		assert.NoError(t, err)
		assert.NotEmpty(t, returnedBots)

		assert.EqualValues(t, bots, returnedBots)

		mockDb.AssertExpectations(t)
	})

	t.Run("empty data returned for non-matching search query", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(btHandler.SearchBots)

		query := "nonexistentbot"

		mockDb.On("SearchBots", query, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return([]db.BotRes{})

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("query", query)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/search/bots/"+query, nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var returnedBots []db.BotRes
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBots)
		assert.NoError(t, err)
		assert.Empty(t, returnedBots)

		mockDb.AssertExpectations(t)
	})
}

func TestDeleteBot(t *testing.T) {
	t.Run("bot can be deleted by the creator of the bot", func(t *testing.T) {
		mockDb := dbMocks.NewDatabase(t)
		mockVerifyTribeUUID := func(uuid string, checkTimestamp bool) (string, error) {
			return "creator-public-key", nil
		}
		mockUUID := "123-456-789"

		btHandler := &botHandler{
			db:              mockDb,
			verifyTribeUUID: mockVerifyTribeUUID,
		}

		ctx := context.WithValue(context.Background(), auth.ContextKey, "creator-public-key")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", mockUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/"+mockUUID, nil)
		assert.NoError(t, err)

		expectedUUID := "123-456-789"

		mockDb.On("UpdateBot", mockUUID, map[string]interface{}{"deleted": true}).Return(true)

		rr := httptest.NewRecorder()

		btHandler.DeleteBot(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		mockDb.AssertCalled(t, "UpdateBot", expectedUUID, map[string]interface{}{"deleted": true})
	})
}
func TestCreateOrEditBot(t *testing.T) {
	mockDb := dbMocks.NewDatabase(t)
	bHandler := NewBotHandler(mockDb)

	t.Run("should test that a 401 error is returned during bot creation if there is no bot uuid", func(t *testing.T) {
		mockUUID := "valid_uuid"

		requestBody := map[string]interface{}{
			"UUID": mockUUID,
		}

		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(requestBodyBytes))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBot)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should test that a 401 error is returned if the user public key can't be verified during bot creation", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.ContextKey, "pubkey")
		mockPubKey := "valid_pubkey"
		mockUUID := "valid_uuid"

		requestBody := map[string]interface{}{
			"UUID": mockUUID,
		}

		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatal(err)
		}

		mockVerifyTribeUUID := func(uuid string, checkTimestamp bool) (string, error) {
			return mockPubKey, nil
		}

		bHandler.verifyTribeUUID = mockVerifyTribeUUID

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBot)

		req, err := http.NewRequestWithContext(ctx, "POST", "/", bytes.NewBuffer(requestBodyBytes))
		if err != nil {
			t.Fatal(err)
		}

		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("uuid", mockUUID)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should test that a bot gets created successfully if an authenticated user sends the right data", func(t *testing.T) {
		mockPubKey := "valid_pubkey"
		mockUUID := "valid_uuid"
		mockName := "Test Bot"
		mockUniqueName := "unique test bot"

		mockVerifyTribeUUID := func(uuid string, checkTimestamp bool) (string, error) {
			return mockPubKey, nil
		}

		bHandler.verifyTribeUUID = mockVerifyTribeUUID

		requestBody := map[string]interface{}{
			"UUID": mockUUID,
			"Name": mockName,
		}
		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatal(err)
		}

		mockDb.On("GetBotByUniqueName", mock.Anything).Return(db.Bot{
			UniqueName: mockUniqueName,
		}, nil)

		mockDb.On("CreateOrEditBot", mock.Anything).Return(db.Bot{
			UUID: mockUUID,
		}, nil)

		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(requestBodyBytes))
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), auth.ContextKey, mockPubKey)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBot)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseData map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.Equal(t, mockUUID, responseData["uuid"])
	})

	t.Run("should test that an existing bot gets updated when passed to POST bots", func(t *testing.T) {
		mockPubKey := "valid_pubkey"
		mockUUID := "valid_uuid"
		mockName := "Updated Test Bot"
		mockUniqueName := "unique test bot"

		mockVerifyTribeUUID := func(uuid string, checkTimestamp bool) (string, error) {
			return mockPubKey, nil
		}
		bHandler.verifyTribeUUID = mockVerifyTribeUUID

		requestBody := map[string]interface{}{
			"UUID": mockUUID,
			"Name": mockName,
		}
		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatal(err)
		}

		mockDb.On("GetBotByUniqueName", mock.Anything).Return(db.Bot{
			UUID:       mockUUID,
			UniqueName: mockUniqueName,
			Name:       "Original Test Bot",
		}, nil)

		mockDb.On("CreateOrEditBot", mock.Anything).Return(db.Bot{
			UUID: mockUUID,
			Name: mockName,
		}, nil)

		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(requestBodyBytes))
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.WithValue(req.Context(), auth.ContextKey, mockPubKey)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBot)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseData map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.Equal(t, mockUUID, responseData["uuid"])
		assert.Equal(t, mockName, responseData["name"])
	})

}

func TestGetBot(t *testing.T) {
	mockDb := dbMocks.NewDatabase(t)
	bHandler := NewBotHandler(mockDb)

	t.Run("should test that a bot can be fetched with its uuid", func(t *testing.T) {

		mockUUID := "valid_uuid"
		mockBot := db.Bot{UUID: mockUUID, Name: "Test Bot"}
		mockDb.On("GetBot", mock.Anything).Return(mockBot).Once()

		rr := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", mockUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/"+mockUUID, nil)

		if err != nil {
			t.Fatal(err)
		}

		handler := http.HandlerFunc(bHandler.GetBot)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var returnedBot db.Bot
		json.Unmarshal(rr.Body.Bytes(), &returnedBot)
		assert.Equal(t, mockBot, returnedBot)
	})
}

func TestGetListedBots(t *testing.T) {
	mockDb := dbMocks.NewDatabase(t)
	bHandler := NewBotHandler(mockDb)

	t.Run("should test that all bots that are not unlisted or deleted get listed", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetListedBots)

		allBots := []db.Bot{
			{UUID: "uuid1", Name: "Bot1", Unlisted: false, Deleted: false},
			{UUID: "uuid2", Name: "Bot2", Unlisted: false, Deleted: true},
			{UUID: "uuid3", Name: "Bot3", Unlisted: true, Deleted: false},
			{UUID: "uuid4", Name: "Bot4", Unlisted: true, Deleted: true},
		}

		expectedBots := []db.Bot{
			{UUID: "uuid1", Name: "Bot1", Unlisted: false, Deleted: false},
		}

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/", nil)
		assert.NoError(t, err)

		mockDb.On("GetListedBots", mock.Anything).Return(allBots)
		handler.ServeHTTP(rr, req)
		var returnedBots []db.Bot
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBots)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)

		var filteredBots []db.Bot
		for _, bot := range returnedBots {
			if !bot.Deleted && !bot.Unlisted {
				filteredBots = append(filteredBots, bot)
			}
		}

		assert.ElementsMatch(t, expectedBots, filteredBots)
		mockDb.AssertExpectations(t)
	})

}
