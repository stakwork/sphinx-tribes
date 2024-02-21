package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
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

func GetAllTribes(w http.ResponseWriter, r *http.Request) {
	tribes := db.DB.GetAllTribes()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribes)
}

func GetTotalribes(w http.ResponseWriter, r *http.Request) {
	tribesTotal := db.DB.GetTribesTotal()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribesTotal)
}

func GetListedTribes(w http.ResponseWriter, r *http.Request) {
	tribes := db.DB.GetListedTribes(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribes)
}

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

func (th *tribeHandler) GetTribesByAppUrl(w http.ResponseWriter, r *http.Request) {
	tribes := []db.Tribe{}
	app_url := chi.URLParam(r, "app_url")
	tribes = th.db.GetTribesByAppUrl(app_url)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribes)
}

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

func PutTribeStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	tribe := db.Tribe{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &tribe)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if tribe.UUID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := auth.VerifyTribeUUID(tribe.UUID, false)
	if err != nil {
		fmt.Println(err)
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

func DeleteTribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	uuid := chi.URLParam(r, "uuid")

	if uuid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := auth.VerifyTribeUUID(uuid, false)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// from token must match
	if pubKeyFromAuth != extractedPubkey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	db.DB.UpdateTribe(uuid, map[string]interface{}{
		"deleted": true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

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

func GetFirstTribeByFeed(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	tribe := db.DB.GetFirstTribeByFeedURL(url)

	if tribe.UUID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var theTribe map[string]interface{}
	j, _ := json.Marshal(tribe)
	json.Unmarshal(j, &theTribe)

	theTribe["channels"] = db.DB.GetChannelsByTribe(tribe.UUID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(theTribe)
}

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

func (th *tribeHandler) CreateOrEditTribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	tribe := db.Tribe{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	err = json.Unmarshal(body, &tribe)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if tribe.UUID == "" {
		fmt.Println("createOrEditTribe no uuid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	now := time.Now() //.Format(time.RFC3339)

	extractedPubkey, err := th.verifyTribeUUID(tribe.UUID, false)
	if err != nil {
		fmt.Println("extract UUID error", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if pubKeyFromAuth == "" {
		tribe.Created = &now
	} else { // IF PUBKEY IN CONTEXT, MUST AUTH!
		if pubKeyFromAuth != extractedPubkey {
			fmt.Println("createOrEditTribe pubkeys dont match")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	existing := th.db.GetTribe(tribe.UUID)
	if existing.UUID == "" { // if doesn't exist already, create unique name
		tribe.UniqueName, _ = th.tribeUniqueNameFromName(tribe.Name)
	} else { // already exists! make sure it's owned
		if existing.OwnerPubKey != extractedPubkey {
			fmt.Println("createOrEditTribe tribe.ownerPubKey not match")
			fmt.Println(existing.OwnerPubKey)
			fmt.Println(extractedPubkey)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	tribe.OwnerPubKey = extractedPubkey
	tribe.Updated = &now
	tribe.LastActive = now.Unix()

	_, err = th.db.CreateOrEditTribe(tribe)
	if err != nil {
		fmt.Println("=> ERR createOrEditTribe", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribe)
}

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
		fmt.Println(err)
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

func SetTribePreview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	extractedPubkey, err := auth.VerifyTribeUUID(uuid, false)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// from token must match
	if pubKeyFromAuth != extractedPubkey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	preview := r.URL.Query().Get("preview")
	db.DB.UpdateTribe(uuid, map[string]interface{}{
		"preview": preview,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

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
		fmt.Println(err)
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
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	_, err = db.DB.CreateLeaderBoard(uuid, leaderBoard)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(true)
}

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
		fmt.Println(err)
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
		fmt.Println(err)
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

func GenerateInvoice(w http.ResponseWriter, r *http.Request) {
	invoice := db.InvoiceRequest{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(body, &invoice)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	pub_key := invoice.User_pubkey
	owner_key := invoice.Owner_pubkey
	date, _ := utils.ConvertStringToInt(invoice.Created)
	memo := invoice.Memo
	invoiceType := invoice.Type
	routeHint := invoice.Route_hint
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
		return
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)

	// Unmarshal result
	invoiceRes := db.InvoiceResponse{}

	err = json.Unmarshal(body, &invoiceRes)

	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	paymentRequest := invoiceRes.Response.Invoice
	now := time.Now()

	newInvoice := db.InvoiceList{
		PaymentRequest: paymentRequest,
		Type:           db.InvoiceType(invoiceType),
		OwnerPubkey:    owner_key,
		Created:        &now,
		Updated:        &now,
		Status:         false,
	}

	db.DB.AddInvoice(newInvoice)

	newInvoiceData := db.UserInvoiceData{
		PaymentRequest: paymentRequest,
		Created:        date,
		Amount:         amount,
		UserPubkey:     pub_key,
		RouteHint:      routeHint,
	}

	db.DB.AddUserInvoiceData(newInvoiceData)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceRes)
}

func GenerateBudgetInvoice(w http.ResponseWriter, r *http.Request) {
	invoice := db.BudgetInvoiceRequest{}
	body, err := io.ReadAll(r.Body)

	r.Body.Close()

	err = json.Unmarshal(body, &invoice)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
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

	// Unmarshal result
	invoiceRes := db.InvoiceResponse{}

	err = json.Unmarshal(body, &invoiceRes)

	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	now := time.Now()
	var paymentHistory = db.PaymentHistory{
		Amount:         invoice.Amount,
		OrgUuid:        invoice.OrgUuid,
		PaymentType:    invoice.PaymentType,
		SenderPubKey:   invoice.SenderPubKey,
		ReceiverPubKey: "",
		Created:        &now,
		Updated:        &now,
		Status:         false,
		BountyId:       0,
	}

	newInvoice := db.InvoiceList{
		PaymentRequest: invoiceRes.Response.Invoice,
		Type:           db.InvoiceType("BUDGET"),
		OwnerPubkey:    invoice.SenderPubKey,
		OrgUuid:        invoice.OrgUuid,
		Created:        &now,
		Updated:        &now,
		Status:         false,
	}

	db.DB.AddPaymentHistory(paymentHistory)
	db.DB.AddInvoice(newInvoice)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceRes)
}

func makeInvoiceRequest(amount string, memo string) (*http.Response, error) {
	url := fmt.Sprintf("%s/invoices", config.RelayUrl)

	bodyData := fmt.Sprintf(`{"amount": %s, "memo": "%s"}`, amount, memo)

	jsonBody := []byte(bodyData)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))

	req.Header.Set("x-user-token", config.RelayAuthKey)
	req.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(req)

	return res, err
}
