package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/blake2b"
)

// Store struct
type Store struct {
	cache *cache.Cache
}

var store Store

func initCache() {
	authTimeout := 120
	store = Store{
		cache: cache.New(
			time.Duration(authTimeout)*time.Second,
			time.Duration(authTimeout*3)*time.Second,
		),
	}
}

// SetChallenge
func (s Store) SetChallenge(key string, value string) error {
	s.cache.Set(key, value, cache.DefaultExpiration)
	return nil
}

// DeleteChallenge
func (s Store) DeleteChallenge(key string) error {
	s.cache.Delete(key)
	return nil
}

// GetChallenge
func (s Store) GetChallenge(key string) (string, error) {
	value, found := s.cache.Get(key)
	c, _ := value.(string)
	if !found || c == "" {
		return "", errors.New("not found")
	}
	return c, nil
}

func ask(w http.ResponseWriter, r *http.Request) {
	ts := strconv.Itoa(int(time.Now().Unix()))
	h := blake2b.Sum256([]byte(ts))
	challenge := base64.URLEncoding.EncodeToString(h[:])

	store.SetChallenge(challenge, ts)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"challenge": challenge,
		"ts":        ts,
	})
}

type VerifyPayload struct {
	ID          uint   `json:"id"`
	Pubkey      string `json:"pubkey"`
	ContactKey  string `json:"contact_key"`
	Alias       string `json:"alias"`
	PhotoURL    string `json:"photo_url"`
	RouteHint   string `json:"route_hint"`
	PriceToMeet uint   `json:"price_to_meet"`
	JWT         string `json:"jwt"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

func verify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	challenge := chi.URLParam(r, "challenge")
	_, err := store.GetChallenge(challenge)
	if err != nil {
		fmt.Println("challenge not found", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload := VerifyPayload{}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &payload)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	payload.Pubkey = pubKeyFromAuth
	marshalled, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("payload unparseable", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// set into the cache
	store.SetChallenge(challenge, string(marshalled))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}

/*
curl localhost:5002/ask
curl localhost:5002/poll/d5SYZNY5pQ7dXwHP-oXh2uSOPUEX0fUJOXI0_5-eOsg=
*/
func poll(w http.ResponseWriter, r *http.Request) {

	challenge := chi.URLParam(r, "challenge")
	res, err := store.GetChallenge(challenge)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if len(res) <= 10 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	pld := VerifyPayload{}
	err = json.Unmarshal([]byte(res), &pld)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if pld.Pubkey == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	existing := DB.getPersonByPubkey(pld.Pubkey)
	if existing.ID > 0 {
		pld.ID = existing.ID // add ID on if exists
		pld.Description = existing.Description
	}

	store.DeleteChallenge(challenge)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pld)
}

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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
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
