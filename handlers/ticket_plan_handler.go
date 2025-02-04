package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stakwork/sphinx-tribes/logger"
	"github.com/stakwork/sphinx-tribes/websocket"
)

type CreateTicketPlanRequest struct {
    FeatureID        string   `json:"feature_id"`
    PhaseID          string   `json:"phase_id"`
    Name             string   `json:"name"`
    Description      string   `json:"description"`
    TicketGroupIDs   []string `json:"ticket_group_ids"`
    SourceWebsocket  string   `json:"source_websocket,omitempty"`
}

type TicketPlanResponse struct {
    Success bool     `json:"success"`
    PlanID  string   `json:"plan_id,omitempty"`
    Message string   `json:"message"`
    Errors  []string `json:"errors,omitempty"`
}

type TicketArrayItem struct {
    TicketName        string `json:"ticket_name"`
    TicketDescription string `json:"ticket_description"`
}

func (th *ticketHandler) CreateTicketPlan(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

    if pubKeyFromAuth == "" {
        logger.Log.Info("[ticket plan] no pubkey from auth")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
        return
    }

    var planRequest CreateTicketPlanRequest
    if err := json.NewDecoder(r.Body).Decode(&planRequest); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Invalid request body",
            Errors:  []string{err.Error()},
        })
        return
    }

    if planRequest.FeatureID == "" || planRequest.PhaseID == "" || planRequest.Name == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Missing required fields",
            Errors:  []string{"feature_id, phase_id, and name are required"},
        })
        return
    }

    newPlan := &db.TicketPlan{
        UUID:          uuid.New(),
        WorkspaceUuid: "",
        FeatureUUID:   planRequest.FeatureID,
        PhaseUUID:     planRequest.PhaseID,
        Name:          planRequest.Name,
        Description:   planRequest.Description,
        TicketGroups:  planRequest.TicketGroupIDs,
        Status:        db.DraftPlan,
        Version:       1,
        CreatedBy:     pubKeyFromAuth,
        UpdatedBy:     pubKeyFromAuth,
    }

    feature := th.db.GetFeatureByUuid(planRequest.FeatureID)
    if feature.Uuid == "" {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Feature not found",
        })
        return
    }
    newPlan.WorkspaceUuid = feature.WorkspaceUuid

    createdPlan, err := th.db.CreateOrEditTicketPlan(newPlan)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Failed to create ticket plan",
            Errors:  []string{err.Error()},
        })
        return
    }

    if planRequest.SourceWebsocket != "" {
        websocket.WebsocketPool.SendTicketMessage(websocket.TicketMessage{
            BroadcastType:   "direct",
            SourceSessionID: planRequest.SourceWebsocket,
            Message:         fmt.Sprintf("Created ticket plan %s", createdPlan.UUID.String()),
            Action:          "TICKET_PLAN_CREATED",
        })
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(TicketPlanResponse{
        Success: true,
        PlanID:  createdPlan.UUID.String(),
        Message: "Ticket plan created successfully",
    })
}

func (th *ticketHandler) GetTicketPlan(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

    if pubKeyFromAuth == "" {
        logger.Log.Info("[ticket plan] no pubkey from auth")
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

    plan, err := th.db.GetTicketPlan(uuid)
    if err != nil {
        status := http.StatusInternalServerError
        if err.Error() == "ticket plan not found" {
            status = http.StatusNotFound
        }
        w.WriteHeader(status)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(plan)
}

func (th *ticketHandler) DeleteTicketPlan(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

    if pubKeyFromAuth == "" {
        logger.Log.Info("[ticket plan] no pubkey from auth")
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

    err := th.db.DeleteTicketPlan(uuid)
    if err != nil {
        status := http.StatusInternalServerError
        if err.Error() == "ticket plan not found" {
            status = http.StatusNotFound
        }
        w.WriteHeader(status)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Ticket plan deleted successfully"})
}

func (th *ticketHandler) GetTicketPlansByFeature(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

    if pubKeyFromAuth == "" {
        logger.Log.Info("[ticket plan] no pubkey from auth")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
        return
    }

    featureUUID := chi.URLParam(r, "feature_uuid")
    if featureUUID == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Feature UUID is required"})
        return
    }

    plans, err := th.db.GetTicketPlansByFeature(featureUUID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(plans)
}

func (th *ticketHandler) GetTicketPlansByPhase(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

    if pubKeyFromAuth == "" {
        logger.Log.Info("[ticket plan] no pubkey from auth")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
        return
    }

    phaseUUID := chi.URLParam(r, "phase_uuid")
    if phaseUUID == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Phase UUID is required"})
        return
    }

    plans, err := th.db.GetTicketPlansByPhase(phaseUUID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(plans)
}

func (th *ticketHandler) GetTicketPlansByWorkspace(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

    if pubKeyFromAuth == "" {
        logger.Log.Info("[ticket plan] no pubkey from auth")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
        return
    }

    workspaceUUID := chi.URLParam(r, "workspace_uuid")
    if workspaceUUID == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Workspace UUID is required"})
        return
    }

    plans, err := th.db.GetTicketPlansByWorkspace(workspaceUUID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(plans)
} 