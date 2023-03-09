package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/xid"
)

const liquidTestModeUrl = "TEST_ASSET_URL"

func createOrEditPerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	person := Person{}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &person)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// ?
	// if person.ID == 0 {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	now := time.Now()

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if pubKeyFromAuth != person.OwnerPubKey {
		fmt.Println(pubKeyFromAuth)
		fmt.Println(person.OwnerPubKey)
		fmt.Println("mismatched pubkey")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	existing := DB.getPersonByPubkey(pubKeyFromAuth)
	if existing.ID == 0 { // new!
		if person.ID != 0 { // cant try to "edit" if not exists already
			fmt.Println("cant edit non existing")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		person.UniqueName, _ = personUniqueNameFromName(person.OwnerAlias)
		person.Created = &now
		person.Uuid = xid.New().String()
	} else { // editing! needs ID
		if person.ID == 0 { // cant create that already exists
			fmt.Println("cant create existing")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if person.ID != existing.ID { // cant edit someone else's
			fmt.Println("cant edit someone else")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	person.OwnerPubKey = pubKeyFromAuth
	person.Updated = &now
	if person.NewTicketTime != 0 {
		go processAlerts(person)
	}
	p, err := DB.createOrEditPerson(person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

func personIsAdmin(pk string) bool {
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

func deleteTicketByAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)
	pubKey := chi.URLParam(r, "pubKey")
	createdStr := chi.URLParam(r, "created")
	created, err := strconv.ParseInt(createdStr, 10, 64)
	if err != nil {
		fmt.Println("Unable to convert created to int64")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if created == 0 || pubKey == "" {
		fmt.Println("Insufficient details to delete ticket")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	existing := DB.getPersonByPubkey(pubKeyFromAuth)
	if existing.ID == 0 {
		fmt.Println("Could not fetch admin details from db")
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if personIsAdmin(existing.OwnerPubKey) == false {
		fmt.Println("Only admin is allowed to delete tickets")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	person := DB.getPersonByPubkey(pubKey)
	if person.ID == 0 {
		fmt.Println("Could not fetch person from db")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	wanteds, ok := person.Extras["wanted"].([]interface{})
	if !ok {
		fmt.Println("No tickets found for person")
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
		fmt.Println("Ticket to delete not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		person.Extras["wanted"] = append(wanteds[:index], wanteds[index+1:]...)
	}

	_, err = DB.createOrEditPerson(person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}

func processTwitterConfirmationsLoop() {
	twitterToken := os.Getenv("TWITTER_TOKEN")
	if twitterToken == "" {
		return
	}
	peeps := DB.getUnconfirmedTwitter()
	for _, p := range peeps {
		twitArray, ok := p.Extras["twitter"].([]interface{})
		if ok {
			if len(twitArray) > 0 {
				twitValue, ok2 := twitArray[0].(map[string]interface{})
				if ok2 {
					username, _ := twitValue["value"].(string)
					if username != "" {
						pubkey, err := ConfirmIdentityTweet(username)
						// fmt.Println("TWitter err", err)
						if err == nil && pubkey != "" {
							if p.OwnerPubKey == pubkey {
								DB.updateTwitterConfirmed(p.ID, true)
							}
						}
					}
				}
			}
		}
	}
	time.Sleep(30 * time.Second)
	processTwitterConfirmationsLoop()
}

func processGithubIssuesLoop() {
	peeps := DB.getListedPeople(nil)

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
			// does githubissue already have a status here, and is it different?
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
			DB.updateGithubIssues(p.ID, clonedGithubIssues)
		}
	}
	time.Sleep(1 * time.Minute)
	processGithubIssuesLoop()
}

func processGithubConfirmationsLoop() {
	peeps := DB.getUnconfirmedGithub()
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
								DB.updateGithubConfirmed(p.ID, true)
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

func personUniqueNameFromName(name string) (string, error) {
	pathOne := strings.ToLower(strings.Join(strings.Fields(name), ""))
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	path := reg.ReplaceAllString(pathOne, "")
	n := 0
	for {
		uniquepath := path
		if n > 0 {
			uniquepath = path + strconv.Itoa(n)
		}
		existing := DB.getPersonByUniqueName(uniquepath)
		if existing.ID != 0 {
			n = n + 1
		} else {
			path = uniquepath
			break
		}
	}
	return path, nil
}

func getPersonByPubkey(w http.ResponseWriter, r *http.Request) {
	pubkey := chi.URLParam(r, "pubkey")
	person := DB.getPersonByPubkey(pubkey)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
}

func getPersonByUuid(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	person := DB.getPersonByUuid(uuid)
	assetBalanceData, err := getAssetByPubkey(person.OwnerPubKey)

	personResponse := make(map[string]interface{})
	personResponse["id"] = person.ID
	personResponse["uuid"] = person.Uuid
	personResponse["owner_alias"] = person.OwnerAlias
	personResponse["unique_name"] = person.UniqueName
	personResponse["description"] = person.Description
	personResponse["tags"] = person.Tags
	personResponse["img"] = person.Img
	personResponse["owner_route_hint"] = person.OwnerRouteHint
	personResponse["owner_contact_key"] = person.OwnerContactKey
	personResponse["price_to_meet"] = person.PriceToMeet
	// personResponse["extras"] = person.Extras
	personResponse["twitter_confirmed"] = person.TwitterConfirmed
	personResponse["github_issues"] = person.GithubIssues
	if err != nil {
		fmt.Println("==> error: ", err)
	} else {
		var badgeSlice []uint
		for i := 0; i < len(assetBalanceData); i++ {
			badgeSlice = append(badgeSlice, assetBalanceData[i].AssetId)
		}
		personResponse["badges"] = badgeSlice
	}
	fmt.Println()
	// FIXME use http to hit sphinx-element server for badges
	// Todo: response should include no pubKey
	// FIXME also filter by the tribe "profile_filters"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(personResponse)
}

func getPersonAssetsByUuid(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	person := DB.getPersonByUuid(uuid)
	assetList, err := getAssetList(person.OwnerPubKey)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	fmt.Println()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(assetList)
}

func getPersonByGithubName(w http.ResponseWriter, r *http.Request) {
	github := chi.URLParam(r, "github")
	person := DB.getPersonByGithubName(github)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if id == 0 {
		fmt.Println("id is 0")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	existing := DB.getPerson(uint(id))
	if existing.ID == 0 {
		fmt.Println("existing id is 0")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if existing.OwnerPubKey != pubKeyFromAuth {
		fmt.Println("keys dont match")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	DB.updatePerson(uint(id), map[string]interface{}{
		"deleted": true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

func getAssetByPubkey(pubkey string) ([]AssetBalanceData, error) {
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
		fmt.Println("GET error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var r AssetResponse
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("json unmarshall error", err)
		return nil, err
	}

	balances := r.Balances

	return balances, nil
}

func getAssetList(pubkey string) ([]AssetListData, error) {
	client := &http.Client{}

	url := os.Getenv("ASSET_LIST_URL")
	if url == "" {
		url = "https://liquid.sphinx.chat/assets"
	}

	url = url + "?pubkey=" + pubkey

	req, err := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("GET error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var r []AssetListData
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("json unmarshall error", err)
		return nil, err
	}

	return r, nil
}

func addOrRemoveBadge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	badgeCreationData := BadgeCreationData{}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &badgeCreationData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if badgeCreationData.Badge == "" {
		fmt.Println("Badge cannot be Empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if badgeCreationData.Action == "" {
		fmt.Println("Action cannot be Empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !(badgeCreationData.Action == "add" || badgeCreationData.Action == "remove") {
		fmt.Println("Invalid action in Request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if badgeCreationData.TribeUUID == "" {
		fmt.Println("tribeId cannot be Empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	extractedPubkey, err := VerifyTribeUUID(badgeCreationData.TribeUUID, false)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tribe := DB.getTribeByIdAndPubkey(badgeCreationData.TribeUUID, extractedPubkey)

	if pubKeyFromAuth != tribe.OwnerPubKey {
		fmt.Println(pubKeyFromAuth)
		fmt.Println("mismatched pubkey")
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
	updatedTribe := DB.updateTribe(tribe.UUID, map[string]interface{}{
		"badges": tribeBadges,
	})

	if updatedTribe {
		tribe = DB.getTribeByIdAndPubkey(badgeCreationData.TribeUUID, extractedPubkey)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tribe)
	}
	w.WriteHeader(http.StatusBadRequest)
	return

}

func migrateBounties(w http.ResponseWriter, r *http.Request) {
	peeps := DB.getListedPeople(nil)

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

			migrateBountyFinal := Bounty{}
			migrateBountyFinal.Title, ok = migrateBounty["title"].(string)

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

			Created, ok7 := migrateBounty["created"].(*time.Time)
			if !ok7 {
				now := time.Now()
				migrateBountyFinal.Created = &now
			} else {
				migrateBountyFinal.Created = Created
			}

			Assignee, ok8 := migrateBounty["assignee"].(Person)
			if !ok8 {
				migrateBountyFinal.Assignee = Person{}
			} else {
				migrateBountyFinal.Assignee = Assignee
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

			CodingLanguage, ok13 := migrateBounty["coding_language"].(PropertyMap)
			if !ok13 {
				migrateBountyFinal.CodingLanguage = ""
			} else {
				migrateBountyFinal.CodingLanguage = CodingLanguage["value"].(string)
			}

			GithuDescription, ok14 := migrateBounty["github_description"].(bool)
			if !ok14 {
				migrateBountyFinal.GithuDescription = false
			} else {
				migrateBountyFinal.GithuDescription = GithuDescription
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
			DB.addBounty(migrateBountyFinal)
			//Migrate the bounties here

		}
	}
	return
}
