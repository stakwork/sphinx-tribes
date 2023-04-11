package handlers

import (
	"encoding/json"
	"net/http"

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
		DeveloperCount uint64              `json:"developer_count"`
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
