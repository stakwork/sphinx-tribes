package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/lib/pq"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/utils"
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

func GetBountiesLeaderboard(w http.ResponseWriter, _ *http.Request) {
	leaderBoard := db.DB.GetBountiesLeaderboard()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(leaderBoard)
}

func DeleteBountyAssignee(w http.ResponseWriter, r *http.Request) {
	invoice := db.DeleteBountyAssignee{}
	body, err := io.ReadAll(r.Body)
	var deletedAssignee bool

	r.Body.Close()

	err = json.Unmarshal(body, &invoice)

	if err != nil {
		utils.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	owner_key := invoice.Owner_pubkey
	date := invoice.Created

	createdUint, _ := strconv.ParseUint(date, 10, 32)
	b, err := db.DB.GetBountyByCreated(uint(createdUint))

	if err == nil && b.OwnerID == owner_key {
		b.Assignee = ""
		b.AssignedHours = 0
		b.CommitmentFee = 0
		b.BountyExpires = ""

		db.DB.UpdateBounty(b)

		deletedAssignee = true
	} else {
		log.Printf("Could not delete bounty assignee")

		deletedAssignee = false

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(deletedAssignee)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deletedAssignee)
}

func MigrateBounties(w http.ResponseWriter, r *http.Request) {
	peeps := db.DB.GetAllPeople()

	for indexPeep, peep := range peeps {
		utils.Log.Info("peep: %d", indexPeep)
		bounties, ok := peep.Extras["wanted"].([]interface{})

		if !ok {
			utils.Log.Info("Wanted not there")
			continue
		}

		for index, bounty := range bounties {

			utils.Log.Info("looping bounties: %d", index)
			migrateBounty := bounty.(map[string]interface{})

			migrateBountyFinal := db.Bounty{}
			migrateBountyFinal.Title, ok = migrateBounty["title"].(string)

			migrateBountyFinal.OwnerID = peep.OwnerPubKey

			Paid, ok1 := migrateBounty["paid"].(bool)
			if !ok1 {
				migrateBountyFinal.Paid = false
			} else {
				migrateBountyFinal.Paid = Paid
			}

			Show, ok2 := migrateBounty["show"].(bool)
			if !ok2 {
				migrateBountyFinal.Show = true
			} else {
				migrateBountyFinal.Show = Show
			}

			Type, ok3 := migrateBounty["type"].(string)
			if !ok3 {
				migrateBountyFinal.Type = ""
			} else {
				migrateBountyFinal.Type = Type
			}

			Award, ok4 := migrateBounty["award"].(string)
			if !ok4 {
				migrateBountyFinal.Award = ""
			} else {
				migrateBountyFinal.Award = Award
			}

			Price, ok5 := migrateBounty["price"].(uint)
			if !ok5 {
				migrateBountyFinal.Price = 0
			} else {
				migrateBountyFinal.Price = Price
			}

			Tribe, ok6 := migrateBounty["tribe"].(string)
			if !ok6 {
				migrateBountyFinal.Tribe = ""
			} else {
				migrateBountyFinal.Tribe = Tribe
			}

			Created, ok7 := migrateBounty["created"].(float64)
			CreatedInt64 := int64(Created)
			if !ok7 {
				migrateBountyFinal.Created = 0
			} else {
				utils.Log.Info("Type: %v", reflect.TypeOf(CreatedInt64))
				utils.Log.Info("Timestamp: %d", CreatedInt64)
				migrateBountyFinal.Created = CreatedInt64
			}

			Assignee, ok8 := migrateBounty["assignee"].(map[string]interface{})
			if !ok8 {
				migrateBountyFinal.Assignee = ""
			} else {
				assigneePubkey := Assignee["owner_pubkey"].(string)
				assigneeId := ""
				for _, peep := range peeps {
					if peep.OwnerPubKey == assigneePubkey {
						assigneeId = peep.OwnerPubKey
					}
				}
				migrateBountyFinal.Assignee = assigneeId
			}

			TicketUrl, ok9 := migrateBounty["ticketUrl"].(string)
			if !ok9 {
				migrateBountyFinal.TicketUrl = ""
			} else {
				migrateBountyFinal.TicketUrl = TicketUrl
			}

			Description, ok10 := migrateBounty["description"].(string)
			if !ok10 {
				migrateBountyFinal.Description = ""
			} else {
				migrateBountyFinal.Description = Description
			}

			WantedType, ok11 := migrateBounty["wanted_type"].(string)
			if !ok11 {
				migrateBountyFinal.WantedType = ""
			} else {
				migrateBountyFinal.WantedType = WantedType
			}

			Deliverables, ok12 := migrateBounty["deliverables"].(string)
			if !ok12 {
				migrateBountyFinal.Deliverables = ""
			} else {
				migrateBountyFinal.Deliverables = Deliverables
			}

			CodingLanguages, ok13 := migrateBounty["coding_language"].(db.PropertyMap)
			if !ok13 {
				migrateBountyFinal.CodingLanguages = pq.StringArray{}
			} else {
				migrateBountyFinal.CodingLanguages = CodingLanguages["value"].(pq.StringArray)
			}

			GithuDescription, ok14 := migrateBounty["github_description"].(bool)
			if !ok14 {
				migrateBountyFinal.GithubDescription = false
			} else {
				migrateBountyFinal.GithubDescription = GithuDescription
			}

			OneSentenceSummary, ok15 := migrateBounty["one_sentence_summary"].(string)
			if !ok15 {
				migrateBountyFinal.OneSentenceSummary = ""
			} else {
				migrateBountyFinal.OneSentenceSummary = OneSentenceSummary
			}

			EstimatedSessionLength, ok16 := migrateBounty["estimated_session_length"].(string)
			if !ok16 {
				migrateBountyFinal.EstimatedSessionLength = ""
			} else {
				migrateBountyFinal.EstimatedSessionLength = EstimatedSessionLength
			}

			EstimatedCompletionDate, ok17 := migrateBounty["estimated_completion_date"].(string)
			if !ok17 {
				migrateBountyFinal.EstimatedCompletionDate = ""
			} else {
				migrateBountyFinal.EstimatedCompletionDate = EstimatedCompletionDate
			}
			utils.Log.Info("Bounty about to be added ")
			db.DB.AddBounty(migrateBountyFinal)
			//Migrate the bounties here
		}
	}
	return
}
