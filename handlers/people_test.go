package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/handlers/mocks"
	dbMocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Person struct {
	ID   string
	Name string
}

func TestGetPersonById(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := dbMocks.NewDatabase(t)
	mockClient := mocks.NewHttpClient(t)
	personHandler := NewPersonyHandler(mockClient, mockDb)
	handler := http.HandlerFunc(personHandler.GetPersonById)

	t.Run("successful retrieval of a person by ID", func(t *testing.T) {
		person := &Person{ID: "123", Name: "John Doe"}
		mockDb.On("GetPersonByID", "123").Return(person, nil)
		

		request, _ := http.NewRequest("GET", "/person?id=123", nil)
		request = request.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responsePerson Person
		json.NewDecoder(rr.Body).Decode(&responsePerson)
		assert.Equal(t, person, &responsePerson)
	})

	t.Run("error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/person?id=invalid-json", bytes.NewBufferString("invalid-json"))
		request = request.WithContext(ctx)

		handler.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("non-existent person ID", func(t *testing.T) {
		mockDb.On("GetPersonByID", "non-existent").Return(nil, errors.New("not found"))

		request, _ := http.NewRequest("GET", "/person?id=non-existent", nil)
		request = request.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("database error", func(t *testing.T) {
		mockDb.On("GetPersonByID", mock.Anything).Return(nil, errors.New("internal server error"))

		request, _ := http.NewRequest("GET", "/person?id=any", nil)
		request = request.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("unauthorized access", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/person?id=123", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
