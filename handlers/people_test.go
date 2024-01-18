package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stakwork/sphinx-tribes/handlers/mocks"
	dbMocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	ID   string
	Name string
}

func TestGetPersonById(t *testing.T) {
	// ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := dbMocks.NewDatabase(t)
	mockClient := mocks.NewHttpClient(t)
	personHandler := NewPersonyHandler(mockClient, mockDb)
	handler := http.HandlerFunc(personHandler.GetPersonById)

	t.Run("successful retrieval of a person by ID", func(t *testing.T) {
		person := &Person{ID: "123", Name: "John Doe"}
		mockDb.On("GetPersonByID", "123").Return(person, nil)

		request, _ := http.NewRequest("GET", "/person/123", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responsePerson Person
		json.NewDecoder(rr.Body).Decode(&responsePerson)
		assert.Equal(t, person, &responsePerson)
	})
}
