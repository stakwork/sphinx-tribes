package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
)

func TestMigrateBounties(t *testing.T) {

	db.InitTestDB()

	defer func() {
		db.CleanTestData()
		db.CloseTestDB()
	}()

	originalDB := db.DB
	defer func() {
		db.DB = originalDB
	}()

	db.DB = db.TestDB

	workspace := db.Workspace{
		Uuid:        "test-workspace-uuid",
		Name:        "Test Workspace",
		OwnerPubKey: "test-pub-key-1",
	}
	_, err := db.TestDB.CreateOrEditWorkspace(workspace)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		setupTestData func(t *testing.T)
		expectedCalls uint64
		wantErr       bool
	}{
		{
			name: "Single Peep with Valid Bounties",
			setupTestData: func(t *testing.T) {
				validBounty := createValidBounty()
				extras := db.PropertyMap{"wanted": []interface{}{validBounty}}

				person := db.Person{
					ID:          1,
					OwnerPubKey: "test-pub-key-1",
					OwnerAlias:  "test-alias-1",
					UniqueName:  "test-unique-name-1",
					Description: "test-description-1",
					Extras:      extras,
				}
				_, err := db.TestDB.CreateOrEditPerson(person)
				assert.NoError(t, err)

				assignee := db.Person{
					ID:          2,
					OwnerPubKey: "test-pub-key-2",
					OwnerAlias:  "test-alias-2",
					UniqueName:  "test-unique-name-2",
					Description: "test-description-2",
				}
				_, err = db.TestDB.CreateOrEditPerson(assignee)
				assert.NoError(t, err)
			},
			expectedCalls: 1,
			wantErr:       false,
		},
		{
			name: "Peep with No Bounties",
			setupTestData: func(t *testing.T) {
				person := db.Person{
					ID:          3,
					OwnerPubKey: "test-pub-key-3",
					OwnerAlias:  "test-alias-3",
					UniqueName:  "test-unique-name-3",
					Description: "test-description-3",
				}
				_, err := db.TestDB.CreateOrEditPerson(person)
				assert.NoError(t, err)
			},
			expectedCalls: 0,
			wantErr:       false,
		},
		{
			name: "Empty Bounties List",
			setupTestData: func(t *testing.T) {
				extras := db.PropertyMap{"wanted": []interface{}{}}
				person := db.Person{
					ID:          4,
					OwnerPubKey: "test-pub-key-4",
					OwnerAlias:  "test-alias-4",
					UniqueName:  "test-unique-name-4",
					Description: "test-description-4",
					Extras:      extras,
				}
				_, err := db.TestDB.CreateOrEditPerson(person)
				assert.NoError(t, err)
			},
			expectedCalls: 0,
			wantErr:       false,
		},
		{
			name: "Invalid Data Types",
			setupTestData: func(t *testing.T) {
				invalidBounty := map[string]interface{}{
					"title":    123,
					"paid":     "true",
					"price":    "1000",
					"created":  "invalid-timestamp",
					"assignee": "invalid-assignee-format",
				}
				extras := db.PropertyMap{"wanted": []interface{}{invalidBounty}}
				person := db.Person{
					ID:          5,
					OwnerPubKey: "test-pub-key-5",
					OwnerAlias:  "test-alias-5",
					UniqueName:  "test-unique-name-5",
					Description: "test-description-5",
					Extras:      extras,
				}
				_, err := db.TestDB.CreateOrEditPerson(person)
				assert.NoError(t, err)
			},
			expectedCalls: 1,
			wantErr:       false,
		},
		{
			name: "Boundary Values for Numeric Fields",
			setupTestData: func(t *testing.T) {
				bounty := createValidBounty()
				bounty["price"] = uint(0)
				bounty["created"] = float64(0)
				extras := db.PropertyMap{"wanted": []interface{}{bounty}}
				person := db.Person{
					ID:          6,
					OwnerPubKey: "test-pub-key-6",
					OwnerAlias:  "test-alias-6",
					UniqueName:  "test-unique-name-6",
					Description: "test-description-6",
					Extras:      extras,
				}
				_, err := db.TestDB.CreateOrEditPerson(person)
				assert.NoError(t, err)
			},
			expectedCalls: 1,
			wantErr:       false,
		},
		{
			name: "Missing Fields",
			setupTestData: func(t *testing.T) {
				minimalBounty := map[string]interface{}{
					"title":          "Minimal Bounty",
					"assignee":       map[string]interface{}{"owner_pubkey": "test-pub-key-2"},
					"workspace_uuid": "test-workspace-uuid",
				}
				extras := db.PropertyMap{"wanted": []interface{}{minimalBounty}}
				person := db.Person{
					ID:          7,
					OwnerPubKey: "test-pub-key-7",
					OwnerAlias:  "test-alias-7",
					UniqueName:  "test-unique-name-7",
					Description: "test-description-7",
					Extras:      extras,
				}
				_, err := db.TestDB.CreateOrEditPerson(person)
				assert.NoError(t, err)
			},
			expectedCalls: 1,
			wantErr:       false,
		},
		{
			name: "Null Values",
			setupTestData: func(t *testing.T) {
				bounty := createValidBounty()
				bounty["description"] = nil
				bounty["estimated_completion_date"] = nil
				bounty["coding_language"] = nil
				extras := db.PropertyMap{"wanted": []interface{}{bounty}}
				person := db.Person{
					ID:          8,
					OwnerPubKey: "test-pub-key-8",
					OwnerAlias:  "test-alias-8",
					UniqueName:  "test-unique-name-8",
					Description: "test-description-8",
					Extras:      extras,
				}
				_, err := db.TestDB.CreateOrEditPerson(person)
				assert.NoError(t, err)
			},
			expectedCalls: 1,
			wantErr:       false,
		},
		{
			name: "Assignee Not Found",
			setupTestData: func(t *testing.T) {
				bounty := createValidBounty()
				bounty["assignee"] = map[string]interface{}{"owner_pubkey": "non-existent-pubkey"}
				extras := db.PropertyMap{"wanted": []interface{}{bounty}}
				person := db.Person{
					ID:          9,
					OwnerPubKey: "test-pub-key-9",
					OwnerAlias:  "test-alias-9",
					UniqueName:  "test-unique-name-9",
					Description: "test-description-9",
					Extras:      extras,
				}
				_, err := db.TestDB.CreateOrEditPerson(person)
				assert.NoError(t, err)
			},
			expectedCalls: 1,
			wantErr:       false,
		},
		{
			name: "Complex Coding Languages Structure",
			setupTestData: func(t *testing.T) {
				bounty := createValidBounty()
				bounty["coding_language"] = db.PropertyMap{
					"value": pq.StringArray{"Go", "Python", "JavaScript", "Ruby", "C++"},
					"extra": "should be ignored",
				}
				extras := db.PropertyMap{"wanted": []interface{}{bounty}}
				person := db.Person{
					ID:          10,
					OwnerPubKey: "test-pub-key-10",
					OwnerAlias:  "test-alias-10",
					UniqueName:  "test-unique-name-10",
					Description: "test-description-10",
					Extras:      extras,
				}
				_, err := db.TestDB.CreateOrEditPerson(person)
				assert.NoError(t, err)
			},
			expectedCalls: 1,
			wantErr:       false,
		},
		{
			name: "Mixed Valid and Invalid Bounties",
			setupTestData: func(t *testing.T) {
				validBounty := createValidBounty()
				invalidBounty := map[string]interface{}{
					"title": 123,
					"paid":  "invalid",
				}
				extras := db.PropertyMap{"wanted": []interface{}{validBounty, invalidBounty}}
				person := db.Person{
					ID:          11,
					OwnerPubKey: "test-pub-key-11",
					OwnerAlias:  "test-alias-11",
					UniqueName:  "test-unique-name-11",
					Description: "test-description-11",
					Extras:      extras,
				}
				_, err := db.TestDB.CreateOrEditPerson(person)
				assert.NoError(t, err)
			},
			expectedCalls: 2,
			wantErr:       false,
		},
		{
			name: "Duplicate Bounties",
			setupTestData: func(t *testing.T) {
				bounty := createValidBounty()
				extras := db.PropertyMap{"wanted": []interface{}{bounty, bounty}}
				person := db.Person{
					ID:          12,
					OwnerPubKey: "test-pub-key-12",
					OwnerAlias:  "test-alias-12",
					UniqueName:  "test-unique-name-12",
					Description: "test-description-12",
					Extras:      extras,
				}
				_, err := db.TestDB.CreateOrEditPerson(person)
				assert.NoError(t, err)
			},
			expectedCalls: 2,
			wantErr:       false,
		},
		{
			name: "Non-String Assignee Owner PubKey",
			setupTestData: func(t *testing.T) {
				bounty := createValidBounty()
				bounty["assignee"] = map[string]interface{}{
					"owner_pubkey": "12345",
				}
				extras := db.PropertyMap{"wanted": []interface{}{bounty}}

				person := db.Person{
					ID:          13,
					OwnerPubKey: "test-pub-key-13",
					OwnerAlias:  "test-alias-13",
					UniqueName:  "test-unique-name-13",
					Description: "test-description-13",
					Extras:      extras,
				}
				_, err := db.TestDB.CreateOrEditPerson(person)
				assert.NoError(t, err)

				assignee := db.Person{
					ID:          14,
					OwnerPubKey: "12345",
					OwnerAlias:  "test-alias-14",
					UniqueName:  "test-unique-name-14",
					Description: "test-description-14",
				}
				_, err = db.TestDB.CreateOrEditPerson(assignee)
				assert.NoError(t, err)
			},
			expectedCalls: 1,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db.CleanTestData()

			_, err := db.TestDB.CreateOrEditWorkspace(workspace)
			assert.NoError(t, err)

			tt.setupTestData(t)

			req := httptest.NewRequest(http.MethodGet, "/migrate_bounties", nil)
			w := httptest.NewRecorder()

			MigrateBounties(w, req)

			if tt.wantErr {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}

			count := db.TestDB.CountBounties()
			assert.Equal(t, tt.expectedCalls, count)
		})
	}
}

func createValidBounty() map[string]interface{} {
	return map[string]interface{}{
		"title":                     "Test Bounty",
		"paid":                      true,
		"show":                      true,
		"type":                      "bug",
		"award":                     "100 sats",
		"price":                     uint(1000),
		"tribe":                     "test-tribe",
		"created":                   float64(1234567890),
		"assignee":                  map[string]interface{}{"owner_pubkey": "test-pub-key-2"},
		"ticketUrl":                 "https://github.com/test",
		"description":               "Test description",
		"wanted_type":               "feature",
		"deliverables":              "Test deliverables",
		"coding_language":           db.PropertyMap{"value": pq.StringArray{"Go", "Python"}},
		"github_description":        true,
		"one_sentence_summary":      "Test summary",
		"estimated_session_length":  "2 hours",
		"estimated_completion_date": time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
		"workspace_uuid":            "test-workspace-uuid",
	}
}

func TestGetWantedsHeader(t *testing.T) {

	db.InitTestDB()
	defer func() {
		db.CleanTestData()
		db.CloseTestDB()
	}()

	originalDB := db.DB
	defer func() {
		db.DB = originalDB
	}()

	db.DB = db.TestDB

	db.CleanTestData()

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
			name:           "No Developers",
			setupTestData:  func(t *testing.T) {},
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
			name: "Large Number of Developers and Bounties",
			setupTestData: func(t *testing.T) {

				for i := 1; i <= 1000; i++ {
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

				for i := 1; i <= 500; i++ {
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
				assert.Equal(t, int64(1000), result.DeveloperCount)
				assert.Equal(t, uint64(500), result.BountiesCount)
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
