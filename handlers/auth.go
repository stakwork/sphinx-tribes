package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/form3tech-oss/jwt-go"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/utils"
)

type authHandler struct {
	db                        db.Database
	makeConnectionCodeRequest func() string
	decodeJwt                 func(token string) (jwt.MapClaims, error)
	encodeJwt                 func(pubkey string) (string, error)
}

func NewAuthHandler(db db.Database) *authHandler {
	return &authHandler{
		db:                        db,
		makeConnectionCodeRequest: MakeConnectionCodeRequest,
		decodeJwt:                 auth.DecodeJwt,
		encodeJwt:                 auth.EncodeJwt,
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
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Not a super admin: handler")
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Log in successful")
	}
}

func (ah *authHandler) CreateConnectionCode(w http.ResponseWriter, r *http.Request) {
	codeBody := db.InviteBody{}
	codeArr := []db.ConnectionCodes{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ReadAll Error", err)
	}
	r.Body.Close()

	err = json.Unmarshal(body, &codeBody)

	if err != nil {
		fmt.Println("Could not umarshal connection code body")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	for i := 0; i < int(codeBody.Number); i++ {
		code := ah.makeConnectionCodeRequest()

		if code != "" {
			newCode := db.ConnectionCodes{
				ConnectionString: code,
				IsUsed:           false,
			}
			codeArr = append(codeArr, newCode)
		}
	}

	_, err = ah.db.CreateConnectionCode(codeArr)

	if err != nil {
		fmt.Println("[auth] => ERR create connection code", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Codes created successfully")
}

func MakeConnectionCodeRequest() string {
	url := fmt.Sprintf("%s/invite", config.V2BotUrl)
	client := http.Client{}

	// Build v2 keysend payment data
	bodyData := utils.BuildV2ConnectionCodes(100, "new_user")
	jsonBody := []byte(bodyData)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("x-admin-token", config.V2BotToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		log.Printf("[Invite] Request Failed: %s", err)
		return ""
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Printf("Could not read invite body: %s", err)
	}

	inviteReponse := db.InviteReponse{}
	err = json.Unmarshal(body, &inviteReponse)

	if err != nil {
		fmt.Println("Could not get connection code")
		return ""
	}

	return inviteReponse.Invite
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

	exVerify, err := auth.VerifyDerSig(sig, k1, userKey)
	if err != nil || !exVerify {
		fmt.Println("[auth] Error signing signature")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	responseMsg := make(map[string]string)

	if userKey != "" {
		// Save in DB if the user does not exists already
		db.DB.CreateLnUser(userKey)

		// Set store data to true
		db.Store.SetLnCache(k1, db.LnStore{K1: k1, Key: userKey, Status: true})

		// Send socket message
		tokenString, err := auth.EncodeJwt(userKey)

		if err != nil {
			fmt.Println("[auth] error creating LNAUTH JWT")
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
			fmt.Println("[auth] Socket Error", err)
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
		fmt.Println("[auth] Failed to parse JWT")
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
			fmt.Println("[auth] error creating  refresh JWT")
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
	user["uuid"] = p.Uuid
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
