package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	mocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateChannel(t *testing.T) {
	mockDb := mocks.NewDatabase(t)
	cHandler := NewChannelHandler(mockDb)

	// Mock data for testing
	mockPubKey := "mock_pubkey"
	mockTribeUUID := "mock_tribe_uuid"
	mockChannelName := "mock_channel"
	mockRequestBody := map[string]interface{}{
		"tribe_uuid": mockTribeUUID,
		"name":       mockChannelName,
	}

	// Mock request body
	requestBodyBytes, err := json.Marshal(mockRequestBody)
	assert.NoError(t, err)

	t.Run("Should test that a user that is not authenticated cannot create a channel", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/channel", bytes.NewBuffer(requestBodyBytes))
		assert.NoError(t, err)
		rr := httptest.NewRecorder()

		mockDb.On("GetTribe", mockTribeUUID).Return(db.Tribe{OwnerPubKey: mockPubKey})

		cHandler.CreateChannel(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that an authenticated user can create a channel", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/channel", bytes.NewBuffer(requestBodyBytes))
		assert.NoError(t, err)
		req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, mockPubKey))
		rr := httptest.NewRecorder()

		mockDb.On("GetTribe", mockTribeUUID).Return(db.Tribe{OwnerPubKey: mockPubKey})
		mockDb.On("GetChannelsByTribe", mockTribeUUID).Return([]db.Channel{})
		mockDb.On("CreateChannel", mock.Anything).Return(db.Channel{}, nil)

		cHandler.CreateChannel(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should test that a user cannot create a channel with a name that already exists", func(t *testing.T) {
		mockDb.ExpectedCalls = nil

		req, err := http.NewRequest("POST", "/channel", bytes.NewBuffer(requestBodyBytes))
		assert.NoError(t, err)
		req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, mockPubKey))
		rr := httptest.NewRecorder()

		mockDb.On("GetTribe", mockTribeUUID).Return(db.Tribe{OwnerPubKey: mockPubKey})
		mockDb.On("GetChannelsByTribe", mockTribeUUID).Return([]db.Channel{{Name: mockChannelName}})

		cHandler.CreateChannel(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)

		// Ensure that the expected methods were called
		mockDb.AssertExpectations(t)
	})
}

func TestDeleteChannel(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "mock_pubkey")
	mockDb := mocks.NewDatabase(t)
	cHandler := NewChannelHandler(mockDb)

	// Mock data for testing
	mockPubKey := "mock_pubkey"
	mockChannelID := uint(1)

	t.Run("Should test that the owner of a channel can delete the channel", func(t *testing.T) {
		mockDb.On("GetChannel", mockChannelID).Return(db.Channel{ID: mockChannelID, TribeUUID: "mock_tribe_uuid"})
		mockDb.On("GetTribe", "mock_tribe_uuid").Return(db.Tribe{OwnerPubKey: mockPubKey})
		mockDb.On("UpdateChannel", mockChannelID, mock.Anything).Return(true)

		// Create and Serve request
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(ctx, "DELETE", "/channel/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		handler := http.HandlerFunc(cHandler.DeleteChannel)
		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should test that non-channel owners cannot delete the channel, it should return a 401 error", func(t *testing.T) {
		mockPubKey := "other_pubkey"

		mockDb.ExpectedCalls = nil
		mockDb.On("GetChannel", mockChannelID).Return(db.Channel{ID: mockChannelID, TribeUUID: "mock_tribe_uuid"})
		mockDb.On("GetTribe", "mock_tribe_uuid").Return(db.Tribe{OwnerPubKey: mockPubKey})

		// Create and Serve request
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(ctx, "DELETE", "/channel/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		handler := http.HandlerFunc(cHandler.DeleteChannel)
		handler.ServeHTTP(rr, req)

		// Verify response
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
