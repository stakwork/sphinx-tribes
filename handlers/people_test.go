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
	"github.com/lib/pq"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	mocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPersonByPuKey(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	pHandler := NewPeopleHandler(db.TestDB)
	t.Run("should return person if present in db", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPersonByPubkey)
		person := db.Person{
			ID:           104,
			Uuid:         "person_104_uuid",
			OwnerPubKey:  "person_104_pubkey",
			OwnerAlias:   "owner",
			UniqueName:   "test_user",
			Description:  "test user",
			Tags:         pq.StringArray{},
			Extras:       db.PropertyMap{},
			GithubIssues: db.PropertyMap{},
		}
		db.TestDB.CreateOrEditPerson(person)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("pubkey", person.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/person/"+person.OwnerPubKey, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var returnedPerson db.Person
		_ = json.Unmarshal(rr.Body.Bytes(), &returnedPerson)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, person, returnedPerson)
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
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	pHandler := NewPeopleHandler(db.TestDB)

	t.Run("successful retrieval", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPersonById)
		person := db.Person{
			ID:           100,
			Uuid:         "perosn_1_uuid",
			OwnerAlias:   "person",
			UniqueName:   "person",
			OwnerPubKey:  "person_1_pubkey",
			PriceToMeet:  0,
			Description:  "this is test user 1",
			Tags:         pq.StringArray{},
			Extras:       db.PropertyMap{},
			GithubIssues: db.PropertyMap{},
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", strconv.Itoa(int(person.ID)))
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/person", nil)
		assert.NoError(t, err)

		db.TestDB.CreateOrEditPerson(person)
		fetchedPerson := db.TestDB.GetPerson(person.ID)

		handler.ServeHTTP(rr, req)

		var returnedPerson db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPerson)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, person, returnedPerson)
		assert.EqualValues(t, person, fetchedPerson)
	})

	t.Run("person not found", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPersonById)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "999")
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/person", nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestDeletePerson(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	pHandler := NewPeopleHandler(mockDb)

	t.Run("successful deletion", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.DeletePerson)
		person := db.Person{
			ID:          1,
			Uuid:        "test-uuid",
			OwnerPubKey: "owner-pub-key",
			OwnerAlias:  "owner",
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")

		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)
		ctx = context.WithValue(ctx, auth.ContextKey, person.OwnerPubKey)

		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "/person", nil)
		assert.NoError(t, err)

		mockDb.On("GetPerson", person.ID).Return(person).Once()
		mockDb.On("UpdatePerson", person.ID, mock.Anything).Return(true).Once()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
	})
}

func TestGetPeopleBySearch(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	pHandler := NewPeopleHandler(mockDb)

	t.Run("should return users that match the search text", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPeopleBySearch)
		expectedPeople := []db.Person{
			{ID: 1, Uuid: "uuid1", OwnerPubKey: "pubkey1", OwnerAlias: "John Doe"},
			{ID: 2, Uuid: "uuid2", OwnerPubKey: "pubkey2", OwnerAlias: "John Smith"},
		}

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/search?search=John", nil)
		assert.NoError(t, err)

		mockDb.On("GetPeopleBySearch", mock.Anything).Return(expectedPeople)
		handler.ServeHTTP(rr, req)

		var returnedPeople []db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPeople)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, expectedPeople, returnedPeople)
		mockDb.AssertExpectations(t)
	})

	t.Run("should return an empty search result when no user matches the search text", func(t *testing.T) {
		mockDb.ExpectedCalls = nil
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPeopleBySearch)
		expectedPeople := []db.Person{}

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/search?search=user not matched", nil)
		assert.NoError(t, err)

		mockDb.On("GetPeopleBySearch", mock.Anything).Return(expectedPeople)
		handler.ServeHTTP(rr, req)

		var returnedPeople []db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPeople)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, expectedPeople, returnedPeople)
		mockDb.AssertExpectations(t)
	})
}

func TestGetListedPeople(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	pHandler := NewPeopleHandler(mockDb)

	t.Run("should return all listed users", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetListedPeople)
		expectedPeople := []db.Person{
			{ID: 1, Uuid: "uuid1", OwnerPubKey: "pubkey1", OwnerAlias: "John Doe"},
			{ID: 2, Uuid: "uuid2", OwnerPubKey: "pubkey2", OwnerAlias: "John Smith"},
		}

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/?page=1&limit=10", nil)
		assert.NoError(t, err)

		mockDb.On("GetListedPeople", mock.Anything).Return(expectedPeople)
		handler.ServeHTTP(rr, req)

		var returnedPeople []db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPeople)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, expectedPeople, returnedPeople)
		mockDb.AssertExpectations(t)
	})

	t.Run("should return only users that match a search text when a search is added to the URL query", func(t *testing.T) {
		mockDb.ExpectedCalls = nil
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetListedPeople)
		expectedPeople := []db.Person{
			{ID: 1, Uuid: "uuid1", OwnerPubKey: "pubkey1", OwnerAlias: "John Doe"},
			{ID: 2, Uuid: "uuid2", OwnerPubKey: "pubkey2", OwnerAlias: "John Smith"},
		}

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/?page=1&limit=10&search=John", nil)
		assert.NoError(t, err)

		mockDb.On("GetListedPeople", mock.Anything).Return(expectedPeople)
		handler.ServeHTTP(rr, req)

		var returnedPeople []db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPeople)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, expectedPeople, returnedPeople)
		mockDb.AssertExpectations(t)
	})

	t.Run("should return only users that match a skill set when languages are passed to the URL query", func(t *testing.T) {
		mockDb.ExpectedCalls = nil
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetListedPeople)
		expectedPeople := []db.Person{
			{ID: 1, Uuid: "uuid1", OwnerPubKey: "pubkey1", OwnerAlias: "John Doe"},
		}

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/?page=1&limit=10&languages=typescript", nil)
		assert.NoError(t, err)

		mockDb.On("GetListedPeople", mock.Anything).Return(expectedPeople)
		handler.ServeHTTP(rr, req)

		var returnedPeople []db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPeople)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, expectedPeople, returnedPeople)
		mockDb.AssertExpectations(t)
	})
}

func TestGetPersonByUuid(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	pHandler := NewPeopleHandler(mockDb)

	t.Run("should return a user with the right UUID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPersonByUuid)
		expectedPerson := db.Person{
			ID:          1,
			Uuid:        uuid.New().String(),
			OwnerPubKey: "person-pub-key",
			OwnerAlias:  "owner",
			UniqueName:  "test_user",
			Description: "test user",
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", "uuid")
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/uuid", nil)
		assert.NoError(t, err)

		mockDb.On("GetPersonByUuid", mock.Anything).Return(expectedPerson)
		handler.ServeHTTP(rr, req)

		var returnedPerson db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPerson)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, expectedPerson, returnedPerson)
		mockDb.AssertExpectations(t)
	})

	t.Run("should return no user for a wrong UUID", func(t *testing.T) {
		mockDb.ExpectedCalls = nil
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPersonByUuid)
		expectedPerson := db.Person{}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", "wrong-uuid")
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/uuid", nil)
		assert.NoError(t, err)

		mockDb.On("GetPersonByUuid", mock.Anything).Return(expectedPerson)
		handler.ServeHTTP(rr, req)

		var returnedPerson db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPerson)
		assert.NoError(t, err)
		assert.EqualValues(t, expectedPerson, returnedPerson)
		mockDb.AssertExpectations(t)
	})
}
