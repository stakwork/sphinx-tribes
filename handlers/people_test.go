package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
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

func TestCreatePerson(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	pHandler := NewPeopleHandler(db.TestDB)

	person := db.Person{
		OwnerAlias:   "person",
		OwnerPubKey:  uuid.New().String(),
		PriceToMeet:  0,
		Description:  "this is test user 1",
		Tags:         pq.StringArray{},
		Extras:       db.PropertyMap{},
		GithubIssues: db.PropertyMap{},
		Img:          "img-url",
	}

	t.Run("should return error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.CreatePerson)

		invalidJson := []byte(`{"key": "value"`)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code, "invalid status received")
	})

	t.Run("should return error if auth pub key not present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.CreatePerson)

		bodyJson := []byte(`{"key": "value"}`)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(bodyJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "invalid status received")
	})

	t.Run("should return error if pub key from auth is different than person", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.CreatePerson)

		bodyJson := []byte(`{"owner_pubkey": "other-key"}`)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(bodyJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "invalid status received")
	})

	t.Run("should create user with unique name from owner_alias", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.CreatePerson)

		requestBody, _ := json.Marshal(person)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		fetchedUpdatedPerson := db.TestDB.GetPersonByPubkey(person.OwnerPubKey)

		person.ID = fetchedUpdatedPerson.ID

		person.Created = fetchedUpdatedPerson.Created
		person.Updated = fetchedUpdatedPerson.Updated
		person.UniqueName = fetchedUpdatedPerson.UniqueName
		person.Uuid = fetchedUpdatedPerson.Uuid

		assert.Equal(t, http.StatusOK, rr.Code, "invalid status received")

		assert.EqualValues(t, person, fetchedUpdatedPerson)
	})

	t.Run("Should return a 200 status code when existing user hits the endpoint", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.CreatePerson)

		requestBody, _ := json.Marshal(person)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		fetchedUpdatedPerson := db.TestDB.GetPersonByPubkey(person.OwnerPubKey)

		assert.Equal(t, http.StatusOK, rr.Code, "invalid status received")

		assert.Equal(t, person.Description, fetchedUpdatedPerson.Description)
		assert.Equal(t, person.OwnerAlias, fetchedUpdatedPerson.OwnerAlias)
		assert.Equal(t, person.PriceToMeet, fetchedUpdatedPerson.PriceToMeet)
	})
}

func TestUpdatePerson(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	pHandler := NewPeopleHandler(db.TestDB)

	t.Run("should return error if trying to update non-existing user", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.UpdatePerson)

		bodyJson := []byte(`{"owner_pubkey": "test-key"}`)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/", bytes.NewReader(bodyJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code, "invalid status received")
	})

	t.Run("should return error if trying to update with user keys not matching", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.UpdatePerson)

		bodyJson := []byte(`{"owner_pubkey": "fake-key"}`)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/", bytes.NewReader(bodyJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "invalid status received")
	})

	t.Run("should return error if trying to update other user", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.UpdatePerson)

		bodyJson := []byte(`{"owner_pubkey": "fake-key", "owner_alias": "test-user"}`)
		ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/", bytes.NewReader(bodyJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "invalid status received")
	})

	t.Run("should update user successfully", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.UpdatePerson)

		now := time.Now()

		// First, create a person
		person := db.Person{
			Uuid:         uuid.New().String(),
			OwnerAlias:   "person-update",
			UniqueName:   "person-update",
			OwnerPubKey:  uuid.New().String(),
			Description:  "this is a update test user",
			Img:          "img-url",
			Tags:         pq.StringArray{},
			Extras:       db.PropertyMap{},
			GithubIssues: db.PropertyMap{},
			PriceToMeet:  40,
			Created:      &now,
			Updated:      &now,
		}

		createdPerson, err := db.TestDB.CreateOrEditPerson(person)
		if err != nil {
			t.Fatal(err)
		}

		// Update the created person
		updatePerson := db.Person{
			ID:           createdPerson.ID,
			Uuid:         createdPerson.Uuid,
			OwnerAlias:   "person-update-after",
			UniqueName:   "person-update-affer",
			OwnerPubKey:  createdPerson.OwnerPubKey,
			Description:  "this is after updated test user",
			Tags:         pq.StringArray{},
			Extras:       db.PropertyMap{},
			GithubIssues: db.PropertyMap{},
			Img:          "img-url",
			PriceToMeet:  100,
		}

		requestBody, _ := json.Marshal(updatePerson)

		ctx := context.WithValue(context.Background(), auth.ContextKey, updatePerson.OwnerPubKey)

		req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		fetchedUpdatedPerson := db.TestDB.GetPersonByUuid(updatePerson.Uuid)

		assert.Equal(t, http.StatusOK, rr.Code, "invalid status received")

		assert.NotEqual(t, fetchedUpdatedPerson.OwnerAlias, person.OwnerAlias)
		assert.NotEqual(t, fetchedUpdatedPerson.UniqueName, person.UniqueName)
		assert.NotEqual(t, fetchedUpdatedPerson.Description, person.Description)

		assert.Equal(t, fetchedUpdatedPerson.OwnerAlias, updatePerson.OwnerAlias)
		assert.Equal(t, fetchedUpdatedPerson.UniqueName, updatePerson.UniqueName)
		assert.Equal(t, fetchedUpdatedPerson.Description, updatePerson.Description)
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
			ID:           300,
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

		person.Created = fetchedPerson.Created
		person.Updated = fetchedPerson.Updated

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
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	pHandler := NewPeopleHandler(db.TestDB)

	t.Run("successful deletion", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.DeletePerson)
		person := db.Person{
			ID:           112,
			Uuid:         "person_112_uuid",
			OwnerPubKey:  "person_112_pubkey",
			OwnerAlias:   "owner",
			UniqueName:   "test_user",
			Description:  "test user",
			Tags:         pq.StringArray{},
			Extras:       db.PropertyMap{},
			GithubIssues: db.PropertyMap{},
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", strconv.Itoa(int(person.ID)))

		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)
		ctx = context.WithValue(ctx, auth.ContextKey, person.OwnerPubKey)

		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "/person", nil)
		assert.NoError(t, err)

		db.TestDB.CreateOrEditPerson(person)
		fetchedPerson := db.TestDB.GetPerson(person.ID)
		assert.Equal(t, person, fetchedPerson)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		deletedPerson := db.TestDB.GetPerson(person.ID)
		assert.Empty(t, deletedPerson)
	})
}

func TestGetPeopleBySearch(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	pHandler := NewPeopleHandler(db.TestDB)

	db.CleanDB()

	personV1 := db.Person{
		ID:             102,
		Uuid:           "perosn_102_uuid",
		OwnerAlias:     "person102",
		UniqueName:     "person102",
		OwnerPubKey:    "person_102_pubkey",
		OwnerRouteHint: "03a6ea2d9ead2120b12bd66292bb4a302c756983dc45dcb2b364b461c66fd53bcb:1099527159809",
		PriceToMeet:    0,
		Description:    "This is test user 102",
		Tags:           pq.StringArray{},
		Extras:         db.PropertyMap{},
		GithubIssues:   db.PropertyMap{},
	}
	personV2 := db.Person{
		ID:             103,
		Uuid:           "perosn_103_uuid",
		OwnerAlias:     "person103",
		UniqueName:     "person103",
		OwnerPubKey:    "person_103_pubkey",
		OwnerRouteHint: "034bcc332390470cc4f9ef7491af1da2ffceefccd39ceb6acd87c83920543013d7_529771090604130310",
		PriceToMeet:    0,
		Description:    "This is test user 103",
		Tags:           pq.StringArray{},
		Extras:         db.PropertyMap{},
		GithubIssues:   db.PropertyMap{},
	}
	db.TestDB.CreateOrEditPerson(personV1)
	db.TestDB.CreateOrEditPerson(personV2)

	t.Run("should return users that V2 pubkeys person and not return V1 pubkeys person", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPeopleBySearch)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodGet, "/search?search=person", nil)
		assert.NoError(t, err)

		fetchedPersonV1 := db.TestDB.GetPerson(personV1.ID)
		fetchedPersonV2 := db.TestDB.GetPerson(personV2.ID)

		// Verify both people exist in the database
		assert.NotEmpty(t, fetchedPersonV1)
		assert.NotEmpty(t, fetchedPersonV2)

		handler.ServeHTTP(rr, req)

		var returnedPeople []db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPeople)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)

		// Verify only V2 person is returned
		assert.Equal(t, 1, len(returnedPeople))
		assert.Equal(t, fetchedPersonV2.ID, returnedPeople[0].ID)

		// Explicitly verify V1 person is not in the results
		for _, person := range returnedPeople {
			assert.NotEqual(t, fetchedPersonV1.ID, person.ID)
		}
	})

	t.Run("should return users that match the search text (only V2 pubkeys)", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPeopleBySearch)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/search?search="+personV2.OwnerAlias, nil)
		assert.NoError(t, err)

		fetchedPersonV2 := db.TestDB.GetPerson(personV2.ID)

		expectedPeople := []db.Person{
			fetchedPersonV2,
		}

		handler.ServeHTTP(rr, req)

		var returnedPeople []db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPeople)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, personV2, fetchedPersonV2)
		assert.EqualValues(t, expectedPeople, returnedPeople)
	})

	t.Run("should return an empty search result when no user matches the search text", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPeopleBySearch)
		expectedPeople := []db.Person{}

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/search?search=user not matched", nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var returnedPeople []db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPeople)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, expectedPeople, returnedPeople)
	})
}

func TestGetListedPeople(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	pHandler := NewPeopleHandler(db.TestDB)

	db.CleanDB()

	person := db.Person{
		ID:           uint(rand.Intn(1000)),
		Uuid:         "person_101_uuid",
		OwnerAlias:   "person101",
		UniqueName:   "person101",
		OwnerPubKey:  "person_101_pubkey",
		PriceToMeet:  0,
		Description:  "this is test user 1",
		Unlisted:     true,
		Tags:         pq.StringArray{},
		GithubIssues: db.PropertyMap{},
		Extras:       db.PropertyMap{"coding_languages": "Typescript"},
	}
	person2 := db.Person{
		ID:           uint(rand.Intn(1000)),
		Uuid:         "person_102_uuid",
		OwnerAlias:   "person102",
		UniqueName:   "person102",
		OwnerPubKey:  "person_102_pubkey",
		PriceToMeet:  0,
		Description:  "This is test user 2",
		Unlisted:     false,
		Tags:         pq.StringArray{},
		GithubIssues: db.PropertyMap{},
		Extras:       db.PropertyMap{"coding_languages": "Golang"},
	}
	person3 := db.Person{
		ID:           uint(rand.Intn(1000)),
		Uuid:         "person_103_uuid",
		OwnerAlias:   "person103",
		UniqueName:   "person103",
		OwnerPubKey:  "person_103_pubkey",
		PriceToMeet:  0,
		Description:  "This is test user 3",
		Unlisted:     false,
		Tags:         pq.StringArray{},
		GithubIssues: db.PropertyMap{},
		Extras:       db.PropertyMap{"coding_languages": "Lightning"},
	}

	db.TestDB.CreateOrEditPerson(person)
	db.TestDB.CreateOrEditPerson(person2)
	db.TestDB.CreateOrEditPerson(person3)

	fetchedPerson2 := db.TestDB.GetPerson(person2.ID)
	fetchedPerson3 := db.TestDB.GetPerson(person3.ID)
	person2.ID = fetchedPerson2.ID
	person3.ID = fetchedPerson3.ID

	t.Run("should return all listed users", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetListedPeople)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/?page=1&limit=10", nil)
		assert.NoError(t, err)

		expectedPeople := []db.Person{
			fetchedPerson2,
			fetchedPerson3,
		}

		handler.ServeHTTP(rr, req)

		var returnedPeople []db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPeople)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, person2, fetchedPerson2)
		assert.EqualValues(t, person3, fetchedPerson3)
		assert.EqualValues(t, expectedPeople, returnedPeople)
	})

	t.Run("should return only users that match a search text when a search is added to the URL query", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetListedPeople)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/?page=1&limit=10&search="+person2.OwnerAlias, nil)
		assert.NoError(t, err)

		expectedPeople := []db.Person{
			fetchedPerson2,
		}

		handler.ServeHTTP(rr, req)

		var returnedPeople []db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPeople)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, person2, fetchedPerson2)
		assert.EqualValues(t, expectedPeople, returnedPeople)
	})

	t.Run("should return only users that match a skill set when languages are passed to the URL query", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetListedPeople)

		rctx := chi.NewRouteContext()
		languages := person2.Extras["coding_languages"].(string)
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodGet,
			"page=1&limit=10&languages="+languages,
			nil,
		)
		assert.NoError(t, err)

		expectedPeople := []db.Person{
			fetchedPerson2,
		}

		handler.ServeHTTP(rr, req)

		var returnedPeople []db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPeople)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, person2, fetchedPerson2)
		assert.EqualValues(t, expectedPeople, returnedPeople)
	})

}
func TestGetPersonByUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	pHandler := NewPeopleHandler(db.TestDB)

	t.Run("should return a user with the right UUID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPersonByUuid)
		person := db.Person{
			ID:           105,
			Uuid:         uuid.New().String(),
			OwnerAlias:   "person101",
			UniqueName:   "person101",
			OwnerPubKey:  "person_101_pubkey",
			PriceToMeet:  0,
			Description:  "this is test user 1",
			Tags:         pq.StringArray{},
			Extras:       db.PropertyMap{},
			GithubIssues: db.PropertyMap{},
		}
		db.TestDB.CreateOrEditPerson(person)
		fetchedPerson := db.TestDB.GetPerson(person.ID)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", person.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/uuid", nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var returnedPerson db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPerson)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)

		if returnedPerson.Extras == nil {
			returnedPerson.Extras = db.PropertyMap{}
		}

		assert.EqualValues(t, fetchedPerson, returnedPerson)
	})

	t.Run("should return no user for a wrong UUID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPersonByUuid)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", "wrong-uuid")
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/uuid", nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)

		var returnedPerson db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &returnedPerson)
		assert.NoError(t, err)
		assert.Empty(t, returnedPerson)
	})
}
