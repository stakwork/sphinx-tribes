package handlers

import (
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

	mockDb := dbMocks.NewDatabase(t)
	btHandler := NewBotHandler(mockDb)

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

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("name", bot.UniqueName)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/bot/"+bot.UniqueName, nil)
		assert.NoError(t, err)

		mockDb.On("GetBotByUniqueName", bot.UniqueName).Return(db.Bot{}, nil).Once()

		handler.ServeHTTP(rr, req)

		var returnedBot db.BotRes
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBot)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
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
