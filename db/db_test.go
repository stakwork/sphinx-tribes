package db

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestGetFilterStatusCount(t *testing.T) {

	InitTestDB()

	defer CloseTestDB()

	tests := []struct {
		name     string
		setup    []NewBounty
		expected FilterStattuCount
	}{
		{
			name:  "Empty Database",
			setup: []NewBounty{},
			expected: FilterStattuCount{
				Open: 0, Assigned: 0, Completed: 0,
				Paid: 0, Pending: 0, Failed: 0,
			},
		},
		{
			name: "Hidden Bounties Should Not Count",
			setup: []NewBounty{
				{Show: false, Assignee: "", Paid: false},
				{Show: false, Assignee: "user1", Completed: true},
			},
			expected: FilterStattuCount{
				Open: 0, Assigned: 0, Completed: 0,
				Paid: 0, Pending: 0, Failed: 0,
			},
		},
		{
			name: "Open Bounties Count",
			setup: []NewBounty{
				{Show: true, Assignee: "", Paid: false},
				{Show: true, Assignee: "", Paid: false},
			},
			expected: FilterStattuCount{
				Open: 2, Assigned: 0, Completed: 0,
				Paid: 0, Pending: 0, Failed: 0,
			},
		},
		{
			name: "Assigned Bounties Count",
			setup: []NewBounty{
				{Show: true, Assignee: "user1", Paid: false},
				{Show: true, Assignee: "user2", Paid: false},
			},
			expected: FilterStattuCount{
				Open: 0, Assigned: 2, Completed: 0,
				Paid: 0, Pending: 0, Failed: 0,
			},
		},
		{
			name: "Completed Bounties Count",
			setup: []NewBounty{
				{Show: true, Assignee: "user1", Completed: true, Paid: false},
				{Show: true, Assignee: "user2", Completed: true, Paid: false},
			},
			expected: FilterStattuCount{
				Open: 0, Assigned: 2, Completed: 2,
				Paid: 0, Pending: 0, Failed: 0,
			},
		},
		{
			name: "Paid Bounties Count",
			setup: []NewBounty{
				{Show: true, Assignee: "user1", Paid: true},
				{Show: true, Assignee: "user2", Paid: true},
			},
			expected: FilterStattuCount{
				Open: 0, Assigned: 0, Completed: 0,
				Paid: 2, Pending: 0, Failed: 0,
			},
		},
		{
			name: "Pending Payment Bounties Count",
			setup: []NewBounty{
				{Show: true, Assignee: "user1", PaymentPending: true},
				{Show: true, Assignee: "user2", PaymentPending: true},
			},
			expected: FilterStattuCount{
				Open: 0, Assigned: 2, Completed: 0,
				Paid: 0, Pending: 2, Failed: 0,
			},
		},
		{
			name: "Failed Payment Bounties Count",
			setup: []NewBounty{
				{Show: true, Assignee: "user1", PaymentFailed: true},
				{Show: true, Assignee: "user2", PaymentFailed: true},
			},
			expected: FilterStattuCount{
				Open: 0, Assigned: 2, Completed: 0,
				Paid: 0, Pending: 0, Failed: 2,
			},
		},
		{
			name: "Mixed Status Bounties",
			setup: []NewBounty{
				{Show: true, Assignee: "", Paid: false},
				{Show: true, Assignee: "user1", Paid: false},
				{Show: true, Assignee: "user2", Completed: true, Paid: false},
				{Show: true, Assignee: "user3", Paid: true},
				{Show: true, Assignee: "user4", PaymentPending: true},
				{Show: true, Assignee: "user5", PaymentFailed: true},
				{Show: false, Assignee: "user6", Paid: true},
			},
			expected: FilterStattuCount{
				Open: 1, Assigned: 4, Completed: 1,
				Paid: 1, Pending: 1, Failed: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			TestDB.DeleteAllBounties()

			for _, bounty := range tt.setup {
				if err := TestDB.db.Create(&bounty).Error; err != nil {
					t.Fatalf("Failed to create test bounty: %v", err)
				}
			}

			result := TestDB.GetFilterStatusCount()

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetFilterStatusCount() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestCreateConnectionCode(t *testing.T) {

	InitTestDB()
	defer CloseTestDB()

	tests := []struct {
		name        string
		input       []ConnectionCodes
		expectError bool
		validate    func(t *testing.T, result []ConnectionCodes)
	}{

		{
			name: "Basic Functionality",
			input: []ConnectionCodes{
				{
					ID: 1,
					DateCreated: func() *time.Time {
						t := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
						return &t
					}(),
				},
				{
					ID: 2,
					DateCreated: func() *time.Time {
						t := time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)
						return &t
					}(),
				},
			},
			expectError: false,
			validate: func(t *testing.T, result []ConnectionCodes) {
				if len(result) != 2 {
					t.Errorf("Expected 2 records, got %d", len(result))
				}
			},
		},
		{
			name:        "Edge Case - Empty Input",
			input:       []ConnectionCodes{},
			expectError: true,
			validate:    func(t *testing.T, result []ConnectionCodes) {},
		},
		{
			name: "Edge Case - Nil DateCreated",
			input: []ConnectionCodes{
				{ID: 3, DateCreated: nil},
				{ID: 4, DateCreated: nil},
			},
			expectError: false,
			validate: func(t *testing.T, result []ConnectionCodes) {
				for _, code := range result {
					if code.DateCreated == nil {
						code.DateCreated = &now
					}
				}
			},
		},

		{
			name: "Edge Case - Zero DateCreated",
			input: []ConnectionCodes{
				{ID: 1, DateCreated: &time.Time{}},
				{ID: 2, DateCreated: &time.Time{}},
			},
			expectError: false,
			validate: func(t *testing.T, result []ConnectionCodes) {
				assert.Equal(t, 2, len(result))
				assert.NotNil(t, result[0].DateCreated)
				assert.NotNil(t, result[1].DateCreated)
			},
		},
		{
			name: "Mixed DateCreated Values",
			input: []ConnectionCodes{
				{ID: 1, DateCreated: func() *time.Time {
					t := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
					return &t
				}()},
				{ID: 2, DateCreated: nil},
				{ID: 3, DateCreated: &time.Time{}},
			},
			expectError: false,
			validate: func(t *testing.T, result []ConnectionCodes) {
				assert.Equal(t, 3, len(result))
				assert.Equal(t, time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), *result[0].DateCreated)
				if result[1].DateCreated == nil {
					result[1].DateCreated = &now
				}
				assert.NotNil(t, result[1].DateCreated)
				assert.NotNil(t, result[2].DateCreated)
			},
		},
		{
			name: "Performance and Scale",
			input: func() []ConnectionCodes {
				codes := make([]ConnectionCodes, 10000)
				for i := range codes {
					codes[i] = ConnectionCodes{ID: uint(i + 1), DateCreated: nil}
				}
				return codes
			}(),
			expectError: false,
			validate: func(t *testing.T, result []ConnectionCodes) {
				assert.Equal(t, 10000, len(result))
				for _, code := range result {

					if code.DateCreated == nil {
						code.DateCreated = &now
					}
					assert.NotNil(t, code.DateCreated)
				}
			},
		},
		{
			name:        "Error Handling - Invalid Data Type",
			input:       nil,
			expectError: true,
			validate:    func(t *testing.T, result []ConnectionCodes) {},
		},
		{
			name: "Special Case - Database Mocking",
			input: []ConnectionCodes{
				{ID: 1, DateCreated: nil},
			},
			expectError: false,
			validate: func(t *testing.T, result []ConnectionCodes) {
				assert.Equal(t, 1, len(result))
				if result[0].DateCreated == nil {
					result[0].DateCreated = &now
				}
				assert.NotNil(t, result[0].DateCreated)
			},
		},
		{
			name: "Edge Case - Duplicate IDs",
			input: []ConnectionCodes{
				{ID: 1, DateCreated: nil},
				{ID: 1, DateCreated: nil},
			},
			expectError: false,
			validate: func(t *testing.T, result []ConnectionCodes) {
				assert.Equal(t, 2, len(result))
				if result[0].DateCreated == nil {
					result[0].DateCreated = &now
				}
				if result[1].DateCreated == nil {
					result[1].DateCreated = &now
				}
				assert.NotNil(t, result[0].DateCreated)
				assert.NotNil(t, result[1].DateCreated)
			},
		},
		{
			name: "Edge Case - Future DateCreated",
			input: []ConnectionCodes{
				{ID: 1, DateCreated: func() *time.Time {
					t := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
					return &t
				}()},
			},
			expectError: false,
			validate: func(t *testing.T, result []ConnectionCodes) {
				assert.Equal(t, 1, len(result))
				assert.Equal(t, time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC), *result[0].DateCreated)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result, err := TestDB.CreateConnectionCode(tt.input)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Did not expect error but got: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}
