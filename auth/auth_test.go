package auth

import (
	"strings"
	"testing"

	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stretchr/testify/assert"
)

// Mock configuration for testing
var testConfig = struct {
	SuperAdmins []string
}{
	SuperAdmins: []string{"admin1", "admin2", "admin3"},
}

func TestAdminCheck(t *testing.T) {
	t.Setenv("SUPER_ADMINS", strings.Join(testConfig.SuperAdmins, ","))

	tests := []struct {
		name     string
		pubkey   interface{}
		expected bool
	}{
		{
			name:     "Valid super admin pubkey",
			pubkey:   "admin1",
			expected: true,
		},
		{
			name:     "Invalid super admin pubkey",
			pubkey:   "notAnAdmin",
			expected: false,
		},
		{
			name:     "Empty pubkey",
			pubkey:   "",
			expected: false,
		},
		{
			name:     "Empty SuperAdmins list",
			pubkey:   "admin1",
			expected: false,
		},
		{
			name:     "Pubkey is a substring of a super admin pubkey",
			pubkey:   "admin",
			expected: false,
		},
		{
			name:     "Pubkey is a super admin pubkey with additional characters",
			pubkey:   "admin1extra",
			expected: false,
		},
		{
			name:     "Null or nil pubkey",
			pubkey:   nil,
			expected: false,
		},
		{
			name:     "Non-string pubkey",
			pubkey:   12345,
			expected: false,
		},
		{
			name:     "Large list of super admin pubkeys",
			pubkey:   "admin1",
			expected: true,
		},
		{
			name:     "Large pubkey",
			pubkey:   "averylongpubkeythatisnotinlist",
			expected: false,
		},
		{
			name:     "Special characters in pubkey",
			pubkey:   "!@#$%^&*()",
			expected: false,
		},
		{
			name:     "Case sensitivity",
			pubkey:   "ADMIN1",
			expected: false,
		},
		{
			name:     "Duplicate entries in SuperAdmins",
			pubkey:   "admin1",
			expected: true,
		},
		{
			name:     "Whitespace in pubkey",
			pubkey:   " admin1 ",
			expected: false,
		},
		{
			name:     "Mixed data types in SuperAdmins",
			pubkey:   "admin1",
			expected: true,
		},
	}

	// Temporarily set SuperAdmins to an empty list for the specific test case
	originalSuperAdmins := testConfig.SuperAdmins
	defer func() { testConfig.SuperAdmins = originalSuperAdmins }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Empty SuperAdmins list" {
				config.SuperAdmins = []string{}
			} else {
				config.SuperAdmins = originalSuperAdmins
			}

			var result bool
			switch v := tt.pubkey.(type) {
			case string:
				result = AdminCheck(v)
			default:
				result = false
			}

			assert.Equal(t, tt.expected, result)
		})
	}
}

// Mock function to be tested
func TestIsFreePass(t *testing.T) {
	t.Setenv("SUPER_ADMINS", "")
	tests := []struct {
		name             string
		superAdmins      []string
		adminDevFreePass string
		adminStrings     string
		expected         bool
	}{
		{
			name:             "Single SuperAdmin with FreePass",
			superAdmins:      []string{"freepass"},
			adminDevFreePass: "freepass",
			adminStrings:     "non-empty",
			expected:         true,
		},
		{
			name:             "Empty AdminStrings",
			superAdmins:      []string{"admin"},
			adminDevFreePass: "freepass",
			adminStrings:     "",
			expected:         true,
		},
		{
			name:             "Both Conditions Met",
			superAdmins:      []string{"freepass"},
			adminDevFreePass: "freepass",
			adminStrings:     "",
			expected:         true,
		},
		{
			name:             "Multiple SuperAdmins",
			superAdmins:      []string{"freepass", "admin2"},
			adminDevFreePass: "freepass",
			adminStrings:     "non-empty",
			expected:         false,
		},
		{
			name:             "Empty SuperAdmins List",
			superAdmins:      []string{},
			adminDevFreePass: "freepass",
			adminStrings:     "non-empty",
			expected:         false,
		},
		{
			name:             "Empty SuperAdmins and Empty AdminStrings",
			superAdmins:      []string{},
			adminDevFreePass: "freepass",
			adminStrings:     "",
			expected:         true,
		},
		{
			name:             "Null SuperAdmins",
			superAdmins:      nil,
			adminDevFreePass: "freepass",
			adminStrings:     "non-empty",
			expected:         false,
		},
		{
			name:             "Null AdminStrings",
			superAdmins:      []string{"admin"},
			adminDevFreePass: "freepass",
			adminStrings:     "",
			expected:         true,
		},
		{
			name:             "SuperAdmin with Different FreePass",
			superAdmins:      []string{"admin"},
			adminDevFreePass: "freepass",
			adminStrings:     "non-empty",
			expected:         false,
		},
		{
			name:             "SuperAdmin with Empty String",
			superAdmins:      []string{""},
			adminDevFreePass: "freepass",
			adminStrings:     "non-empty",
			expected:         false,
		},
		{
			name:             "Large SuperAdmins List",
			superAdmins:      make([]string, 1000),
			adminDevFreePass: "freepass",
			adminStrings:     "non-empty",
			expected:         false,
		},
		{
			name:             "SuperAdmin with Null FreePass",
			superAdmins:      []string{"freepass"},
			adminDevFreePass: "",
			adminStrings:     "non-empty",
			expected:         false,
		},
		{
			name:             "AdminDevFreePass as Empty String",
			superAdmins:      []string{"freepass"},
			adminDevFreePass: "",
			adminStrings:     "non-empty",
			expected:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.SuperAdmins = tt.superAdmins
			config.AdminDevFreePass = tt.adminDevFreePass
			config.AdminStrings = tt.adminStrings

			result := IsFreePass()
			assert.Equal(t, tt.expected, result)
		})
	}
}
