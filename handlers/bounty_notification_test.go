package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContactKeyAvailable(t *testing.T) {
	empty := ""
	blank := "   "
	contactKey := "contact-key"

	assert.False(t, contactKeyAvailable(nil))
	assert.False(t, contactKeyAvailable(&empty))
	assert.False(t, contactKeyAvailable(&blank))
	assert.True(t, contactKeyAvailable(&contactKey))
}

func TestGetContactKeyHandlesEmptyContactKey(t *testing.T) {
	oldURL := config.V2BotUrl
	oldToken := config.V2BotToken
	defer func() {
		config.V2BotUrl = oldURL
		config.V2BotToken = oldToken
	}()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/get_contact/test-pubkey", r.URL.Path)
		assert.Equal(t, "test-token", r.Header.Get("x-admin-token"))
		_ = json.NewEncoder(w).Encode(map[string]string{"contact_key": ""})
	}))
	defer ts.Close()

	config.V2BotUrl = ts.URL
	config.V2BotToken = "test-token"

	contactKey, err := getContactKey("test-pubkey")
	require.NoError(t, err)
	require.NotNil(t, contactKey)
	assert.False(t, contactKeyAvailable(contactKey))
}

func TestGetContactKeyRejectsNonSuccessResponse(t *testing.T) {
	oldURL := config.V2BotUrl
	oldToken := config.V2BotToken
	defer func() {
		config.V2BotUrl = oldURL
		config.V2BotToken = oldToken
	}()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	config.V2BotUrl = ts.URL
	config.V2BotToken = "test-token"

	contactKey, err := getContactKey("test-pubkey")
	require.Error(t, err)
	assert.Nil(t, contactKey)
}

func TestSendNotificationHandlesV2BotResponses(t *testing.T) {
	oldURL := config.V2BotUrl
	oldToken := config.V2BotToken
	defer func() {
		config.V2BotUrl = oldURL
		config.V2BotToken = oldToken
	}()

	tests := []struct {
		name       string
		statusCode int
		response   map[string]string
		expected   string
	}{
		{
			name:       "returns send status",
			statusCode: http.StatusOK,
			response:   map[string]string{"status": string(db.NotificationStatusComplete)},
			expected:   string(db.NotificationStatusComplete),
		},
		{
			name:       "maps missing status to failed",
			statusCode: http.StatusOK,
			response:   map[string]string{},
			expected:   string(db.NotificationStatusFailed),
		},
		{
			name:       "maps non success response to failed",
			statusCode: http.StatusBadGateway,
			response:   map[string]string{"status": string(db.NotificationStatusComplete)},
			expected:   string(db.NotificationStatusFailed),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/send", r.URL.Path)
				assert.Equal(t, "test-token", r.Header.Get("x-admin-token"))
				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(tt.response)
			}))
			defer ts.Close()

			config.V2BotUrl = ts.URL
			config.V2BotToken = "test-token"

			assert.Equal(t, tt.expected, sendNotification("test-pubkey", "hello"))
		})
	}
}
