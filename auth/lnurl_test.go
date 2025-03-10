package auth

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestEncodeLNURL(t *testing.T) {
	tests := []struct {
		name          string
		host          string
		expectedError bool
	}{
		{
			name:          "Standard Hostname",
			host:          "example.com",
			expectedError: false,
		},
		{
			name:          "Localhost Hostname",
			host:          "localhost",
			expectedError: false,
		},
		{
			name:          "Empty Hostname",
			host:          "",
			expectedError: false,
		},
		{
			name:          "Hostname with Special Characters",
			host:          "example.com/path?query=1",
			expectedError: false,
		},
		{
			name:          "Invalid Hostname",
			host:          "invalid_host",
			expectedError: false,
		},
		{
			name:          "Very Long Hostname",
			host:          strings.Repeat("a", 1000) + ".com",
			expectedError: false,
		},
		{
			name:          "Hostname with Port",
			host:          "example.com:8080",
			expectedError: false,
		},
		{
			name:          "Hostname with Subdomain",
			host:          "sub.example.com",
			expectedError: false,
		},
		{
			name:          "Hostname with Trailing Slash",
			host:          "example.com/",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EncodeLNURL(tt.host)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result.Encode)
				assert.Len(t, result.K1, 64)
			}
		})
	}
}
