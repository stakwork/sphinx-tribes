package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
	"github.com/stakwork/sphinx-tribes/utils"
	"github.com/stakwork/sphinx-tribes/websocket"
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

type StakworkResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ProjectID int64 `json:"project_id"`
	} `json:"data"`
}

type CreateBountyResponse struct {
	BountyID uint   `json:"bounty_id"`
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
}

func NewTicketHandler(httpClient HttpClient, database db.Database) *ticketHandler {
	return &ticketHandler{
		httpClient: httpClient,
		db:         database,
	}
}

type UpdateTicketRequest struct {
	Metadata struct {
		Source string `json:"source"`
		ID     string `json:"id"`
	} `json:"metadata"`
	Ticket *db.Tickets `json:"ticket"`
}

func (th *ticketHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "UUID is required"})
		return
	}

	ticket, err := th.db.GetTicket(uuid)
	if err != nil {
		if err.Error() == "ticket not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Ticket not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to get ticket: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ticket)
}

func (th *ticketHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[ticket] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	uuidStr := chi.URLParam(r, "uuid")
	if uuidStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "UUID is required"})
		return
	}

	ticketUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid UUID format"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error reading request body"})
		return
	}
	defer r.Body.Close()

	var updateRequest UpdateTicketRequest
	if err := json.Unmarshal(body, &updateRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error parsing request body"})
		return
	}

	if updateRequest.Ticket == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ticket data is required"})
		return
	}

	existingTicket, err := th.db.GetTicket(ticketUUID.String())
	var newTicket db.Tickets

	if err != nil {
		newTicket = db.Tickets{
			UUID:        updateRequest.Ticket.UUID,
			TicketGroup: func() *uuid.UUID { id := uuid.New(); return &id }(),
			FeatureUUID: updateRequest.Ticket.FeatureUUID,
			PhaseUUID:   updateRequest.Ticket.PhaseUUID,
			Name:        updateRequest.Ticket.Name,
			Sequence:    updateRequest.Ticket.Sequence,
			Dependency:  updateRequest.Ticket.Dependency,
			Description: updateRequest.Ticket.Description,
			Status:      updateRequest.Ticket.Status,
			Version:     1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	} else {

		newTicket = db.Tickets{
			UUID:        uuid.New(),
			TicketGroup: existingTicket.TicketGroup,
			FeatureUUID: updateRequest.Ticket.FeatureUUID,
			PhaseUUID:   updateRequest.Ticket.PhaseUUID,
			Name:        updateRequest.Ticket.Name,
			Sequence:    updateRequest.Ticket.Sequence,
			Dependency:  updateRequest.Ticket.Dependency,
			Description: updateRequest.Ticket.Description,
			Status:      updateRequest.Ticket.Status,
			Version:     existingTicket.Version + 1,
			Author:      updateRequest.Ticket.Author,
			AuthorID:    updateRequest.Ticket.AuthorID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	}

	createdTicket, err := th.db.CreateOrEditTicket(&newTicket)
	if err != nil {
		if err.Error() == "feature_uuid, phase_uuid, and name are required" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to update ticket: %v", err)})
		return
	}

	if updateRequest.Metadata.Source == "websocket" && updateRequest.Metadata.ID != "" {
		ticketMsg := websocket.TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: updateRequest.Metadata.ID,
			Message:         fmt.Sprintf("Hive has successfully updated your ticket %s", updateRequest.Ticket.Name),
			Action:          "message",
			TicketDetails: websocket.TicketData{
				FeatureUUID:       updateRequest.Ticket.FeatureUUID,
				PhaseUUID:         updateRequest.Ticket.PhaseUUID,
				TicketUUID:        updateRequest.Ticket.UUID.String(),
				TicketDescription: updateRequest.Ticket.Description,
			},
		}

		if err := websocket.WebsocketPool.SendTicketMessage(ticketMsg); err != nil {
			log.Printf("Failed to send websocket message: %v", err)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"ticket":          createdTicket,
				"websocket_error": err.Error(),
			})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(createdTicket)
}

func (th *ticketHandler) DeleteTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[ticket] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "UUID is required"})
		return
	}

	err := th.db.DeleteTicket(uuid)
	if err != nil {
		if err.Error() == "ticket not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Ticket not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to delete ticket: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Ticket deleted successfully"})
}

func (th *ticketHandler) PostTicketDataToStakwork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		logger.Log.Info("[ticket] no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TicketResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  []string{"Error reading request body"},
		})
		return
	}
	defer r.Body.Close()

	var ticketRequest UpdateTicketRequest
	if err := json.Unmarshal(body, &ticketRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TicketResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  []string{"Error parsing request body: " + err.Error()},
		})
		return
	}

	if ticketRequest.Ticket == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TicketResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  []string{"Ticket data is required"},
		})
		return
	}

	ticket := ticketRequest.Ticket
	var validationErrors []string
	if ticket.UUID == uuid.Nil {
		validationErrors = append(validationErrors, "UUID is required")
	} else if _, err := uuid.Parse(ticket.UUID.String()); err != nil {
		validationErrors = append(validationErrors, "Invalid UUID format")
	}

	if len(validationErrors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TicketResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  validationErrors,
		})
		return
	}

	var (
		productBrief, featureBrief, codeGraphURL string
		feature                                  db.WorkspaceFeatures
	)

	if ticket.FeatureUUID != "" {
		feature = th.db.GetFeatureByUuid(ticket.FeatureUUID)
		if feature.Uuid == "" {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(TicketResponse{
				Success: false,
				Message: "Error retrieving feature details",
				Errors:  []string{"Feature not found with the provided UUID"},
			})
			return
		}

		var err error
		productBrief, err = th.db.GetProductBrief(feature.WorkspaceUuid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(TicketResponse{
				Success: false,
				Message: "Error retrieving product brief",
				Errors:  []string{err.Error()},
			})
			return
		}

		featureBrief, err = th.db.GetFeatureBrief(ticket.FeatureUUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(TicketResponse{
				Success: false,
				Message: "Error retrieving feature brief",
				Errors:  []string{err.Error()},
			})
			return
		}
	}

	host := os.Getenv("HOST")
	if host == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TicketResponse{
			Success: false,
			Message: "HOST environment variable not set",
		})
		return
	}

	webhookURL := fmt.Sprintf("%s/bounties/ticket/review", host)

	var schematicURL string
	if feature.WorkspaceUuid != "" {

		workspace := th.db.GetWorkspaceByUuid(feature.WorkspaceUuid)

		if workspace.Uuid == "" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(TicketResponse{
				Success: false,
				Message: "Workspace not found",
			})
			return
		}

		schematicURL = workspace.SchematicUrl

		codeGraph, err := th.db.GetCodeGraphByUUID(feature.WorkspaceUuid)
		if err == nil {
			codeGraphURL = codeGraph.Url
		} else {
			codeGraphURL = ""
		}

	}

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
						"sourceWebsocket":   ticketRequest.Metadata.ID,
						"webhook_url":       webhookURL,
						"phaseSchematic":    schematicURL,
						"codeGraph":         codeGraphURL,
					},
				},
			},
		},
	}

	stakworkPayloadJSON, err := json.Marshal(stakworkPayload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TicketResponse{
			Success: false,
			Message: "Error encoding payload",
			Errors:  []string{err.Error()},
		})
		return
	}

	apiKey := os.Getenv("SWWFKEY")
	if apiKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TicketResponse{
			Success: false,
			Message: "API key not set in environment",
		})
		return
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.stakwork.com/api/v1/projects", bytes.NewBuffer(stakworkPayloadJSON))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TicketResponse{
			Success: false,
			Message: "Error creating request",
			Errors:  []string{err.Error()},
		})
		return
	}

	req.Header.Set("Authorization", "Token token="+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := th.httpClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TicketResponse{
			Success: false,
			Message: "Error sending request to Stakwork",
			Errors:  []string{err.Error()},
		})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TicketResponse{
			Success: false,
			Message: "Error reading response from Stakwork",
			Errors:  []string{err.Error()},
		})
		return
	}

	var stakworkResp StakworkResponse
	if err := json.Unmarshal(respBody, &stakworkResp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TicketResponse{
			Success: false,
			Message: "Error parsing Stakwork response",
			Errors:  []string{err.Error()},
		})
		return
	}

	if resp.StatusCode != http.StatusOK || !stakworkResp.Success {
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(TicketResponse{
			Success: false,
			Message: string(respBody),
			Errors:  []string{fmt.Sprintf("Stakwork API returned status code: %d", resp.StatusCode)},
		})
		return
	}

	if ticketRequest.Metadata.Source == "websocket" && ticketRequest.Metadata.ID != "" {
		ticketMsg := websocket.TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: ticketRequest.Metadata.ID,
			Message:         fmt.Sprintf("I have your updates and I'm rewriting ticket %s now", ticketRequest.Ticket.Name),
			Action:          "message",
			TicketDetails: websocket.TicketData{
				FeatureUUID:       ticketRequest.Ticket.FeatureUUID,
				PhaseUUID:         ticketRequest.Ticket.PhaseUUID,
				TicketUUID:        ticketRequest.Ticket.UUID.String(),
				TicketDescription: ticketRequest.Ticket.Description,
			},
		}

		if err := websocket.WebsocketPool.SendTicketMessage(ticketMsg); err != nil {
			log.Printf("Failed to send websocket message: %v", err)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"ticket":          ticketRequest,
				"websocket_error": err.Error(),
			})
			return
		}

		projectMsg := websocket.TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: ticketRequest.Metadata.ID,
			Message:         fmt.Sprintf("https://jobs.stakwork.com/admin/projects/%d", stakworkResp.Data.ProjectID),
			Action:          "swrun",
			TicketDetails: websocket.TicketData{
				FeatureUUID:       ticketRequest.Ticket.FeatureUUID,
				PhaseUUID:         ticketRequest.Ticket.PhaseUUID,
				TicketUUID:        ticketRequest.Ticket.UUID.String(),
				TicketDescription: ticketRequest.Ticket.Description,
			},
		}

		if err := websocket.WebsocketPool.SendTicketMessage(projectMsg); err != nil {
			log.Printf("Failed to send project ID websocket message: %v", err)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"ticket":          ticketRequest,
				"websocket_error": err.Error(),
			})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TicketResponse{
		Success:  true,
		Message:  string(respBody),
		TicketID: ticket.UUID.String(),
	})
}

func (th *ticketHandler) ProcessTicketReview(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error reading request body"})
		return
	}
	defer r.Body.Close()

	var reviewReq utils.TicketReviewRequest
	if err := json.Unmarshal(body, &reviewReq); err != nil {
		log.Printf("Error parsing request JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error parsing request body"})
		return
	}

	if err := utils.ValidateTicketReviewRequest(&reviewReq); err != nil {
		log.Printf("Invalid request data: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	existingTicket, err := th.db.GetTicket(reviewReq.Value.TicketUUID)
	if err != nil {
		log.Printf("Error fetching ticket: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ticket not found"})
		return
	}

	newTicket := db.Tickets{
		UUID:        uuid.New(),
		TicketGroup: existingTicket.TicketGroup,
		FeatureUUID: existingTicket.FeatureUUID,
		PhaseUUID:   existingTicket.PhaseUUID,
		Name:        existingTicket.Name,
		Sequence:    existingTicket.Sequence,
		Dependency:  existingTicket.Dependency,
		Description: reviewReq.Value.TicketDescription,
		Status:      existingTicket.Status,
		Version:     existingTicket.Version + 1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdTicket, err := th.db.CreateOrEditTicket(&newTicket)
	if err != nil {
		log.Printf("Error creating new ticket: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create new ticket"})
		return
	}

	ticketMsg := websocket.TicketMessage{
		BroadcastType:   "direct",
		SourceSessionID: reviewReq.SourceWebsocket,
		Message:         fmt.Sprintf("Successfully created new version for ticket %s", createdTicket.UUID.String()),
		Action:          "process",
		TicketDetails: websocket.TicketData{
			FeatureUUID:       reviewReq.Value.FeatureUUID,
			PhaseUUID:         reviewReq.Value.PhaseUUID,
			TicketUUID:        createdTicket.UUID.String(),
			TicketDescription: reviewReq.Value.TicketDescription,
		},
	}

	if err := websocket.WebsocketPool.SendTicketMessage(ticketMsg); err != nil {
		log.Printf("Failed to send websocket message: %v", err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ticket":          createdTicket,
			"websocket_error": err.Error(),
		})
		return
	}

	log.Printf("Successfully created new ticket version %s", createdTicket.UUID.String())
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(createdTicket)
}

func (th *ticketHandler) GetTicketsByPhaseUUID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	featureUUID := chi.URLParam(r, "feature_uuid")
	phaseUUID := chi.URLParam(r, "phase_uuid")

	if featureUUID == "" {
		log.Println("feature uuid is missing")
		http.Error(w, "Missing feature uuid", http.StatusBadRequest)
		return
	}

	if phaseUUID == "" {
		log.Println("phase uuid is missing")
		http.Error(w, "Missing phase uuid", http.StatusBadRequest)
		return
	}

	feature := th.db.GetFeatureByUuid(featureUUID)
	if feature.Uuid == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "feature not found"})
		return
	}

	_, err := th.db.GetFeaturePhaseByUuid(featureUUID, phaseUUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Phase not found"})
		return
	}

	tickets, err := th.db.GetTicketsByPhaseUUID(featureUUID, phaseUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tickets)
}

func (th *ticketHandler) TicketToBounty(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	if pubKeyFromAuth == "" {
		logger.Log.Info("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ticketUUID := chi.URLParam(r, "ticket_uuid")
	if ticketUUID == "" {
		http.Error(w, "ticket UUID is required", http.StatusBadRequest)
		return
	}

	ticket, err := th.db.GetTicket(ticketUUID)
	if err != nil {
		http.Error(w, "failed to fetch ticket", http.StatusNotFound)
		return
	}

	bounty, err := th.db.CreateBountyFromTicket(ticket)
	if err != nil {
		http.Error(w, "failed to create bounty", http.StatusInternalServerError)
		return
	}

	response := CreateBountyResponse{
		BountyID: bounty.ID,
		Success:  true,
		Message:  "Bounty created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
