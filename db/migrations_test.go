package db

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestTicketsMigrationPostgres(t *testing.T) {
	// init test db
	InitTestDB()

	defer CloseTestDB()

	// Check if table exists
	if !TestDB.db.Migrator().HasTable(&Tickets{}) {
		t.Fatalf("Table 'tickets' does not exist after migration")
	}

	// Check if columns exist
	columns := []string{
		"uuid",
		"ticket_group",
		"feature_uuid",
		"phase_uuid",
		"name",
		"sequence",
		"dependency",
		"description",
		"status",
		"version",
		"author",
		"author_id",
		"created_at",
		"updated_at",
	}

	for _, column := range columns {
		if !TestDB.db.Migrator().HasColumn(&Tickets{}, column) {
			t.Errorf("Column %s is missing in the 'tickets' table", column)
		}
	}

	indexes := []string{
		"group_index",
		"composite_index",
		"phase_index",
	}

	for _, index := range indexes {
		hasIndex := TestDB.db.Migrator().HasIndex(&Tickets{}, index)
		if !hasIndex {
			t.Errorf("Index %s is missing in the 'tickets' table", index)
		}
	}

	t.Log("Migration test for Tickets struct with PostgreSQL passed")
}

type TestTickets struct {
	UUID        uuid.UUID    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TicketGroup *uuid.UUID   `gorm:"type:uuid;index:group_index" json:"ticket_group,omitempty"`
	FeatureUUID string       `gorm:"type:uuid;not null" json:"feature_uuid" validate:"required"`
	PhaseUUID   string       `gorm:"type:uuid;not null" json:"phase_uuid" validate:"required"`
	Name        string       `gorm:"type:varchar(255);not null"`
	Sequence    int          `gorm:"type:integer;not null;index:composite_index"`
	Dependency  []int        `gorm:"type:integer[]"`
	Description string       `gorm:"type:text"`
	Status      TicketStatus `gorm:"type:varchar(50);not null;default:'draft'"`
	Version     *int         `gorm:"type:integer" json:"version,omitempty"`
	Author      *Author      `gorm:"type:varchar(50)" json:"author,omitempty"`
	AuthorID    *string      `gorm:"type:varchar(255)" json:"author_id,omitempty"`
	CreatedAt   time.Time    `gorm:"type:timestamp;not null;default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"type:timestamp;not null;default:current_timestamp" json:"updated_at"`
}

func TestTicketsIndexPerformance(t *testing.T) {
	InitTestDB()

	defer CloseTestDB()

	// Ensure cleanup
	defer func() {
		t.Log("Cleaning up...")
		TestDB.db.Migrator().DropTable(&TestTickets{})
	}()

	// Prepare the environment
	setupGormTestTickets(TestDB.db, t)

	// Measure performance without explicit index
	t.Log("Measuring performance without explicit index...")
	noIndexTime := measureTicketsQueryPerformance(TestDB.db, t)

	// Add explicit index
	t.Log("Adding explicit index on feature_uuid...")
	if err := TestDB.db.Exec(`CREATE INDEX idx_feature_uuid ON test_tickets (feature_uuid);`).Error; err != nil {
		t.Fatalf("Failed to create index: %v", err)
	}

	// Measure performance with explicit index
	t.Log("Measuring performance with explicit index...")
	withIndexTime := measureTicketsQueryPerformance(TestDB.db, t)

	// Compare results
	t.Logf("Query time without index: %v ms", noIndexTime)
	t.Logf("Query time with index: %v ms", withIndexTime)

	if withIndexTime >= noIndexTime {
		t.Errorf("Index did not improve performance. Time without index: %v ms, Time with index: %v ms", noIndexTime, withIndexTime)
	} else {
		t.Log("Index improved query performance!")
	}
}

// Setup test data for Tickets table
func setupGormTestTickets(db *gorm.DB, t *testing.T) {
	TestDB.db.AutoMigrate(&TestTickets{})

	// Insert test data
	t.Log("Inserting test data...")
	tickets := make([]TestTickets, 100000)
	for i := range tickets {
		tickets[i] = TestTickets{
			FeatureUUID: fmt.Sprintf("00000000-0000-0000-0000-%012d", i%1000),
			PhaseUUID:   fmt.Sprintf("00000000-0000-0000-0000-%012d", i%100),
			Name:        fmt.Sprintf("Ticket %d", i),
			Sequence:    i % 1000,
			Status:      "draft",
		}
	}
	if err := db.CreateInBatches(tickets, 1000).Error; err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
}

// Measure query performance on Tickets table
func measureTicketsQueryPerformance(db *gorm.DB, t *testing.T) int64 {
	query := "00000000-0000-0000-0000-000000000001"
	start := time.Now()

	var count int64
	if err := db.Model(&TestTickets{}).Where("feature_uuid = ?", query).Count(&count).Error; err != nil {
		t.Fatalf("Query execution failed: %v", err)
	}

	duration := time.Since(start).Milliseconds()

	t.Logf("Query execution time: %v ms, Result count: %v", duration, count)
	return duration
}
