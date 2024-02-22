package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	mocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAdminPubkeys(t *testing.T) {
	t.Run("Should test that all admin pubkeys is returned", func(t *testing.T) {
		// set the admins and init the config to update superadmins
		os.Setenv("ADMINS", "test")
		os.Setenv("RELAY_URL", "RelayUrl")
		os.Setenv("RELAY_AUTH_KEY", "RelayAuthKey")
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
	})
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

func TestGetIsAdmin(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	aHandler := NewAuthHandler(mockDb)

	t.Run("Should test that GetIsAdmin returns a 401 error if the user is not an admin", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		pubKey := "non_admin_pubkey"
		ctx := context.WithValue(req.Context(), auth.ContextKey, pubKey)
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that a 200 status code is returned if the user is an admin", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		adminPubKey := config.SuperAdmins[0]
		ctx := context.WithValue(req.Context(), auth.ContextKey, adminPubKey)
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestRefreshToken(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	aHandler := NewAuthHandler(mockDb)

	t.Run("Should test that a user token can be refreshed", func(t *testing.T) {
		mockToken := "mock_token"
		mockUserPubkey := "mock_pubkey"
		mockPerson := db.Person{
			ID:          1,
			OwnerPubKey: mockUserPubkey,
		}
		mockDb.On("GetLnUser", mockUserPubkey).Return(int64(1)).Once()
		mockDb.On("GetPersonByPubkey", mockUserPubkey).Return(mockPerson).Once()

		// Mock JWT decoding
		mockClaims := jwt.MapClaims{
			"pubkey": mockUserPubkey,
		}
		mockDecodeJwt := func(token string) (jwt.MapClaims, error) {
			return mockClaims, nil
		}
		aHandler.decodeJwt = mockDecodeJwt

		// Mock JWT encoding
		mockEncodedToken := "encoded_mock_token"
		mockEncodeJwt := func(pubkey string) (string, error) {
			return mockEncodedToken, nil
		}
		aHandler.encodeJwt = mockEncodeJwt

		// Create request with mock token in header
		req, err := http.NewRequest("GET", "/refresh_jwt", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("x-jwt", mockToken)

		// Serve request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.RefreshToken)
		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		var responseData map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &responseData)
		if err != nil {
			t.Fatalf("Error decoding JSON response: %s", err)
		}
		assert.Equal(t, true, responseData["status"])
		assert.Equal(t, mockEncodedToken, responseData["jwt"])
	})
}
