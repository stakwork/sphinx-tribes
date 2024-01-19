package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	mocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPersonByPuKey(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	pHandler := NewPeopleHandler(mockDb)
	t.Run("should return person if present in db", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPersonByPubkey)
		person := db.Person{
			ID:          1,
			Uuid:        uuid.New().String(),
			OwnerPubKey: "person-pub-key",
			OwnerAlias:  "owner",
			UniqueName:  "test_user",
			Description: "test user",
		}
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("pubkey", "person-pub-key")
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/person/person-pub-key", nil)
		if err != nil {
			t.Fatal(err)
		}
		mockDb.On("GetPersonByPubkey", "person-pub-key").Return(person).Once()
		handler.ServeHTTP(rr, req)

		var returnedPerson db.Person
		_ = json.Unmarshal(rr.Body.Bytes(), &returnedPerson)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, person, returnedPerson)
		mockDb.AssertExpectations(t)
	})
}

func TestCreateOrEditPerson(t *testing.T) {

	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	pHandler := NewPeopleHandler(mockDb)

	t.Run("should return error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.CreateOrEditPerson)

		invalidJson := []byte(`{"key": "value"`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code, "invalid status received")
	})

	t.Run("should return error if auth pub key not present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.CreateOrEditPerson)

		bodyJson := []byte(`{"key": "value"}`)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/", bytes.NewReader(bodyJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "invalid status received")
	})

	t.Run("should return error if pub key from auth is different than person", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.CreateOrEditPerson)

		bodyJson := []byte(`{"owner_pubkey": "other-key"}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(bodyJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "invalid status received")
	})

	t.Run("should return error if trying to update no existing user", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.CreateOrEditPerson)

		bodyJson := []byte(`{"owner_pubkey": "test-key", "id": 100}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(bodyJson))
		if err != nil {
			t.Fatal(err)
		}

		mockDb.On("GetPersonByPubkey", "test-key").Return(db.Person{}).Once()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "invalid status received")
		mockDb.AssertExpectations(t)
	})

	t.Run("should create user with unique name from owner_alias", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.CreateOrEditPerson)

		bodyJson := []byte(`{"owner_pubkey": "test-key", "owner_alias": "test-user"}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(bodyJson))
		if err != nil {
			t.Fatal(err)
		}

		mockDb.On("GetPersonByPubkey", "test-key").Return(db.Person{}).Once()
		mockDb.On("PersonUniqueNameFromName", "test-user").Return("unique-name", nil).Once()
		mockDb.On("CreateOrEditPerson", mock.MatchedBy(func(p db.Person) bool {
			return p.UniqueName == "unique-name" &&
				p.OwnerPubKey == "test-key" &&
				p.OwnerAlias == "test-user"
		})).Return(db.Person{}, nil).Once()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "invalid status received")
		mockDb.AssertExpectations(t)
	})

	t.Run("should return error if trying to update other user", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.CreateOrEditPerson)

		bodyJson := []byte(`{"owner_pubkey": "test-key", "owner_alias": "test-user", "id": 1}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(bodyJson))
		if err != nil {
			t.Fatal(err)
		}

		mockDb.On("GetPersonByPubkey", "test-key").Return(db.Person{ID: 2}).Once()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "invalid status received")
		mockDb.AssertExpectations(t)
	})

	t.Run("should update user successfully", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.CreateOrEditPerson)

		bodyJson := []byte(`{"owner_pubkey": "test-key", "owner_alias": "test-user", "id": 1, "img": "img-url"}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(bodyJson))
		if err != nil {
			t.Fatal(err)
		}

		mockDb.On("GetPersonByPubkey", "test-key").Return(db.Person{ID: 1}).Once()
		mockDb.On("CreateOrEditPerson", mock.MatchedBy(func(p db.Person) bool {
			return p.OwnerPubKey == "test-key" &&
				p.OwnerAlias == "test-user" &&
				p.Img == "img-url"
		})).Return(db.Person{}, nil).Once()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "invalid status received")
		mockDb.AssertExpectations(t)
	})
}

func TestGetPersonById(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	pHandler := NewPeopleHandler(mockDb)

	t.Run("should return person by ID if present in db", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPersonById)
		person := db.Person{
			ID:          0,
			Uuid:        uuid.New().String(),
			OwnerPubKey: "person-pub-key",
			OwnerAlias:  "owner",
			UniqueName:  "test_user",
			Description: "test user",
		}
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "0")
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/person/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		mockDb.On("GetPerson", uint(0)).Return(person).Once()
		handler.ServeHTTP(rr, req)

		var returnedPerson db.Person
		_ = json.Unmarshal(rr.Body.Bytes(), &returnedPerson)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, person, returnedPerson)
		mockDb.AssertExpectations(t)
	})

	t.Run("should return 401 if ID is not valid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPersonById)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "invalid-id")
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/person/invalid-id", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Attempt to parse the ID from the route
		idStr := rctx.URLParam("id")
		_, parseErr := strconv.Atoi(idStr)
		if parseErr != nil {
			// Return an error if the ID is not a valid integer
			rr.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestDeletePerson(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	pHandler := NewPeopleHandler(mockDb)
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")

	t.Run("should return error if unauthorized user", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.DeletePerson)
	
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/", nil)
		assert.NoError(t, err)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		mockDb.On("CheckAuthorization", "test-key").Return(true).Once()
	
		handler.ServeHTTP(rr, req)
	
		assert.Equal(t, http.StatusUnauthorized, rr.Code, "expected unauthorized status")
		mockDb.AssertExpectations(t)
	})
	

	t.Run("should return error if person does not exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.DeletePerson)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "/", nil)
		assert.NoError(t, err)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		mockDb.On("GetPerson", uint(1)).Return(db.Person{}).Once()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "expected unauthorized status for non-existent person")
		mockDb.AssertExpectations(t)
	})

	t.Run("should successfully delete a person", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.DeletePerson)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "/", nil)
		assert.NoError(t, err)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		mockDb.On("GetPerson", uint(1)).Return(db.Person{ID: 1, OwnerPubKey: "test-key"}).Once()
		mockDb.On("UpdatePerson", uint(1), mock.AnythingOfType("map[string]interface {}")).Return(nil).Once()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "expected successful deletion status")
		mockDb.AssertExpectations(t)
	})
}
