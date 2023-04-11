package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
)

func GetAdminPubkeys(w http.ResponseWriter, r *http.Request) {
	adminPubKeys := os.Getenv("ADMIN_PUBKEYS")
	admins := strings.Split(adminPubKeys, ",")
	type PubKeysReturn struct {
		Pubkeys []string `json:"pubkeys"`
	}
	pubkeys := PubKeysReturn{}
	if adminPubKeys != "" {
		for _, admin := range admins {
			pubkeys.Pubkeys = append(pubkeys.Pubkeys, admin)
		}
	}
	json.NewEncoder(w).Encode(pubkeys)
	w.WriteHeader(http.StatusOK)
}

func CreateConnectionCode(w http.ResponseWriter, r *http.Request) {
	code := db.ConnectionCodes{}
	now := time.Now()

	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(body, &code)

	code.IsUsed = false
	code.DateCreated = &now

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	_, err = db.DB.CreateConnectionCode(code)

	if err != nil {
		fmt.Println("=> ERR create connection code", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func GetConnectionCode(w http.ResponseWriter, _ *http.Request) {
	connectionCode := db.DB.GetConnectionCode()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(connectionCode)
}

func GetLnurlAuth(w http.ResponseWriter, _ *http.Request) {
	encodeData, err := auth.EncodeLNURL()
	responseData := make(map[string]string)

	if err != nil {
		responseData["k1"] = ""
		responseData["encode"] = ""

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Could not generate LNURL AUTH")
	}

	db.Store.SetLnCache(encodeData.K1, db.LnStore{encodeData.K1, "", false})

	responseData["k1"] = encodeData.K1
	responseData["encode"] = encodeData.Encode

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

func PollLnurlAuth(w http.ResponseWriter, r *http.Request) {
	k1 := r.URL.Query().Get("k1")
	responseData := make(map[string]interface{})

	res, err := db.Store.GetLnCache(k1)

	if err != nil {
		responseData["k1"] = ""
		responseData["status"] = false

		fmt.Println("=> ERR polling LNURL AUTH", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("LNURL auth data not found")
	}

	tokenString, err := auth.EncodeToken(res.Key)

	if err != nil {
		fmt.Println("error creating JWT")
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	person := db.DB.GetPersonByPubkey(res.Key)
	user := returnUserMap(person)

	responseData["k1"] = res.K1
	responseData["status"] = res.Status
	responseData["jwt"] = tokenString
	responseData["user"] = user

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

func ReceiveLnAuthData(w http.ResponseWriter, r *http.Request) {
	userKey := r.URL.Query().Get("key")
	k1 := r.URL.Query().Get("k1")

	responseMsg := make(map[string]string)

	if userKey != "" {
		// Save in DB if the user does not exists already
		db.DB.CreateLnUser(userKey)

		// Set store data to true
		db.Store.SetLnCache(k1, db.LnStore{k1, userKey, true})

		responseMsg["status"] = "OK"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseMsg)
	}

	responseMsg["status"] = "ERROR"
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(responseMsg)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-jwt")

	responseData := make(map[string]interface{})
	claims, err := auth.DecodeToken(token)

	if err != nil {
		fmt.Println("Failed to parse JWT")
		http.Error(w, http.StatusText(401), 401)
		return
	}

	pubkey := fmt.Sprint(claims["pubkey"])

	userCount := db.DB.GetLnUser(pubkey)

	if userCount > 0 {
		// Generate a new token
		tokenString, err := auth.EncodeToken(pubkey)

		if err != nil {
			fmt.Println("error creating  refresh JWT")
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		person := db.DB.GetPersonByPubkey(pubkey)
		user := returnUserMap(person)

		responseData["k1"] = ""
		responseData["status"] = true
		responseData["jwt"] = tokenString
		responseData["user"] = user

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseData)
	}
}

func returnUserMap(p db.Person) map[string]interface{} {
	user := make(map[string]interface{})

	user["id"] = p.ID
	user["created"] = p.Created
	user["owner_pubkey"] = p.OwnerPubKey
	user["owner_alias"] = p.OwnerAlias
	user["contact_key"] = p.OwnerContactKey
	user["img"] = p.Img
	user["description"] = p.Description
	user["tags"] = p.Tags
	user["unique_name"] = p.UniqueName
	user["pubkey"] = p.OwnerPubKey
	user["extras"] = p.Extras
	user["last_login"] = p.LastLogin
	user["price_to_meet"] = p.PriceToMeet
	user["alias"] = p.OwnerAlias
	user["url"] = config.Host

	return user
}
