package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/xid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
	"github.com/stakwork/sphinx-tribes/utils"
)

const liquidTestModeUrl = "TEST_ASSET_URL"

type peopleHandler struct {
	db db.Database
}

func NewPeopleHandler(db db.Database) *peopleHandler {
	return &peopleHandler{db: db}
}

// CreatePerson godoc
//
//	@Summary		Create Person
//	@Description	Create a new person
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			referred_by	query		string		false	"Referred By"
//	@Param			person		body		db.Person	true	"Person"
//	@Success		200			{object}	db.Person
//	@Router			/people [post]
func (ph *peopleHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	keys := r.URL.Query()
	referredBy := keys.Get("referred_by")

	person := db.Person{}
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Sent wrong body data")
		return
	}

	r.Body.Close()
	err = json.Unmarshal(body, &person)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	now := time.Now()

	if pubKeyFromAuth == "" {
		log.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if pubKeyFromAuth != person.OwnerPubKey {
		log.Println(pubKeyFromAuth)
		log.Println(person.OwnerPubKey)
		log.Println("mismatched pubkey")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	existing := ph.db.GetPersonByPubkey(pubKeyFromAuth)
	if existing.ID == 0 {

		person.UniqueName, _ = ph.db.PersonUniqueNameFromName(person.OwnerAlias)
		person.Created = &now
		person.Uuid = xid.New().String()

		if referredBy != "" {
			// get the referral and populate the pubkey
			referral := db.DB.GetPersonByUuid(referredBy)
			// if referral exists
			if referral.ID != 0 {
				person.ReferredBy = referral.ID
			}
		}
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existing)
		return
	}

	person.OwnerPubKey = pubKeyFromAuth

	if person.NewTicketTime != 0 {
		go ph.db.ProcessAlerts(person)
	}

	b := new(bytes.Buffer)
	decodeErr := json.NewEncoder(b).Encode(person.Extras)

	if decodeErr != nil {
		log.Printf("Could not encode extras json data")
	}

	p, err := ph.db.CreateOrEditPerson(person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send email notification for new user profile creation
	// This is only reached for new users since existing users return early
	go ph.db.SendNewUserNotification(p)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

// UpdatePerson godoc
//
//	@Summary		Update Person
//	@Description	Update an existing person
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			person	body		db.Person	true	"Person"
//	@Success		200		{object}	db.Person
//	@Router			/people [put]
func (ph *peopleHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	person := db.Person{}
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Sent wrong body data")
		return
	}

	r.Body.Close()
	err = json.Unmarshal(body, &person)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	now := time.Now()

	if pubKeyFromAuth == "" {
		log.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if pubKeyFromAuth != person.OwnerPubKey {
		log.Println(pubKeyFromAuth)
		log.Println(person.OwnerPubKey)
		log.Println("mismatched pubkey")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	existing := ph.db.GetPersonByPubkey(pubKeyFromAuth)
	if existing.ID == 0 {
		msg := fmt.Sprintf("User does not exists: %s", pubKeyFromAuth)
		log.Println(msg)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(msg)
		return
	} else {
		if person.OwnerPubKey != existing.OwnerPubKey && person.OwnerAlias != existing.OwnerAlias {
			// can't edit someone else's
			logger.Log.Info("cant edit someone else")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	person.Updated = &now

	if person.NewTicketTime != 0 {
		go ph.db.ProcessAlerts(person)
	}

	b := new(bytes.Buffer)
	decodeErr := json.NewEncoder(b).Encode(person.Extras)

	if decodeErr != nil {
		log.Printf("Could not encode extras json data")
	}

	p, err := ph.db.CreateOrEditPerson(person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

// UpsertLogin godoc
//
//	@Summary		Upsert Login
//	@Description	Upsert login for a person
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Security		CypressAuth
//	@Param			person	body		db.Person	true	"Person"
//	@Success		200		{string}	string		"JWT Token"
//	@Router			/person/login [post]
func (ph *peopleHandler) UpsertLogin(w http.ResponseWriter, r *http.Request) {
	person := db.Person{}
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Sent wrong body data")
		return
	}

	r.Body.Close()
	err = json.Unmarshal(body, &person)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	now := time.Now()

	pubKeyFromAuth := person.OwnerPubKey

	existing := ph.db.GetPersonByPubkey(pubKeyFromAuth)
	if existing.ID == 0 {
		if person.ID != 0 {
			// cant try to "edit" if not exists already
			logger.Log.Info("cant edit non existing")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		person.UniqueName, _ = ph.db.PersonUniqueNameFromName(person.OwnerAlias)
		person.Created = &now
		person.Uuid = xid.New().String()

	} else { // editing! needs ID
		if person.ID != 0 && person.ID != existing.ID { // can't edit someone else's
			logger.Log.Info("cant edit someone else")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	person.OwnerPubKey = pubKeyFromAuth
	person.Updated = &now

	if person.NewTicketTime != 0 {
		go ph.db.ProcessAlerts(person)
	}

	b := new(bytes.Buffer)
	decodeErr := json.NewEncoder(b).Encode(person.Extras)

	if decodeErr != nil {
		log.Printf("Could not encode extras json data")
	}

	p, err := ph.db.CreateOrEditPerson(person)
	_ = p

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseData := make(map[string]interface{})
	tokenString, err := auth.EncodeJwt(person.OwnerPubKey)

	if err != nil {
		logger.Log.Info("Cannot generate jwt token")
	}

	responseData["jwt"] = tokenString

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))
}

func PersonIsAdmin(pk string) bool {
	adminPubkeys := os.Getenv("ADMIN_PUBKEYS")
	if adminPubkeys == "" {
		return false
	}
	admins := strings.Split(adminPubkeys, ",")
	for _, admin := range admins {
		if admin == pk {
			return true
		}
	}
	return false
}

// DeleteTicketByAdmin godoc
//
//	@Summary		Delete Ticket by Admin
//	@Description	Delete a ticket by admin
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			pubKey	path		string	true	"Public Key"
//	@Param			created	path		int64	true	"Created Timestamp"
//	@Success		200		{string}	string	"Ticket deleted successfully"
//	@Router			/ticket/{pubKey}/{created} [delete]
func DeleteTicketByAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	pubKey := chi.URLParam(r, "pubKey")
	createdStr := chi.URLParam(r, "created")
	created, err := strconv.ParseInt(createdStr, 10, 64)
	if err != nil {
		logger.Log.Info("Unable to convert created to int64")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if created == 0 || pubKey == "" {
		logger.Log.Info("Insufficient details to delete ticket")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	existing := db.DB.GetPersonByPubkey(pubKeyFromAuth)
	if existing.ID == 0 {
		logger.Log.Info("Could not fetch admin details from db")
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if PersonIsAdmin(existing.OwnerPubKey) == false {
		logger.Log.Info("Only admin is allowed to delete tickets")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	person := db.DB.GetPersonByPubkey(pubKey)
	if person.ID == 0 {
		logger.Log.Info("Could not fetch person from db")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	wanteds, ok := person.Extras["wanted"].([]interface{})
	if !ok {
		logger.Log.Info("No tickets found for person")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	index := -1
	for i, wanted := range wanteds {
		w, ok2 := wanted.(map[string]interface{})
		if !ok2 {
			continue
		}
		timeF, ok3 := w["created"].(float64)
		if !ok3 {
			continue
		}
		t := int64(timeF)
		if t == created {
			index = i
			break
		}
	}

	if index == -1 {
		logger.Log.Info("Ticket to delete not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		person.Extras["wanted"] = append(wanteds[:index], wanteds[index+1:]...)
	}

	b := new(bytes.Buffer)
	decodeErr := json.NewEncoder(b).Encode(person.Extras)

	if decodeErr != nil {
		log.Printf("Could not encode extras json data")
	}

	_, err = db.DB.CreateOrEditPerson(person)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func ProcessTwitterConfirmationsLoop() {
	twitterToken := os.Getenv("TWITTER_TOKEN")
	if twitterToken == "" {
		return
	}
	peeps := db.DB.GetUnconfirmedTwitter()
	for _, p := range peeps {
		twitArray, ok := p.Extras["twitter"].([]interface{})
		if ok {
			if len(twitArray) > 0 {
				twitValue, ok2 := twitArray[0].(map[string]interface{})
				if ok2 {
					username, _ := twitValue["value"].(string)
					if username != "" {
						pubkey, err := utils.ConfirmIdentityTweet(username)
						// fmt.Println("Twitter err", err)
						if err == nil && pubkey != "" {
							if p.OwnerPubKey == pubkey {
								db.DB.UpdateTwitterConfirmed(p.ID, true)
							}
						}
					}
				}
			}
		}
	}
	time.Sleep(30 * time.Second)
	ProcessTwitterConfirmationsLoop()
}

func ProcessGithubIssuesLoop() {
	peeps := db.DB.GetListedPeople(nil)

	for _, p := range peeps {
		wanteds, ok := p.Extras["wanted"].([]interface{})
		if !ok {
			continue // next person
		}
		for _, wanted := range wanteds {
			w, ok2 := wanted.(map[string]interface{})
			if !ok2 {
				continue // next wanted
			}
			repo, ok3 := w["repo"].(string)
			issnum, ok4 := w["issue"].(string)
			if !ok3 || !ok4 {
				continue
			}
			if !strings.Contains(repo, "/") {
				continue
			}
			arr := strings.Split(repo, "/")
			owner := arr[0]
			reponame := arr[1]
			issint, err := strconv.Atoi(issnum)
			if issint < 1 || err != nil {
				continue
			}
			issue, issueErr := GetIssue(owner, reponame, issint)
			if issueErr != nil {
				continue
			}
			fullissuename := owner + "/" + reponame + "/" + issnum

			// scan original github issue and replace existing or add, if no new info then don't update
			// does github issue already have a status here, and is it different?
			if _, ok5 := p.GithubIssues[fullissuename]; ok5 {

				if w, ok6 := p.GithubIssues[fullissuename].(map[string]interface{}); ok6 {

					assignee, ok7 := w["assignee"].(string)
					status, ok8 := w["status"].(string)

					if ok7 || ok8 {
						//if there are no changes to this ticket, then skip it
						if status == issue.Status && assignee == issue.Assignee {
							continue
						}
					}
				}
			}

			clonedGithubIssues := p.GithubIssues
			// map new values to proper key
			clonedGithubIssues[fullissuename] = map[string]string{
				"assignee": issue.Assignee,
				"status":   issue.Status,
			}

			// update with altered record
			db.DB.UpdateGithubIssues(p.ID, clonedGithubIssues)
		}
	}
	time.Sleep(1 * time.Minute)
	ProcessGithubIssuesLoop()
}

func processGithubConfirmationsLoop() {
	peeps := db.DB.GetUnconfirmedGithub()
	for _, p := range peeps {
		gitArray, ok := p.Extras["twitter"].([]interface{})
		if ok {
			if len(gitArray) > 0 {
				gitValue, ok2 := gitArray[0].(map[string]interface{})
				if ok2 {
					username, _ := gitValue["value"].(string)
					if username != "" {
						pubkey, err := PubkeyForGithubUser(username)
						if err == nil && pubkey != "" {
							if p.OwnerPubKey == pubkey {
								db.DB.UpdateGithubConfirmed(p.ID, true)
							}
						}
					}
				}
			}
		}
	}
	time.Sleep(30 * time.Second)
	processGithubConfirmationsLoop()
}

// GetPersonByPubkey godoc
//
//	@Summary		Get Person by Pubkey
//	@Description	Get a person by their public key
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Param			pubkey	path		string	true	"Public Key"
//	@Success		200		{object}	db.Person
//	@Router			/person/{pubkey} [get]
func (ph *peopleHandler) GetPersonByPubkey(w http.ResponseWriter, r *http.Request) {
	pubkey := chi.URLParam(r, "pubkey")

	person := ph.db.GetPersonByPubkey(pubkey)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
}

// GetPersonById godoc
//
//	@Summary		Get Person by ID
//	@Description	Get a person by their ID
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint	true	"ID"
//	@Success		200	{object}	db.Person
//	@Router			/person/id/{id} [get]
func (ph *peopleHandler) GetPersonById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(idParam, 10, 32)

	person := ph.db.GetPerson(uint(id))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
}

// GetPersonByUuid godoc
//
//	@Summary		Get Person by UUID
//	@Description	Get a person by their UUID
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string	true	"UUID"
//	@Success		200		{object}	map[string]interface{}
//	@Router			/person/uuid/{uuid} [get]
func (ph *peopleHandler) GetPersonByUuid(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	person := ph.db.GetPersonByUuid(uuid)
	assetBalanceData, err := GetAssetByPubkey(person.OwnerPubKey)

	personResponse := make(map[string]interface{})
	personResponse["id"] = person.ID
	personResponse["uuid"] = person.Uuid
	personResponse["owner_pubkey"] = person.OwnerPubKey
	personResponse["owner_alias"] = person.OwnerAlias
	personResponse["unique_name"] = person.UniqueName
	personResponse["description"] = person.Description
	personResponse["tags"] = person.Tags
	personResponse["img"] = person.Img
	personResponse["owner_route_hint"] = person.OwnerRouteHint
	personResponse["owner_contact_key"] = person.OwnerContactKey
	personResponse["price_to_meet"] = person.PriceToMeet
	personResponse["twitter_confirmed"] = person.TwitterConfirmed
	personResponse["github_issues"] = person.GithubIssues
	if err != nil {
		logger.Log.Error("==> error: %v", err)
	} else {
		var badgeSlice []uint
		for i := 0; i < len(assetBalanceData); i++ {
			badgeSlice = append(badgeSlice, assetBalanceData[i].AssetId)
		}
		personResponse["badges"] = badgeSlice
	}
	logger.Log.Info("")
	// FIXME use http to hit sphinx-element server for badges
	// Todo: response should include no pubKey
	// FIXME also filter by the tribe "profile_filters"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(personResponse)
}

// GetPersonAssetsByUuid godoc
//
//	@Summary		Get Person Assets by UUID
//	@Description	Get assets of a person by their UUID
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path	string	true	"UUID"
//	@Success		200		{array}	db.AssetListData
//	@Router			/person/assets/{uuid} [get]
func GetPersonAssetsByUuid(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	person := db.DB.GetPersonByUuid(uuid)
	assetList, err := GetAssetList(person.OwnerPubKey)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	logger.Log.Info("")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(assetList)
}

// GetPersonByGithubName godoc
//
//	@Summary		Get Person by Github Name
//	@Description	Get a person by their Github name
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Param			github	path		string	true	"Github Name"
//	@Success		200		{object}	db.Person
//	@Router			/person/github/{github} [get]
func GetPersonByGithubName(w http.ResponseWriter, r *http.Request) {
	github := chi.URLParam(r, "github")
	person := db.DB.GetPersonByGithubName(github)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
}

// DeletePerson godoc
//
//	@Summary		Delete Person
//	@Description	Delete a person by their ID
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path		int		true	"ID"
//	@Success		200	{string}	string	"Person deleted successfully"
//	@Router			/person/{id} [delete]
func (ph *peopleHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if id == 0 {
		logger.Log.Info("id is 0")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	existing := ph.db.GetPerson(uint(id))
	if existing.ID == 0 {
		logger.Log.Info("existing id is 0")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if existing.OwnerPubKey != pubKeyFromAuth {
		logger.Log.Info("keys dont match")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ph.db.UpdatePerson(uint(id), map[string]interface{}{
		"deleted": true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

func GetAssetByPubkey(pubkey string) ([]db.AssetBalanceData, error) {
	client := &http.Client{}
	testMode, err := strconv.ParseBool(os.Getenv("TEST_MODE"))
	if err != nil {
		testMode = false
	}

	url := os.Getenv(liquidTestModeUrl)
	if testMode && (url != "") {
		url = os.Getenv(liquidTestModeUrl)
	} else {
		url = "https://liquid.sphinx.chat/balances?pubkey=" + pubkey
	}

	req, err := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)

	if err != nil {
		logger.Log.Error("GET error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var r db.AssetResponse
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		logger.Log.Error("json unmarshall error: %v", err)
		return nil, err
	}

	balances := r.Balances

	return balances, nil
}

func GetAssetList(pubkey string) ([]db.AssetListData, error) {
	client := &http.Client{}

	url := os.Getenv("ASSET_LIST_URL")
	if url == "" {
		url = "https://liquid.sphinx.chat/assets"
	}

	url = url + "?pubkey=" + pubkey

	req, err := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)

	if err != nil {
		logger.Log.Error("GET error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var r []db.AssetListData
	body, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &r)
	if err != nil {
		logger.Log.Error("json unmarshall error: %v", err)
		return nil, err
	}

	return r, nil
}

// AddOrRemoveBadge godoc
//
//	@Summary		Add or Remove Badge
//	@Description	Add or remove a badge for a person
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			badgeCreationData	body		db.BadgeCreationData	true	"Badge Creation Data"
//	@Success		200					{object}	db.Tribe
//	@Router			/badges [post]
func AddOrRemoveBadge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	badgeCreationData := db.BadgeCreationData{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &badgeCreationData)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if badgeCreationData.Badge == "" {
		logger.Log.Info("Badge cannot be Empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if badgeCreationData.Action == "" {
		logger.Log.Info("Action cannot be Empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !(badgeCreationData.Action == "add" || badgeCreationData.Action == "remove") {
		logger.Log.Info("Invalid action in Request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if badgeCreationData.TribeUUID == "" {
		logger.Log.Info("tribeId cannot be Empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	extractedPubkey, err := auth.VerifyTribeUUID(badgeCreationData.TribeUUID, false)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tribe := db.DB.GetTribeByIdAndPubkey(badgeCreationData.TribeUUID, extractedPubkey)

	if pubKeyFromAuth != tribe.OwnerPubKey {
		logger.Log.Info("%s", pubKeyFromAuth)
		logger.Log.Info("mismatched pubkey")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tribeBadges := tribe.Badges
	if tribeBadges == nil {
		tribeBadges = []string{}
	}
	if badgeCreationData.Action == "add" {
		badges := append(tribeBadges, badgeCreationData.Badge)
		tribeBadges = badges
	}

	if badgeCreationData.Action == "remove" {
		for i, v := range tribeBadges {
			if strings.ToLower(v) == strings.ToLower(badgeCreationData.Badge) {
				tribeBadges = append(tribeBadges[:i], tribeBadges[i+1:]...)
				break
			}
		}

	}

	tribe.Badges = tribeBadges
	updatedTribe := db.DB.UpdateTribe(tribe.UUID, map[string]interface{}{
		"badges": tribeBadges,
	})

	if updatedTribe {
		tribe = db.DB.GetTribeByIdAndPubkey(badgeCreationData.TribeUUID, extractedPubkey)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tribe)
	}
	w.WriteHeader(http.StatusBadRequest)
	return
}

// GetPeopleShortList godoc

// @Summary		Get People Short List
// @Description	Get a short list of people
// @Tags			People
// @Accept			json
// @Produce		json
// @Success		200	{array}	db.Person
// @Router			/person/short [get]
func GetPeopleShortList(w http.ResponseWriter, r *http.Request) {
	var maxCount uint32 = 10000
	people := db.DB.GetPeopleListShort(maxCount)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
}

// GetPeopleBySearch godoc
//
//	@Summary		Get People by Search
//	@Description	Get people by search query
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	db.Person
//	@Router			/person/search [get]
func (ph *peopleHandler) GetPeopleBySearch(w http.ResponseWriter, r *http.Request) {
	people := ph.db.GetPeopleBySearch(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
}

// GetListedPeople godoc
//
//	@Summary		Get Listed People
//	@Description	Get listed people
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	db.Person
//	@Router			/people [get]
func (ph *peopleHandler) GetListedPeople(w http.ResponseWriter, r *http.Request) {
	people := ph.db.GetListedPeople(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
}

// GetListedPosts godoc
//
//	@Summary		Get Listed Posts
//	@Description	Get listed posts
//	@Tags			People
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	db.PeopleExtra
//	@Router			/person/posts [get]
func GetListedPosts(w http.ResponseWriter, r *http.Request) {
	people, err := db.DB.GetListedPosts(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(people)
	}
}
