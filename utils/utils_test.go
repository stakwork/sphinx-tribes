package utils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildV2KeysendBodyData(t *testing.T) {
	tests := []struct {
		name           string
		amount         uint
		receiverPubkey string
		routeHint      string
		memo           string
		expected       string
		expectError    bool
	}{
		{
			name:           "Standard Input with Route Hint",
			amount:         100,
			receiverPubkey: "abcdef123456",
			routeHint:      "hint123",
			memo:           "Test transaction",
			expected:       `{"amt_msat": 100000, "dest": "abcdef123456", "route_hint": "hint123", "data": "Test transaction", "wait": true}`,
		},
		{
			name:           "Standard Input without Route Hint",
			amount:         100,
			receiverPubkey: "abcdef123456",
			routeHint:      "",
			memo:           "Test transaction",
			expected:       `{"amt_msat": 100000, "dest": "abcdef123456", "route_hint": "", "data": "Test transaction", "wait": true}`,
		},
		{
			name:           "Minimum Amount",
			amount:         0,
			receiverPubkey: "abcdef123456",
			routeHint:      "hint123",
			memo:           "Test transaction",
			expected:       `{"amt_msat": 0, "dest": "abcdef123456", "route_hint": "hint123", "data": "Test transaction", "wait": true}`,
		},
		{
			name:           "Empty Strings for All String Parameters",
			amount:         100,
			receiverPubkey: "",
			routeHint:      "",
			memo:           "",
			expected:       `{"amt_msat": 100000, "dest": "", "route_hint": "", "data": "", "wait": true}`,
		},
		{
			name:           "Large Amount",
			amount:         4294967295,
			receiverPubkey: "abcdef123456",
			routeHint:      "hint123",
			memo:           "Test transaction",
			expected:       `{"amt_msat": 4294967295000, "dest": "abcdef123456", "route_hint": "hint123", "data": "Test transaction", "wait": true}`,
		},
		{
			name:           "Long Strings",
			amount:         100,
			receiverPubkey: strings.Repeat("a", 1000),
			routeHint:      strings.Repeat("b", 1000),
			memo:           strings.Repeat("c", 1000),
			expected:       fmt.Sprintf(`{"amt_msat": 100000, "dest": "%s", "route_hint": "%s", "data": "%s", "wait": true}`, strings.Repeat("a", 1000), strings.Repeat("b", 1000), strings.Repeat("c", 1000)),
		},
		{
			name:           "Special Characters in Strings",
			amount:         100,
			receiverPubkey: "abc!@#123",
			routeHint:      "hint$%^",
			memo:           "Test &*() transaction",
			expected:       `{"amt_msat": 100000, "dest": "abc!@#123", "route_hint": "hint$%^", "data": "Test &*() transaction", "wait": true}`,
		},
		{
			name:           "Whitespace in Strings",
			amount:         100,
			receiverPubkey: "abcdef123456",
			routeHint:      "hint123  ",
			memo:           "  Test transaction  ",
			expected:       `{"amt_msat": 100000, "dest": "abcdef123456", "route_hint": "hint123", "data": "Test transaction", "wait": true}`,
		},
		{
			name:           "Non-ASCII Characters in Strings",
			amount:         100,
			receiverPubkey: "abcñ123",
			routeHint:      "hintü123",
			memo:           "Test transaction with emoji 😊",
			expected:       `{"amt_msat": 100000, "dest": "abcñ123", "route_hint": "hintü123", "data": "Test transaction with emoji 😊", "wait": true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildV2KeysendBodyData(tt.amount, tt.receiverPubkey, tt.routeHint, tt.memo)
			assert.JSONEq(t, tt.expected, result)
		})
	}
}
