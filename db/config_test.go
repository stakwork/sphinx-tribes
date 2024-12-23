package db

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRolesCheck_UserHasRole(t *testing.T) {
	// Mock user roles
	userRoles := []WorkspaceUserRoles{
		{Role: "ADD BOUNTY", OwnerPubKey: "user1", WorkspaceUuid: "org1", Created: &time.Time{}},
	}

	// Role to check
	roleToCheck := "ADD BOUNTY"

	// Call the function
	result := RolesCheck(userRoles, roleToCheck)

	// Assert that it returns true
	if !result {
		t.Errorf("Expected RolesCheck to return true for user with role %s, got false", roleToCheck)
	}
}

func TestRolesCheck_UserDoesNotHaveRole(t *testing.T) {
	// Mock user roles
	userRoles := []WorkspaceUserRoles{
		{Role: "DELETE BOUNTY", OwnerPubKey: "user2", WorkspaceUuid: "org1", Created: &time.Time{}},
	}

	// Role to check
	roleToCheck := "ADD BOUNTY"

	// Call the function
	result := RolesCheck(userRoles, roleToCheck)

	// Assert that it returns false
	if result {
		t.Errorf("Expected RolesCheck to return false for user without role %s, got true", roleToCheck)
	}
}

func TestCheckUser(t *testing.T) {
	userRoles := []WorkspaceUserRoles{
		{OwnerPubKey: "userPublicKey"},
	}

	// if in the user roles, one of the owner_pubkey belongs to the user return true else return false
	assert.True(t, CheckUser(userRoles, "userPublicKey"))
	assert.False(t, CheckUser(userRoles, "anotherPublicKey"))
}

func TestUserHasAccess(t *testing.T) {
	mockGetWorkspaceByUuid := func(uuid string) Workspace {
		return Workspace{
			Uuid:        uuid,
			OwnerPubKey: "org_admin",
		}
	}

	mockGetUserRoles := func(uuid string, pubkey string) []WorkspaceUserRoles {
		return []WorkspaceUserRoles{
			{Role: "ADD BOUNTY", OwnerPubKey: pubkey, WorkspaceUuid: uuid, Created: &time.Time{}},
		}
	}

	mockDB := &gorm.DB{}

	databaseConfig := NewDatabaseConfig(mockDB)
	databaseConfig.getWorkspaceByUuid = mockGetWorkspaceByUuid
	databaseConfig.getUserRoles = mockGetUserRoles

	t.Run("Should test that if the user is the admin of an workspace returns true", func(t *testing.T) {
		result := databaseConfig.UserHasAccess("org_admin", "workspace_uuid", "ADD BOUNTY")

		// Assert that it returns true since the user is the org admin
		if !result {
			t.Errorf("Expected UserHasAccess to return true for workspace admin, got false")
		}
	})

	t.Run("Should test that if the user is not the workspace admin, and the user has the required role it should return true", func(t *testing.T) {
		result := databaseConfig.UserHasAccess("user_pubkey", "workspace_uuid", "ADD BOUNTY")

		// Assert that it returns true since the user has the required role
		if !result {
			t.Errorf("Expected UserHasAccess to return true for user with required role, got false")
		}
	})

	t.Run("Should test that if the user is not the workspace admin, and the user has not the required role it should return false", func(t *testing.T) {
		result := databaseConfig.UserHasAccess("user_pubkey", "workspace_uuid", "DELETE BOUNTY")

		// Assert that it returns false since the user does not have the required role
		if result {
			t.Errorf("Expected UserHasAccess to return false for user without required role, got true")
		}
	})
}

func TestUserHasManageBountyRoles(t *testing.T) {
	mockGetWorkspaceByUuid := func(uuid string) Workspace {
		return Workspace{
			Uuid:        uuid,
			OwnerPubKey: "org_admin",
		}
	}

	mockGetUserRoles := func(uuid string, pubkey string) []WorkspaceUserRoles {
		if uuid == "workspace_uuid" {
			return []WorkspaceUserRoles{
				{Role: "ADD BOUNTY", OwnerPubKey: pubkey, WorkspaceUuid: uuid, Created: &time.Time{}},
			}
		} else {
			return []WorkspaceUserRoles{
				{Role: "ADD BOUNTY", OwnerPubKey: pubkey, WorkspaceUuid: uuid, Created: &time.Time{}},
				{Role: "UPDATE BOUNTY", OwnerPubKey: pubkey, WorkspaceUuid: uuid, Created: &time.Time{}},
				{Role: "DELETE BOUNTY", OwnerPubKey: pubkey, WorkspaceUuid: uuid, Created: &time.Time{}},
				{Role: "PAY BOUNTY", OwnerPubKey: pubkey, WorkspaceUuid: uuid, Created: &time.Time{}},
			}
		}
	}

	mockDB := &gorm.DB{}

	databaseConfig := NewDatabaseConfig(mockDB)
	databaseConfig.getWorkspaceByUuid = mockGetWorkspaceByUuid
	databaseConfig.getUserRoles = mockGetUserRoles

	t.Run("Should test that if the user is the workspace admin return true", func(t *testing.T) {
		result := databaseConfig.UserHasManageBountyRoles("org_admin", "workspace_uuid")

		// Assert that it returns true since the user is the org admin
		assert.True(t, result, "Expected UserHasManageBountyRoles to return true for workspace admin")
	})

	t.Run("Should test that if the user has all bounty roles return true", func(t *testing.T) {
		result := databaseConfig.UserHasManageBountyRoles("user_pubkey", "workspace_uuid2")

		// Assert that it returns true since the user has all bounty roles
		assert.True(t, result, "Expected UserHasManageBountyRoles to return true for user with all bounty roles")
	})

	t.Run("Should test that if the user don't have all bounty roles return false.", func(t *testing.T) {
		result := databaseConfig.UserHasManageBountyRoles("user_pubkey", "workspace_uuid")

		// Assert that it returns false since the user does not have all bounty roles
		assert.False(t, result, "Expected UserHasManageBountyRoles to return false for user without all bounty roles")
	})
}

func TestProcessUpdateTicketsWithoutGroup(t *testing.T) {
	InitTestDB()

	// create person
	now := time.Now()

	person := Person{
		Uuid:        uuid.New().String(),
		OwnerPubKey: "testfeaturepubkeyProcess",
		OwnerAlias:  "testfeaturealiasProcess",
		Description: "testfeaturedescriptionProcess",
		Created:     &now,
		Updated:     &now,
		Deleted:     false,
	}

	// create person
	TestDB.CreateOrEditPerson(person)

	workspace := Workspace{
		Uuid:    uuid.New().String(),
		Name:    "Test tickets process space",
		Created: &now,
		Updated: &now,
	}

	// create workspace
	TestDB.CreateOrEditWorkspace(workspace)

	workspaceFeatures := WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test process feature",
		Brief:         "test get process brief",
		Requirements:  "Test get process requirements",
		Architecture:  "Test get process architecture",
		Url:           "Test get process url",
		Priority:      1,
		Created:       &now,
		Updated:       &now,
		CreatedBy:     "test",
		UpdatedBy:     "test",
	}

	// create WorkspaceFeatures
	TestDB.CreateOrEditFeature(workspaceFeatures)

	featurePhase := FeaturePhase{
		Uuid:        uuid.New().String(),
		FeatureUuid: workspaceFeatures.Uuid,
		Name:        "test get process feature phase",
		Priority:    1,
		Created:     &now,
		Updated:     &now,
	}

	// create FeaturePhase
	TestDB.CreateOrEditFeaturePhase(featurePhase)

	ticket := Tickets{
		UUID:        uuid.New(),
		FeatureUUID: workspaceFeatures.Uuid,
		PhaseUUID:   featurePhase.Uuid,
		Name:        "test get process ticket",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// create ticket
	TestDB.CreateOrEditTicket(&ticket)

	// process update tickets without group
	TestDB.ProcessUpdateTicketsWithoutGroup()

	// get tickets without group and assert that there is 0
	tickets, err := TestDB.GetTicketsWithoutGroup()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tickets))

	// get ticket and assert that the ticket group is the same as the ticket uuid
	ticket, err = TestDB.GetTicket(ticket.UUID.String())

	utils.Log.Info("tickets: %v", tickets)

	ticketUuid := ticket.UUID
	ticketAuthorID := "12345"
	ticketAuthor := Author("HUMAN")

	assert.NoError(t, err)
	assert.Equal(t, ticket.TicketGroup, &ticketUuid)
	assert.Equal(t, ticket.AuthorID, &ticketAuthorID)
	assert.Equal(t, ticket.Author, &ticketAuthor)
}

func TestGetRolesMap(t *testing.T) {

	originalRoles := ConfigBountyRoles
	defer func() {

		ConfigBountyRoles = originalRoles
	}()
	tests := []struct {
		name     string
		roles    []BountyRoles
		expected map[string]string
	}{
		{
			name: "Basic Functionality: Standard Input",
			roles: []BountyRoles{
				{Name: "ADD BOUNTY"},
				{Name: "UPDATE BOUNTY"},
				{Name: "DELETE BOUNTY"},
			},
			expected: map[string]string{
				"ADD BOUNTY":    "ADD BOUNTY",
				"UPDATE BOUNTY": "UPDATE BOUNTY",
				"DELETE BOUNTY": "DELETE BOUNTY",
			},
		},
		{
			name:     "Edge Case: Empty Input",
			roles:    []BountyRoles{},
			expected: map[string]string{},
		},
		{
			name: "Edge Case: Single Role",
			roles: []BountyRoles{
				{Name: "ADD BOUNTY"},
			},
			expected: map[string]string{
				"ADD BOUNTY": "ADD BOUNTY",
			},
		},
		{
			name: "Edge Case: Duplicate Role Names",
			roles: []BountyRoles{
				{Name: "ADD BOUNTY"},
				{Name: "ADD BOUNTY"},
			},
			expected: map[string]string{
				"ADD BOUNTY": "ADD BOUNTY",
			},
		},
		{
			name: "Special Case: Role Names with Special Characters",
			roles: []BountyRoles{
				{Name: "ROLE@SPECIAL"},
				{Name: "ROLE#TEST"},
				{Name: "ROLE!ADMIN"},
			},
			expected: map[string]string{
				"ROLE@SPECIAL": "ROLE@SPECIAL",
				"ROLE#TEST":    "ROLE#TEST",
				"ROLE!ADMIN":   "ROLE!ADMIN",
			},
		},
		{
			name: "Special Case: Role Names with Spaces",
			roles: []BountyRoles{
				{Name: "ADD USER ROLE"},
				{Name: "MANAGE USER ROLE"},
			},
			expected: map[string]string{
				"ADD USER ROLE":    "ADD USER ROLE",
				"MANAGE USER ROLE": "MANAGE USER ROLE",
			},
		},
		{
			name:     "Error Condition: Nil Input",
			roles:    nil,
			expected: map[string]string{},
		},
		{
			name: "Performance and Scale: Large Number of Roles",
			roles: func() []BountyRoles {
				roles := make([]BountyRoles, 1000)
				for i := 0; i < 1000; i++ {
					roles[i] = BountyRoles{Name: fmt.Sprintf("ROLE_%d", i)}
				}
				return roles
			}(),
			expected: func() map[string]string {
				expected := make(map[string]string)
				for i := 0; i < 1000; i++ {
					roleName := fmt.Sprintf("ROLE_%d", i)
					expected[roleName] = roleName
				}
				return expected
			}(),
		},
		{
			name: "Special Case: Role Names with Unicode Characters",
			roles: []BountyRoles{
				{Name: "管理员角色"},
				{Name: "用户角色"},
			},
			expected: map[string]string{
				"管理员角色": "管理员角色",
				"用户角色":  "用户角色",
			},
		},
		{
			name: "Special Case: Role Names with Numeric Characters",
			roles: []BountyRoles{
				{Name: "ROLE_123"},
				{Name: "ROLE_456"},
			},
			expected: map[string]string{
				"ROLE_123": "ROLE_123",
				"ROLE_456": "ROLE_456",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ConfigBountyRoles = tt.roles
			result := GetRolesMap()
			assert.Equal(t, tt.expected, result, "Maps should be equal")
			assert.Equal(t, len(tt.expected), len(result), "Map lengths should match")
			if tt.roles != nil {
				for _, role := range tt.roles {
					mappedRole, exists := result[role.Name]
					assert.True(t, exists, "Role should exist in map")
					assert.Equal(t, role.Name, mappedRole, "Role should be mapped to itself")
				}
			}
		})
	}
}
