package auth

import (
	"math"
	"testing"
	"time"

	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stretchr/testify/assert"
)

func TestInitJwt(t *testing.T) {

	config.InitConfig()
	InitJwt()

	if TokenAuth == nil {
		t.Error("Could not init JWT")
	} else {
		t.Log("JWT inited successfully")
	}
}

func TestExpireInsHours(t *testing.T) {
	tests := []struct {
		name     string
		hours    int
		expected int64
	}{
		{
			name:     "Basic Functionality",
			hours:    24,
			expected: time.Now().Add(24 * time.Hour).Unix(),
		},
		{
			name:     "Zero Hours",
			hours:    0,
			expected: time.Now().Unix(),
		},
		{
			name:     "Negative Hours",
			hours:    -5,
			expected: time.Now().Add(-5 * time.Hour).Unix(),
		},
		{
			name:     "Maximum Integer Hours",
			hours:    math.MaxInt32 / 1000,
			expected: time.Now().Add(time.Duration(math.MaxInt32/1000) * time.Hour).Unix(),
		},
		{
			name:     "Minimum Integer Hours",
			hours:    math.MinInt32 / 1000,
			expected: time.Now().Add(time.Duration(math.MinInt32/1000) * time.Hour).Unix(),
		},
		{
			name:     "Small Positive Hours",
			hours:    1,
			expected: time.Now().Add(time.Hour).Unix(),
		},
		{
			name:     "Large Positive Hours",
			hours:    10000,
			expected: time.Now().Add(10000 * time.Hour).Unix(),
		},
		{
			name:     "Boundary Condition at 1 Hour",
			hours:    1,
			expected: time.Now().Add(time.Hour).Unix(),
		},
		{
			name:     "Boundary Condition at Maximum Duration",
			hours:    int(math.MaxInt64 / int64(time.Hour)),
			expected: time.Now().Add(time.Duration(math.MaxInt64/int64(time.Hour)) * time.Hour).Unix(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpireInHours(tt.hours)

			assert.InDelta(t, tt.expected, result, 2,
				"ExpireInHours(%d) = %d; want approximately %d",
				tt.hours, result, tt.expected)
		})
	}
}
