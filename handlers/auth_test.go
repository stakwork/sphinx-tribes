package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

	t.Run("Should handle multiple admin pubkeys", func(t *testing.T) {
		os.Setenv("ADMINS", "test1,test2,test3")
		os.Setenv("RELAY_URL", "RelayUrl")
		os.Setenv("RELAY_AUTH_KEY", "RelayAuthKey")
		config.InitConfig()

		req, err := http.NewRequest("GET", "/admin_pubkeys", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetAdminPubkeys)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		expected := `{"pubkeys":["test1","test2","test3"]}`
		assert.JSONEq(t, expected, strings.TrimRight(rr.Body.String(), "\n"))
	})

	t.Run("Should handle empty admin pubkeys", func(t *testing.T) {
		os.Setenv("ADMINS", "")
		os.Setenv("RELAY_URL", "RelayUrl")
		os.Setenv("RELAY_AUTH_KEY", "RelayAuthKey")
		config.InitConfig()

		req, err := http.NewRequest("GET", "/admin_pubkeys", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetAdminPubkeys)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		expected := `{"pubkeys":[]}`
		assert.JSONEq(t, expected, strings.TrimRight(rr.Body.String(), "\n"))
	})

	t.Run("Should handle admin pubkeys with special characters", func(t *testing.T) {
		os.Setenv("ADMINS", "test@123,test#456,test$789")
		os.Setenv("RELAY_URL", "RelayUrl")
		os.Setenv("RELAY_AUTH_KEY", "RelayAuthKey")
		config.InitConfig()

		req, err := http.NewRequest("GET", "/admin_pubkeys", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetAdminPubkeys)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		expected := `{"pubkeys":["test@123","test#456","test$789"]}`
		assert.JSONEq(t, expected, strings.TrimRight(rr.Body.String(), "\n"))
	})

	t.Run("Should handle admin pubkeys with spaces", func(t *testing.T) {
		os.Setenv("ADMINS", "test 123, test 456 , test 789")
		os.Setenv("RELAY_URL", "RelayUrl")
		os.Setenv("RELAY_AUTH_KEY", "RelayAuthKey")
		config.InitConfig()

		req, err := http.NewRequest("GET", "/admin_pubkeys", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetAdminPubkeys)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		expected := `{"pubkeys":["test 123","test 456","test 789"]}`
		assert.JSONEq(t, expected, strings.TrimRight(rr.Body.String(), "\n"))
	})

	t.Run("Should handle invalid HTTP method", func(t *testing.T) {
		os.Setenv("ADMINS", "test")
		os.Setenv("RELAY_URL", "RelayUrl")
		os.Setenv("RELAY_AUTH_KEY", "RelayAuthKey")
		config.InitConfig()

		req, err := http.NewRequest(http.MethodPost, "/admin_pubkeys", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetAdminPubkeys)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		expected := `{"pubkeys":["test"]}`
		assert.JSONEq(t, expected, strings.TrimRight(rr.Body.String(), "\n"))
	})
	t.Run("Maximum Number of Admin Keys", func(t *testing.T) {

		var keys []string
		for i := 0; i < 1000; i++ {
			keys = append(keys, fmt.Sprintf("key%d", i))
		}
		os.Setenv("ADMINS", strings.Join(keys, ","))
		os.Setenv("RELAY_URL", "RelayUrl")
		os.Setenv("RELAY_AUTH_KEY", "RelayAuthKey")
		config.InitConfig()

		req, err := http.NewRequest("GET", "/admin_pubkeys", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetAdminPubkeys)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response map[string][]string
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1000, len(response["pubkeys"]))
	})

	t.Run("Null Admin Keys List", func(t *testing.T) {
		os.Unsetenv("ADMINS")
		os.Setenv("RELAY_URL", "RelayUrl")
		os.Setenv("RELAY_AUTH_KEY", "RelayAuthKey")
		config.InitConfig()

		req, err := http.NewRequest("GET", "/admin_pubkeys", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetAdminPubkeys)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		expected := `{"pubkeys":[]}`
		assert.JSONEq(t, expected, strings.TrimRight(rr.Body.String(), "\n"))
	})

	t.Run("Unicode Characters in Admin Keys", func(t *testing.T) {
		os.Setenv("ADMINS", "ключ1,キー2,钥匙3")
		os.Setenv("RELAY_URL", "RelayUrl")
		os.Setenv("RELAY_AUTH_KEY", "RelayAuthKey")
		config.InitConfig()

		req, err := http.NewRequest("GET", "/admin_pubkeys", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetAdminPubkeys)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		expected := `{"pubkeys":["ключ1","キー2","钥匙3"]}`
		assert.JSONEq(t, expected, strings.TrimRight(rr.Body.String(), "\n"))
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

		req, err := http.NewRequest(http.MethodGet, "/connectioncodes", nil)
		if err != nil {
			t.Fatal(err)
		}

		fetchedCodes := db.TestDB.GetConnectionCode()

		fmt.Println("fetchedCodes", fetchedCodes)

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
	t.Run("should return empty fields if no connection codes exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetConnectionCode)

		req, err := http.NewRequest(http.MethodGet, "/connectioncodes", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response db.ConnectionCodesShort
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatal("Failed to unmarshal response:", err)
		}

		assert.Empty(t, response.ConnectionString)
		assert.Nil(t, response.DateCreated)
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

	t.Run("Should test that empty public key returns unauthorized", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, "")
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		var responseBody string
		json.NewDecoder(rr.Body).Decode(&responseBody)
		assert.Equal(t, "Not a super admin: handler", responseBody)
	})

	t.Run("Should test that nil context value returns unauthorized", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, nil)
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		var responseBody string
		json.NewDecoder(rr.Body).Decode(&responseBody)
		assert.Equal(t, "Not a super admin: handler", responseBody)
	})

	t.Run("Should test that free pass enabled allows any user", func(t *testing.T) {

		originalAdmins := config.SuperAdmins
		config.SuperAdmins = []string{config.AdminDevFreePass}
		defer func() {
			config.SuperAdmins = originalAdmins
		}()

		req, err := http.NewRequest("GET", "/admin/auth", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, "any_pubkey")
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody string
		json.NewDecoder(rr.Body).Decode(&responseBody)
		assert.Equal(t, "Log in successful", responseBody)
	})

	t.Run("Should test that invalid context value type returns unauthorized", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, 12345)
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		var responseBody string
		json.NewDecoder(rr.Body).Decode(&responseBody)
		assert.Equal(t, "Not a super admin: handler", responseBody)
	})

	t.Run("Should test multiple admins configuration", func(t *testing.T) {

		originalAdmins := config.SuperAdmins
		config.SuperAdmins = []string{"admin1", "admin2", "admin3"}
		defer func() {
			config.SuperAdmins = originalAdmins
		}()

		req, err := http.NewRequest("GET", "/admin/auth", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, "admin2")
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody string
		json.NewDecoder(rr.Body).Decode(&responseBody)
		assert.Equal(t, "Log in successful", responseBody)
	})

	t.Run("Admin User with Free Pass Disabled", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, config.SuperAdmins[0])
		req = req.WithContext(ctx)

		config.AdminDevFreePass = ""
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Non-Admin User with Free Pass Disabled", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, "non_admin_pubkey")
		req = req.WithContext(ctx)

		config.AdminDevFreePass = ""
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Admin User with Free Pass Enabled", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, config.SuperAdmins[0])
		req = req.WithContext(ctx)

		config.AdminDevFreePass = "freepass"
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Non-Admin User with Free Pass Enabled", func(t *testing.T) {

		originalAdmins := config.SuperAdmins
		config.SuperAdmins = []string{config.AdminDevFreePass}
		defer func() {
			config.SuperAdmins = originalAdmins
		}()

		req, err := http.NewRequest("GET", "/admin/auth", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, "any_pubkey")
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody string
		json.NewDecoder(rr.Body).Decode(&responseBody)
		assert.Equal(t, "Log in successful", responseBody)
	})

	t.Run("Empty Public Key in Context", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, "")
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Nil Context", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid Data Type for Public Key", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, 12345)
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Missing Context Key", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Large Number of Admins in Configuration", func(t *testing.T) {
		originalAdmins := config.SuperAdmins
		config.SuperAdmins = make([]string, 1000)
		for i := 0; i < 1000; i++ {
			config.SuperAdmins[i] = fmt.Sprintf("admin%d", i)
		}
		defer func() { config.SuperAdmins = originalAdmins }()

		req, err := http.NewRequest("GET", "/admin/auth", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, "admin500")
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Admin User with Invalid Free Pass State", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, config.SuperAdmins[0])
		req = req.WithContext(ctx)

		config.AdminDevFreePass = "invalid"
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Non-Admin User with Invalid Free Pass State", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/admin/auth", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.GetIsAdmin)

		ctx := context.WithValue(req.Context(), auth.ContextKey, "non_admin_pubkey")
		req = req.WithContext(ctx)

		config.AdminDevFreePass = "invalid"
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
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

	t.Run("Empty JWT Token", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/refresh_jwt", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.RefreshToken)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("JWT Token with Missing pubkey Claim", func(t *testing.T) {
		mockToken := "mock_token"
		aHandler.decodeJwt = func(token string) (jwt.MapClaims, error) {
			return jwt.MapClaims{}, nil
		}

		req, err := http.NewRequest("GET", "/refresh_jwt", nil)
		assert.NoError(t, err)
		req.Header.Set("x-jwt", mockToken)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.RefreshToken)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid JWT Token", func(t *testing.T) {
		mockToken := "invalid_token"
		aHandler.decodeJwt = func(token string) (jwt.MapClaims, error) {
			return nil, fmt.Errorf("invalid token")
		}

		req, err := http.NewRequest("GET", "/refresh_jwt", nil)
		assert.NoError(t, err)
		req.Header.Set("x-jwt", mockToken)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.RefreshToken)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Error During JWT Encoding", func(t *testing.T) {
		mockToken := "mock_token"
		person := db.Person{
			Uuid:        uuid.New().String(),
			OwnerPubKey: "your_pubkey",
		}
		db.TestDB.CreateOrEditPerson(person)

		aHandler.decodeJwt = func(token string) (jwt.MapClaims, error) {
			return jwt.MapClaims{"pubkey": person.OwnerPubKey}, nil
		}
		aHandler.encodeJwt = func(pubkey string) (string, error) {
			return "", fmt.Errorf("encoding error")
		}

		req, err := http.NewRequest("GET", "/refresh_jwt", nil)
		assert.NoError(t, err)
		req.Header.Set("x-jwt", mockToken)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(aHandler.RefreshToken)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

}

func TestCreateConnectionCode(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	aHandler := NewAuthHandler(db.TestDB)

	aHandler.makeConnectionCodeRequest = func(inviter_pubkey string, inviter_route_hint string, msats_amount uint64) string {
		return "22222222222222222"
	}

	tests := []struct {
		name           string
		input          db.InviteBody
		expectedStatus int
		expectedBody   string
		mockDBError    error
	}{
		{
			name:           "Valid Input with Pubkey and RouteHint",
			input:          db.InviteBody{Number: 2, Pubkey: "Test_pubkey", RouteHint: "Test_Route_hint", SatsAmount: 21},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "Codes created successfully"}`,
		},
		{
			name:           "Valid Input without Pubkey and RouteHint",
			input:          db.InviteBody{SatsAmount: 21, Number: 2},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "Codes created successfully"}`,
		},
		{
			name:           "Zero SatsAmount",
			input:          db.InviteBody{SatsAmount: 0, Number: 1},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "Codes created successfully"}`,
		},
		{
			name:           "Maximum Number of Codes",
			input:          db.InviteBody{SatsAmount: 21, Number: 1000},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "Codes created successfully"}`,
		},
		{
			name:           "Invalid JSON Body",
			input:          db.InviteBody{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing RouteHint with Pubkey",
			input:          db.InviteBody{Pubkey: "Test_pubkey", SatsAmount: 21, Number: 1},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Route hint is required when pubkey is provided"}`,
		},
		{
			name:           "Missing Pubkey with RouteHint",
			input:          db.InviteBody{RouteHint: "Test_Route_hint", SatsAmount: 21, Number: 1},
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   `{"error": "Pubkey is required when route hint is provided"}`,
		},
		{
			name:           "Empty Request Body",
			input:          db.InviteBody{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Non-integer Number of Codes",
			input:          db.InviteBody{SatsAmount: 21, Number: 0},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Boundary Test for Number of Codes",
			input:          db.InviteBody{SatsAmount: 21, Number: 1},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "Codes created successfully"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, _ := json.Marshal(tt.input)

			req, err := http.NewRequest(http.MethodPost, "/connectioncodes", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(aHandler.CreateConnectionCode)

			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestReturnUserMap(t *testing.T) {

	createdTimeStr := "2023-10-01T12:00:00Z"
	createdTime, _ := time.Parse(time.RFC3339, createdTimeStr)

	tests := []struct {
		name     string
		input    db.Person
		expected map[string]interface{}
	}{
		{
			name: "Standard Input Test",
			input: db.Person{
				ID:              1,
				Uuid:            "123e4567-e89b-12d3-a456-426614174000",
				Created:         &createdTime,
				OwnerPubKey:     "owner_pubkey_value",
				OwnerAlias:      "owner_alias_value",
				OwnerContactKey: "contact_key_value",
				Img:             "img_url",
				Description:     "A sample description",
				Tags:            []string{"tag1", "tag2"},
				UniqueName:      "unique_name_value",
				Extras:          map[string]interface{}{"key": "value"},
				LastLogin:       123124353534,
				PriceToMeet:     100,
			},
			expected: map[string]interface{}{
				"id":            1,
				"uuid":          "123e4567-e89b-12d3-a456-426614174000",
				"created":       "2023-10-01T12:00:00Z",
				"owner_pubkey":  "owner_pubkey_value",
				"owner_alias":   "owner_alias_value",
				"contact_key":   "contact_key_value",
				"img":           "img_url",
				"description":   "A sample description",
				"tags":          []string{"tag1", "tag2"},
				"unique_name":   "unique_name_value",
				"pubkey":        "owner_pubkey_value",
				"extras":        map[string]interface{}{"key": "value"},
				"last_login":    "2023-10-02T12:00:00Z",
				"price_to_meet": 100,
				"alias":         "owner_alias_value",
				"url":           config.Host,
			},
		},
		{
			name:  "Empty Fields Test",
			input: db.Person{},
			expected: map[string]interface{}{
				"id":            nil,
				"uuid":          "",
				"created":       "",
				"owner_pubkey":  "",
				"owner_alias":   "",
				"contact_key":   "",
				"img":           "",
				"description":   "",
				"tags":          nil,
				"unique_name":   "",
				"pubkey":        "",
				"extras":        nil,
				"last_login":    "",
				"price_to_meet": 0,
				"alias":         "",
				"url":           config.Host,
			},
		},
		{
			name: "Maximum Length Strings Test",
			input: db.Person{
				Uuid:            strings.Repeat("a", 255),
				OwnerPubKey:     strings.Repeat("b", 255),
				OwnerAlias:      strings.Repeat("c", 255),
				OwnerContactKey: strings.Repeat("d", 255),
				Img:             strings.Repeat("e", 255),
				Description:     strings.Repeat("f", 255),
				UniqueName:      strings.Repeat("g", 255),
			},
			expected: map[string]interface{}{
				"uuid":         strings.Repeat("a", 255),
				"owner_pubkey": strings.Repeat("b", 255),
				"owner_alias":  strings.Repeat("c", 255),
				"contact_key":  strings.Repeat("d", 255),
				"img":          strings.Repeat("e", 255),
				"description":  strings.Repeat("f", 255),
				"unique_name":  strings.Repeat("g", 255),
				"pubkey":       strings.Repeat("b", 255),
				"alias":        strings.Repeat("c", 255),
				"url":          config.Host,
			},
		},
		{
			name:  "Nil Input Test",
			input: db.Person{},
			expected: map[string]interface{}{
				"id":            nil,
				"uuid":          "",
				"created":       "",
				"owner_pubkey":  "",
				"owner_alias":   "",
				"contact_key":   "",
				"img":           "",
				"description":   "",
				"tags":          nil,
				"unique_name":   "",
				"pubkey":        "",
				"extras":        nil,
				"last_login":    "",
				"price_to_meet": 0,
				"alias":         "",
				"url":           config.Host,
			},
		},
		{
			name: "Invalid Data Types Test",
			input: db.Person{
				Uuid: "12345",
			},
			expected: map[string]interface{}{
				"uuid": "12345",
				"url":  config.Host,
			},
		},
		{
			name: "Large Number of Tags Test",
			input: db.Person{
				Tags: make([]string, 10000),
			},
			expected: map[string]interface{}{
				"tags": make([]string, 10000),
				"url":  config.Host,
			},
		},
		{
			name: "Special Characters in Strings Test",
			input: db.Person{
				Uuid:        "123e4567-e89b-12d3-a456-426614174000",
				OwnerAlias:  "owner_alias_!@#$%^&*()",
				Description: "Description with special characters: !@#$%^&*()",
			},
			expected: map[string]interface{}{
				"uuid":        "123e4567-e89b-12d3-a456-426614174000",
				"owner_alias": "owner_alias_!@#$%^&*()",
				"description": "Description with special characters: !@#$%^&*()",
				"url":         config.Host,
			},
		},
		{
			name:  "Config Dependency Test",
			input: db.Person{},
			expected: map[string]interface{}{
				"url": config.Host,
			},
		},
		{
			name: "Null Values in Map Test",
			input: db.Person{
				Extras: map[string]interface{}{"key1": nil, "key2": "value"},
			},
			expected: map[string]interface{}{
				"extras": map[string]interface{}{"key1": nil, "key2": "value"},
				"url":    config.Host,
			},
		},
		{
			name: "Negative Price Test",
			input: db.Person{
				PriceToMeet: -50,
			},
			expected: map[string]interface{}{
				"price_to_meet": -50,
				"url":           config.Host,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := returnUserMap(tt.input)
			if expectedUuid, ok := tt.expected["uuid"]; ok {
				assert.Equal(t, expectedUuid, result["uuid"])
			}

		})
	}
}

func TestGetLnurlAuth(t *testing.T) {

	originalStore := db.Store
	defer func() {
		db.Store = originalStore
	}()

	db.InitCache()

	t.Run("Valid Request with Existing Socket Key", func(t *testing.T) {

		existingSocketKey := "existing123"
		db.Store.SetSocketConnections(db.Client{
			Host: existingSocketKey,
			Conn: nil,
		})

		req, err := http.NewRequest("GET", "/lnauth?socketKey="+existingSocketKey, nil)
		assert.NoError(t, err)
		req.Host = "test.com"

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetLnurlAuth)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]string
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response["k1"])
		assert.NotEmpty(t, response["encode"])

		socket, err := db.Store.GetSocketConnections(response["k1"][0:20])
		assert.NoError(t, err)
		assert.Equal(t, response["k1"][0:20], socket.Host)
	})

	t.Run("Valid Request with Non-Existing Socket Key", func(t *testing.T) {
		nonExistingKey := "nonexistent456"
		req, err := http.NewRequest("GET", "/lnauth?socketKey="+nonExistingKey, nil)
		assert.NoError(t, err)
		req.Host = "test.com"

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetLnurlAuth)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]string
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response["k1"])
		assert.NotEmpty(t, response["encode"])

		socket, err := db.Store.GetSocketConnections(response["k1"][0:20])
		assert.NoError(t, err)
		assert.Equal(t, response["k1"][0:20], socket.Host)
	})

	t.Run("Handles missing socketKey parameter", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/lnauth", nil)
		assert.NoError(t, err)
		req.Host = "test.com"

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetLnurlAuth)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]string
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.NotEmpty(t, response["k1"])
		assert.NotEmpty(t, response["encode"])
	})

	t.Run("Successfully generates LNURL AUTH", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/lnauth?socketKey=test123", nil)
		assert.NoError(t, err)

		req.Host = "test.com"

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetLnurlAuth)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]string
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.NotEmpty(t, response["k1"])
		assert.NotEmpty(t, response["encode"])

		lnStore, err := db.Store.GetLnCache(response["k1"])
		assert.NoError(t, err)
		assert.Equal(t, response["k1"], lnStore.K1)
		assert.Empty(t, lnStore.Key)
		assert.False(t, lnStore.Status)

		socket, err := db.Store.GetSocketConnections(response["k1"][0:20])
		assert.NoError(t, err)
		assert.Equal(t, response["k1"][0:20], socket.Host)
	})

	t.Run("Handles empty server host", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/lnauth?socketKey=test123", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetLnurlAuth)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]string
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		assert.NotEmpty(t, response["k1"])
		assert.NotEmpty(t, response["encode"])
	})

	t.Run("Verifies cache storage", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/lnauth?socketKey=test123", nil)
		assert.NoError(t, err)
		req.Host = "test.com"

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetLnurlAuth)
		handler.ServeHTTP(rr, req)

		var response map[string]string
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)

		lnStore, err := db.Store.GetLnCache(response["k1"])
		assert.NoError(t, err)
		assert.Equal(t, response["k1"], lnStore.K1)
		assert.Empty(t, lnStore.Key)
		assert.False(t, lnStore.Status)

		socket, err := db.Store.GetSocketConnections(response["k1"][0:20])
		assert.NoError(t, err)
		assert.Equal(t, response["k1"][0:20], socket.Host)
	})

	t.Run("Empty Socket Key", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/lnauth?socketKey=", nil)
		assert.NoError(t, err)
		req.Host = "test.com"

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetLnurlAuth)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]string
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response["k1"])
		assert.NotEmpty(t, response["encode"])

		socket, err := db.Store.GetSocketConnections(response["k1"][0:20])
		assert.NoError(t, err)
		assert.Equal(t, response["k1"][0:20], socket.Host)
	})

	t.Run("Handles LNURL encode failure", func(t *testing.T) {

		originalEncodeLNURL := auth.EncodeLNURLFunc
		auth.EncodeLNURLFunc = func(host string) (auth.LnEncodeData, error) {
			return auth.LnEncodeData{}, fmt.Errorf("encoding failed")
		}
		defer func() {
			auth.EncodeLNURLFunc = originalEncodeLNURL
		}()

		req, err := http.NewRequest("GET", "/lnauth?socketKey=test123", nil)
		assert.NoError(t, err)
		req.Host = "test.com"

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetLnurlAuth)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var response map[string]string
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Empty(t, response["k1"])
		assert.Empty(t, response["encode"])
	})
}
