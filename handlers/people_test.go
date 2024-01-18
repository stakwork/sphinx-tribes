package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	dbMocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeletePerson(t *testing.T) {

	mockDb := dbMocks.NewDatabase(t)

	mockUser := &db.Person{
		ID:          1,
		OwnerPubKey: "authorized-key",
	}
	mockDb.On("GetPerson", uint(1)).Return(mockUser, nil)

	r := chi.NewRouter()
	r.Delete("/{id}", DeletePerson)

	t.Run("unauthorized deletion due to pubkey mismatch", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		rr := httptest.NewRecorder()
		ctx := context.WithValue(req.Context(), "db", mockDb)
		req = req.WithContext(ctx)
		ctx = context.WithValue(ctx, auth.ContextKey, "unauthorized-key")
		req = req.WithContext(ctx)
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("authorized deletion", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		rr := httptest.NewRecorder()
		ctx := context.WithValue(req.Context(), "db", mockDb)
		req = req.WithContext(ctx)
		ctx = context.WithValue(ctx, auth.ContextKey, "authorized-key")
		req = req.WithContext(ctx)
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	mockDb.AssertExpectations(t)
}
