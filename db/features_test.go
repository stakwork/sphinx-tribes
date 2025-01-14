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

func TestGetProductBrief(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	CleanTestData()

	tests := []struct {
		name          string
		workspaceUuid string
		setup         func() error
		expected      string
		expectError   bool
	}{
		{
			name:          "Valid UUID with Complete Data",
			workspaceUuid: "valid-uuid-1",
			setup: func() error {
				workspace := Workspace{
					Uuid:    "valid-uuid-1",
					Name:    "Product1",
					Mission: "Mission1",
					Tactics: "Tactics1",
				}
				return TestDB.db.Create(&workspace).Error
			},
			expected:    "Product: Product1. Product Brief:\n Mission: Mission1.\n\n Objectives: Tactics1",
			expectError: false,
		},
		{
			name:          "Valid UUID with Partial Data",
			workspaceUuid: "valid-uuid-2",
			setup: func() error {
				workspace := Workspace{
					Uuid:    "valid-uuid-2",
					Name:    "Product2",
					Mission: "",
					Tactics: "Tactics2",
				}
				return TestDB.db.Create(&workspace).Error
			},
			expected:    "Product: Product2. Product Brief:\n Mission: .\n\n Objectives: Tactics2",
			expectError: false,
		},
		{
			name:          "Empty UUID",
			workspaceUuid: "",
			setup:         func() error { return nil },
			expected:      "",
			expectError:   true,
		},
		{
			name:          "Non-Existent UUID",
			workspaceUuid: "non-existent-uuid",
			setup:         func() error { return nil },
			expected:      "",
			expectError:   true,
		},
		{
			name:          "Case Sensitivity",
			workspaceUuid: "VALID-UUID-5",
			setup: func() error {
				workspace := Workspace{
					Uuid:    "valid-uuid-5",
					Name:    "Product5",
					Mission: "Mission5",
					Tactics: "Tactics5",
				}
				return TestDB.db.Create(&workspace).Error
			},
			expected:    "",
			expectError: true,
		},
		{
			name:          "Whitespace in UUID",
			workspaceUuid: " valid-uuid-6 ",
			setup: func() error {
				workspace := Workspace{
					Uuid:    "valid-uuid-6",
					Name:    "Product6",
					Mission: "Mission6",
					Tactics: "Tactics6",
				}
				return TestDB.db.Create(&workspace).Error
			},
			expected:    "",
			expectError: true,
		},
		{
			name:          "UUID with Special Characters",
			workspaceUuid: "valid-uuid-special-123",
			setup: func() error {
				workspace := Workspace{
					Uuid:    "valid-uuid-special-123",
					Name:    "ProductSpecial",
					Mission: "MissionSpecial",
					Tactics: "TacticsSpecial",
				}
				return TestDB.db.Create(&workspace).Error
			},
			expected:    "Product: ProductSpecial. Product Brief:\n Mission: MissionSpecial.\n\n Objectives: TacticsSpecial",
			expectError: false,
		},
		{
			name:          "Null UUID",
			workspaceUuid: "",
			setup:         func() error { return nil },
			expected:      "",
			expectError:   true,
		},
		{
			name:          "Invalid UUID Format",
			workspaceUuid: "invalid-uuid-!@#",
			setup:         func() error { return nil },
			expected:      "",
			expectError:   true,
		},
		{
			name:          "Concurrent Access",
			workspaceUuid: "valid-uuid-7",
			setup: func() error {
				workspace := Workspace{
					Uuid:    "valid-uuid-7",
					Name:    "Product7",
					Mission: "Mission7",
					Tactics: "Tactics7",
				}
				return TestDB.db.Create(&workspace).Error
			},
			expected:    "Product: Product7. Product Brief:\n Mission: Mission7.\n\n Objectives: Tactics7",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.setup(); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			result, err := TestDB.GetProductBrief(tt.workspaceUuid)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
		CleanTestData()
	}
}

func TestGetFeatureBrief(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	CleanTestData()

	tests := []struct {
		name        string
		featureUuid string
		setup       func() error
		expected    string
		expectError bool
	}{
		{
			name:        "Valid UUID with Complete Data",
			featureUuid: "valid-uuid-1",
			setup: func() error {
				feature := WorkspaceFeatures{
					Uuid:  "valid-uuid-1",
					Name:  "Feature1",
					Brief: "This is a test feature",
				}
				return TestDB.db.Create(&feature).Error
			},
			expected:    "Feature: Feature1. Brief: This is a test feature",
			expectError: false,
		},
		{
			name:        "Empty UUID",
			featureUuid: "",
			setup:       func() error { return nil },
			expected:    "",
			expectError: true,
		},
		{
			name:        "Non-Existent UUID",
			featureUuid: "non-existent-uuid",
			setup:       func() error { return nil },
			expected:    "",
			expectError: true,
		},
		{
			name:        "Case Sensitivity",
			featureUuid: "VALID-UUID-2",
			setup: func() error {
				feature := WorkspaceFeatures{
					Uuid:  "valid-uuid-2",
					Name:  "Feature2",
					Brief: "This is another test feature",
				}
				return TestDB.db.Create(&feature).Error
			},
			expected:    "",
			expectError: true,
		},
		{
			name:        "Whitespace in UUID",
			featureUuid: " valid-uuid-3 ",
			setup: func() error {
				feature := WorkspaceFeatures{
					Uuid:  "valid-uuid-3",
					Name:  "Feature3",
					Brief: "Feature brief with spaces",
				}
				return TestDB.db.Create(&feature).Error
			},
			expected:    "",
			expectError: true,
		},
		{
			name:        "UUID with Special Characters",
			featureUuid: "valid-uuid-special-123",
			setup: func() error {
				feature := WorkspaceFeatures{
					Uuid:  "valid-uuid-special-123",
					Name:  "SpecialFeature",
					Brief: "Feature with special characters",
				}
				return TestDB.db.Create(&feature).Error
			},
			expected:    "Feature: SpecialFeature. Brief: Feature with special characters",
			expectError: false,
		},
		{
			name:        "Feature with Empty Fields",
			featureUuid: "uuid-with-empty-fields",
			setup: func() error {
				feature := WorkspaceFeatures{
					Uuid:  "uuid-with-empty-fields",
					Name:  "",
					Brief: "",
				}
				return TestDB.db.Create(&feature).Error
			},
			expected:    "Feature: . Brief: ",
			expectError: false,
		},
		{
			name:        "Invalid UUID Format",
			featureUuid: "invalid-uuid-!@#",
			setup:       func() error { return nil },
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.setup(); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			result, err := TestDB.GetFeatureBrief(tt.featureUuid)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
		CleanTestData()
	}
}
