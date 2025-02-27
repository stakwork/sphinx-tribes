package db

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateArtifact(t *testing.T) {
	InitTestDB()

	tests := []struct {
		name        string
		input       *Artifact
		setup       func()
		expected    *Artifact
		expectError bool
	}{
		{
			name: "Successfully create text artifact",
			input: &Artifact{
				MessageID: "msg123",
				Type:      TextArtifact,
				Content: PropertyMap{
					"text_type": "code",
					"content":   "print('hello world')",
				},
			},
			setup: func() {
				TestDB.db.Create(&ChatMessage{
					ID:      "msg123",
					ChatID:  "chat123",
					Message: "Test Message",
				})
			},
			expected: &Artifact{
				MessageID: "msg123",
				Type:      TextArtifact,
				Content: PropertyMap{
					"text_type": "code",
					"content":   "print('hello world')",
				},
			},
			expectError: false,
		},
		{
			name: "Successfully create visual artifact",
			input: &Artifact{
				MessageID: "msg456",
				Type:      VisualArtifact,
				Content: PropertyMap{
					"text_type": "img",
					"url":       "https://example.com/image.png",
				},
			},
			setup: func() {
				TestDB.db.Create(&ChatMessage{
					ID:      "msg456",
					ChatID:  "chat456",
					Message: "Test Message",
				})
			},
			expected: &Artifact{
				MessageID: "msg456",
				Type:      VisualArtifact,
				Content: PropertyMap{
					"text_type": "img",
					"url":       "https://example.com/image.png",
				},
			},
			expectError: false,
		},
		{
			name: "Successfully create action artifact",
			input: &Artifact{
				MessageID: "msg789",
				Type:      ActionArtifact,
				Content: PropertyMap{
					"action_text": "Choose an option",
					"options": []interface{}{
						map[string]interface{}{
							"action_type":     "button",
							"option_label":    "Click me",
							"option_response": "clicked",
						},
					},
				},
			},
			setup: func() {
				TestDB.db.Create(&ChatMessage{
					ID:      "msg789",
					ChatID:  "chat789",
					Message: "Test Message",
				})
			},
			expected: &Artifact{
				MessageID: "msg789",
				Type:      ActionArtifact,
				Content: PropertyMap{
					"action_text": "Choose an option",
					"options": []interface{}{
						map[string]interface{}{
							"action_type":     "button",
							"option_label":    "Click me",
							"option_response": "clicked",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Create artifact with non-existent message ID",
			input: &Artifact{
				MessageID: "nonexistent",
				Type:      TextArtifact,
				Content: PropertyMap{
					"text_type": "code",
					"content":   "test",
				},
			},
			setup:       func() {},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Create artifact with empty message ID",
			input: &Artifact{
				MessageID: "",
				Type:      TextArtifact,
				Content: PropertyMap{
					"text_type": "code",
					"content":   "test",
				},
			},
			setup:       func() {},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Create artifact with invalid type",
			input: &Artifact{
				MessageID: "msg123",
				Type:      "invalid_type",
				Content: PropertyMap{
					"text_type": "code",
					"content":   "test",
				},
			},
			setup:       func() {},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Create artifact with empty content",
			input: &Artifact{
				MessageID: "msg123",
				Type:      TextArtifact,
				Content:   PropertyMap{},
			},
			setup: func() {
				TestDB.db.Create(&ChatMessage{
					ID:      "msg123",
					ChatID:  "chat123",
					Message: "Test Message",
				})
			},
			expected: &Artifact{
				MessageID: "msg123",
				Type:      TextArtifact,
				Content:   PropertyMap{},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestDB.db.Exec("DELETE FROM artifacts")
			TestDB.db.Exec("DELETE FROM chat_messages")

			tt.setup()

			result, err := TestDB.CreateArtifact(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEqual(t, uuid.Nil, result.ID)
			assert.Equal(t, tt.expected.MessageID, result.MessageID)
			assert.Equal(t, tt.expected.Type, result.Type)
			assert.Equal(t, tt.expected.Content, result.Content)

			assert.WithinDuration(t, time.Now(), result.CreatedAt, time.Second)
			assert.WithinDuration(t, time.Now(), result.UpdatedAt, time.Second)

			var savedArtifact Artifact
			err = TestDB.db.First(&savedArtifact, "id = ?", result.ID).Error
			assert.NoError(t, err)
			assert.Equal(t, result.ID, savedArtifact.ID)
			assert.Equal(t, tt.expected.MessageID, savedArtifact.MessageID)
			assert.Equal(t, tt.expected.Type, savedArtifact.Type)
			assert.Equal(t, tt.expected.Content, savedArtifact.Content)
		})
	}
}

func TestGetArtifactByID(t *testing.T) {
	InitTestDB()
	currentTime := time.Now()

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		input       uuid.UUID
		expected    *Artifact
		expectError bool
	}{
		{
			name: "Successfully get text artifact",
			setup: func() uuid.UUID {
				id := uuid.New()
				TestDB.db.Create(&ChatMessage{
					ID:      "msg123",
					ChatID:  "chat123",
					Message: "Test Message",
				})
				artifact := &Artifact{
					ID:        id,
					MessageID: "msg123",
					Type:      TextArtifact,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "print('hello world')",
					},
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}
				TestDB.db.Create(artifact)
				return id
			},
			expected: &Artifact{
				MessageID: "msg123",
				Type:      TextArtifact,
				Content: PropertyMap{
					"text_type": "code",
					"content":   "print('hello world')",
				},
			},
			expectError: false,
		},
		{
			name: "Successfully get visual artifact",
			setup: func() uuid.UUID {
				id := uuid.New()
				TestDB.db.Create(&ChatMessage{
					ID:      "msg456",
					ChatID:  "chat456",
					Message: "Test Message",
				})
				artifact := &Artifact{
					ID:        id,
					MessageID: "msg456",
					Type:      VisualArtifact,
					Content: PropertyMap{
						"text_type": "img",
						"url":       "https://example.com/image.png",
					},
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}
				TestDB.db.Create(artifact)
				return id
			},
			expected: &Artifact{
				MessageID: "msg456",
				Type:      VisualArtifact,
				Content: PropertyMap{
					"text_type": "img",
					"url":       "https://example.com/image.png",
				},
			},
			expectError: false,
		},
		{
			name: "Successfully get action artifact",
			setup: func() uuid.UUID {
				id := uuid.New()
				TestDB.db.Create(&ChatMessage{
					ID:      "msg789",
					ChatID:  "chat789",
					Message: "Test Message",
				})
				artifact := &Artifact{
					ID:        id,
					MessageID: "msg789",
					Type:      ActionArtifact,
					Content: PropertyMap{
						"action_text": "Choose an option",
						"options": []interface{}{
							map[string]interface{}{
								"action_type":     "button",
								"option_label":    "Click me",
								"option_response": "clicked",
							},
						},
					},
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}
				TestDB.db.Create(artifact)
				return id
			},
			expected: &Artifact{
				MessageID: "msg789",
				Type:      ActionArtifact,
				Content: PropertyMap{
					"action_text": "Choose an option",
					"options": []interface{}{
						map[string]interface{}{
							"action_type":     "button",
							"option_label":    "Click me",
							"option_response": "clicked",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Non-existent ID",
			setup: func() uuid.UUID {
				return uuid.New() 
			},
			expected:    nil,
			expectError: false, 
		},
		{
			name: "Nil UUID",
			setup: func() uuid.UUID {
				return uuid.Nil
			},
			expected:    nil,
			expectError: false,
		},
		{
			name: "Get deleted artifact",
			setup: func() uuid.UUID {
				id := uuid.New()

				TestDB.db.Create(&ChatMessage{
					ID:      "msg_deleted",
					ChatID:  "chat_deleted",
					Message: "Test Message",
				})
				artifact := &Artifact{
					ID:        id,
					MessageID: "msg_deleted",
					Type:      TextArtifact,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "deleted content",
					},
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}
				TestDB.db.Create(artifact)
				TestDB.db.Delete(artifact) 
				return id
			},
			expected:    nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestDB.db.Exec("DELETE FROM artifacts")
			TestDB.db.Exec("DELETE FROM chat_messages")

			tt.input = tt.setup()

			result, err := TestDB.GetArtifactByID(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.input, result.ID)
			assert.Equal(t, tt.expected.MessageID, result.MessageID)
			assert.Equal(t, tt.expected.Type, result.Type)
			assert.Equal(t, tt.expected.Content, result.Content)
			assert.WithinDuration(t, currentTime, result.CreatedAt, time.Second)
			assert.WithinDuration(t, currentTime, result.UpdatedAt, time.Second)
		})
	}
}

func TestGetArtifactsByMessageID(t *testing.T) {
	InitTestDB()
	currentTime := time.Now()

	tests := []struct {
		name        string
		setup       func() string 
		expected    []Artifact
		expectError bool
	}{
		{
			name: "Successfully get multiple artifacts for message",
			setup: func() string {
				messageID := "msg123"

				TestDB.db.Create(&ChatMessage{
					ID:      messageID,
					ChatID:  "chat123",
					Message: "Test Message",
				})

				artifacts := []Artifact{
					{
						ID:        uuid.New(),
						MessageID: messageID,
						Type:      TextArtifact,
						Content: PropertyMap{
							"text_type": "code",
							"content":   "print('hello')",
						},
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
					{
						ID:        uuid.New(),
						MessageID: messageID,
						Type:      VisualArtifact,
						Content: PropertyMap{
							"text_type": "img",
							"url":       "https://example.com/image.png",
						},
						CreatedAt: currentTime.Add(time.Second),
						UpdatedAt: currentTime.Add(time.Second),
					},
					{
						ID:        uuid.New(),
						MessageID: messageID,
						Type:      ActionArtifact,
						Content: PropertyMap{
							"action_text": "Choose an option",
							"options": []interface{}{
								map[string]interface{}{
									"action_type":     "button",
									"option_label":    "Click me",
									"option_response": "clicked",
								},
							},
						},
						CreatedAt: currentTime.Add(2 * time.Second),
						UpdatedAt: currentTime.Add(2 * time.Second),
					},
				}

				for _, a := range artifacts {
					TestDB.db.Create(&a)
				}

				return messageID
			},
			expected: []Artifact{
				{
					Type:      ActionArtifact,
					MessageID: "msg123",
					Content: PropertyMap{
						"action_text": "Choose an option",
						"options": []interface{}{
							map[string]interface{}{
								"action_type":     "button",
								"option_label":    "Click me",
								"option_response": "clicked",
							},
						},
					},
				},
				{
					Type:      VisualArtifact,
					MessageID: "msg123",
					Content: PropertyMap{
						"text_type": "img",
						"url":       "https://example.com/image.png",
					},
				},
				{
					Type:      TextArtifact,
					MessageID: "msg123",
					Content: PropertyMap{
						"text_type": "code",
						"content":   "print('hello')",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Message with no artifacts",
			setup: func() string {
				messageID := "msg_empty"
				TestDB.db.Create(&ChatMessage{
					ID:      messageID,
					ChatID:  "chat_empty",
					Message: "Message with no artifacts",
				})
				return messageID
			},
			expected:    []Artifact{},
			expectError: false,
		},
		{
			name: "Non-existent message ID",
			setup: func() string {
				return "nonexistent_msg"
			},
			expected:    []Artifact{},
			expectError: false,
		},
		{
			name: "Empty message ID",
			setup: func() string {
				return ""
			},
			expected:    []Artifact{},
			expectError: false,
		},
		{
			name: "Message with soft-deleted artifacts",
			setup: func() string {
				messageID := "msg_deleted"

				TestDB.db.Create(&ChatMessage{
					ID:      messageID,
					ChatID:  "chat_deleted",
					Message: "Test Message",
				})

				artifact := Artifact{
					ID:        uuid.New(),
					MessageID: messageID,
					Type:      TextArtifact,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "deleted content",
					},
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}
				TestDB.db.Create(&artifact)
				TestDB.db.Delete(&artifact)

				return messageID
			},
			expected:    []Artifact{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestDB.db.Exec("DELETE FROM artifacts")
			TestDB.db.Exec("DELETE FROM chat_messages")

			messageID := tt.setup()

			results, err := TestDB.GetArtifactsByMessageID(messageID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, results)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, results, len(tt.expected))

			if len(tt.expected) > 0 {

				for i, expected := range tt.expected {
					assert.Equal(t, expected.MessageID, results[i].MessageID)
					assert.Equal(t, expected.Type, results[i].Type)
					assert.Equal(t, expected.Content, results[i].Content)
					assert.NotEqual(t, uuid.Nil, results[i].ID)
					assert.False(t, results[i].CreatedAt.IsZero())
					assert.False(t, results[i].UpdatedAt.IsZero())
				}

				for i := 0; i < len(results)-1; i++ {
					assert.True(t, results[i].CreatedAt.After(results[i+1].CreatedAt) || 
						results[i].CreatedAt.Equal(results[i+1].CreatedAt))
				}
			}
		})
	}
}

func TestGetAllArtifactsByChatID(t *testing.T) {
	InitTestDB()
	currentTime := time.Now()

	tests := []struct {
		name        string
		setup       func() string 
		expected    []Artifact
		expectError bool
	}{
		{
			name: "Successfully get multiple artifacts across messages",
			setup: func() string {
				chatID := "chat123"

				messages := []ChatMessage{
					{
						ID:      "msg1",
						ChatID:  chatID,
						Message: "First Message",
					},
					{
						ID:      "msg2",
						ChatID:  chatID,
						Message: "Second Message",
					},
				}
				for _, msg := range messages {
					TestDB.db.Create(&msg)
				}

				artifacts1 := []Artifact{
					{
						ID:        uuid.New(),
						MessageID: "msg1",
						Type:      TextArtifact,
						Content: PropertyMap{
							"text_type": "code",
							"content":   "print('first')",
						},
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
					{
						ID:        uuid.New(),
						MessageID: "msg1",
						Type:      VisualArtifact,
						Content: PropertyMap{
							"text_type": "img",
							"url":       "https://example.com/first.png",
						},
						CreatedAt: currentTime.Add(time.Second),
						UpdatedAt: currentTime.Add(time.Second),
					},
				}

				artifacts2 := []Artifact{
					{
						ID:        uuid.New(),
						MessageID: "msg2",
						Type:      ActionArtifact,
						Content: PropertyMap{
							"action_text": "Choose option",
							"options": []interface{}{
								map[string]interface{}{
									"action_type":     "button",
									"option_label":    "Click",
									"option_response": "clicked",
								},
							},
						},
						CreatedAt: currentTime.Add(2 * time.Second),
						UpdatedAt: currentTime.Add(2 * time.Second),
					},
				}

				for _, a := range artifacts1 {
					TestDB.db.Create(&a)
				}
				for _, a := range artifacts2 {
					TestDB.db.Create(&a)
				}

				return chatID
			},
			expected: []Artifact{
				{
					MessageID: "msg2",
					Type:      ActionArtifact,
					Content: PropertyMap{
						"action_text": "Choose option",
						"options": []interface{}{
							map[string]interface{}{
								"action_type":     "button",
								"option_label":    "Click",
								"option_response": "clicked",
							},
						},
					},
				},
				{
					MessageID: "msg1",
					Type:      VisualArtifact,
					Content: PropertyMap{
						"text_type": "img",
						"url":       "https://example.com/first.png",
					},
				},
				{
					MessageID: "msg1",
					Type:      TextArtifact,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "print('first')",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Chat with no messages",
			setup: func() string {
				chatID := "empty_chat"
				return chatID
			},
			expected:    []Artifact{},
			expectError: false,
		},
		{
			name: "Chat with messages but no artifacts",
			setup: func() string {
				chatID := "chat_no_artifacts"
				TestDB.db.Create(&ChatMessage{
					ID:      "msg_empty",
						ChatID:  chatID,
					Message: "Message without artifacts",
				})
				return chatID
			},
			expected:    []Artifact{},
			expectError: false,
		},
		{
			name: "Non-existent chat ID",
			setup: func() string {
				return "nonexistent_chat"
			},
			expected:    []Artifact{},
			expectError: false,
		},
		{
			name: "Empty chat ID",
			setup: func() string {
				return ""
			},
			expected:    []Artifact{},
			expectError: false,
		},
		{
			name: "Chat with soft-deleted artifacts",
			setup: func() string {
				chatID := "chat_deleted"
				TestDB.db.Create(&ChatMessage{
					ID:      "msg_deleted",
						ChatID:  chatID,
					Message: "Message with deleted artifact",
				})

				artifact := Artifact{
					ID:        uuid.New(),
					MessageID: "msg_deleted",
					Type:      TextArtifact,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "deleted content",
					},
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}
				TestDB.db.Create(&artifact)
				TestDB.db.Delete(&artifact)

				return chatID
			},
			expected:    []Artifact{},
			expectError: false,
		},
		{
			name: "Chat with mixed active and deleted artifacts",
			setup: func() string {
				chatID := "chat_mixed"
				TestDB.db.Create(&ChatMessage{
					ID:      "msg_mixed",
						ChatID:  chatID,
					Message: "Message with mixed artifacts",
				})

				active := Artifact{
						ID:        uuid.New(),
					MessageID: "msg_mixed",
						Type:      TextArtifact,
						Content: PropertyMap{
							"text_type": "code",
						"content":   "active content",
						},
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
				}
				TestDB.db.Create(&active)

				deleted := Artifact{
						ID:        uuid.New(),
					MessageID: "msg_mixed",
					Type:      TextArtifact,
						Content: PropertyMap{
						"text_type": "code",
						"content":   "deleted content",
					},
					CreatedAt: currentTime.Add(-time.Second),
					UpdatedAt: currentTime.Add(-time.Second),
				}
				TestDB.db.Create(&deleted)
				TestDB.db.Delete(&deleted)

				return chatID
			},
			expected: []Artifact{
				{
					MessageID: "msg_mixed",
					Type:      TextArtifact,
						Content: PropertyMap{
						"text_type": "code",
						"content":   "active content",
								},
							},
						},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			TestDB.db.Exec("DELETE FROM artifacts")
			TestDB.db.Exec("DELETE FROM chat_messages")

			chatID := tt.setup()

			results, err := TestDB.GetAllArtifactsByChatID(chatID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, results)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, results, len(tt.expected))

			if len(tt.expected) > 0 {

				for i, expected := range tt.expected {
					assert.Equal(t, expected.MessageID, results[i].MessageID)
					assert.Equal(t, expected.Type, results[i].Type)
					assert.Equal(t, expected.Content, results[i].Content)
					assert.NotEqual(t, uuid.Nil, results[i].ID)
					assert.False(t, results[i].CreatedAt.IsZero())
					assert.False(t, results[i].UpdatedAt.IsZero())
				}

				for i := 0; i < len(results)-1; i++ {
					assert.True(t, results[i].CreatedAt.After(results[i+1].CreatedAt) || 
						results[i].CreatedAt.Equal(results[i+1].CreatedAt))
				}
			}
		})
	}
}

func TestUpdateArtifact(t *testing.T) {
	InitTestDB()
	currentTime := time.Now()

	tests := []struct {
		name        string
		setup       func() *Artifact
		input       func(*Artifact) *Artifact
		expected    *Artifact
		expectError bool
	}{
		{
			name: "Successfully update text artifact content",
			setup: func() *Artifact {
				TestDB.db.Create(&ChatMessage{
					ID:      "msg123",
					ChatID:  "chat123",
					Message: "Test Message",
				})

				artifact := &Artifact{
					ID:        uuid.New(),
					MessageID: "msg123",
					Type:      TextArtifact,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "original content",
					},
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}
				TestDB.db.Create(artifact)
				return artifact
			},
			input: func(original *Artifact) *Artifact {
				return &Artifact{
					ID:        original.ID,
					MessageID: original.MessageID,
					Type:      original.Type,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "updated content",
					},
					CreatedAt: original.CreatedAt,
				}
			},
			expected: &Artifact{
				Type:      TextArtifact,
				MessageID: "msg123",
				Content: PropertyMap{
					"text_type": "code",
					"content":   "updated content",
				},
			},
			expectError: false,
		},
		{
			name: "Successfully update visual artifact content",
			setup: func() *Artifact {
				TestDB.db.Create(&ChatMessage{
					ID:      "msg456",
					ChatID:  "chat456",
					Message: "Test Message",
				})

				artifact := &Artifact{
					ID:        uuid.New(),
					MessageID: "msg456",
					Type:      VisualArtifact,
					Content: PropertyMap{
						"text_type": "img",
						"url":       "https://old.example.com/image.png",
					},
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}
				TestDB.db.Create(artifact)
				return artifact
			},
			input: func(original *Artifact) *Artifact {
				return &Artifact{
					ID:        original.ID,
					MessageID: original.MessageID,
					Type:      original.Type,
					Content: PropertyMap{
						"text_type": "img",
						"url":       "https://new.example.com/image.png",
					},
					CreatedAt: original.CreatedAt,
				}
			},
			expected: &Artifact{
				Type:      VisualArtifact,
				MessageID: "msg456",
				Content: PropertyMap{
					"text_type": "img",
					"url":       "https://new.example.com/image.png",
				},
			},
			expectError: false,
		},
		{
			name: "Successfully update action artifact content",
			setup: func() *Artifact {
				TestDB.db.Create(&ChatMessage{
					ID:      "msg789",
					ChatID:  "chat789",
					Message: "Test Message",
				})

				artifact := &Artifact{
					ID:        uuid.New(),
					MessageID: "msg789",
					Type:      ActionArtifact,
					Content: PropertyMap{
						"action_text": "Old option",
						"options": []interface{}{
							map[string]interface{}{
								"action_type":     "button",
								"option_label":    "Old button",
								"option_response": "old_click",
							},
						},
					},
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}
				TestDB.db.Create(artifact)
				return artifact
			},
			input: func(original *Artifact) *Artifact {
				return &Artifact{
					ID:        original.ID,
					MessageID: original.MessageID,
					Type:      original.Type,
					Content: PropertyMap{
						"action_text": "New option",
						"options": []interface{}{
							map[string]interface{}{
								"action_type":     "button",
								"option_label":    "New button",
								"option_response": "new_click",
							},
						},
					},
					CreatedAt: original.CreatedAt,
				}
			},
			expected: &Artifact{
				Type:      ActionArtifact,
				MessageID: "msg789",
				Content: PropertyMap{
					"action_text": "New option",
					"options": []interface{}{
						map[string]interface{}{
							"action_type":     "button",
							"option_label":    "New button",
							"option_response": "new_click",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Update with nil UUID",
			setup: func() *Artifact {
				return nil
			},
			input: func(original *Artifact) *Artifact {
				return &Artifact{
					ID: uuid.Nil,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "test content",
					},
				}
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Update non-existent artifact",
			setup: func() *Artifact {
				return nil
			},
			input: func(original *Artifact) *Artifact {
				return &Artifact{
					ID: uuid.New(),
					Content: PropertyMap{
						"text_type": "code",
						"content":   "test content",
					},
				}
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Update deleted artifact",
			setup: func() *Artifact {
				TestDB.db.Create(&ChatMessage{
					ID:      "msg_deleted",
					ChatID:  "chat_deleted",
					Message: "Test Message",
				})

				artifact := &Artifact{
					ID:        uuid.New(),
					MessageID: "msg_deleted",
					Type:      TextArtifact,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "original content",
					},
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}
				TestDB.db.Create(artifact)
				TestDB.db.Delete(artifact)
				return artifact
			},
			input: func(original *Artifact) *Artifact {
				return &Artifact{
					ID:        original.ID,
					MessageID: original.MessageID,
					Type:      original.Type,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "updated content",
					},
					CreatedAt: original.CreatedAt,
				}
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			TestDB.db.Exec("DELETE FROM artifacts")
			TestDB.db.Exec("DELETE FROM chat_messages")

			original := tt.setup()
			input := tt.input(original)

			result, err := TestDB.UpdateArtifact(input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			if original != nil {
				assert.Equal(t, original.ID, result.ID)
				assert.Equal(t, original.CreatedAt, result.CreatedAt)
				assert.True(t, result.UpdatedAt.After(original.UpdatedAt))
			}
			assert.Equal(t, tt.expected.MessageID, result.MessageID)
			assert.Equal(t, tt.expected.Type, result.Type)
			assert.Equal(t, tt.expected.Content, result.Content)

			var savedArtifact Artifact
			err = TestDB.db.First(&savedArtifact, "id = ?", result.ID).Error
			assert.NoError(t, err)
			assert.Equal(t, result.ID, savedArtifact.ID)
			assert.Equal(t, tt.expected.MessageID, savedArtifact.MessageID)
			assert.Equal(t, tt.expected.Type, savedArtifact.Type)
			assert.Equal(t, tt.expected.Content, savedArtifact.Content)
			if original != nil {
				assert.Equal(t, original.CreatedAt, savedArtifact.CreatedAt)
				assert.True(t, savedArtifact.UpdatedAt.After(original.UpdatedAt))
			}
		})
	}
}

func TestDeleteArtifactByID(t *testing.T) {
	InitTestDB()
	currentTime := time.Now()

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		input       uuid.UUID
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successfully delete text artifact",
			setup: func() uuid.UUID {

				TestDB.db.Create(&ChatMessage{
					ID:      "msg123",
					ChatID:  "chat123",
					Message: "Test Message",
				})

				artifact := &Artifact{
					ID:        uuid.New(),
					MessageID: "msg123",
					Type:      TextArtifact,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "test content",
					},
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}
				TestDB.db.Create(artifact)
				return artifact.ID
			},
			expectError: false,
		},
		{
			name: "Delete with nil UUID",
			setup: func() uuid.UUID {
				return uuid.Nil
			},
			expectError: true,
			errorMsg:    "artifact not found",
		},
		{
			name: "Delete non-existent artifact",
			setup: func() uuid.UUID {
				return uuid.New() 
			},
			expectError: true,
			errorMsg:    "artifact not found",
		},
		{
			name: "Delete already deleted artifact",
			setup: func() uuid.UUID {

				TestDB.db.Create(&ChatMessage{
					ID:      "msg_deleted",
					ChatID:  "chat_deleted",
					Message: "Test Message",
				})

				artifact := &Artifact{
					ID:        uuid.New(),
					MessageID: "msg_deleted",
					Type:      TextArtifact,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "deleted content",
					},
					CreatedAt: currentTime,
					UpdatedAt: currentTime,
				}
				TestDB.db.Create(artifact)
				TestDB.db.Delete(&Artifact{}, "id = ?", artifact.ID)
				return artifact.ID
			},
			expectError: true,
			errorMsg:    "artifact not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			TestDB.db.Exec("DELETE FROM artifacts")
			TestDB.db.Exec("DELETE FROM chat_messages")

			id := tt.setup()

			err := TestDB.DeleteArtifactByID(id)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)

				if tt.errorMsg == "artifact not found" {
					var count int64
					TestDB.db.Model(&Artifact{}).Where("id = ?", id).Count(&count)
					assert.Equal(t, int64(0), count)
				}
				return
			}

			assert.NoError(t, err)

			var count int64

			TestDB.db.Model(&Artifact{}).Where("id = ?", id).Count(&count)
			assert.Equal(t, int64(0), count)

			TestDB.db.Unscoped().Model(&Artifact{}).Where("id = ?", id).Count(&count)
			assert.Equal(t, int64(0), count)
		})
	}
}

func TestDeleteAllArtifactsByChatID(t *testing.T) {
	InitTestDB()
	currentTime := time.Now()

	tests := []struct {
		name        string
		setup       func() string
		input       string
		expectError bool
		verifyFunc  func(*testing.T, string)
	}{
		{
			name: "Successfully delete multiple artifacts from multiple messages",
			setup: func() string {
				chatID := "chat123"

				messages := []ChatMessage{
					{
						ID:      "msg1",
						ChatID:  chatID,
						Message: "First Message",
					},
					{
						ID:      "msg2",
						ChatID:  chatID,
						Message: "Second Message",
					},
				}
				for _, msg := range messages {
					TestDB.db.Create(&msg)
				}

				artifacts1 := []Artifact{
					{
						ID:        uuid.New(),
						MessageID: "msg1",
						Type:      TextArtifact,
						Content: PropertyMap{
							"text_type": "code",
							"content":   "first message code",
						},
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
					{
						ID:        uuid.New(),
						MessageID: "msg1",
						Type:      VisualArtifact,
						Content: PropertyMap{
							"text_type": "img",
							"url":       "https://example.com/first.png",
						},
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
				}

				artifacts2 := []Artifact{
					{
						ID:        uuid.New(),
						MessageID: "msg2",
						Type:      ActionArtifact,
						Content: PropertyMap{
							"action_text": "Choose option",
							"options": []interface{}{
								map[string]interface{}{
									"action_type":     "button",
									"option_label":    "Click",
									"option_response": "clicked",
								},
							},
						},
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
					},
				}

				for _, a := range artifacts1 {
					TestDB.db.Create(&a)
				}
				for _, a := range artifacts2 {
					TestDB.db.Create(&a)
				}

				return chatID
			},
			verifyFunc: func(t *testing.T, chatID string) {

				var count int64
				TestDB.db.Model(&Artifact{}).
					Joins("JOIN chat_messages ON artifacts.message_id = chat_messages.id").
					Where("chat_messages.chat_id = ?", chatID).
					Count(&count)
				assert.Equal(t, int64(0), count)

				TestDB.db.Model(&ChatMessage{}).Where("chat_id = ?", chatID).Count(&count)
				assert.Equal(t, int64(2), count)
			},
		},
		{
			name: "Delete from chat with no messages",
			setup: func() string {
				return "empty_chat"
			},
			verifyFunc: func(t *testing.T, chatID string) {
				var count int64
				TestDB.db.Model(&Artifact{}).
					Joins("JOIN chat_messages ON artifacts.message_id = chat_messages.id").
					Where("chat_messages.chat_id = ?", chatID).
					Count(&count)
				assert.Equal(t, int64(0), count)
			},
		},
		{
			name: "Delete from chat with messages but no artifacts",
			setup: func() string {
				chatID := "chat_no_artifacts"
				messages := []ChatMessage{
					{
						ID:      "msg_empty1",
						ChatID:  chatID,
						Message: "Message without artifacts 1",
					},
					{
						ID:      "msg_empty2",
						ChatID:  chatID,
						Message: "Message without artifacts 2",
					},
				}
				for _, msg := range messages {
					TestDB.db.Create(&msg)
				}
				return chatID
			},
			verifyFunc: func(t *testing.T, chatID string) {
				var count int64
				TestDB.db.Model(&Artifact{}).
					Joins("JOIN chat_messages ON artifacts.message_id = chat_messages.id").
					Where("chat_messages.chat_id = ?", chatID).
					Count(&count)
				assert.Equal(t, int64(0), count)

				TestDB.db.Model(&ChatMessage{}).Where("chat_id = ?", chatID).Count(&count)
				assert.Equal(t, int64(2), count)
			},
		},
		{
			name: "Delete with empty chat ID",
			setup: func() string {
				return ""
			},
			verifyFunc: func(t *testing.T, chatID string) {
				var count int64
				TestDB.db.Model(&Artifact{}).Count(&count)
				assert.True(t, count >= 0)
			},
		},
		{
			name: "Verify other chats' artifacts remain untouched",
			setup: func() string {
				targetChatID := "chat_to_delete"
				otherChatID := "chat_to_keep"

				TestDB.db.Create(&ChatMessage{
					ID:      "msg_delete",
					ChatID:  targetChatID,
					Message: "Message to delete",
				})
				TestDB.db.Create(&Artifact{
					ID:        uuid.New(),
					MessageID: "msg_delete",
					Type:      TextArtifact,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "delete me",
					},
				})

				TestDB.db.Create(&ChatMessage{
					ID:      "msg_keep",
					ChatID:  otherChatID,
					Message: "Message to keep",
				})
				TestDB.db.Create(&Artifact{
					ID:        uuid.New(),
					MessageID: "msg_keep",
					Type:      TextArtifact,
					Content: PropertyMap{
						"text_type": "code",
						"content":   "keep me",
					},
				})

				return targetChatID
			},
			verifyFunc: func(t *testing.T, chatID string) {
				var count int64
				TestDB.db.Model(&Artifact{}).
					Joins("JOIN chat_messages ON artifacts.message_id = chat_messages.id").
					Where("chat_messages.chat_id = ?", chatID).
					Count(&count)
				assert.Equal(t, int64(0), count)

				TestDB.db.Model(&Artifact{}).
					Joins("JOIN chat_messages ON artifacts.message_id = chat_messages.id").
					Where("chat_messages.chat_id = ?", "chat_to_keep").
					Count(&count)
				assert.Equal(t, int64(1), count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestDB.db.Exec("DELETE FROM artifacts")
			TestDB.db.Exec("DELETE FROM chat_messages")

			chatID := tt.setup()

			err := TestDB.DeleteAllArtifactsByChatID(chatID)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			tt.verifyFunc(t, chatID)
		})
	}
}