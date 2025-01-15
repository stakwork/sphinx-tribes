package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateNotification(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	CleanTestData()

	currentTime := time.Now()

	tests := []struct {
		name         string
		notification *Notification
		expectError  bool
		errorMsg     string
	}{
		{
			name: "Successfully create notification with all fields",
			notification: &Notification{
				UUID:      uuid.New().String(),
				Event:     "test_event",
				PubKey:    "test_pub_key",
				Content:   "test_content",
				Status:    NotificationStatusPending,
				CreatedAt: &currentTime,
				UpdatedAt: &currentTime,
			},
			expectError: false,
		},
		{
			name: "Missing event",
			notification: &Notification{
				UUID:    uuid.New().String(),
				PubKey:  "test_pub_key",
				Content: "test_content",
			},
			expectError: true,
			errorMsg:    "event is required",
		},
		{
			name: "Missing public key",
			notification: &Notification{
				UUID:    uuid.New().String(),
				Event:   "test_event",
				Content: "test_content",
			},
			expectError: true,
			errorMsg:    "public key is required",
		},
		{
			name: "Missing content",
			notification: &Notification{
				UUID:   uuid.New().String(),
				Event:  "test_event",
				PubKey: "test_pub_key",
			},
			expectError: true,
			errorMsg:    "content is required",
		},
		{
			name: "Auto-generate UUID when empty",
			notification: &Notification{
				Event:   "test_event",
				PubKey:  "test_pub_key",
				Content: "test_content",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TestDB.CreateNotification(tt.notification)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tt.notification.UUID)

				var saved Notification
				result := TestDB.db.Where("uuid = ?", tt.notification.UUID).First(&saved)
				assert.NoError(t, result.Error)
				assert.Equal(t, tt.notification.Event, saved.Event)
				assert.Equal(t, tt.notification.PubKey, saved.PubKey)
				assert.Equal(t, tt.notification.Content, saved.Content)
				assert.Equal(t, NotificationStatusPending, saved.Status)
				assert.NotNil(t, saved.CreatedAt)
				assert.NotNil(t, saved.UpdatedAt)
			}
		})
	}
}

func TestGetNotification(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	testNotif := Notification{
		UUID:    uuid.New().String(),
		Event:   "test_event",
		PubKey:  "test_pub_key",
		Content: "test_content",
		Status:  NotificationStatusPending,
	}
	err := TestDB.db.Create(&testNotif).Error
	assert.NoError(t, err, "Failed to create test notification")

	tests := []struct {
		name        string
		uuid        string
		expectError bool
		expectNil   bool
	}{
		{
			name:        "Successfully get existing notification",
			uuid:        testNotif.UUID,
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "Notification not found",
			uuid:        uuid.New().String(),
			expectError: false,
			expectNil:   true,
		},
		{
			name:        "Empty UUID",
			uuid:        "",
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notification, err := TestDB.GetNotification(tt.uuid)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, notification)
			} else {
				assert.NoError(t, err)
				if tt.expectNil {
					assert.Nil(t, notification)
				} else {
					assert.NotNil(t, notification)
					assert.Equal(t, testNotif.UUID, notification.UUID)
					assert.Equal(t, testNotif.Event, notification.Event)
					assert.Equal(t, testNotif.PubKey, notification.PubKey)
					assert.Equal(t, testNotif.Content, notification.Content)
				}
			}
		})
	}
}

func TestUpdateNotification(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name        string
		setup       func() string
		updates     map[string]interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successfully update notification",
			setup: func() string {
				n := &Notification{
					Event:   "test_event",
					PubKey:  "test_pub_key",
					Content: "test_content",
				}
				err := TestDB.CreateNotification(n)
				assert.NoError(t, err)
				return n.UUID
			},
			updates: map[string]interface{}{
				"content": "updated_content",
				"status":  NotificationStatusComplete,
			},
			expectError: false,
		},
		{
			name: "Update non-existent notification",
			setup: func() string {
				return uuid.New().String()
			},
			updates: map[string]interface{}{
				"content": "updated_content",
			},
			expectError: true,
			errorMsg:    "notification not found",
		},
		{
			name: "Empty UUID",
			setup: func() string {
				return ""
			},
			updates: map[string]interface{}{
				"content": "updated_content",
			},
			expectError: true,
			errorMsg:    "notification UUID is required",
		},
		{
			name: "No updates provided",
			setup: func() string {
				n := &Notification{
					Event:   "test_event",
					PubKey:  "test_pub_key",
					Content: "test_content",
				}
				err := TestDB.CreateNotification(n)
				assert.NoError(t, err)
				return n.UUID
			},
			updates:     map[string]interface{}{},
			expectError: true,
			errorMsg:    "no updates provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanTestData()
			setup := tt.setup()

			err := TestDB.UpdateNotification(setup, tt.updates)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)

				updated, err := TestDB.GetNotification(setup)
				assert.NoError(t, err)
				assert.NotNil(t, updated)

				for key, value := range tt.updates {
					switch key {
					case "content":
						assert.Equal(t, value, updated.Content)
					case "status":
						assert.Equal(t, value, updated.Status)
					}
				}
			}
		})
	}
}

func TestDeleteNotification(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successfully delete notification",
			setup: func() string {
				n := &Notification{
					Event:   "test_event",
					PubKey:  "test_pub_key",
					Content: "test_content",
				}
				_ = TestDB.CreateNotification(n)
				return n.UUID
			},
			expectError: false,
		},
		{
			name: "Empty UUID",
			setup: func() string {
				return ""
			},
			expectError: true,
			errorMsg:    "uuid is required",
		},
		{
			name: "Non-existent notification",
			setup: func() string {
				return uuid.New().String()
			},
			expectError: true,
			errorMsg:    "notification not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanTestData()
			setup := tt.setup()

			err := TestDB.DeleteNotification(setup)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				var count int64
				TestDB.db.Model(&Notification{}).Where("setup = ?", setup).Count(&count)
				assert.Equal(t, int64(0), count)
			}
		})
	}
}

func TestGetPendingNotifications(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	currentTime := time.Now()

	tests := []struct {
		name          string
		setup         func()
		expectedCount int
		expectError   bool
		expectedLen   int
		maxRetries    int
	}{
		{
			name:       "Get failed notifications within retry limit",
			maxRetries: 3,
			setup: func() {
				notifications := []Notification{
					{
						UUID:      uuid.New().String(),
						Event:     "event1",
						PubKey:    "pub_key1",
						Content:   "content1",
						Status:    NotificationStatusFailed,
						Retries:   1,
						CreatedAt: &currentTime,
					},
					{
						UUID:      uuid.New().String(),
						Event:     "event2",
						PubKey:    "pub_key2",
						Content:   "content2",
						Status:    NotificationStatusFailed,
						Retries:   2,
						CreatedAt: &currentTime,
					},
				}
				for _, n := range notifications {
					err := TestDB.db.Create(&n).Error
					assert.NoError(t, err)
				}
			},
			expectedLen: 2,
			expectError: false,
		},
		{
			name: "No pending notifications",
			setup: func() {
				n := &Notification{
					UUID:      uuid.New().String(),
					Event:     "event1",
					PubKey:    "pub_key1",
					Content:   "content1",
					Status:    NotificationStatusComplete,
					CreatedAt: &currentTime,
				}
				err := TestDB.db.Create(n).Error
				assert.NoError(t, err)
			},
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			CleanTestData()

			notifications, err := TestDB.GetPendingNotifications()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(notifications))
				for _, n := range notifications {
					assert.Equal(t, NotificationStatusPending, n.Status)
				}
			}
		})
	}
}

func TestGetFailedNotifications(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	CleanTestData()

	notifications := []Notification{
		{
			UUID:    uuid.New().String(),
			Event:   "event1",
			PubKey:  "pub_key1",
			Content: "content1",
			Status:  NotificationStatusFailed,
			Retries: 1,
		},
		{
			UUID:    uuid.New().String(),
			Event:   "event2",
			PubKey:  "pub_key2",
			Content: "content2",
			Status:  NotificationStatusFailed,
			Retries: 2,
		},
	}

	for _, n := range notifications {
		err := TestDB.db.Create(&n).Error
		assert.NoError(t, err)
	}

	var count int64
	err := TestDB.db.Model(&Notification{}).Where("status = ?", NotificationStatusFailed).Count(&count).Error
	assert.NoError(t, err)

	maxRetries := 3
	results, err := TestDB.GetFailedNotifications(maxRetries)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(results), "Expected 2 failed notifications, got %d", len(results))

	for _, n := range results {
		assert.Equal(t, NotificationStatusFailed, n.Status)
		assert.Less(t, n.Retries, maxRetries)
	}
}

func TestGetNotificationsByPubKey(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	CleanTestData()

	testPubKey := "test_pub_key"

	notifications := []Notification{
		{
			UUID:    uuid.New().String(),
			Event:   "event1",
			PubKey:  testPubKey,
			Content: "content1",
			Status:  NotificationStatusPending,
		},
		{
			UUID:    uuid.New().String(),
			Event:   "event2",
			PubKey:  testPubKey,
			Content: "content2",
			Status:  NotificationStatusPending,
		},
		{
			UUID:    uuid.New().String(),
			Event:   "event3",
			PubKey:  testPubKey,
			Content: "content3",
			Status:  NotificationStatusPending,
		},
	}

	for _, n := range notifications {
		err := TestDB.db.Create(&n).Error
		assert.NoError(t, err)
	}

	var count int64
	err := TestDB.db.Model(&Notification{}).Where("pub_key = ?", testPubKey).Count(&count).Error
	assert.NoError(t, err)

	results, err := TestDB.GetNotificationsByPubKey(testPubKey, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(results), "Expected 3 notifications, got %d", len(results))

	for _, n := range results {
		assert.Equal(t, testPubKey, n.PubKey)
	}
}

func TestIncrementRetryCount(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successfully increment retry count",
			setup: func() string {
				n := &Notification{
					UUID:    uuid.New().String(),
					Event:   "test_event",
					PubKey:  "test_pub_key",
					Content: "test_content",
					Retries: 0,
				}
				err := TestDB.db.Create(n).Error
				assert.NoError(t, err)
				return n.UUID
			},
			expectError: false,
		},
		{
			name: "Empty UUID",
			setup: func() string {
				return ""
			},
			expectError: true,
			errorMsg:    "uuid is required",
		},
		{
			name: "Non-existent notification",
			setup: func() string {
				return uuid.New().String()
			},
			expectError: true,
			errorMsg:    "notification not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CleanTestData()
			setup := tt.setup()

			err := TestDB.IncrementRetryCount(setup)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)

				var notification Notification
				err = TestDB.db.Where("uuid = ?", setup).First(&notification).Error
				assert.NoError(t, err)
				assert.Equal(t, 1, notification.Retries)
			}
		})
	}
}

func TestGetNotificationCount(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name        string
		pubKey      string
		setup       func(t *testing.T)
		expected    int64
		expectError bool
		errorMsg    string
	}{
		{
			name:   "Get count for existing notifications",
			pubKey: "test_pub_key",
			setup: func(t *testing.T) {
				CleanTestData()

				for i := 0; i < 3; i++ {
					notification := &Notification{
						UUID:    uuid.New().String(),
						Event:   fmt.Sprintf("event%d", i),
						PubKey:  "test_pub_key",
						Content: fmt.Sprintf("content%d", i),
						Status:  NotificationStatusPending,
					}

					err := TestDB.db.Create(notification).Error
					assert.NoError(t, err)
				}

				var count int64
				err := TestDB.db.Model(&Notification{}).
					Where("pub_key = ?", "test_pub_key").
					Count(&count).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(3), count, "Failed to create test notifications")
			},
			expected:    3,
			expectError: false,
		},
		{
			name:        "Empty public key",
			pubKey:      "",
			setup:       func(t *testing.T) {},
			expected:    0,
			expectError: true,
			errorMsg:    "public key is required",
		},
		{
			name:        "No notifications for public key",
			pubKey:      "non_existent_key",
			setup:       func(t *testing.T) {},
			expected:    0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)

			count, err := TestDB.GetNotificationCount(tt.pubKey)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, count)
			}
		})
	}
}
