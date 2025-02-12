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

type BountyTimer struct {
	StartTime   *time.Time `json:"start_time"`
	PausedTime  *time.Time `json:"paused_time"`
	TotalPaused time.Duration `json:"total_paused"`
}

var BountyTimers = make(map[string]*BountyTimer)

func startBountyTimer(bountyID string) {
	startTime := time.Now()
	BountyTimers[bountyID] = &BountyTimer{StartTime: &startTime}
}

func pauseBountyTimer(bountyID string) {
	if timer, exists := BountyTimers[bountyID]; exists {
		pausedTime := time.Now()
		timer.PausedTime = &pausedTime
	}
}

func resumeBountyTimer(bountyID string) {
	if timer, exists := BountyTimers[bountyID]; exists && timer.PausedTime != nil {
		pausedDuration := time.Since(*timer.PausedTime)
		timer.TotalPaused += pausedDuration
		timer.PausedTime = nil
	}
}

func closeBountyTimer(bountyID string) {
	delete(BountyTimers, bountyID)
}

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
	peeps := db.DB.GetAllPeople()

	for indexPeep, peep := range peeps {
		logger.Log.Info("peep: %d", indexPeep)
		bounties, ok := peep.Extras["wanted"].([]interface{})

		if !ok {
			logger.Log.Info("Wanted not there")
			continue
		}

		for index, bounty := range bounties {
			logger.Log.Info("looping bounties: %d", index)
			migrateBounty := bounty.(map[string]interface{})

			migrateBountyFinal := db.Bounty{}
			migrateBountyFinal.Title, ok = migrateBounty["title"].(string)
			migrateBountyFinal.OwnerID = peep.OwnerPubKey

			if Paid, ok := migrateBounty["paid"].(bool); ok {
				migrateBountyFinal.Paid = Paid
			} else {
				migrateBountyFinal.Paid = false
			}

			if Show, ok := migrateBounty["show"].(bool); ok {
				migrateBountyFinal.Show = Show
			} else {
				migrateBountyFinal.Show = true
			}

			if Assignee, ok := migrateBounty["assignee"].(map[string]interface{}); ok {
				assigneePubkey := Assignee["owner_pubkey"].(string)
				assigneeId := ""
				for _, peep := range peeps {
					if peep.OwnerPubKey == assigneePubkey {
						assigneeId = peep.OwnerPubKey
					}
				}
				migrateBountyFinal.Assignee = assigneeId
				if assigneeId != "" {
					startBountyTimer(migrateBountyFinal.OwnerID)
				}
			} else {
				migrateBountyFinal.Assignee = ""
			}

			if migrateBountyFinal.Assignee == "" {
				closeBountyTimer(migrateBountyFinal.OwnerID)
			}

			if _, ok := migrateBounty["proof_submitted"].(bool); ok {
				pauseBountyTimer(migrateBountyFinal.OwnerID)
			}

			if _, ok := migrateBounty["feedback_given"].(bool); ok {
				resumeBountyTimer(migrateBountyFinal.OwnerID)
			}

			if _, ok := migrateBounty["completed"].(bool); ok {
				closeBountyTimer(migrateBountyFinal.OwnerID)
			}

			logger.Log.Info("Bounty about to be added ")
			db.DB.AddBounty(migrateBountyFinal)
		}
	}
	return
}
