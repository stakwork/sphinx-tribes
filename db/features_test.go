package db

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteFeatureStoryByUuid(t *testing.T) {

	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	person := Person{
		Uuid:        uuid.New().String(),
		OwnerPubKey: "test-pubkey",
		OwnerAlias:  "test-alias",
	}
	TestDB.CreateOrEditPerson(person)

	currentTime := time.Now()

	tests := []struct {
		name        string
		setup       func() (string, string)
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successfully delete existing feature story",
			setup: func() (string, string) {
				featureUuid := uuid.New().String()
				storyUuid := uuid.New().String()

				story := FeatureStory{
					Uuid:        storyUuid,
					FeatureUuid: featureUuid,
					Description: "Test story",
					Priority:    1,
					Created:     &currentTime,
					Updated:     &currentTime,
					CreatedBy:   person.OwnerPubKey,
					UpdatedBy:   person.OwnerPubKey,
				}

				result := TestDB.db.Create(&story)
				assert.NoError(t, result.Error)

				return featureUuid, storyUuid
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "Try to delete non-existent feature story",
			setup: func() (string, string) {
				return uuid.New().String(), uuid.New().String()
			},
			expectError: true,
			errorMsg:    "no story found to delete",
		},
		{
			name: "Try to delete with empty UUIDs",
			setup: func() (string, string) {
				return "", ""
			},
			expectError: true,
			errorMsg:    "no story found to delete",
		},
		{
			name: "Try to delete with valid feature UUID but invalid story UUID",
			setup: func() (string, string) {
				featureUuid := uuid.New().String()
				story := FeatureStory{
					Uuid:        uuid.New().String(),
					FeatureUuid: featureUuid,
					Description: "Test story",
					Priority:    1,
					Created:     &currentTime,
					CreatedBy:   person.OwnerPubKey,
				}

				result := TestDB.db.Create(&story)
				assert.NoError(t, result.Error)

				return featureUuid, uuid.New().String()
			},
			expectError: true,
			errorMsg:    "no story found to delete",
		},
		{
			name: "Try to delete with invalid feature UUID but valid story UUID",
			setup: func() (string, string) {
				featureUuid := uuid.New().String()
				storyUuid := uuid.New().String()

				story := FeatureStory{
					Uuid:        storyUuid,
					FeatureUuid: featureUuid,
					Description: "Test story",
					Priority:    1,
					Created:     &currentTime,
					CreatedBy:   person.OwnerPubKey,
				}

				result := TestDB.db.Create(&story)
				assert.NoError(t, result.Error)

				return uuid.New().String(), storyUuid
			},
			expectError: true,
			errorMsg:    "no story found to delete",
		},
		{
			name: "Delete story and verify it's gone",
			setup: func() (string, string) {
				featureUuid := uuid.New().String()
				storyUuid := uuid.New().String()

				story := FeatureStory{
					Uuid:        storyUuid,
					FeatureUuid: featureUuid,
					Description: "Test story",
					Priority:    1,
					Created:     &currentTime,
					CreatedBy:   person.OwnerPubKey,
				}

				result := TestDB.db.Create(&story)
				assert.NoError(t, result.Error)

				return featureUuid, storyUuid
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "Try to delete already deleted story",
			setup: func() (string, string) {
				featureUuid := uuid.New().String()
				storyUuid := uuid.New().String()

				story := FeatureStory{
					Uuid:        storyUuid,
					FeatureUuid: featureUuid,
					Description: "Test story",
					Priority:    1,
					Created:     &currentTime,
					CreatedBy:   person.OwnerPubKey,
				}

				TestDB.db.Create(&story)
				TestDB.db.Delete(&story)

				return featureUuid, storyUuid
			},
			expectError: true,
			errorMsg:    "no story found to delete",
		},
		{
			name: "Case Sensitivity Test",
			setup: func() (string, string) {
				featureUuid := uuid.New().String()
				storyUuid := uuid.New().String()

				story := FeatureStory{
					Uuid:        storyUuid,
					FeatureUuid: featureUuid,
					Description: "Case sensitivity test story",
					Priority:    1,
					Created:     &currentTime,
					CreatedBy:   person.OwnerPubKey,
				}

				result := TestDB.db.Create(&story)
				assert.NoError(t, result.Error)

				return strings.ToUpper(featureUuid), strings.ToUpper(storyUuid)
			},
			expectError: true,
			errorMsg:    "no story found to delete",
		},
		{
			name: "SQL Injection Attempt",
			setup: func() (string, string) {
				return "' OR '1'='1", "' OR '1'='1"
			},
			expectError: true,
			errorMsg:    "no story found to delete",
		},
		{
			name: "Invalid UUID Format",
			setup: func() (string, string) {
				return "invalid-uuid-format", "another-invalid-uuid"
			},
			expectError: true,
			errorMsg:    "no story found to delete",
		},
		{
			name: "Multiple Stories with Same Feature UUID",
			setup: func() (string, string) {
				featureUuid := uuid.New().String()
				storyUuid1 := uuid.New().String()
				storyUuid2 := uuid.New().String()

				story1 := FeatureStory{
					Uuid:        storyUuid1,
					FeatureUuid: featureUuid,
					Description: "First story",
					Priority:    1,
					Created:     &currentTime,
					CreatedBy:   person.OwnerPubKey,
				}

				story2 := FeatureStory{
					Uuid:        storyUuid2,
					FeatureUuid: featureUuid,
					Description: "Second story",
					Priority:    2,
					Created:     &currentTime,
					CreatedBy:   person.OwnerPubKey,
				}

				TestDB.db.Create(&story1)
				TestDB.db.Create(&story2)

				return featureUuid, storyUuid1
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "Special Characters in UUID",
			setup: func() (string, string) {
				return "!@#$%^&*()", "!@#$%^&*()"
			},
			expectError: true,
			errorMsg:    "no story found to delete",
		},
		{
			name: "Very Long UUID Values",
			setup: func() (string, string) {
				return strings.Repeat("a", 1000), strings.Repeat("b", 1000)
			},
			expectError: true,
			errorMsg:    "no story found to delete",
		},
		{
			name: "Unicode Characters in UUID",
			setup: func() (string, string) {
				return "测试UUID", "テストUUID"
			},
			expectError: true,
			errorMsg:    "no story found to delete",
		},
		{
			name: "Concurrent Deletion Attempt",
			setup: func() (string, string) {
				featureUuid := uuid.New().String()
				storyUuid := uuid.New().String()

				story := FeatureStory{
					Uuid:        storyUuid,
					FeatureUuid: featureUuid,
					Description: "Concurrent deletion test",
					Priority:    1,
					Created:     &currentTime,
					CreatedBy:   person.OwnerPubKey,
				}

				TestDB.db.Create(&story)

				go func() {
					TestDB.DeleteFeatureStoryByUuid(featureUuid, storyUuid)
				}()
				time.Sleep(10 * time.Millisecond)

				return featureUuid, storyUuid
			},
			expectError: true,
			errorMsg:    "no story found to delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			featureUuid, storyUuid := tt.setup()

			err := TestDB.DeleteFeatureStoryByUuid(featureUuid, storyUuid)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)

				var story FeatureStory
				result := TestDB.db.Where("feature_uuid = ? AND uuid = ?", featureUuid, storyUuid).First(&story)
				assert.Error(t, result.Error)
				assert.True(t, errors.Is(result.Error, gorm.ErrRecordNotFound))
			}
		})
	}
}
