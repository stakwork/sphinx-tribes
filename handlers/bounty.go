package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/stakwork/sphinx-tribes/db"
)

func GetAllBounties(w http.ResponseWriter, r *http.Request) {

	bounties := db.DB.GetAllBounties()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bounties)

}

func CreateOrEditBounty(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	//pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	bounty := db.Bounty{}
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

	b, err := db.DB.CreateOrEditBounty(bounty)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(b)
}
