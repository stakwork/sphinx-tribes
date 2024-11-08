package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/utils"
	"gorm.io/gorm"
)

type bountyHandler struct {
	httpClient               HttpClient
	db                       db.Database
	getSocketConnections     func(host string) (db.Client, error)
	generateBountyResponse   func(bounties []db.NewBounty) []db.BountyResponse
	userHasAccess            func(pubKeyFromAuth string, uuid string, role string) bool
	getInvoiceStatusByTag    func(tag string) db.V2TagRes
	getHoursDifference       func(createdDate int64, endDate *time.Time) int64
	userHasManageBountyRoles func(pubKeyFromAuth string, uuid string) bool
	m                        sync.Mutex
}

func NewBountyHandler(httpClient HttpClient, database db.Database) *bountyHandler {
	dbConf := db.NewDatabaseConfig(&gorm.DB{})
	return &bountyHandler{
		httpClient:               httpClient,
		db:                       database,
		getSocketConnections:     db.Store.GetSocketConnections,
		userHasAccess:            dbConf.UserHasAccess,
		getInvoiceStatusByTag:    GetInvoiceStatusByTag,
		getHoursDifference:       utils.GetHoursDifference,
		userHasManageBountyRoles: dbConf.UserHasManageBountyRoles,
	}
}

func (h *bountyHandler) GetAllBounties(w http.ResponseWriter, r *http.Request) {
	bounties := h.db.GetAllBounties(r)
	var bountyResponse []db.BountyResponse = h.GenerateBountyResponse(bounties)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyResponse)
}

func (h *bountyHandler) GetBountyById(w http.ResponseWriter, r *http.Request) {
	bountyId := chi.URLParam(r, "bountyId")
	if bountyId == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bounties, err := h.db.GetBountyById(bountyId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("[bounty] Error", err)
	} else {
		var bountyResponse []db.BountyResponse = h.GenerateBountyResponse(bounties)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

func (h *bountyHandler) GetNextBountyByCreated(w http.ResponseWriter, r *http.Request) {
	bounties, err := h.db.GetNextBountyByCreated(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("[bounty] Error", err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bounties)
	}
}

func (h *bountyHandler) GetPreviousBountyByCreated(w http.ResponseWriter, r *http.Request) {
	bounties, err := h.db.GetPreviousBountyByCreated(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("[bounty] Error", err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bounties)
	}
}

func (h *bountyHandler) GetWorkspaceNextBountyByCreated(w http.ResponseWriter, r *http.Request) {
	bounties, err := h.db.GetNextWorkspaceBountyByCreated(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("[bounty] Error", err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bounties)
	}
}

func (h *bountyHandler) GetWorkspacePreviousBountyByCreated(w http.ResponseWriter, r *http.Request) {
	bounties, err := h.db.GetPreviousWorkspaceBountyByCreated(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("[bounty] Error", err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bounties)
	}
}

func (h *bountyHandler) GetBountyIndexById(w http.ResponseWriter, r *http.Request) {
	bountyId := chi.URLParam(r, "bountyId")
	if bountyId == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bountyIndex := h.db.GetBountyIndexById(bountyId)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyIndex)
}

func (h *bountyHandler) GetBountyByCreated(w http.ResponseWriter, r *http.Request) {
	created := chi.URLParam(r, "created")
	if created == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bounties, err := h.db.GetBountyDataByCreated(created)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("[bounty] Error", err)
	} else {
		var bountyResponse []db.BountyResponse = h.GenerateBountyResponse(bounties)

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

func (h *bountyHandler) GetPersonCreatedBounties(w http.ResponseWriter, r *http.Request) {
	bounties, err := h.db.GetCreatedBounties(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("[bounty] Error", err)
	} else {
		var bountyResponse []db.BountyResponse = h.GenerateBountyResponse(bounties)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

func (h *bountyHandler) GetPersonAssignedBounties(w http.ResponseWriter, r *http.Request) {
	bounties, err := h.db.GetAssignedBounties(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("[bounty] Error", err)
	} else {
		var bountyResponse []db.BountyResponse = h.GenerateBountyResponse(bounties)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

func (h *bountyHandler) CreateOrEditBounty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	bounty := db.NewBounty{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		fmt.Println("[bounty read]", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = json.Unmarshal(body, &bounty)
	if err != nil {
		fmt.Println("[bounty]", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	now := time.Now()

	if bounty.WorkspaceUuid == "" && bounty.OrgUuid != "" {
		bounty.WorkspaceUuid = bounty.OrgUuid
	}

	//Check if bounty exists
	bounty.Updated = &now

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

	if !bounty.Show && bounty.ID != 0 {
		h.db.UpdateBountyBoolColumn(bounty, "show")
	}

	if bounty.Title != "" && bounty.Assignee == "" {
		h.db.UpdateBountyNullColumn(bounty, "assignee")
	}

	if bounty.ID == 0 && bounty.Created == 0 {
		bounty.Created = time.Now().Unix()
	}

	if bounty.Title != "" && bounty.ID != 0 {
		// get bounty from DB
		dbBounty := h.db.GetBounty(bounty.ID)

		// check if the bounty has a pending payment
		if dbBounty.PaymentPending {
			msg := "You cannot update a bounty with a pending payment"
			fmt.Println("[bounty]", msg)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(msg)
			return
		}

		// trying to update
		// check if bounty belongs to user
		if pubKeyFromAuth != dbBounty.OwnerID {
			if bounty.WorkspaceUuid != "" {
				hasBountyRoles := h.userHasManageBountyRoles(pubKeyFromAuth, bounty.WorkspaceUuid)
				if !hasBountyRoles {
					msg := "You don't have the right permission ton update bounty"
					fmt.Println("[bounty]", msg)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(msg)
					return
				}
			} else {
				msg := "Cannot edit another user's bounty"
				fmt.Println("[bounty]", msg)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(msg)
				return
			}
		}
	}

	if bounty.PhaseUuid != "" {
		phase, err := h.db.GetPhaseByUuid(bounty.PhaseUuid)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Phase Error")
			return
		}
		if bounty.PhaseUuid != phase.Uuid {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Not a valid phase")
			return
		}
	}

	b, err := h.db.CreateOrEditBounty(bounty)
	if err != nil {
		fmt.Println("[bounty]", err)
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
		fmt.Println("[bounty] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	created := chi.URLParam(r, "created")
	pubkey := chi.URLParam(r, "pubkey")

	if pubkey == "" {
		fmt.Println("[bounty] no pubkey from route")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if created == "" {
		fmt.Println("[bounty] no created timestamp from route")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get bounty by created
	createdUint, _ := utils.ConvertStringToUint(created)
	createdBounty, err := h.db.GetBountyByCreated(createdUint)
	if err != nil {
		fmt.Println("[bounty] failed to delete bounty", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("failed to delete bounty")
		return
	}

	if createdBounty.ID == 0 {
		fmt.Println("[bounty] failed to delete bounty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("failed to delete bounty")
		return
	}

	b, err := h.db.DeleteBounty(pubkey, created)
	if err != nil {
		fmt.Println("[bounty] failed to delete bounty", err.Error())
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
	if bounty.PaymentPending {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode("Cannot update a bounty with a pending payment")
		return
	}

	if bounty.ID != 0 && bounty.Created == int64(created) {
		bounty.Paid = !bounty.Paid
		now := time.Now()
		// if setting paid as true by mark as paid
		// set completion date and mark as paid
		if bounty.Paid {
			bounty.Completed = true
			bounty.CompletionDate = &now
			bounty.MarkAsPaidDate = &now

			if bounty.PaidDate == nil {
				bounty.PaidDate = &now
			}
		}
		db.DB.UpdateBountyPayment(bounty)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bounty)
}

func UpdateCompletedStatus(w http.ResponseWriter, r *http.Request) {
	createdParam := chi.URLParam(r, "created")
	created, _ := strconv.ParseUint(createdParam, 10, 32)
	bounty, _ := db.DB.GetBountyByCreated(uint(created))

	if bounty.PaymentPending {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode("Cannot update a bounty with a pending payment")
		return
	}

	if bounty.ID != 0 && bounty.Created == int64(created) {
		now := time.Now()
		// set bounty as completed
		if !bounty.Paid && !bounty.Completed {
			bounty.CompletionDate = &now
			bounty.Completed = true
		}
		db.DB.UpdateBountyCompleted(bounty)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bounty)
}

func GetPaymentByBountyId(w http.ResponseWriter, r *http.Request) {
	bountyIdParam := chi.URLParam(r, "bountyId")
	bountyId, _ := strconv.ParseUint(bountyIdParam, 10, 32)
	payment := db.DB.GetPaymentByBountyId(uint(bountyId))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

func (h *bountyHandler) GenerateBountyResponse(bounties []db.NewBounty) []db.BountyResponse {
	var bountyResponse []db.BountyResponse

	for i := 0; i < len(bounties); i++ {
		bounty := bounties[i]

		owner := h.db.GetPersonByPubkey(bounty.OwnerID)
		assignee := h.db.GetPersonByPubkey(bounty.Assignee)
		workspace := h.db.GetWorkspaceByUuid(bounty.WorkspaceUuid)

		b := db.BountyResponse{
			Bounty: db.NewBounty{
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
				OrgUuid:                 bounty.WorkspaceUuid,
				WorkspaceUuid:           bounty.WorkspaceUuid,
				Updated:                 bounty.Updated,
				CodingLanguages:         bounty.CodingLanguages,
				Completed:               bounty.Completed,
				PaymentPending:          bounty.PaymentPending,
				PaymentFailed:           bounty.PaymentFailed,
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
			Organization: db.WorkspaceShort{
				Name: workspace.Name,
				Uuid: workspace.Uuid,
				Img:  workspace.Img,
			},
			Workspace: db.WorkspaceShort{
				Name: workspace.Name,
				Uuid: workspace.Uuid,
				Img:  workspace.Img,
			},
		}
		bountyResponse = append(bountyResponse, b)
	}

	return bountyResponse
}

func (h *bountyHandler) MakeBountyPayment(w http.ResponseWriter, r *http.Request) {
	h.m.Lock()

	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	idParam := chi.URLParam(r, "id")

	id, err := utils.ConvertStringToUint(idParam)
	if err != nil {
		fmt.Println("[bounty] could not parse id")
		w.WriteHeader(http.StatusForbidden)
		h.m.Unlock()
		return
	}

	if pubKeyFromAuth == "" {
		fmt.Println("[bounty] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		h.m.Unlock()
		return
	}

	bounty := h.db.GetBounty(id)
	amount := bounty.Price

	if bounty.WorkspaceUuid == "" && bounty.OrgUuid != "" {
		bounty.WorkspaceUuid = bounty.OrgUuid
	}

	if bounty.ID != id {
		w.WriteHeader(http.StatusNotFound)
		h.m.Unlock()
		return
	}

	// check if the bounty has been paid already to avoid double payment
	if bounty.Paid {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode("Bounty has already been paid")
		h.m.Unlock()
		return
	}

	if bounty.PaymentPending {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Bounty payemnt is pending, cannot retry payment")
		h.m.Unlock()
		return
	}

	// check if user is the admin of the workspace
	// or has a pay bounty role
	hasRole := h.userHasAccess(pubKeyFromAuth, bounty.WorkspaceUuid, db.PayBounty)
	if !hasRole {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("You don't have appropriate permissions to pay bounties")
		h.m.Unlock()
		return
	}

	// check if the workspace bounty balance
	// is greater than the amount
	orgBudget := h.db.GetWorkspaceBudget(bounty.WorkspaceUuid)
	if orgBudget.TotalBudget < amount {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("workspace budget is not enough to pay the amount")
		h.m.Unlock()
		return
	}

	request := db.BountyPayRequest{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Println("[read body]", err)
		w.WriteHeader(http.StatusNotAcceptable)
		h.m.Unlock()
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("[bounty]", err)
		w.WriteHeader(http.StatusNotAcceptable)
		h.m.Unlock()
		return
	}

	// Get Bounty Assignee
	assignee := h.db.GetPersonByPubkey(bounty.Assignee)

	memoData := fmt.Sprintf("Payment For: %ss", bounty.Title)
	memoText := url.QueryEscape(memoData)
	now := time.Now()

	// If the v2contactkey is present
	if config.IsV2Payment {
		url := fmt.Sprintf("%s/pay", config.V2BotUrl)

		fmt.Println("IS V2 PAYMENT ====")

		// Build v2 keysend payment data
		bodyData := utils.BuildV2KeysendBodyData(amount, assignee.OwnerPubKey, assignee.OwnerRouteHint, memoText)
		jsonBody := []byte(bodyData)

		log.Println("Payment Body Data", bodyData)

		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("x-admin-token", config.V2BotToken)
		req.Header.Set("Content-Type", "application/json")
		log.Printf("[bounty] Making Bounty V2 Payment: amount: %d, pubkey: %s, route_hint: %s", amount, assignee.OwnerPubKey, assignee.OwnerRouteHint)

		res, err := h.httpClient.Do(req)

		if err != nil {
			log.Printf("[bounty] Request Failed: %s", err)
			h.m.Unlock()
			return
		}

		log.Printf("[bounty] After Making Bounty V2 Payment: amount: %d, pubkey: %s, route_hint: %s", amount, assignee.OwnerPubKey, assignee.OwnerRouteHint)

		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Println("[read body failed]", err)
			w.WriteHeader(http.StatusNotAcceptable)
			h.m.Unlock()
			return
		}

		log.Println("[bounty] After Reading Keysend V2 Payment Body ===")

		msg := make(map[string]interface{})
		// payment is successful add to payment history
		// and reduce workspaces budget

		paymentHistory := db.NewPaymentHistory{
			Amount:         amount,
			SenderPubKey:   pubKeyFromAuth,
			ReceiverPubKey: assignee.OwnerPubKey,
			WorkspaceUuid:  bounty.WorkspaceUuid,
			BountyId:       id,
			Created:        &now,
			Updated:        &now,
			Status:         false,
			PaymentType:    "payment",
			Tag:            "",
			PaymentStatus:  "FAILED",
		}

		if res.StatusCode == 200 {
			// Unmarshal result
			v2KeysendRes := db.V2SendOnionRes{}
			err = json.Unmarshal(body, &v2KeysendRes)

			if err != nil {
				fmt.Println("[Unmarshal failed]", err)
				w.WriteHeader(http.StatusNotAcceptable)
				h.m.Unlock()
				return
			}

			log.Printf("[bounty] V2 Status After Making Bounty V2 Payment: amount: %d, pubkey: %s, route_hint: %s is : %s", amount, assignee.OwnerPubKey, assignee.OwnerRouteHint, v2KeysendRes.Status)

			// if the payment has a completed status
			if v2KeysendRes.Status == db.PaymentComplete {
				bounty.Paid = true
				bounty.PaymentFailed = false
				bounty.PaymentPending = false
				bounty.PaidDate = &now
				bounty.Completed = true
				bounty.CompletionDate = &now
				paymentHistory.Status = true
				paymentHistory.PaymentStatus = db.PaymentComplete
				paymentHistory.Tag = v2KeysendRes.Tag

				h.db.ProcessBountyPayment(paymentHistory, bounty)

				msg["msg"] = "keysend_success"
				msg["invoice"] = ""

				socket, err := h.getSocketConnections(request.Websocket_token)
				if err == nil {
					socket.Conn.WriteJSON(msg)
				}

				h.m.Unlock()

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(msg)
				return
			} else if v2KeysendRes.Status == db.PaymentPending {
				// Send payment status
				log.Printf("[bounty] V2 Status is pending:  %s", v2KeysendRes.Status)
				bounty.Paid = false
				bounty.PaymentPending = true
				bounty.PaymentFailed = false
				bounty.PaidDate = &now
				bounty.Completed = true
				bounty.CompletionDate = &now
				paymentHistory.Status = true
				paymentHistory.PaymentStatus = db.PaymentPending
				paymentHistory.Tag = v2KeysendRes.Tag

				h.db.ProcessBountyPayment(paymentHistory, bounty)

				msg["msg"] = "keysend_pending"
				msg["invoice"] = ""

				socket, err := h.getSocketConnections(request.Websocket_token)
				if err == nil {
					socket.Conn.WriteJSON(msg)
				}

				h.m.Unlock()

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(msg)
				return
			} else {
				// Send payment status
				log.Printf("[bounty] V2 Status Was not completed:  %s", v2KeysendRes.Status)

				bounty.Paid = false
				bounty.PaymentPending = false
				bounty.PaymentFailed = true

				// set the error message
				paymentHistory.Error = v2KeysendRes.Message
				paymentHistory.PaymentStatus = db.PaymentFailed
				paymentHistory.Tag = v2KeysendRes.Tag

				h.db.AddPaymentHistory(paymentHistory)
				h.db.UpdateBounty(bounty)

				log.Println("Keysend payment not completed ===")
				msg["msg"] = "keysend_error"
				msg["invoice"] = ""

				socket, err := h.getSocketConnections(request.Websocket_token)
				if err == nil {
					socket.Conn.WriteJSON(msg)
				}

				h.m.Unlock()

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(msg)
				return
			}
		} else { // Send Payment error
			log.Println("Keysend payment error: Failed to send ===")
			msg["msg"] = "keysend_error"
			msg["invoice"] = ""

			bounty.Paid = false
			bounty.PaymentPending = false
			bounty.PaymentFailed = true

			// set the error message
			paymentHistory.Error = "Payment Request Failed"

			h.db.AddPaymentHistory(paymentHistory)
			h.db.UpdateBounty(bounty)

			socket, err := h.getSocketConnections(request.Websocket_token)
			if err == nil {
				socket.Conn.WriteJSON(msg)
			}

			h.m.Unlock()

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(msg)
			return
		}
	} else { // Process v1 payment
		url := fmt.Sprintf("%s/payment", config.RelayUrl)

		bodyData := utils.BuildKeysendBodyData(amount, assignee.OwnerPubKey, assignee.OwnerRouteHint, memoText)
		jsonBody := []byte(bodyData)

		log.Println("Payment Body Data", bodyData)

		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("x-user-token", config.RelayAuthKey)
		req.Header.Set("Content-Type", "application/json")
		log.Printf("[bounty] Making Bounty Payment: amount: %d, pubkey: %s, route_hint: %s", amount, assignee.OwnerPubKey, assignee.OwnerRouteHint)
		res, err := h.httpClient.Do(req)

		if err != nil {
			log.Printf("[bounty] Request Failed: %s", err)
			h.m.Unlock()
			return
		}

		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("[read body]", err)
			w.WriteHeader(http.StatusNotAcceptable)
			h.m.Unlock()
			return
		}

		msg := make(map[string]interface{})

		// payment is successful add to payment history
		// and reduce workspaces budget
		if res.StatusCode == 200 {
			// Unmarshal result
			keysendRes := db.KeysendSuccess{}
			err = json.Unmarshal(body, &keysendRes)

			if err != nil {
				fmt.Println("[Unmarshal]", err)
				w.WriteHeader(http.StatusNotAcceptable)
				h.m.Unlock()
				return
			}

			now := time.Now()

			paymentHistory := db.NewPaymentHistory{
				Amount:         amount,
				SenderPubKey:   pubKeyFromAuth,
				ReceiverPubKey: assignee.OwnerPubKey,
				WorkspaceUuid:  bounty.WorkspaceUuid,
				BountyId:       id,
				Created:        &now,
				Updated:        &now,
				Status:         true,
				PaymentType:    "payment",
			}

			bounty.Paid = true
			bounty.PaidDate = &now
			bounty.Completed = true
			bounty.CompletionDate = &now

			h.db.ProcessBountyPayment(paymentHistory, bounty)

			msg["msg"] = "keysend_success"
			msg["invoice"] = ""

			socket, err := h.getSocketConnections(request.Websocket_token)
			if err == nil {
				socket.Conn.WriteJSON(msg)
			}
			h.m.Unlock()
			return
		} else {
			msg["msg"] = "keysend_error"
			msg["invoice"] = ""

			socket, err := h.getSocketConnections(request.Websocket_token)
			if err == nil {
				socket.Conn.WriteJSON(msg)
			}

			h.m.Unlock()
			return
		}
	}
}

func (h *bountyHandler) GetBountyPaymentStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	idParam := chi.URLParam(r, "id")

	id, err := utils.ConvertStringToUint(idParam)
	if err != nil {
		fmt.Println("[bounty] could not parse id")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if pubKeyFromAuth == "" {
		fmt.Println("[bounty] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bounty := h.db.GetBounty(id)

	// check if the bounty has been paid already to avoid double payment
	if bounty.Paid {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode("Bounty has already been paid")
		return
	}

	payment := h.db.GetPaymentByBountyId(bounty.ID)

	if payment.Tag == "" {
		w.WriteHeader(http.StatusBadRequest)
		res := db.NewPaymentHistory{
			Status:        false,
			PaymentStatus: db.PaymentNotFound,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

func (h *bountyHandler) UpdateBountyPaymentStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	idParam := chi.URLParam(r, "id")

	id, err := utils.ConvertStringToUint(idParam)
	if err != nil {
		fmt.Println("[bounty] could not parse id")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if pubKeyFromAuth == "" {
		fmt.Println("[bounty] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bounty := h.db.GetBounty(id)

	if bounty.WorkspaceUuid == "" && bounty.OrgUuid != "" {
		bounty.WorkspaceUuid = bounty.OrgUuid
	}

	if bounty.ID != id {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// check if the bounty has been paid already to avoid double payment
	if bounty.Paid {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode("Bounty has already been paid")
		return
	}

	payment := h.db.GetPaymentByBountyId(bounty.ID)

	if payment.Tag != "" {
		tag := payment.Tag

		tagResult := h.getInvoiceStatusByTag(tag)

		msg := map[string]string{
			"payment_status": tagResult.Status,
		}

		if tagResult.Status == db.PaymentComplete {
			// Update only if it is still pending
			if payment.PaymentStatus == db.PaymentPending {
				h.db.SetPaymentAsComplete(tag)
			}

			now := time.Now()

			bounty.PaymentPending = false
			bounty.PaymentFailed = false
			bounty.Paid = true
			bounty.PaidDate = &now
			bounty.Completed = true
			bounty.CompletionDate = &now

			h.db.UpdateBounty(bounty)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(msg)
			return
		} else if tagResult.Status == db.PaymentFailed {
			// Handle failed payments
			h.db.SetPaymentStatusByBountyId(bounty.ID, tagResult)

			bounty.Paid = false
			bounty.PaymentPending = false
			bounty.PaymentFailed = true

			h.db.UpdateBounty(bounty)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(msg)
			return
		}
	}

	msg := map[string]string{
		"payment_status": db.PaymentNotFound,
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(msg)
}

// Todo: change back to BountyBudgetWithdraw
func (h *bountyHandler) BountyBudgetWithdraw(w http.ResponseWriter, r *http.Request) {
	h.m.Lock()

	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("[bounty] no pubkey from auth")
		h.m.Unlock()

		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := db.NewWithdrawBudgetRequest{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		h.m.Unlock()

		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		h.m.Unlock()

		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	lastWithdrawal := h.db.GetLastWithdrawal(request.WorkspaceUuid)

	if lastWithdrawal.ID > 0 {
		now := time.Now()
		withdrawCreated := lastWithdrawal.Created
		withdrawTime := utils.ConvertTimeToTimestamp(withdrawCreated.String())

		hoursDiff := h.getHoursDifference(int64(withdrawTime), &now)

		// Check that last withdraw time is greater than 1
		if hoursDiff < 1 {
			h.m.Unlock()

			w.WriteHeader(http.StatusUnauthorized)
			errMsg := formatPayError("Your last withdrawal is  not more than an hour ago")
			log.Println("Your last withdrawal is not more than an hour ago", hoursDiff, lastWithdrawal.Created, request.WorkspaceUuid)
			json.NewEncoder(w).Encode(errMsg)
			return
		}
	}

	log.Printf("[bounty] [BountyBudgetWithdraw] Logging body: workspace_uuid: %s, pubkey: %s, invoice: %s", request.WorkspaceUuid, pubKeyFromAuth, request.PaymentRequest)

	// check if user is the admin of the workspace
	// or has a withdraw bounty budget role
	hasRole := h.userHasAccess(pubKeyFromAuth, request.WorkspaceUuid, db.WithdrawBudget)
	if !hasRole {
		h.m.Unlock()

		w.WriteHeader(http.StatusUnauthorized)
		errMsg := formatPayError("You don't have appropriate permissions to withdraw bounty budget")
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	amount := utils.GetInvoiceAmount(request.PaymentRequest)

	if amount > 0 {
		// check if the workspace bounty balance
		// is greater than the amount
		orgBudget := h.db.GetWorkspaceBudget(request.WorkspaceUuid)
		if amount > orgBudget.TotalBudget {
			h.m.Unlock()

			w.WriteHeader(http.StatusForbidden)
			errMsg := formatPayError("Workspace budget is not enough to withdraw the amount")
			json.NewEncoder(w).Encode(errMsg)
			return
		}

		// Check that the deposit is more than the withdrawal plus amount to withdraw
		sumOfWithdrawals := h.db.GetSumOfWithdrawal(request.WorkspaceUuid)
		sumOfDeposits := h.db.GetSumOfDeposits(request.WorkspaceUuid)

		if sumOfDeposits < sumOfWithdrawals+amount {
			h.m.Unlock()

			w.WriteHeader(http.StatusUnauthorized)
			errMsg := formatPayError("Your deposits is lesser than your withdral")
			json.NewEncoder(w).Encode(errMsg)
			return
		}

		paymentSuccess, paymentError := h.PayLightningInvoice(request.PaymentRequest)
		if paymentSuccess.Success {
			// withdraw amount from workspace budget
			h.db.WithdrawBudget(pubKeyFromAuth, request.WorkspaceUuid, amount)

			h.m.Unlock()

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(paymentSuccess)
		} else {
			h.m.Unlock()

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(paymentError)
		}
	} else {
		h.m.Unlock()

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

func (h *bountyHandler) GetLightningInvoice(payment_request string) (db.InvoiceResult, db.InvoiceError) {
	if config.IsV2Payment {
		return h.GetV2LightningInvoice(payment_request)
	} else {
		return h.GetV1LightningInvoice(payment_request)
	}
}

func (h *bountyHandler) GetV1LightningInvoice(payment_request string) (db.InvoiceResult, db.InvoiceError) {
	url := fmt.Sprintf("%s/invoice?payment_request=%s", config.RelayUrl, payment_request)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set("x-user-token", config.RelayAuthKey)
	req.Header.Set("Content-Type", "application/json")
	res, _ := h.httpClient.Do(req)

	if err != nil {
		log.Printf("[bounty] Request Failed: %s", err)
		return db.InvoiceResult{}, db.InvoiceError{}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Printf("Error reading: %s", err)
		return db.InvoiceResult{}, db.InvoiceError{Success: false, Error: err.Error()}
	}

	if res.StatusCode != 200 {
		// Unmarshal result
		invoiceErr := db.InvoiceError{}
		err = json.Unmarshal(body, &invoiceErr)

		if err != nil {
			log.Printf("[bounty] Reading Invoice body failed: %s", err)
			return db.InvoiceResult{}, invoiceErr
		}

		return db.InvoiceResult{}, invoiceErr
	} else {
		// Unmarshal result
		invoiceRes := db.InvoiceResult{}
		err = json.Unmarshal(body, &invoiceRes)

		if err != nil {
			log.Printf("[bounty] Reading Invoice body failed: %s", err)
			return invoiceRes, db.InvoiceError{}
		}

		return invoiceRes, db.InvoiceError{}
	}
}

func (h *bountyHandler) GetV2LightningInvoice(payment_request string) (db.InvoiceResult, db.InvoiceError) {
	url := fmt.Sprintf("%s/check_invoice", config.V2BotUrl)

	invoiceBody := db.V2InvoiceBody{
		Bolt11: payment_request,
	}

	jsonBody, _ := json.Marshal(invoiceBody)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("x-admin-token", config.V2BotToken)
	req.Header.Set("Content-Type", "application/json")
	res, _ := h.httpClient.Do(req)

	if err != nil {
		log.Printf("[bounty] Request Failed: %s", err)
		return db.InvoiceResult{}, db.InvoiceError{Success: false, Error: err.Error()}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Printf("[bounty] Reading Invoice body failed: %s", err)
		return db.InvoiceResult{}, db.InvoiceError{Success: false, Error: err.Error()}
	}

	if res.StatusCode != 200 {
		// Unmarshal result
		invoiceErr := db.InvoiceError{}
		err = json.Unmarshal(body, &invoiceErr)

		if err != nil {
			log.Printf("[bounty] Unmarshalling Invoice body failed: %s", err)
			return db.InvoiceResult{}, invoiceErr
		}

		return db.InvoiceResult{}, invoiceErr
	} else {
		// Unmarshal result
		invoiceRes := db.V2InvoiceResponse{}
		err = json.Unmarshal(body, &invoiceRes)

		if err != nil {
			log.Printf("[bounty] Reading Invoice body failed: %s", err)
			return db.InvoiceResult{}, db.InvoiceError{}
		}

		invoiceResult := db.InvoiceResult{
			Success: false,
			Response: db.InvoiceCheckResponse{
				Settled:         false,
				Payment_request: payment_request,
				Payment_hash:    "",
				Preimage:        "",
			},
		}

		if invoiceRes.Status == db.InvoicePaid {
			invoiceResult.Success = true
			invoiceResult.Response.Settled = true
			return invoiceResult, db.InvoiceError{}
		}
		return invoiceResult, db.InvoiceError{}
	}
}

func (h *bountyHandler) PayLightningInvoice(payment_request string) (db.InvoicePaySuccess, db.InvoicePayError) {
	if config.IsV2Payment {
		return h.PayV2LightningInvoice(payment_request)
	} else {
		return h.PayV1LightningInvoice(payment_request)
	}
}

func (h *bountyHandler) PayV1LightningInvoice(payment_request string) (db.InvoicePaySuccess, db.InvoicePayError) {
	url := fmt.Sprintf("%s/invoices", config.RelayUrl)
	bodyData := fmt.Sprintf(`{"payment_request": "%s"}`, payment_request)
	jsonBody := []byte(bodyData)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))

	if err != nil {
		log.Printf("Error paying invoice: %s", err)
	}

	req.Header.Set("x-user-token", config.RelayAuthKey)
	req.Header.Set("Content-Type", "application/json")
	res, err := h.httpClient.Do(req)

	if err != nil {
		log.Printf("[bounty] Request Failed: %s", err)
		return db.InvoicePaySuccess{}, db.InvoicePayError{}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Printf("Error could not read body: %s", err)
	}

	if res.StatusCode != 200 {
		invoiceError := db.InvoicePayError{}
		err = json.Unmarshal(body, &invoiceError)

		if err != nil {
			log.Printf("[bounty] Reading Invoice pay error body failed: %s", err)
			return db.InvoicePaySuccess{}, db.InvoicePayError{}
		}

		return db.InvoicePaySuccess{}, invoiceError
	} else {
		invoiceSuccess := db.InvoicePaySuccess{}
		err = json.Unmarshal(body, &invoiceSuccess)

		if err != nil {
			log.Printf("[bounty] Reading Invoice pay success body failed: %s", err)
			return db.InvoicePaySuccess{}, db.InvoicePayError{}
		}

		return invoiceSuccess, db.InvoicePayError{}
	}
}

func (h *bountyHandler) PayV2LightningInvoice(payment_request string) (db.InvoicePaySuccess, db.InvoicePayError) {
	url := fmt.Sprintf("%s/pay_invoice", config.V2BotUrl)
	bodyData := fmt.Sprintf(`{"bolt11": "%s", "wait": true}`, payment_request)
	jsonBody := []byte(bodyData)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))

	if err != nil {
		log.Printf("Error paying invoice: %s", err)
		return db.InvoicePaySuccess{}, db.InvoicePayError{}
	}

	req.Header.Set("x-admin-token", config.V2BotToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := h.httpClient.Do(req)

	if err != nil {
		log.Printf("[bounty] Request Failed: %s", err)
		return db.InvoicePaySuccess{}, db.InvoicePayError{}
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Printf("Error could not read body: %s", err)
	}

	if res.StatusCode != 200 {
		invoiceError := db.InvoicePayError{}
		err = json.Unmarshal(body, &invoiceError)

		if err != nil {
			log.Printf("[bounty] Reading Invoice pay error body failed: %s", err)
			return db.InvoicePaySuccess{}, db.InvoicePayError{}
		}

		return db.InvoicePaySuccess{}, invoiceError
	} else {
		invoiceRes := db.V2InvoiceResponse{}
		err = json.Unmarshal(body, &invoiceRes)

		if err != nil {
			log.Printf("[bounty] Reading Invoice pay success body failed: %s", err)
			return db.InvoicePaySuccess{}, db.InvoicePayError{}
		}

		invoiceResult := db.InvoicePaySuccess{
			Success: false,
			Response: db.InvoiceCheckResponse{
				Settled:         false,
				Payment_request: payment_request,
				Payment_hash:    "",
				Preimage:        "",
			},
		}

		if invoiceRes.Status == db.PaymentComplete {
			invoiceResult.Success = true
			invoiceResult.Response.Settled = true
			return invoiceResult, db.InvoicePayError{}
		}

		return invoiceResult, db.InvoicePayError{}
	}
}

func (h *bountyHandler) GetInvoiceData(w http.ResponseWriter, r *http.Request) {
	paymentRequest := chi.URLParam(r, "paymentRequest")
	invoiceData, invoiceErr := h.GetLightningInvoice(paymentRequest)

	if invoiceErr.Error != "" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(invoiceErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceData)
}

func (h *bountyHandler) PollInvoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	paymentRequest := chi.URLParam(r, "paymentRequest")

	if pubKeyFromAuth == "" {
		fmt.Println("[bounty] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	invoiceRes, invoiceErr := h.GetLightningInvoice(paymentRequest)

	if invoiceErr.Error != "" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(invoiceErr)
		return
	}

	if invoiceRes.Response.Settled {
		// Todo if an invoice is settled
		invoice := h.db.GetInvoice(paymentRequest)
		dbInvoice := h.db.GetInvoice(paymentRequest)

		// Make any change only if the invoice has not been settled
		if !dbInvoice.Status {
			if invoice.Type == "BUDGET" {
				h.db.AddAndUpdateBudget(invoice)
			}
			// Update the invoice status
			h.db.UpdateInvoice(paymentRequest)
		}
	} else {
		// Cheeck if time has expired
		isInvoiceExpired := utils.GetInvoiceExpired(paymentRequest)
		// If the invoice has expired and it is not paid delete from the DB
		if isInvoiceExpired {
			h.db.DeleteInvoice(paymentRequest)
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
