package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
)

func TestGetWantedsHeader(t *testing.T) {
	db.InitTestDB()

	db.CleanTestData()

	defer func() {
		db.CleanTestData()
		db.CloseTestDB()
	}()

	db.DB = db.TestDB

	tests := []struct {
		name           string
		setupTestData  func(t *testing.T)
		expectedStatus int
		validate       func(t *testing.T, response []byte)
	}{
		{
			name: "Standard Case",
			setupTestData: func(t *testing.T) {

				for i := 1; i <= 5; i++ {
					person := db.Person{
						ID:          uint(i),
						Uuid:        fmt.Sprintf("uuid-%d", i),
						OwnerPubKey: fmt.Sprintf("test-pub-key-%d", i),
						OwnerAlias:  fmt.Sprintf("test-alias-%d", i),
						UniqueName:  fmt.Sprintf("test-name-%d", i),
						Img:         fmt.Sprintf("test-img-%d", i),
					}
					_, err := db.TestDB.CreateOrEditPerson(person)
					assert.NoError(t, err)
				}

				for i := 1; i <= 3; i++ {
					bounty := db.Bounty{
						Title:   fmt.Sprintf("Test Bounty %d", i),
						OwnerID: fmt.Sprintf("test-pub-key-%d", i),
					}
					_, err := db.TestDB.AddBounty(bounty)
					assert.NoError(t, err)
				}
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, response []byte) {
				var result struct {
					DeveloperCount int64               `json:"developer_count"`
					BountiesCount  uint64              `json:"bounties_count"`
					People         *[]db.PersonInShort `json:"people"`
				}
				err := json.Unmarshal(response, &result)
				assert.NoError(t, err)
				assert.Equal(t, int64(5), result.DeveloperCount)
				assert.Equal(t, uint64(3), result.BountiesCount)
				assert.NotNil(t, result.People)
				assert.Equal(t, 3, len(*result.People))
			},
		},
		{
			name: "No Developers",
			setupTestData: func(t *testing.T) {
				db.DeleteAllBounties()
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, response []byte) {
				var result struct {
					DeveloperCount int64               `json:"developer_count"`
					BountiesCount  uint64              `json:"bounties_count"`
					People         *[]db.PersonInShort `json:"people"`
				}
				err := json.Unmarshal(response, &result)
				assert.NoError(t, err)
				assert.Equal(t, int64(0), result.DeveloperCount)
				assert.Equal(t, uint64(0), result.BountiesCount)
				assert.NotNil(t, result.People)
				assert.Equal(t, 0, len(*result.People))
			},
		},
		{
			name: "No Bounties",
			setupTestData: func(t *testing.T) {

				for i := 1; i <= 3; i++ {
					person := db.Person{
						ID:          uint(i),
						Uuid:        fmt.Sprintf("uuid-%d", i),
						OwnerPubKey: fmt.Sprintf("test-pub-key-%d", i),
						OwnerAlias:  fmt.Sprintf("test-alias-%d", i),
						UniqueName:  fmt.Sprintf("test-name-%d", i),
						Img:         fmt.Sprintf("test-img-%d", i),
					}
					_, err := db.TestDB.CreateOrEditPerson(person)
					assert.NoError(t, err)
				}
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, response []byte) {
				var result struct {
					DeveloperCount int64               `json:"developer_count"`
					BountiesCount  uint64              `json:"bounties_count"`
					People         *[]db.PersonInShort `json:"people"`
				}
				err := json.Unmarshal(response, &result)
				assert.NoError(t, err)
				assert.Equal(t, int64(3), result.DeveloperCount)
				assert.Equal(t, uint64(0), result.BountiesCount)
				assert.NotNil(t, result.People)
				assert.Equal(t, 3, len(*result.People))
			},
		},
		{
			name: "Maximum People",
			setupTestData: func(t *testing.T) {

				for i := 1; i <= 10; i++ {
					person := db.Person{
						ID:          uint(i),
						Uuid:        fmt.Sprintf("uuid-%d", i),
						OwnerPubKey: fmt.Sprintf("test-pub-key-%d", i),
						OwnerAlias:  fmt.Sprintf("test-alias-%d", i),
						UniqueName:  fmt.Sprintf("test-name-%d", i),
						Img:         fmt.Sprintf("test-img-%d", i),
					}
					_, err := db.TestDB.CreateOrEditPerson(person)
					assert.NoError(t, err)
				}
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, response []byte) {
				var result struct {
					DeveloperCount int64               `json:"developer_count"`
					BountiesCount  uint64              `json:"bounties_count"`
					People         *[]db.PersonInShort `json:"people"`
				}
				err := json.Unmarshal(response, &result)
				assert.NoError(t, err)
				assert.Equal(t, int64(10), result.DeveloperCount)
				assert.Equal(t, uint64(0), result.BountiesCount)
				assert.NotNil(t, result.People)
				assert.Equal(t, 3, len(*result.People))
			},
		},
		{
			name: "Large Number Developers and Bounties",
			setupTestData: func(t *testing.T) {

				for i := 1; i <= 500; i++ {
					person := db.Person{
						ID:          uint(i),
						Uuid:        fmt.Sprintf("uuid-%d", i),
						OwnerPubKey: fmt.Sprintf("test-pub-key-%d", i),
						OwnerAlias:  fmt.Sprintf("test-alias-%d", i),
						UniqueName:  fmt.Sprintf("test-name-%d", i),
						Img:         fmt.Sprintf("test-img-%d", i),
					}
					_, err := db.TestDB.CreateOrEditPerson(person)
					assert.NoError(t, err)
				}

				for i := 1; i <= 250; i++ {
					bounty := db.Bounty{
						Title:   fmt.Sprintf("Test Bounty %d", i),
						OwnerID: fmt.Sprintf("test-pub-key-%d", i),
					}
					_, err := db.TestDB.AddBounty(bounty)
					assert.NoError(t, err)
				}
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, response []byte) {

				expectedDeveloperCount := db.TestDB.CountDevelopers()
				expectedBountiesCount := db.TestDB.CountBounties()

				var result struct {
					DeveloperCount int64               `json:"developer_count"`
					BountiesCount  uint64              `json:"bounties_count"`
					People         *[]db.PersonInShort `json:"people"`
				}
				err := json.Unmarshal(response, &result)
				assert.NoError(t, err)
				assert.Equal(t, expectedDeveloperCount, result.DeveloperCount)
				assert.Equal(t, expectedBountiesCount, result.BountiesCount)
				assert.NotNil(t, result.People)
				assert.Equal(t, 3, len(*result.People))
			},
		},
		{
			name: "No People",
			setupTestData: func(t *testing.T) {

				for i := 1; i <= 3; i++ {
					bounty := db.Bounty{
						Title:   fmt.Sprintf("Test Bounty %d", i),
						OwnerID: fmt.Sprintf("test-pub-key-%d", i),
					}
					_, err := db.TestDB.AddBounty(bounty)
					assert.NoError(t, err)
				}
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, response []byte) {
				var result struct {
					DeveloperCount int64               `json:"developer_count"`
					BountiesCount  uint64              `json:"bounties_count"`
					People         *[]db.PersonInShort `json:"people"`
				}
				err := json.Unmarshal(response, &result)
				assert.NoError(t, err)
				assert.Equal(t, int64(0), result.DeveloperCount)
				assert.Equal(t, uint64(3), result.BountiesCount)
				assert.NotNil(t, result.People)
				assert.Equal(t, 0, len(*result.People))
			},
		},
		{
			name: "Negative Developer Count",
			setupTestData: func(t *testing.T) {

				person := db.Person{
					ID:          1,
					Uuid:        "uuid-1",
					OwnerPubKey: "test-pub-key-1",
					OwnerAlias:  "test-alias-1",
					UniqueName:  "test-name-1",
					Img:         "test-img-1",
					Deleted:     true,
				}
				_, err := db.TestDB.CreateOrEditPerson(person)
				assert.NoError(t, err)

				db.DeleteAllBounties()
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, response []byte) {
				var result struct {
					DeveloperCount int64               `json:"developer_count"`
					BountiesCount  uint64              `json:"bounties_count"`
					People         *[]db.PersonInShort `json:"people"`
				}
				err := json.Unmarshal(response, &result)
				assert.NoError(t, err)
				assert.Equal(t, int64(0), result.DeveloperCount)
				assert.Equal(t, uint64(0), result.BountiesCount)
				assert.NotNil(t, result.People)
				assert.Equal(t, 0, len(*result.People))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setupTestData(t)

			req := httptest.NewRequest(http.MethodGet, "/wanteds/header", nil)
			w := httptest.NewRecorder()

			GetWantedsHeader(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			tt.validate(t, w.Body.Bytes())
		})
		db.CleanTestData()
	}
}
