package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
	"github.com/stakwork/sphinx-tribes/utils"
	"gorm.io/gorm"
)

type BountyTimingResponse struct {
	TotalWorkTimeSeconds    int        `json:"total_work_time_seconds"`
	TotalDurationSeconds    int        `json:"total_duration_seconds"`
	TotalAttempts           int        `json:"total_attempts"`
	FirstAssignedAt         *time.Time `json:"first_assigned_at"`
	LastPoWAt               *time.Time `json:"last_pow_at"`
	ClosedAt                *time.Time `json:"closed_at"`
	IsPaused                bool       `json:"is_paused"`
	LastPausedAt            *time.Time `json:"last_paused_at"`
	AccumulatedPauseSeconds int        `json:"accumulated_pause_seconds"`
}

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

type TimingError struct {
	Operation string `json:"operation"`
	Error     string `json:"error"`
}

func handleTimingError(w http.ResponseWriter, operation string, err error) {
	logger.Log.Error("[bounty_timing] %s failed: %v", operation, err)
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(TimingError{
		Operation: operation,
		Error:     err.Error(),
	})
}

// GetAllBounties godoc
//
//	@Summary		Get all bounties
//	@Description	Get a list of all bounties
//	@Tags			Bounties
//	@Success		200	{array}	db.Bounty
//	@Router			/gobounties/all [get]
func (h *bountyHandler) GetAllBounties(w http.ResponseWriter, r *http.Request) {
	bounties := h.db.GetAllBounties(r)
	var bountyResponse []db.BountyResponse = h.GenerateBountyResponse(bounties)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyResponse)
}

// GetBountyById godoc
//
//	@Summary		Get a bounty
//	@Description	Get a bounty by ID
//	@Tags			Bounties
//	@Param			id	path		string	true	"Bounty ID"
//	@Success		200	{object}	db.Bounty
//	@Router			/gobounties/id/{bountyId} [get]
func (h *bountyHandler) GetBountyById(w http.ResponseWriter, r *http.Request) {
	bountyId := chi.URLParam(r, "bountyId")
	if bountyId == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bounties, err := h.db.GetBountyById(bountyId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("[bounty] Error: %v", err)
	} else {
		var bountyResponse []db.BountyResponse = h.GenerateBountyResponse(bounties)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

// GetNextBountyByCreated godoc
//
//	@Summary		Get next bounty
//	@Description	Get next bounty by created date
//	@Tags			Bounties
//	@Param			created	path		string	true	"Created date"
//	@Success		200		{object}	db.Bounty
//	@Router			/gobounties/next/{created} [get]
func (h *bountyHandler) GetNextBountyByCreated(w http.ResponseWriter, r *http.Request) {
	bounties, err := h.db.GetNextBountyByCreated(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("[bounty] Error: %v", err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bounties)
	}
}

// GetPreviousBountyByCreated godoc
//
//	@Summary		Get previous bounty
//	@Description	Get previous bounty by created date
//	@Tags			Bounties
//	@Param			created	path		string	true	"Created date"
//	@Success		200		{object}	db.Bounty
//	@Router			/gobounties/previous/{created} [get]
func (h *bountyHandler) GetPreviousBountyByCreated(w http.ResponseWriter, r *http.Request) {
	bounties, err := h.db.GetPreviousBountyByCreated(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("[bounty] Error: %v", err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bounties)
	}
}

// GetWorkspaceNextBountyByCreated godoc
//
//	@Summary		Get next workspace bounty
//	@Description	Get next workspace bounty by created date
//	@Tags			Bounties
//	@Param			uuid	path		string	true	"Workspace UUID"
//	@Param			created	path		string	true	"Created date"
//	@Success		200		{object}	db.Bounty
//	@Router			/gobounties/org/next/{uuid}/{created} [get]
func (h *bountyHandler) GetWorkspaceNextBountyByCreated(w http.ResponseWriter, r *http.Request) {
	bounties, err := h.db.GetNextWorkspaceBountyByCreated(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("[bounty] Error: %v", err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bounties)
	}
}

// GetWorkspacePreviousBountyByCreated godoc
//
//	@Summary		Get previous workspace bounty
//	@Description	Get previous workspace bounty by created date
//	@Tags			Bounties
//	@Param			uuid	path		string	true	"Workspace UUID"
//	@Param			created	path		string	true	"Created date"
//	@Success		200		{object}	db.Bounty
//	@Router			/gobounties/org/previous/{uuid}/{created} [get]
func (h *bountyHandler) GetWorkspacePreviousBountyByCreated(w http.ResponseWriter, r *http.Request) {
	bounties, err := h.db.GetPreviousWorkspaceBountyByCreated(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("[bounty] Error: %v", err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bounties)
	}
}

// GetBountyIndexById godoc
//
//	@Summary		Get bounty index
//	@Description	Get bounty index by ID
//	@Tags			Bounties
//	@Param			id	path		string	true	"Bounty ID"
//	@Success		200	{object}	int
//	@Router			/gobounties/index/{bountyId} [get]
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

// GetBountyByCreated godoc
//
//	@Summary		Get bounty by created date
//	@Description	Get bounty by created date
//	@Tags			Bounties
//	@Param			created	path		string	true	"Created date"
//	@Success		200		{object}	db.Bounty
//	@Router			/gobounties/created/{created} [get]
func (h *bountyHandler) GetBountyByCreated(w http.ResponseWriter, r *http.Request) {
	created := chi.URLParam(r, "created")
	if created == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bounties, err := h.db.GetBountyDataByCreated(created)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("[bounty] Error: %v", err)
	} else {
		var bountyResponse []db.BountyResponse = h.GenerateBountyResponse(bounties)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

// GetUserBountyCount godoc
//
//	@Summary		Get user bounty count
//	@Description	Get user bounty count by person key and tab type
//	@Tags			Bounties
//	@Param			personKey	path		string	true	"Person Key"
//	@Param			tabType		path		string	true	"Tab Type"
//	@Success		200			{object}	int
//	@Router			/gobounties/count/{personKey}/{tabType} [get]
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

// GetBountyCount godoc
//
//	@Summary		Get bounty count
//	@Description	Get bounty count
//	@Tags			Bounties
//	@Success		200	{object}	int
//	@Router			/gobounties/count [get]
func GetBountyCount(w http.ResponseWriter, r *http.Request) {
	bountyCount := db.DB.GetBountiesCount(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyCount)
}

func (h *bountyHandler) GetPersonCreatedBounties(w http.ResponseWriter, r *http.Request) {
	bounties, err := h.db.GetCreatedBounties(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("[bounty] Error: %v", err)
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
		logger.Log.Error("[bounty] Error: %v", err)
	} else {
		var bountyResponse []db.BountyResponse = h.GenerateBountyResponse(bounties)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

func getContactKey(pubkey string) (*string, error) {
	url := fmt.Sprintf("%s/get_contact/%s", config.V2BotUrl, pubkey)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("x-admin-token", config.V2BotToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching contact: %v", err)
	}
	defer resp.Body.Close()

	var contactResp struct {
		ContactKey *string `json:"contact_key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&contactResp); err != nil {
		return nil, fmt.Errorf("error decoding contact response: %v", err)
	}

	return contactResp.ContactKey, nil
}

func ProcessWaitingNotifications() {
	notifications := db.DB.GetNotificationsByStatus("WAITING_KEY_EXCHANGE")

	for _, n := range notifications {
		contactKey, err := getContactKey(n.PubKey)
		if err != nil {
			logger.Log.Error("Error checking contact key for pubkey %s: %v", n.PubKey, err)
			db.DB.IncrementNotificationRetry(n.UUID)
			continue
		}

		if contactKey == nil {
			db.DB.IncrementNotificationRetry(n.UUID)
			continue
		}

		// Contact key is available, proceed with sending
		sendRespStatus := sendNotification(n.PubKey, n.Content)
		db.DB.UpdateNotificationStatus(n.UUID, sendRespStatus)
	}
}

func sendNotification(pubkey, content string) string {
	sendURL := fmt.Sprintf("%s/send", config.V2BotUrl)
	msgBody, _ := json.Marshal(map[string]interface{}{
		"dest":     pubkey,
		"amt_msat": 0,
		"content":  content,
		"is_tribe": false,
		"wait":     true,
	})

	req, err := http.NewRequest(http.MethodPost, sendURL, bytes.NewBuffer(msgBody))
	if err != nil {
		logger.Log.Error("Error creating send request: %v", err)
		return "FAILED"
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-admin-token", config.V2BotToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Error("Error sending notification: %v", err)
		return "FAILED"
	}
	defer resp.Body.Close()

	var sendResp struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&sendResp); err != nil {
		logger.Log.Error("Error decoding send response: %v", err)
		return "FAILED"
	}

	return sendResp.Status
}

func processNotification(pubkey, event, content, alias string, route_hint string) string {
	contactKey, err := getContactKey(pubkey)
	if err != nil {
		logger.Log.Error("Error checking contact key: %v", err)
		return "FAILED"
	}

	if contactKey == nil {

		contact_info := fmt.Sprintf("%s_%s", pubkey, route_hint)
		logger.Log.Info("Sending contact info: %v", contact_info)
		addContactURL := fmt.Sprintf("%s/add_contact", config.V2BotUrl)
		body, _ := json.Marshal(map[string]string{"contact_info": contact_info, "alias": alias})
		req, _ := http.NewRequest(http.MethodPost, addContactURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-admin-token", config.V2BotToken)
		http.DefaultClient.Do(req)

		contactKey, err = getContactKey(pubkey)
		if err != nil || contactKey == nil {
			db.DB.SaveNotification(pubkey, event, content, "WAITING_KEY_EXCHANGE")
			return "FAILED"
		}
	}

	return sendNotification(pubkey, content)
}

// CreateOrEditBounty godoc
//
//	@Summary		Create or edit a bounty
//	@Description	Create or edit a bounty
//	@Tags			Bounties
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			bounty	body		db.NewBounty	true	"Bounty object"
//	@Success		200		{object}	db.NewBounty
//	@Router			/gobounties [post]
func (h *bountyHandler) CreateOrEditBounty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	// return 401 if pubKeyFromAuth is empty
	if pubKeyFromAuth == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check if  use exists
	user := h.db.GetPersonByPubkey(pubKeyFromAuth)
	if user.OwnerPubKey == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bounty := db.NewBounty{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		logger.Log.Error("[bounty] Read error: %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = json.Unmarshal(body, &bounty)
	if err != nil {
		logger.Log.Error("[bounty] Unmarshal error: %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if bounty.ID != 0 {
		existingBounty := h.db.GetBounty(bounty.ID)
		if existingBounty.UnlockCode != nil {
			bounty.UnlockCode = existingBounty.UnlockCode
		}
	}

	if bounty.UnlockCode == nil {
		code := generateUnlockCode()
		bounty.UnlockCode = &code
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

		if bounty.ID != 0 {
			if err := h.db.StartBountyTiming(bounty.ID); err != nil {
				handleTimingError(w, "start_timing", err)
			}
		}

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
			logger.Log.Info("[bounty]: %v", msg)
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
					logger.Log.Info("[bounty]: %v", msg)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(msg)
					return
				}
			} else {
				msg := "Cannot edit another user's bounty"
				logger.Log.Info("[bounty]: %v", msg)
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
	existingBounty := h.db.GetBounty(bounty.ID)
	b, err := h.db.CreateOrEditBounty(bounty)
	if err != nil {
		logger.Log.Error("[bounty] Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if bounty.ID == 0 && bounty.Assignee != "" {
		if err := h.db.StartBountyTiming(b.ID); err != nil {
			handleTimingError(w, "start_timing", err)
		}
	}

	if bounty.Assignee != "" {
		msg := fmt.Sprintf("You have been assigned a new ticket: %s. %s/bounty/%d", bounty.Title, os.Getenv("HOST"), b.ID)
		assigneePubkey := bounty.Assignee
		if bounty.ID != 0 {
			if existingBounty.Assignee != "" && existingBounty.Assignee == bounty.Assignee {
				assigneePubkey = ""
			}
		}

		if assigneePubkey != "" {
			person := db.DB.GetPersonByPubkey(assigneePubkey)
			processNotification(assigneePubkey, "bounty_assigned", msg, person.OwnerAlias, person.OwnerRouteHint)
		}

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(b)
}

func generateUnlockCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// DeleteBounty godoc
//
//	@Summary		Delete a bounty
//	@Description	Delete a bounty by ID
//	@Tags			Bounties
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			pubkey	path		string	true	"PubKey"
//	@Param			created	path		string	true	"Created"
//	@Success		200		{object}	bool
//	@Router			/gobounties/{pubkey}/{created} [delete]
func (h *bountyHandler) DeleteBounty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	created := chi.URLParam(r, "created")
	pubkey := chi.URLParam(r, "pubkey")

	if pubkey == "" {
		logger.Log.Error("[bounty] no pubkey from route")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if created == "" {
		logger.Log.Error("[bounty] no created timestamp from route")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get bounty by created
	createdUint, _ := utils.ConvertStringToUint(created)
	createdBounty, err := h.db.GetBountyByCreated(createdUint)
	if err != nil {
		logger.Log.Error("[bounty] failed to delete bounty: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("failed to delete bounty")
		return
	}

	if createdBounty.ID == 0 {
		logger.Log.Error("[bounty] failed to delete bounty")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("failed to delete bounty")
		return
	}

	b, err := h.db.DeleteBounty(pubkey, created)
	if err != nil {
		logger.Log.Error("[bounty] failed to delete bounty: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("failed to delete bounty")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(b)
}

// UpdatePaymentStatus godoc
//
//	@Summary		Update payment status
//	@Description	Update payment status by created date
//	@Tags			Bounties - Payment
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			created	path		string	true	"Created"
//	@Success		200		{object}	db.NewBounty
//	@Router			/gobounties/paymentstatus/{created} [post]
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

// UpdateCompletedStatus godoc
//
//	@Summary		Update completed status
//	@Description	Update completed status by created date
//	@Tags			Bounties
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			created	path		string	true	"Created"
//	@Success		200		{object}	db.NewBounty
//	@Router			/gobounties/completedstatus/{created} [post]
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

// GetPaymentByBountyId godoc
//
//	@Summary		Get payment by bounty ID
//	@Description	Get payment by bounty ID
//	@Tags			Bounties - Payment
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			bountyId	path		string	true	"Bounty ID"
//	@Success		200			{object}	db.NewPaymentHistory
//	@Router			/gobounties/payment/{bountyId} [get]
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

		proofs := h.db.GetProofsByBountyID(bounty.ID)

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
				PhaseUuid:               bounty.PhaseUuid,
				FeatureUuid:             bounty.FeatureUuid,
				PhasePriority:           bounty.PhasePriority,
				ProofOfWorkCount:        bounty.ProofOfWorkCount,
				UnlockCode:              bounty.UnlockCode,
				AccessRestriction:       bounty.AccessRestriction,
				IsStakable:              bounty.IsStakable,
				StakeMin:                bounty.StakeMin,
				MaxStakers:              bounty.MaxStakers,
				CurrentStakers:          bounty.CurrentStakers,
				Stakes:                  bounty.Stakes,
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
			Pow: bounty.ProofOfWorkCount,
		}

		if len(proofs) > 0 {
			b.Proofs = proofs
		}

		bountyResponse = append(bountyResponse, b)
	}

	return bountyResponse
}

// MakeBountyPayment godoc
//
//	@Summary		Make a bounty payment
//	@Description	Make a bounty payment
//	@Tags			Bounties - Payment
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path		string	true	"Bounty ID"
//	@Success		200	{object}	db.NewBounty
//	@Router			/gobounties/pay/{id} [post]
func (h *bountyHandler) MakeBountyPayment(w http.ResponseWriter, r *http.Request) {
	h.m.Lock()

	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	idParam := chi.URLParam(r, "id")

	id, err := utils.ConvertStringToUint(idParam)
	if err != nil {
		logger.Log.Error("[bounty] could not parse id")
		w.WriteHeader(http.StatusForbidden)
		h.m.Unlock()
		return
	}

	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty] no pubkey from auth")
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
		logger.Log.Error("[bounty] Read body error: %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		h.m.Unlock()
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Log.Error("[bounty] Unmarshal error: %v", err)
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

		logger.Log.Info("IS V2 PAYMENT ====")

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
				logger.Log.Error("[Unmarshal failed]: %v", err)
				w.WriteHeader(http.StatusNotAcceptable)
				h.m.Unlock()
				return
			}

			log.Printf("[bounty] V2 Status After Making Bounty V2 Payment: amount: %d, pubkey: %s, route_hint: %s is : %s", amount, assignee.OwnerPubKey, assignee.OwnerRouteHint, v2KeysendRes.Status)

			// if the payment has a completed status
			if v2KeysendRes.Status == db.PaymentComplete {
				bounty.PaymentFailed = false
				bounty.PaymentPending = false
				bounty.Paid = true
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
				bounty.PaymentFailed = false
				bounty.PaymentPending = true
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
				msg["msg"] = "keysend_failed"
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
			logger.Log.Error("[bounty] Read body error: %v", err)
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
				logger.Log.Error("[bounty] Unmarshal error: %v", err)
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

// GetBountyPaymentStatus godoc
//
//	@Summary		Get bounty payment status
//	@Description	Get bounty payment status by ID
//	@Tags			Bounties - Payment
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path		string	true	"Bounty ID"
//	@Success		200	{object}	db.NewPaymentHistory
//	@Router			/gobounties/payment/status/{id} [get]
func (h *bountyHandler) GetBountyPaymentStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	idParam := chi.URLParam(r, "id")

	id, err := utils.ConvertStringToUint(idParam)
	if err != nil {
		logger.Log.Error("[bounty] could not parse id")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty] no pubkey from auth")
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

// UpdateBountyPaymentStatus godoc
//
//	@Summary		Update bounty payment status
//	@Description	Update bounty payment status by ID
//	@Tags			Bounties - Payment
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path		string	true	"Bounty ID"
//	@Success		200	{object}	db.NewPaymentHistory
//	@Router			/gobounties/payment/status/{id} [put]
func (h *bountyHandler) UpdateBountyPaymentStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	idParam := chi.URLParam(r, "id")

	id, err := utils.ConvertStringToUint(idParam)
	if err != nil {
		logger.Log.Error("[bounty] could not parse id")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty] no pubkey from auth")
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

	if bounty.PaymentFailed {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode("Bounty payment has failed, have to make payment again")
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

			h.db.UpdateBountyPaymentStatuses(bounty)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(msg)
			return
		} else if tagResult.Status == db.PaymentFailed {

			err = h.db.ProcessReversePayments(payment.ID)

			if err != nil {
				log.Printf("Could not reverse bounty payment : Bounty ID - %d, Payment ID - %d, Error - %s", bounty.ID, payment.ID, err)
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(msg)
			return
		} else if tagResult.Status == db.PaymentPending {
			if payment.PaymentStatus == db.PaymentPending {
				created := utils.ConvertTimeToTimestamp(payment.Created.String())

				now := time.Now()
				daysDiff := utils.GetDateDaysDifference(int64(created), &now)

				if daysDiff >= 7 {

					err = h.db.ProcessReversePayments(payment.ID)
					if err != nil {
						log.Printf("Could not reverse bounty payment after 7 days : Bounty ID - %d, Payment ID - %d, Error - %s", bounty.ID, payment.ID, err)
					}
				}
			}
		}
	}

	msg := map[string]string{
		"payment_status": db.PaymentNotFound,
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(msg)
}

// Todo: change back to BountyBudgetWithdraw
// BountyBudgetWithdraw godoc
//
//	@Summary		Withdraw bounty budget
//	@Description	Withdraw bounty budget
//	@Tags			Bounties - Payment
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			request	body		db.NewWithdrawBudgetRequest	true	"Withdraw Budget Request"
//	@Success		200		{object}	db.InvoicePaySuccess
//	@Router			/gobounties/budget/withdraw [post]
func (h *bountyHandler) BountyBudgetWithdraw(w http.ResponseWriter, r *http.Request) {
	h.m.Lock()

	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty] no pubkey from auth")
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

// GetInvoiceData godoc
//
//	@Summary		Get invoice data
//	@Description	Get invoice data by payment request
//	@Tags			Bounties - Payment
//	@Param			paymentRequest	path		string	true	"Payment Request"
//	@Success		200				{object}	db.InvoiceResult
//	@Router			/gobounties/invoice/{paymentRequest} [get]
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

// PollInvoice godoc
//
//	@Summary		Poll invoice
//	@Description	Poll invoice by payment request
//	@Tags			Bounties - Payment
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			paymentRequest	path		string	true	"Payment Request"
//	@Success		200				{object}	db.InvoiceResult
//	@Router			/poll/invoice/{paymentRequest} [get]
func (h *bountyHandler) PollInvoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	paymentRequest := chi.URLParam(r, "paymentRequest")

	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty] no pubkey from auth")
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

// GetFilterCount godoc
//
//	@Summary		Get filter count
//	@Description	Get filter count
//	@Tags			Bounties
//	@Success		200	{object}	db.FilterStatusCount
//	@Router			/gobounties/filter/count [get]
func (h *bountyHandler) GetFilterCount(w http.ResponseWriter, r *http.Request) {
	filterCount := h.db.GetFilterStatusCount()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filterCount)
}

// GetBountyCards godoc
//
//	@Summary		Get bounty cards
//	@Description	Get bounty cards
//	@Tags			Bounties
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{array}	db.BountyCard
//	@Router			/gobounties/bounty-cards [get]
func (h *bountyHandler) GetBountyCards(w http.ResponseWriter, r *http.Request) {
	workspaceUuid := r.URL.Query().Get("workspace_uuid")
	var bounties []db.NewBounty

	if workspaceUuid != "" {
		bounties = h.db.GetWorkspaceBountyCardsData(r)
	} else {
		bounties = h.db.GetAllBounties(r)
	}

	bountyCardResponse := h.GenerateBountyCardResponse(bounties)

	ticketCards, err := h.GenerateTicketCardResponse(workspaceUuid)
	if err != nil {
		logger.Log.Error("failed to generate ticket cards", "error", err)
	} else {
		bountyCardResponse = append(bountyCardResponse, ticketCards...)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyCardResponse)
}

func (h *bountyHandler) GenerateTicketCardResponse(workspaceUuid string) ([]db.BountyCard, error) {
	var ticketCards []db.BountyCard

	ticketGroups, err := h.db.GetAllTicketGroups(workspaceUuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket groups: %w", err)
	}

	for _, group := range ticketGroups {
		ticket, err := h.db.GetLatestTicketByGroup(group)
		if err != nil {
			logger.Log.Error("failed to get latest ticket", "group", group, "error", err)
			continue
		}

		feature := h.db.GetFeatureByUuid(ticket.FeatureUUID)
		if feature.WorkspaceUuid != workspaceUuid {
			continue
		}

		phase, _ := h.db.GetFeaturePhaseByUuid(ticket.FeatureUUID, ticket.PhaseUUID)
		workspace := h.db.GetWorkspaceByUuid(feature.WorkspaceUuid)
		bountyID := uint(ticket.UUID.ID())

		ticketCard := db.BountyCard{
			BountyID:     bountyID,
			TicketUUID:   &ticket.UUID,
			TicketGroup:  ticket.TicketGroup,
			Title:        ticket.Name,
			AssigneePic:  "",
			Assignee:     "",
			AssigneeName: "",
			Features:     feature,
			Phase:        phase,
			Workspace:    workspace,
			Status:       db.StatusDraft,
		}

		ticketCards = append(ticketCards, ticketCard)
	}

	return ticketCards, nil
}

func (h *bountyHandler) GenerateBountyCardResponse(bounties []db.NewBounty) []db.BountyCard {
	var bountyCardResponse []db.BountyCard

	for i := 0; i < len(bounties); i++ {
		bounty := bounties[i]

		var assigneePic, assigneeName, assigneePubkey string
		if bounty.Assignee != "" {
			assignee := h.db.GetPersonByPubkey(bounty.Assignee)
			if assignee.OwnerPubKey != "" {
				assigneePic = assignee.Img
				assigneeName = assignee.OwnerAlias
				assigneePubkey = assignee.OwnerPubKey
			}
		}

		workspace := h.db.GetWorkspaceByUuid(bounty.WorkspaceUuid)

		var phase db.FeaturePhase
		var feature db.WorkspaceFeatures

		if bounty.PhaseUuid != "" {
			phase, _ = h.db.GetPhaseByUuid(bounty.PhaseUuid)
		}

		if phase.FeatureUuid != "" {
			feature = h.db.GetFeatureByUuid(phase.FeatureUuid)
		}

		status := calculateBountyStatus(bounty)

		b := db.BountyCard{
			BountyID:     bounty.ID,
			Title:        bounty.Title,
			AssigneePic:  assigneePic,
			Assignee:     assigneePubkey,
			AssigneeName: assigneeName,
			Features:     feature,
			Phase:        phase,
			Workspace:    workspace,
			Status:       status,
		}

		bountyCardResponse = append(bountyCardResponse, b)
	}

	return bountyCardResponse
}

func calculateBountyStatus(bounty db.NewBounty) db.BountyStatus {
	if bounty.Paid {
		return db.StatusPaid
	}
	if bounty.Completed || bounty.PaymentPending {
		return db.StatusComplete
	}
	if bounty.Assignee == "" {
		return db.StatusTodo
	}
	if bounty.Assignee != "" && bounty.ProofOfWorkCount == 0 {
		return db.StatusInProgress
	}
	if bounty.Assignee != "" && bounty.ProofOfWorkCount > 0 {
		return db.StatusInReview
	}

	return db.StatusTodo
}

// AddProofOfWork godoc
//
//	@Summary		Add proof of work
//	@Description	Add proof of work to a bounty
//	@Tags			Bounties - Proof of Work
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id		path		string			true	"Bounty ID"
//	@Param			proof	body		db.ProofOfWork	true	"Proof of Work object"
//	@Success		201		{object}	db.ProofOfWork
//	@Router			/gobounties/{id}/proof [post]
func (h *bountyHandler) AddProofOfWork(w http.ResponseWriter, r *http.Request) {
	bountyID := chi.URLParam(r, "id")
	var proof db.ProofOfWork

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &proof)
	if err != nil || proof.Description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	proof.ID = uuid.New()
	proof.BountyID, _ = utils.ConvertStringToUint(bountyID)
	proof.CreatedAt = time.Now()
	proof.SubmittedAt = time.Now()

	if err := h.db.CreateProof(proof); err != nil {
		http.Error(w, "Failed to create proof", http.StatusInternalServerError)
		return
	}

	if err := h.db.PauseBountyTiming(proof.BountyID); err != nil {
		handleTimingError(w, "pause_timing", err)
	}

	if err := h.db.UpdateBountyTimingOnProof(proof.BountyID); err != nil {
		handleTimingError(w, "update_timing_on_proof", err)
	}

	if err := h.db.IncrementProofCount(proof.BountyID); err != nil {
		http.Error(w, "Failed to update bounty proof count", http.StatusInternalServerError)
		return
	}

	bounties, err := h.db.GetBountyById(bountyID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("[bounty] Error: %v", err)
	} else {
		var bountyResponse []db.BountyResponse = h.GenerateBountyResponse(bounties)

		ownerPubKey := bountyResponse[0].Owner.OwnerPubKey
		ownerAlias := bountyResponse[0].Owner.OwnerAlias
		ownerRouteHint := bountyResponse[0].Owner.OwnerRouteHint
		assineeAlias := bountyResponse[0].Assignee.OwnerAlias
		bountyTitle := bountyResponse[0].Bounty.Title
		bountyId := bountyResponse[0].Bounty.ID

		msg := fmt.Sprintf("%s has submitted PoW on Bounty %s/bounty/%d. %s", assineeAlias, os.Getenv("HOST"), bountyId, bountyTitle)

		if ownerPubKey != "" {
			processNotification(ownerPubKey, "bounty_assigned", msg, ownerAlias, ownerRouteHint)
		}

	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(proof)
}

// GetProofsByBounty godoc
//
//	@Summary		Get proofs by bounty
//	@Description	Get proofs by bounty ID
//	@Tags			Bounties - Proof of Work
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path	string	true	"Bounty ID"
//	@Success		200	{array}	db.ProofOfWork
//	@Router			/gobounties/{id}/proofs [get]
func (h *bountyHandler) GetProofsByBounty(w http.ResponseWriter, r *http.Request) {
	bountyID := chi.URLParam(r, "id")

	bountyUUID, err := utils.ConvertStringToUint(bountyID)
	if err != nil {
		http.Error(w, "Invalid bounty ID", http.StatusBadRequest)
		return
	}

	proofs := h.db.GetProofsByBountyID(bountyUUID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proofs)
}

// DeleteProof godoc
//
//	@Summary		Delete proof
//	@Description	Delete proof by ID
//	@Tags			Bounties - Proof of Work
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id		path	string	true	"Bounty ID"
//	@Param			proofId	path	string	true	"Proof ID"
//	@Success		204
//	@Router			/gobounties/{id}/proofs/{proofId} [delete]
func (h *bountyHandler) DeleteProof(w http.ResponseWriter, r *http.Request) {
	bountyID := chi.URLParam(r, "id")
	proofID := chi.URLParam(r, "proofId")

	if _, err := uuid.Parse(proofID); err != nil {
		http.Error(w, "Invalid proof ID", http.StatusBadRequest)
		return
	}

	bountyIDUint, err := utils.ConvertStringToUint(bountyID)
	if err != nil {
		http.Error(w, "Invalid bounty ID", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteProof(proofID); err != nil {
		http.Error(w, "Failed to delete proof", http.StatusInternalServerError)
		return
	}

	if err := h.db.DecrementProofCount(bountyIDUint); err != nil {
		http.Error(w, "Failed to update bounty proof count", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdateProofStatusResponse struct {
	Status db.ProofOfWorkStatus `json:"status"`
}

// UpdateProofStatus godoc
//
//	@Summary		Update the status of a proof of work
//	@Description	Update the status of a proof of work for a specific bounty. Valid statuses are "accepted", "rejected", and "change_requested".
//	@Tags			Bounties - Proof of Work
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Bounty ID"
//	@Param			proofId	path		string						true	"Proof of Work ID"
//	@Param			status	body		UpdateProofStatusResponse	true	"New status for the proof of work"
//	@Success		200		{object}	UpdateProofStatusResponse	"Status updated successfully"
//	@Failure		400		{string}	string						"Bad request: Invalid proof ID, bounty ID, or status"
//	@Failure		500		{string}	string						"Internal server error: Failed to update status"
//	@Router			/bounty/{id}/proof/{proofId}/status [put]
func (h *bountyHandler) UpdateProofStatus(w http.ResponseWriter, r *http.Request) {
	proofID := chi.URLParam(r, "proofId")
	bountyID := chi.URLParam(r, "id")

	var statusUpdate UpdateProofStatusResponse

	if _, err := uuid.Parse(proofID); err != nil {
		http.Error(w, "Invalid proof ID", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &statusUpdate); err != nil || !isValidProofStatus(statusUpdate.Status) {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	switch statusUpdate.Status {
	case db.RejectedStatus, db.ChangeRequestedStatus:
		id, err := utils.ConvertStringToUint(bountyID)
		if err != nil {
			http.Error(w, "Invalid bounty ID", http.StatusBadRequest)
			return
		}

		if err := h.db.ResumeBountyTiming(id); err != nil {
			logger.Log.Error(fmt.Sprintf("Failed to resume timing for bounty ID %d: %v", id, err))
		}

	case db.AcceptedStatus:
		id, err := utils.ConvertStringToUint(bountyID)
		if err != nil {
			http.Error(w, "Invalid bounty ID", http.StatusBadRequest)
			return
		}

		if err := h.db.CloseBountyTiming(id); err != nil {
			logger.Log.Error(fmt.Sprintf("Failed to close timing for bounty ID %d: %v", id, err))
		}
	}

	if err := h.db.UpdateProofStatus(proofID, statusUpdate.Status); err != nil {
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isValidProofStatus(status db.ProofOfWorkStatus) bool {
	switch status {
	case db.NewStatus, db.AcceptedStatus, db.RejectedStatus, db.ChangeRequestedStatus:
		return true
	}
	return false
}

// DeleteBountyAssignee godoc
//
//	@Summary		Delete a bounty assignee
//	@Description	Delete the assignee of a bounty. Only the bounty owner can perform this action.
//	@Tags			Bounties
//	@Accept			json
//	@Produce		json
//	@Param			request	body		db.DeleteBountyAssignee	true	"Request body containing owner_pubkey and created timestamp"
//	@Success		200		{boolean}	boolean					"Assignee deleted successfully"
//	@Failure		400		{string}	string					"Bad request: Missing or invalid parameters"
//	@Failure		406		{string}	string					"Not acceptable: Invalid request body"
//	@Failure		500		{string}	string					"Internal server error: Failed to delete assignee"
//	@Router			/bounty/assignee [delete]
func (h *bountyHandler) DeleteBountyAssignee(w http.ResponseWriter, r *http.Request) {
	invoice := db.DeleteBountyAssignee{}
	body, err := io.ReadAll(r.Body)
	var deletedAssignee bool

	r.Body.Close()

	err = json.Unmarshal(body, &invoice)

	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	owner_key := invoice.Owner_pubkey
	date := invoice.Created

	if owner_key == "" || date == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(false)
		return
	}

	createdUint, _ := strconv.ParseUint(date, 10, 32)
	b, err := h.db.GetBountyByCreated(uint(createdUint))

	if err == nil && b.OwnerID == owner_key {
		b.Assignee = ""
		b.AssignedHours = 0
		b.CommitmentFee = 0
		b.BountyExpires = ""

		h.db.UpdateBounty(b)

		if err := h.db.CloseBountyTiming(b.ID); err != nil {
			handleTimingError(w, "close_timing", err)
		}

		deletedAssignee = true
	} else {
		log.Printf("Could not delete bounty assignee")

		deletedAssignee = false

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(deletedAssignee)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deletedAssignee)

}

// GetBountyTimingStats godoc
//
//	@Summary		Get bounty timing stats
//	@Description	Get bounty timing stats by ID
//	@Tags			Bounties - Timing
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path		string	true	"Bounty ID"
//	@Success		200	{object}	BountyTimingResponse
//	@Router			/gobounties/{id}/timing [get]
func (h *bountyHandler) GetBountyTimingStats(w http.ResponseWriter, r *http.Request) {
	bountyID := chi.URLParam(r, "id")
	id, err := utils.ConvertStringToUint(bountyID)
	if err != nil {
		http.Error(w, "Invalid bounty ID", http.StatusBadRequest)
		return
	}

	timing, err := h.db.GetBountyTiming(id)
	if err != nil {
		http.Error(w, "Failed to get timing stats", http.StatusInternalServerError)
		return
	}

	response := BountyTimingResponse{
		TotalWorkTimeSeconds:    timing.TotalWorkTimeSeconds,
		TotalDurationSeconds:    timing.TotalDurationSeconds,
		TotalAttempts:           timing.TotalAttempts,
		FirstAssignedAt:         timing.FirstAssignedAt,
		LastPoWAt:               timing.LastPoWAt,
		ClosedAt:                timing.ClosedAt,
		IsPaused:                timing.IsPaused,
		LastPausedAt:            timing.LastPausedAt,
		AccumulatedPauseSeconds: timing.AccumulatedPauseSeconds,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// StartBountyTiming godoc
//
//	@Summary		Start bounty timing
//	@Description	Start bounty timing by ID
//	@Tags			Bounties - Timing
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path	string	true	"Bounty ID"
//	@Success		200
//	@Router			/gobounties/{id}/timing/start [post]
func (h *bountyHandler) StartBountyTiming(w http.ResponseWriter, r *http.Request) {
	bountyID := chi.URLParam(r, "id")

	id, err := utils.ConvertStringToUint(bountyID)
	if err != nil {
		http.Error(w, "Invalid bounty ID", http.StatusBadRequest)
		return
	}

	if err := h.db.StartBountyTiming(id); err != nil {
		http.Error(w, "Failed to start timing", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CloseBountyTiming godoc
//
//	@Summary		Close bounty timing
//	@Description	Close bounty timing by ID
//	@Tags			Bounties - Timing
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path	string	true	"Bounty ID"
//	@Success		200
//	@Router			/gobounties/{id}/timing/close [post]
func (h *bountyHandler) CloseBountyTiming(w http.ResponseWriter, r *http.Request) {
	bountyID := chi.URLParam(r, "id")
	id, err := utils.ConvertStringToUint(bountyID)
	if err != nil {
		http.Error(w, "Invalid bounty ID", http.StatusBadRequest)
		return
	}

	if err := h.db.CloseBountyTiming(id); err != nil {
		http.Error(w, "Failed to close timing", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetBountiesLeaderboard godoc
//
//	@Summary		Get bounties leaderboard
//	@Description	Get bounties leaderboard
//	@Tags			Bounties
//	@Success		200	{array}	db.LeaderData
//	@Router			/bounty/leaderboard [get]
func (h *bountyHandler) GetBountiesLeaderboard(w http.ResponseWriter, _ *http.Request) {
	leaderBoard := h.db.GetBountiesLeaderboard()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(leaderBoard)
}

func (h *bountyHandler) GetDailyEarnings(w http.ResponseWriter, _ *http.Request) {
	dailyEarnings := h.db.GetDailyEarnings()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dailyEarnings)
}

// GetAllFeaturedBounties godoc
//
//	@Summary		Get all featured bounties
//	@Description	Get all featured bounties
//	@Tags			Featured Bounties
//	@Success		200	{array}	db.FeaturedBounty
//	@Router			/gobounties/featured/all [get]
func (h *bountyHandler) GetAllFeaturedBounties(w http.ResponseWriter, r *http.Request) {
	bounties, err := h.db.GetAllFeaturedBounties()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bounties)
}

// CreateFeaturedBounty godoc
//
//	@Summary		Create a featured bounty
//	@Description	Create a featured bounty
//	@Tags			Featured Bounties
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			bounty	body		db.FeaturedBounty	true	"Featured Bounty object"
//	@Success		201		{object}	db.FeaturedBounty
//	@Router			/gobounties/featured/create [post]
func (h *bountyHandler) CreateFeaturedBounty(w http.ResponseWriter, r *http.Request) {
	var bounty db.FeaturedBounty
	if err := json.NewDecoder(r.Body).Decode(&bounty); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid payload"})
		return
	}

	if bounty.BountyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "BountyID is required"})
		return
	}

	if bounty.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "URL is required"})
		return
	}

	if _, err := url.ParseRequestURI(bounty.URL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid URL format"})
		return
	}

	existingBounty, _ := h.db.GetFeaturedBountyById(bounty.BountyID)
	if existingBounty.BountyID != "" {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "BountyID already exists"})
		return
	}

	bounty.AddedAt = time.Now().UnixMilli()
	bounty.CreatedAt = time.Now()
	bounty.UpdatedAt = time.Now()

	if err := h.db.CreateFeaturedBounty(bounty); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": map[string]interface{}{
			"bountyId": bounty.BountyID,
			"url":      bounty.URL,
			"addedAt":  bounty.AddedAt,
			"title":    bounty.Title,
		},
	})
}

// UpdateFeaturedBounty godoc
//
//	@Summary		Update a featured bounty
//	@Description	Update a featured bounty
//	@Tags			Featured Bounties
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			bounty	body		db.FeaturedBounty	true	"Featured Bounty object"
//	@Success		200		{object}	db.FeaturedBounty
//	@Router			/gobounties/featured/update [put]
func (h *bountyHandler) UpdateFeaturedBounty(w http.ResponseWriter, r *http.Request) {

	var bounty db.FeaturedBounty
	if err := json.NewDecoder(r.Body).Decode(&bounty); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid payload"})
		return
	}

	if bounty.BountyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "BountyID is required in the request body"})
		return
	}
	if bounty.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "URL is required"})
		return
	}

	if _, err := url.ParseRequestURI(bounty.URL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid URL format"})
		return
	}

	existingBounty, _ := h.db.GetFeaturedBountyById(bounty.BountyID)
	if existingBounty.BountyID == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bounty not found"})
		return
	}

	existingBounty.URL = bounty.URL
	existingBounty.Title = bounty.Title
	existingBounty.UpdatedAt = time.Now()

	if err := h.db.UpdateFeaturedBounty(bounty.BountyID, existingBounty); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": map[string]interface{}{
			"bountyId": existingBounty.BountyID,
			"url":      existingBounty.URL,
			"addedAt":  existingBounty.AddedAt,
			"title":    existingBounty.Title,
		},
	})
}

// DeleteFeaturedBounty godoc
//
//	@Summary		Delete a featured bounty
//	@Description	Delete a featured bounty by ID
//	@Tags			Featured Bounties
//	@Param			bountyId	path	string	true	"Bounty ID"
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		204
//	@Router			/gobounties/featured/delete/{bountyId} [delete]
func (h *bountyHandler) DeleteFeaturedBounty(w http.ResponseWriter, r *http.Request) {
	bountyID := chi.URLParam(r, "bountyId")

	if bountyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "BountyID is required in the URL"})
		return
	}

	existingBounty, _ := h.db.GetFeaturedBountyById(bountyID)
	if existingBounty.BountyID == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bounty not found"})
		return
	}

	if err := h.db.DeleteFeaturedBounty(bountyID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteBountyTiming godoc
//
//	@Summary		Delete bounty timing
//	@Description	Delete bounty timing by ID
//	@Tags			Bounties - Timing
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path	string	true	"Bounty ID"
//	@Success		204
//	@Router			/gobounties/{id}/timing [delete]
func (h *bountyHandler) DeleteBountyTiming(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bountyID := chi.URLParam(r, "id")
	id, err := utils.ConvertStringToUint(bountyID)
	if err != nil {
		http.Error(w, "Invalid bounty ID", http.StatusBadRequest)
		return
	}

	_, err = h.db.GetBountyTiming(id)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("No bounty timing found for bounty ID %d: %v", id, err))
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No timing record found"})
		return
	}

	if err := h.db.DeleteBountyTiming(id); err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to delete bounty timing for bounty ID %d: %v", id, err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete bounty timing"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetBountiesByWorkspaceTime godoc
//
//	@Summary		Get bounties by workspace time
//	@Description	Get bounties by workspace time range
//	@Tags			Bounties
//	@Param			workspaceId	path		string	true	"Workspace ID"
//	@Param			daysStart	path		string	true	"Days Start"
//	@Param			daysEnd		path		string	true	"Days End"
//	@Success		200			{object}	db.NodeListResponse
//	@Router			/gobounties/org/timerange/{workspaceId}/{daysStart}/{daysEnd} [get]
func (h *bountyHandler) GetBountiesByWorkspaceTime(w http.ResponseWriter, r *http.Request) {

	workspaceId := chi.URLParam(r, "workspaceId")
	daysStartStr := chi.URLParam(r, "daysStart")
	daysEndStr := chi.URLParam(r, "daysEnd")

	daysStart, err := strconv.Atoi(daysStartStr)
	if err != nil {
		http.Error(w, "Invalid daysStart parameter", http.StatusBadRequest)
		return
	}

	daysEnd, err := strconv.Atoi(daysEndStr)
	if err != nil {
		http.Error(w, "Invalid daysEnd parameter", http.StatusBadRequest)
		return
	}

	if workspaceId == "" {
		http.Error(w, "Workspace ID is required", http.StatusBadRequest)
		return
	}

	now := time.Now()
	var endDate time.Time

	if daysStart == 0 {
		endDate = now.AddDate(0, 0, 1)
	} else {
		endDate = now.AddDate(0, 0, -daysStart)
	}

	startDate := now.AddDate(0, 0, -daysEnd)

	bounties, err := h.db.GetBountiesByWorkspaceAndTimeRange(workspaceId, startDate, endDate)
	if err != nil {
		logger.Log.Error("[bounty] Error retrieving bounties: %v", err)
		http.Error(w, "Error retrieving bounties", http.StatusInternalServerError)
		return
	}

	nodes := make([]db.Node, len(bounties))
	for i, bounty := range bounties {
		nodes[i] = db.Node{
			NodeType: "Bounty",
			NodeData: db.NodeData{
				BountyID:    bounty.ID,
				Title:       bounty.Title,
				Description: bounty.Description,
			},
		}
	}

	response := db.NodeListResponse{
		NodeList: nodes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateBountyStake godoc
//
//	@Summary		Create a bounty stake
//	@Description	Create a new stake for a bounty
//	@Tags			Bounties - Stakes
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			stake	body		db.BountyStake	true	"Stake object"
//	@Success		201		{object}	db.BountyStake
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/gobounties/stake [post]
func (h *bountyHandler) CreateBountyStake(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty_stake] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	
	var stake db.BountyStake
	if err := json.NewDecoder(r.Body).Decode(&stake); err != nil {
		logger.Log.Error("[bounty_stake] invalid request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	
	stake.HunterPubKey = pubKeyFromAuth
	
	createdStake, err := h.db.CreateBountyStake(stake)
	if err != nil {
		logger.Log.Error("[bounty_stake] failed to create stake: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdStake)
}

// GetAllBountyStakes godoc
//
//	@Summary		Get all bounty stakes
//	@Description	Get a list of all bounty stakes
//	@Tags			Bounties - Stakes
//	@Produce		json
//	@Success		200	{array}		db.BountyStake
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/gobounties/stakes [get]
func (h *bountyHandler) GetAllBountyStakes(w http.ResponseWriter, r *http.Request) {
	stakes, err := h.db.GetAllBountyStakes()
	if err != nil {
		logger.Log.Error("[bounty_stake] failed to get all stakes: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve stakes"})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stakes)
}

// GetBountyStakesByBountyID godoc
//
//	@Summary		Get bounty stakes by bounty ID
//	@Description	Get stakes associated with a specific bounty
//	@Tags			Bounties - Stakes
//	@Produce		json
//	@Param			bountyId	path		string	true	"Bounty ID"
//	@Success		200			{array}		db.BountyStake
//	@Failure		400			{string}	string	"Bad request"
//	@Failure		500			{string}	string	"Internal server error"
//	@Router			/gobounties/stake/bounty/{bountyId} [get]
func (h *bountyHandler) GetBountyStakesByBountyID(w http.ResponseWriter, r *http.Request) {
	bountyIDStr := chi.URLParam(r, "bountyId")
	bountyID, err := utils.ConvertStringToUint(bountyIDStr)
	if err != nil {
		logger.Log.Error("[bounty_stake] invalid bounty ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid bounty ID"})
		return
	}
	
	stakes, err := h.db.GetBountyStakesByBountyID(bountyID)
	if err != nil {
		logger.Log.Error("[bounty_stake] failed to get stakes by bounty ID: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve stakes"})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stakes)
}

// GetBountyStakeByID godoc
//
//	@Summary		Get bounty stake by ID
//	@Description	Get a specific stake by its ID
//	@Tags			Bounties - Stakes
//	@Produce		json
//	@Param			id	path		string	true	"Stake ID"
//	@Success		200	{object}	db.BountyStake
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		404	{string}	string	"Not found"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/gobounties/stake/{id} [get]
func (h *bountyHandler) GetBountyStakeByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Log.Error("[bounty_stake] invalid stake ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid stake ID"})
		return
	}
	
	stake, err := h.db.GetBountyStakeByID(id)
	if err != nil {
		logger.Log.Error("[bounty_stake] failed to get stake by ID: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Stake not found"})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stake)
}

// GetBountyStakesByHunterPubKey godoc
//
//	@Summary		Get bounty stakes by hunter public key
//	@Description	Get stakes associated with a specific hunter
//	@Tags			Bounties - Stakes
//	@Produce		json
//	@Param			hunterPubKey	path		string	true	"Hunter Public Key"
//	@Success		200				{array}		db.BountyStake
//	@Failure		500				{string}	string	"Internal server error"
//	@Router			/gobounties/stake/hunter/{hunterPubKey} [get]
func (h *bountyHandler) GetBountyStakesByHunterPubKey(w http.ResponseWriter, r *http.Request) {
	hunterPubKey := chi.URLParam(r, "hunterPubKey")
	
	stakes, err := h.db.GetBountyStakesByHunterPubKey(hunterPubKey)
	if err != nil {
		logger.Log.Error("[bounty_stake] failed to get stakes by hunter pubkey: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve stakes"})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stakes)
}

// UpdateBountyStake godoc
//
//	@Summary		Update a bounty stake
//	@Description	Update a specific stake by its ID
//	@Tags			Bounties - Stakes
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id		path		string					true	"Stake ID"
//	@Param			updates	body		map[string]interface{}	true	"Fields to update"
//	@Success		200		{object}	db.BountyStake
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		404		{string}	string	"Not found"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/gobounties/stake/{id} [put]
func (h *bountyHandler) UpdateBountyStake(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty_stake] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Log.Error("[bounty_stake] invalid stake ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid stake ID"})
		return
	}
	
	existingStake, err := h.db.GetBountyStakeByID(id)
	if err != nil {
		logger.Log.Error("[bounty_stake] stake not found: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Stake not found"})
		return
	}
	
	bounty := h.db.GetBounty(existingStake.BountyID)
	if existingStake.HunterPubKey != pubKeyFromAuth && bounty.OwnerID != pubKeyFromAuth {
		logger.Log.Error("[bounty_stake] unauthorized update attempt")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to update this stake"})
		return
	}
	
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		logger.Log.Error("[bounty_stake] invalid request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	
	updatedStake, err := h.db.UpdateBountyStake(id, updates)
	if err != nil {
		logger.Log.Error("[bounty_stake] failed to update stake: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedStake)
}

// DeleteBountyStake godoc
//
//	@Summary		Delete a bounty stake
//	@Description	Delete a specific stake by its ID
//	@Tags			Bounties - Stakes
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path		string	true	"Stake ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		404	{string}	string	"Not found"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/gobounties/stake/{id} [delete]
func (h *bountyHandler) DeleteBountyStake(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty_stake] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Log.Error("[bounty_stake] invalid stake ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid stake ID"})
		return
	}
	
	existingStake, err := h.db.GetBountyStakeByID(id)
	if err != nil {
		logger.Log.Error("[bounty_stake] stake not found: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Stake not found"})
		return
	}
	
	bounty := h.db.GetBounty(existingStake.BountyID)
	if existingStake.HunterPubKey != pubKeyFromAuth && bounty.OwnerID != pubKeyFromAuth {
		logger.Log.Error("[bounty_stake] unauthorized delete attempt")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to delete this stake"})
		return
	}
	
	err = h.db.DeleteBountyStake(id)
	if err != nil {
		logger.Log.Error("[bounty_stake] failed to delete stake: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Stake deleted successfully"})
}

// CreateBountyStakeProcess godoc
//
//	@Summary		Create a bounty stake process
//	@Description	Create a new stake process for a bounty
//	@Tags			Bounties - Stakes
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			process	body		db.BountyStakeProcess	true	"Stake process object"
//	@Success		201		{object}	db.BountyStakeProcess
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/gobounties/stake/stakeprocessing [post]
func (h *bountyHandler) CreateBountyStakeProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty_stake_process] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	
	var process db.BountyStakeProcess
	if err := json.NewDecoder(r.Body).Decode(&process); err != nil {
		logger.Log.Error("[bounty_stake_process] invalid request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	
	process.HunterPubKey = pubKeyFromAuth
	
	process.Status = db.StakeProcessStatusNew
	
	createdProcess, err := h.db.CreateBountyStakeProcess(&process)
	if err != nil {
		logger.Log.Error("[bounty_stake_process] failed to create stake process: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProcess)
}

// GetAllBountyStakeProcesses godoc
//
//	@Summary		Get all bounty stake processes
//	@Description	Get all stake processes sorted by newest first
//	@Tags			Bounties - Stakes
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Success		200	{array}		db.BountyStakeProcess
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/gobounties/stake/stakeprocessing [get]
func (h *bountyHandler) GetAllBountyStakeProcesses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty_stake_process] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	
	processes, err := h.db.GetAllBountyStakeProcesses()
	if err != nil {
		logger.Log.Error("[bounty_stake_process] failed to get stake processes: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve stake processes"})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(processes)
}

// GetBountyStakeProcessByID godoc
//
//	@Summary		Get bounty stake process by ID
//	@Description	Get a specific stake process by its ID
//	@Tags			Bounties - Stakes
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path		string	true	"Stake Process ID"
//	@Success		200	{object}	db.BountyStakeProcess
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		404	{string}	string	"Not found"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/gobounties/stake/stakeprocessing/{id} [get]
func (h *bountyHandler) GetBountyStakeProcessByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty_stake_process] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Log.Error("[bounty_stake_process] invalid process ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid process ID"})
		return
	}
	
	process, err := h.db.GetBountyStakeProcessByID(id)
	if err != nil {
		logger.Log.Error("[bounty_stake_process] failed to get process by ID: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Stake process not found"})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(process)
}

// UpdateBountyStakeProcess godoc
//
//	@Summary		Update a bounty stake process
//	@Description	Update a specific stake process by its ID
//	@Tags			Bounties - Stakes
//	@Accept			json
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id		path		string	true	"Stake Process ID"
//	@Param			updates	body		object	true	"Fields to update"
//	@Success		200		{object}	db.BountyStakeProcess
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		404		{string}	string	"Not found"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/gobounties/stake/stakeprocessing/{id} [put]
func (h *bountyHandler) UpdateBountyStakeProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty_stake_process] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Log.Error("[bounty_stake_process] invalid process ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid process ID"})
		return
	}
	
	existingProcess, err := h.db.GetBountyStakeProcessByID(id)
	if err != nil {
		logger.Log.Error("[bounty_stake_process] process not found: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Stake process not found"})
		return
	}
	
	bounty := h.db.GetBounty(existingProcess.BountyID)
	if existingProcess.HunterPubKey != pubKeyFromAuth && bounty.OwnerID != pubKeyFromAuth {
		logger.Log.Error("[bounty_stake_process] unauthorized update attempt")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to update this stake process"})
		return
	}
	
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		logger.Log.Error("[bounty_stake_process] invalid request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	
	updatedProcess, err := h.db.UpdateBountyStakeProcess(id, updates)
	if err != nil {
		logger.Log.Error("[bounty_stake_process] failed to update process: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedProcess)
}

// DeleteBountyStakeProcess godoc
//
//	@Summary		Delete a bounty stake process
//	@Description	Delete a specific stake process by its ID
//	@Tags			Bounties - Stakes
//	@Produce		json
//	@Security		PubKeyContextAuth
//	@Param			id	path		string	true	"Stake Process ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		404	{string}	string	"Not found"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/gobounties/stake/stakeprocessing/{id} [delete]
func (h *bountyHandler) DeleteBountyStakeProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	
	if pubKeyFromAuth == "" {
		logger.Log.Error("[bounty_stake_process] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Log.Error("[bounty_stake_process] invalid process ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid process ID"})
		return
	}
	
	existingProcess, err := h.db.GetBountyStakeProcessByID(id)
	if err != nil {
		logger.Log.Error("[bounty_stake_process] process not found: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Stake process not found"})
		return
	}
	
	bounty := h.db.GetBounty(existingProcess.BountyID)
	if existingProcess.HunterPubKey != pubKeyFromAuth && bounty.OwnerID != pubKeyFromAuth {
		logger.Log.Error("[bounty_stake_process] unauthorized delete attempt")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "You are not authorized to delete this stake process"})
		return
	}
	
	err = h.db.DeleteBountyStakeProcess(id)
	if err != nil {
		logger.Log.Error("[bounty_stake_process] failed to delete process: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Stake process deleted successfully"})
}

func (h *bountyHandler) GetBountyStakeProcessesByBountyID(w http.ResponseWriter, r *http.Request) {
	bountyIDStr := chi.URLParam(r, "bountyId")
	bountyID, err := strconv.ParseUint(bountyIDStr, 10, 32)
	if err != nil {
		logger.Log.Error("[bounty_stake_process] invalid bounty ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid bounty ID format"})
		return
	}

	processes, err := h.db.GetBountyStakeProcessesByBountyID(uint(bountyID))
	if err != nil {
		logger.Log.Error("[bounty_stake_process] failed to get stake processes: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(processes)
}
