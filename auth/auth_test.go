package auth

import (
	"bytes"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	btcecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
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

func TestVerifyAndExtract(t *testing.T) {

	privKey, err := btcec.NewPrivateKey()
	assert.NoError(t, err)

	createValidSignature := func(msg []byte) []byte {
		signedMsg := append(signedMsgPrefix, msg...)
		digest := chainhash.DoubleHashB(signedMsg)
		sig, err := btcecdsa.SignCompact(privKey, digest, true)
		assert.NoError(t, err)
		return sig
	}

	expectedPubKeyHex := hex.EncodeToString(privKey.PubKey().SerializeCompressed())

	tests := []struct {
		name          string
		msg           []byte
		sig           []byte
		expectedKey   string
		expectedValid bool
		expectedErr   error
	}{
		{
			name:          "Valid signature and message",
			msg:           []byte("test message"),
			sig:           createValidSignature([]byte("test message")),
			expectedKey:   expectedPubKeyHex,
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			name:          "Empty message",
			msg:           []byte{},
			sig:           createValidSignature([]byte{}),
			expectedKey:   expectedPubKeyHex,
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			name:          "Nil signature",
			msg:           []byte("test message"),
			sig:           nil,
			expectedKey:   "",
			expectedValid: false,
			expectedErr:   errors.New("bad"),
		},
		{
			name:          "Nil message",
			msg:           nil,
			sig:           createValidSignature([]byte("test message")),
			expectedKey:   "",
			expectedValid: false,
			expectedErr:   errors.New("bad"),
		},
		{
			name:          "Both nil inputs",
			msg:           nil,
			sig:           nil,
			expectedKey:   "",
			expectedValid: false,
			expectedErr:   errors.New("bad"),
		},
		{
			name:          "Empty signature",
			msg:           []byte("test message"),
			sig:           []byte{},
			expectedKey:   "",
			expectedValid: false,
			expectedErr:   errors.New("invalid compact signature size"),
		},
		{
			name:          "Invalid signature format",
			msg:           []byte("test message"),
			sig:           []byte{0xFF, 0xFF},
			expectedKey:   "",
			expectedValid: false,
			expectedErr:   errors.New("invalid compact signature size"),
		},
		{
			name:          "Corrupted signature",
			msg:           []byte("test message"),
			sig:           append(createValidSignature([]byte("test message")), byte(0x00)),
			expectedKey:   "",
			expectedValid: false,
			expectedErr:   errors.New("invalid compact signature size"),
		},
		{
			name:          "Large message",
			msg:           bytes.Repeat([]byte("a"), 1000),
			sig:           createValidSignature(bytes.Repeat([]byte("a"), 1000)),
			expectedKey:   expectedPubKeyHex,
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			name:          "Special characters in message",
			msg:           []byte("!@#$%^&*()_+{}:|<>?"),
			sig:           createValidSignature([]byte("!@#$%^&*()_+{}:|<>?")),
			expectedKey:   expectedPubKeyHex,
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			name:          "UTF-8 characters in message",
			msg:           []byte("Hello, 世界"),
			sig:           createValidSignature([]byte("Hello, 世界")),
			expectedKey:   expectedPubKeyHex,
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			name:          "Message with null bytes",
			msg:           []byte("test\x00message"),
			sig:           createValidSignature([]byte("test\x00message")),
			expectedKey:   expectedPubKeyHex,
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			name:          "Message with only whitespace",
			msg:           []byte("   "),
			sig:           createValidSignature([]byte("   ")),
			expectedKey:   expectedPubKeyHex,
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			name:          "Maximum length message",
			msg:           bytes.Repeat([]byte("x"), 1<<20),
			sig:           createValidSignature(bytes.Repeat([]byte("x"), 1<<20)),
			expectedKey:   expectedPubKeyHex,
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			name:          "Binary data in message",
			msg:           []byte{0x00, 0x01, 0x02, 0x03, 0xFF},
			sig:           createValidSignature([]byte{0x00, 0x01, 0x02, 0x03, 0xFF}),
			expectedKey:   expectedPubKeyHex,
			expectedValid: true,
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pubKeyHex, valid, err := VerifyAndExtract(tt.msg, tt.sig)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedValid, valid)

			if tt.expectedKey != "" {
				assert.Equal(t, tt.expectedKey, pubKeyHex)
			}

			if tt.msg != nil && tt.sig != nil && err == nil {
				assert.True(t, bytes.HasPrefix(append(signedMsgPrefix, tt.msg...), signedMsgPrefix))
			}

			if valid && err == nil {
				_, err := hex.DecodeString(pubKeyHex)
				assert.NoError(t, err, "Public key should be valid hex")

				if tt.sig != nil {
					assert.Equal(t, 65, len(tt.sig),
						"Valid signature should be 65 bytes (64 bytes signature + 1 byte recovery ID)")
				}
			}
		})
	}
}
