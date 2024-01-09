package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/utils"
)

type bountyHandler struct {
	httpClient HttpClient
	db         db.Database
}

func NewBountyHandler(httpClient HttpClient, db db.Database) *bountyHandler {
	return &bountyHandler{
		httpClient: httpClient,
		db:         db,
	}
}

func GetAllBounties(w http.ResponseWriter, r *http.Request) {
	bounties := db.DB.GetAllBounties(r)
	var bountyResponse []db.BountyResponse = generateBountyResponse(bounties)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyResponse)
}

func GetBountyById(w http.ResponseWriter, r *http.Request) {
	bountyId := chi.URLParam(r, "bountyId")
	if bountyId == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bounties, err := db.DB.GetBountyById(bountyId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error", err)
	} else {
		var bountyResponse []db.BountyResponse = generateBountyResponse(bounties)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

func GetBountyIndexById(w http.ResponseWriter, r *http.Request) {
	bountyId := chi.URLParam(r, "bountyId")
	if bountyId == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bountyIndex := db.DB.GetBountyIndexById(bountyId)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyIndex)
}

func GetBountyByCreated(w http.ResponseWriter, r *http.Request) {
	created := chi.URLParam(r, "created")
	if created == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bounties, err := db.DB.GetBountyDataByCreated(created)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error", err)
	} else {
		var bountyResponse []db.BountyResponse = generateBountyResponse(bounties)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

func GetUserBountyCount(w http.ResponseWriter, r *http.Request) {
	personKey := chi.URLParam(r, "personKey")
	tabType := chi.URLParam(r, "tabType")

	if personKey == "" || tabType == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bountyCount := db.DB.GetUserBountiesCount(personKey, tabType)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyCount)
}

func GetBountyCount(w http.ResponseWriter, r *http.Request) {
	bountyCount := db.DB.GetBountiesCount(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyCount)
}

func GetPersonCreatedBounties(w http.ResponseWriter, r *http.Request) {
	bounties, err := db.DB.GetCreatedBounties(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error", err)
	} else {
		var bountyResponse []db.BountyResponse = generateBountyResponse(bounties)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

func GetPersonAssignedBounties(w http.ResponseWriter, r *http.Request) {
	bounties, err := db.DB.GetAssignedBounties(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error", err)
	} else {
		var bountyResponse []db.BountyResponse = generateBountyResponse(bounties)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

func (h *bountyHandler) CreateOrEditBounty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	bounty := db.Bounty{}
	body, err := io.ReadAll(r.Body)

	r.Body.Close()
	err = json.Unmarshal(body, &bounty)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	now := time.Now()

	//Check if bounty exists
	bounty.Updated = &now
	if bounty.Created == 0 {
		bounty.Created = time.Now().Unix()
	}

	if bounty.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Type is a required field")
		return
	}

	if bounty.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Title is a required field")
		return
	}

	if bounty.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Description is a required field")
		return
	}

	if bounty.Assignee != "" {
		now := time.Now()
		bounty.AssignedDate = &now
	}

	if bounty.Tribe == "" {
		bounty.Tribe = "None"
	}

	if bounty.Show == false && bounty.ID != 0 {
		h.db.UpdateBountyBoolColumn(bounty, "show")
	}

	if bounty.Title != "" && bounty.Assignee == "" {
		h.db.UpdateBountyNullColumn(bounty, "assignee")
	}

	if bounty.Title != "" && bounty.ID != 0 {
		// get bounty from DB
		dbBounty := h.db.GetBounty(bounty.ID)

		// trying to update
		// check if bounty belongs to user
		if pubKeyFromAuth != dbBounty.OwnerID {
			if bounty.OrgUuid != "" {
				hasBountyRoles := h.db.UserHasManageBountyRoles(pubKeyFromAuth, bounty.OrgUuid)
				if !hasBountyRoles {
					msg := "You don't have a=the right permission ton update bounty"
					fmt.Println(msg)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(msg)
					return
				}
			} else {
				msg := "Cannot edit another user's bounty"
				fmt.Println(msg)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(msg)
				return
			}
		}
	}

	b, err := h.db.CreateOrEditBounty(bounty)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(b)
}

func (h *bountyHandler) DeleteBounty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	created := chi.URLParam(r, "created")
	pubkey := chi.URLParam(r, "pubkey")

	if pubkey == "" {
		fmt.Println("no pubkey from route")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if created == "" {
		fmt.Println("no created timestamp from route")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	b, err := h.db.DeleteBounty(pubkey, created)
	if err != nil {
		fmt.Println("failed to delete bounty", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("failed to delete bounty")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(b)
}

func UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	createdParam := chi.URLParam(r, "created")
	created, _ := strconv.ParseUint(createdParam, 10, 32)

	bounty, _ := db.DB.GetBountyByCreated(uint(created))
	if bounty.ID != 0 && bounty.Created == int64(created) {
		bounty.Paid = !bounty.Paid
		now := time.Now()
		// if setting paid as true by mark as paid
		// set completion date and mark as paid
		if bounty.Paid {
			bounty.CompletionDate = &now
			bounty.MarkAsPaidDate = &now
		}
		db.DB.UpdateBountyPayment(bounty)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bounty)
}

func generateBountyResponse(bounties []db.Bounty) []db.BountyResponse {
	var bountyResponse []db.BountyResponse

	for i := 0; i < len(bounties); i++ {
		bounty := bounties[i]

		owner := db.DB.GetPersonByPubkey(bounty.OwnerID)
		assignee := db.DB.GetPersonByPubkey(bounty.Assignee)
		organization := db.DB.GetOrganizationByUuid(bounty.OrgUuid)

		b := db.BountyResponse{
			Bounty: db.Bounty{
				ID:                      bounty.ID,
				OwnerID:                 bounty.OwnerID,
				Paid:                    bounty.Paid,
				Show:                    bounty.Show,
				Type:                    bounty.Type,
				Award:                   bounty.Award,
				AssignedHours:           bounty.AssignedHours,
				BountyExpires:           bounty.BountyExpires,
				CommitmentFee:           bounty.CommitmentFee,
				Price:                   bounty.Price,
				Title:                   bounty.Title,
				Tribe:                   bounty.Tribe,
				Created:                 bounty.Created,
				Assignee:                bounty.Assignee,
				TicketUrl:               bounty.TicketUrl,
				Description:             bounty.Description,
				WantedType:              bounty.WantedType,
				Deliverables:            bounty.Deliverables,
				GithubDescription:       bounty.GithubDescription,
				OneSentenceSummary:      bounty.OneSentenceSummary,
				EstimatedSessionLength:  bounty.EstimatedSessionLength,
				EstimatedCompletionDate: bounty.EstimatedCompletionDate,
				OrgUuid:                 bounty.OrgUuid,
				Updated:                 bounty.Updated,
				CodingLanguages:         bounty.CodingLanguages,
			},
			Assignee: db.Person{
				ID:               assignee.ID,
				Uuid:             assignee.Uuid,
				OwnerPubKey:      assignee.OwnerPubKey,
				OwnerAlias:       assignee.OwnerAlias,
				UniqueName:       assignee.UniqueName,
				Description:      assignee.Description,
				Tags:             assignee.Tags,
				Img:              assignee.Img,
				Created:          assignee.Created,
				Updated:          assignee.Updated,
				LastLogin:        assignee.LastLogin,
				OwnerRouteHint:   assignee.OwnerRouteHint,
				OwnerContactKey:  assignee.OwnerContactKey,
				PriceToMeet:      assignee.PriceToMeet,
				TwitterConfirmed: assignee.TwitterConfirmed,
			},
			Owner: db.Person{
				ID:               owner.ID,
				Uuid:             owner.Uuid,
				OwnerPubKey:      owner.OwnerPubKey,
				OwnerAlias:       owner.OwnerAlias,
				UniqueName:       owner.UniqueName,
				Description:      owner.Description,
				Tags:             owner.Tags,
				Img:              owner.Img,
				Created:          owner.Created,
				Updated:          owner.Updated,
				LastLogin:        owner.LastLogin,
				OwnerRouteHint:   owner.OwnerRouteHint,
				OwnerContactKey:  owner.OwnerContactKey,
				PriceToMeet:      owner.PriceToMeet,
				TwitterConfirmed: owner.TwitterConfirmed,
			},
			Organization: db.OrganizationShort{
				Name: organization.Name,
				Uuid: organization.Uuid,
				Img:  organization.Img,
			},
		}
		bountyResponse = append(bountyResponse, b)
	}

	return bountyResponse
}

func MakeBountyPayment(w http.ResponseWriter, r *http.Request) {
	var m sync.Mutex
	m.Lock()

	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	idParam := chi.URLParam(r, "id")

	id, err := utils.ConvertStringToUint(idParam)
	if err != nil {
		fmt.Println("could not parse id")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bounty := db.DB.GetBounty(id)
	amount := bounty.Price

	if bounty.ID != id {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// check if the bounty has been paid already to avoid double payment
	if bounty.Paid {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode("Bounty has already been paid")
	}

	// check if user is the admin of the organization
	// or has a pay bounty role
	hasRole := db.UserHasAccess(pubKeyFromAuth, bounty.OrgUuid, db.PayBounty)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("You don't have appropriate permissions to pay bounties")
		return
	}

	// check if the organization bounty balance
	// is greater than the amount
	orgBudget := db.DB.GetOrganizationBudget(bounty.OrgUuid)
	if orgBudget.TotalBudget < amount {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("organization budget is not enough to pay the amount")
		return
	}

	request := db.BountyPayRequest{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	url := fmt.Sprintf("%s/payment", config.RelayUrl)

	bodyData := utils.BuildKeysendBodyData(amount, request.ReceiverPubKey, request.RouteHint)

	jsonBody := []byte(bodyData)

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("x-user-token", config.RelayAuthKey)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}

	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	msg := make(map[string]interface{})

	// payment is successful add to payment history
	// and reduce organizations budget
	if res.StatusCode == 200 {
		// Unmarshal result
		keysendRes := db.KeysendSuccess{}
		err = json.Unmarshal(body, &keysendRes)

		now := time.Now()
		paymentHistory := db.PaymentHistory{
			Amount:         amount,
			SenderPubKey:   pubKeyFromAuth,
			ReceiverPubKey: request.ReceiverPubKey,
			OrgUuid:        bounty.OrgUuid,
			BountyId:       id,
			Created:        &now,
			Updated:        &now,
			Status:         true,
			PaymentType:    "payment",
		}
		db.DB.AddPaymentHistory(paymentHistory)

		bounty.Paid = true
		bounty.PaidDate = &now
		bounty.CompletionDate = &now
		db.DB.UpdateBounty(bounty)

		msg["msg"] = "keysend_success"
		msg["invoice"] = ""

		socket, err := db.Store.GetSocketConnections(request.Websocket_token)
		if err == nil {
			socket.Conn.WriteJSON(msg)
		}
	} else {
		msg["msg"] = "keysend_error"
		msg["invoice"] = ""

		socket, err := db.Store.GetSocketConnections(request.Websocket_token)
		if err == nil {
			socket.Conn.WriteJSON(msg)
		}
	}

	m.Unlock()
}

func (h *bountyHandler) BountyBudgetWithdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := db.WithdrawBudgetRequest{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// check if user is the admin of the organization
	// or has a withdraw bounty budget role
	hasRole := db.UserHasAccess(pubKeyFromAuth, request.OrgUuid, db.WithdrawBudget)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		errMsg := formatPayError("You don't have appropriate permissions to withdraw bounty budget")
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	amount := utils.GetInvoiceAmount(request.PaymentRequest)

	if err == nil && amount > 0 {
		// check if the organization bounty balance
		// is greater than the amount
		orgBudget := db.DB.GetOrganizationBudget(request.OrgUuid)
		if amount > orgBudget.TotalBudget {
			w.WriteHeader(http.StatusForbidden)
			errMsg := formatPayError("Organization budget is not enough to withdraw the amount")
			json.NewEncoder(w).Encode(errMsg)
			return
		}
		paymentSuccess, paymentError := h.PayLightningInvoice(request.PaymentRequest)
		if paymentSuccess.Success {
			// withdraw amount from organization budget
			db.DB.WithdrawBudget(pubKeyFromAuth, request.OrgUuid, amount)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(paymentSuccess)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(paymentError)
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
		errMsg := formatPayError("Could not pay lightning invoice")
		json.NewEncoder(w).Encode(errMsg)
	}
}

func formatPayError(errorMsg string) db.InvoicePayError {
	return db.InvoicePayError{
		Success: false,
		Error:   errorMsg,
	}
}

func GetLightningInvoice(payment_request string) (db.InvoiceResult, db.InvoiceError) {
	url := fmt.Sprintf("%s/invoice?payment_request=%s", config.RelayUrl, payment_request)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set("x-user-token", config.RelayAuthKey)
	req.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(req)

	if err != nil {
		log.Printf("Request Failed: %s", err)
		return db.InvoiceResult{}, db.InvoiceError{}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		// Unmarshal result
		invoiceErr := db.InvoiceError{}
		err = json.Unmarshal(body, &invoiceErr)

		if err != nil {
			log.Printf("Reading Invoice body failed: %s", err)
			return db.InvoiceResult{}, invoiceErr
		}

		return db.InvoiceResult{}, invoiceErr
	} else {
		// Unmarshal result
		invoiceRes := db.InvoiceResult{}
		err = json.Unmarshal(body, &invoiceRes)

		if err != nil {
			log.Printf("Reading Invoice body failed: %s", err)
			return invoiceRes, db.InvoiceError{}
		}

		return invoiceRes, db.InvoiceError{}
	}
}

func (h *bountyHandler) PayLightningInvoice(payment_request string) (db.InvoicePaySuccess, db.InvoicePayError) {
	url := fmt.Sprintf("%s/invoices", config.RelayUrl)
	bodyData := fmt.Sprintf(`{"payment_request": "%s"}`, payment_request)
	jsonBody := []byte(bodyData)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))

	req.Header.Set("x-user-token", config.RelayAuthKey)
	req.Header.Set("Content-Type", "application/json")
	res, err := h.httpClient.Do(req)

	if err != nil {
		log.Printf("Request Failed: %s", err)
		return db.InvoicePaySuccess{}, db.InvoicePayError{}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		invoiceError := db.InvoicePayError{}
		err = json.Unmarshal(body, &invoiceError)

		if err != nil {
			log.Printf("Reading Invoice pay error body failed: %s", err)
			return db.InvoicePaySuccess{}, db.InvoicePayError{}
		}

		return db.InvoicePaySuccess{}, invoiceError
	} else {
		invoiceSuccess := db.InvoicePaySuccess{}
		err = json.Unmarshal(body, &invoiceSuccess)

		if err != nil {
			log.Printf("Reading Invoice pay success body failed: %s", err)
			return db.InvoicePaySuccess{}, db.InvoicePayError{}
		}

		return invoiceSuccess, db.InvoicePayError{}
	}
}

func GetInvoiceData(w http.ResponseWriter, r *http.Request) {
	paymentRequest := chi.URLParam(r, "paymentRequest")
	invoiceData, invoiceErr := GetLightningInvoice(paymentRequest)

	if invoiceErr.Error != "" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(invoiceErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceData)
}

func PollInvoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	paymentRequest := chi.URLParam(r, "paymentRequest")
	var err error

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	invoiceRes, invoiceErr := GetLightningInvoice(paymentRequest)

	if invoiceErr.Error != "" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(invoiceErr)
		return
	}

	if invoiceRes.Response.Settled {
		// Todo if an invoice is settled
		invoice := db.DB.GetInvoice(paymentRequest)
		invData := db.DB.GetUserInvoiceData(paymentRequest)
		dbInvoice := db.DB.GetInvoice(paymentRequest)

		// Make any change only if the invoice has not been settled
		if !dbInvoice.Status {
			if invoice.Type == "BUDGET" {
				db.DB.AddAndUpdateBudget(invoice)
			} else if invoice.Type == "KEYSEND" {
				url := fmt.Sprintf("%s/payment", config.RelayUrl)

				amount := invData.Amount

				bodyData := utils.BuildKeysendBodyData(amount, invData.UserPubkey, invData.RouteHint)

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

				body, _ := io.ReadAll(res.Body)

				if res.StatusCode == 200 {
					// Unmarshal result
					keysendRes := db.KeysendSuccess{}
					err = json.Unmarshal(body, &keysendRes)

					bounty, err := db.DB.GetBountyByCreated(uint(invData.Created))

					if err == nil {
						bounty.Paid = true
					}

					db.DB.UpdateBounty(bounty)
				} else {
					// Unmarshal result
					keysendError := db.KeysendError{}
					err = json.Unmarshal(body, &keysendError)
					log.Printf("Keysend Payment to %s Failed, with Error: %s", invData.UserPubkey, keysendError.Error)
				}
			}
			// Update the invoice status
			db.DB.UpdateInvoice(paymentRequest)
		}

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceRes)
}

func GetFilterCount(w http.ResponseWriter, r *http.Request) {
	filterCount := db.DB.GetFilterStatusCount()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filterCount)
}
