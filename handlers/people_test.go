package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/go-chi/chi"
	dbMocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeletePerson(t *testing.T) {
    mockDb := dbMocks.NewDatabase(t)
    ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")

    t.Run("unauthorized deletion due to pubkey mismatch", func(t *testing.T) {
        rr := httptest.NewRecorder()

        // Mocking the database response
        mockUser := db.Person{
            ID: 1,
            OwnerPubKey: "authorized-key",
        }
        mockDb.On("GetPerson", uint(1)).Return(mockUser)

        // Creating the request with the appropriate context
        rctx := chi.NewRouteContext()
        rctx.URLParams.Add("id", "1")
        req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/", nil)
        if err != nil {
            t.Fatal(err)
        }

        // Calling the handler function
        DeletePerson(rr, req)

        // Asserting the response
        assert.Equal(t, http.StatusUnauthorized, rr.Code)
    })
}

