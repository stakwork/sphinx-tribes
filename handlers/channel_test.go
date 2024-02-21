package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	mocks "github.com/stakwork/sphinx-tribes/mocks"
)

func TestCreateChannel(t *testing.T) {
	mockDB := new(mocks.Database)
	handler := NewChannelHandler(mockDB)

	authPubKey := "authPubKey"
	ctx := context.WithValue(context.Background(), auth.ContextKey, authPubKey)

	tribeUUID := uuid.New().String()
	mockChannel := db.Channel{Name: "TestChannel", TribeUUID: tribeUUID}
	mockTribe := db.Tribe{UUID: tribeUUID, OwnerPubKey: authPubKey}

	mockDB.On("GetTribe", tribeUUID).Return(mockTribe, nil)
	mockDB.On("CreateChannel", mock.AnythingOfType("db.Channel")).Return(mockChannel, nil)

	channelData, _ := json.Marshal(mockChannel)
	req, err := http.NewRequest("POST", "/channel", bytes.NewBuffer(channelData))
	assert.NoError(t, err)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/channel", handler.CreateChannel)

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var responseChannel db.Channel
	err = json.Unmarshal(rr.Body.Bytes(), &responseChannel)
	assert.NoError(t, err)
	assert.Equal(t, mockChannel.Name, responseChannel.Name, "Channel names should match")

	mockDB.AssertExpectations(t)
}
