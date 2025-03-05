package db

import (
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateFeatureCall(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	workspace := Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace",
		OwnerPubKey: "test-pubkey",
	}
	TestDB.CreateOrEditWorkspace(workspace)

	tests := []struct {
		name        string
		setup       func() (string, string)
		expectError bool
		errorMsg    string
		validate    func(t *testing.T, result *FeatureCall)
	}{
		{
			name: "Successfully create new feature call",
			setup: func() (string, string) {
				return workspace.Uuid, "https://example.com/api"
			},
			expectError: false,
			validate: func(t *testing.T, result *FeatureCall) {
				assert.NotNil(t, result)
				assert.Equal(t, workspace.Uuid, result.WorkspaceID)
				assert.Equal(t, "https://example.com/api", result.URL)
				assert.NotNil(t, result.CreatedAt)
				assert.NotNil(t, result.UpdatedAt)
			},
		},
		{
			name: "Successfully update existing feature call",
			setup: func() (string, string) {
				featureCall := &FeatureCall{
					ID:          uuid.New(),
					WorkspaceID: workspace.Uuid,
					URL:         "https://old.example.com/api",
				}
				TestDB.db.Create(featureCall)
				return workspace.Uuid, "https://new.example.com/api"
			},
			expectError: false,
			validate: func(t *testing.T, result *FeatureCall) {
				assert.NotNil(t, result)
				assert.Equal(t, workspace.Uuid, result.WorkspaceID)
				assert.Equal(t, "https://new.example.com/api", result.URL)
			},
		},
		{
			name: "Empty workspace ID",
			setup: func() (string, string) {
				return "", "https://example.com/api"
			},
			expectError: true,
			errorMsg:    "workspace_id is required",
		},
		{
			name: "Empty URL",
			setup: func() (string, string) {
				return workspace.Uuid, ""
			},
			expectError: true,
			errorMsg:    "url is required",
		},
		{
			name: "Non-existent workspace",
			setup: func() (string, string) {
				return uuid.New().String(), "https://example.com/api"
			},
			expectError: true,
			errorMsg:    "workspace not found",
		},
		{
			name: "URL with special characters",
			setup: func() (string, string) {
				return workspace.Uuid, "https://example.com/api?key=value&special=!@#$%^&*()"
			},
			expectError: false,
			validate: func(t *testing.T, result *FeatureCall) {
				assert.NotNil(t, result)
				assert.Equal(t, "https://example.com/api?key=value&special=!@#$%^&*()", result.URL)
			},
		},
		{
			name: "URL with Unicode characters",
			setup: func() (string, string) {
				return workspace.Uuid, "https://例子.com/测试"
			},
			expectError: false,
			validate: func(t *testing.T, result *FeatureCall) {
				assert.NotNil(t, result)
				assert.Equal(t, "https://例子.com/测试", result.URL)
			},
		},
		{
			name: "Very long URL",
			setup: func() (string, string) {
				return workspace.Uuid, "https://example.com/" + strings.Repeat("a", 1000)
			},
			expectError: false,
			validate: func(t *testing.T, result *FeatureCall) {
				assert.NotNil(t, result)
				assert.Equal(t, "https://example.com/"+strings.Repeat("a", 1000), result.URL)
			},
		},
		{
			name: "SQL injection attempt in URL",
			setup: func() (string, string) {
				return workspace.Uuid, "'; DROP TABLE feature_calls; --"
			},
			expectError: false,
			validate: func(t *testing.T, result *FeatureCall) {
				assert.NotNil(t, result)
				assert.Equal(t, "'; DROP TABLE feature_calls; --", result.URL)
			},
		},
		{
			name: "Concurrent create/update",
			setup: func() (string, string) {
				var wg sync.WaitGroup
				for i := 0; i < 5; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						TestDB.CreateOrUpdateFeatureCall(workspace.Uuid, "https://concurrent.example.com/api")
					}()
				}
				wg.Wait()
				return workspace.Uuid, "https://final.example.com/api"
			},
			expectError: false,
			validate: func(t *testing.T, result *FeatureCall) {
				assert.NotNil(t, result)
				assert.Equal(t, "https://final.example.com/api", result.URL)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspaceID, url := tt.setup()

			result, err := TestDB.CreateOrUpdateFeatureCall(workspaceID, url)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestGetFeatureCallByWorkspaceID(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	DeleteAllFeatureCalls()

	workspace := Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace",
		OwnerPubKey: "test-pubkey",
	}
	TestDB.CreateOrEditWorkspace(workspace)

	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		errorMsg    string
		validate    func(t *testing.T, result *FeatureCall)
	}{
		{
			name: "Successfully retrieve existing feature call",
			setup: func() string {
				featureCall := &FeatureCall{
					ID:          uuid.New(),
					WorkspaceID: workspace.Uuid,
					URL:         "https://example.com/api",
				}
				TestDB.db.Create(featureCall)
				return workspace.Uuid
			},
			expectError: false,
			validate: func(t *testing.T, result *FeatureCall) {
				assert.NotNil(t, result)
				assert.Equal(t, workspace.Uuid, result.WorkspaceID)
				assert.Equal(t, "https://example.com/api", result.URL)
			},
		},
		{
			name: "Empty workspace ID",
			setup: func() string {
				return ""
			},
			expectError: true,
			errorMsg:    "workspace_id is required",
		},
		{
			name: "Non-existent workspace ID",
			setup: func() string {
				return uuid.New().String()
			},
			expectError: true,
			errorMsg:    "record not found",
		},
		{
			name: "Retrieve after soft delete",
			setup: func() string {
				featureCall := &FeatureCall{
					ID:          uuid.New(),
					WorkspaceID: workspace.Uuid,
					URL:         "https://example.com/api",
				}
				TestDB.db.Create(featureCall)
				TestDB.db.Delete(featureCall)
				return workspace.Uuid
			},
			expectError: true,
			errorMsg:    "record not found",
		},
		{
			name: "Case sensitivity test",
			setup: func() string {
				featureCall := &FeatureCall{
					ID:          uuid.New(),
					WorkspaceID: workspace.Uuid,
					URL:         "https://example.com/api",
				}
				TestDB.db.Create(featureCall)
				return strings.ToUpper(workspace.Uuid)
			},
			expectError: true,
			errorMsg:    "record not found",
		},
		{
			name: "SQL injection attempt",
			setup: func() string {
				return "' OR '1'='1"
			},
			expectError: true,
			errorMsg:    "record not found",
		},
		{
			name: "Multiple feature calls for same workspace",
			setup: func() string {
				DeleteAllFeatureCalls()

				featureCall := &FeatureCall{
					ID:          uuid.New(),
					WorkspaceID: workspace.Uuid,
					URL:         "https://first.example.com/api",
				}
				err := TestDB.db.Create(featureCall).Error
				assert.NoError(t, err)

				featureCall.URL = "https://second.example.com/api"
				err = TestDB.db.Save(featureCall).Error
				assert.NoError(t, err)

				return workspace.Uuid
			},
			expectError: false,
			validate: func(t *testing.T, result *FeatureCall) {
				assert.NotNil(t, result)
				assert.Equal(t, workspace.Uuid, result.WorkspaceID)
				assert.Equal(t, "https://second.example.com/api", result.URL)
			},
		},
		{
			name: "Unicode workspace ID",
			setup: func() string {
				return "测试workspace"
			},
			expectError: true,
			errorMsg:    "record not found",
		},
		{
			name: "Very long workspace ID",
			setup: func() string {
				return strings.Repeat("a", 1000)
			},
			expectError: true,
			errorMsg:    "record not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestDB.db.Where("1 = 1").Delete(&FeatureCall{})

			workspaceID := tt.setup()

			result, err := TestDB.GetFeatureCallByWorkspaceID(workspaceID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestDeleteFeatureCall(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	DeleteAllFeatureCalls()
	TestDB.db.Exec("DELETE FROM workspaces")

	workspace := Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace-" + uuid.New().String(),
		OwnerPubKey: "test-pubkey",
	}
	result := TestDB.db.Create(&workspace)
	assert.NoError(t, result.Error)

	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		errorMsg    string
		validate    func(t *testing.T)
	}{
		{
			name: "Successfully delete existing feature call",
			setup: func() string {
				featureCall := &FeatureCall{
					ID:          uuid.New(),
					WorkspaceID: workspace.Uuid,
					URL:         "https://example.com/api",
				}
				err := TestDB.db.Create(featureCall).Error
				assert.NoError(t, err)
				assert.NotNil(t, featureCall)

				created, err := TestDB.GetFeatureCallByWorkspaceID(workspace.Uuid)
				assert.NoError(t, err)
				assert.NotNil(t, created)
				
				return workspace.Uuid
			},
			expectError: false,
			validate: func(t *testing.T) {
				result, err := TestDB.GetFeatureCallByWorkspaceID(workspace.Uuid)
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "record not found", err.Error())
			},
		},
		{
			name: "Empty workspace ID",
			setup: func() string {
				return ""
			},
			expectError: true,
			errorMsg:    "workspace_id is required",
		},
		{
			name: "Non-existent workspace ID",
			setup: func() string {
				return uuid.New().String()
			},
			expectError: true,
			errorMsg:    "feature call not found",
		},
		{
			name: "Delete already deleted feature call",
			setup: func() string {
				featureCall := &FeatureCall{
					ID:          uuid.New(),
					WorkspaceID: workspace.Uuid,
					URL:         "https://example.com/api",
				}
				err := TestDB.db.Create(featureCall).Error
				assert.NoError(t, err)
				
				err = TestDB.DeleteFeatureCall(workspace.Uuid)
				assert.NoError(t, err)
				
				return workspace.Uuid
			},
			expectError: true,
			errorMsg:    "feature call not found",
		},
		{
			name: "SQL injection attempt",
			setup: func() string {
				return "'; DROP TABLE feature_calls; --"
			},
			expectError: true,
			errorMsg:    "feature call not found",
		},
		{
			name: "Delete with special characters in workspace ID",
			setup: func() string {
				specialID := "test!@#$%^&*()"
				featureCall := &FeatureCall{
					ID:          uuid.New(),
					WorkspaceID: specialID,
					URL:         "https://example.com/api",
				}
				TestDB.db.Create(featureCall)
				return specialID
			},
			expectError: false,
			validate: func(t *testing.T) {
				result, err := TestDB.GetFeatureCallByWorkspaceID("test!@#$%^&*()")
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "record not found", err.Error())
			},
		},
		{
			name: "Delete with very long workspace ID",
			setup: func() string {
				longID := strings.Repeat("a", 255)
				featureCall := &FeatureCall{
					ID:          uuid.New(),
					WorkspaceID: longID,
					URL:         "https://example.com/api",
				}
				TestDB.db.Create(featureCall)
				return longID
			},
			expectError: false,
			validate: func(t *testing.T) {
				result, err := TestDB.GetFeatureCallByWorkspaceID(strings.Repeat("a", 255))
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, "record not found", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeleteAllFeatureCalls()

			workspaceID := tt.setup()

			err := TestDB.DeleteFeatureCall(workspaceID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t)
				}
			}
		})
	}
}
