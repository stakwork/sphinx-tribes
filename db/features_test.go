package db

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteFeatureByUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	person := Person{
		Uuid:        uuid.New().String(),
		OwnerPubKey: "test-pubkey",
		OwnerAlias:  "test-alias",
	}
	TestDB.CreateOrEditPerson(person)

	workspace := Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace",
		OwnerPubKey: person.OwnerPubKey,
	}
	TestDB.CreateOrEditWorkspace(workspace)

	currentTime := time.Now()

	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successfully delete existing feature",
			setup: func() string {
				feature := WorkspaceFeatures{
					Uuid:          uuid.New().String(),
					WorkspaceUuid: workspace.Uuid,
					Name:          "Test Feature",
					Brief:         "Test Brief",
					Priority:      1,
					Created:       &currentTime,
					Updated:       &currentTime,
					CreatedBy:     person.OwnerPubKey,
					UpdatedBy:     person.OwnerPubKey,
					FeatStatus:    ActiveFeature,
				}
				result := TestDB.db.Create(&feature)
				assert.NoError(t, result.Error)
				return feature.Uuid
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "Try to delete non-existent feature",
			setup: func() string {
				return uuid.New().String()
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "Try to delete with empty UUID",
			setup: func() string {
				return ""
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "Delete archived feature",
			setup: func() string {
				feature := WorkspaceFeatures{
					Uuid:          uuid.New().String(),
					WorkspaceUuid: workspace.Uuid,
					Name:          "Archived Feature",
					FeatStatus:    ArchivedFeature,
					Created:       &currentTime,
					CreatedBy:     person.OwnerPubKey,
				}
				result := TestDB.db.Create(&feature)
				assert.NoError(t, result.Error)
				return feature.Uuid
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "Try to delete already deleted feature",
			setup: func() string {
				feature := WorkspaceFeatures{
					Uuid:          uuid.New().String(),
					WorkspaceUuid: workspace.Uuid,
					Name:          "To Be Deleted",
					Created:       &currentTime,
					CreatedBy:     person.OwnerPubKey,
				}
				TestDB.db.Create(&feature)
				TestDB.DeleteFeatureByUuid(feature.Uuid)
				return feature.Uuid
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "Case sensitivity test",
			setup: func() string {
				featureUuid := uuid.New().String()
				feature := WorkspaceFeatures{
					Uuid:          featureUuid,
					WorkspaceUuid: workspace.Uuid,
					Name:          "Case Sensitivity Test",
					Created:       &currentTime,
					CreatedBy:     person.OwnerPubKey,
				}
				TestDB.db.Create(&feature)
				return strings.ToUpper(featureUuid)
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "SQL injection attempt",
			setup: func() string {
				return "' OR '1'='1"
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "Invalid UUID format",
			setup: func() string {
				return "invalid-uuid-format"
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "Feature with special characters in UUID",
			setup: func() string {
				return "!@#$%^&*()"
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "Feature with very long UUID",
			setup: func() string {
				return strings.Repeat("a", 1000)
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "Feature with Unicode characters in UUID",
			setup: func() string {
				return "测试UUID"
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "Concurrent deletion attempt",
			setup: func() string {
				featureUuid := uuid.New().String()
				feature := WorkspaceFeatures{
					Uuid:          featureUuid,
					WorkspaceUuid: workspace.Uuid,
					Name:          "Concurrent Test",
					Created:       &currentTime,
					CreatedBy:     person.OwnerPubKey,
				}
				TestDB.db.Create(&feature)

				go func() {
					TestDB.DeleteFeatureByUuid(featureUuid)
				}()
				time.Sleep(10 * time.Millisecond)

				return featureUuid
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "Feature with maximum field values",
			setup: func() string {
				feature := WorkspaceFeatures{
					Uuid:          uuid.New().String(),
					WorkspaceUuid: workspace.Uuid,
					Name:          strings.Repeat("a", 255),
					Brief:         strings.Repeat("b", 1000),
					Requirements:  strings.Repeat("c", 1000),
					Architecture:  strings.Repeat("d", 1000),
					Priority:      999999,
					Created:       &currentTime,
					CreatedBy:     person.OwnerPubKey,
				}
				TestDB.db.Create(&feature)
				return feature.Uuid
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "Large Number of Features",
			setup: func() string {
				targetUuid := uuid.New().String()

				for i := 0; i < 100; i++ {
					feature := WorkspaceFeatures{
						Uuid:          uuid.New().String(),
						WorkspaceUuid: workspace.Uuid,
						Name:          fmt.Sprintf("Bulk Feature %d", i),
						Created:       &currentTime,
						CreatedBy:     person.OwnerPubKey,
						FeatStatus:    ActiveFeature,
					}
					if i == 50 {
						feature.Uuid = targetUuid
					}
					TestDB.db.Create(&feature)
				}
				return targetUuid
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "UUID with Maximum Length",
			setup: func() string {
				maxUuid := strings.Repeat("a", 36)
				feature := WorkspaceFeatures{
					Uuid:          maxUuid,
					WorkspaceUuid: workspace.Uuid,
					Name:          "Max Length UUID Feature",
					Created:       &currentTime,
					CreatedBy:     person.OwnerPubKey,
				}
				TestDB.db.Create(&feature)
				return maxUuid
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "Null UUID",
			setup: func() string {
				return "00000000-0000-0000-0000-000000000000"
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "Feature with Related Data",
			setup: func() string {
				featureUuid := uuid.New().String()
				feature := WorkspaceFeatures{
					Uuid:          featureUuid,
					WorkspaceUuid: workspace.Uuid,
					Name:          "Feature with Relations",
					Created:       &currentTime,
					CreatedBy:     person.OwnerPubKey,
				}
				TestDB.db.Create(&feature)

				story := FeatureStory{
					Uuid:        uuid.New().String(),
					FeatureUuid: featureUuid,
					Description: "Related Story",
					Created:     &currentTime,
					CreatedBy:   person.OwnerPubKey,
				}
				TestDB.db.Create(&story)

				return featureUuid
			},
			expectError: false,
			errorMsg:    "",
		},
		{
			name: "Multiple Concurrent Deletions",
			setup: func() string {
				featureUuid := uuid.New().String()
				feature := WorkspaceFeatures{
					Uuid:          featureUuid,
					WorkspaceUuid: workspace.Uuid,
					Name:          "Concurrent Deletion Test",
					Created:       &currentTime,
					CreatedBy:     person.OwnerPubKey,
				}
				TestDB.db.Create(&feature)

				var wg sync.WaitGroup
				for i := 0; i < 5; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						TestDB.DeleteFeatureByUuid(featureUuid)
					}()
				}
				wg.Wait()

				return featureUuid
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "UUID with Mixed Case",
			setup: func() string {
				originalUuid := uuid.New().String()
				feature := WorkspaceFeatures{
					Uuid:          originalUuid,
					WorkspaceUuid: workspace.Uuid,
					Name:          "Mixed Case UUID Test",
					Created:       &currentTime,
					CreatedBy:     person.OwnerPubKey,
				}
				TestDB.db.Create(&feature)
				return strings.ToUpper(originalUuid[0:18]) + strings.ToLower(originalUuid[18:])
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "UUID with Whitespace",
			setup: func() string {
				return "  " + uuid.New().String() + "  "
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
		{
			name: "UUID with Line Breaks",
			setup: func() string {
				return "test\nuuid\r\nvalue"
			},
			expectError: true,
			errorMsg:    "no feature found to delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuid := tt.setup()

			err := TestDB.DeleteFeatureByUuid(uuid)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)

				var feature WorkspaceFeatures
				result := TestDB.db.Where("uuid = ?", uuid).First(&feature)
				assert.Error(t, result.Error)
				assert.True(t, errors.Is(result.Error, gorm.ErrRecordNotFound))
			}
		})
	}
}
