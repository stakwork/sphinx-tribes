package db

import (
	"reflect"
	"testing"
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
