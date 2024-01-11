package db

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var now = time.Now()

var code = ConnectionCodes{
	ID:               1,
	ConnectionString: "2222222",
	IsUsed:           false,
	DateCreated:      &now,
}

func TestCodeGet(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gorm.Open("postgres", db)
	rows := sqlmock.NewRows([]string{"connection_string", "date_created", "is_used", "date_created"}).AddRow(code.ID, code.ConnectionString, code.IsUsed, code.DateCreated)
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT connection_string, date_created FROM connectioncodes WHERE is_used = ? ORDER BY id DESC LIMIT 1`)).
		WithArgs(false).
		WillReturnRows(rows)

	assert.Nil(t, err)
}
