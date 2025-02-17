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

type SendTicketPlanRequest struct {
    FeatureID       string   `json:"feature_id"`
    PhaseID         string   `json:"phase_id"`
    TicketGroupIDs  []string `json:"ticket_group_ids"`
    SourceWebsocket string   `json:"source_websocket"`
    RequestUUID     string   `json:"request_uuid"`
}

type SendTicketPlanResponse struct {
    Success     bool     `json:"success"`
    Message     string   `json:"message"`
    RequestUUID string   `json:"request_uuid"`
    Errors      []string `json:"errors,omitempty"`
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

    feature := th.db.GetFeatureByUuid(planRequest.FeatureID)
    if feature.Uuid == "" {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Feature not found",
        })
        return
    }

    phase, err := th.db.GetPhaseByUuid(planRequest.PhaseID)
    if err != nil || phase.Uuid == "" {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Phase not found",
        })
        return
    }

    newPlan := &db.TicketPlan{
        UUID:          uuid.New(),
        WorkspaceUuid: feature.WorkspaceUuid,
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
        websocketErr := websocket.WebsocketPool.SendTicketMessage(websocket.TicketMessage{
            BroadcastType:   "direct",
            SourceSessionID: planRequest.SourceWebsocket,
            Action:          "TICKET_PLAN_CREATED",
            Message:         fmt.Sprintf("Created ticket plan %s", createdPlan.UUID.String()),
        })

        if websocketErr != nil {
            logger.Log.Error("Failed to send websocket message", "error", websocketErr)
        }
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

func (th *ticketHandler) SendTicketPlanToStakwork(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

    if pubKeyFromAuth == "" {
        logger.Log.Info("[ticket plan] no pubkey from auth")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
        return
    }

    user := th.db.GetPersonByPubkey(pubKeyFromAuth)
    if user.OwnerPubKey != pubKeyFromAuth {
        logger.Log.Info("Person not exists")
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Validation failed",
            Errors:  []string{"Error reading request body"},
        })
        return
    }
    defer r.Body.Close()

    var planRequest SendTicketPlanRequest
    if err := json.Unmarshal(body, &planRequest); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Validation failed",
            Errors:  []string{"Error parsing request body: " + err.Error()},
        })
        return
    }

    if planRequest.FeatureID == "" || planRequest.PhaseID == "" || len(planRequest.TicketGroupIDs) == 0 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Validation failed",
            Errors:  []string{"feature_id, phase_id, and ticket_group_ids are required"},
        })
        return
    }

    var (
        productBrief, featureBrief, phaseDesign, codeGraphURL, codeGraphAlias string
        feature                                                                db.WorkspaceFeatures
    )

    feature = th.db.GetFeatureByUuid(planRequest.FeatureID)
    if feature.Uuid == "" {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Error retrieving feature details",
            Errors:  []string{"Feature not found with the provided UUID"},
        })
        return
    }

    productBrief, err = th.db.GetProductBrief(feature.WorkspaceUuid)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Error retrieving product brief",
            Errors:  []string{err.Error()},
        })
        return
    }

    featureBrief, err = th.db.GetFeatureBrief(planRequest.FeatureID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Error retrieving feature brief",
            Errors:  []string{err.Error()},
        })
        return
    }

    phaseDesign, err = th.db.GetPhaseDesign(planRequest.PhaseID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Error retrieving phase design",
            Errors:  []string{err.Error()},
        })
        return
    }

    host := os.Getenv("HOST")
    if host == "" {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "HOST environment variable not set",
        })
        return
    }

    webhookURL := fmt.Sprintf("%s/bounties/ticket/plan/review", host)

    var schematicURL string
    if feature.WorkspaceUuid != "" {
        workspace := th.db.GetWorkspaceByUuid(feature.WorkspaceUuid)
        if workspace.Uuid == "" {
            w.WriteHeader(http.StatusNotFound)
            json.NewEncoder(w).Encode(TicketPlanResponse{
                Success: false,
                Message: "Workspace not found",
            })
            return
        }

        schematicURL = workspace.SchematicUrl

        codeGraph, err := th.db.GetCodeGraphByWorkspaceUuid(feature.WorkspaceUuid)
        if err == nil {
            codeGraphURL = codeGraph.Url
            codeGraphAlias = codeGraph.SecretAlias
        } else {
            codeGraphURL = ""
            codeGraphAlias = ""
        }
    }

    phase, err := th.db.GetPhaseByUuid(planRequest.PhaseID)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    ticketArray := th.db.BuildTicketArray(planRequest.TicketGroupIDs)

    stakworkPayload := map[string]interface{}{
        "name":        "Ticket Plan Builder",
        "workflow_id": 42472,
        "workflow_params": map[string]interface{}{
            "set_var": map[string]interface{}{
                "attributes": map[string]interface{}{
                    "vars": map[string]interface{}{
                        "featureUUID":     planRequest.FeatureID,
                        "phaseUUID":       planRequest.PhaseID,
                        "ticketPlanUUID":  uuid.New().String(),
                        "phaseOutcome":    phase.PhaseOutcome,
                        "phasePurpose":    phase.PhasePurpose,
                        "phaseScope":      phase.PhaseScope,
                        "phaseDesign":     phaseDesign,
                        "ticketArray":     ticketArray,
                        "productBrief":    productBrief,
                        "featureBrief":    featureBrief,
                        "sourceWebsocket": planRequest.SourceWebsocket,
                        "webhook_url":     webhookURL,
                        "phaseSchematic":  schematicURL,
                        "codeGraph":       codeGraphURL,
                        "alias":           user.OwnerAlias,
                        "requestUUID":     planRequest.RequestUUID,
                        "codeGraphAlias":  codeGraphAlias,
                    },
                },
            },
        },
    }

    stakworkPayloadJSON, err := json.Marshal(stakworkPayload)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Error encoding payload",
            Errors:  []string{err.Error()},
        })
        return
    }

    apiKey := os.Getenv("SWWFKEY")
    if apiKey == "" {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "API key not set in environment",
        })
        return
    }

    req, err := http.NewRequest(http.MethodPost, "https://api.stakwork.com/api/v1/projects", bytes.NewBuffer(stakworkPayloadJSON))
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(TicketPlanResponse{
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
        json.NewEncoder(w).Encode(TicketPlanResponse{
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
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Error reading response from Stakwork",
            Errors:  []string{err.Error()},
        })
        return
    }

    var stakworkResp StakworkResponse
    if err := json.Unmarshal(respBody, &stakworkResp); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: "Error parsing Stakwork response",
            Errors:  []string{err.Error()},
        })
        return
    }

    if resp.StatusCode != http.StatusOK || !stakworkResp.Success {
        w.WriteHeader(resp.StatusCode)
        json.NewEncoder(w).Encode(TicketPlanResponse{
            Success: false,
            Message: string(respBody),
            Errors:  []string{fmt.Sprintf("Stakwork API returned status code: %d", resp.StatusCode)},
        })
        return
    }

    if planRequest.SourceWebsocket != "" {
        ticketMsg := websocket.TicketPlanMessage{
            BroadcastType:   "direct",
            SourceSessionID: planRequest.SourceWebsocket,
            Message:         "Processing ticket plan generation",
            Action:          "TICKET_PLAN_PROCESSING",
            PlanDetails: websocket.TicketPlanDetails{
                RequestUUID:  planRequest.RequestUUID,
                FeatureUUID: planRequest.FeatureID,
                PhaseUUID:   planRequest.PhaseID,
            },
        }

        if err := websocket.WebsocketPool.SendTicketPlanMessage(ticketMsg); err != nil {
            log.Printf("Failed to send websocket message: %v", err)
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "plan":            planRequest,
                "websocket_error": err.Error(),
            })
            return
        }

        projectMsg := websocket.TicketPlanMessage{
            BroadcastType:   "direct",
            SourceSessionID: planRequest.SourceWebsocket,
            Message:         fmt.Sprintf("https://jobs.stakwork.com/admin/projects/%d", stakworkResp.Data.ProjectID),
            Action:          "swrun",
            PlanDetails: websocket.TicketPlanDetails{
                RequestUUID:  planRequest.RequestUUID,
                FeatureUUID: planRequest.FeatureID,
                PhaseUUID:   planRequest.PhaseID,
            },
        }

        if err := websocket.WebsocketPool.SendTicketPlanMessage(projectMsg); err != nil {
            log.Printf("Failed to send project ID websocket message: %v", err)
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "plan":            planRequest,
                "websocket_error": err.Error(),
            })
            return
        }
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(SendTicketPlanResponse{
        Success:     true,
        Message:     string(respBody),
        RequestUUID: planRequest.RequestUUID,
    })
}

func (th *ticketHandler) ProcessTicketPlanReview(w http.ResponseWriter, r *http.Request) {

    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("Error reading request body: %v", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(db.TicketPlanReviewResponse{
            Success: false,
            Message: "Error reading request body",
            Errors:  []string{err.Error()},
        })
        return
    }
    defer r.Body.Close()

    var planReview db.TicketPlanReviewRequest
    if err := json.Unmarshal(body, &planReview); err != nil {
        log.Printf("Error parsing request JSON: %v", err)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(db.TicketPlanReviewResponse{
            Success: false,
            Message: "Error parsing request body",
            Errors:  []string{err.Error()},
        })
        return
    }

    feature := th.db.GetFeatureByUuid(planReview.Value.FeatureUUID)
    if feature.Uuid == "" {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(db.TicketPlanReviewResponse{
            Success: false,
            Message: "Feature not found",
        })
        return
    }

    phase, err := th.db.GetPhaseByUuid(planReview.Value.PhaseUUID)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(db.TicketPlanReviewResponse{
            Success: false,
            Message: "Phase not found",
        })
        return
    }

    var createdTickets []db.Tickets
    for i, stub := range planReview.Value.PhasePlan.StubTickets {
        ticketGroup := uuid.New()
        description := fmt.Sprintf("%s\n\nReasoning: %s", stub.TicketDescription, stub.Reasoning)

        ticket := db.Tickets{
            UUID:          uuid.New(),
            TicketGroup:   &ticketGroup,
            WorkspaceUuid: feature.WorkspaceUuid,
            FeatureUUID:   planReview.Value.FeatureUUID,
            PhaseUUID:     planReview.Value.PhaseUUID,
            Name:          stub.TicketName,
            Sequence:      i,
            Description:   description,
            Status:        db.DraftTicket,
            Version:       1,
            Author:        (*db.Author)(nil),
            CreatedAt:     time.Now(),
            UpdatedAt:     time.Now(),
        }

        createdTicket, err := th.db.CreateOrEditTicket(&ticket)
        if err != nil {
            log.Printf("Error creating ticket: %v", err)
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(db.TicketPlanReviewResponse{
                Success: false,
                Message: "Error creating ticket",
                Errors:  []string{err.Error()},
            })
            return
        }
        createdTickets = append(createdTickets, createdTicket)

        if planReview.SourceWebsocket != "" {
            ticketMsg := websocket.TicketMessage{
                BroadcastType:   "direct",
                SourceSessionID: planReview.SourceWebsocket,
                Message:         fmt.Sprintf("Created ticket: %s", stub.TicketName),
                Action:          "process",
                TicketDetails: websocket.TicketData{
                    FeatureUUID:       planReview.Value.FeatureUUID,
                    PhaseUUID:         planReview.Value.PhaseUUID,
                    TicketUUID:        createdTicket.UUID.String(),
                    TicketDescription: description,
                    TicketName:        stub.TicketName,
                },
            }

            if err := websocket.WebsocketPool.SendTicketMessage(ticketMsg); err != nil {
                log.Printf("Failed to send ticket websocket message: %v", err)
            }
        }
    }

    if planReview.SourceWebsocket != "" {
        completionMsg := websocket.TicketPlanMessage{
            BroadcastType:   "direct",
            SourceSessionID: planReview.SourceWebsocket,
            Message:         fmt.Sprintf("Successfully created %d tickets for phase %s", len(createdTickets), phase.Name),
            Action:          "TICKET_PLAN_COMPLETED",
            PlanDetails: websocket.TicketPlanDetails{
                RequestUUID:  planReview.RequestUUID,
                FeatureUUID: planReview.Value.FeatureUUID,
                PhaseUUID:   planReview.Value.PhaseUUID,
            },
        }

        if err := websocket.WebsocketPool.SendTicketPlanMessage(completionMsg); err != nil {
            log.Printf("Failed to send completion websocket message: %v", err)
        }
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(db.TicketPlanReviewResponse{
        Success: true,
        Message: fmt.Sprintf("Successfully created %d tickets", len(createdTickets)),
    })
}
