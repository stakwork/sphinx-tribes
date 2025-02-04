package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTicketPlan(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tHandler := NewTicketHandler(&http.Client{}, db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	_, err := db.TestDB.CreateOrEditPerson(person)
	require.NoError(t, err)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace-" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	_, err = db.TestDB.CreateOrEditWorkspace(workspace)
	require.NoError(t, err)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "https://github.com/test-feature",
		Priority:      0,
		CreatedBy:     person.OwnerPubKey,
	}
	_, err = db.TestDB.CreateOrEditFeature(feature)
	require.NoError(t, err)

	phase := db.FeaturePhase{
		Uuid:        uuid.New().String(),
		FeatureUuid: feature.Uuid,
		Name:        "test-phase",
		Priority:    0,
	}
	_, err = db.TestDB.CreateOrEditFeaturePhase(phase)
	require.NoError(t, err)

	tests := []struct {
		name           string
		requestBody    CreateTicketPlanRequest
		auth           string
		expectedStatus int
		validateFunc   func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			name: "successful ticket plan creation",
			requestBody: CreateTicketPlanRequest{
				FeatureID:    feature.Uuid,
				PhaseID:      phase.Uuid,
				Name:         "Test Ticket Plan",
				Description:  "A test ticket plan description",
			},
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp TicketPlanResponse
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.True(t, resp.Success)
				assert.NotEmpty(t, resp.PlanID)
				assert.Equal(t, "Ticket plan created successfully", resp.Message)

				createdPlan, err := db.TestDB.GetTicketPlan(resp.PlanID)
				require.NoError(t, err)
				assert.Equal(t, feature.Uuid, createdPlan.FeatureUUID)
				assert.Equal(t, phase.Uuid, createdPlan.PhaseUUID)
				assert.Equal(t, "Test Ticket Plan", createdPlan.Name)
			},
		},
		{
			name: "unauthorized no auth token",
			requestBody: CreateTicketPlanRequest{
				FeatureID: feature.Uuid,
				PhaseID:   phase.Uuid,
				Name:      "Test Ticket Plan",
			},
			auth:           "",
			expectedStatus: http.StatusUnauthorized,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
				assert.Equal(t, "Unauthorized", resp["error"])
			},
		},
		{
			name: "bad request - missing feature ID",
			requestBody: CreateTicketPlanRequest{
				PhaseID: phase.Uuid,
				Name:    "Test Ticket Plan",
			},
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp TicketPlanResponse
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.False(t, resp.Success)
				assert.Contains(t, resp.Message, "Missing required fields")
			},
		},
		{
			name: "bad request - missing phase ID",
			requestBody: CreateTicketPlanRequest{
				FeatureID: feature.Uuid,
				Name:      "Test Ticket Plan",
			},
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp TicketPlanResponse
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.False(t, resp.Success)
				assert.Contains(t, resp.Message, "Missing required fields")
			},
		},
		{
			name: "bad request - missing name",
			requestBody: CreateTicketPlanRequest{
				FeatureID: feature.Uuid,
				PhaseID:   phase.Uuid,
			},
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp TicketPlanResponse
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.False(t, resp.Success)
				assert.Contains(t, resp.Message, "Missing required fields")
			},
		},
		{
			name: "not found -- invalid feature",
			requestBody: CreateTicketPlanRequest{
				FeatureID: uuid.New().String(),
				PhaseID:   phase.Uuid,
				Name:      "Test Ticket Plan",
			},
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusNotFound,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp TicketPlanResponse
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.False(t, resp.Success)
				assert.Equal(t, "Feature not found", resp.Message)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/bounties/ticket/plan", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.auth != "" {
				req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, tt.auth))
			}

			rr := httptest.NewRecorder()

			tHandler.CreateTicketPlan(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.validateFunc != nil {
				tt.validateFunc(t, rr)
			}
		})
	}
}

func TestGetTicketPlan(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tHandler := NewTicketHandler(&http.Client{}, db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	_, err := db.TestDB.CreateOrEditPerson(person)
	require.NoError(t, err)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace-" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	_, err = db.TestDB.CreateOrEditWorkspace(workspace)
	require.NoError(t, err)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "https://github.com/test-feature",
		Priority:      0,
		CreatedBy:     person.OwnerPubKey,
	}
	_, err = db.TestDB.CreateOrEditFeature(feature)
	require.NoError(t, err)

	phase := db.FeaturePhase{
		Uuid:        uuid.New().String(),
		FeatureUuid: feature.Uuid,
		Name:        "test-phase",
		Priority:    0,
	}
	_, err = db.TestDB.CreateOrEditFeaturePhase(phase)
	require.NoError(t, err)

	ticketPlan := &db.TicketPlan{
		UUID:          uuid.New(),
		WorkspaceUuid: workspace.Uuid,
		FeatureUUID:   feature.Uuid,
		PhaseUUID:     phase.Uuid,
		Name:          "Test Ticket Plan",
		Description:   "Test Description",
		Status:        db.DraftPlan,
		CreatedBy:     person.OwnerPubKey,
	}
	createdPlan, err := db.TestDB.CreateOrEditTicketPlan(ticketPlan)
	require.NoError(t, err)

	tests := []struct {
		name           string
		planUUID       string
		auth           string
		expectedStatus int
		validateFunc   func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			name:           "successful ticket plan retrieval",
			planUUID:       createdPlan.UUID.String(),
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var retrievedPlan db.TicketPlan
				err := json.Unmarshal(response.Body.Bytes(), &retrievedPlan)
				require.NoError(t, err)

				assert.Equal(t, createdPlan.UUID, retrievedPlan.UUID)
				assert.Equal(t, "Test Ticket Plan", retrievedPlan.Name)
				assert.Equal(t, feature.Uuid, retrievedPlan.FeatureUUID)
				assert.Equal(t, phase.Uuid, retrievedPlan.PhaseUUID)
			},
		},
		{
			name:           "unauthorized - no auth token",
			planUUID:       createdPlan.UUID.String(),
			auth:           "",
			expectedStatus: http.StatusUnauthorized,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
				assert.Equal(t, "Unauthorized", resp["error"])
			},
		},
		{
			name:           "bad request - empty UUID",
			planUUID:       "",
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
				assert.Equal(t, "UUID is required", resp["error"])
			},
		},
		{
			name:           "not found - non-existent UUID",
			planUUID:       uuid.New().String(),
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusNotFound,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
				assert.Contains(t, resp["error"], "ticket plan not found")
			},
		},
		{
			name:           "invalid UUID format",
			planUUID:       "invalid-uuid",
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusInternalServerError,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/bounties/ticket/plan/"+tt.planUUID, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("uuid", tt.planUUID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.auth != "" {
				req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, tt.auth))
			}

			rr := httptest.NewRecorder()

			tHandler.GetTicketPlan(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.validateFunc != nil {
				tt.validateFunc(t, rr)
			}
		})
	}
}

func TestDeleteTicketPlan(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tHandler := NewTicketHandler(&http.Client{}, db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	_, err := db.TestDB.CreateOrEditPerson(person)
	require.NoError(t, err)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace-" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	_, err = db.TestDB.CreateOrEditWorkspace(workspace)
	require.NoError(t, err)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "https://github.com/test-feature",
		Priority:      0,
		CreatedBy:     person.OwnerPubKey,
	}
	_, err = db.TestDB.CreateOrEditFeature(feature)
	require.NoError(t, err)

	phase := db.FeaturePhase{
		Uuid:        uuid.New().String(),
		FeatureUuid: feature.Uuid,
		Name:        "test-phase",
		Priority:    0,
	}
	_, err = db.TestDB.CreateOrEditFeaturePhase(phase)
	require.NoError(t, err)

	tests := []struct {
		name           string
		createPlan     bool
		planUUID       string
		auth           string
		expectedStatus int
		validateFunc   func(t *testing.T, response *httptest.ResponseRecorder, planUUID string)
	}{
		{
			name:           "successful ticket plan deletion",
			createPlan:     true,
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder, planUUID string) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, "Ticket plan deleted successfully", resp["message"])

				_, err = db.TestDB.GetTicketPlan(planUUID)
				assert.Error(t, err)
				assert.Contains(t, strings.ToLower(err.Error()), "not found")
			},
		},
		{
			name:           "unauthorized - no auth token",
			createPlan:     true,
			auth:           "",
			expectedStatus: http.StatusUnauthorized,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder, planUUID string) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
				assert.Equal(t, "Unauthorized", resp["error"])
			},
		},
		{
			name:           "bad request - empty UUID",
			createPlan:     false,
			planUUID:       "",
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder, planUUID string) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
				assert.Equal(t, "UUID is required", resp["error"])
			},
		},
		{
			name:           "not found - non-existent UUID",
			createPlan:     false,
			planUUID:       uuid.New().String(),
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusNotFound,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder, planUUID string) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
				assert.Contains(t, resp["error"], "ticket plan not found")
			},
		},
		{
			name:           "invalid UUID format",
			createPlan:     false,
			planUUID:       "invalid-uuid",
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusInternalServerError,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder, planUUID string) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var planUUID string
			if tt.createPlan {
				ticketPlan := &db.TicketPlan{
					UUID:          uuid.New(),
					WorkspaceUuid: workspace.Uuid,
					FeatureUUID:   feature.Uuid,
					PhaseUUID:     phase.Uuid,
					Name:          "Test Ticket Plan",
					Description:   "Test Description",
					Status:        db.DraftPlan,
					CreatedBy:     person.OwnerPubKey,
				}
				createdPlan, err := db.TestDB.CreateOrEditTicketPlan(ticketPlan)
				require.NoError(t, err)
				planUUID = createdPlan.UUID.String()
			} else {
				planUUID = tt.planUUID
			}

			req := httptest.NewRequest(http.MethodDelete, "/bounties/ticket/plan/"+planUUID, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("uuid", planUUID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.auth != "" {
				req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, tt.auth))
			}

			rr := httptest.NewRecorder()

			tHandler.DeleteTicketPlan(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.validateFunc != nil {
				tt.validateFunc(t, rr, planUUID)
			}
		})
	}
}

func TestGetTicketPlansByFeature(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tHandler := NewTicketHandler(&http.Client{}, db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	_, err := db.TestDB.CreateOrEditPerson(person)
	require.NoError(t, err)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace-" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	_, err = db.TestDB.CreateOrEditWorkspace(workspace)
	require.NoError(t, err)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "https://github.com/test-feature",
		Priority:      0,
		CreatedBy:     person.OwnerPubKey,
	}
	_, err = db.TestDB.CreateOrEditFeature(feature)
	require.NoError(t, err)

	phase := db.FeaturePhase{
		Uuid:        uuid.New().String(),
		FeatureUuid: feature.Uuid,
		Name:        "test-phase",
		Priority:    0,
	}
	_, err = db.TestDB.CreateOrEditFeaturePhase(phase)
	require.NoError(t, err)

	ticketPlans := []*db.TicketPlan{
		{
			UUID:          uuid.New(),
			WorkspaceUuid: workspace.Uuid,
			FeatureUUID:   feature.Uuid,
			PhaseUUID:     phase.Uuid,
			Name:          "Test Ticket Plan 1",
			Description:   "Test Description 1",
			Status:        db.DraftPlan,
			CreatedBy:     person.OwnerPubKey,
		},
		{
			UUID:          uuid.New(),
			WorkspaceUuid: workspace.Uuid,
			FeatureUUID:   feature.Uuid,
			PhaseUUID:     phase.Uuid,
			Name:          "Test Ticket Plan 2",
			Description:   "Test Description 2",
			Status:        db.DraftPlan,
			CreatedBy:     person.OwnerPubKey,
		},
	}

	var createdPlans []string
	for _, plan := range ticketPlans {
		createdPlan, err := db.TestDB.CreateOrEditTicketPlan(plan)
		require.NoError(t, err)
		createdPlans = append(createdPlans, createdPlan.UUID.String())
	}

	tests := []struct {
		name           string
		featureUUID    string
		auth           string
		expectedStatus int
		validateFunc   func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			name:           "successful retrieval of ticket plans by feature",
			featureUUID:    feature.Uuid,
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var plans []db.TicketPlan
				err := json.Unmarshal(response.Body.Bytes(), &plans)
				require.NoError(t, err)
				
				assert.Equal(t, 2, len(plans), "Should return 2 ticket plans")
				
				planNames := make(map[string]bool)
				for _, plan := range plans {
					assert.Equal(t, feature.Uuid, plan.FeatureUUID, "Feature UUID should match")
					planNames[plan.Name] = true
				}
				assert.Contains(t, planNames, "Test Ticket Plan 1")
				assert.Contains(t, planNames, "Test Ticket Plan 2")
			},
		},
		{
			name:           "unauthorized - no auth token",
			featureUUID:    feature.Uuid,
			auth:           "",
			expectedStatus: http.StatusUnauthorized,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
				assert.Equal(t, "Unauthorized", resp["error"])
			},
		},
		{
			name:           "bad request - empty feature UUID",
			featureUUID:    "",
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
				assert.Equal(t, "Feature UUID is required", resp["error"])
			},
		},
		{
			name:           "no plans for feature",
			featureUUID:    uuid.New().String(),
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var plans []db.TicketPlan
				err := json.Unmarshal(response.Body.Bytes(), &plans)
				require.NoError(t, err)
				assert.Equal(t, 0, len(plans), "Should return empty list for non-existent feature")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/bounties/ticket/plan/feature/"+tt.featureUUID, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("feature_uuid", tt.featureUUID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.auth != "" {
				req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, tt.auth))
			}

			rr := httptest.NewRecorder()

			tHandler.GetTicketPlansByFeature(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.validateFunc != nil {
				tt.validateFunc(t, rr)
			}
		})
	}
}

func TestGetTicketPlansByPhase(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tHandler := NewTicketHandler(&http.Client{}, db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	_, err := db.TestDB.CreateOrEditPerson(person)
	require.NoError(t, err)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace-" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	_, err = db.TestDB.CreateOrEditWorkspace(workspace)
	require.NoError(t, err)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "https://github.com/test-feature",
		Priority:      0,
		CreatedBy:     person.OwnerPubKey,
	}
	_, err = db.TestDB.CreateOrEditFeature(feature)
	require.NoError(t, err)

	phase := db.FeaturePhase{
		Uuid:        uuid.New().String(),
		FeatureUuid: feature.Uuid,
		Name:        "test-phase",
		Priority:    0,
	}
	_, err = db.TestDB.CreateOrEditFeaturePhase(phase)
	require.NoError(t, err)

	ticketPlans := []*db.TicketPlan{
		{
			UUID:          uuid.New(),
			WorkspaceUuid: workspace.Uuid,
			FeatureUUID:   feature.Uuid,
			PhaseUUID:     phase.Uuid,
			Name:          "Test Ticket Plan 1",
			Description:   "Test Description 1",
			Status:        db.DraftPlan,
			CreatedBy:     person.OwnerPubKey,
		},
		{
			UUID:          uuid.New(),
			WorkspaceUuid: workspace.Uuid,
			FeatureUUID:   feature.Uuid,
			PhaseUUID:     phase.Uuid,
			Name:          "Test Ticket Plan 2",
			Description:   "Test Description 2",
			Status:        db.DraftPlan,
			CreatedBy:     person.OwnerPubKey,
		},
	}

	var createdPlans []string
	for _, plan := range ticketPlans {
		createdPlan, err := db.TestDB.CreateOrEditTicketPlan(plan)
		require.NoError(t, err)
		createdPlans = append(createdPlans, createdPlan.UUID.String())
	}

	tests := []struct {
		name           string
		phaseUUID      string
		auth           string
		expectedStatus int
		validateFunc   func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			name:           "successful retrieval of ticket plans by phase",
			phaseUUID:      phase.Uuid,
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var plans []db.TicketPlan
				err := json.Unmarshal(response.Body.Bytes(), &plans)
				require.NoError(t, err)
				
				assert.Equal(t, 2, len(plans), "Should return 2 ticket plans")
				
				planNames := make(map[string]bool)
				for _, plan := range plans {
					assert.Equal(t, phase.Uuid, plan.PhaseUUID, "Phase UUID should match")
					planNames[plan.Name] = true
				}
				assert.Contains(t, planNames, "Test Ticket Plan 1")
				assert.Contains(t, planNames, "Test Ticket Plan 2")
			},
		},
		{
			name:           "unauthorized - no auth token",
			phaseUUID:      phase.Uuid,
			auth:           "",
			expectedStatus: http.StatusUnauthorized,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
				assert.Equal(t, "Unauthorized", resp["error"])
			},
		},
		{
			name:           "bad request - empty phase UUID",
			phaseUUID:      "",
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
				assert.Equal(t, "Phase UUID is required", resp["error"])
			},
		},
		{
			name:           "no plans for phase",
			phaseUUID:      uuid.New().String(),
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var plans []db.TicketPlan
				err := json.Unmarshal(response.Body.Bytes(), &plans)
				require.NoError(t, err)
				assert.Equal(t, 0, len(plans), "Should return empty list for non-existent phase")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/bounties/ticket/plan/phase/"+tt.phaseUUID, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("phase_uuid", tt.phaseUUID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.auth != "" {
				req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, tt.auth))
			}

			rr := httptest.NewRecorder()

			tHandler.GetTicketPlansByPhase(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.validateFunc != nil {
				tt.validateFunc(t, rr)
			}
		})
	}
}

func TestGetTicketPlansByWorkspace(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	tHandler := NewTicketHandler(&http.Client{}, db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	_, err := db.TestDB.CreateOrEditPerson(person)
	require.NoError(t, err)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace-" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	_, err = db.TestDB.CreateOrEditWorkspace(workspace)
	require.NoError(t, err)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "https://github.com/test-feature",
		Priority:      0,
		CreatedBy:     person.OwnerPubKey,
	}
	_, err = db.TestDB.CreateOrEditFeature(feature)
	require.NoError(t, err)

	phase := db.FeaturePhase{
		Uuid:        uuid.New().String(),
		FeatureUuid: feature.Uuid,
		Name:        "test-phase",
		Priority:    0,
	}
	_, err = db.TestDB.CreateOrEditFeaturePhase(phase)
	require.NoError(t, err)

	ticketPlans := []*db.TicketPlan{
		{
			UUID:          uuid.New(),
			WorkspaceUuid: workspace.Uuid,
			FeatureUUID:   feature.Uuid,
			PhaseUUID:     phase.Uuid,
			Name:          "Test Ticket Plan 1",
			Description:   "Test Description 1",
			Status:        db.DraftPlan,
			CreatedBy:     person.OwnerPubKey,
		},
		{
			UUID:          uuid.New(),
			WorkspaceUuid: workspace.Uuid,
			FeatureUUID:   feature.Uuid,
			PhaseUUID:     phase.Uuid,
			Name:          "Test Ticket Plan 2",
			Description:   "Test Description 2",
			Status:        db.DraftPlan,
			CreatedBy:     person.OwnerPubKey,
		},
	}

	var createdPlans []string
	for _, plan := range ticketPlans {
		createdPlan, err := db.TestDB.CreateOrEditTicketPlan(plan)
		require.NoError(t, err)
		createdPlans = append(createdPlans, createdPlan.UUID.String())
	}

	tests := []struct {
		name           string
		workspaceUUID  string
		auth           string
		expectedStatus int
		validateFunc   func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			name:           "successful retrieval of ticket plans by workspace",
			workspaceUUID:  workspace.Uuid,
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var plans []db.TicketPlan
				err := json.Unmarshal(response.Body.Bytes(), &plans)
				require.NoError(t, err)
				
				assert.Equal(t, 2, len(plans), "Should return 2 ticket plans")
				
				planNames := make(map[string]bool)
				for _, plan := range plans {
					assert.Equal(t, workspace.Uuid, plan.WorkspaceUuid, "Workspace UUID should match")
					planNames[plan.Name] = true
				}
				assert.Contains(t, planNames, "Test Ticket Plan 1")
				assert.Contains(t, planNames, "Test Ticket Plan 2")
			},
		},
		{
			name:           "unauthorized - no auth token",
			workspaceUUID:  workspace.Uuid,
			auth:           "",
			expectedStatus: http.StatusUnauthorized,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
				assert.Equal(t, "Unauthorized", resp["error"])
			},
		},
		{
			name:           "bad request - empty workspace UUID",
			workspaceUUID:  "",
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var resp map[string]string
				err := json.Unmarshal(response.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Contains(t, resp, "error")
				assert.Equal(t, "Workspace UUID is required", resp["error"])
			},
		},
		{
			name:           "no plans for workspace",
			workspaceUUID:  uuid.New().String(),
			auth:           person.OwnerPubKey,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var plans []db.TicketPlan
				err := json.Unmarshal(response.Body.Bytes(), &plans)
				require.NoError(t, err)
				assert.Equal(t, 0, len(plans), "Should return empty list for non-existent workspace")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/bounties/ticket/plan/workspace/"+tt.workspaceUUID, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("workspace_uuid", tt.workspaceUUID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.auth != "" {
				req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, tt.auth))
			}

			rr := httptest.NewRecorder()

			tHandler.GetTicketPlansByWorkspace(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.validateFunc != nil {
				tt.validateFunc(t, rr)
			}
		})
	}
}
