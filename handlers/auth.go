package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
)

type authHandler struct {
	db        db.Database
	decodeJwt func(token string) (jwt.MapClaims, error)
	encodeJwt func(pubkey string) (string, error)
}

func NewAuthHandler(db db.Database) *authHandler {
	return &authHandler{
		db:        db,
		decodeJwt: auth.DecodeJwt,
		encodeJwt: auth.EncodeJwt,
	}
}

func GetAdminPubkeys(w http.ResponseWriter, r *http.Request) {
	type PubKeysReturn struct {
		Pubkeys []string `json:"pubkeys"`
	}
	pubkeys := PubKeysReturn{
		Pubkeys: config.SuperAdmins,
	}
	json.NewEncoder(w).Encode(pubkeys)
	w.WriteHeader(http.StatusOK)
}

func (ah *authHandler) GetIsAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	isAdmin := auth.AdminCheck(pubKeyFromAuth)

	if !auth.IsFreePass() && !isAdmin {
		fmt.Println("Not a super admin: handler")
		http.Error(w, http.StatusText(401), 401)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Log in successful")
	}
}

func (ah *authHandler) CreateConnectionCode(w http.ResponseWriter, r *http.Request) {
	code := db.ConnectionCodes{}
	now := time.Now()

	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(body, &code)

	code.IsUsed = false
	code.DateCreated = &now

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	_, err = ah.db.CreateConnectionCode(code)

	if err != nil {
		fmt.Println("=> ERR create connection code", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (ah *authHandler) GetConnectionCode(w http.ResponseWriter, _ *http.Request) {
	connectionCode := ah.db.GetConnectionCode()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(connectionCode)
}

func GetLnurlAuth(w http.ResponseWriter, r *http.Request) {
	socketKey := r.URL.Query().Get("socketKey")
	socket, _ := db.Store.GetSocketConnections(socketKey)
	serverHost := r.Host

	encodeData, err := auth.EncodeLNURL(serverHost)
	responseData := make(map[string]string)

	if err != nil {
		responseData["k1"] = ""
		responseData["encode"] = ""

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Could not generate LNURL AUTH")
	}

	db.Store.SetLnCache(encodeData.K1, db.LnStore{K1: encodeData.K1, Key: "", Status: false})

	// add socket to store with K1, so the LNURL return data can use it
	db.Store.SetSocketConnections(db.Client{
		Host: encodeData.K1[0:20],
		Conn: socket.Conn,
	})

	responseData["k1"] = encodeData.K1
	responseData["encode"] = encodeData.Encode

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

func ReceiveLnAuthData(w http.ResponseWriter, r *http.Request) {
	userKey := r.URL.Query().Get("key")
	k1 := r.URL.Query().Get("k1")
	sig := r.URL.Query().Get("sig")

	fmt.Println("LN Sig ====", sig)

	urlSig := base64.URLEncoding.EncodeToString([]byte(sig))

	// abVerify, _ := auth.VerifyArbitrary(userKey, sig)
	exVerify := auth.VerifySig(sig, k1)

	// fmt.Println("AB Verify ====", abVerify)
	fmt.Println("Ex Verify ====", urlSig, exVerify)

	responseMsg := make(map[string]string)

	if userKey != "" {
		// Save in DB if the user does not exists already
		db.DB.CreateLnUser(userKey)

		// Set store data to true
		db.Store.SetLnCache(k1, db.LnStore{K1: k1, Key: userKey, Status: true})

		// Send socket message
		tokenString, err := auth.EncodeJwt(userKey)

		if err != nil {
			fmt.Println("error creating LNAUTH JWT")
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		person := db.DB.GetPersonByPubkey(userKey)
		user := returnUserMap(person)

		socketMsg := make(map[string]interface{})

		// Send socket message
		socketMsg["k1"] = k1
		socketMsg["status"] = true
		socketMsg["jwt"] = tokenString
		socketMsg["user"] = user
		socketMsg["msg"] = "lnauth_success"

		socket, err := db.Store.GetSocketConnections(k1[0:20])

		if err == nil {
			socket.Conn.WriteJSON(socketMsg)
			db.Store.DeleteCache(k1[0:20])
		} else {
			fmt.Println("Socket Error", err)
		}

		responseMsg["status"] = "OK"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseMsg)
	}

	responseMsg["status"] = "ERROR"
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(responseMsg)
}

func (ah *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-jwt")

	responseData := make(map[string]interface{})
	claims, err := ah.decodeJwt(token)

	if err != nil {
		fmt.Println("Failed to parse JWT")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	pubkey := fmt.Sprint(claims["pubkey"])

	userCount := ah.db.GetLnUser(pubkey)

	if userCount > 0 {
		// Generate a new token
		tokenString, err := ah.encodeJwt(pubkey)

		if err != nil {
			fmt.Println("error creating  refresh JWT")
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		person := ah.db.GetPersonByPubkey(pubkey)
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
