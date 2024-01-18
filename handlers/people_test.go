package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stakwork/sphinx-tribes/handlers/mocks"
	dbMocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Person struct {
	ID   string
	Name string
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetPerson(id uint) (*Person, error) {
	args := m.Called(id)
	return args.Get(0).(*Person), args.Error(1)
}

func TestGetPersonById(t *testing.T) {
	// ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := dbMocks.NewDatabase(t)
	mockClient := mocks.NewHttpClient(t)
	personHandler := NewPersonyHandler(mockClient, mockDb)
	handler := http.HandlerFunc(personHandler.GetPersonById)

	t.Run("successful retrieval of a person by ID", func(t *testing.T) {
		person := &Person{ID: "123", Name: "John Doe"}
		mockDb.On("GetPerson", uint(123)).Return(person, nil)

		request, _ := http.NewRequest("GET", "/person/123", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responsePerson Person
		json.NewDecoder(rr.Body).Decode(&responsePerson)
		assert.Equal(t, person, &responsePerson)
	})
}
