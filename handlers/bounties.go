package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

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
