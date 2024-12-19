package db

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateAverageDays(t *testing.T) {
	const SecondsToDateConversion = 86400
	tests := []struct {
		name      string
		paidCount int64
		paidSum   uint
		expected  uint
	}{
		{
			name:      "Standard Input",
			paidCount: 10,
			paidSum:   1000,
			expected:  uint(math.Round(float64(1000/10) / float64(SecondsToDateConversion))),
		},
		{
			name:      "Zero Paid Count",
			paidCount: 0,
			paidSum:   1000,
			expected:  0,
		},
		{
			name:      "Zero Paid Sum",
			paidCount: 10,
			paidSum:   0,
			expected:  0,
		},
		{
			name:      "Both Zero",
			paidCount: 0,
			paidSum:   0,
			expected:  0,
		},
		{
			name:      "Minimum Positive Values",
			paidCount: 1,
			paidSum:   1,
			expected:  0,
		},
		{
			name:      "Large Paid Count and Sum",
			paidCount: 1000000,
			paidSum:   1000000000,
			expected:  uint(math.Round(float64(1000000000/1000000) / float64(SecondsToDateConversion))),
		},
		{
			name:      "Paid Count Greater than Paid Sum",
			paidCount: 1000,
			paidSum:   500,
			expected:  0,
		},
		{
			name:      "Paid Sum Exactly Divisible by Paid Count",
			paidCount: 5,
			paidSum:   100,
			expected:  uint(math.Round(float64(100/5) / float64(SecondsToDateConversion))),
		},
		{
			name:      "Paid Sum Not Divisible by Paid Count",
			paidCount: 3,
			paidSum:   100,
			expected:  uint(math.Round(float64(100/3) / float64(SecondsToDateConversion))),
		},
		{
			name:      "Zero Result Due to Seconds Conversion",
			paidCount: 10,
			paidSum:   100,
			expected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateAverageDays(tt.paidCount, tt.paidSum)
			assert.Equal(t, tt.expected, result)
		})
	}
}
