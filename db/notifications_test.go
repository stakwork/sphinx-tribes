package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNotificationCRUD(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	cleanup := func() {
		TestDB.db.Exec("DELETE FROM notifications")
	}

	t.Run("CreateNotification", func(t *testing.T) {
		defer cleanup()

		tests := []struct {
			name         string
			notification *Notification
			expectError  error
		}{
			{
				name: "Valid Notification",
				notification: &Notification{
					Event:   "test_event",
					PubKey:  "test_pubkey",
					Content: "test_content",
				},
				expectError: nil,
			},
			{
				name: "Missing Event",
				notification: &Notification{
					PubKey:  "test_pubkey",
					Content: "test_content",
				},
				expectError: ErrMissingEvent,
			},
			{
				name: "Missing PubKey",
				notification: &Notification{
					Event:   "test_event",
					Content: "test_content",
				},
				expectError: ErrMissingPubKey,
			},
			{
				name: "Missing Content",
				notification: &Notification{
					Event:  "test_event",
					PubKey: "test_pubkey",
				},
				expectError: ErrMissingContent,
			},
			{
				name: "With Custom UUID",
				notification: &Notification{
					UUID:    uuid.New().String(),
					Event:   "test_event",
					PubKey:  "test_pubkey",
					Content: "test_content",
				},
				expectError: nil,
			},
			{
				name: "With Custom Status",
				notification: &Notification{
					Event:   "test_event",
					PubKey:  "test_pubkey",
					Content: "test_content",
					Status:  NotificationStatusComplete,
				},
				expectError: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := TestDB.CreateNotification(tt.notification)
				assert.Equal(t, tt.expectError, err)

				if err == nil {
					assert.NotEmpty(t, tt.notification.UUID)
					assert.NotNil(t, tt.notification.CreatedAt)
					assert.NotNil(t, tt.notification.UpdatedAt)
					if tt.notification.Status == "" {
						assert.Equal(t, NotificationStatusPending, tt.notification.Status)
					}
				}
			})
		}
	})

	t.Run("GetNotification", func(t *testing.T) {
		defer cleanup()

		testNotif := &Notification{
			Event:   "test_event",
			PubKey:  "test_pubkey",
			Content: "test_content",
		}
		err := TestDB.CreateNotification(testNotif)
		assert.NoError(t, err)

		tests := []struct {
			name        string
			uuid        string
			expectError error
			expectNil   bool
		}{
			{
				name:        "Valid UUID",
				uuid:        testNotif.UUID,
				expectError: nil,
				expectNil:   false,
			},
			{
				name:        "Empty UUID",
				uuid:        "",
				expectError: ErrMissingUUID,
				expectNil:   true,
			},
			{
				name:        "Non-existent UUID",
				uuid:        uuid.New().String(),
				expectError: nil,
				expectNil:   true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				notification, err := TestDB.GetNotification(tt.uuid)
				assert.Equal(t, tt.expectError, err)
				if tt.expectNil {
					assert.Nil(t, notification)
				} else {
					assert.NotNil(t, notification)
					assert.Equal(t, testNotif.UUID, notification.UUID)
				}
			})
		}
	})

	t.Run("UpdateNotification", func(t *testing.T) {
		defer cleanup()

		testNotif := &Notification{
			Event:   "test_event",
			PubKey:  "test_pubkey",
			Content: "test_content",
		}
		err := TestDB.CreateNotification(testNotif)
		assert.NoError(t, err)

		tests := []struct {
			name        string
			uuid        string
			updates     map[string]interface{}
			expectError error
		}{
			{
				name: "Valid Update",
				uuid: testNotif.UUID,
				updates: map[string]interface{}{
					"content": "updated_content",
				},
				expectError: nil,
			},
			{
				name:        "Empty UUID",
				uuid:        "",
				updates:     map[string]interface{}{},
				expectError: ErrMissingUUID,
			},
			{
				name: "Empty Event",
				uuid: testNotif.UUID,
				updates: map[string]interface{}{
					"event": "",
				},
				expectError: ErrMissingEvent,
			},
			{
				name: "Empty PubKey",
				uuid: testNotif.UUID,
				updates: map[string]interface{}{
					"pub_key": "",
				},
				expectError: ErrMissingPubKey,
			},
			{
				name: "Empty Content",
				uuid: testNotif.UUID,
				updates: map[string]interface{}{
					"content": "",
				},
				expectError: ErrMissingContent,
			},
			{
				name:        "Non-existent UUID",
				uuid:        uuid.New().String(),
				updates:     map[string]interface{}{},
				expectError: gorm.ErrRecordNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := TestDB.UpdateNotification(tt.uuid, tt.updates)
				assert.Equal(t, tt.expectError, err)
			})
		}
	})

	t.Run("DeleteNotification", func(t *testing.T) {
		defer cleanup()

		testNotif := &Notification{
			Event:   "test_event",
			PubKey:  "test_pubkey",
			Content: "test_content",
		}
		err := TestDB.CreateNotification(testNotif)
		assert.NoError(t, err)

		tests := []struct {
			name        string
			uuid        string
			expectError error
		}{
			{
				name:        "Valid UUID",
				uuid:        testNotif.UUID,
				expectError: nil,
			},
			{
				name:        "Empty UUID",
				uuid:        "",
				expectError: ErrMissingUUID,
			},
			{
				name:        "Non-existent UUID",
				uuid:        uuid.New().String(),
				expectError: gorm.ErrRecordNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := TestDB.DeleteNotification(tt.uuid)
				assert.Equal(t, tt.expectError, err)
			})
		}
	})

	t.Run("GetPendingNotifications", func(t *testing.T) {
		defer cleanup()

		notifications := []*Notification{
			{Event: "event1", PubKey: "key1", Content: "content1", Status: NotificationStatusPending},
			{Event: "event2", PubKey: "key2", Content: "content2", Status: NotificationStatusComplete},
			{Event: "event3", PubKey: "key3", Content: "content3", Status: NotificationStatusPending},
		}

		for _, n := range notifications {
			err := TestDB.CreateNotification(n)
			assert.NoError(t, err)
		}

		result, err := TestDB.GetPendingNotifications()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		for _, n := range result {
			assert.Equal(t, NotificationStatusPending, n.Status)
		}
	})

	t.Run("GetFailedNotifications", func(t *testing.T) {
		defer cleanup()

		notifications := []*Notification{
			{Event: "event1", PubKey: "key1", Content: "content1", Status: NotificationStatusFailed, Retries: 1},
			{Event: "event2", PubKey: "key2", Content: "content2", Status: NotificationStatusFailed, Retries: 2},
			{Event: "event3", PubKey: "key3", Content: "content3", Status: NotificationStatusFailed, Retries: 3},
		}

		for _, n := range notifications {
			err := TestDB.CreateNotification(n)
			assert.NoError(t, err)
		}

		tests := []struct {
			name          string
			maxRetries    int
			expectedCount int
		}{
			{
				name:          "Max Retries 0",
				maxRetries:    0,
				expectedCount: 0,
			},
			{
				name:          "Negative Max Retries",
				maxRetries:    -1,
				expectedCount: 0,
			},
			{
				name:          "Max Retries 5",
				maxRetries:    5,
				expectedCount: 3,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := TestDB.GetFailedNotifications(tt.maxRetries)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(result))
			})
		}
	})

	t.Run("GetNotificationsByPubKey", func(t *testing.T) {
		defer cleanup()

		pubKey := "test_key"
		notifications := []*Notification{
			{Event: "event1", PubKey: pubKey, Content: "content1"},
			{Event: "event2", PubKey: pubKey, Content: "content2"},
			{Event: "event3", PubKey: "other_key", Content: "content3"},
		}

		for _, n := range notifications {
			err := TestDB.CreateNotification(n)
			assert.NoError(t, err)
		}

		tests := []struct {
			name          string
			pubKey        string
			limit         int
			offset        int
			expectedCount int
			expectError   error
		}{
			{
				name:          "Valid PubKey",
				pubKey:        pubKey,
				limit:         10,
				offset:        0,
				expectedCount: 2,
				expectError:   nil,
			},
			{
				name:          "Empty PubKey",
				pubKey:        "",
				limit:         10,
				offset:        0,
				expectedCount: 0,
				expectError:   ErrMissingPubKey,
			},
			{
				name:          "With Limit",
				pubKey:        pubKey,
				limit:         1,
				offset:        0,
				expectedCount: 1,
				expectError:   nil,
			},
			{
				name:          "With Offset",
				pubKey:        pubKey,
				limit:         10,
				offset:        1,
				expectedCount: 1,
				expectError:   nil,
			},
			{
				name:          "Negative Limit",
				pubKey:        pubKey,
				limit:         -1,
				offset:        0,
				expectedCount: 2,
				expectError:   nil,
			},
			{
				name:          "Negative Offset",
				pubKey:        pubKey,
				limit:         10,
				offset:        -1,
				expectedCount: 2,
				expectError:   nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := TestDB.GetNotificationsByPubKey(tt.pubKey, tt.limit, tt.offset)
				assert.Equal(t, tt.expectError, err)
				if err == nil {
					assert.Equal(t, tt.expectedCount, len(result))
				}
			})
		}
	})

	t.Run("IncrementRetryCount", func(t *testing.T) {
		defer cleanup()

		testNotif := &Notification{
			Event:   "test_event",
			PubKey:  "test_pubkey",
			Content: "test_content",
			Retries: 0,
		}
		err := TestDB.CreateNotification(testNotif)
		assert.NoError(t, err)

		tests := []struct {
			name        string
			uuid        string
			expectError error
		}{
			{
				name:        "Valid UUID",
				uuid:        testNotif.UUID,
				expectError: nil,
			},
			{
				name:        "Empty UUID",
				uuid:        "",
				expectError: ErrMissingUUID,
			},
			{
				name:        "Non-existent UUID",
				uuid:        uuid.New().String(),
				expectError: gorm.ErrRecordNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := TestDB.IncrementRetryCount(tt.uuid)
				assert.Equal(t, tt.expectError, err)

				if err == nil {
					notification, _ := TestDB.GetNotification(tt.uuid)
					assert.Equal(t, testNotif.Retries+1, notification.Retries)
				}
			})
		}
	})

	t.Run("GetNotificationCount", func(t *testing.T) {
		defer cleanup()

		pubKey := "test_key"
		notifications := []*Notification{
			{Event: "event1", PubKey: pubKey, Content: "content1"},
			{Event: "event2", PubKey: pubKey, Content: "content2"},
			{Event: "event3", PubKey: "other_key", Content: "content3"},
		}

		for _, n := range notifications {
			err := TestDB.CreateNotification(n)
			assert.NoError(t, err)
		}

		tests := []struct {
			name          string
			pubKey        string
			expectedCount int64
			expectError   error
		}{
			{
				name:          "Valid PubKey",
				pubKey:        pubKey,
				expectedCount: 2,
				expectError:   nil,
			},
			{
				name:          "Empty PubKey",
				pubKey:        "",
				expectedCount: 0,
				expectError:   ErrMissingPubKey,
			},
			{
				name:          "Non-existent PubKey",
				pubKey:        "non_existent",
				expectedCount: 0,
				expectError:   nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				count, err := TestDB.GetNotificationCount(tt.pubKey)
				assert.Equal(t, tt.expectError, err)
				if err == nil {
					assert.Equal(t, tt.expectedCount, count)
				}
			})
		}
	})
}
