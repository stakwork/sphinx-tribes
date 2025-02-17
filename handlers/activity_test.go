package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
)

const validHumanAuthorRef = "02abc123456789abcdef0123456789abcdef0123"

func TestReceiveActivity(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	ah := NewActivityHandler(&http.Client{}, db.TestDB)

	tests := []struct {
		name           string
		payload        WebhookActivityRequest
		expectedStatus int
		expectedError  string
		validateFunc   func(t *testing.T, resp WebhookResponse)
	}{
		{
			name: "successful new activity creation",
			payload: WebhookActivityRequest{
				ContentType: "general_update",
				Content:     "Test content",
				Workspace:   "test-workspace",
				Author:      db.HumansAuthor,
				AuthorRef:   validHumanAuthorRef,
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp WebhookResponse) {
				assert.True(t, resp.Success)
				assert.NotEmpty(t, resp.ActivityID)

				activity, err := db.TestDB.GetActivity(resp.ActivityID)
				assert.NoError(t, err)
				assert.Equal(t, "Test content", activity.Content)
				assert.Equal(t, "test-workspace", activity.Workspace)
			},
		},
		{
			name: "successful thread activity creation",
			payload: WebhookActivityRequest{
				ContentType: "general_update",
				Content:     "Thread reply",
				Workspace:   "test-workspace",
				ThreadID:    uuid.New().String(),
				Author:      db.HumansAuthor,
				AuthorRef:   validHumanAuthorRef,
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp WebhookResponse) {
				assert.True(t, resp.Success)
				assert.NotEmpty(t, resp.ActivityID)

				activity, err := db.TestDB.GetActivity(resp.ActivityID)
				assert.NoError(t, err)
				assert.Equal(t, "Thread reply", activity.Content)
				assert.Equal(t, "test-workspace", activity.Workspace)
			},
		},
		{
			name: "invalid content - empty",
			payload: WebhookActivityRequest{
				ContentType: "general_update",
				Content:     "",
				Workspace:   "test-workspace",
				Author:      db.HumansAuthor,
				AuthorRef:   validHumanAuthorRef,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "content must not be empty and must be less than 10000 characters",
		},
		{
			name: "invalid author type",
			payload: WebhookActivityRequest{
				ContentType: "general_update",
				Content:     "Test content",
				Workspace:   "test-workspace",
				Author:      "invalid",
				AuthorRef:   validHumanAuthorRef,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid author type",
		},
		{
			name: "invalid human author ref",
			payload: WebhookActivityRequest{
				ContentType: "general_update",
				Content:     "Test content",
				Workspace:   "test-workspace",
				Author:      db.HumansAuthor,
				AuthorRef:   "short",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid public key format for human author",
		},
		{
			name: "successful hive activity with URL",
			payload: WebhookActivityRequest{
				ContentType: "general_update",
				Content:     "Test content",
				Workspace:   "test-workspace",
				Author:      db.HiveAuthor,
				AuthorRef:   "https://example.com/hive",
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, resp WebhookResponse) {
				assert.True(t, resp.Success)
				assert.NotEmpty(t, resp.ActivityID)

				activity, err := db.TestDB.GetActivity(resp.ActivityID)
				assert.NoError(t, err)
				assert.Equal(t, "Test content", activity.Content)
				assert.Equal(t, "https://example.com/hive", activity.AuthorRef)
			},
		},
		{
			name: "invalid thread ID format",
			payload: WebhookActivityRequest{
				ContentType: "general_update",
				Content:     "Test content",
				Workspace:   "test-workspace",
				ThreadID:    "not-a-uuid",
				Author:      db.HumansAuthor,
				AuthorRef:   validHumanAuthorRef,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid source ID format",
		},
		{
			name: "missing workspace",
			payload: WebhookActivityRequest{
				ContentType: "general_update",
				Content:     "Test content",
				Author:      db.HumansAuthor,
				AuthorRef:   validHumanAuthorRef,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "workspace is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/activities/receive", bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			ah.ReceiveActivity(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response WebhookResponse
			err = json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			if tt.expectedError != "" {
				assert.False(t, response.Success)
				assert.Equal(t, tt.expectedError, response.Error)
			} else if tt.validateFunc != nil {
				tt.validateFunc(t, response)
			}
		})
	}
}

func TestReceiveActivity_InvalidJSON(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	ah := NewActivityHandler(&http.Client{}, db.TestDB)

	req := httptest.NewRequest(http.MethodPost, "/activities/receive", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ah.ReceiveActivity(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response WebhookResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Invalid request payload", response.Error)
}
