package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/utils"
)

type ticketHandler struct {
	db db.Database
}

func NewTicketHandler(database db.Database) *ticketHandler {
	return &ticketHandler{
		db: database,
	}
}

type TicketReviewRequest struct {
	FeatureUUID       string `json:"featureUUID" validate:"required"`
	PhaseUUID         string `json:"phaseUUID" validate:"required"`
	TicketUUID        string `json:"ticketUUID" validate:"required"`
	TicketDescription string `json:"ticketDescription" validate:"required"`
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

	var ticket db.Tickets
	if err := json.Unmarshal(body, &ticket); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	ticket.UUID = ticketUUID

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

func (th *ticketHandler) ProcessTicketReview(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

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
