package db

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
	"github.com/stakwork/sphinx-tribes/auth"
	"golang.org/x/crypto/blake2b"
)

// Store struct
type StoreData struct {
	Cache *cache.Cache
}

type LnStore struct {
	K1     string
	Key    string
	Status bool
}

var Store StoreData

func InitCache() {
	authTimeout := 120
	Store = StoreData{
		Cache: cache.New(
			time.Duration(authTimeout)*time.Second,
			time.Duration(authTimeout*3)*time.Second,
		),
	}
}

// SetCache
func (s StoreData) SetCache(key string, value string) error {
	s.Cache.Set(key, value, cache.DefaultExpiration)
	return nil
}

// DeleteCache
func (s StoreData) DeleteCache(key string) error {
	s.Cache.Delete(key)
	return nil
}

// GetCache
func (s StoreData) GetCache(key string) (string, error) {
	value, found := s.Cache.Get(key)
	c, _ := value.(string)
	if !found || c == "" {
		return "", errors.New("not found")
	}
	return c, nil
}

// SetCache
func (s StoreData) SetLnCache(key string, value LnStore) error {
	s.Cache.Set(key, value, cache.DefaultExpiration)
	return nil
}

// GetCache
func (s StoreData) GetLnCache(key string) (LnStore, error) {
	value, found := s.Cache.Get(key)
	c, _ := value.(LnStore)
	if !found {
		return LnStore{}, errors.New("not found")
	}
	return c, nil
}

func Ask(w http.ResponseWriter, r *http.Request) {
	ts := strconv.Itoa(int(time.Now().Unix()))
	h := blake2b.Sum256([]byte(ts))
	challenge := base64.URLEncoding.EncodeToString(h[:])

	Store.SetCache(challenge, ts)

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

func Verify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	challenge := chi.URLParam(r, "challenge")
	_, err := Store.GetCache(challenge)
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
	Store.SetCache(challenge, string(marshalled))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}

/*
curl localhost:5002/ask
curl localhost:5002/poll/d5SYZNY5pQ7dXwHP-oXh2uSOPUEX0fUJOXI0_5-eOsg=
*/
func Poll(w http.ResponseWriter, r *http.Request) {

	challenge := chi.URLParam(r, "challenge")
	res, err := Store.GetCache(challenge)
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

	existing := DB.GetPersonByPubkey(pld.Pubkey)
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

	// update LastLogin for user
	DB.UpdatePerson(pld.ID, map[string]interface{}{
		"last_login": time.Now().Unix(),
	})

	// store.DeleteChallenge(challenge)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pld)
}

type Save struct {
	Key    string `json:"key"`
	Body   string `json:"body"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

func PostSave(w http.ResponseWriter, r *http.Request) {

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

	Store.SetCache(save.Key, string(s))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"key": save.Key,
	})
}

func PollSave(w http.ResponseWriter, r *http.Request) {

	key := chi.URLParam(r, "key")
	res, err := Store.GetCache(key)
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
