package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/stakwork/sphinx-tribes/db"
)

// WantedsHeaderResponse represents the response structure for the wanteds header
type WantedsHeaderResponse struct {
	DeveloperCount int64               `json:"developer_count"`
	BountiesCount  uint64              `json:"bounties_count"`
	People         *[]db.PersonInShort `json:"people"`
}

// GetWantedsHeader godoc
//
//	@Summary		Get wanteds header
//	@Description	Get the header information for wanteds
//	@Tags			People
//	@Success		200	{object}	WantedsHeaderResponse
//	@Router			/people/wanteds/header [get]
func GetWantedsHeader(w http.ResponseWriter, r *http.Request) {
	var ret WantedsHeaderResponse
	ret.DeveloperCount = db.DB.CountDevelopers()
	ret.BountiesCount = db.DB.CountBounties()
	ret.People = db.DB.GetPeopleListShort(3)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)
}

// GetListedOffers godoc
//
//	@Summary		Get listed offers
//	@Description	Get a list of listed offers
//	@Tags			People
//	@Success		200	{array}	db.Person
//	@Router			/people/offers [get]
func GetListedOffers(w http.ResponseWriter, r *http.Request) {
	people, err := db.DB.GetListedOffers(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(people)
	}
}
