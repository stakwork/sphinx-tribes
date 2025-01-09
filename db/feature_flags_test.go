package db

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateFeatureFlag(t *testing.T) {

	InitTestDB()
	defer CloseTestDB()

	TestDB.db.Exec("DELETE FROM feature_flags")

	existingFlag := FeatureFlag{
		UUID:        uuid.New(),
		Name:        "Test Feature " + uuid.New().String(),
		Description: "A test feature flag",
		Enabled:     false,
		UpdatedAt:   time.Now(),
	}

	if err := TestDB.db.Create(&existingFlag).Error; err != nil {
		t.Fatalf("Failed to create test feature flag: %v", err)
	}

	tests := []struct {
		name          string
		inputFlag     FeatureFlag
		expectedError string
	}{
		{
			name: "Update Existing Feature Flag",
			inputFlag: FeatureFlag{
				UUID:        existingFlag.UUID,
				Name:        "Updated Feature " + uuid.New().String(),
				Description: "Updated description",
				Enabled:     true,
			},
			expectedError: "",
		},
		{
			name: "Missing UUID",
			inputFlag: FeatureFlag{
				UUID:        uuid.Nil,
				Name:        "No UUID Feature " + uuid.New().String(),
				Description: "Should fail",
				Enabled:     false,
			},
			expectedError: "feature flag UUID is required",
		},
		{
			name: "Non-existent Feature Flag",
			inputFlag: FeatureFlag{
				UUID:        uuid.New(),
				Name:        "Non-existent Feature " + uuid.New().String(),
				Description: "Should fail",
				Enabled:     false,
			},
			expectedError: "feature flag not found",
		},
		{
			name: "Update with Empty Name",
			inputFlag: FeatureFlag{
				UUID:        existingFlag.UUID,
				Name:        "",
				Description: "Empty name update",
				Enabled:     true,
			},
			expectedError: "feature flag name cannot be empty",
		},
		{
			name: "Update with Long Description",
			inputFlag: FeatureFlag{
				UUID:        existingFlag.UUID,
				Name:        "Long Description Feature",
				Description: string(make([]byte, 1001)),
				Enabled:     true,
			},
			expectedError: "description too long",
		},
		{
			name: "Toggle Enabled State",
			inputFlag: FeatureFlag{
				UUID:        existingFlag.UUID,
				Name:        "Toggle Enabled Feature",
				Description: "Toggle enabled state",
				Enabled:     !existingFlag.Enabled,
			},
			expectedError: "",
		},
		{
			name: "Update with No Changes",
			inputFlag: FeatureFlag{
				UUID:        existingFlag.UUID,
				Name:        existingFlag.Name,
				Description: existingFlag.Description,
				Enabled:     existingFlag.Enabled,
			},
			expectedError: "",
		},
		{
			name: "Update with Special Characters in Name",
			inputFlag: FeatureFlag{
				UUID:        existingFlag.UUID,
				Name:        "Special!@#$%^&*()_+",
				Description: "Special characters in name",
				Enabled:     true,
			},
			expectedError: "",
		},
		{
			name: "Update with Null Description",
			inputFlag: FeatureFlag{
				UUID:        existingFlag.UUID,
				Name:        "Null Description Feature",
				Description: "",
				Enabled:     true,
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := TestDB.UpdateFeatureFlag(&tt.inputFlag)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
