package handlers

import (
	"encoding/json"
	"github.com/lib/pq"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
	"net/http"
	"reflect"
	"time"
)

func GetWantedsHeader(w http.ResponseWriter, r *http.Request) {
	var ret struct {
		DeveloperCount int64               `json:"developer_count"`
		BountiesCount  uint64              `json:"bounties_count"`
		People         *[]db.PersonInShort `json:"people"`
	}
	ret.DeveloperCount = db.DB.CountDevelopers()
	ret.BountiesCount = db.DB.CountBounties()
	ret.People = db.DB.GetPeopleListShort(3)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

func GetListedOffers(w http.ResponseWriter, r *http.Request) {
	people, err := db.DB.GetListedOffers(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(people)
	}
}

func MigrateBounties(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger.Log.Info("Starting bounty migration at %s", startTime.Format(time.RFC3339))

	peeps := db.DB.GetAllPeople()

	for _, peep := range peeps {
		bounties, ok := peep.Extras["wanted"].([]interface{})
		if !ok {
			continue
		}

		for _, bounty := range bounties {
			migrateBounty, ok := bounty.(map[string]interface{})
			if !ok {
				continue
			}

			migrateBountyFinal := db.Bounty{
				Title:                 getString(migrateBounty, "title"),
				OwnerID:               peep.OwnerPubKey,
				Paid:                  getBool(migrateBounty, "paid", false),
				Show:                  getBool(migrateBounty, "show", true),
				Type:                  getString(migrateBounty, "type"),
				Award:                 getString(migrateBounty, "award"),
				Price:                 getUint(migrateBounty, "price"),
				Tribe:                 getString(migrateBounty, "tribe"),
				Created:               getInt64(migrateBounty, "created"),
				Assignee:              getAssignee(peeps, migrateBounty),
				TicketUrl:             getString(migrateBounty, "ticketUrl"),
				Description:           getString(migrateBounty, "description"),
				WantedType:            getString(migrateBounty, "wanted_type"),
				Deliverables:          getString(migrateBounty, "deliverables"),
				CodingLanguages:       getStringArray(migrateBounty, "coding_language"),
				GithubDescription:     getBool(migrateBounty, "github_description", false),
				OneSentenceSummary:    getString(migrateBounty, "one_sentence_summary"),
				EstimatedSessionLength: getString(migrateBounty, "estimated_session_length"),
				EstimatedCompletionDate: getString(migrateBounty, "estimated_completion_date"),
			}

			logger.Log.Info("Adding bounty: %s", migrateBountyFinal.Title)
			db.DB.AddBounty(migrateBountyFinal)
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	logger.Log.Info("Bounty migration completed in %s", duration)
}

func getString(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}

func getBool(data map[string]interface{}, key string, defaultValue bool) bool {
	if value, ok := data[key].(bool); ok {
		return value
	}
	return defaultValue
}

func getUint(data map[string]interface{}, key string) uint {
	if value, ok := data[key].(uint); ok {
		return value
	}
	return 0
}

func getInt64(data map[string]interface{}, key string) int64 {
	if value, ok := data[key].(float64); ok {
		return int64(value)
	}
	return 0
}

func getAssignee(peeps []db.Person, data map[string]interface{}) string {
	if assignee, ok := data["assignee"].(map[string]interface{}); ok {
		if pubkey, exists := assignee["owner_pubkey"].(string); exists {
			for _, peep := range peeps {
				if peep.OwnerPubKey == pubkey {
					return peep.OwnerPubKey
				}
			}
		}
	}
	return ""
}

func getStringArray(data map[string]interface{}, key string) pq.StringArray {
	if value, ok := data[key].(db.PropertyMap); ok {
		if array, exists := value["value"].(pq.StringArray); exists {
			return array
		}
	}
	return pq.StringArray{}
}
