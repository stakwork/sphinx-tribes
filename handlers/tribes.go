package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
)

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

func GetTribesByOwner(w http.ResponseWriter, r *http.Request) {
	all := r.URL.Query().Get("all")
	tribes := []db.Tribe{}
	pubkey := chi.URLParam(r, "pubkey")
	if all == "true" {
		tribes = db.DB.GetAllTribesByOwner(pubkey)
	} else {
		tribes = db.DB.GetTribesByOwner(pubkey)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tribes)
}

func PutTribeStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	tribe := db.Tribe{}
	body, err := ioutil.ReadAll(r.Body)
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

func GetTribe(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	tribe := db.DB.GetTribe(uuid)

	var theTribe map[string]interface{}
	j, _ := json.Marshal(tribe)
	json.Unmarshal(j, &theTribe)

	theTribe["channels"] = db.DB.GetChannelsByTribe(uuid)

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

func GetTribeByUniqueName(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "un")
	tribe := db.DB.GetTribeByUniqueName(uuid)

	var theTribe map[string]interface{}
	j, _ := json.Marshal(tribe)
	json.Unmarshal(j, &theTribe)

	theTribe["channels"] = db.DB.GetChannelsByTribe(tribe.UUID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(theTribe)
}

func CreateOrEditTribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	tribe := db.Tribe{}
	body, err := ioutil.ReadAll(r.Body)
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

	extractedPubkey, err := auth.VerifyTribeUUID(tribe.UUID, false)
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

	existing := db.DB.GetTribe(tribe.UUID)
	if existing.UUID == "" { // doesnt exist already, create unique name
		tribe.UniqueName, _ = TribeUniqueNameFromName(tribe.Name)
	} else { // already exists! make sure its owned
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

	_, err = db.DB.CreateOrEditTribe(tribe)
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

	body, err := ioutil.ReadAll(r.Body)
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

	body, err := ioutil.ReadAll(r.Body)
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
	body, err := ioutil.ReadAll(r.Body)

	r.Body.Close()

	err = json.Unmarshal(body, &invoice)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	pub_key := invoice.User_pubkey
	owner_key := invoice.Owner_pubkey
	amount := invoice.Amount
	date := invoice.Created
	memo := invoice.Memo

	url := fmt.Sprintf("%s/invoices", config.RelayUrl)

	bodyData := fmt.Sprintf(`{"amount": %s, "memo": "%s"}`, amount, memo)

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

	body, err = ioutil.ReadAll(res.Body)

	// Unmarshal result
	invoiceRes := db.InvoiceResponse{}

	err = json.Unmarshal(body, &invoiceRes)

	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	var invoiceCache, _ = db.Store.GetInvoiceCache()
	var invoiceData = db.InvoiceStoreData{
		Amount:       amount,
		Created:      date,
		Invoice:      invoiceRes.Response.Invoice,
		Owner_pubkey: owner_key,
		User_pubkey:  pub_key,
	}

	var invoiceList = append(invoiceCache, invoiceData)

	// save the invoice to store
	db.Store.SetInvoiceCache(invoiceList)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceRes)
}

func GetInvoiceStatus(w http.ResponseWriter, r *http.Request) {
	payment_request := chi.URLParam(r, "payment_request")

	if payment_request == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var invoiceState bool
	var bountyPaid bool

	/**
	  if invoice is still in the store
	  It means the invoice has not been paid
	  else it has been paid
	*/
	invoiceList, _ := db.Store.GetInvoiceCache()
	invoiceLength := len(invoiceList)

	if invoiceLength > 0 {

		for _, invoice := range invoiceList {
			if invoice.Invoice == payment_request {
				invoiceState = false
				bountyPaid = false
			} else {
				invoiceState = true
				bountyPaid = true
			}
		}
	} else {
		invoiceState = true
		bountyPaid = true
	}

	invoiceData := db.InvoiceStatus{
		Status:          invoiceState,
		Payment_request: payment_request,
	}

	invoiceResult := make(map[string]interface{})

	invoiceResult["status"] = invoiceData.Status
	invoiceResult["payment_request"] = invoiceData.Payment_request
	invoiceResult["bounty_paid"] = bountyPaid

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceResult)
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
