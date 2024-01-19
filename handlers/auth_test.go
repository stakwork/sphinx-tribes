package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	mocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAdminPubkeys(t *testing.T) {
	// set the admins and init the config to update superadmins
	os.Setenv("ADMINS", "test")
	config.InitConfig()

	req, err := http.NewRequest("GET", "/admin_pubkeys", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAdminPubkeys)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"pubkeys":["test"]}`
	if strings.TrimRight(rr.Body.String(), "\n") != expected {

		t.Errorf("handler returned unexpected body: expected %s pubkeys %s is there a space after?", expected, rr.Body.String())
	}
}

func TestCreateConnectionCode(t *testing.T) {

	mockDb := mocks.NewDatabase(t)
	aHandler := NewAuthHandler(mockDb)
	t.Run("should create connection code successful", func(t *testing.T) {
		codeToBeInserted := db.ConnectionCodes{
			ConnectionString: "custom connection string",
		}
		mockDb.On("CreateConnectionCode", mock.MatchedBy(func(code db.ConnectionCodes) bool {
			return code.IsUsed == false && code.ConnectionString == codeToBeInserted.ConnectionString
		})).Return(codeToBeInserted, nil).Once()

		body, _ := json.Marshal(codeToBeInserted)
		req, err := http.NewRequest("POST", "/connectioncodes", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.CreateConnectionCode)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("should return error if failed to add connection code", func(t *testing.T) {
		codeToBeInserted := db.ConnectionCodes{
			ConnectionString: "custom connection string",
		}
		mockDb.On("CreateConnectionCode", mock.MatchedBy(func(code db.ConnectionCodes) bool {
			return code.IsUsed == false && code.ConnectionString == codeToBeInserted.ConnectionString
		})).Return(codeToBeInserted, errors.New("failed to create connection")).Once()

		body, _ := json.Marshal(codeToBeInserted)
		req, err := http.NewRequest("POST", "/connectioncodes", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.CreateConnectionCode)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("should return error for malformed request body", func(t *testing.T) {
		body := []byte(`{"id":0,"connection_string":"string","is_used":false,"date_created":"5T11:50:00Z"}`)
		req, err := http.NewRequest("POST", "/connectioncodes", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.CreateConnectionCode)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
		mockDb.AssertExpectations(t)
	})

	t.Run("should return error for invalid json", func(t *testing.T) {
		body := []byte(`{"id":0,"connection_string":"string"`)
		req, err := http.NewRequest("POST", "/connectioncodes", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.CreateConnectionCode)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
		mockDb.AssertExpectations(t)
	})
}

func TestGetConnectionCode(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	aHandler := NewAuthHandler(mockDb)

	t.Run("should return connection code from db", func(t *testing.T) {
		creationDate, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")
		existingConnectionCode := db.ConnectionCodesShort{
			ConnectionString: "test",
			DateCreated:      &creationDate,
		}
		mockDb.On("GetConnectionCode").Return(existingConnectionCode).Once()
		req, err := http.NewRequest("GET", "/connectioncodes", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetConnectionCode)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		expected := `{"connection_string":"test","date_created":"2000-01-01T00:00:00Z"}`
		assert.EqualValues(t, expected, strings.TrimRight(rr.Body.String(), "\n"))
	})

}
