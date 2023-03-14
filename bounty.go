package main

import (
	"encoding/json"
	"net/http"
)

func getAllBounties(w http.ResponseWriter, r *http.Request) {

	bounties := DB.getAllBounties()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bounties)

}
