package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

// SetCache
func (s Store) SetCache(key string, value string) error {
	s.cache.Set(key, value, cache.DefaultExpiration)
	return nil
}

// DeleteCache
func (s Store) DeleteCache(key string) error {
	s.cache.Delete(key)
	return nil
}

// GetCache
func (s Store) GetCache(key string) (string, error) {
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

	store.SetCache(challenge, ts)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"challenge": challenge,
		"ts":        ts,
	})
}

type VerifyPayload struct {
	ID                    uint                   `json:"id"`
	Pubkey                string                 `json:"pubkey"`
	ContactKey            string                 `json:"contact_key"`
	Alias                 string                 `json:"alias"`
	PhotoURL              string                 `json:"photo_url"`
	RouteHint             string                 `json:"route_hint"`
	PriceToMeet           uint                   `json:"price_to_meet"`
	JWT                   string                 `json:"jwt"`
	URL                   string                 `json:"url"`
	Description           string                 `json:"description"`
	VerificationSignature string                 `json:"verification_signature"`
	Extras                map[string]interface{} `json:"extras"`
}

func verify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(ContextKey).(string)

	challenge := chi.URLParam(r, "challenge")
	_, err := store.GetCache(challenge)
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
	store.SetCache(challenge, string(marshalled))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}

/*
curl localhost:5002/ask
curl localhost:5002/poll/d5SYZNY5pQ7dXwHP-oXh2uSOPUEX0fUJOXI0_5-eOsg=
*/
func poll(w http.ResponseWriter, r *http.Request) {

	challenge := chi.URLParam(r, "challenge")
	res, err := store.GetCache(challenge)
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
		pld.Extras = existing.Extras
		// standardize language for frontend, retrun photo_url for img
		if existing.Img != "" {
			pld.PhotoURL = existing.Img
		}
		// standardize language for frontend, return alias for img
		if existing.OwnerAlias != "" {
			pld.Alias = existing.OwnerAlias
		}
	}

	// store.DeleteChallenge(challenge)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pld)
}

type Save struct {
	Key  string `json:"key"`
	Body string `json:"body"`
	Path string `json:"path"`
}

func postSave(w http.ResponseWriter, r *http.Request) {
	save := Save{}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &save)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	s, err := json.Marshal(save)
	if err != nil {
		fmt.Println("save payload unparseable", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	store.SetCache(save.Key, string(s))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"key": save.Key,
	})
}

func pollSave(w http.ResponseWriter, r *http.Request) {

	key := chi.URLParam(r, "key")
	res, err := store.GetCache(key)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if len(res) <= 10 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s := Save{}
	err = json.Unmarshal([]byte(res), &s)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}
