package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/db"
	mocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetPersonById(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	pHandler := NewPeopleHandler(mockDb)

	t.Run("should return person by ID if present in db", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPersonById)

		personID := uint(0)
		person := db.Person{
			ID:          personID,
			Uuid:        uuid.New().String(),
			OwnerPubKey: "person-pub-key",
			OwnerAlias:  "owner",
			UniqueName:  "test_user",
			Description: "test user",
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/person/%d", personID), nil)
		if err != nil {
			t.Fatal(err)
		}

		mockDb.On("GetPerson", personID).Return(person).Once()
		handler.ServeHTTP(rr, req)

		var returnedPerson db.Person
		_ = json.Unmarshal(rr.Body.Bytes(), &returnedPerson)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, person, returnedPerson)
		mockDb.AssertExpectations(t)
	})

	t.Run("should return error if person not found by ID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(pHandler.GetPersonById)

		nonExistentID := uint(99999999)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/person/%d", nonExistentID), nil)
		if err != nil {
			t.Fatal(err)
		}

		mockDb.On("GetPerson", nonExistentID).Return(db.Person{}).Once()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		mockDb.AssertExpectations(t)
	})
}
