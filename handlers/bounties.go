package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/stakwork/sphinx-tribes/db"
)

func GetListedWanteds(w http.ResponseWriter, r *http.Request) {
	people, err := db.DB.GetListedWanteds(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(people)
	}
}

func GetPersonAssignedWanteds(w http.ResponseWriter, r *http.Request) {
	people, err := db.DB.GetListedWanteds(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(people)
	}
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

func GetBountiesLeaderboard(w http.ResponseWriter, _ *http.Request) {
	leaderBoard := db.DB.GetBountiesLeaderboard()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(leaderBoard)
}

func DeleteBountyAssignee(w http.ResponseWriter, r *http.Request) {
	invoice := db.DeleteBountyAssignee{}
	body, err := ioutil.ReadAll(r.Body)
	var deletedAssignee bool

	r.Body.Close()

	err = json.Unmarshal(body, &invoice)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	owner_key := invoice.Owner_pubkey
	date := invoice.Created

	var p = db.DB.GetPersonByPubkey(owner_key)

	wanteds, _ := p.Extras["wanted"].([]interface{})

	for _, wanted := range wanteds {
		w, ok2 := wanted.(map[string]interface{})
		if !ok2 {
			continue
		}

		created, ok3 := w["created"].(float64)
		createdArr := strings.Split(fmt.Sprintf("%f", created), ".")
		createdString := createdArr[0]
		createdInt, _ := strconv.ParseInt(createdString, 10, 32)

		dateInt, _ := strconv.ParseInt(date, 10, 32)

		if !ok3 {
			continue
		}

		if createdInt == dateInt {
			delete(w, "assignee")
		}
	}
	p.Extras["wanted"] = wanteds

	b := new(bytes.Buffer)
	decodeErr := json.NewEncoder(b).Encode(p.Extras)

	if decodeErr != nil {
		log.Printf("Could not encode extras json data")

		deletedAssignee = false

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(deletedAssignee)
	} else {
		db.DB.UpdatePerson(p.ID, map[string]interface{}{
			"extras": b,
		})

		deletedAssignee = true

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deletedAssignee)
	}
}

func MigrateBounties(w http.ResponseWriter, r *http.Request) {
	peeps := db.DB.GetAllPeople()

	for indexPeep, peep := range peeps {
		fmt.Println("peep: ", indexPeep)
		bounties, ok := peep.Extras["wanted"].([]interface{})

		if !ok {
			fmt.Println("Wanted not there")
			continue
		}

		for index, bounty := range bounties {

			fmt.Println("looping bounties: ", index)
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

			Price, ok5 := migrateBounty["price"].(string)
			if !ok5 {
				migrateBountyFinal.Price = "0"
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
				now := time.Now().Unix()
				migrateBountyFinal.Created = now
			} else {
				fmt.Println(reflect.TypeOf(CreatedInt64))
				fmt.Println("Timestamp:", CreatedInt64)
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

			CodingLanguage, ok13 := migrateBounty["coding_language"].(db.PropertyMap)
			if !ok13 {
				migrateBountyFinal.CodingLanguage = ""
			} else {
				migrateBountyFinal.CodingLanguage = CodingLanguage["value"].(string)
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
			fmt.Println("Bounty about to be added ")
			db.DB.AddBounty(migrateBountyFinal)
			//Migrate the bounties here
		}
	}
	return
}
