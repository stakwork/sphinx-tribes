package db

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var now = time.Now()
var person = Person{
	ID:               1,
	Uuid:             xid.New().String(),
	UniqueName:       "test-name",
	OwnerPubKey:      "000000333366666555555555",
	Description:      "Trying this",
	OwnerAlias:       "test",
	Img:              "",
	Unlisted:         false,
	Deleted:          false,
	Created:          &now,
	Updated:          &now,
	Tags:             pq.StringArray{},
	Extras:           map[string]interface{}{},
	OwnerRouteHint:   "00000000000000",
	OwnerContactKey:  "00000000000000",
	PriceToMeet:      10,
	GithubIssues:     map[string]interface{}{},
	TwitterConfirmed: false,
	LastLogin:        20000000,
	NewTicketTime:    10000,
}

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	people *Person
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)

	s.DB.LogMode(true)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestPersonGet() {

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM people WHERE id = ? AND (deleted = 'f' OR deleted is null)`)).
		WithArgs(person.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "unique_name"}).
			AddRow(person))

	res := DB.GetPerson(uint(person.ID))

	require.Nil(s.T(), deep.Equal(&Person{UniqueName: person.UniqueName}, res))
}
