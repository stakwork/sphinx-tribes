package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	datamocks "github.com/stakwork/sphinx-tribes/mocks"

	"github.com/form3tech-oss/jwt-go"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
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
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	aHandler := NewAuthHandler(db.TestDB)

	t.Run("should create connection code successful", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.CreateConnectionCode)
		data := db.InviteBody{
			Number:     2,
			Pubkey:     "Test_pubkey",
			RouteHint:  "Test_Route_hint",
			SatsAmount: 21,
		}

		aHandler.makeConnectionCodeRequest = func(inviter_pubkey string, inviter_route_hint string, msats_amount uint64) string {
			return "22222222222222222"
		}

		body, _ := json.Marshal(data)

		req, err := http.NewRequest(http.MethodPost, "/connectioncodes", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		codes := db.TestDB.GetConnectionCode()
		assert.NotEmpty(t, codes)
	})

	t.Run("should return error if failed to add connection code", func(t *testing.T) {
		data := db.InviteBody{
			Number: 0,
		}

		body, _ := json.Marshal(data)

		req, err := http.NewRequest("POST", "/connectioncodes", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.CreateConnectionCode)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return error for malformed request body", func(t *testing.T) {
		body := []byte(`{"number": "0"}`)
		req, err := http.NewRequest("POST", "/connectioncodes", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.CreateConnectionCode)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("should return error for invalid json", func(t *testing.T) {
		body := []byte(`{"nonumber":0`)

		req, err := http.NewRequest("POST", "/connectioncodes", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.CreateConnectionCode)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})
}

func TestGetConnectionCode(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	aHandler := NewAuthHandler(db.TestDB)

	t.Run("should return connection code from db", func(t *testing.T) {

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetConnectionCode)

		codeStrArr := []string{"sampleCode1"}

		codeArr := []db.ConnectionCodes{}
		now := time.Now()

		for i, code := range codeStrArr {
			code := db.ConnectionCodes{
				ID:               uint(i),
				ConnectionString: code,
				IsUsed:           false,
				DateCreated:      &now,
			}

			codeArr = append(codeArr, code)
		}

		// Ensure codeArr has at least one element
		codeShort := db.ConnectionCodesShort{
			ConnectionString: codeArr[0].ConnectionString,
			DateCreated:      codeArr[0].DateCreated,
		}

		db.TestDB.CreateConnectionCode(codeArr)

		req, err := http.NewRequest("GET", "/connectioncodes", nil)
		if err != nil {
			t.Fatal(err)
		}

		fetchedCodes := db.TestDB.GetConnectionCode()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.EqualValues(t, codeShort.ConnectionString, fetchedCodes.ConnectionString)
		tolerance := time.Millisecond
		timeDifference := codeShort.DateCreated.Sub(*fetchedCodes.DateCreated)
		if timeDifference < 0 {
			timeDifference = -timeDifference
		}
		assert.True(t, timeDifference <= tolerance, "Expected DateCreated to be within tolerance")

	})

}

func TestGetIsAdmin(t *testing.T) {
	mockDb := datamocks.NewDatabase(t)
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
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	aHandler := NewAuthHandler(db.TestDB)

	t.Run("Should test that a user token can be refreshed", func(t *testing.T) {
		mockToken := "mock_token"
		person := db.Person{
			Uuid:         uuid.New().String(),
			OwnerPubKey:  "your_pubkey",
			OwnerAlias:   "your_owner",
			UniqueName:   "your_user",
			Description:  "your user",
			Tags:         pq.StringArray{},
			Extras:       db.PropertyMap{},
			GithubIssues: db.PropertyMap{},
		}
		db.TestDB.CreateOrEditPerson(person)

		// Mock JWT decoding
		mockClaims := jwt.MapClaims{
			"pubkey": person.OwnerPubKey,
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

		fetchedPerson := db.TestDB.GetPersonByUuid(person.Uuid)
		person.ID = fetchedPerson.ID
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
		assert.EqualValues(t, person, fetchedPerson)
	})
}
