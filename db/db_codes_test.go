package db

import (
	"gorm.io/driver/postgres"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	gormIo "gorm.io/gorm"
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

func TestGetAllBounties(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gormIo.Open(dialector, &gormIo.Config{})
	DB = database{db: db}

	t.Run("should return all bounties, not query parameters", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM public.bounty WHERE show != false     ORDER BY created desc LIMIT -1  OFFSET 0`)).WithArgs(false).WillReturnRows(sqlmock.NewRows([]string{}))

		req, _ := http.NewRequest("GET", "/gobounties/all", nil)
		DB.GetAllBounties(req)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
