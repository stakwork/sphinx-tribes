package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
	"github.com/stakwork/sphinx-tribes/utils"
)

type tribeHandler struct {
	db                      db.Database
	verifyTribeUUID         func(uuid string, checkTimestamp bool) (string, error)
	tribeUniqueNameFromName func(name string) (string, error)
}

func NewTribeHandler(db db.Database) *tribeHandler {
	return &tribeHandler{
		db:                      db,
		verifyTribeUUID:         auth.VerifyTribeUUID,
		tribeUniqueNameFromName: TribeUniqueNameFromName,
	}
}

// GetAllTribes godoc
//
//	@Summary		Get all tribes
//	@Description	Get a list of all tribes
//	@Tags			Tribes
//	@Success		200	{array}	db.Tribe
//	@Router			/tribes [get]
func (th *tribeHandler) GetAllTribes(w http.ResponseWriter, r *http.Request) {
	tribes := th.db.GetAllTribes()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribes)
}

// GetTotalTribes godoc
//
//	@Summary		Get total number of tribes
//	@Description	Get the total number of tribes
//	@Tags			Tribes
//	@Success		200	{object}	int
//	@Router			/tribes/total [get]
func (th *tribeHandler) GetTotalTribes(w http.ResponseWriter, r *http.Request) {
	tribesTotal := th.db.GetTribesTotal()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribesTotal)
}

// GetListedTribes godoc
//
//	@Summary		Get listed tribes
//	@Description	Get a list of listed tribes
//	@Tags			Tribes
//	@Success		200	{array}	db.Tribe
//	@Router			/tribes [get]
func (th *tribeHandler) GetListedTribes(w http.ResponseWriter, r *http.Request) {
	tribes := th.db.GetListedTribes(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribes)
}

// GetTribesByOwner godoc
//
//	@Summary		Get tribes by owner
//	@Description	Get a list of tribes by owner public key
//	@Tags			Tribes
//	@Param			pubkey	path	string	true	"Owner public key"
//	@Param			all		query	string	false	"Include all tribes"
//	@Success		200		{array}	db.Tribe
//	@Router			/tribes_by_owner/{pubkey} [get]
func (th *tribeHandler) GetTribesByOwner(w http.ResponseWriter, r *http.Request) {
	all := r.URL.Query().Get("all")
	tribes := []db.Tribe{}
	pubkey := chi.URLParam(r, "pubkey")
	if all == "true" {
		tribes = th.db.GetAllTribesByOwner(pubkey)
	} else {
		tribes = th.db.GetTribesByOwner(pubkey)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribes)
}

// GetTribesByAppUrl godoc
//
//	@Summary		Get tribes by app URL
//	@Description	Get a list of tribes by app URL
//	@Tags			Tribes
//	@Param			app_url	path	string	true	"App URL"
//	@Success		200		{array}	db.Tribe
//	@Router			/tribes/app_url/{app_url} [get]
func (th *tribeHandler) GetTribesByAppUrl(w http.ResponseWriter, r *http.Request) {
	tribes := []db.Tribe{}
	app_url := chi.URLParam(r, "app_url")
	tribes = th.db.GetTribesByAppUrl(app_url)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribes)
}

// GetTribesByAppUrls godoc
//
//	@Summary		Get tribes by multiple app URLs
//	@Description	Get a list of tribes by multiple app URLs
//	@Tags			Tribes
//	@Param			app_urls	path		string	true	"Comma-separated list of app URLs"
//	@Success		200			{object}	map[string][]db.Tribe
//	@Router			/tribes/apps_urls/{app_urls} [get]
func GetTribesByAppUrls(w http.ResponseWriter, r *http.Request) {
	app_urls := chi.URLParam(r, "app_urls")
	app_url_list := strings.Split(app_urls, ",")
	m := make(map[string][]db.Tribe)
	for _, app_url := range app_url_list {
		tribes := db.DB.GetTribesByAppUrl(app_url)
		m[app_url] = tribes
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m)
}

// PutTribeStats godoc
//
//	@Summary		Update tribe stats
//	@Description	Update the stats of a tribe
//	@Tags			Tribes
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			tribe	body		db.Tribe	true	"Tribe object"
//	@Success		200		{object}	bool
//	@Router			/tribestats [put]
func PutTribeStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	tribe := db.Tribe{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &tribe)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if tribe.UUID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := auth.VerifyTribeUUID(tribe.UUID, false)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// from token must match
	if pubKeyFromAuth != extractedPubkey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	now := time.Now()
	tribe.Updated = &now
	db.DB.UpdateTribe(tribe.UUID, map[string]interface{}{
		"member_count": tribe.MemberCount,
		"updated":      &now,
		"bots":         tribe.Bots,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

// DeleteTribe godoc
//
//	@Summary		Delete a tribe
//	@Description	Delete a tribe by UUID
//	@Tags			Tribes
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path		string	true	"Tribe UUID"
//	@Success		200		{object}	bool
//	@Router			/tribe/{uuid} [delete]
func (th *tribeHandler) DeleteTribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	uuid := chi.URLParam(r, "uuid")

	if uuid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := th.verifyTribeUUID(uuid, false)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// from token must match
	if pubKeyFromAuth != extractedPubkey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	th.db.UpdateTribe(uuid, map[string]interface{}{
		"deleted": true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

// GetTribe godoc
//
//	@Summary		Get a tribe
//	@Description	Get a tribe by UUID
//	@Tags			Tribes
//	@Param			uuid	path		string	true	"Tribe UUID"
//	@Success		200		{object}	map[string]interface{}
//	@Router			/tribes/{uuid} [get]
func (th *tribeHandler) GetTribe(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	tribe := th.db.GetTribe(uuid)

	var theTribe map[string]interface{}
	j, _ := json.Marshal(tribe)
	json.Unmarshal(j, &theTribe)

	theTribe["channels"] = th.db.GetChannelsByTribe(uuid)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(theTribe)
}

// GetFirstTribeByFeed godoc
//
//	@Summary		Get the first tribe by feed URL
//	@Description	Get the first tribe by feed URL
//	@Tags			Tribes
//	@Param			url	query		string	true	"Feed URL"
//	@Success		200	{object}	map[string]interface{}
//	@Router			/tribe_by_feed [get]
func (th *tribeHandler) GetFirstTribeByFeed(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	tribe := th.db.GetFirstTribeByFeedURL(url)

	if tribe.UUID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var theTribe map[string]interface{}
	j, _ := json.Marshal(tribe)
	json.Unmarshal(j, &theTribe)

	theTribe["channels"] = th.db.GetChannelsByTribe(tribe.UUID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(theTribe)
}

// GetTribeByUniqueName godoc
//
//	@Summary		Get a tribe by unique name
//	@Description	Get a tribe by unique name
//	@Tags			Tribes
//	@Param			un	path		string	true	"Unique name"
//	@Success		200	{object}	map[string]interface{}
//	@Router			/tribe_by_un/{un} [get]
func (th *tribeHandler) GetTribeByUniqueName(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "un")
	tribe := th.db.GetTribeByUniqueName(uuid)

	var theTribe map[string]interface{}
	j, _ := json.Marshal(tribe)
	json.Unmarshal(j, &theTribe)

	theTribe["channels"] = th.db.GetChannelsByTribe(tribe.UUID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(theTribe)
}

// CreateOrEditTribe godoc
//
//	@Summary		Create or edit a tribe
//	@Description	Create or edit a tribe
//	@Tags			Tribes
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			tribe	body		db.Tribe	true	"Tribe object"
//	@Success		200		{object}	db.Tribe
//	@Router			/tribe [post]
func (th *tribeHandler) CreateOrEditTribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	tribe := db.Tribe{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &tribe)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if tribe.UUID == "" {
		logger.Log.Info("createOrEditTribe no uuid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	now := time.Now() //.Format(time.RFC3339)

	extractedPubkey, err := th.verifyTribeUUID(tribe.UUID, false)
	if err != nil {
		logger.Log.Error("extract UUID error: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if pubKeyFromAuth == "" {
		tribe.Created = &now
	} else { // IF PUBKEY IN CONTEXT, MUST AUTH!
		if pubKeyFromAuth != extractedPubkey {
			logger.Log.Info("createOrEditTribe pubkeys dont match")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	existing := th.db.GetTribe(tribe.UUID)
	if existing.UUID == "" { // if doesn't exist already, create unique name
		tribe.UniqueName, _ = th.tribeUniqueNameFromName(tribe.Name)
	} else { // already exists! make sure it's owned
		if existing.OwnerPubKey != extractedPubkey {
			logger.Log.Info("createOrEditTribe tribe.ownerPubKey not match")
			logger.Log.Info("existing owner: %s", existing.OwnerPubKey)
			logger.Log.Info("extracted pubkey: %s", extractedPubkey)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	tribe.OwnerPubKey = extractedPubkey
	tribe.Updated = &now
	tribe.LastActive = now.Unix()

	_, err = th.db.CreateOrEditTribe(tribe)
	if err != nil {
		logger.Log.Error("=> ERR createOrEditTribe: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribe)
}

// PutTribeActivity godoc
//
//	@Summary		Update tribe activity
//	@Description	Update the activity of a tribe
//	@Tags			Tribes
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			uuid	path		string	true	"Tribe UUID"
//	@Success		200		{object}	bool
//	@Router			/tribeactivity/{uuid} [put]
func PutTribeActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := auth.VerifyTribeUUID(uuid, false)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// from token must match
	if pubKeyFromAuth != extractedPubkey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	now := time.Now().Unix()
	db.DB.UpdateTribe(uuid, map[string]interface{}{
		"last_active": now,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

// SetTribePreview godoc
//
//	@Summary		Set tribe preview
//	@Description	Set the preview of a tribe
//	@Tags			Tribes
//	@Param			uuid	path		string	true	"Tribe UUID"
//	@Param			preview	query		string	true	"Preview URL"
//	@Success		200		{object}	bool
//	@Router			/tribepreview/{uuid} [put]
func (th *tribeHandler) SetTribePreview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := th.verifyTribeUUID(uuid, false)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// from token must match
	if pubKeyFromAuth != extractedPubkey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	preview := r.URL.Query().Get("preview")
	th.db.UpdateTribe(uuid, map[string]interface{}{
		"preview": preview,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

// CreateLeaderBoard godoc
//
//	@Summary		Create a leaderboard
//	@Description	Create a leaderboard for a tribe
//	@Tags			Tribes
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			tribe_uuid	path		string				true	"Tribe UUID"
//	@Param			leaderboard	body		[]db.LeaderBoard	true	"Leaderboard object"
//	@Success		200			{object}	bool
//	@Router			/leaderboard/{tribe_uuid} [post]
func CreateLeaderBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "tribe_uuid")

	leaderBoard := []db.LeaderBoard{}

	if uuid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := auth.VerifyTribeUUID(uuid, false)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//from token must match
	if pubKeyFromAuth != extractedPubkey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &leaderBoard)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	_, err = db.DB.CreateLeaderBoard(uuid, leaderBoard)

	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

// GetLeaderBoard godoc
//
//	@Summary		Get a leaderboard
//	@Description	Get a leaderboard for a tribe
//	@Tags			Tribes
//	@Param			tribe_uuid	path		string	true	"Tribe UUID"
//	@Param			alias		query		string	false	"Alias"
//	@Success		200			{object}	interface{}
//	@Router			/leaderboard/{tribe_uuid} [get]
func GetLeaderBoard(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "tribe_uuid")
	alias := r.URL.Query().Get("alias")

	if alias == "" {
		leaderBoards := db.DB.GetLeaderBoard(uuid)

		var board = []db.LeaderBoard{}
		for _, leaderboard := range leaderBoards {
			leaderboard.TribeUuid = ""
			board = append(board, leaderboard)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(board)
	} else {
		leaderBoardFromDb := db.DB.GetLeaderBoardByUuidAndAlias(uuid, alias)

		if leaderBoardFromDb.Alias != alias {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(leaderBoardFromDb)
	}
}

// UpdateLeaderBoard godoc
//
//	@Summary		Update a leaderboard
//	@Description	Update a leaderboard for a tribe
//	@Tags			Tribes
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			tribe_uuid	path		string			true	"Tribe UUID"
//	@Param			leaderboard	body		db.LeaderBoard	true	"Leaderboard object"
//	@Success		200			{object}	bool
//	@Router			/leaderboard/{tribe_uuid} [put]
func UpdateLeaderBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	uuid := chi.URLParam(r, "tribe_uuid")

	if uuid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := auth.VerifyTribeUUID(uuid, false)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//from token must match
	if pubKeyFromAuth != extractedPubkey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	leaderBoard := db.LeaderBoard{}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &leaderBoard)
	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	leaderBoardFromDb := db.DB.GetLeaderBoardByUuidAndAlias(uuid, leaderBoard.Alias)

	if leaderBoardFromDb.Alias != leaderBoard.Alias {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	leaderBoard.TribeUuid = leaderBoardFromDb.TribeUuid

	db.DB.UpdateLeaderBoard(leaderBoardFromDb.TribeUuid, leaderBoardFromDb.Alias, map[string]interface{}{
		"spent":      leaderBoard.Spent,
		"earned":     leaderBoard.Earned,
		"reputation": leaderBoard.Reputation,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

// GenerateInvoice godoc
//
//	@Summary		Generate an invoice
//	@Description	Generate an invoice for a tribe
//	@Tags			Tribes
//	@Param			invoice	body		db.InvoiceRequest	true	"Invoice request"
//	@Success		200		{object}	db.InvoiceResponse
//	@Router			/invoice [post]
func GenerateInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceRes, invoiceErr := db.InvoiceResponse{}, db.InvoiceError{}

	if config.IsV2Payment {
		invoiceRes, invoiceErr = GenerateV2Invoice(w, r)
	} else {
		invoiceRes, invoiceErr = GenerateV1Invoice(w, r)
	}

	if invoiceErr.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(invoiceErr)
	}

	invoice := db.InvoiceRequest{}
	body, err := io.ReadAll(r.Body)

	r.Body.Close()

	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = json.Unmarshal(body, &invoice)

	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	pub_key := invoice.User_pubkey
	owner_key := invoice.Owner_pubkey
	date, _ := utils.ConvertStringToInt(invoice.Created)
	invoiceType := invoice.Type
	routeHint := invoice.Route_hint
	amount, _ := utils.ConvertStringToUint(invoice.Amount)

	paymentRequest := invoiceRes.Response.Invoice
	now := time.Now()

	newInvoice := db.NewInvoiceList{
		PaymentRequest: paymentRequest,
		Type:           db.InvoiceType(invoiceType),
		OwnerPubkey:    owner_key,
		Created:        &now,
		Updated:        &now,
		Status:         false,
	}

	newInvoiceData := db.UserInvoiceData{
		PaymentRequest: paymentRequest,
		Created:        date,
		Amount:         amount,
		UserPubkey:     pub_key,
		RouteHint:      routeHint,
	}

	db.DB.ProcessAddInvoice(newInvoice, newInvoiceData)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceRes)
}

func GenerateV1Invoice(w http.ResponseWriter, r *http.Request) (db.InvoiceResponse, db.InvoiceError) {
	invoice := db.InvoiceRequest{}
	body, err := io.ReadAll(r.Body)

	r.Body.Close()

	if err != nil {
		logger.Log.Error("%v", err)
		return db.InvoiceResponse{}, db.InvoiceError{Success: false, Error: err.Error()}
	}

	err = json.Unmarshal(body, &invoice)

	if err != nil {
		logger.Log.Error("%v", err)
		return db.InvoiceResponse{}, db.InvoiceError{Success: false, Error: err.Error()}
	}

	memo := invoice.Memo
	amount, _ := utils.ConvertStringToUint(invoice.Amount)

	url := fmt.Sprintf("%s/invoices", config.RelayUrl)

	bodyData := fmt.Sprintf(`{"amount": %d, "memo": "%s"}`, amount, memo)

	jsonBody := []byte(bodyData)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))

	req.Header.Set("x-user-token", config.RelayAuthKey)
	req.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(req)

	if err != nil {
		log.Printf("Request Failed: %s", err)
		return db.InvoiceResponse{}, db.InvoiceError{Success: false, Error: err.Error()}
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)

	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return db.InvoiceResponse{}, db.InvoiceError{Success: false, Error: err.Error()}
	}

	// Unmarshal result
	invoiceRes := db.InvoiceResponse{}

	err = json.Unmarshal(body, &invoiceRes)

	if err != nil {
		log.Printf("Unmarshal body failed: %s", err)
		return db.InvoiceResponse{}, db.InvoiceError{Success: false, Error: err.Error()}
	}

	return invoiceRes, db.InvoiceError{Success: true}
}

// GenerateInvoice godoc
//
//	@Summary		Generate an invoice
//	@Description	Generate an invoice for a tribe
//	@Tags			Tribes
//	@Param			invoice	body		db.InvoiceRequest	true	"Invoice request"
//	@Success		200		{object}	db.InvoiceResponse
//	@Router			/tribes/invoice [post]
func GenerateV2Invoice(w http.ResponseWriter, r *http.Request) (db.InvoiceResponse, db.InvoiceError) {
	invoice := db.InvoiceRequest{}

	var err error
	body, err := io.ReadAll(r.Body)

	r.Body.Close()

	if err != nil {
		logger.Log.Error("%v", err)
		return db.InvoiceResponse{}, db.InvoiceError{Success: false, Error: err.Error()}
	}

	err = json.Unmarshal(body, &invoice)

	if err != nil {
		logger.Log.Error("%v", err)
		return db.InvoiceResponse{}, db.InvoiceError{Success: false, Error: err.Error()}
	}

	url := fmt.Sprintf("%s/invoice", config.V2BotUrl)

	amount, _ := utils.ConvertStringToUint(invoice.Amount)

	amountMsat := amount * 1000

	bodyData := fmt.Sprintf(`{"amt_msat": %d}`, amountMsat)

	jsonBody := []byte(bodyData)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))

	req.Header.Set("x-admin-token", config.V2BotToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		log.Printf("Client Request Failed: %s", err)
		return db.InvoiceResponse{}, db.InvoiceError{Success: false, Error: err.Error()}
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)

	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return db.InvoiceResponse{}, db.InvoiceError{Success: false, Error: err.Error()}
	}

	// Unmarshal result
	v2InvoiceRes := db.V2CreateInvoiceResponse{}
	err = json.Unmarshal(body, &v2InvoiceRes)

	if err != nil {
		log.Printf("Json Unmarshal failed: %s", err)
		return db.InvoiceResponse{}, db.InvoiceError{Success: false, Error: err.Error()}
	}
	return db.InvoiceResponse{
		Response: db.Invoice{
			Invoice: v2InvoiceRes.Bolt11,
		},
	}, db.InvoiceError{Success: true}
}

// GenerateBudgetInvoice godoc
//
//	@Summary		Generate a budget invoice
//	@Description	Generate a budget invoice for a tribe
//	@Tags			Tribes
//	@Param			invoice	body		db.BudgetInvoiceRequest	true	"Budget invoice request"
//	@Success		200		{object}	db.InvoiceResponse
//	@Router			/tribes/budget_invoice [post]
func (th *tribeHandler) GenerateBudgetInvoice(w http.ResponseWriter, r *http.Request) {
	if config.IsV2Payment {
		th.GenerateV2BudgetInvoice(w, r)
	} else {
		th.GenerateV1BudgetInvoice(w, r)
	}
}

func (th *tribeHandler) ProcessStake(w http.ResponseWriter, r *http.Request) {
	var stakeReq db.StakeInvoiceRequest

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Log.Error("Failed reading request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &stakeReq)
	if err != nil {
		logger.Log.Error("Failed unmarshaling request: %v", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	bountyIDStr := chi.URLParam(r, "bountyId")
	bountyIDUint, err := strconv.ParseUint(bountyIDStr, 10, 64)
	if err != nil {
		logger.Log.Error("Invalid bountyID: %v", err)
		http.Error(w, "Invalid bounty ID", http.StatusBadRequest)
		return
	}
	stakeReq.BountyID = uint(bountyIDUint)

	if !stakeReq.StakeOperation {
		http.Error(w, "Stake operation flag not set", http.StatusBadRequest)
		return
	}

	invoiceReq := db.BudgetInvoiceRequest{
		Amount:        stakeReq.Amount,
		SenderPubKey:  stakeReq.SenderPubKey,
		WorkspaceUuid: stakeReq.WorkspaceUuid,
		PaymentType:   stakeReq.PaymentType,
		BountyID:      stakeReq.BountyID,
	}

	modifiedBody, err := json.Marshal(invoiceReq)
	if err != nil {
		logger.Log.Error("Failed to marshal invoice request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(bytes.NewBuffer(modifiedBody))
	if stakeReq.StakeOperation {
		th.GenerateBudgetInvoice(w, r)
	}
}

func (th *tribeHandler) GenerateV1BudgetInvoice(w http.ResponseWriter, r *http.Request) {
	invoice := db.BudgetInvoiceRequest{}

	var err error
	body, err := io.ReadAll(r.Body)

	r.Body.Close()

	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = json.Unmarshal(body, &invoice)

	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if invoice.WorkspaceUuid == "" && invoice.OrgUuid != "" {
		invoice.WorkspaceUuid = invoice.OrgUuid
	}

	url := fmt.Sprintf("%s/invoices", config.RelayUrl)

	bodyData := fmt.Sprintf(`{"amount": %d, "memo": "%s"}`, invoice.Amount, "Budget Invoice")

	jsonBody := []byte(bodyData)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))

	req.Header.Set("x-user-token", config.RelayAuthKey)
	req.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(req)

	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)

	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	// Unmarshal result
	invoiceRes := db.InvoiceResponse{}

	err = json.Unmarshal(body, &invoiceRes)

	if err != nil {
		log.Printf("Json Unmarshal failed: %s", err)
		return
	}

	now := time.Now()
	var paymentHistory = db.NewPaymentHistory{
		Amount:         invoice.Amount,
		WorkspaceUuid:  invoice.WorkspaceUuid,
		PaymentType:    invoice.PaymentType,
		SenderPubKey:   invoice.SenderPubKey,
		ReceiverPubKey: "",
		Created:        &now,
		Updated:        &now,
		Status:         false,
		BountyId:       0,
	}

	newInvoice := db.NewInvoiceList{
		PaymentRequest: invoiceRes.Response.Invoice,
		Type:           db.InvoiceType("BUDGET"),
		OwnerPubkey:    invoice.SenderPubKey,
		WorkspaceUuid:  invoice.WorkspaceUuid,
		Created:        &now,
		Updated:        &now,
		Status:         false,
	}

	th.db.ProcessBudgetInvoice(paymentHistory, newInvoice)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceRes)
}

func (th *tribeHandler) GenerateV2BudgetInvoice(w http.ResponseWriter, r *http.Request) {
	invoice := db.BudgetInvoiceRequest{}

	var err error
	body, err := io.ReadAll(r.Body)

	r.Body.Close()

	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = json.Unmarshal(body, &invoice)

	if err != nil {
		logger.Log.Error("%v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if invoice.WorkspaceUuid == "" && invoice.OrgUuid != "" {
		invoice.WorkspaceUuid = invoice.OrgUuid
	}

	url := fmt.Sprintf("%s/invoice", config.V2BotUrl)

	amountMsat := invoice.Amount * 1000

	bodyData := fmt.Sprintf(`{"amt_msat": %d}`, amountMsat)

	jsonBody := []byte(bodyData)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))

	req.Header.Set("x-admin-token", config.V2BotToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		log.Printf("Client Request Failed: %s", err)
		return
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)

	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	// Unmarshal result
	v2InvoiceRes := db.V2CreateInvoiceResponse{}
	err = json.Unmarshal(body, &v2InvoiceRes)

	if err != nil {
		log.Printf("Json Unmarshal failed: %s", err)
		return
	}

	now := time.Now()
	var paymentHistory = db.NewPaymentHistory{
		Amount:         invoice.Amount,
		WorkspaceUuid:  invoice.WorkspaceUuid,
		PaymentType:    invoice.PaymentType,
		SenderPubKey:   invoice.SenderPubKey,
		ReceiverPubKey: "",
		Created:        &now,
		Updated:        &now,
		Status:         false,
		BountyId:       0,
	}

	newInvoice := db.NewInvoiceList{
		PaymentRequest: v2InvoiceRes.Bolt11,
		Type:           db.InvoiceType("BUDGET"),
		OwnerPubkey:    invoice.SenderPubKey,
		WorkspaceUuid:  invoice.WorkspaceUuid,
		Created:        &now,
		Updated:        &now,
		Status:         false,
	}

	th.db.ProcessBudgetInvoice(paymentHistory, newInvoice)

	invoiceRes := db.InvoiceResponse{
		Succcess: true,
		Response: db.Invoice{
			Invoice: v2InvoiceRes.Bolt11,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceRes)
}
