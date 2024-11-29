package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/utils"
)

type ticketHandler struct {
	httpClient HttpClient
	db         db.Database
}

type TicketResponse struct {
	Success  bool     `json:"success"`
	TicketID string   `json:"ticket_id,omitempty"`
	Message  string   `json:"message"`
	Errors   []string `json:"errors,omitempty"`
}

func NewTicketHandler(httpClient HttpClient, database db.Database) *ticketHandler {
	return &ticketHandler{
		httpClient: httpClient,
		db:         database,
	}
}

func (th *ticketHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		http.Error(w, "UUID is required", http.StatusBadRequest)
		return
	}

	ticket, err := th.db.GetTicket(uuid)
	if err != nil {
		if err.Error() == "ticket not found" {
			http.Error(w, "Ticket not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get ticket: %v", err), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, ticket)
}

func (th *ticketHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("[ticket] no pubkey from auth")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	uuidStr := chi.URLParam(r, "uuid")
	if uuidStr == "" {
		http.Error(w, "UUID is required", http.StatusBadRequest)
		return
	}

	ticketUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var ticket db.Tickets
	if err := json.Unmarshal(body, &ticket); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	ticket.UUID = ticketUUID

	if ticket.Status != "" && !db.IsValidTicketStatus(ticket.Status) {
		http.Error(w, "Invalid ticket status", http.StatusBadRequest)
		return
	}

	updatedTicket, err := th.db.UpdateTicket(ticket)
	if err != nil {
		if err.Error() == "feature_uuid, phase_uuid, and name are required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to update ticket: %v", err), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedTicket)
}
func (th *ticketHandler) DeleteTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("[ticket] no pubkey from auth")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		http.Error(w, "UUID is required", http.StatusBadRequest)
		return
	}

	err := th.db.DeleteTicket(uuid)
	if err != nil {
		if err.Error() == "ticket not found" {
			http.Error(w, "Ticket not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to delete ticket: %v", err), http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (th *ticketHandler) PostTicketDataToStakwork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("[ticket] no pubkey from auth")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, TicketResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  []string{"Error reading request body"},
		})
		return
	}
	defer r.Body.Close()

	var ticket db.Tickets
	if err := json.Unmarshal(body, &ticket); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, TicketResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  []string{"Error parsing request body: " + err.Error()},
		})
		return
	}

	var validationErrors []string
	if ticket.UUID == uuid.Nil {
		validationErrors = append(validationErrors, "UUID is required")
	} else {
		if _, err := uuid.Parse(ticket.UUID.String()); err != nil {
			validationErrors = append(validationErrors, "Invalid UUID format")
		}
	}
	if ticket.FeatureUUID == "" {
		validationErrors = append(validationErrors, "FeatureUUID is required")
	}
	if ticket.PhaseUUID == "" {
		validationErrors = append(validationErrors, "PhaseUUID is required")
	}
	if ticket.Name == "" {
		validationErrors = append(validationErrors, "Name is required")
	}

	if len(validationErrors) > 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, TicketResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  validationErrors,
		})
		return
	}

	db := th.db
	feature := db.GetFeatureByUuid(ticket.FeatureUUID)
	if feature.Uuid == "" {
		utils.RespondWithJSON(w, http.StatusInternalServerError, TicketResponse{
			Success: false,
			Message: "Error retrieving feature details",
			Errors:  []string{"Feature not found with the provided UUID"},
		})
		return
	}

	productBrief, err := db.GetProductBrief(feature.WorkspaceUuid)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, TicketResponse{
			Success: false,
			Message: "Error retrieving product brief",
			Errors:  []string{err.Error()},
		})
		return
	}

	featureBrief, err := db.GetFeatureBrief(ticket.FeatureUUID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, TicketResponse{
			Success: false,
			Message: "Error retrieving feature brief",
			Errors:  []string{err.Error()},
		})
		return
	}

	host := os.Getenv("HOST")
	if host == "" {
		http.Error(w, "HOST environment variable not set", http.StatusInternalServerError)
		return
	}

	webhookURL := fmt.Sprintf("%s/bounty/ticket/review/", host)

	stakworkPayload := map[string]interface{}{
		"name":        "Hive Ticket Builder",
		"workflow_id": 37324,
		"workflow_params": map[string]interface{}{
			"set_var": map[string]interface{}{
				"attributes": map[string]interface{}{
					"vars": map[string]interface{}{
						"featureUUID":       ticket.FeatureUUID,
						"phaseUUID":         ticket.PhaseUUID,
						"ticketUUID":        ticket.UUID.String(),
						"ticketName":        ticket.Name,
						"ticketDescription": ticket.Description,
						"productBrief":      productBrief,
						"featureBrief":      featureBrief,
						"examples":          "",
						"webhook_url":       webhookURL,
					},
				},
			},
		},
	}

	stakworkPayloadJSON, err := json.Marshal(stakworkPayload)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, TicketResponse{
			Success: false,
			Message: "Error encoding payload",
			Errors:  []string{err.Error()},
		})
		return
	}

	apiKey := os.Getenv("SWWFKEY")
	if apiKey == "" {
		utils.RespondWithJSON(w, http.StatusInternalServerError, TicketResponse{
			Success: false,
			Message: "API key not set in environment",
		})
		return
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.stakwork.com/api/v1/projects", bytes.NewBuffer(stakworkPayloadJSON))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, TicketResponse{
			Success: false,
			Message: "Error creating request to Stakwork API",
			Errors:  []string{err.Error()},
		})
		return
	}

	req.Header.Set("Authorization", "Token token="+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, TicketResponse{
			Success: false,
			Message: "Error sending request to Stakwork API",
			Errors:  []string{err.Error()},
		})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, TicketResponse{
			Success: false,
			Message: "Error reading response from Stakwork API",
			Errors:  []string{err.Error()},
		})
		return
	}

	utils.RespondWithJSON(w, resp.StatusCode, TicketResponse{
		Success:  resp.StatusCode == http.StatusOK,
		Message:  string(respBody),
		TicketID: ticket.UUID.String(),
	})
}

func (th *ticketHandler) ProcessTicketReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("[ticket] no pubkey from auth")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var reviewReq utils.TicketReviewRequest
	if err := json.Unmarshal(body, &reviewReq); err != nil {
		log.Printf("Error parsing request JSON: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateTicketReviewRequest(&reviewReq); err != nil {
		log.Printf("Invalid request data: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticket, err := th.db.GetTicket(reviewReq.TicketUUID)
	if err != nil {
		log.Printf("Error fetching ticket: %v", err)
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	ticket.Description = reviewReq.TicketDescription

	updatedTicket, err := th.db.UpdateTicket(ticket)
	if err != nil {
		log.Printf("Error updating ticket: %v", err)
		http.Error(w, "Failed to update ticket", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully updated ticket %s", reviewReq.TicketUUID)
	utils.RespondWithJSON(w, http.StatusOK, updatedTicket)
}