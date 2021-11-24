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
)

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

	p, err := DB.createOrEditPerson(person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
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
			DB.updateGithubIssues(p.ID, map[string]interface{}{
				fullissuename: map[string]string{
					"assignee": issue.Assignee,
					"status":   issue.Status,
				},
			})
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
