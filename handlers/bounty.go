package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/utils"
)

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

func GetBountyCount(w http.ResponseWriter, r *http.Request) {
	personKey := chi.URLParam(r, "personKey")
	tabType := chi.URLParam(r, "tabType")

	if personKey == "" || tabType == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bountyCount := db.DB.GetBountiesCount(personKey, tabType)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyCount)
}

func GetPersonCreatedBounties(w http.ResponseWriter, r *http.Request) {
	pubkey := chi.URLParam(r, "pubkey")
	if pubkey == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bounties, err := db.DB.GetCreatedBounties(pubkey)
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
	pubkey := chi.URLParam(r, "pubkey")
	if pubkey == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bounties, err := db.DB.GetAssignedBounties(pubkey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error", err)
	} else {
		var bountyResponse []db.BountyResponse = generateBountyResponse(bounties)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

func CreateOrEditBounty(w http.ResponseWriter, r *http.Request) {
	bounty := db.Bounty{}
	body, err := ioutil.ReadAll(r.Body)

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
	if bounty.Tribe == "" {
		bounty.Tribe = "None"
	}

	if bounty.Show == false && bounty.ID != 0 {
		db.DB.UpdateBountyBoolColumn(bounty, "show")
	}

	if bounty.Title != "" && bounty.Assignee == "" {
		db.DB.UpdateBountyNullColumn(bounty, "assignee")
	}

	b, err := db.DB.CreateOrEditBounty(bounty)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(b)
}

func DeleteBounty(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	//pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
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

	b, _ := db.DB.DeleteBounty(pubkey, created)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(b)
}

func UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	createdParam := chi.URLParam(r, "created")
	created, _ := strconv.ParseUint(createdParam, 10, 32)

	bounty, _ := db.DB.GetBountyByCreated(uint(created))
	if bounty.ID != 0 && bounty.Created == int64(created) {
		bounty.Paid = !bounty.Paid
		db.DB.UpdateBountyPayment(bounty)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bounty)
}

func generateBountyResponse(bounties []db.BountyData) []db.BountyResponse {
	var bountyResponse []db.BountyResponse

	for i := 0; i < len(bounties); i++ {
		bounty := bounties[i]
		b := db.BountyResponse{
			Bounty: db.Bounty{
				ID:                      bounty.BountyId,
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
				Created:                 bounty.BountyCreated,
				Assignee:                bounty.Assignee,
				TicketUrl:               bounty.TicketUrl,
				Description:             bounty.BountyDescription,
				WantedType:              bounty.WantedType,
				Deliverables:            bounty.Deliverables,
				GithubDescription:       bounty.GithubDescription,
				OneSentenceSummary:      bounty.OneSentenceSummary,
				EstimatedSessionLength:  bounty.EstimatedSessionLength,
				EstimatedCompletionDate: bounty.EstimatedCompletionDate,
				OrgUuid:                 bounty.OrgUuid,
				Updated:                 bounty.BountyUpdated,
				CodingLanguages:         bounty.CodingLanguages,
			},
			Assignee: db.Person{
				ID:               bounty.AssigneeId,
				Uuid:             bounty.Uuid,
				OwnerPubKey:      bounty.OwnerPubKey,
				OwnerAlias:       bounty.AssigneeAlias,
				UniqueName:       bounty.UniqueName,
				Description:      bounty.AssigneeDescription,
				Tags:             bounty.Tags,
				Img:              bounty.Img,
				Created:          bounty.AssigneeCreated,
				Updated:          bounty.AssigneeUpdated,
				LastLogin:        bounty.LastLogin,
				OwnerRouteHint:   bounty.AssigneeRouteHint,
				OwnerContactKey:  bounty.OwnerContactKey,
				PriceToMeet:      bounty.PriceToMeet,
				TwitterConfirmed: bounty.TwitterConfirmed,
			},
			Owner: db.Person{
				ID:               bounty.BountyOwnerId,
				Uuid:             bounty.OwnerUuid,
				OwnerPubKey:      bounty.OwnerKey,
				OwnerAlias:       bounty.OwnerAlias,
				UniqueName:       bounty.OwnerUniqueName,
				Description:      bounty.OwnerDescription,
				Tags:             bounty.OwnerTags,
				Img:              bounty.OwnerImg,
				Created:          bounty.OwnerCreated,
				Updated:          bounty.OwnerUpdated,
				LastLogin:        bounty.OwnerLastLogin,
				OwnerRouteHint:   bounty.OwnerRouteHint,
				OwnerContactKey:  bounty.OwnerContactKey,
				PriceToMeet:      bounty.OwnerPriceToMeet,
				TwitterConfirmed: bounty.OwnerTwitterConfirmed,
			},
			Organization: db.OrganizationShort{
				Name: bounty.OrganizationName,
				Uuid: bounty.OrganizationUuid,
				Img:  bounty.OrganizationImg,
			},
		}
		bountyResponse = append(bountyResponse, b)
	}

	return bountyResponse
}

func MakeBountyPayment(w http.ResponseWriter, r *http.Request) {
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
	amount, _ := utils.ConvertStringToUint(bounty.Price)

	if bounty.ID != id {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// check if user is the admin of the organization
	// or has a pay bounty role
	hasRole := db.UserHasAccess(pubKeyFromAuth, bounty.OrgUuid, db.PayBounty)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("You don't have appropriate permissions to pay bounties")
		return
	}

	// check if the orgnization bounty balance
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
		}
		db.DB.AddPaymentHistory(paymentHistory)
		bounty.Paid = true
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
}

func PollInvoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	paymentRequest := chi.URLParam(r, "paymentRequest")

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	url := fmt.Sprintf("%s/invoice?payment_request=%s", config.RelayUrl, paymentRequest)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set("x-user-token", config.RelayAuthKey)
	req.Header.Set("Content-Type", "application/json")
	res, _ := client.Do(req)

	if err != nil {
		log.Printf("Request Failed: %s", err)
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode("could not decode invoice")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	// Unmarshal result
	invoiceRes := db.InvoiceResult{}

	err = json.Unmarshal(body, &invoiceRes)

	if err != nil {
		log.Printf("Reading Invoice body failed: %s", err)
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode("could not decode invoice")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceRes)
}
