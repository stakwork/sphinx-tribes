package db

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *mockDatabase) AverageCompletedTime(r PaymentDateRange, workspace string) uint {
	paidList := m.CompletedDifference(r, workspace)
	paidCount := int64(len(paidList))
	var paidSum uint
	for _, diff := range paidList {
		paidSum += uint(math.Round(diff.Diff))
	}
	return CalculateAverageDays(paidCount, paidSum)
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

func TestAverageCompletedTime(t *testing.T) {

	tests := []struct {
		name      string
		dateRange PaymentDateRange
		workspace string
		mockData  []DateDifference
		expected  uint
	}{
		{
			name: "Standard Case - Multiple Entries",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "workspace1",
			mockData:  []DateDifference{{Diff: 86400}, {Diff: 172800}},
			expected:  2,
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
			name: "Single Entry",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "workspace3",
			mockData:  []DateDifference{{Diff: 86400}},
			expected:  1,
		},
		{
			name: "Very Small Time Difference",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "workspace4",
			mockData:  []DateDifference{{Diff: 3600}},
			expected:  0,
		},
		{
			name: "Large Time Difference",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "workspace5",
			mockData:  []DateDifference{{Diff: 864000}},
			expected:  10,
		},
		{
			name: "Mixed Time Differences",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "workspace6",
			mockData: []DateDifference{
				{Diff: 86400},
				{Diff: 172800},
				{Diff: 259200},
			},
			expected: 2,
		},
		{
			name: "No Workspace Specified",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "",
			mockData:  []DateDifference{{Diff: 86400}},
			expected:  1,
		},
		{
			name: "Boundary Values - Zero",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "workspace7",
			mockData:  []DateDifference{{Diff: 0}},
			expected:  0,
		},
		{
			name: "Boundary Values - Very Large",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "workspace8",
			mockData:  []DateDifference{{Diff: 8640000}},
			expected:  100,
		},
		{
			name: "Invalid Date Range",
			dateRange: PaymentDateRange{
				StartDate: "2023-12-31",
				EndDate:   "2023-01-01",
			},
			workspace: "workspace9",
			mockData:  []DateDifference{},
			expected:  0,
		},
		{
			name: "Large Dataset",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "workspace10",
			mockData: func() []DateDifference {
				diffs := make([]DateDifference, 1000)
				for i := range diffs {
					diffs[i] = DateDifference{Diff: 86400}
				}
				return diffs
			}(),
			expected: 1,
		},
		{
			name: "Fractional Days",
			dateRange: PaymentDateRange{
				StartDate: "2023-01-01",
				EndDate:   "2023-12-31",
			},
			workspace: "workspace11",
			mockData: []DateDifference{
				{Diff: 129600},
				{Diff: 172800},
			},
			expected: 2,
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

			result := mockDB.AverageCompletedTime(tt.dateRange, tt.workspace)
			assert.Equal(t, tt.expected, result)
		})
	}
}

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) TotalPeopleByPeriod(r PaymentDateRange) int64 {
	args := m.Called(r)
	return args.Get(0).(int64)
}

func TestTotalPeopleByPeriod(t *testing.T) {
	mockDB := new(MockDatabase)

	tests := []struct {
		name          string
		input         PaymentDateRange
		mockReturn    int64
		expectPanic   bool
		expectedValue int64
	}{
		{
			name:          "Standard Input with Valid Return Value",
			input:         PaymentDateRange{StartDate: "2023-01-01", EndDate: "2023-12-31"},
			mockReturn:    100,
			expectPanic:   false,
			expectedValue: 100,
		},
		{
			name:          "Empty Date Range",
			input:         PaymentDateRange{StartDate: "", EndDate: ""},
			mockReturn:    0,
			expectPanic:   false,
			expectedValue: 0,
		},
		{
			name:          "Single Day Date Range",
			input:         PaymentDateRange{StartDate: "2023-01-01", EndDate: "2023-01-01"},
			mockReturn:    1,
			expectPanic:   false,
			expectedValue: 1,
		},
		{
			name:          "Maximum Date Range",
			input:         PaymentDateRange{StartDate: "1900-01-01", EndDate: "2100-12-31"},
			mockReturn:    10000,
			expectPanic:   false,
			expectedValue: 10000,
		},
		{
			name:          "Invalid Date Range (End Date Before Start Date)",
			input:         PaymentDateRange{StartDate: "2023-12-31", EndDate: "2023-01-01"},
			mockReturn:    0,
			expectPanic:   false,
			expectedValue: 0,
		},
		{
			name:          "Invalid Date Format",
			input:         PaymentDateRange{StartDate: "2023-31-12", EndDate: "2023-01-01"},
			mockReturn:    0,
			expectPanic:   false,
			expectedValue: 0,
		},
		{
			name:          "Large Volume of Data",
			input:         PaymentDateRange{StartDate: "2000-01-01", EndDate: "2023-12-31"},
			mockReturn:    5000,
			expectPanic:   false,
			expectedValue: 5000,
		},
		{
			name:          "Leap Year Date Range",
			input:         PaymentDateRange{StartDate: "2020-02-28", EndDate: "2020-03-01"},
			mockReturn:    2,
			expectPanic:   false,
			expectedValue: 2,
		},
		{
			name:          "Boundary Date Range",
			input:         PaymentDateRange{StartDate: "2023-12-31", EndDate: "2024-01-01"},
			mockReturn:    1,
			expectPanic:   false,
			expectedValue: 1,
		},
		{
			name:          "Non-Existent Date Range",
			input:         PaymentDateRange{StartDate: "2023-02-30", EndDate: "2023-03-01"},
			mockReturn:    0,
			expectPanic:   false,
			expectedValue: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() {
					mockDB.TotalPeopleByPeriod(tt.input)
				})
			} else {
				mockDB.On("TotalPeopleByPeriod", tt.input).Return(tt.mockReturn)
				result := mockDB.TotalPeopleByPeriod(tt.input)
				assert.Equal(t, tt.expectedValue, result)
				mockDB.AssertExpectations(t)
			}
		})
	}
}
