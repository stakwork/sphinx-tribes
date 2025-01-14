package db

import (
	"errors"
	"fmt"
	"strconv"
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

func TestGetPhaseByUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	currentTime := time.Now()
	tests := []struct {
		name          string
		setup         func() string
		expectedPhase *FeaturePhase
		expectError   bool
		errorMessage  string
	}{
		{
			name: "Successfully get phase with all fields populated",
			setup: func() string {
				phase := FeaturePhase{
					Uuid:         uuid.New().String(),
					FeatureUuid:  uuid.New().String(),
					Name:         "Test Phase",
					Priority:     1,
					PhasePurpose: "Test Purpose",
					PhaseOutcome: "Test Outcome",
					PhaseScope:   "Test Scope",
					Created:      &currentTime,
					Updated:      &currentTime,
					CreatedBy:    "test-user",
					UpdatedBy:    "test-user",
				}
				TestDB.db.Create(&phase)
				return phase.Uuid
			},
			expectedPhase: &FeaturePhase{
				Name:         "Test Phase",
				Priority:     1,
				PhasePurpose: "Test Purpose",
				PhaseOutcome: "Test Outcome",
				PhaseScope:   "Test Scope",
				CreatedBy:    "test-user",
				UpdatedBy:    "test-user",
			},
			expectError: false,
		},
		{
			name: "Phase not found",
			setup: func() string {
				return uuid.New().String()
			},
			expectedPhase: nil,
			expectError:   true,
			errorMessage:  "no phase found",
		},
		{
			name: "Empty UUID",
			setup: func() string {
				return ""
			},
			expectedPhase: nil,
			expectError:   true,
			errorMessage:  "no phase found",
		},
		{
			name: "Phase with special characters in fields",
			setup: func() string {
				phase := FeaturePhase{
					Uuid:         uuid.New().String(),
					FeatureUuid:  uuid.New().String(),
					Name:         "Special !@#$%^&*()",
					Priority:     1,
					PhasePurpose: "Purpose !@#$%^&*()",
					PhaseOutcome: "Outcome !@#$%^&*()",
					PhaseScope:   "Scope !@#$%^&*()",
					Created:      &currentTime,
					Updated:      &currentTime,
					CreatedBy:    "test-user!@#",
					UpdatedBy:    "test-user!@#",
				}
				TestDB.db.Create(&phase)
				return phase.Uuid
			},
			expectedPhase: &FeaturePhase{
				Name:         "Special !@#$%^&*()",
				Priority:     1,
				PhasePurpose: "Purpose !@#$%^&*()",
				PhaseOutcome: "Outcome !@#$%^&*()",
				PhaseScope:   "Scope !@#$%^&*()",
				CreatedBy:    "test-user!@#",
				UpdatedBy:    "test-user!@#",
			},
			expectError: false,
		},
		{
			name: "Phase with Unicode characters",
			setup: func() string {
				phase := FeaturePhase{
					Uuid:         uuid.New().String(),
					FeatureUuid:  uuid.New().String(),
					Name:         "测试阶段",
					Priority:     1,
					PhasePurpose: "目的 🎯",
					PhaseOutcome: "结果 ✨",
					PhaseScope:   "范围 🌟",
					Created:      &currentTime,
					Updated:      &currentTime,
					CreatedBy:    "用户",
					UpdatedBy:    "用户",
				}
				TestDB.db.Create(&phase)
				return phase.Uuid
			},
			expectedPhase: &FeaturePhase{
				Name:         "测试阶段",
				Priority:     1,
				PhasePurpose: "目的 🎯",
				PhaseOutcome: "结果 ✨",
				PhaseScope:   "范围 🌟",
				CreatedBy:    "用户",
				UpdatedBy:    "用户",
			},
			expectError: false,
		},
		{
			name: "Phase with maximum length strings",
			setup: func() string {
				longString := strings.Repeat("a", 255)
				phase := FeaturePhase{
					Uuid:         uuid.New().String(),
					FeatureUuid:  uuid.New().String(),
					Name:         longString,
					Priority:     1,
					PhasePurpose: longString,
					PhaseOutcome: longString,
					PhaseScope:   longString,
					Created:      &currentTime,
					Updated:      &currentTime,
					CreatedBy:    longString,
					UpdatedBy:    longString,
				}
				TestDB.db.Create(&phase)
				return phase.Uuid
			},
			expectedPhase: &FeaturePhase{
				Name:         strings.Repeat("a", 255),
				Priority:     1,
				PhasePurpose: strings.Repeat("a", 255),
				PhaseOutcome: strings.Repeat("a", 255),
				PhaseScope:   strings.Repeat("a", 255),
				CreatedBy:    strings.Repeat("a", 255),
				UpdatedBy:    strings.Repeat("a", 255),
			},
			expectError: false,
		},
		{
			name: "UUID with Leading and Trailing Spaces",
			setup: func() string {
				phase := FeaturePhase{
					Uuid:        uuid.New().String(),
					FeatureUuid: uuid.New().String(),
					Name:        "Space Test Phase",
					Priority:    1,
					Created:     &currentTime,
					Updated:     &currentTime,
					CreatedBy:   "test-user",
				}
				TestDB.db.Create(&phase)
				return "  " + phase.Uuid + "  "
			},
			expectedPhase: &FeaturePhase{
				Name:      "Space Test Phase",
				Priority:  1,
				CreatedBy: "test-user",
			},
			expectError:  true,
			errorMessage: "no phase found",
		},
		{
			name: "Case Sensitivity in UUID",
			setup: func() string {
				originalUuid := "12345678-ABCD-EFGH-IJKL-MNOPQRSTUVWX"
				phase := FeaturePhase{
					Uuid:        originalUuid,
					FeatureUuid: uuid.New().String(),
					Name:        "Case Sensitivity Test",
					Priority:    1,
					Created:     &currentTime,
					Updated:     &currentTime,
					CreatedBy:   "test-user",
				}
				TestDB.db.Create(&phase)
				return strings.ToLower(originalUuid)
			},
			expectedPhase: &FeaturePhase{
				Name:      "Case Sensitivity Test",
				Priority:  1,
				CreatedBy: "test-user",
			},
			expectError:  true,
			errorMessage: "no phase found",
		},
		{
			name: "With Large UUID String",
			setup: func() string {
				largeUuid := strings.Repeat("a", 1000)
				return largeUuid
			},
			expectedPhase: nil,
			expectError:   true,
			errorMessage:  "no phase found",
		},
		{
			name: "UUID with Invalid Format",
			setup: func() string {
				return "not-a-valid-uuid-format"
			},
			expectedPhase: nil,
			expectError:   true,
			errorMessage:  "no phase found",
		},
		{
			name: "UUID with SQL Injection Attempt",
			setup: func() string {
				return "' OR '1'='1"
			},
			expectedPhase: nil,
			expectError:   true,
			errorMessage:  "no phase found",
		},
		{
			name: "UUID with Null Characters",
			setup: func() string {
				return "12345678-abcd-efgh-ijkl-mnopqrstuvwx\x00"
			},
			expectedPhase: nil,
			expectError:   true,
			errorMessage:  "no phase found",
		},
		{
			name: "UUID with HTML Characters",
			setup: func() string {
				return "<script>alert('test')</script>"
			},
			expectedPhase: nil,
			expectError:   true,
			errorMessage:  "no phase found",
		},
		{
			name: "Multiple Identical UUIDs Attempt",
			setup: func() string {
				sharedUuid := uuid.New().String()

				phase1 := FeaturePhase{
					Uuid:        sharedUuid,
					FeatureUuid: uuid.New().String(),
					Name:        "First Phase",
					Priority:    1,
					Created:     &currentTime,
					Updated:     &currentTime,
					CreatedBy:   "test-user",
					UpdatedBy:   "test-user",
				}
				err := TestDB.db.Create(&phase1).Error
				if err != nil {
					return sharedUuid
				}

				phase2 := FeaturePhase{
					Uuid:        sharedUuid,
					FeatureUuid: uuid.New().String(),
					Name:        "Second Phase",
					Priority:    2,
					Created:     &currentTime,
					Updated:     &currentTime,
					CreatedBy:   "test-user-2",
					UpdatedBy:   "test-user-2",
				}
				_ = TestDB.db.Create(&phase2)

				return sharedUuid
			},
			expectedPhase: &FeaturePhase{
				Name:      "First Phase",
				Priority:  1,
				CreatedBy: "test-user",
				UpdatedBy: "test-user",
			},
			expectError: false,
		},
		{
			name: "UUID with Mixed Case and Special Characters",
			setup: func() string {
				complexUuid := "Ab#12$Cd-EfGh-IjKl-MnOp-QrStUvWxYz"
				phase := FeaturePhase{
					Uuid:        complexUuid,
					FeatureUuid: uuid.New().String(),
					Name:        "Complex UUID Test",
					Priority:    1,
					Created:     &currentTime,
					CreatedBy:   "test-user",
				}
				TestDB.db.Create(&phase)
				return "ab#12$cd-efgh-ijkl-mnop-qrstuvwxyz"
			},
			expectedPhase: nil,
			expectError:   true,
			errorMessage:  "no phase found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			TestDB.db.Exec("DELETE FROM feature_phases")

			phaseUuid := tt.setup()
			phase, err := TestDB.GetPhaseByUuid(phaseUuid)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, phaseUuid, phase.Uuid)
				assert.Equal(t, tt.expectedPhase.Name, phase.Name)
				assert.Equal(t, tt.expectedPhase.Priority, phase.Priority)
				assert.Equal(t, tt.expectedPhase.PhasePurpose, phase.PhasePurpose)
				assert.Equal(t, tt.expectedPhase.PhaseOutcome, phase.PhaseOutcome)
				assert.Equal(t, tt.expectedPhase.PhaseScope, phase.PhaseScope)
				assert.Equal(t, tt.expectedPhase.CreatedBy, phase.CreatedBy)
				assert.Equal(t, tt.expectedPhase.UpdatedBy, phase.UpdatedBy)
				assert.NotNil(t, phase.Created)
				assert.NotNil(t, phase.Updated)
			}
		})
	}
}

func TestDeleteFeaturePhase(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()
	CleanTestData()

	tests := []struct {
		name                 string
		featureUuid          string
		phaseUuid            string
		setup                func(featureUuid, phaseUuid string)
		expectError          bool
		expectedRowsAffected int64
	}{
		{
			name:        "Valid feature and phase UUIDs",
			featureUuid: uuid.New().String(),
			phaseUuid:   uuid.New().String(),
			setup: func(featureUuid, phaseUuid string) {

				feature := WorkspaceFeatures{
					Uuid:          featureUuid,
					Name:          "Test Feature",
					WorkspaceUuid: "workspace1",
				}
				TestDB.CreateOrEditFeature(feature)

				phase := FeaturePhase{
					Uuid:        phaseUuid,
					FeatureUuid: featureUuid,
					Name:        "Phase 1",
				}
				TestDB.CreateOrEditFeaturePhase(phase)
			},
			expectError:          false,
			expectedRowsAffected: 0,
		},
		{
			name:        "Non-Existent Phase",
			featureUuid: uuid.New().String(),
			phaseUuid:   uuid.New().String(),
			setup: func(featureUuid, _ string) {

				feature := WorkspaceFeatures{
					Uuid: featureUuid,
					Name: "Test Feature",
				}
				TestDB.CreateOrEditFeature(feature)
			},
			expectError:          true,
			expectedRowsAffected: 0,
		},
		{
			name:                 "Non-Existent Feature",
			featureUuid:          uuid.New().String(),
			phaseUuid:            uuid.New().String(),
			setup:                func(_, _ string) {},
			expectError:          true,
			expectedRowsAffected: 0,
		},
		{
			name:                 "Both UUIDs Non-Existent",
			featureUuid:          uuid.New().String(),
			phaseUuid:            uuid.New().String(),
			setup:                func(_, _ string) {},
			expectError:          true,
			expectedRowsAffected: 0,
		},
		{
			name:                 "Empty UUIDs",
			featureUuid:          "",
			phaseUuid:            "",
			setup:                func(featureUuid, phaseUuid string) {},
			expectError:          true,
			expectedRowsAffected: 0,
		},
		{
			name:                 "Invalid UUID Format",
			featureUuid:          "invalid-uuid-format",
			phaseUuid:            "invalid-uuid-format",
			setup:                func(featureUuid, phaseUuid string) {},
			expectError:          true,
			expectedRowsAffected: 0,
		},
		{
			name:        "Case Sensitivity",
			featureUuid: "VALID-FEATURE-UUID",
			phaseUuid:   "VALID-PHASE-UUID",
			setup: func(featureUuid, phaseUuid string) {

				feature := WorkspaceFeatures{
					Uuid: "valid-feature-uuid",
					Name: "Test Feature",
				}
				TestDB.CreateOrEditFeature(feature)
				phase := FeaturePhase{
					Uuid:        "valid-phase-uuid",
					FeatureUuid: "valid-feature-uuid",
					Name:        "Phase 1",
				}
				TestDB.CreateOrEditFeaturePhase(phase)
			},
			expectError:          true,
			expectedRowsAffected: 0,
		},
		{
			name:        "SQL Injection Attempt",
			featureUuid: "valid-feature-uuid'; DROP TABLE FeaturePhase; --",
			phaseUuid:   "valid-phase-uuid",
			setup: func(featureUuid, phaseUuid string) {
				feature := WorkspaceFeatures{
					Uuid: "valid-feature-uuid4",
					Name: "Test Feature",
				}
				TestDB.CreateOrEditFeature(feature)

				phase := FeaturePhase{
					Uuid:        "valid-phase-uuid4",
					FeatureUuid: "valid-feature-uuid",
					Name:        "Phase 1",
				}
				TestDB.CreateOrEditFeaturePhase(phase)
			},
			expectError:          true,
			expectedRowsAffected: 0,
		},
		{
			name:        "Large Number of Phases",
			featureUuid: "valid-feature-uuid",
			phaseUuid:   "valid-phase-uuid",
			setup: func(featureUuid, phaseUuid string) {
				phases := make(map[string]bool)
				for i := 0; i < 1000000; i++ {
					phases[strconv.Itoa(i)] = true
				}
				phases["valid-phase-uuid"] = true
				TestDB.CreateOrEditFeaturePhase(FeaturePhase{Uuid: "valid-phase-uuid", FeatureUuid: featureUuid, Name: "Phase 1"})
			},
			expectError:          false,
			expectedRowsAffected: 0,
		},
		{
			name:        "Concurrent Deletion",
			featureUuid: "valid-feature-uuid",
			phaseUuid:   "valid-phase-uuid",
			setup: func(featureUuid, phaseUuid string) {

				TestDB.CreateOrEditFeaturePhase(FeaturePhase{Uuid: phaseUuid, FeatureUuid: featureUuid, Name: "Phase 1"})
			},
			expectError:          false,
			expectedRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.featureUuid, tt.phaseUuid)

			err := TestDB.DeleteFeaturePhase(tt.featureUuid, tt.phaseUuid)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				var count int64
				TestDB.db.Model(&FeaturePhase{}).Where("uuid = ? AND feature_uuid = ?", tt.phaseUuid, tt.featureUuid).Count(&count)
				assert.Equal(t, tt.expectedRowsAffected, count)
			}
		})
		CleanTestData()
	}
}
