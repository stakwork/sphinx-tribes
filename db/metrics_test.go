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

type mockDatabase struct {
	completedDifferenceFn func(PaymentDateRange, string) []DateDifference
}

func (m *mockDatabase) CompletedDifference(r PaymentDateRange, workspace string) []DateDifference {
	return m.completedDifferenceFn(r, workspace)
}

func (m *mockDatabase) CompletedDifferenceCount(r PaymentDateRange, workspace string) int64 {
	list := m.CompletedDifference(r, workspace)
	return int64(len(list))
}

func TestCompletedDifferenceCount(t *testing.T) {

	tests := []struct {
		name      string
		dateRange PaymentDateRange
		workspace string
		mockData  []DateDifference
		expected  int64
	}{
		{
			name: "Standard Case",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "workspace1",
			mockData:  []DateDifference{{Diff: 100}, {Diff: 200}, {Diff: 300}},
			expected:  3,
		},
		{
			name: "Empty Result",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "workspace2",
			mockData:  []DateDifference{},
			expected:  0,
		},
		{
			name: "No Workspace Specified",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "",
			mockData:  []DateDifference{{Diff: 100}, {Diff: 200}},
			expected:  2,
		},
		{
			name: "Single Day Range",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-01-01",
			},
			workspace: "workspace1",
			mockData:  []DateDifference{{Diff: 100}},
			expected:  1,
		},
		{
			name: "Large Dataset",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "workspace1",
			mockData:  make([]DateDifference, 1000),
			expected:  1000,
		},
		{
			name: "Invalid Date Range",
			dateRange: PaymentDateRange{
				StartDate: "2023-12-31",
				EndDate:   "2023-01-01",
			},
			workspace: "workspace1",
			mockData:  []DateDifference{},
			expected:  0,
		},
		{
			name:      "Null Date Range",
			dateRange: PaymentDateRange{},
			workspace: "workspace1",
			mockData:  []DateDifference{},
			expected:  0,
		},
		{
			name: "Special Characters in Workspace",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "workspace!@#$%",
			mockData:  []DateDifference{{Diff: 100}},
			expected:  1,
		},
		{
			name: "Maximum Date Range",
			dateRange: PaymentDateRange{
				StartDate: "1970-01-01",
				EndDate:   "2099-12-31",
			},
			workspace: "workspace1",
			mockData:  []DateDifference{{Diff: 100}, {Diff: 200}, {Diff: 300}},
			expected:  3,
		},
		{
			name: "Boundary Values",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-01-01",
			},
			workspace: "",
			mockData:  []DateDifference{{Diff: 0}, {Diff: math.MaxFloat64}},
			expected:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := &mockDatabase{
				completedDifferenceFn: func(r PaymentDateRange, workspace string) []DateDifference {
					assert.Equal(t, tt.dateRange, r)
					assert.Equal(t, tt.workspace, workspace)
					return tt.mockData
				},
			}

			result := mockDB.CompletedDifferenceCount(tt.dateRange, tt.workspace)

			assert.Equal(t, tt.expected, result)
		})
	}
}
