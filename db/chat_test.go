package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetChatsForWorkspace(t *testing.T) {
	InitTestDB()

	chat := Chat{
		WorkspaceID: "workspace123", Status: ActiveStatus,
	}

	TestDB.db.Create(&chat)

	tests := []struct {
		name        string
		workspaceID string
		mockChats   []Chat
		mockError   error
		expected    []Chat
		expectError bool
	}{
		{
			name:        "Basic Functionality",
			workspaceID: "workspace123",
			expected: []Chat{
				{WorkspaceID: "workspace123", Status: ActiveStatus},
			},
			expectError: false,
		},
		{
			name:        "No Chats for Workspace",
			workspaceID: "emptyWorkspace",
			expected:    []Chat{},
			expectError: false,
		},
		{
			name:        "Invalid Workspace ID",
			workspaceID: "nonExistentWorkspace",
			expected:    []Chat{},
			expectError: false,
		},
		{
			name:        "Empty Workspace ID",
			workspaceID: "",
			expected:    []Chat{},
			expectError: true,
		},
		{
			name:        "Null Workspace ID",
			workspaceID: "",
			expected:    []Chat{},
			expectError: true,
		},
		{
			name:        "Special Characters in Workspace ID",
			workspaceID: "special!@#Workspace",
			expected:    []Chat{},
			expectError: false,
		},
		{
			name:        "SQL Injection Attempt",
			workspaceID: "workspace123'; DROP TABLE chats; --",
			expected:    []Chat{},
			expectError: false,
		},
		{
			name:        "Case Sensitivity",
			workspaceID: "Workspace123",
			expected:    []Chat{},
			expectError: false,
		},
		{
			name:        "Unicode Characters in Workspace ID",
			workspaceID: "工作区123",
			expected:    []Chat{},
			expectError: false,
		},
		{
			name:        "Maximum Length Workspace ID",
			workspaceID: "asdfasfdaasdfasdfasdfasdfasdfasdfasdfjlkajsldkfjalsdflakjsdlkfjalsdkjfal",
			expected:    []Chat{},
			expectError: false,
		},
		{
			name:        "Minimum Length Workspace ID",
			workspaceID: "a",
			expected:    []Chat{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			chats, err := TestDB.GetChatsForWorkspace(tt.workspaceID, "")

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if len(tt.expected) > 0 {
					assert.Equal(t, tt.expected[0].WorkspaceID, chats[0].WorkspaceID)
				} else {
					assert.Equal(t, tt.expected, chats)
				}
			}
		})
	}
}
