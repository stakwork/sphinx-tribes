package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	mocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUnitCreateOrEditOrganization(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	oHandler := NewOrganizationHandler(mockDb)

	t.Run("should return error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditOrganization)

		invalidJson := []byte(`{"key": "value"`)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("should return error if public key not present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditOrganization)

		invalidJson := []byte(`{"name": "value"}`)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return error org name is empty", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditOrganization)

		invalidJson := []byte(`{"name": ""}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return error org name is more than 20", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditOrganization)

		invalidJson := []byte(`{"name": "DemoTestingOrganization"}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should successfully add organization if request is valid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.CreateOrEditOrganization)

		mockDb.On("GetOrganizationByUuid", mock.AnythingOfType("string")).Return(db.Organization{}).Once()
		mockDb.On("GetOrganizationByName", "TestOrganization").Return(db.Organization{}).Once()
		mockDb.On("CreateOrEditOrganization", mock.MatchedBy(func(org db.Organization) bool {
			return org.Name == "TestOrganization" && org.Uuid != "" && org.Updated != nil && org.Created != nil
		})).Return(db.Organization{}, nil).Once()

		validJson := []byte(`{"name": "TestOrganization", "owner_pubkey": "test-key"}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(validJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
