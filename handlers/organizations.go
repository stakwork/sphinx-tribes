package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/xid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
)

func CreateOrEditOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	now := time.Now()

	org := db.Organization{}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &org)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if pubKeyFromAuth != org.OwnerPubKey {
		fmt.Println(pubKeyFromAuth)
		fmt.Println(org.OwnerPubKey)
		fmt.Println("mismatched pubkey")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	existing := db.DB.GetOrganizationByUuid(org.Uuid)
	if existing.ID == 0 { // new!
		if org.ID != 0 { // cant try to "edit" if not exists already
			fmt.Println("cant edit non existing")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		name := strings.ToLower(org.Name)

		// check if the organization name already exists
		orgName := db.DB.GetOrganizationByName(name)

		if orgName.Name == name {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Organization name alreday exists")
			return
		} else {
			org.Created = &now
			org.Updated = &now
			org.Uuid = xid.New().String()
			org.Name = name
		}
	} else {
		if org.ID == 0 {
			// cant create that already exists
			fmt.Println("cant create existing")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if org.ID != existing.ID { // cant edit someone else's
			fmt.Println("cant edit someone else")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	p, err := db.DB.CreateOrEditOrganization(org)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

func GetOrganizations(w http.ResponseWriter, r *http.Request) {
	orgs := db.DB.GetOrganizations(r)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orgs)
}

func GetOrganizationByUuid(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	org := db.DB.GetOrganizationByUuid(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(org)
}

func CreateOrganizationUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	orgUser := db.OrganizationUsers{}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &orgUser)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	org := db.DB.GetOrganizationByUuid(orgUser.Organization)

	// if not the orgnization admin
	if pubKeyFromAuth != org.OwnerPubKey {
		// todo check if the user as create user access
		fmt.Println("don't have access to create auser")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check if user already exists
	userExists := db.DB.GetOrganizationUser(orgUser.OwnerPubKey, orgUser.Organization)

	if userExists.ID != 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("User already exists")
		return
	}

	// create user
	user, _ := db.DB.CreateOrganizationUser(orgUser)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func GetOrganizationUsers(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	orgUsers := db.DB.GetOrganizationUsers(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orgUsers)
}
