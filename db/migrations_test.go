package db

import (
	"testing"
)

func TestTicketsMigrationPostgres(t *testing.T) {
	// init test db
	InitTestDB()

	// Check if table exists
	if !TestDB.db.Migrator().HasTable(&Tickets{}) {
		t.Fatalf("Table 'tickets' does not exist after migration")
	}

	// Check if columns exist
	columns := []string{"uuid", "feature_uuid", "phase_uuid", "name", "sequence", "dependency", "description", "status", "created_at", "updated_at"}
	for _, column := range columns {
		if !TestDB.db.Migrator().HasColumn(&Tickets{}, column) {
			t.Errorf("Column %s is missing in the 'tickets' table", column)
		}
	}

	t.Log("Migration test for Tickets struct with PostgreSQL passed")
}
