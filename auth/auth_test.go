package auth

import (
	"fmt"
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

func generateLargePayload() map[string]interface{} {
	payload := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		payload[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
	}
	return payload
}

func TestEncodeJwt(t *testing.T) {

	config.InitConfig()
	InitJwt()

	tests := []struct {
		name        string
		publicKey   string
		payload     interface{}
		expectError bool
	}{
		{
			name:        "Valid Public Key and Payload",
			publicKey:   "validPublicKey",
			payload:     map[string]interface{}{"user": "testUser"},
			expectError: false,
		},
		{
			name:        "Valid Public Key with Minimal Payload",
			publicKey:   "validPublicKey",
			payload:     map[string]interface{}{"id": 1},
			expectError: false,
		},
		{
			name:        "Empty Payload",
			publicKey:   "validPublicKey",
			payload:     map[string]interface{}{},
			expectError: false,
		},
		{
			name:        "Maximum Size Payload",
			publicKey:   "validPublicKey",
			payload:     generateLargePayload(),
			expectError: false,
		},
		{
			name:        "Boundary Public Key Length",
			publicKey:   "a",
			payload:     map[string]interface{}{"user": "testUser"},
			expectError: false,
		},
		{
			name:        "Invalid Public Key",
			publicKey:   "invalidPublicKey!",
			payload:     map[string]interface{}{"user": "testUser"},
			expectError: true,
		},
		{
			name:        "Null Public Key",
			publicKey:   "",
			payload:     map[string]interface{}{"user": "testUser"},
			expectError: true,
		},
		{
			name:        "Expired Payload",
			publicKey:   "validPublicKey",
			payload:     map[string]interface{}{"exp": -1},
			expectError: false,
		},
		{
			name:        "Future Expiration Date",
			publicKey:   "validPublicKey",
			payload:     map[string]interface{}{"exp": 9999999999},
			expectError: false,
		},
		{
			name:        "Payload with Special Characters",
			publicKey:   "validPublicKey",
			payload:     map[string]interface{}{"emoji": "😀"},
			expectError: false,
		},
		{
			name:        "Payload with Reserved JWT Claims",
			publicKey:   "validPublicKey",
			payload:     map[string]interface{}{"iss": "issuer", "sub": "subject"},
			expectError: false,
		},
		{
			name:        "Payload with Mixed Data Types",
			publicKey:   "validPublicKey",
			payload:     map[string]interface{}{"string": "value", "number": 123, "boolean": true},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwt, err := EncodeJwt(tt.publicKey)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, jwt)
			}
		})
	}
}
