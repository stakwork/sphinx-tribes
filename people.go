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
	authTimeout := 60
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
	MemeToken   string `json:"memeToken"`
	TribesToken string `json:"tribesToken"`
	Pubkey      string `json:"pubkey"`
	ContactKey  string `json:"contactKey"`
	Alias       string `json:"alias"`
	PhotoURL    string `json:"photoUrl"`
	RouteHint   string `json:"routeHint"`
}

func verify(w http.ResponseWriter, r *http.Request) {
	challenge := chi.URLParam(r, "challenge")
	_, err := store.GetChallenge(challenge)
	if err != nil {
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

	if payload.MemeToken == "" || payload.TribesToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	pubkey, err := VerifyTribeUUID(payload.TribesToken, false)
	if pubkey == "" || err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload.Pubkey = pubkey
	marshalled, err := json.Marshal(payload)
	if err != nil {
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

	if person.ID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	now := time.Now()

	if pubKeyFromAuth == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		person.Created = &now
	}

	person.OwnerPubKey = pubKeyFromAuth
	person.Updated = &now
	person.UniqueName, _ = botUniqueNameFromName(person.OwnerAlias)

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
