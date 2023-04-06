package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func getAllBounties(w http.ResponseWriter, r *http.Request) {

	bounties := DB.getAllBounties()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bounties)

}

func createOrEditBounty(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	//pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	bounty := Bounty{}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &bounty)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	now := time.Now()

	/*if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if pubKeyFromAuth != bounty.OwnerID {
		fmt.Println(pubKeyFromAuth)
		fmt.Println(bounty.OwnerID)
		fmt.Println("mismatched pubkey")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}*/
	bounty.Updated = &now
	bounty.Created = time.Now().Unix()

	if bounty.Type == "" {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Type is a required field")
		return
	}
	if bounty.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Title is a required field")
		return
	}
	if bounty.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Description is a required field")
		return
	}
	if bounty.Tribe == "" {
		bounty.Tribe = "None"
	}

	b, err := DB.createOrEditBounty(bounty)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(b)
}
