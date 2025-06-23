package db

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateSSEMessageLog(t *testing.T) {
	InitTestDB()

	tests := []struct {
		name        string
		event       map[string]interface{}
		chatID      string
		from        string
		to          string
		expected    *SSEMessageLog
		expectError bool
	}{
		{
			name: "Valid Message Log Creation",
			event: map[string]interface{}{
				"type":    "message",
				"content": "Hello world",
			},
			chatID: "chat123",
			from:   "https://source.com/sse",
			to:     "https://target.com/webhook",
			expected: &SSEMessageLog{
				Event: map[string]interface{}{
					"type":    "message",
					"content": "Hello world",
				},
				ChatID: "chat123",
				From:   "https://source.com/sse",
				To:     "https://target.com/webhook",
				Status: SSEStatusNew,
			},
			expectError: false,
		},
		{
			name:        "Empty Chat ID",
			event:       map[string]interface{}{"type": "message"},
			chatID:      "",
			from:        "https://source.com/sse",
			to:          "https://target.com/webhook",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Empty Source URL",
			event:       map[string]interface{}{"type": "message"},
			chatID:      "chat123",
			from:        "",
			to:          "https://target.com/webhook",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Empty Target URL",
			event:       map[string]interface{}{"type": "message"},
			chatID:      "chat123",
			from:        "https://source.com/sse",
			to:          "",
			expected:    nil,
			expectError: true,
		},
		{
			name:  "Nil Event",
			event: nil,
			chatID: "chat123",
			from:   "https://source.com/sse",
			to:     "https://target.com/webhook",
			expected: &SSEMessageLog{
				Event:  nil,
				ChatID: "chat123",
				From:   "https://source.com/sse",
				To:     "https://target.com/webhook",
				Status: SSEStatusNew,
			},
			expectError: false,
		},
		{
			name: "Complex Event Object",
			event: map[string]interface{}{
				"type": "message",
				"content": map[string]interface{}{
					"text":     "Hello world",
					"metadata": []interface{}{"tag1", "tag2"},
					"nested": map[string]interface{}{
						"key": "value",
					},
				},
			},
			chatID: "chat123",
			from:   "https://source.com/sse",
			to:     "https://target.com/webhook",
			expected: &SSEMessageLog{
				Event: map[string]interface{}{
					"type": "message",
					"content": map[string]interface{}{
						"text":     "Hello world",
						"metadata": []interface{}{"tag1", "tag2"},
						"nested": map[string]interface{}{
							"key": "value",
						},
					},
				},
				ChatID: "chat123",
				From:   "https://source.com/sse",
				To:     "https://target.com/webhook",
				Status: SSEStatusNew,
			},
			expectError: false,
		},
		{
			name: "Special Characters in URLs",
			event: map[string]interface{}{
				"type": "message",
			},
			chatID: "chat123",
			from:   "https://source.com/sse?param=value&special=!@#$%^",
			to:     "https://target.com/webhook?callback=true&id=123",
			expected: &SSEMessageLog{
				Event: map[string]interface{}{
					"type": "message",
				},
				ChatID: "chat123",
				From:   "https://source.com/sse?param=value&special=!@#$%^",
				To:     "https://target.com/webhook?callback=true&id=123",
				Status: SSEStatusNew,
			},
			expectError: false,
		},
		{
			name: "Unicode Characters in Chat ID",
			event: map[string]interface{}{
				"type": "message",
			},
			chatID: "chat-你好-123",
			from:   "https://source.com/sse",
			to:     "https://target.com/webhook",
			expected: &SSEMessageLog{
				Event: map[string]interface{}{
					"type": "message",
				},
				ChatID: "chat-你好-123",
				From:   "https://source.com/sse",
				To:     "https://target.com/webhook",
				Status: SSEStatusNew,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestDB.db.Exec("DELETE FROM sse_message_logs")

			messageLog, err := TestDB.CreateSSEMessageLog(tt.event, tt.chatID, tt.from, tt.to)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, messageLog)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, messageLog)
				assert.NotEqual(t, uuid.Nil, messageLog.ID)
				assert.Equal(t, tt.expected.ChatID, messageLog.ChatID)
				assert.Equal(t, tt.expected.From, messageLog.From)
				assert.Equal(t, tt.expected.To, messageLog.To)
				assert.Equal(t, tt.expected.Status, messageLog.Status)
				assert.WithinDuration(t, time.Now(), messageLog.CreatedAt, 2*time.Second)
				assert.WithinDuration(t, time.Now(), messageLog.UpdatedAt, 2*time.Second)
			}
		})
	}
}

func TestDeleteSSEMessageLog(t *testing.T) {
	InitTestDB()

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		expectError bool
	}{
		{
			name: "Delete Existing Message Log",
			setup: func() uuid.UUID {
				messageLog := &SSEMessageLog{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Event:     map[string]interface{}{"type": "message"},
					ChatID:    "chat123",
					From:      "https://source.com/sse",
					To:        "https://target.com/webhook",
					Status:    SSEStatusNew,
				}
				TestDB.db.Create(messageLog)
				return messageLog.ID
			},
			expectError: false,
		},
		{
			name: "Delete Non-Existent Message Log",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			expectError: true,
		},
		{
			name: "Delete with Nil UUID",
			setup: func() uuid.UUID {
				return uuid.Nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestDB.db.Exec("DELETE FROM sse_message_logs")

			id := tt.setup()
			err := TestDB.DeleteSSEMessageLog(id)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				var count int64
				TestDB.db.Model(&SSEMessageLog{}).Where("id = ?", id).Count(&count)
				assert.Equal(t, int64(0), count)
			}
		})
	}
}

func TestUpdateSSEMessageLogStatusBatch(t *testing.T) {
	InitTestDB()

	tests := []struct {
		name        string
		setup       func() []uuid.UUID
		expectError bool
	}{
		{
			name: "Update Multiple Message Logs",
			setup: func() []uuid.UUID {
				ids := make([]uuid.UUID, 3)
				for i := 0; i < 3; i++ {
					messageLog := &SSEMessageLog{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Event:     map[string]interface{}{"type": "message"},
						ChatID:    "chat123",
						From:      "https://source.com/sse",
						To:        "https://target.com/webhook",
						Status:    SSEStatusNew,
					}
					TestDB.db.Create(messageLog)
					ids[i] = messageLog.ID
				}
				return ids
			},
			expectError: false,
		},
		{
			name: "Update Empty ID List",
			setup: func() []uuid.UUID {
				return []uuid.UUID{}
			},
			expectError: true,
		},
		{
			name: "Update Non-Existent IDs",
			setup: func() []uuid.UUID {
				return []uuid.UUID{uuid.New(), uuid.New()}
			},
			expectError: true,
		},
		{
			name: "Update Mix of Existing and Non-Existing IDs",
			setup: func() []uuid.UUID {
				ids := make([]uuid.UUID, 2)
				messageLog := &SSEMessageLog{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Event:     map[string]interface{}{"type": "message"},
					ChatID:    "chat123",
					From:      "https://source.com/sse",
					To:        "https://target.com/webhook",
					Status:    SSEStatusNew,
				}
				TestDB.db.Create(messageLog)
				ids[0] = messageLog.ID
				ids[1] = uuid.New()
				return ids
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestDB.db.Exec("DELETE FROM sse_message_logs")

			ids := tt.setup()
			err := TestDB.UpdateSSEMessageLogStatusBatch(ids)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				var logs []SSEMessageLog
				TestDB.db.Where("id IN ?", ids).Find(&logs)

				for _, log := range logs {
					assert.Equal(t, SSEStatusSent, log.Status)
				}
			}
		})
	}
}

func TestUpdateSSEMessageLog(t *testing.T) {
	InitTestDB()

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		updates     map[string]interface{}
		verify      func(t *testing.T, log *SSEMessageLog)
		expectError bool
	}{
		{
			name: "Update Status Only",
			setup: func() uuid.UUID {
				messageLog := &SSEMessageLog{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Event:     PropertyMap{"type": "message"},
					ChatID:    "chat123",
					From:      "https://source.com/sse",
					To:        "https://target.com/webhook",
					Status:    SSEStatusNew,
				}
				TestDB.db.Create(messageLog)
				return messageLog.ID
			},
			updates: map[string]interface{}{
				"status": SSEStatusSent,
			},
			verify: func(t *testing.T, log *SSEMessageLog) {
				assert.Equal(t, SSEStatusSent, log.Status)
			},
			expectError: false,
		},
		{
			name: "Update From and To URLs",
			setup: func() uuid.UUID {
				messageLog := &SSEMessageLog{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Event:     PropertyMap{"type": "message"},
					ChatID:    "chat123",
					From:      "https://source.com/sse",
					To:        "https://target.com/webhook",
					Status:    SSEStatusNew,
				}
				TestDB.db.Create(messageLog)
				return messageLog.ID
			},
			updates: map[string]interface{}{
				"from": "https://newsource.com/sse",
				"to":   "https://newtarget.com/webhook",
			},
			verify: func(t *testing.T, log *SSEMessageLog) {
				assert.Equal(t, "https://newsource.com/sse", log.From)
				assert.Equal(t, "https://newtarget.com/webhook", log.To)
			},
			expectError: false,
		},
		{
			name: "Update Multiple Fields",
			setup: func() uuid.UUID {
				messageLog := &SSEMessageLog{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Event:     PropertyMap{"type": "message"},
					ChatID:    "chat123",
					From:      "https://source.com/sse",
					To:        "https://target.com/webhook",
					Status:    SSEStatusNew,
				}
				TestDB.db.Create(messageLog)
				return messageLog.ID
			},
			updates: map[string]interface{}{
				"status": SSEStatusSent,
				"from":   "https://newsource.com/sse",
				"to":     "https://newtarget.com/webhook",
			},
			verify: func(t *testing.T, log *SSEMessageLog) {
				assert.Equal(t, SSEStatusSent, log.Status)
				assert.Equal(t, "https://newsource.com/sse", log.From)
				assert.Equal(t, "https://newtarget.com/webhook", log.To)
			},
			expectError: false,
		},
		{
			name: "Non-Existent ID",
			setup: func() uuid.UUID {
				return uuid.New() 
			},
			updates: map[string]interface{}{
				"status": SSEStatusSent,
			},
			verify:      func(t *testing.T, log *SSEMessageLog) {},
			expectError: true,
		},
		{
			name: "Empty Updates",
			setup: func() uuid.UUID {
				messageLog := &SSEMessageLog{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Event:     PropertyMap{"type": "message"},
					ChatID:    "chat123",
					From:      "https://source.com/sse",
					To:        "https://target.com/webhook",
					Status:    SSEStatusNew,
				}
				TestDB.db.Create(messageLog)
				return messageLog.ID
			},
			updates:     map[string]interface{}{},
			verify:      func(t *testing.T, log *SSEMessageLog) {},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestDB.db.Exec("DELETE FROM sse_message_logs")

			id := tt.setup()
			updatedLog, err := TestDB.UpdateSSEMessageLog(id, tt.updates)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, updatedLog)
				
				tt.verify(t, updatedLog)
				
				var originalLog SSEMessageLog
				TestDB.db.First(&originalLog, "id = ?", id)
				assert.True(t, updatedLog.UpdatedAt.After(originalLog.CreatedAt))
			}
		})
	}
}

func TestGetSSEMessageLogByID(t *testing.T) {
	InitTestDB()

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		expectError bool
	}{
		{
			name: "Get Existing Message Log",
			setup: func() uuid.UUID {
				messageLog := &SSEMessageLog{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Event:     map[string]interface{}{"type": "message"},
					ChatID:    "chat123",
					From:      "https://source.com/sse",
					To:        "https://target.com/webhook",
					Status:    SSEStatusNew,
				}
				TestDB.db.Create(messageLog)
				return messageLog.ID
			},
			expectError: false,
		},
		{
			name: "Get Non-Existent Message Log",
			setup: func() uuid.UUID {
				return uuid.New() 
			},
			expectError: true,
		},
		{
			name: "Get with Nil UUID",
			setup: func() uuid.UUID {
				return uuid.Nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestDB.db.Exec("DELETE FROM sse_message_logs")

			id := tt.setup()
			messageLog, err := TestDB.GetSSEMessageLogByID(id)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, messageLog)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, messageLog)
				assert.Equal(t, id, messageLog.ID)
			}
		})
	}
}

func TestGetSSEMessageLogsByChatID(t *testing.T) {
	InitTestDB()

	tests := []struct {
		name        string
		setup       func() string
		expectedLen int
		expectError bool
	}{
		{
			name: "Get Multiple Message Logs",
			setup: func() string {
				chatID := "chat123"
				for i := 0; i < 3; i++ {
					messageLog := &SSEMessageLog{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Event:     map[string]interface{}{"type": "message"},
						ChatID:    chatID,
						From:      "https://source.com/sse",
						To:        "https://target.com/webhook",
						Status:    SSEStatusNew,
					}
					TestDB.db.Create(messageLog)
				}
				return chatID
			},
			expectedLen: 3,
			expectError: false,
		},
		{
			name: "Get No Message Logs",
			setup: func() string {
				return "nonexistent-chat"
			},
			expectedLen: 0,
			expectError: false,
		},
		{
			name: "Empty Chat ID",
			setup: func() string {
				return ""
			},
			expectedLen: 0,
			expectError: true,
		},
		{
			name: "Mixed Status Message Logs",
			setup: func() string {
				chatID := "mixed-status-chat"
				statuses := []SSEMessageStatus{SSEStatusNew, SSEStatusSent, SSEStatusNew}
				for i := 0; i < 3; i++ {
					messageLog := &SSEMessageLog{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Event:     map[string]interface{}{"type": "message"},
						ChatID:    chatID,
						From:      "https://source.com/sse",
						To:        "https://target.com/webhook",
						Status:    statuses[i],
					}
					TestDB.db.Create(messageLog)
				}
				return chatID
			},
			expectedLen: 3,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestDB.db.Exec("DELETE FROM sse_message_logs")

			chatID := tt.setup()
			messageLogs, err := TestDB.GetSSEMessageLogsByChatID(chatID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedLen, len(messageLogs))
				
				if len(messageLogs) > 0 {
					for _, log := range messageLogs {
						assert.Equal(t, chatID, log.ChatID)
					}
					
					for i := 1; i < len(messageLogs); i++ {
						assert.True(t, !messageLogs[i].CreatedAt.After(messageLogs[i-1].CreatedAt))
					}
				}
			}
		})
	}
}

func TestGetNewSSEMessageLogsByChatID(t *testing.T) {
	InitTestDB()

	tests := []struct {
		name        string
		setup       func() string
		expectedLen int
		expectError bool
	}{
		{
			name: "Get Only New Message Logs",
			setup: func() string {
				chatID := "chat-with-new"
				statuses := []SSEMessageStatus{SSEStatusNew, SSEStatusSent, SSEStatusNew}
				for i := 0; i < 3; i++ {
					messageLog := &SSEMessageLog{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Event:     map[string]interface{}{"type": "message"},
						ChatID:    chatID,
						From:      "https://source.com/sse",
						To:        "https://target.com/webhook",
						Status:    statuses[i],
					}
					TestDB.db.Create(messageLog)
				}
				return chatID
			},
			expectedLen: 2,
			expectError: false,
		},
		{
			name: "No New Message Logs",
			setup: func() string {
				chatID := "chat-no-new"
				for i := 0; i < 2; i++ {
					messageLog := &SSEMessageLog{
						ID:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Event:     map[string]interface{}{"type": "message"},
						ChatID:    chatID,
						From:      "https://source.com/sse",
						To:        "https://target.com/webhook",
						Status:    SSEStatusSent,
					}
					TestDB.db.Create(messageLog)
				}
				return chatID
			},
			expectedLen: 0,
			expectError: false,
		},
		{
			name: "Empty Chat ID",
			setup: func() string {
				return ""
			},
			expectedLen: 0,
			expectError: true,
		},
		{
			name: "Non-Existent Chat ID",
			setup: func() string {
				return "nonexistent-chat"
			},
			expectedLen: 0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestDB.db.Exec("DELETE FROM sse_message_logs")

			chatID := tt.setup()
			messageLogs, err := TestDB.GetNewSSEMessageLogsByChatID(chatID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedLen, len(messageLogs))
				
				for _, log := range messageLogs {
					assert.Equal(t, chatID, log.ChatID)
					assert.Equal(t, SSEStatusNew, log.Status)
				}
				
				for i := 1; i < len(messageLogs); i++ {
					assert.True(t, !messageLogs[i].CreatedAt.After(messageLogs[i-1].CreatedAt))
				}
			}
		})
	}
} 