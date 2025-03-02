package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/form3tech-oss/jwt-go"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
	"github.com/stakwork/sphinx-tribes/utils"
)

// AuthHandler struct
type AuthHandler struct {
	db                        db.Database
	makeConnectionCodeRequest func(inviter_pubkey string, inviter_route_hint string, msats_amount uint64) string
	decodeJwt                 func(token string) (jwt.MapClaims, error)
	encodeJwt                 func(pubkey string) (string, error)
}

func NewAuthHandler(db db.Database) *AuthHandler {
	return &AuthHandler{
		db:                        db,
		makeConnectionCodeRequest: MakeConnectionCodeRequest,
		decodeJwt:                 auth.DecodeJwt,
		encodeJwt:                 auth.EncodeJwt,
	}
}

type LnAuthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	JWT     string `json:"jwt,omitempty"`
}

type RefreshTokenResponse struct {
	K1     string    `json:"k1,omitempty"`
	Status bool      `json:"status"`
	JWT    string    `json:"jwt"`
	User   db.Person `json:"user"`
}

type ConnectionCodesListResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Codes []db.ConnectionCodesList `json:"codes"`
		Total int64                    `json:"total"`
	} `json:"data"`
}

// GetAdminPubkeys godoc
//
//	@Summary		Get admin pubkeys
//	@Description	Get a list of admin pubkeys
//	@Tags			Auth
//	@Success		200	{array}	string
//	@Router			/admin_pubkeys [get]
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

// GetIsAdmin godoc
//
//	@Summary		Check if user is admin
//	@Description	Check if the user is an admin. Requires a valid JWT token in the request context.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{string}	string	"Returns a success message if the user is an admin"
//	@Failure		401	{string}	string	"Unauthorized: User is not an admin or missing/invalid JWT token"
//	@Router			/admin/auth [get]
func (ah *AuthHandler) GetIsAdmin(w http.ResponseWriter, r *http.Request) {
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

// CreateConnectionCode godoc
//
//	@Summary		Create connection codes
//	@Description	Create one or more connection codes for a user. Requires a valid pubkey and route hint.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body	db.InviteBody	true	"Request body containing pubkey, route hint, sats amount, and number of codes"
//	@Security		SuperAdminAuth
//	@Success		200	{string}	string	"Connection codes created successfully"
//	@Failure		400	{string}	string	"Bad request: Missing or invalid parameters"
//	@Failure		406	{string}	string	"Not acceptable: Invalid request body"
//	@Router			/connectioncodes [post]
func (ah *AuthHandler) CreateConnectionCode(w http.ResponseWriter, r *http.Request) {
	codeBody := db.InviteBody{}
	codeArr := []db.ConnectionCodes{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("ReadAll Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()

	err = json.Unmarshal(body, &codeBody)

	if err != nil {
		logger.Log.Error("Could not unmarshal connection code body")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if codeBody.Pubkey != "" && codeBody.RouteHint == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Route hint is required when pubkey is provided")
		return
	}

	if codeBody.RouteHint != "" && codeBody.Pubkey == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode("pubkey is required when Route hint is provided")
		return
	}

	if codeBody.SatsAmount == 0 {
		codeBody.SatsAmount = 100
	} else {
		codeBody.SatsAmount = utils.ConvertSatsToMsats(codeBody.SatsAmount)
	}

	for i := 0; i < int(codeBody.Number); i++ {
		code := ah.makeConnectionCodeRequest(codeBody.Pubkey, codeBody.RouteHint, codeBody.SatsAmount)

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
		logger.Log.Error("[auth] => ERR create connection code: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Codes created successfully")
}

func MakeConnectionCodeRequest(inviter_pubkey string, inviter_route_hint string, msats_amount uint64) string {
	url := fmt.Sprintf("%s/invite", config.V2BotUrl)
	client := http.Client{}

	// Build v2 keysend payment data
	bodyData := utils.BuildV2ConnectionCodes(msats_amount, "new_user", inviter_pubkey, inviter_route_hint)
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
		logger.Log.Error("Could not get connection code")
		return ""
	}

	return inviteReponse.Invite
}

// GetConnectionCode godoc
//
//	@Summary		Get a connection code
//	@Description	Retrieve a single connection code from the database.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	db.ConnectionCodesShort	"Connection code retrieved successfully"
//	@Router			/connectioncodes [get]
func (ah *AuthHandler) GetConnectionCode(w http.ResponseWriter, _ *http.Request) {
	connectionCode := ah.db.GetConnectionCode()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(connectionCode)
}

// GetLnurlAuth godoc
//
//	@Summary		Generate LNURL-auth data
//	@Description	Generate a unique LNURL-auth challenge (k1) and encoded LNURL for authentication.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			socketKey	query		string				true	"Socket connection key to associate with the LNURL-auth challenge"
//	@Success		200			{object}	auth.LnEncodeData	"LNURL-auth data generated successfully"
//	@Failure		400			{object}	auth.LnEncodeData	"Bad request: Failed to generate LNURL-auth data"
//	@Router			/lnurl_auth [get]
func GetLnurlAuth(w http.ResponseWriter, r *http.Request) {
	socketKey := r.URL.Query().Get("socketKey")
	socket, _ := db.Store.GetSocketConnections(socketKey)
	serverHost := r.Host

	encodeData, err := auth.EncodeLNURLFunc(serverHost)
	responseData := make(map[string]string)

	if err != nil {
		responseData["k1"] = ""
		responseData["encode"] = ""

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseData)
		return
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

// ReceiveLnAuthData godoc
//
//	@Summary		Receive LNURL auth data
//	@Description	Receive LNURL auth data and authenticate the user using LNURL-auth protocol.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			key	query		string			true	"User's public key"
//	@Param			k1	query		string			true	"Unique challenge string (k1)"
//	@Param			sig	query		string			true	"Signature of the challenge signed by the user's private key"
//	@Success		200	{object}	LnAuthResponse	"Authentication successful"
//	@Failure		400	{object}	LnAuthResponse	"Invalid request or missing parameters"
//	@Failure		401	{object}	LnAuthResponse	"Unauthorized: Signature verification failed"
//	@Failure		406	{object}	LnAuthResponse	"Not Acceptable: JWT creation failed"
//	@Router			/lnauth_login [get]
func ReceiveLnAuthData(w http.ResponseWriter, r *http.Request) {
	userKey := r.URL.Query().Get("key")
	k1 := r.URL.Query().Get("k1")
	sig := r.URL.Query().Get("sig")

	exVerify, err := auth.VerifyDerSig(sig, k1, userKey)
	if err != nil || !exVerify {
		logger.Log.Error("[auth] Error signing signature")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	responseMsg := LnAuthResponse{}

	if userKey != "" {
		// Save in DB if the user does not exists already
		db.DB.CreateLnUser(userKey)

		// Set store data to true
		db.Store.SetLnCache(k1, db.LnStore{K1: k1, Key: userKey, Status: true})

		// Send socket message
		tokenString, err := auth.EncodeJwt(userKey)

		if err != nil {
			logger.Log.Error("[auth] error creating LNAUTH JWT")
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
			logger.Log.Error("[auth] Socket Error: %v", err)
		}

		responseMsg.Status = "OK"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseMsg)
	}

	responseMsg.Status = "ERROR"
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(responseMsg)
}

// RefreshToken godoc
//
//	@Summary		Refresh JWT token
//	@Description	Refresh the JWT token using a valid existing token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			x-jwt	header		string					true	"Existing JWT token"
//	@Success		200		{object}	RefreshTokenResponse	"Token refreshed successfully"
//	@Failure		401		{object}	string					"Unauthorized: Missing or invalid JWT token"
//	@Failure		406		{object}	string					"Not Acceptable: Failed to create a new JWT token"
//	@Router			/refresh_jwt [get]
func (ah *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-jwt")

	if token == "" {
		logger.Log.Error("[auth] Missing JWT token")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Missing JWT token")
		return
	}

	responseData := make(map[string]interface{})
	claims, err := ah.decodeJwt(token)

	if err != nil {
		logger.Log.Error("[auth] Failed to parse JWT", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	pubkey, ok := claims["pubkey"].(string)
	if !ok || pubkey == "" {
		logger.Log.Error("[auth] Missing pubkey claim in JWT")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Missing pubkey claim in JWT")
		return
	}

	userCount := ah.db.GetLnUser(pubkey)

	if userCount > 0 {
		// Generate a new token
		tokenString, err := ah.encodeJwt(pubkey)

		if err != nil {
			logger.Log.Error("[auth] error creating refresh JWT")
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

// ListConnectionCodes godoc
//
//	@Summary		List connection codes
//	@Description	List all connection codes with pagination support.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		SuperAdminAuth
//	@Param			page	query		int							false	"Page number (default: 1)"
//	@Param			limit	query		int							false	"Number of items per page (default: 20)"
//	@Success		200		{object}	ConnectionCodesListResponse	"Connection codes retrieved successfully"
//	@Failure		400		{object}	ConnectionCodesListResponse	"Bad request: Invalid pagination parameters"
//	@Failure		500		{object}	ConnectionCodesListResponse	"Internal server error: Failed to retrieve connection codes"
//	@Router			/connectioncodes/list [get]
func (ah *AuthHandler) ListConnectionCodes(w http.ResponseWriter, r *http.Request) {

	page := 1
	limit := 20

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	codes, total, err := ah.db.GetConnectionCodesList(page, limit)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ConnectionCodesListResponse{
			Success: false,
			Data: struct {
				Codes []db.ConnectionCodesList `json:"codes"`
				Total int64                    `json:"total"`
			}{
				Codes: []db.ConnectionCodesList{},
				Total: 0,
			},
		})
		return
	}

	response := ConnectionCodesListResponse{
		Success: true,
		Data: struct {
			Codes []db.ConnectionCodesList `json:"codes"`
			Total int64                    `json:"total"`
		}{
			Codes: codes,
			Total: total,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
