package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/lib/pq"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"

	"github.com/stretchr/testify/assert"
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
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	btHandler := NewBotHandler(db.TestDB)

	t.Run("empty list is returned when a user has no bots", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(btHandler.GetBotsByOwner)

		person := db.Person{
			Uuid:         uuid.New().String(),
			OwnerAlias:   "person",
			UniqueName:   "person",
			OwnerPubKey:  uuid.New().String(),
			PriceToMeet:  0,
			Description:  "this is test user 1",
			Tags:         pq.StringArray{},
			Extras:       db.PropertyMap{},
			GithubIssues: db.PropertyMap{},
		}
		db.TestDB.CreateOrEditPerson(person)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("pubkey", person.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/bots/owner/"+person.OwnerPubKey, nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var returnedBot []db.BotRes
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBot)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, returnedBot)
	})

	t.Run("retrieval all bots by an owner", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(btHandler.GetBotsByOwner)

		person := db.Person{
			Uuid:         uuid.New().String(),
			OwnerAlias:   "person",
			UniqueName:   "person",
			OwnerPubKey:  uuid.New().String(),
			PriceToMeet:  0,
			Description:  "this is test user 1",
			Tags:         pq.StringArray{},
			Extras:       db.PropertyMap{},
			GithubIssues: db.PropertyMap{},
		}
		db.TestDB.CreateOrEditPerson(person)

		bot := db.Bot{
			UUID:           "bot1_uuid",
			OwnerPubKey:    person.OwnerPubKey,
			OwnerAlias:     person.OwnerAlias,
			Name:           "test_bot_owner",
			UniqueName:     "test_bot_owner",
			Description:    "bot description",
			Tags:           pq.StringArray{},
			Img:            "bot-img-url",
			PricePerUse:    100,
			Created:        nil,
			Updated:        nil,
			Unlisted:       false,
			Deleted:        false,
			MemberCount:    10,
			OwnerRouteHint: "route-hint",
		}

		bot2 := db.Bot{
			UUID:           "bot2_uuid",
			OwnerPubKey:    person.OwnerPubKey,
			OwnerAlias:     person.OwnerAlias,
			Name:           "test_bot_owner2",
			UniqueName:     "test_bot_owner2",
			Description:    "bot description",
			Tags:           pq.StringArray{},
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
		db.TestDB.CreateOrEditBot(bot2)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("pubkey", person.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/bots/owner/"+person.OwnerPubKey, nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var returnedBots []db.BotRes
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBots)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Len(t, returnedBots, 2)
		assert.ElementsMatch(t, []string{bot.UUID, bot2.UUID}, []string{returnedBots[0].UUID, returnedBots[1].UUID})
	})
}

func TestSearchBots(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	btHandler := NewBotHandler(db.TestDB)

	t.Run("successful search query returns data", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(btHandler.SearchBots)

		bot := db.Bot{
			UUID:        "uuid-1",
			OwnerPubKey: "owner-pubkey-1",
			OwnerAlias:  "owner-alias-1",
			Name:        "Bot 1",
			UniqueName:  "unique-bot-1",
			Description: "Description for Bot 1",
			Tags:        pq.StringArray{"tag1", "tag2"},
			Img:         "bot-img-url-1",
			PricePerUse: 100,
		}

		bot, err := db.TestDB.CreateOrEditBot(bot)
		if err != nil {
			t.Fatalf("Failed to create bot: %v", err)
		}

		query := bot.Name

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
	})

	t.Run("empty data returned for non-matching search query", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(btHandler.SearchBots)

		query := "nonexistentbot"

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
	})
}

func TestDeleteBot(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	btHandler := NewBotHandler(db.TestDB)

	btHandler.verifyTribeUUID = func(uuid string, checkTimestamp bool) (string, error) {
		return "owner-pubkey-123", nil
	}

	t.Run("bot can be deleted by the creator of the bot", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(btHandler.DeleteBot)

		bot := db.Bot{
			UUID:        uuid.New().String(),
			OwnerPubKey: "owner-pubkey-123",
			Name:        "bot-name",
			UniqueName:  "unique-bot-name",
			Description: "bot-description",
			Tags:        pq.StringArray{"tag1", "tag2"},
			Img:         "bot-img-url",
			PricePerUse: 100,
		}

		db.TestDB.CreateOrEditBot(bot)

		ctx := context.WithValue(context.Background(), auth.ContextKey, bot.OwnerPubKey)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", bot.UUID)
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "/"+bot.UUID, nil)
		assert.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		deletedBot := db.TestDB.GetBot(bot.UUID)
		assert.Empty(t, deletedBot)
	})
}

func TestCreateOrEditBot(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	bHandler := NewBotHandler(db.TestDB)

	t.Run("should test that a 401 error is returned during bot creation if there is no bot uuid", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"Name": "Test Bot",
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
		ctx := context.WithValue(context.Background(), auth.ContextKey, "invalid_pubkey")

		requestBody := map[string]interface{}{
			"UUID": uuid.New().String(),
			"Name": "Test Bot",
		}

		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.CreateOrEditBot)

		req, err := http.NewRequestWithContext(ctx, "POST", "/", bytes.NewBuffer(requestBodyBytes))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should test that a bot gets created successfully if an authenticated user sends the right data", func(t *testing.T) {
		bHandler.verifyTribeUUID = func(uuid string, checkTimestamp bool) (string, error) {
			return "owner-pubkey-123", nil
		}

		botUUID := uuid.New().String()
		uniqueName := "testbot"

		requestBody := map[string]interface{}{
			"UUID":        botUUID,
			"OwnerPubKey": "owner-pubkey-123",
			"Name":        "test_bot",
			"UniqueName":  uniqueName,
			"Description": "bot description",
		}
		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatal(err)
		}

		bot := db.Bot{
			UUID:           botUUID,
			OwnerPubKey:    "owner-pubkey-123",
			OwnerAlias:     "your_owner",
			Name:           "test_bot",
			UniqueName:     uniqueName,
			Description:    "bot description",
			Tags:           pq.StringArray{},
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

		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(requestBodyBytes))
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), auth.ContextKey, bot.OwnerPubKey)
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
		assert.Equal(t, bot.UUID, responseData["uuid"])

		createdBot := db.TestDB.GetBot(bot.UUID)
		assert.Equal(t, bot.UUID, createdBot.UUID)
		assert.Equal(t, bot.OwnerPubKey, createdBot.OwnerPubKey)
		assert.Equal(t, bot.OwnerAlias, createdBot.OwnerAlias)
		assert.Equal(t, bot.Name, createdBot.Name)
		assert.Equal(t, bot.Description, createdBot.Description)
		assert.Equal(t, bot.Tags, createdBot.Tags)
		assert.Equal(t, bot.Img, createdBot.Img)
		assert.Equal(t, bot.PricePerUse, createdBot.PricePerUse)
		assert.Equal(t, bot.Unlisted, createdBot.Unlisted)
		assert.Equal(t, bot.Deleted, createdBot.Deleted)
		assert.Equal(t, bot.MemberCount, createdBot.MemberCount)
		assert.Equal(t, bot.OwnerRouteHint, createdBot.OwnerRouteHint)
	})

	t.Run("should test that an existing bot gets updated when passed to POST bots", func(t *testing.T) {
		bHandler.verifyTribeUUID = func(uuid string, checkTimestamp bool) (string, error) {
			return "owner-pubkey-123", nil
		}

		requestBody := map[string]interface{}{
			"UUID": "bot_uuid",
			"Name": "Updated Test Bot",
		}
		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatal(err)
		}

		bot := db.Bot{
			UUID:        "bot_uuid",
			OwnerPubKey: "owner-pubkey-123",
			Name:        "test_bot",
			UniqueName:  "test_bot",
			Description: "bot description",
			Tags:        pq.StringArray{},
		}
		db.TestDB.CreateOrEditBot(bot)

		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(requestBodyBytes))
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.WithValue(req.Context(), auth.ContextKey, bot.OwnerPubKey)
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
		assert.Equal(t, "bot_uuid", responseData["uuid"])
		assert.Equal(t, "Updated Test Bot", responseData["name"])

		updatedBot := db.TestDB.GetBot("bot_uuid")
		assert.Equal(t, "Updated Test Bot", updatedBot.Name)
	})
}

func TestGetBot(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	bHandler := NewBotHandler(db.TestDB)

	t.Run("should test that a bot can be fetched with its uuid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetBot)

		bot := db.Bot{
			UUID:        uuid.New().String(),
			OwnerPubKey: "owner-pubkey-123",
			Name:        "bot-name",
			UniqueName:  "unique-bot-name",
			Description: "bot-description",
			Tags:        pq.StringArray{"tag1", "tag2"},
			Img:         "bot-img-url",
			PricePerUse: 100,
			Tsv:         "'bot':2A,5B 'bot-descript':4B 'bot-nam':1A 'descript':6B 'name':3A 'tag1' 'tag2'",
		}

		db.TestDB.CreateOrEditBot(bot)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", bot.UUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/"+bot.UUID, nil)

		if err != nil {
			t.Fatal(err)
		}

		fetchedBot := db.TestDB.GetBot(bot.UUID)

		handler.ServeHTTP(rr, req)

		var returnedBot db.Bot
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBot)
		assert.Equal(t, http.StatusOK, rr.Code)

		assert.Equal(t, bot, returnedBot)
		assert.Equal(t, bot, fetchedBot)
	})
}

func TestGetListedBots(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	bHandler := NewBotHandler(db.TestDB)

	t.Run("should test that all bots that are not unlisted or deleted get listed", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bHandler.GetListedBots)

		db.TestDB.DeleteBot()

		bot := db.Bot{
			UUID:        "bot_uuid1",
			OwnerPubKey: "your_pubkey1",
			OwnerAlias:  "your_owner1",
			Name:        "test_bot1",
			UniqueName:  "test_bot1",
			Description: "bot description 1",
			Tags:        pq.StringArray{},
			Unlisted:    false,
			Tsv:         "'1':5B 'bot':3B 'bot1':2A 'descript':4B 'test':1A",
		}

		bot2 := db.Bot{
			UUID:        "bot_uuid2",
			OwnerPubKey: "your_pubkey2",
			OwnerAlias:  "your_owner2",
			Name:        "test_bot2",
			UniqueName:  "test_bot2",
			Description: "bot description 2",
			Tags:        pq.StringArray{},
			Unlisted:    true,
		}

		db.TestDB.CreateOrEditBot(bot)
		db.TestDB.CreateOrEditBot(bot2)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/", nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		var returnedBots []db.Bot
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBots)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.ElementsMatch(t, []db.Bot{bot}, returnedBots)
	})

}
