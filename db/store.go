package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/patrickmn/go-cache"
	"github.com/rs/xid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
)

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

func (s StoreData) SetCache(key string, value string) error {
	s.Cache.Set(key, value, cache.DefaultExpiration)
	return nil
}

func (s StoreData) DeleteCache(key string) error {
	s.Cache.Delete(key)
	return nil
}

func (s StoreData) GetCache(key string) (string, error) {
	value, found := s.Cache.Get(key)
	c, _ := value.(string)
	if !found || c == "" {
		return "", errors.New("not found")
	}
	return c, nil
}

func (s StoreData) SetLnCache(key string, value LnStore) error {
	s.Cache.Set(key, value, cache.DefaultExpiration)
	return nil
}

func (s StoreData) GetLnCache(key string) (LnStore, error) {
	value, found := s.Cache.Get(key)
	c, _ := value.(LnStore)
	if !found {
		return LnStore{}, errors.New("not found")
	}
	return c, nil
}

func (s StoreData) SetInvoiceCache(value []InvoiceStoreData) error {
	// The invoice should expire every 6 minutes
	s.Cache.Set(config.InvoiceList, value, 6*time.Minute)
	return nil
}

func (s StoreData) GetInvoiceCache() ([]InvoiceStoreData, error) {
	value, found := s.Cache.Get(config.InvoiceList)
	c, _ := value.([]InvoiceStoreData)
	if !found {
		return []InvoiceStoreData{}, errors.New("Invoice Cache not found")
	}
	return c, nil
}

func (s StoreData) SetBudgetInvoiceCache(value []BudgetStoreData) error {
	// The invoice should expire every 6 minutes
	s.Cache.Set(config.BudgetInvoiceList, value, 6*time.Minute)
	return nil
}

func (s StoreData) GetBudgetInvoiceCache() ([]BudgetStoreData, error) {
	value, found := s.Cache.Get(config.BudgetInvoiceList)
	c, _ := value.([]BudgetStoreData)
	if !found {
		return []BudgetStoreData{}, errors.New("Budget Invoice Cache not found")
	}
	return c, nil
}

func (s StoreData) SetSocketConnections(value Client) error {
	// The websocket in cache should not expire unless when deleted
	s.Cache.Set(value.Host, value, cache.NoExpiration)
	return nil
}

func (s StoreData) GetSocketConnections(host string) (Client, error) {
	value, found := s.Cache.Get(host)
	c, _ := value.(Client)
	if !found {
		return Client{}, errors.New("Socket Cache not found")
	}
	return c, nil
}

func (s StoreData) SetChallengeCache(key string, value string) error {
	// The challenge should expire every 10 minutes
	s.Cache.Set(key, value, 10*time.Minute)
	return nil
}

func (s StoreData) GetChallengeCache(key string) (string, error) {
	value, found := s.Cache.Get(key)
	c, _ := value.(string)
	if !found {
		return "", errors.New("Challenge Cache not found")
	}
	return c, nil
}

func Ask(w http.ResponseWriter, r *http.Request) {
	var m sync.Mutex
	m.Lock()

	ts := strconv.Itoa(int(time.Now().Unix()))
	challenge := xid.New().String()

	Store.SetChallengeCache(challenge, ts)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"challenge": challenge,
		"ts":        ts,
	})
	m.Unlock()
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
	TribeJWT              string                 `json:"tribe_jwt"`
}

func Verify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	challenge := chi.URLParam(r, "challenge")
	_, err := Store.GetChallengeCache(challenge)
	if err != nil {
		fmt.Println("challenge not found", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	payload := VerifyPayload{}
	body, err := io.ReadAll(r.Body)
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
	Store.SetChallengeCache(challenge, string(marshalled))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}

func Poll(w http.ResponseWriter, r *http.Request) {

	challenge := chi.URLParam(r, "challenge")
	res, err := Store.GetChallengeCache(challenge)
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

	tribeJWT, _ := auth.EncodeJwt(pld.Pubkey)
	pld.TribeJWT = tribeJWT

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
	body, err := io.ReadAll(r.Body)
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
