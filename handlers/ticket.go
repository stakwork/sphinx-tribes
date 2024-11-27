package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/db"
)

type ticketHandler struct {
	db db.Database
}

func NewTicketHandler(database db.Database) *ticketHandler {
	return &ticketHandler{
		db: database,
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

	respondWithJSON(w, http.StatusOK, ticket)
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

	respondWithJSON(w, http.StatusOK, updatedTicket)
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

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
