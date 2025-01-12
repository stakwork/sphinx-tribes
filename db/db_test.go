package db

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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

	cleanup := func() {
		TestDB.db.Exec("DELETE FROM connectioncodes")
	}

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
				{ID: 5, DateCreated: &time.Time{}},
				{ID: 6, DateCreated: &time.Time{}},
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
				{ID: 7, DateCreated: func() *time.Time {
					t := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
					return &t
				}()},
				{ID: 8, DateCreated: nil},
				{ID: 9, DateCreated: &time.Time{}},
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

			cleanup()
		})
	}
}

func TestIncrementProofCount(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	tests := []struct {
		name          string
		bountyID      uint
		initialBounty *NewBounty
		expectedCount int
		expectedError error
	}{
		{
			name:     "Valid Bounty ID with Existing Record",
			bountyID: 1,
			initialBounty: &NewBounty{
				ID:               1,
				ProofOfWorkCount: 5,
			},
			expectedCount: 6,
			expectedError: nil,
		},
		{
			name:          "Minimum Bounty ID",
			bountyID:      0,
			initialBounty: nil,
			expectedCount: 0,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "Maximum Bounty ID",
			bountyID:      uint(2147483647),
			initialBounty: nil,
			expectedCount: 0,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:          "Non-Existent Bounty ID",
			bountyID:      9999,
			initialBounty: nil,
			expectedCount: 0,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Bounty with Maximum Proof of Work Count",
			bountyID: 2,
			initialBounty: &NewBounty{
				ID:               2,
				ProofOfWorkCount: 21,
			},
			expectedCount: 22,
			expectedError: nil,
		},
		{
			name:     "Bounty with Null Updated Timestamp",
			bountyID: 3,
			initialBounty: &NewBounty{
				ID:               3,
				ProofOfWorkCount: 10,
				Updated:          nil,
			},
			expectedCount: 11,
			expectedError: nil,
		},
		{
			name:     "Bounty with Negative Proof of Work Count",
			bountyID: 4,
			initialBounty: &NewBounty{
				ID:               4,
				ProofOfWorkCount: -5,
			},
			expectedCount: -4,
			expectedError: nil,
		},
		{
			name:     "Bounty with Maximum Updated Timestamp",
			bountyID: 5,
			initialBounty: &NewBounty{
				ID:               5,
				ProofOfWorkCount: 15,
				Updated:          &time.Time{},
			},
			expectedCount: 16,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			TestDB.DeleteAllBounties()

			if tt.initialBounty != nil {
				if err := TestDB.db.Create(tt.initialBounty).Error; err != nil {
					t.Fatalf("Failed to create test bounty: %v", err)
				}
			}

			err := TestDB.IncrementProofCount(tt.bountyID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}

			assert.NoError(t, err)

			var bounty NewBounty
			if err := TestDB.db.Where("id = ?", tt.bountyID).First(&bounty).Error; err != nil {
				t.Fatalf("Failed to retrieve bounty: %v", err)
			}

			assert.Equal(t, tt.expectedCount, bounty.ProofOfWorkCount)

			if bounty.Updated != nil {
				assert.WithinDuration(t, time.Now(), *bounty.Updated, time.Second)
			} else {
				t.Error("Updated timestamp is nil")
			}
		})
	}
}

func parseDate(dateStr string) *time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic("Error parsing date: " + err.Error())
	}
	return &date
}

func TestGetConnectionCode(t *testing.T) {
	InitTestDB()

	defer CloseTestDB()

	cleanup := func() {
		TestDB.db.Exec("DELETE FROM connectioncodes")
	}

	tests := []struct {
		name             string
		connectionCodes  []ConnectionCodes
		expected         ConnectionCodesShort
		expectUpdateCall bool
	}{
		{
			name: "Basic Functionality",
			connectionCodes: []ConnectionCodes{
				{ConnectionString: "code1", DateCreated: parseDate("2006-01-02"), IsUsed: false},
				{ConnectionString: "code2", DateCreated: parseDate("2023-10-02"), IsUsed: false},
			},
			expected:         ConnectionCodesShort{ConnectionString: "code2", DateCreated: parseDate("2023-10-02")},
			expectUpdateCall: true,
		},
		{
			name: "No Unused Connection Codes",
			connectionCodes: []ConnectionCodes{
				{ConnectionString: "code1", DateCreated: parseDate("2023-10-01"), IsUsed: true},
			},
			expected:         ConnectionCodesShort{},
			expectUpdateCall: false,
		},
		{
			name: "Single Unused Connection Code",
			connectionCodes: []ConnectionCodes{
				{ConnectionString: "code1", DateCreated: parseDate("2023-10-01"), IsUsed: false},
			},
			expected:         ConnectionCodesShort{ConnectionString: "code1", DateCreated: parseDate("2023-10-01")},
			expectUpdateCall: true,
		},
		{
			name: "Multiple Unused Connection Codes",
			connectionCodes: []ConnectionCodes{
				{ConnectionString: "code1", DateCreated: parseDate("2023-10-01"), IsUsed: false},
				{ConnectionString: "code2", DateCreated: parseDate("2023-10-02"), IsUsed: false},
			},
			expected:         ConnectionCodesShort{ConnectionString: "code2", DateCreated: parseDate("2023-10-02")},
			expectUpdateCall: true,
		},
		{
			name:             "Edge Case: Empty Database",
			connectionCodes:  []ConnectionCodes{},
			expected:         ConnectionCodesShort{},
			expectUpdateCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			TestDB.CreateConnectionCode(tt.connectionCodes)

			result := TestDB.GetConnectionCode()

			assert.Equal(t, tt.expected.ConnectionString, result.ConnectionString)

			cleanup()
		})
	}
}

func TestGetBountiesLeaderboard(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	tests := []struct {
		name     string
		setup    []NewBounty
		expected []LeaderData
	}{
		{
			name: "Standard Input with Multiple Users",
			setup: []NewBounty{
				{
					Assignee: "user1", Price: uint(100), Paid: true,
					Type: "coding", Title: "Test Bounty 1",
				},
				{
					Assignee: "user1", Price: uint(200), Paid: true,
					Type: "coding", Title: "Test Bounty 2",
				},
				{
					Assignee: "user2", Price: uint(150), Paid: true,
					Type: "coding", Title: "Test Bounty 3",
				},
			},
			expected: []LeaderData{
				{"owner_pubkey": "user1", "total_bounties_completed": uint(2), "total_sats_earned": uint(300)},
				{"owner_pubkey": "user2", "total_bounties_completed": uint(1), "total_sats_earned": uint(150)},
			},
		},
		{
			name: "Single User with Completed Bounties",
			setup: []NewBounty{
				{
					Assignee: "user1", Price: uint(100), Paid: true,
					Type: "coding", Title: "Test Bounty 1",
				},
				{
					Assignee: "user1", Price: uint(200), Paid: true,
					Type: "coding", Title: "Test Bounty 2",
				},
			},
			expected: []LeaderData{
				{"owner_pubkey": "user1", "total_bounties_completed": uint(2), "total_sats_earned": uint(300)},
			},
		},
		{
			name: "No Completed Bounties",
			setup: []NewBounty{
				{
					Assignee: "user1", Price: uint(100), Paid: false,
					Type: "coding", Title: "Test Bounty",
				},
			},
			expected: []LeaderData{},
		},
		{
			name: "Users with Zero Sats Earned",
			setup: []NewBounty{
				{
					Assignee: "user1", Price: uint(0), Paid: true,
					Type: "coding", Title: "Test Bounty 1",
				},
				{
					Assignee: "user2", Price: uint(0), Paid: true,
					Type: "coding", Title: "Test Bounty 2",
				},
			},
			expected: []LeaderData{
				{"owner_pubkey": "user1", "total_bounties_completed": uint(1), "total_sats_earned": uint(0)},
				{"owner_pubkey": "user2", "total_bounties_completed": uint(1), "total_sats_earned": uint(0)},
			},
		},
		{
			name: "Maximum Integer Values for Sats",
			setup: []NewBounty{
				{
					Assignee: "user1", Price: uint(2147483647), Paid: true,
					Type: "coding", Title: "Test Bounty",
				},
			},
			expected: []LeaderData{
				{"owner_pubkey": "user1", "total_bounties_completed": uint(1), "total_sats_earned": uint(2147483647)},
			},
		},
		{
			name: "Invalid Data Types in Database",
			setup: []NewBounty{
				{
					Assignee: "user1", Price: uint(0), Paid: true,
					Type: "coding", Title: "Test Bounty 1",
				},
				{
					Assignee: "user1", Price: uint(100), Paid: true,
					Type: "coding", Title: "Test Bounty 2",
				},
			},
			expected: []LeaderData{
				{"owner_pubkey": "user1", "total_bounties_completed": uint(2), "total_sats_earned": uint(100)},
			},
		},
		{
			name:  "Large Number of Users",
			setup: generateLargeUserSet(1000),
			expected: []LeaderData{
				{"owner_pubkey": "user999", "total_bounties_completed": uint(1), "total_sats_earned": uint(1999)},
				{"owner_pubkey": "user998", "total_bounties_completed": uint(1), "total_sats_earned": uint(1998)},
				{"owner_pubkey": "user997", "total_bounties_completed": uint(1), "total_sats_earned": uint(1997)},
				{"owner_pubkey": "user996", "total_bounties_completed": uint(1), "total_sats_earned": uint(1996)},
				{"owner_pubkey": "user995", "total_bounties_completed": uint(1), "total_sats_earned": uint(1995)},
			},
		},
		{
			name: "Duplicate Users with Different Bounties",
			setup: []NewBounty{
				{
					Assignee: "user1", Price: uint(100), Paid: true,
					Type: "coding", Title: "Test Bounty 1",
				},
				{
					Assignee: "user1", Price: uint(100), Paid: true,
					Type: "coding", Title: "Test Bounty 2",
				},
				{
					Assignee: "user1", Price: uint(100), Paid: false,
					Type: "coding", Title: "Test Bounty 3",
				},
			},
			expected: []LeaderData{
				{"owner_pubkey": "user1", "total_bounties_completed": uint(2), "total_sats_earned": uint(200)},
			},
		},
		{
			name: "Users with Identical Sats Earned",
			setup: []NewBounty{
				{
					Assignee: "user1", Price: uint(100), Paid: true,
					Type: "coding", Title: "Test Bounty 1",
				},
				{
					Assignee: "user2", Price: uint(100), Paid: true,
					Type: "coding", Title: "Test Bounty 2",
				},
			},
			expected: []LeaderData{
				{"owner_pubkey": "user1", "total_bounties_completed": uint(1), "total_sats_earned": uint(100)},
				{"owner_pubkey": "user2", "total_bounties_completed": uint(1), "total_sats_earned": uint(100)},
			},
		},
		{
			name:     "Empty Database",
			setup:    []NewBounty{},
			expected: []LeaderData{},
		},
		{
			name: "Zero Value for Negative Input",
			setup: []NewBounty{
				{
					Assignee: "user1", Price: uint(0), Paid: true,
					Type: "coding", Title: "Test Bounty",
				},
			},
			expected: []LeaderData{
				{"owner_pubkey": "user1", "total_bounties_completed": uint(1), "total_sats_earned": uint(0)},
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

			result := TestDB.GetBountiesLeaderboard()

			if tt.name == "Large Number of Users" {
				if len(result) < len(tt.expected) {
					t.Errorf("Expected at least %d results, got %d", len(tt.expected), len(result))
					return
				}
			} else {
				if len(result) != len(tt.expected) {
					t.Errorf("Expected %d results, got %d", len(tt.expected), len(result))
					return
				}
			}

			if tt.name == "Large Number of Users" {
				for i, expected := range tt.expected {
					actual := result[i]
					if actual["owner_pubkey"] != expected["owner_pubkey"] {
						t.Errorf("Expected owner_pubkey %v, got %v", expected["owner_pubkey"], actual["owner_pubkey"])
					}
					expectedSats := uint(1000 + 999 - i)
					if actual["total_sats_earned"] != expectedSats {
						t.Errorf("Expected total_sats_earned %v, got %v", expectedSats, actual["total_sats_earned"])
					}
					if actual["total_bounties_completed"] != uint(1) {
						t.Errorf("Expected total_bounties_completed 1, got %v", actual["total_bounties_completed"])
					}
				}
			} else if tt.name == "Users with Zero Sats Earned" || tt.name == "Users with Identical Sats Earned" {
				for _, expected := range tt.expected {
					found := false
					for _, actual := range result {
						if actual["owner_pubkey"] == expected["owner_pubkey"] &&
							actual["total_bounties_completed"] == expected["total_bounties_completed"] &&
							actual["total_sats_earned"] == expected["total_sats_earned"] {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected to find user %v with bounties %v and sats %v",
							expected["owner_pubkey"],
							expected["total_bounties_completed"],
							expected["total_sats_earned"])
					}
				}
			} else {
				for i, expected := range tt.expected {
					if i >= len(result) {
						t.Errorf("Missing expected result at index %d", i)
						continue
					}

					actual := result[i]
					if actual["owner_pubkey"] != expected["owner_pubkey"] {
						t.Errorf("Expected owner_pubkey %v, got %v", expected["owner_pubkey"], actual["owner_pubkey"])
					}
					if actual["total_bounties_completed"] != expected["total_bounties_completed"] {
						t.Errorf("Expected total_bounties_completed %v, got %v",
							expected["total_bounties_completed"], actual["total_bounties_completed"])
					}
					if actual["total_sats_earned"] != expected["total_sats_earned"] {
						t.Errorf("Expected total_sats_earned %v, got %v",
							expected["total_sats_earned"], actual["total_sats_earned"])
					}
				}
			}
		})
	}
}

func generateLargeUserSet(count int) []NewBounty {
	bounties := make([]NewBounty, count)
	for i := 0; i < count; i++ {
		bounties[i] = NewBounty{
			Assignee: fmt.Sprintf("user%d", i),
			Price:    uint(1000 + i),
			Paid:     true,
			Type:     "coding",
			Title:    fmt.Sprintf("Test Bounty %d", i),
		}
	}
	return bounties
}
