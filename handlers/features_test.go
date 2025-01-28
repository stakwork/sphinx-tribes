package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrEditFeatures(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	fHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)
	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "https://github.com/test-feature",
		Priority:      0,
	}

	t.Run("should return 401 error if not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeatures)

		requestBody, _ := json.Marshal(feature)
		req, err := http.NewRequest(http.MethodPost, "/features", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return 406 error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeatures)

		invalidJson := []byte(`{"key": "value"`)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("should return 401 error if workspace UUID does not exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeatures)

		feature.WorkspaceUuid = "non-existent-uuid"
		requestBody, _ := json.Marshal(feature)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should successfully add feature if request is valid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeatures)

		feature.WorkspaceUuid = workspace.Uuid
		requestBody, _ := json.Marshal(feature)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		createdFeature := db.TestDB.GetFeatureByUuid(feature.Uuid)
		assert.Equal(t, feature.Name, createdFeature.Name)
		assert.Equal(t, feature.Url, createdFeature.Url)
		assert.Equal(t, feature.Priority, createdFeature.Priority)
	})
}

func TestDeleteFeature(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	oHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)
	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "https://github.com/test-feature",
		Priority:      0,
	}
	db.TestDB.CreateOrEditFeature(feature)
	feature = db.TestDB.GetFeatureByUuid(feature.Uuid)

	ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)

	t.Run("should return error if not authorized", func(t *testing.T) {
		featureUUID := feature.Uuid
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteFeature)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", featureUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodDelete, "/features/"+featureUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should delete feature on successful delete", func(t *testing.T) {
		featureUUID := feature.Uuid

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.DeleteFeature)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", featureUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/features/"+featureUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		deletedFeature := db.TestDB.GetFeatureByUuid(featureUUID)
		assert.Equal(t, db.WorkspaceFeatures{}, deletedFeature)
	})
}

func TestGetWorkspaceFeaturesCount(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewFeatureHandler(db.TestDB)

	db.CleanTestData()

	person := db.Person{
		Uuid:        "person-uuid",
		OwnerAlias:  "alias",
		UniqueName:  "unique_name",
		OwnerPubKey: "validPubKey",
		PriceToMeet: 0,
		Description: "description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "workspace_name",
		OwnerPubKey: person.OwnerPubKey,
		Github:      "github",
		Website:     "website",
		Description: "description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	features := []db.WorkspaceFeatures{
		{
			Uuid:          "feature-uuid-1",
			WorkspaceUuid: workspace.Uuid,
			Name:          "feature_1",
			Url:           "url_1",
			Priority:      1,
		},
		{
			Uuid:          "feature-uuid-2",
			WorkspaceUuid: workspace.Uuid,
			Name:          "feature_2",
			Url:           "url_2",
			Priority:      2,
		},
	}

	for _, feature := range features {
		db.TestDB.CreateOrEditFeature(feature)
	}

	tests := []struct {
		name           string
		uuid           string
		pubKeyFromAuth string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Request with Authenticated User",
			uuid:           workspace.Uuid,
			pubKeyFromAuth: person.OwnerPubKey,
			expectedStatus: http.StatusOK,
			expectedBody:   "2",
		},
		{
			name:           "Missing UUID Parameter",
			uuid:           "",
			pubKeyFromAuth: person.OwnerPubKey,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty UUID Parameter",
			uuid:           "",
			pubKeyFromAuth: person.OwnerPubKey,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unauthenticated User",
			uuid:           workspace.Uuid,
			pubKeyFromAuth: "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid UUID Format",
			uuid:           "invalid-uuid@#",
			pubKeyFromAuth: person.OwnerPubKey,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Large Number of Features",
			uuid:           workspace.Uuid,
			pubKeyFromAuth: person.OwnerPubKey,
			expectedStatus: http.StatusOK,
			expectedBody:   "2",
		},
		{
			name:           "UUID with Maximum Integer Features",
			uuid:           workspace.Uuid,
			pubKeyFromAuth: person.OwnerPubKey,
			expectedStatus: http.StatusOK,
			expectedBody:   "2",
		},
		{
			name:           "UUID with No Features",
			uuid:           workspace.Uuid,
			pubKeyFromAuth: person.OwnerPubKey,
			expectedStatus: http.StatusOK,
			expectedBody:   "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("uuid", tt.uuid)

			req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/workspace/count/"+tt.uuid, nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(req.Context(), auth.ContextKey, tt.pubKeyFromAuth)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			http.HandlerFunc(oHandler.GetWorkspaceFeaturesCount).ServeHTTP(rr, req)

			if tt.name == "UUID with Maximum Integer Features" {
				db.CleanTestData()
			}

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
	db.CleanTestData()
}

func TestGetFeatureByUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()

	fHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)
	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

	feature := db.WorkspaceFeatures{
		ID:            1,
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature1",
		Brief:         "brief",
		Architecture:  "architecture",
		Requirements:  "requirements",
		Url:           "https://github.com/test-feature",
		Priority:      0,
	}
	db.TestDB.CreateOrEditFeature(feature)
	feature = db.TestDB.GetFeatureByUuid(feature.Uuid)

	tests := []struct {
		name           string
		contextKey     interface{}
		contextValue   interface{}
		uuid           string
		expectedStatus int
		expectedBody   db.WorkspaceFeatures
		dbError        error
	}{
		{
			name:           "Valid UUID with Authorization",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           feature.Uuid,
			expectedStatus: http.StatusOK,
			expectedBody:   feature,
		},
		{
			name:           "Empty UUID",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   db.WorkspaceFeatures{},
		},
		{
			name:           "Non-Existent UUID",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           "non-existent-uuid",
			expectedStatus: http.StatusOK,
			expectedBody:   db.WorkspaceFeatures{},
		},
		{
			name:           "Missing Authorization",
			contextKey:     auth.ContextKey,
			contextValue:   "",
			uuid:           feature.Uuid,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   db.WorkspaceFeatures{},
		},
		{
			name:           "Invalid UUID Format",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           "invalid-uuid",
			expectedStatus: http.StatusOK,
			expectedBody:   db.WorkspaceFeatures{},
		},
		{
			name:           "Authorization Context Key Mismatch",
			contextKey:     "wrongKey",
			contextValue:   person.OwnerPubKey,
			uuid:           feature.Uuid,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   db.WorkspaceFeatures{},
		},
		{
			name:           "Null person owner pubkey",
			contextKey:     auth.ContextKey,
			contextValue:   "",
			uuid:           "valid-uuid",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   db.WorkspaceFeatures{},
		},
		{
			name:           "Slow Database Response",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           feature.Uuid,
			expectedStatus: http.StatusOK,
			expectedBody:   feature,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/features/"+tt.uuid, nil)
			ctx := context.WithValue(req.Context(), tt.contextKey, tt.contextValue)
			req = req.WithContext(ctx)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("uuid", tt.uuid)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

			rr := httptest.NewRecorder()
			http.HandlerFunc(fHandler.GetFeatureByUuid).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var responseBody db.WorkspaceFeatures
			json.NewDecoder(rr.Body).Decode(&responseBody)

			assert.Equal(t, tt.expectedBody.Uuid, responseBody.Uuid)

		})
	}
}

func TestCreateOrEditFeaturePhase(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	fHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        "uuid",
		OwnerAlias:  "alias",
		UniqueName:  "unique_name",
		OwnerPubKey: "pubkey",
		PriceToMeet: 0,
		Description: "description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        "workspace_uuid",
		Name:        "workspace_name",
		OwnerPubKey: "person.OwnerPubkey",
		Github:      "gtihub",
		Website:     "website",
		Description: "description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	feature := db.WorkspaceFeatures{
		Uuid:          "feature_uuid",
		WorkspaceUuid: workspace.Uuid,
		Name:          "feature_name",
		Url:           "feature_url",
		Priority:      0,
	}
	db.TestDB.CreateOrEditFeature(feature)

	t.Run("should return 401 error if not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		featurePhase := db.FeaturePhase{
			Uuid:        "feature_phase_uuid",
			FeatureUuid: feature.Uuid,
			Name:        "feature_phase_name",
			Priority:    0,
		}

		requestBody, _ := json.Marshal(featurePhase)
		req, err := http.NewRequest(http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return 406 error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		invalidJson := []byte(`{"key": "value"`)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("should return 401 error if a Feature UUID that does not exist Is passed to the API body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		featurePhase := db.FeaturePhase{
			Uuid:        "feature_phase_uuid",
			FeatureUuid: "non-existent-uuid",
			Name:        "feature_phase_name",
			Priority:    0,
		}

		requestBody, _ := json.Marshal(featurePhase)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should successfully user can add a feature phase when the right conditions are met", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		featurePhase := db.FeaturePhase{
			Uuid:        "feature_phase_uuid",
			FeatureUuid: feature.Uuid,
			Name:        "feature_phase_name",
			Priority:    0,
		}

		requestBody, _ := json.Marshal(featurePhase)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		createdFeaturePhase, _ := db.TestDB.GetFeaturePhaseByUuid(feature.Uuid, featurePhase.Uuid)

		assert.Equal(t, featurePhase.Name, createdFeaturePhase.Name)
		assert.Equal(t, featurePhase.FeatureUuid, createdFeaturePhase.FeatureUuid)
		assert.Equal(t, featurePhase.Priority, createdFeaturePhase.Priority)
	})

	t.Run("should successfully create a feature phase with all new fields", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		featurePhase := db.FeaturePhase{
			Uuid:         "feature_phase_uuid_full",
			FeatureUuid:  feature.Uuid,
			Name:         "feature_phase_name",
			Priority:     0,
			PhasePurpose: "Test phase purpose",
			PhaseOutcome: "Expected test outcome",
			PhaseScope:   "Test phase scope",
		}

		requestBody, _ := json.Marshal(featurePhase)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		createdFeaturePhase, _ := db.TestDB.GetFeaturePhaseByUuid(feature.Uuid, featurePhase.Uuid)

		assert.Equal(t, featurePhase.Name, createdFeaturePhase.Name)
		assert.Equal(t, featurePhase.FeatureUuid, createdFeaturePhase.FeatureUuid)
		assert.Equal(t, featurePhase.Priority, createdFeaturePhase.Priority)
		assert.Equal(t, featurePhase.PhasePurpose, createdFeaturePhase.PhasePurpose)
		assert.Equal(t, featurePhase.PhaseOutcome, createdFeaturePhase.PhaseOutcome)
		assert.Equal(t, featurePhase.PhaseScope, createdFeaturePhase.PhaseScope)
	})

	t.Run("should successfully create a feature phase with partial new fields", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		featurePhase := db.FeaturePhase{
			Uuid:         "feature_phase_uuid_partial",
			FeatureUuid:  feature.Uuid,
			Name:         "feature_phase_name",
			Priority:     0,
			PhasePurpose: "Test phase purpose",
			PhaseScope:   "Test phase scope",
		}

		requestBody, _ := json.Marshal(featurePhase)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		createdFeaturePhase, _ := db.TestDB.GetFeaturePhaseByUuid(feature.Uuid, featurePhase.Uuid)

		assert.Equal(t, featurePhase.Name, createdFeaturePhase.Name)
		assert.Equal(t, featurePhase.PhasePurpose, createdFeaturePhase.PhasePurpose)
		assert.Equal(t, "", createdFeaturePhase.PhaseOutcome)
		assert.Equal(t, featurePhase.PhaseScope, createdFeaturePhase.PhaseScope)
	})

	t.Run("should handle empty request body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("should handle invalid priority value", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		featurePhase := db.FeaturePhase{
			Uuid:        "feature_phase_uuid_priority",
			FeatureUuid: feature.Uuid,
			Name:        "feature_phase_name",
			Priority:    -1,
		}

		requestBody, _ := json.Marshal(featurePhase)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("should handle updating existing phase", func(t *testing.T) {

		existingPhase := db.FeaturePhase{
			Uuid:        "existing_phase_uuid",
			FeatureUuid: feature.Uuid,
			Name:        "original_name",
			Priority:    0,
		}
		db.TestDB.CreateOrEditFeaturePhase(existingPhase)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		updatedPhase := existingPhase
		updatedPhase.Name = "updated_name"
		updatedPhase.Priority = 1

		requestBody, _ := json.Marshal(updatedPhase)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		updatedFeaturePhase, _ := db.TestDB.GetFeaturePhaseByUuid(feature.Uuid, existingPhase.Uuid)
		assert.Equal(t, "updated_name", updatedFeaturePhase.Name)
		assert.Equal(t, 1, updatedFeaturePhase.Priority)
	})

	t.Run("should handle extremely long field values", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		longString := strings.Repeat("a", 1000)
		featurePhase := db.FeaturePhase{
			Uuid:         "feature_phase_uuid_long",
			FeatureUuid:  feature.Uuid,
			Name:         longString,
			Priority:     0,
			PhasePurpose: longString,
			PhaseOutcome: longString,
			PhaseScope:   longString,
		}

		requestBody, _ := json.Marshal(featurePhase)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("should handle special characters in fields", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		specialChars := "!@#$%^&*()?><,./;'[]\\{}|`~"
		featurePhase := db.FeaturePhase{
			Uuid:         "feature_phase_uuid_special",
			FeatureUuid:  feature.Uuid,
			Name:         "Name with " + specialChars,
			Priority:     0,
			PhasePurpose: "Purpose with " + specialChars,
			PhaseOutcome: "Outcome with " + specialChars,
			PhaseScope:   "Scope with " + specialChars,
		}

		requestBody, _ := json.Marshal(featurePhase)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		createdFeaturePhase, _ := db.TestDB.GetFeaturePhaseByUuid(feature.Uuid, featurePhase.Uuid)
		assert.Equal(t, featurePhase.Name, createdFeaturePhase.Name)
	})

	t.Run("should handle invalid auth token type", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		featurePhase := db.FeaturePhase{
			Uuid:        "feature_phase_uuid_auth",
			FeatureUuid: feature.Uuid,
			Name:        "feature_phase_name",
			Priority:    0,
		}

		requestBody, _ := json.Marshal(featurePhase)

		ctx := context.WithValue(context.Background(), auth.ContextKey, 12345)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Valid Input with New Phase", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		newPhase := db.FeaturePhase{
			FeatureUuid:  feature.Uuid,
			Name:         "new_phase_name",
			Priority:     1,
			PhasePurpose: "New phase purpose",
			PhaseOutcome: "New phase outcome",
			PhaseScope:   "New phase scope",
		}

		requestBody, _ := json.Marshal(newPhase)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		var response db.FeaturePhase
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Uuid)
	})

	t.Run("Empty UUID in Input", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		emptyUUIDPhase := db.FeaturePhase{
			Uuid:        "",
			FeatureUuid: feature.Uuid,
			Name:        "empty_uuid_phase",
			Priority:    0,
		}

		requestBody, _ := json.Marshal(emptyUUIDPhase)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		var response db.FeaturePhase
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Uuid)
	})

	t.Run("No Public Key in Context", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		phase := db.FeaturePhase{
			Uuid:        uuid.New().String(),
			FeatureUuid: feature.Uuid,
			Name:        "no_pubkey_phase",
			Priority:    0,
		}

		requestBody, _ := json.Marshal(phase)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid JSON Body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		invalidJSON := []byte(`{
        "uuid": "invalid_json,
        "feature_uuid": "missing_quote
    }`)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(invalidJSON))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("Feature Does Not Exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditFeaturePhase)

		nonExistentFeaturePhase := db.FeaturePhase{
			Uuid:        uuid.New().String(),
			FeatureUuid: uuid.New().String(),
			Name:        "non_existent_feature_phase",
			Priority:    0,
		}

		requestBody, _ := json.Marshal(nonExistentFeaturePhase)

		ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/phase", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

}

func TestGetFeaturePhases(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        "uuid",
		OwnerAlias:  "alias",
		UniqueName:  "unique_name",
		OwnerPubKey: "pubkey",
		PriceToMeet: 0,
		Description: "description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        "workspace_uuid",
		Name:        "workspace_name",
		OwnerPubKey: "person.OwnerPubkey",
		Github:      "gtihub",
		Website:     "website",
		Description: "description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	feature := db.WorkspaceFeatures{
		Uuid:          "feature_uuid",
		WorkspaceUuid: workspace.Uuid,
		Name:          "feature_name",
		Url:           "feature_url",
		Priority:      0,
	}
	db.TestDB.CreateOrEditFeature(feature)

	featurePhase := db.FeaturePhase{
		Uuid:        "feature_phase_uuid",
		FeatureUuid: feature.Uuid,
		Name:        "feature_phase_name",
		Priority:    0,
	}
	db.TestDB.CreateOrEditFeaturePhase(featurePhase)

	ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

	t.Run("Should test that it throws a 401 error if a user is not authorized", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/features/"+feature.Uuid+"/phase", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that the workspace features phases array returned from the API has the feature phases created", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/features/"+feature.Uuid+"/phase", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

		var returnedFeaturePhases []db.FeaturePhase
		err = json.Unmarshal(rr.Body.Bytes(), &returnedFeaturePhases)
		assert.NoError(t, err)

		updatedFeaturePhases := db.TestDB.GetPhasesByFeatureUuid(feature.Uuid)

		for i := range updatedFeaturePhases {
			created := updatedFeaturePhases[i].Created.In(time.UTC)
			updated := updatedFeaturePhases[i].Updated.In(time.UTC)
			updatedFeaturePhases[i].Created = &created
			updatedFeaturePhases[i].Updated = &updated
		}

		assert.Equal(t, returnedFeaturePhases, updatedFeaturePhases)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should handle non-existent feature UUID", func(t *testing.T) {
		nonExistentUUID := uuid.New().String()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", nonExistentUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/"+nonExistentUUID+"/phase", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

		var returnedFeaturePhases []db.FeaturePhase
		err = json.Unmarshal(rr.Body.Bytes(), &returnedFeaturePhases)
		assert.NoError(t, err)
		assert.Empty(t, returnedFeaturePhases)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should handle multiple phases for a feature", func(t *testing.T) {

		for i := 1; i <= 3; i++ {
			phase := db.FeaturePhase{
				Uuid:        fmt.Sprintf("multi_phase_uuid_%d", i),
				FeatureUuid: feature.Uuid,
				Name:        fmt.Sprintf("Phase %d", i),
				Priority:    i,
			}
			db.TestDB.CreateOrEditFeaturePhase(phase)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/"+feature.Uuid+"/phase", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

		var returnedFeaturePhases []db.FeaturePhase
		err = json.Unmarshal(rr.Body.Bytes(), &returnedFeaturePhases)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(returnedFeaturePhases), 3)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should handle empty feature UUID", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", "")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features//phase", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

		var returnedFeaturePhases []db.FeaturePhase
		err = json.Unmarshal(rr.Body.Bytes(), &returnedFeaturePhases)
		assert.NoError(t, err)
		assert.Empty(t, returnedFeaturePhases)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should handle feature with no phases", func(t *testing.T) {

		emptyFeature := db.WorkspaceFeatures{
			Uuid:          "empty_feature_uuid",
			WorkspaceUuid: workspace.Uuid,
			Name:          "Empty Feature",
			Url:           "empty_feature_url",
			Priority:      0,
		}
		db.TestDB.CreateOrEditFeature(emptyFeature)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", emptyFeature.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/"+emptyFeature.Uuid+"/phase", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

		var returnedFeaturePhases []db.FeaturePhase
		err = json.Unmarshal(rr.Body.Bytes(), &returnedFeaturePhases)
		assert.NoError(t, err)
		assert.Empty(t, returnedFeaturePhases)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should handle invalid authorization token", func(t *testing.T) {
		invalidCtx := context.WithValue(context.Background(), auth.ContextKey, "invalid-token")
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(invalidCtx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/"+feature.Uuid+"/phase", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should handle concurrent requests", func(t *testing.T) {
		numRequests := 5
		var wg sync.WaitGroup
		wg.Add(numRequests)

		for i := 0; i < numRequests; i++ {
			go func() {
				defer wg.Done()
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("feature_uuid", feature.Uuid)
				req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
					http.MethodGet, "/features/"+feature.Uuid+"/phase", nil)
				if err != nil {
					t.Error(err)
					return
				}

				rr := httptest.NewRecorder()
				http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

				assert.Equal(t, http.StatusOK, rr.Code)
			}()
		}

		wg.Wait()
	})

	t.Run("Valid Request with Authenticated User", func(t *testing.T) {
		db.CleanTestData()

		db.TestDB.CreateOrEditPerson(person)
		db.TestDB.CreateOrEditWorkspace(workspace)
		db.TestDB.CreateOrEditFeature(feature)
		db.TestDB.CreateOrEditFeaturePhase(featurePhase)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/"+feature.Uuid+"/phase", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

		var returnedFeaturePhases []db.FeaturePhase
		err = json.Unmarshal(rr.Body.Bytes(), &returnedFeaturePhases)
		assert.NoError(t, err)
		assert.NotEmpty(t, returnedFeaturePhases)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Large Number of Phases", func(t *testing.T) {

		db.CleanTestData()
		db.TestDB.CreateOrEditPerson(person)
		db.TestDB.CreateOrEditWorkspace(workspace)
		db.TestDB.CreateOrEditFeature(feature)

		for i := 0; i < 100; i++ {
			phase := db.FeaturePhase{
				Uuid:        fmt.Sprintf("large_phase_uuid_%d", i),
				FeatureUuid: feature.Uuid,
				Name:        fmt.Sprintf("Large Phase %d", i),
				Priority:    i,
			}
			db.TestDB.CreateOrEditFeaturePhase(phase)
		}

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/"+feature.Uuid+"/phase", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

		var returnedFeaturePhases []db.FeaturePhase
		err = json.Unmarshal(rr.Body.Bytes(), &returnedFeaturePhases)
		assert.NoError(t, err)
		assert.Equal(t, 100, len(returnedFeaturePhases))
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Missing Public Key in Context", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		emptyCtx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)
		req, err := http.NewRequestWithContext(emptyCtx,
			http.MethodGet, "/features/"+feature.Uuid+"/phase", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid Feature UUID Format", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", "invalid-uuid-format-123")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/invalid-uuid-format-123/phase", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

		var returnedFeaturePhases []db.FeaturePhase
		err = json.Unmarshal(rr.Body.Bytes(), &returnedFeaturePhases)
		assert.NoError(t, err)
		assert.Empty(t, returnedFeaturePhases)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Context Cancellation", func(t *testing.T) {

		ctxWithCancel, cancel := context.WithCancel(ctx)
		defer cancel()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctxWithCancel, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/"+feature.Uuid+"/phase", nil)
		if err != nil {
			t.Fatal(err)
		}

		cancel()

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Invalid Context Key Type", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)

		invalidCtx := context.WithValue(context.Background(), auth.ContextKey, 123)
		req, err := http.NewRequestWithContext(context.WithValue(invalidCtx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/"+feature.Uuid+"/phase", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetFeaturePhases).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestGetFeaturePhaseByUUID(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()

	dbHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        "test-person-uuid",
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        "test-workspace-uuid",
		Name:        "test-workspace",
		OwnerPubKey: person.OwnerPubKey,
		Github:      "test-github",
		Website:     "test-website",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	feature := db.WorkspaceFeatures{
		Uuid:          "test-feature-uuid",
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "test-url",
		Priority:      0,
	}
	db.TestDB.CreateOrEditFeature(feature)

	featurePhase := db.FeaturePhase{
		Uuid:         "test-feature-phase-uuid",
		FeatureUuid:  feature.Uuid,
		Name:         "test-phase",
		Priority:     0,
		PhasePurpose: "Initial test purpose",
		PhaseOutcome: "Expected initial outcome",
		PhaseScope:   "Initial scope",
	}
	db.TestDB.CreateOrEditFeaturePhase(featurePhase)

	fullFeaturePhase := db.FeaturePhase{
		Uuid:         "feature_phase_uuid_full_get",
		FeatureUuid:  feature.Uuid,
		Name:         "Full Feature Phase",
		Priority:     1,
		PhasePurpose: "Test phase purpose",
		PhaseOutcome: "Expected test outcome",
		PhaseScope:   "Test phase scope",
	}
	db.TestDB.CreateOrEditFeaturePhase(fullFeaturePhase)

	minimalFeaturePhase := db.FeaturePhase{
		Uuid:        "feature_phase_uuid_minimal_get",
		FeatureUuid: feature.Uuid,
		Name:        "Minimal Feature Phase",
		Priority:    2,
	}
	db.TestDB.CreateOrEditFeaturePhase(minimalFeaturePhase)

	ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

	tests := []struct {
		name           string
		featureUuid    string
		phaseUuid      string
		pubKeyFromAuth string
		expectedStatus int
		mockReturn     db.FeaturePhase
		mockError      error
		validateFunc   func(t *testing.T, body []byte)
	}{
		{
			name:           "Valid Request with Existing Feature and Phase UUIDs",
			featureUuid:    feature.Uuid,
			phaseUuid:      featurePhase.Uuid,
			pubKeyFromAuth: workspace.OwnerPubKey,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, body []byte) {
				var returnedFeaturePhase db.FeaturePhase
				assert.NoError(t, json.Unmarshal(body, &returnedFeaturePhase))
				assert.Equal(t, featurePhase.Uuid, returnedFeaturePhase.Uuid)
			},
		},
		{
			name:           "Unauthorized Request",
			featureUuid:    feature.Uuid,
			phaseUuid:      featurePhase.Uuid,
			pubKeyFromAuth: "",
			expectedStatus: http.StatusUnauthorized,
			validateFunc:   nil,
		},
		{
			name:           "Invalid Authentication Key",
			featureUuid:    feature.Uuid,
			phaseUuid:      featurePhase.Uuid,
			pubKeyFromAuth: "Invalid-Authentication-Key",
			expectedStatus: http.StatusUnauthorized,
			validateFunc:   nil,
		},
		{
			name:           "Feature UUID Does Not Exist",
			featureUuid:    "non-existent-feature-uuid",
			phaseUuid:      featurePhase.Uuid,
			pubKeyFromAuth: workspace.OwnerPubKey,
			expectedStatus: http.StatusNotFound,
			validateFunc:   nil,
		},
		{
			name:           "Phase UUID Does Not Exist",
			featureUuid:    feature.Uuid,
			phaseUuid:      "non-existent-phase-uuid",
			pubKeyFromAuth: workspace.OwnerPubKey,
			expectedStatus: http.StatusNotFound,
			validateFunc:   nil,
		},
		{
			name:           "Both Feature and Phase UUIDs Do Not Exist",
			featureUuid:    "non-existent-feature-uuid",
			phaseUuid:      "non-existent-phase-uuid",
			pubKeyFromAuth: workspace.OwnerPubKey,
			expectedStatus: http.StatusNotFound,
			validateFunc:   nil,
		},
		{
			name:           "Invalid UUID Format for Feature or Phase",
			featureUuid:    "invalid-format",
			phaseUuid:      "test-feature-phase-uuid",
			pubKeyFromAuth: workspace.OwnerPubKey,
			expectedStatus: http.StatusNotFound,
			validateFunc:   nil,
		},
		{
			name:           "Feature and Phase UUIDs Are the Same",
			featureUuid:    featurePhase.Uuid,
			phaseUuid:      featurePhase.Uuid,
			pubKeyFromAuth: workspace.OwnerPubKey,
			expectedStatus: http.StatusNotFound,
			validateFunc:   nil,
		},
		{
			name:           "Empty UUIDs",
			featureUuid:    "",
			phaseUuid:      "",
			pubKeyFromAuth: workspace.OwnerPubKey,
			expectedStatus: http.StatusNotFound,
			validateFunc:   nil,
		},
		{
			name:           "Should return feature phase with all new fields",
			featureUuid:    feature.Uuid,
			phaseUuid:      fullFeaturePhase.Uuid,
			pubKeyFromAuth: workspace.OwnerPubKey,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, body []byte) {
				var returnedFeaturePhase db.FeaturePhase
				assert.NoError(t, json.Unmarshal(body, &returnedFeaturePhase))

				assert.Equal(t, fullFeaturePhase.Uuid, returnedFeaturePhase.Uuid)
				assert.Equal(t, fullFeaturePhase.PhasePurpose, returnedFeaturePhase.PhasePurpose)
				assert.Equal(t, fullFeaturePhase.PhaseOutcome, returnedFeaturePhase.PhaseOutcome)
				assert.Equal(t, fullFeaturePhase.PhaseScope, returnedFeaturePhase.PhaseScope)
			},
		},
		{
			name:           "Should handle empty optional fields correctly",
			featureUuid:    feature.Uuid,
			phaseUuid:      minimalFeaturePhase.Uuid,
			pubKeyFromAuth: workspace.OwnerPubKey,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, body []byte) {
				var returnedFeaturePhase db.FeaturePhase
				assert.NoError(t, json.Unmarshal(body, &returnedFeaturePhase))

				assert.Equal(t, minimalFeaturePhase.Uuid, returnedFeaturePhase.Uuid)
				assert.Empty(t, returnedFeaturePhase.PhasePurpose)
				assert.Empty(t, returnedFeaturePhase.PhaseOutcome)
				assert.Empty(t, returnedFeaturePhase.PhaseScope)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("feature_uuid", tt.featureUuid)
			rctx.URLParams.Add("phase_uuid", tt.phaseUuid)

			req := httptest.NewRequest(http.MethodGet, "/features/"+tt.featureUuid+"/phase/"+tt.phaseUuid, nil)
			req = req.WithContext(context.WithValue(ctx, auth.ContextKey, tt.pubKeyFromAuth))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			http.HandlerFunc(dbHandler.GetFeaturePhaseByUUID).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.validateFunc != nil {
				tt.validateFunc(t, rr.Body.Bytes())
			}
		})
	}
	db.CleanTestData()
}

func TestDeleteFeaturePhase(t *testing.T) {
	tests := []struct {
		name           string
		featureUuid    string
		phaseUuid      string
		pubKeyFromAuth interface{}
		dbError        error
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "Valid Request with Existing Feature and Phase UUID",
			featureUuid:    "validFeatureUuid",
			phaseUuid:      "validPhaseUuid",
			pubKeyFromAuth: "validPubKey",
			dbError:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "Phase deleted successfully"},
		},
		{
			name:           "Valid Request with Non-Existing Feature UUID",
			featureUuid:    "nonExistingFeatureUuid",
			phaseUuid:      "validPhaseUuid",
			pubKeyFromAuth: "validPubKey",
			dbError:        errors.New("Feature not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"error": "no phase found to delete"},
		},
		{
			name:           "Valid Request with Non-Existing Phase UUID",
			featureUuid:    "validFeatureUuids",
			phaseUuid:      "nonExistingPhaseUuid",
			pubKeyFromAuth: "validPubKey",
			dbError:        errors.New("no phase found to delete"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"error": "no phase found to delete"},
		},
		{
			name:           "Valid Request with Both Non-Existing Feature and Phase UUIDs",
			featureUuid:    "nonExistingFeatureUuid",
			phaseUuid:      "nonExistingPhaseUuid",
			pubKeyFromAuth: "validPubKey",
			dbError:        errors.New("no phase found to delete"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"error": "no phase found to delete"},
		},
		{
			name:           "Missing Authorization Key",
			featureUuid:    "validFeatureUuid",
			phaseUuid:      "validPhaseUuid",
			pubKeyFromAuth: nil,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
		{
			name:           "Invalid Authorization Key Type",
			featureUuid:    "validFeatureUuid",
			phaseUuid:      "validPhaseUuid",
			pubKeyFromAuth: 12345, // invalid pub key
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
		{
			name:           "Feature UUID and Phase UUID as Empty Strings",
			featureUuid:    "",
			phaseUuid:      "",
			pubKeyFromAuth: "validPubKey",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Malformed UUIDs"},
		},
		{
			name:           "Feature UUID and Phase UUID with Special Characters",
			featureUuid:    "!@#$%",
			phaseUuid:      "^&*()",
			pubKeyFromAuth: "validPubKey",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Malformed UUIDs"},
		},
		{
			name:           "Malformed UUIDs",
			featureUuid:    "malformedFeatureUuid",
			phaseUuid:      "malformedPhaseUuid@",
			pubKeyFromAuth: "validPubKey",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Malformed UUIDs"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardownSuite := SetupSuite(t)
			defer teardownSuite(t)

			fHandler := NewFeatureHandler(db.TestDB)

			if tt.featureUuid == "validFeatureUuid" {
				person := db.Person{
					Uuid:        uuid.New().String(),
					OwnerAlias:  "test-alias",
					UniqueName:  "test-unique-name",
					OwnerPubKey: "validPubKey",
					PriceToMeet: 0,
					Description: "test-description",
				}
				db.TestDB.CreateOrEditPerson(person)

				workspace := db.Workspace{
					Uuid:        uuid.New().String(),
					Name:        "test-workspace" + uuid.New().String(),
					OwnerPubKey: person.OwnerPubKey,
					Github:      "https://github.com/test",
					Website:     "https://www.testwebsite.com",
					Description: "test-description",
				}
				db.TestDB.CreateOrEditWorkspace(workspace)
				workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

				feature := db.WorkspaceFeatures{
					Uuid:          tt.featureUuid,
					WorkspaceUuid: workspace.Uuid,
					Name:          "test-feature",
					Url:           "https://github.com/test-feature",
					Priority:      0,
				}
				db.TestDB.CreateOrEditFeature(feature)

				featurePhase := db.FeaturePhase{
					Uuid:        tt.phaseUuid,
					FeatureUuid: feature.Uuid,
					Name:        "test-feature-phase",
					Priority:    0,
				}
				db.TestDB.CreateOrEditFeaturePhase(featurePhase)
			}

			encodedFeatureUuid := url.QueryEscape(tt.featureUuid)
			encodedPhaseUuid := url.QueryEscape(tt.phaseUuid)

			req := httptest.NewRequest(http.MethodDelete, "/features/"+encodedFeatureUuid+"/phase/"+encodedPhaseUuid, nil)
			w := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), auth.ContextKey, tt.pubKeyFromAuth)
			req = req.WithContext(ctx)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("feature_uuid", encodedFeatureUuid)
			rctx.URLParams.Add("phase_uuid", encodedPhaseUuid)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			fHandler.DeleteFeaturePhase(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var responseBody map[string]string
				err := json.NewDecoder(w.Body).Decode(&responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}

func TestCreateOrEditStory(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	fHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)
	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "https://github.com/test-feature",
		Priority:      0,
	}
	db.TestDB.CreateOrEditFeature(feature)

	featureStory := db.FeatureStory{
		Uuid:        uuid.New().String(),
		FeatureUuid: feature.Uuid,
		Description: "test-description",
		Priority:    0,
	}

	t.Run("should return 401 error if not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditStory)

		requestBody, _ := json.Marshal(featureStory)
		req, err := http.NewRequest(http.MethodPost, "/features/story", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return 406 error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditStory)

		invalidJson := []byte(`{"key": "value"`)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/story", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("should successfully add feature story if request is valid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.CreateOrEditStory)

		requestBody, _ := json.Marshal(featureStory)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/features/story", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		createdStory, err := db.TestDB.GetFeatureStoryByUuid(featureStory.FeatureUuid, featureStory.Uuid)
		assert.NoError(t, err)
		assert.Equal(t, featureStory.Description, createdStory.Description)
		assert.Equal(t, featureStory.Priority, createdStory.Priority)
		assert.Equal(t, featureStory.FeatureUuid, createdStory.FeatureUuid)
	})
}

func TestGetStoriesByFeatureUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        "uuid",
		OwnerAlias:  "alias",
		UniqueName:  "unique_name",
		OwnerPubKey: "pubkey",
		PriceToMeet: 0,
		Description: "description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        "workspace_uuid",
		Name:        "workspace_name",
		OwnerPubKey: "person.OwnerPubkey",
		Github:      "gtihub",
		Website:     "website",
		Description: "description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)

	feature := db.WorkspaceFeatures{
		Uuid:          "feature_uuid",
		WorkspaceUuid: workspace.Uuid,
		Name:          "feature_name",
		Url:           "feature_url",
		Priority:      0,
	}
	db.TestDB.CreateOrEditFeature(feature)

	newStory := db.FeatureStory{
		Uuid:        "feature_story_uuid",
		FeatureUuid: feature.Uuid,
		Description: "feature_story_description",
		Priority:    0,
	}
	db.TestDB.CreateOrEditFeatureStory(newStory)

	ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

	t.Run("Should test that it throws a 401 error if a user is not authorized", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/"+feature.Uuid+"/story", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetStoriesByFeatureUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that the workspace features stories array returned from the API has the feature stories created", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/"+feature.Uuid+"/story", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetStoriesByFeatureUuid).ServeHTTP(rr, req)

		var returnedFeatureStory []db.FeatureStory
		err = json.Unmarshal(rr.Body.Bytes(), &returnedFeatureStory)
		assert.NoError(t, err)

		updatedFeatureStory, err := db.TestDB.GetFeatureStoriesByFeatureUuid(feature.Uuid)
		if err != nil {
			t.Fatal(err)
		}

		for i := range updatedFeatureStory {
			created := updatedFeatureStory[i].Created.In(time.UTC)
			updated := updatedFeatureStory[i].Updated.In(time.UTC)
			updatedFeatureStory[i].Created = &created
			updatedFeatureStory[i].Updated = &updated
		}

		assert.Equal(t, returnedFeatureStory, updatedFeatureStory)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestGetStoryByUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	fHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)
	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "https://github.com/test-feature",
		Priority:      0,
	}
	db.TestDB.CreateOrEditFeature(feature)

	featureStory := db.FeatureStory{
		Uuid:        uuid.New().String(),
		FeatureUuid: feature.Uuid,
		Description: "test-description",
		Priority:    0,
	}
	db.TestDB.CreateOrEditFeatureStory(featureStory)

	t.Run("should return 401 error if not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.GetStoryByUuid)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("story_uuid", featureStory.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/features/"+feature.Uuid+"/story/"+featureStory.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return feature story if user is authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.GetStoryByUuid)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("story_uuid", featureStory.Uuid)
		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/features/"+feature.Uuid+"/story/"+featureStory.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var returnedStory db.FeatureStory
		err = json.Unmarshal(rr.Body.Bytes(), &returnedStory)
		assert.NoError(t, err)
		assert.Equal(t, featureStory.Description, returnedStory.Description)
		assert.Equal(t, featureStory.Priority, returnedStory.Priority)
		assert.Equal(t, featureStory.FeatureUuid, returnedStory.FeatureUuid)
	})
}

func TestDeleteStory(t *testing.T) {

	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	fHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)
	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "https://github.com/test-feature",
		Priority:      0,
	}
	db.TestDB.CreateOrEditFeature(feature)

	featureStory := db.FeatureStory{
		Uuid:        uuid.New().String(),
		FeatureUuid: feature.Uuid,
		Description: "test-description",
		Priority:    0,
	}
	db.TestDB.CreateOrEditFeatureStory(featureStory)

	t.Run("should return 401 error if user not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		req, err := http.NewRequest(http.MethodDelete, "/"+feature.Uuid+"/story/"+featureStory.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should successfully delete feature story if request is valid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("story_uuid", featureStory.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodDelete, "/"+feature.Uuid+"/story/"+featureStory.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		deletedFeatureStory, _ := db.TestDB.GetFeatureStoryByUuid(feature.Uuid, featureStory.Uuid)
		assert.Equal(t, db.FeatureStory{}, deletedFeatureStory)

	})

	t.Run("should handle non-existent feature UUID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		nonExistentFeatureUUID := uuid.New().String()
		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", nonExistentFeatureUUID)
		rctx.URLParams.Add("story_uuid", featureStory.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodDelete,
			fmt.Sprintf("/%s/story/%s", nonExistentFeatureUUID, featureStory.Uuid),
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("should handle non-existent story UUID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		nonExistentStoryUUID := uuid.New().String()
		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("story_uuid", nonExistentStoryUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodDelete,
			fmt.Sprintf("/%s/story/%s", feature.Uuid, nonExistentStoryUUID),
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("should handle empty feature UUID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", "")
		rctx.URLParams.Add("story_uuid", featureStory.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodDelete,
			fmt.Sprintf("/%s/story/%s", "", featureStory.Uuid),
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("should handle empty story UUID", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("story_uuid", "")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodDelete,
			fmt.Sprintf("/%s/story/%s", feature.Uuid, ""),
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("should handle invalid UUID format", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", "invalid-uuid")
		rctx.URLParams.Add("story_uuid", "invalid-uuid")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodDelete,
			"/invalid-uuid/story/invalid-uuid",
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("should handle missing URL parameters", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(ctx,
			http.MethodDelete,
			"/story/",
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("should handle invalid auth token format", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		ctx := context.WithValue(context.Background(), auth.ContextKey, 12345)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("story_uuid", featureStory.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodDelete,
			fmt.Sprintf("/%s/story/%s", feature.Uuid, featureStory.Uuid),
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should handle concurrent delete requests", func(t *testing.T) {

		concurrentStory := db.FeatureStory{
			Uuid:        uuid.New().String(),
			FeatureUuid: feature.Uuid,
			Description: "concurrent-test-description",
			Priority:    0,
		}
		db.TestDB.CreateOrEditFeatureStory(concurrentStory)

		var wg sync.WaitGroup
		numRequests := 5
		wg.Add(numRequests)

		for i := 0; i < numRequests; i++ {
			go func() {
				defer wg.Done()
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(fHandler.DeleteStory)

				ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("feature_uuid", feature.Uuid)
				rctx.URLParams.Add("story_uuid", concurrentStory.Uuid)
				req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
					http.MethodDelete,
					fmt.Sprintf("/%s/story/%s", feature.Uuid, concurrentStory.Uuid),
					nil)
				if err != nil {
					t.Error(err)
					return
				}

				handler.ServeHTTP(rr, req)
				// First request should succeed, others should fail with NotFound
				assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, rr.Code)
			}()
		}

		wg.Wait()
	})

	t.Run("Valid Request with Existing Story", func(t *testing.T) {

		testStory := db.FeatureStory{
			Uuid:        uuid.New().String(),
			FeatureUuid: feature.Uuid,
			Description: "test-valid-story",
			Priority:    1,
		}
		db.TestDB.CreateOrEditFeatureStory(testStory)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("story_uuid", testStory.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodDelete,
			fmt.Sprintf("/%s/story/%s", feature.Uuid, testStory.Uuid),
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		deletedStory, _ := db.TestDB.GetFeatureStoryByUuid(feature.Uuid, testStory.Uuid)
		assert.Equal(t, db.FeatureStory{}, deletedStory)
	})

	t.Run("Valid Request with Non-Existing Story", func(t *testing.T) {
		nonExistentUUID := uuid.New().String()
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("story_uuid", nonExistentUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodDelete,
			fmt.Sprintf("/%s/story/%s", feature.Uuid, nonExistentUUID),
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Empty feature_uuid and story_uuid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", "")
		rctx.URLParams.Add("story_uuid", "")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodDelete,
			"/story/",
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Unauthorized Request with Valid UUIDs", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("story_uuid", featureStory.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodDelete,
			fmt.Sprintf("/%s/story/%s", feature.Uuid, featureStory.Uuid),
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid UUID Format for Both Parameters", func(t *testing.T) {
		invalidUUIDs := []string{
			"invalid-uuid",
			"123-456-789",
			"not-a-uuid-at-all",
			"12345",
			"",
		}

		for _, invalidUUID := range invalidUUIDs {
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(fHandler.DeleteStory)

			ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("feature_uuid", invalidUUID)
			rctx.URLParams.Add("story_uuid", invalidUUID)
			req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
				http.MethodDelete,
				fmt.Sprintf("/%s/story/%s", invalidUUID, invalidUUID),
				nil)
			assert.NoError(t, err)

			handler.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusNotFound, rr.Code)
		}
	})

	t.Run("Missing UUID Parameters in Context", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)

		req, err := http.NewRequestWithContext(ctx,
			http.MethodDelete,
			"/story/",
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Valid feature_uuid with Invalid story_uuid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("story_uuid", "invalid-story-uuid")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodDelete,
			fmt.Sprintf("/%s/story/invalid-story-uuid", feature.Uuid),
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Invalid feature_uuid with Valid story_uuid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", "invalid-feature-uuid")
		rctx.URLParams.Add("story_uuid", featureStory.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodDelete,
			fmt.Sprintf("/invalid-feature-uuid/story/%s", featureStory.Uuid),
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Case Sensitivity in UUIDs", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.DeleteStory)

		upperFeatureUUID := strings.ToUpper(feature.Uuid)
		upperStoryUUID := strings.ToUpper(featureStory.Uuid)

		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", upperFeatureUUID)
		rctx.URLParams.Add("story_uuid", upperStoryUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodDelete,
			fmt.Sprintf("/%s/story/%s", upperFeatureUUID, upperStoryUUID),
			nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

}

func TestGetBountiesByFeatureAndPhaseUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	fHandler := NewFeatureHandler(db.TestDB)

	db.CleanTestData()

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)
	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "https://github.com/test-feature",
		Priority:      0,
	}
	db.TestDB.CreateOrEditFeature(feature)

	featurePhase := db.FeaturePhase{
		Uuid:        uuid.New().String(),
		FeatureUuid: feature.Uuid,
		Name:        "test-feature-phase",
		Priority:    0,
	}
	db.TestDB.CreateOrEditFeaturePhase(featurePhase)

	bounty := db.NewBounty{
		OwnerID:       person.OwnerPubKey,
		WorkspaceUuid: workspace.Uuid,
		Title:         "test-bounty",
		PhaseUuid:     featurePhase.Uuid,
		Description:   "test-description",
		Price:         1000,
		Type:          "coding_task",
		Assignee:      "",
	}
	db.TestDB.CreateOrEditBounty(bounty)

	ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

	t.Run("should return 401 error if not authorized", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/features/"+feature.Uuid+"/phase/"+featurePhase.Uuid+"/bounty", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return the correct bounty if user is authorized", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/features/"+feature.Uuid+"/phase/"+featurePhase.Uuid+"/bounty", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		var returnedBounties []db.BountyResponse
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBounties)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 1, len(returnedBounties))
		assert.Equal(t, bounty.Title, returnedBounties[0].Bounty.Title)
		assert.Equal(t, bounty.Description, returnedBounties[0].Bounty.Description)
		assert.Equal(t, bounty.Price, returnedBounties[0].Bounty.Price)
		assert.Equal(t, bounty.PhaseUuid, returnedBounties[0].Bounty.PhaseUuid)
	})

	t.Run("should phase return the correct bounty response structure", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/features/"+feature.Uuid+"/phase/"+featurePhase.Uuid+"/bounty", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		var returnedBounties []db.BountyResponse
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBounties)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 1, len(returnedBounties))

		expectedBounty := db.BountyResponse{
			Bounty: db.NewBounty{
				ID:                      returnedBounties[0].Bounty.ID,
				OwnerID:                 person.OwnerPubKey,
				Paid:                    false,
				Show:                    false,
				Completed:               false,
				Type:                    "coding_task",
				Award:                   "",
				AssignedHours:           0,
				CommitmentFee:           0,
				Price:                   1000,
				Title:                   "test-bounty",
				Tribe:                   "",
				Assignee:                "",
				TicketUrl:               "",
				OrgUuid:                 workspace.Uuid,
				WorkspaceUuid:           workspace.Uuid,
				Description:             "test-description",
				WantedType:              "",
				Deliverables:            "",
				GithubDescription:       false,
				OneSentenceSummary:      "",
				EstimatedSessionLength:  "",
				EstimatedCompletionDate: "",
				Created:                 0,
				Updated:                 nil,
				PhaseUuid:               featurePhase.Uuid,
				PhasePriority:           0,
				PaymentPending:          false,
				PaymentFailed:           false,
			},
			Assignee: db.Person{},
			Owner: db.Person{
				ID:          returnedBounties[0].Owner.ID,
				Uuid:        person.Uuid,
				OwnerPubKey: person.OwnerPubKey,
				OwnerAlias:  person.OwnerAlias,
				UniqueName:  person.UniqueName,
				Description: person.Description,
				Tags:        pq.StringArray{},
				Img:         "",
			},
			Organization: db.WorkspaceShort{
				Uuid: workspace.Uuid,
				Name: workspace.Name,
				Img:  "",
			},
			Workspace: db.WorkspaceShort{
				Uuid: workspace.Uuid,
				Name: workspace.Name,
				Img:  "",
			},
		}

		assert.Equal(t, expectedBounty, returnedBounties[0])
	})

	t.Run("should return 404 if feature or phase UUID is invalid", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", "invalid-feature-uuid")
		rctx.URLParams.Add("phase_uuid", "invalid-phase-uuid")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/features/invalid-feature-uuid/phase/invalid-phase-uuid/bounty", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("should handle pagination correctly", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty?limit=1&offset=0", feature.Uuid, featurePhase.Uuid),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should handle search parameter", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty?search=test", feature.Uuid, featurePhase.Uuid),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should handle status filters", func(t *testing.T) {
		testCases := []struct {
			name   string
			query  string
			status string
		}{
			{"open bounties", "Open=true", ""},
			{"assigned bounties", "Assigned=true", person.OwnerPubKey},
			{"completed bounties", "Completed=true", person.OwnerPubKey},
			{"paid bounties", "Paid=true", person.OwnerPubKey},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				statusBounty := bounty
				statusBounty.Assignee = tc.status
				if tc.name == "completed bounties" {
					statusBounty.Completed = true
				}
				if tc.name == "paid bounties" {
					statusBounty.Paid = true
				}
				db.TestDB.CreateOrEditBounty(statusBounty)

				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("feature_uuid", feature.Uuid)
				rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
				req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
					http.MethodGet,
					fmt.Sprintf("/features/%s/phase/%s/bounty?%s", feature.Uuid, featurePhase.Uuid, tc.query),
					nil)
				assert.NoError(t, err)

				rr := httptest.NewRecorder()
				http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

				assert.Equal(t, http.StatusOK, rr.Code)
			})
		}
	})

	t.Run("should handle language filters", func(t *testing.T) {

		langBounty := bounty
		langBounty.CodingLanguages = pq.StringArray{"golang", "javascript"}
		db.TestDB.CreateOrEditBounty(langBounty)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty?languages=golang", feature.Uuid, featurePhase.Uuid),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should handle tags filter", func(t *testing.T) {

		taggedBounty := db.NewBounty{
			OwnerID:       person.OwnerPubKey,
			WorkspaceUuid: workspace.Uuid,
			Title:         "tagged-test-bounty",
			PhaseUuid:     featurePhase.Uuid,
			Description:   "test-description-with-tags",
			Price:         1000,
			Type:          "coding_task",
			Assignee:      "",
		}
		db.TestDB.CreateOrEditBounty(taggedBounty)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty?tags=urgent", feature.Uuid, featurePhase.Uuid),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should handle sorting parameters", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty?sortBy=price&direction=DESC", feature.Uuid, featurePhase.Uuid),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should handle multiple filters together", func(t *testing.T) {

		multiBounty := db.NewBounty{
			ID:              1,
			OwnerID:         person.OwnerPubKey,
			WorkspaceUuid:   workspace.Uuid,
			Title:           "multi-filter-test-bounty",
			PhaseUuid:       featurePhase.Uuid,
			FeatureUuid:     feature.Uuid,
			Description:     "test-description-multi-filter",
			Price:           1000,
			Type:            "coding_task",
			Assignee:        "",
			CodingLanguages: pq.StringArray{"golang", "javascript"},
			Paid:            false,
			Completed:       false,
		}
		_, err := db.TestDB.CreateOrEditBounty(multiBounty)
		assert.NoError(t, err)

		bounties, err := db.TestDB.GetBountiesByFeatureAndPhaseUuid(feature.Uuid, featurePhase.Uuid, &http.Request{URL: &url.URL{RawQuery: ""}})
		assert.NoError(t, err)
		assert.NotEmpty(t, bounties)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)

		reqURL := fmt.Sprintf("/features/%s/phase/%s/bounty?Open=true&languages=golang&sortBy=price&direction=DESC",
			feature.Uuid, featurePhase.Uuid)
		req, err := http.NewRequestWithContext(
			context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			reqURL,
			nil,
		)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var returnedBounties []db.BountyResponse
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBounties)
		assert.NoError(t, err)
		assert.NotEmpty(t, returnedBounties)
		assert.Equal(t, multiBounty.Title, returnedBounties[0].Bounty.Title)
	})

	t.Run("Valid Feature and Phase UUIDs with No Bounties", func(t *testing.T) {
		emptyFeature := db.WorkspaceFeatures{
			Uuid:          uuid.New().String(),
			WorkspaceUuid: workspace.Uuid,
			Name:          "empty-feature",
		}
		db.TestDB.CreateOrEditFeature(emptyFeature)

		emptyPhase := db.FeaturePhase{
			Uuid:        uuid.New().String(),
			FeatureUuid: emptyFeature.Uuid,
			Name:        "empty-phase",
		}
		db.TestDB.CreateOrEditFeaturePhase(emptyPhase)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", emptyFeature.Uuid)
		rctx.URLParams.Add("phase_uuid", emptyPhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty", emptyFeature.Uuid, emptyPhase.Uuid),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Empty Feature UUID", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", "")
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty", "", featurePhase.Uuid),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Empty Phase UUID", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", "")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty", feature.Uuid, ""),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Invalid UUID Format", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", "invalid-uuid-format")
		rctx.URLParams.Add("phase_uuid", "invalid-uuid-format")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			"/features/invalid-uuid-format/phase/invalid-uuid-format/bounty",
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Feature and Phase UUIDs with Special Characters", func(t *testing.T) {
		specialFeatureUUID := "test!@#$%^&*()"
		specialPhaseUUID := "phase!@#$%^&*()"

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", specialFeatureUUID)
		rctx.URLParams.Add("phase_uuid", specialPhaseUUID)

		encodedFeatureUUID := url.QueryEscape(specialFeatureUUID)
		encodedPhaseUUID := url.QueryEscape(specialPhaseUUID)

		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty", encodedFeatureUUID, encodedPhaseUUID),
			nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Null Feature and Phase UUIDs", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			"/features/null/phase/null/bounty",
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
	t.Run("Valid Request with Bounties", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty", feature.Uuid, featurePhase.Uuid),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		var returnedBounties []db.BountyResponse
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBounties)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotEmpty(t, returnedBounties)
	})

	t.Run("Non-Existent Feature UUID", func(t *testing.T) {
		nonExistentUUID := uuid.New().String()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", nonExistentUUID)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty", nonExistentUUID, featurePhase.Uuid),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Non-Existent Phase UUID", func(t *testing.T) {
		nonExistentUUID := uuid.New().String()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", nonExistentUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty", feature.Uuid, nonExistentUUID),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Valid UUIDs but Different Case", func(t *testing.T) {
		upperCaseFeatureUUID := strings.ToUpper(feature.Uuid)
		upperCasePhaseUUID := strings.ToUpper(featurePhase.Uuid)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", upperCaseFeatureUUID)
		rctx.URLParams.Add("phase_uuid", upperCasePhaseUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty", upperCaseFeatureUUID, upperCasePhaseUUID),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Malformed UUIDs", func(t *testing.T) {
		malformedUUIDs := []string{
			"not-a-uuid",
			"123e4567-e89b-12d3-a456",
			"123e4567-e89b-12d3-a456-426614174000-extra",
			"123e4567-e89b-12d3-a456-42661417400g",
		}

		for _, malformedUUID := range malformedUUIDs {
			t.Run(fmt.Sprintf("Malformed UUID: %s", malformedUUID), func(t *testing.T) {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("feature_uuid", malformedUUID)
				rctx.URLParams.Add("phase_uuid", malformedUUID)
				req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
					http.MethodGet,
					fmt.Sprintf("/features/%s/phase/%s/bounty", malformedUUID, malformedUUID),
					nil)
				assert.NoError(t, err)

				rr := httptest.NewRecorder()
				http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

				assert.Equal(t, http.StatusNotFound, rr.Code)
			})
		}
	})

	t.Run("Valid Request with No Bounties", func(t *testing.T) {

		emptyFeature := db.WorkspaceFeatures{
			Uuid:          uuid.New().String(),
			WorkspaceUuid: workspace.Uuid,
			Name:          "feature-without-bounties",
		}
		db.TestDB.CreateOrEditFeature(emptyFeature)

		emptyPhase := db.FeaturePhase{
			Uuid:        uuid.New().String(),
			FeatureUuid: emptyFeature.Uuid,
			Name:        "phase-without-bounties",
		}
		db.TestDB.CreateOrEditFeaturePhase(emptyPhase)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", emptyFeature.Uuid)
		rctx.URLParams.Add("phase_uuid", emptyPhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty", emptyFeature.Uuid, emptyPhase.Uuid),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Empty Feature UUID", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", "")
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			"/features//phase/"+featurePhase.Uuid+"/bounty",
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Empty Phase UUID", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", "")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet,
			"/features/"+feature.Uuid+"/phase//bounty",
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Unauthorized Access", func(t *testing.T) {

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
			http.MethodGet,
			fmt.Sprintf("/features/%s/phase/%s/bounty", feature.Uuid, featurePhase.Uuid),
			nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid UUID Format", func(t *testing.T) {
		invalidUUIDs := []struct {
			name        string
			featureUUID string
			phaseUUID   string
		}{
			{
				name:        "Invalid Feature UUID",
				featureUUID: "invalid-uuid",
				phaseUUID:   featurePhase.Uuid,
			},
			{
				name:        "Invalid Phase UUID",
				featureUUID: feature.Uuid,
				phaseUUID:   "invalid-uuid",
			},
			{
				name:        "Both Invalid UUIDs",
				featureUUID: "invalid-uuid-1",
				phaseUUID:   "invalid-uuid-2",
			},
			{
				name:        "UUID with Special Characters",
				featureUUID: "123e4567-e89b-12d3-a456-426614174000!",
				phaseUUID:   "123e4567-e89b-12d3-a456-426614174000@",
			},
		}

		for _, tc := range invalidUUIDs {
			t.Run(tc.name, func(t *testing.T) {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("feature_uuid", tc.featureUUID)
				rctx.URLParams.Add("phase_uuid", tc.phaseUUID)
				req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
					http.MethodGet,
					fmt.Sprintf("/features/%s/phase/%s/bounty", tc.featureUUID, tc.phaseUUID),
					nil)
				assert.NoError(t, err)

				rr := httptest.NewRecorder()
				http.HandlerFunc(fHandler.GetBountiesByFeatureAndPhaseUuid).ServeHTTP(rr, req)

				assert.Equal(t, http.StatusNotFound, rr.Code)
			})
		}
	})
}

func TestGetBountiesCountByFeatureAndPhaseUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	fHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)
	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature",
		Url:           "https://github.com/test-feature",
		Priority:      0,
	}
	db.TestDB.CreateOrEditFeature(feature)

	featurePhase := db.FeaturePhase{
		Uuid:        uuid.New().String(),
		FeatureUuid: feature.Uuid,
		Name:        "test-feature-phase",
		Priority:    0,
	}
	db.TestDB.CreateOrEditFeaturePhase(featurePhase)

	bounty := db.NewBounty{
		OwnerID:       person.OwnerPubKey,
		WorkspaceUuid: workspace.Uuid,
		Title:         "test-bounty",
		PhaseUuid:     featurePhase.Uuid,
		Description:   "test-description",
		Price:         1000,
		Type:          "coding_task",
		Assignee:      "",
	}
	db.TestDB.CreateOrEditBounty(bounty)

	ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

	t.Run("should return 401 error if not authorized", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/features/"+feature.Uuid+"/phase/"+featurePhase.Uuid+"/bounty/count", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return the correct bounty count if user is authorized", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/features/"+feature.Uuid+"/phase/"+featurePhase.Uuid+"/bounty/count", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		var returnedBountiesCount int64
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBountiesCount)
		assert.NoError(t, err)

		bountiesCount := db.TestDB.GetBountiesCountByFeatureAndPhaseUuid(feature.Uuid, featurePhase.Uuid, req)

		assert.Equal(t, returnedBountiesCount, bountiesCount)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should handle invalid feature UUID format", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", "invalid-uuid-format")
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/invalid-uuid/phase/"+featurePhase.Uuid+"/bounty/count", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		var returnedBountiesCount int64
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBountiesCount)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), returnedBountiesCount)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should handle invalid phase UUID format", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", "invalid-phase-uuid")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/"+feature.Uuid+"/phase/invalid-phase-uuid/bounty/count", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		var returnedBountiesCount int64
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBountiesCount)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), returnedBountiesCount)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should handle missing feature UUID parameter", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features//phase/"+featurePhase.Uuid+"/bounty/count", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		var returnedBountiesCount int64
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBountiesCount)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), returnedBountiesCount)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should handle missing phase UUID parameter", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/"+feature.Uuid+"/phase//bounty/count", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		var returnedBountiesCount int64
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBountiesCount)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), returnedBountiesCount)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should handle all status filters", func(t *testing.T) {
		statuses := []string{"Open", "Assigned", "Completed", "Paid"}
		for _, status := range statuses {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("feature_uuid", feature.Uuid)
			rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
			req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
				http.MethodGet, fmt.Sprintf("/features/%s/phase/%s/bounty/count?%s=true",
					feature.Uuid, featurePhase.Uuid, status), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
		}
	})

	t.Run("should handle multiple status filters simultaneously", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, fmt.Sprintf("/features/%s/phase/%s/bounty/count?Open=true&Assigned=true",
				feature.Uuid, featurePhase.Uuid), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should handle non-existent feature", func(t *testing.T) {
		nonExistentUUID := uuid.New().String()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", nonExistentUUID)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/"+nonExistentUUID+"/phase/"+featurePhase.Uuid+"/bounty/count", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		var returnedBountiesCount int64
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBountiesCount)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), returnedBountiesCount)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should handle invalid status filter values", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, fmt.Sprintf("/features/%s/phase/%s/bounty/count?Open=invalid",
				feature.Uuid, featurePhase.Uuid), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should handle concurrent requests", func(t *testing.T) {
		numRequests := 5
		var wg sync.WaitGroup
		wg.Add(numRequests)

		for i := 0; i < numRequests; i++ {
			go func() {
				defer wg.Done()
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("feature_uuid", feature.Uuid)
				rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
				req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
					http.MethodGet, "/features/"+feature.Uuid+"/phase/"+featurePhase.Uuid+"/bounty/count", nil)
				if err != nil {
					t.Error(err)
					return
				}

				rr := httptest.NewRecorder()
				http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

				assert.Equal(t, http.StatusOK, rr.Code)
			}()
		}

		wg.Wait()
	})

	t.Run("Empty Feature UUID", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", "")
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features//phase/"+featurePhase.Uuid+"/bounty/count", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		var returnedBountiesCount int64
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBountiesCount)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), returnedBountiesCount)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Empty Phase UUID", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		rctx.URLParams.Add("phase_uuid", "")
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/"+feature.Uuid+"/phase//bounty/count", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		var returnedBountiesCount int64
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBountiesCount)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), returnedBountiesCount)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Non-Existent Feature and Phase UUIDs", func(t *testing.T) {
		nonExistentFeatureUUID := uuid.New().String()
		nonExistentPhaseUUID := uuid.New().String()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", nonExistentFeatureUUID)
		rctx.URLParams.Add("phase_uuid", nonExistentPhaseUUID)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, "/features/"+nonExistentFeatureUUID+"/phase/"+nonExistentPhaseUUID+"/bounty/count", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		var returnedBountiesCount int64
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBountiesCount)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), returnedBountiesCount)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Feature and Phase UUIDs with Special Characters", func(t *testing.T) {
		specialChars := "!@#$%^&*()"
		encodedSpecialChars := url.QueryEscape(specialChars)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("feature_uuid", specialChars)
		rctx.URLParams.Add("phase_uuid", specialChars)

		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodGet, fmt.Sprintf("/features/%s/phase/%s/bounty/count",
				encodedSpecialChars, encodedSpecialChars), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetBountiesCountByFeatureAndPhaseUuid).ServeHTTP(rr, req)

		var returnedBountiesCount int64
		err = json.Unmarshal(rr.Body.Bytes(), &returnedBountiesCount)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), returnedBountiesCount)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

}

func TestGetFeatureStories(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()

	fHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-get-feature-stories-alias",
		UniqueName:  "test-get-feature-stories-unique-name",
		OwnerPubKey: "test-get-feature-stories-pubkey",
		PriceToMeet: 0,
		Description: "test-get-feature-stories-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-get-feature-stories-workspace-name",
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-get-feature-stories-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)
	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

	feature := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-get-feature-stories-feature-name",
		Url:           "https://github.com/test-get-feature-stories-feature-url",
		Priority:      0,
	}

	feature2 := db.WorkspaceFeatures{
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-get-feature-stories-feature-name-2",
		Url:           "https://github.com/test-get-feature-stories-feature-url-2",
		Priority:      0,
	}

	db.TestDB.CreateOrEditFeature(feature)
	db.TestDB.CreateOrEditFeature(feature2)

	ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

	story := db.FeatureStories{
		UserStory: "This is a test user story",
		Rationale: "This is a test rationale",
		Order:     1,
	}

	story2 := db.FeatureStories{
		UserStory: "This is a test user story 2",
		Rationale: "This is a test rationale 2",
		Order:     2,
	}

	story3 := db.FeatureStories{
		UserStory: "This is a test user story 3",
		Rationale: "This is a test rationale 3",
		Order:     3,
	}

	story4 := db.FeatureStories{
		UserStory: "This is a test user story 4",
		Rationale: "This is a test rationale 4",
		Order:     4,
	}

	story5 := db.FeatureStories{
		UserStory: "This is a test user story 5",
		Rationale: "This is a test rationale 5",
		Order:     5,
	}

	story6 := db.FeatureStories{
		UserStory: "This is a test user story 6",
		Rationale: "This is a test rationale 6",
		Order:     6,
	}

	stories := []db.FeatureStories{
		story,
		story2,
		story3,
	}

	stories2 := []db.FeatureStories{
		story4,
		story5,
		story6,
	}

	featureStories := db.FeatureStoriesReponse{
		Output: db.FeatureOutput{
			FeatureUuid:    feature.Uuid,
			FeatureContext: "Feature Context",
			Stories:        stories,
		},
	}

	featureStories2 := db.FeatureStoriesReponse{
		Output: db.FeatureOutput{
			FeatureUuid:    "Fake-feature-uuid",
			FeatureContext: "Feature Context",
			Stories:        stories2,
		},
	}

	requestBody, _ := json.Marshal(featureStories)
	requestBody2, _ := json.Marshal(featureStories2)

	t.Run("Should add user stories from stakwork to the feature stories table", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/features/stories", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		var featureStoriesReponse string
		err = json.Unmarshal(rr.Body.Bytes(), &featureStoriesReponse)
		assert.NoError(t, err)

		featureStories, _ := db.TestDB.GetFeatureStoriesByFeatureUuid(feature.Uuid)
		featureStoriesCount := len(featureStories)

		assert.Equal(t, int64(featureStoriesCount), int64(3))
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should not add user stories from stakwork to the feature stories table if the feature uuid is not found", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodPost, "/features/stories", bytes.NewReader(requestBody2))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		var featureStoriesReponse string
		err = json.Unmarshal(rr.Body.Bytes(), &featureStoriesReponse)
		assert.NoError(t, err)

		featureStories, _ := db.TestDB.GetFeatureStoriesByFeatureUuid(feature2.Uuid)
		featureStoriesCount := len(featureStories)

		assert.Equal(t, int64(featureStoriesCount), int64(0))
		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("Should not add the user stories if request body is empty", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodPost, "/features/stories", bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("Should not add user stories if request body is malformed JSON", func(t *testing.T) {
		malformedJSON := []byte(`{"invalid json`)
		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodPost, "/features/stories", bytes.NewReader(malformedJSON))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("Should not add user stories if stories array is empty", func(t *testing.T) {

		db.CleanTestData()
		db.TestDB.CreateOrEditPerson(person)
		db.TestDB.CreateOrEditWorkspace(workspace)
		db.TestDB.CreateOrEditFeature(feature)

		emptyStories := db.FeatureStoriesReponse{
			Output: db.FeatureOutput{
				FeatureUuid:    feature.Uuid,
				FeatureContext: "Feature Context",
				Stories:        []db.FeatureStories{},
			},
		}
		requestBody, _ := json.Marshal(emptyStories)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodPost, "/features/stories", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		featureStories, _ := db.TestDB.GetFeatureStoriesByFeatureUuid(feature.Uuid)
		assert.Equal(t, 0, len(featureStories))
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should handle invalid feature UUID format", func(t *testing.T) {
		invalidUUIDStories := db.FeatureStoriesReponse{
			Output: db.FeatureOutput{
				FeatureUuid:    "invalid-uuid-format",
				FeatureContext: "Feature Context",
				Stories:        stories,
			},
		}
		requestBody, _ := json.Marshal(invalidUUIDStories)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodPost, "/features/stories", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("Should handle maximum story length", func(t *testing.T) {
		longStory := db.FeatureStories{
			UserStory: strings.Repeat("a", 10000),
			Rationale: "Test rationale",
			Order:     1,
		}
		longStories := db.FeatureStoriesReponse{
			Output: db.FeatureOutput{
				FeatureUuid:    feature.Uuid,
				FeatureContext: "Feature Context",
				Stories:        []db.FeatureStories{longStory},
			},
		}
		requestBody, _ := json.Marshal(longStories)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodPost, "/features/stories", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		savedStories, _ := db.TestDB.GetFeatureStoriesByFeatureUuid(feature.Uuid)
		assert.Greater(t, len(savedStories), 0)
	})

	t.Run("Should handle concurrent requests", func(t *testing.T) {
		numRequests := 5
		var wg sync.WaitGroup
		wg.Add(numRequests)

		for i := 0; i < numRequests; i++ {
			go func() {
				defer wg.Done()
				rctx := chi.NewRouteContext()
				req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
					http.MethodPost, "/features/stories", bytes.NewReader(requestBody))
				if err != nil {
					t.Error(err)
					return
				}

				rr := httptest.NewRecorder()
				http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

				assert.Equal(t, http.StatusOK, rr.Code)
			}()
		}

		wg.Wait()

		savedStories, _ := db.TestDB.GetFeatureStoriesByFeatureUuid(feature.Uuid)
		assert.Greater(t, len(savedStories), 0)
	})

	t.Run("Missing Feature UUID", func(t *testing.T) {
		missingUUIDStories := db.FeatureStoriesReponse{
			Output: db.FeatureOutput{
				FeatureContext: "Feature Context",
				Stories:        stories,
			},
		}
		requestBody, _ := json.Marshal(missingUUIDStories)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodPost, "/features/stories", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("Non-String Feature UUID", func(t *testing.T) {
		nonStringUUID := []byte(`{"Output":{"FeatureUuid":123,"FeatureContext":"Test","Stories":[]}}`)
		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodPost, "/features/stories", bytes.NewReader(nonStringUUID))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("Large Number of Stories", func(t *testing.T) {
		largeStories := make([]db.FeatureStories, 1000)
		for i := 0; i < 1000; i++ {
			largeStories[i] = db.FeatureStories{
				UserStory: fmt.Sprintf("Story %d", i),
				Rationale: fmt.Sprintf("Rationale %d", i),
				Order:     uint(i + 1),
			}
		}

		largeStoriesReq := db.FeatureStoriesReponse{
			Output: db.FeatureOutput{
				FeatureUuid:    feature.Uuid,
				FeatureContext: "Feature Context",
				Stories:        largeStories,
			},
		}
		requestBody, _ := json.Marshal(largeStoriesReq)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodPost, "/features/stories", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Feature UUID with Special Characters", func(t *testing.T) {
		specialCharUUID := db.FeatureStoriesReponse{
			Output: db.FeatureOutput{
				FeatureUuid:    "!@#$%^&*()",
				FeatureContext: "Feature Context",
				Stories:        stories,
			},
		}
		requestBody, _ := json.Marshal(specialCharUUID)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodPost, "/features/stories", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("Duplicate Stories", func(t *testing.T) {
		db.CleanTestData()
		db.TestDB.CreateOrEditPerson(person)
		db.TestDB.CreateOrEditWorkspace(workspace)

		feature := db.WorkspaceFeatures{
			Uuid:          uuid.New().String(),
			WorkspaceUuid: workspace.Uuid,
			Name:          "test-get-feature-stories-duplicate-content",
			Url:           "https://github.com/test-get-feature-stories-feature-url",
			Priority:      0,
		}

		db.TestDB.CreateOrEditFeature(feature)

		duplicateStory := db.FeatureStories{
			UserStory: "Duplicate story content",
			Rationale: "Duplicate rationale content",
			Order:     1,
		}
		duplicateStories := []db.FeatureStories{duplicateStory, duplicateStory, duplicateStory}

		duplicateStoriesReq := db.FeatureStoriesReponse{
			Output: db.FeatureOutput{
				FeatureUuid:    feature.Uuid,
				FeatureContext: "Feature Context",
				Stories:        duplicateStories,
			},
		}
		requestBody, _ := json.Marshal(duplicateStoriesReq)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodPost, "/features/stories", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		savedStories, _ := db.TestDB.GetFeatureStoriesByFeatureUuid(feature.Uuid)
		assert.Equal(t, len(duplicateStories), len(savedStories))

		for _, story := range savedStories {
			assert.Equal(t, "Duplicate story content", story.Description)
		}
	})

	t.Run("Feature UUID as Null", func(t *testing.T) {
		nullUUID := []byte(`{"Output":{"FeatureUuid":null,"FeatureContext":"Test","Stories":[]}}`)
		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodPost, "/features/stories", bytes.NewReader(nullUUID))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("Empty Feature UUID String", func(t *testing.T) {
		emptyUUIDStories := db.FeatureStoriesReponse{
			Output: db.FeatureOutput{
				FeatureUuid:    "",
				FeatureContext: "Feature Context",
				Stories:        stories,
			},
		}
		requestBody, _ := json.Marshal(emptyUUIDStories)

		rctx := chi.NewRouteContext()
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx),
			http.MethodPost, "/features/stories", bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeatureStories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})
}

func TestGetFeaturesByWorkspaceUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewWorkspaceHandler(db.TestDB)

	db.CleanTestData()

	tests := []struct {
		name           string
		pubKeyFromAuth string
		workspaceUUID  string
		setupMocks     func() (string, interface{})
		expectedStatus int
	}{
		{
			name:           "should return error if a user is not authorized",
			pubKeyFromAuth: "",
			workspaceUUID:  "valid-uuid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing UUID Parameter",
			pubKeyFromAuth: "validPubKey",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Special Characters in UUID",
			pubKeyFromAuth: "validPubKey",
			workspaceUUID:  "1234-5678-!@#$",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Non-Existent UUID",
			pubKeyFromAuth: "validPubKey",
			workspaceUUID:  "nonexistentuuid",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "created feature should be present in the returned array",
			pubKeyFromAuth: "validPubKey",
			workspaceUUID:  "valid-uuid",
			setupMocks: func() (string, interface{}) {
				person := db.Person{
					Uuid:        uuid.New().String(),
					OwnerAlias:  "alias",
					UniqueName:  "unique_name",
					OwnerPubKey: "validPubKey",
					PriceToMeet: 0,
					Description: "description",
				}
				db.TestDB.CreateOrEditPerson(person)

				workspace := db.Workspace{
					Uuid:        uuid.New().String(),
					Name:        "unique_workspace_name",
					OwnerPubKey: person.OwnerPubKey,
					Github:      "github",
					Website:     "website",
					Description: "description",
				}
				db.TestDB.CreateOrEditWorkspace(workspace)

				feature := db.WorkspaceFeatures{
					Uuid:          "mock-feature-uuid1",
					WorkspaceUuid: workspace.Uuid,
					Name:          "feature_name",
					Url:           "https://www.bountieswebsite.com",
					Priority:      0,
				}
				db.TestDB.CreateOrEditFeature(feature)

				return workspace.Uuid, []db.WorkspaceFeatures{
					{
						Uuid:     "mock-feature-uuid1",
						Name:     "feature_name",
						Url:      "https://www.bountieswebsite.com",
						Priority: 0,
					},
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Large Number of Features",
			pubKeyFromAuth: "validPubKey",
			workspaceUUID:  "valid-uuid",
			setupMocks: func() (string, interface{}) {
				person := db.Person{
					Uuid:        uuid.New().String(),
					OwnerAlias:  "alias",
					UniqueName:  "unique_name",
					OwnerPubKey: "validPubKey",
					PriceToMeet: 0,
					Description: "description",
				}
				db.TestDB.CreateOrEditPerson(person)

				workspace := db.Workspace{
					Uuid:        uuid.New().String(),
					Name:        "unique_workspace_name",
					OwnerPubKey: person.OwnerPubKey,
					Github:      "github",
					Website:     "website",
					Description: "description",
				}
				db.TestDB.CreateOrEditWorkspace(workspace)

				var features []db.WorkspaceFeatures
				for i := 0; i < 1000; i++ {
					feature := db.WorkspaceFeatures{
						Uuid:          uuid.New().String(),
						WorkspaceUuid: workspace.Uuid,
						Name:          fmt.Sprintf("feature_%d", i),
						Url:           fmt.Sprintf("https://example.com/feature_%d", i),
						Priority:      i,
					}
					db.TestDB.CreateOrEditFeature(feature)
					features = append(features, feature)
				}

				return workspace.Uuid, features
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty Workspace with Status Filter",
			pubKeyFromAuth: "validPubKey",
			workspaceUUID:  "valid-uuid",
			setupMocks: func() (string, interface{}) {
				person := db.Person{
					Uuid:        uuid.New().String(),
					OwnerPubKey: "validPubKey",
				}
				db.TestDB.CreateOrEditPerson(person)

				workspace := db.Workspace{
					Uuid:        uuid.New().String(),
					OwnerPubKey: person.OwnerPubKey,
				}
				db.TestDB.CreateOrEditWorkspace(workspace)

				return workspace.Uuid, []db.WorkspaceFeatures{}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Features with Different Statuses",
			pubKeyFromAuth: "validPubKey",
			workspaceUUID:  "valid-uuid",
			setupMocks: func() (string, interface{}) {
				person := db.Person{
					Uuid:        uuid.New().String(),
					OwnerPubKey: "validPubKey",
				}
				db.TestDB.CreateOrEditPerson(person)

				workspace := db.Workspace{
					Uuid:        uuid.New().String(),
					OwnerPubKey: person.OwnerPubKey,
				}
				db.TestDB.CreateOrEditWorkspace(workspace)

				activeFeature := db.WorkspaceFeatures{
					Uuid:          uuid.New().String(),
					WorkspaceUuid: workspace.Uuid,
					Name:          "Active Feature",
					FeatStatus:    db.ActiveFeature,
				}
				db.TestDB.CreateOrEditFeature(activeFeature)

				archivedFeature := db.WorkspaceFeatures{
					Uuid:          uuid.New().String(),
					WorkspaceUuid: workspace.Uuid,
					Name:          "Archived Feature",
					FeatStatus:    db.ArchivedFeature,
				}
				db.TestDB.CreateOrEditFeature(archivedFeature)

				return workspace.Uuid, []db.WorkspaceFeatures{activeFeature}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Features with Custom Sorting",
			pubKeyFromAuth: "validPubKey",
			workspaceUUID:  "valid-uuid",
			setupMocks: func() (string, interface{}) {
				person := db.Person{
					Uuid:        uuid.New().String(),
					OwnerPubKey: "validPubKey",
				}
				db.TestDB.CreateOrEditPerson(person)

				workspace := db.Workspace{
					Uuid:        uuid.New().String(),
					OwnerPubKey: person.OwnerPubKey,
				}
				db.TestDB.CreateOrEditWorkspace(workspace)

				features := []db.WorkspaceFeatures{
					{
						Uuid:          uuid.New().String(),
						WorkspaceUuid: workspace.Uuid,
						Name:          "Feature A",
						Priority:      3,
					},
					{
						Uuid:          uuid.New().String(),
						WorkspaceUuid: workspace.Uuid,
						Name:          "Feature B",
						Priority:      1,
					},
					{
						Uuid:          uuid.New().String(),
						WorkspaceUuid: workspace.Uuid,
						Name:          "Feature C",
						Priority:      2,
					},
				}

				for _, f := range features {
					db.TestDB.CreateOrEditFeature(f)
				}

				return workspace.Uuid, features
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Features with Pagination",
			pubKeyFromAuth: "validPubKey",
			workspaceUUID:  "valid-uuid",
			setupMocks: func() (string, interface{}) {
				person := db.Person{
					Uuid:        uuid.New().String(),
					OwnerPubKey: "validPubKey",
				}
				db.TestDB.CreateOrEditPerson(person)

				workspace := db.Workspace{
					Uuid:        uuid.New().String(),
					OwnerPubKey: person.OwnerPubKey,
				}
				db.TestDB.CreateOrEditWorkspace(workspace)

				var features []db.WorkspaceFeatures
				for i := 0; i < 25; i++ {
					feature := db.WorkspaceFeatures{
						Uuid:          uuid.New().String(),
						WorkspaceUuid: workspace.Uuid,
						Name:          fmt.Sprintf("Feature %d", i),
						Priority:      i,
					}
					db.TestDB.CreateOrEditFeature(feature)
					features = append(features, feature)
				}

				return workspace.Uuid, features
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Features with Special Characters in Names",
			pubKeyFromAuth: "validPubKey",
			workspaceUUID:  "valid-uuid",
			setupMocks: func() (string, interface{}) {
				person := db.Person{
					Uuid:        uuid.New().String(),
					OwnerPubKey: "validPubKey",
				}
				db.TestDB.CreateOrEditPerson(person)

				workspace := db.Workspace{
					Uuid:        uuid.New().String(),
					OwnerPubKey: person.OwnerPubKey,
				}
				db.TestDB.CreateOrEditWorkspace(workspace)

				feature := db.WorkspaceFeatures{
					Uuid:          uuid.New().String(),
					WorkspaceUuid: workspace.Uuid,
					Name:          "Feature !@#$%^&*()",
					Priority:      1,
				}
				db.TestDB.CreateOrEditFeature(feature)

				return workspace.Uuid, []db.WorkspaceFeatures{feature}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Features with Unicode Characters",
			pubKeyFromAuth: "validPubKey",
			workspaceUUID:  "valid-uuid",
			setupMocks: func() (string, interface{}) {
				person := db.Person{
					Uuid:        uuid.New().String(),
					OwnerPubKey: "validPubKey",
				}
				db.TestDB.CreateOrEditPerson(person)

				workspace := db.Workspace{
					Uuid:        uuid.New().String(),
					OwnerPubKey: person.OwnerPubKey,
				}
				db.TestDB.CreateOrEditWorkspace(workspace)

				feature := db.WorkspaceFeatures{
					Uuid:          uuid.New().String(),
					WorkspaceUuid: workspace.Uuid,
					Name:          "测试Feature テスト",
					Priority:      1,
				}
				db.TestDB.CreateOrEditFeature(feature)

				return workspace.Uuid, []db.WorkspaceFeatures{feature}
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var workspaceUUID string
			var expectedBody interface{}

			if tt.setupMocks != nil {
				workspaceUUID, expectedBody = tt.setupMocks()
			}

			req := httptest.NewRequest(http.MethodGet, "/features", nil)
			req = req.WithContext(context.WithValue(req.Context(), auth.ContextKey, tt.pubKeyFromAuth))
			rctx := chi.NewRouteContext()

			if workspaceUUID != "" {
				rctx.URLParams.Add("workspace_uuid", workspaceUUID)
			} else {
				rctx.URLParams.Add("workspace_uuid", tt.workspaceUUID)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rec := httptest.NewRecorder()

			oHandler.GetFeaturesByWorkspaceUuid(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if expectedBody != nil {
				var actualBody []db.WorkspaceFeatures
				err := json.NewDecoder(rec.Body).Decode(&actualBody)
				assert.NoError(t, err)

				sort.Slice(expectedBody, func(i, j int) bool {
					return expectedBody.([]db.WorkspaceFeatures)[i].Uuid < expectedBody.([]db.WorkspaceFeatures)[j].Uuid
				})
				sort.Slice(actualBody, func(i, j int) bool {
					return actualBody[i].Uuid < actualBody[j].Uuid
				})

				assert.Len(t, actualBody, len(expectedBody.([]db.WorkspaceFeatures)))

				for i, expectedFeature := range expectedBody.([]db.WorkspaceFeatures) {
					actualFeature := actualBody[i]
					assert.Equal(t, expectedFeature.Uuid, actualFeature.Uuid)
					assert.Equal(t, expectedFeature.Name, actualFeature.Name)
					assert.Equal(t, expectedFeature.Url, actualFeature.Url)
					assert.Equal(t, expectedFeature.Priority, actualFeature.Priority)

					assert.NotNil(t, actualFeature.Created)
					assert.NotNil(t, actualFeature.Updated)
				}
			}
		})
	}
}

func TestUpdateFeatureStatus(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()

	fHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)
	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

	feature := db.WorkspaceFeatures{
		ID:            1,
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature1",
		Brief:         "brief",
		Architecture:  "architecture",
		Requirements:  "requirements",
		Url:           "https://github.com/test-feature",
		Priority:      0,
	}
	db.TestDB.CreateOrEditFeature(feature)

	tests := []struct {
		name           string
		contextKey     interface{}
		contextValue   interface{}
		uuid           string
		body           map[string]interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Request with Active Status",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           feature.Uuid,
			body:           map[string]interface{}{"status": db.ActiveFeature},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"active"}`,
		},
		{
			name:           "Valid Request with Archived Status",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           feature.Uuid,
			body:           map[string]interface{}{"status": db.ArchivedFeature},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"archived"}`,
		},
		{
			name:           "Missing UUID Parameter",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           "",
			body:           map[string]interface{}{"status": db.ActiveFeature},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "Empty Request Body",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           feature.Uuid,
			body:           nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "Invalid Feature Status",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           feature.Uuid,
			body:           map[string]interface{}{"status": 999},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "Unauthorized Access - Missing PubKey",
			contextKey:     auth.ContextKey,
			contextValue:   "",
			uuid:           feature.Uuid,
			body:           map[string]interface{}{"status": db.ActiveFeature},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name:           "Unauthorized Access - Invalid PubKey",
			contextKey:     auth.ContextKey,
			contextValue:   "invalidPubKey",
			uuid:           feature.Uuid,
			body:           map[string]interface{}{"status": db.ActiveFeature},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name:           "Invalid JSON in Request Body",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           feature.Uuid,
			body:           map[string]interface{}{"": ""},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "Valid Request with Non-Existent UUID",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           "nonExistentUUID",
			body:           map[string]interface{}{"status": db.ActiveFeature},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/features/"+tt.uuid+"/status", bytes.NewBuffer(reqBody))
			req = req.WithContext(context.WithValue(req.Context(), tt.contextKey, tt.contextValue))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("uuid", tt.uuid)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			http.HandlerFunc(fHandler.UpdateFeatureStatus).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != "" {
				var responseBody map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&responseBody)
				assert.NoError(t, err, "Failed to decode response body")

				if status, exists := responseBody["feat_status"]; exists {
					var expectedBody map[string]interface{}
					err := json.Unmarshal([]byte(tt.expectedBody), &expectedBody)
					assert.NoError(t, err, "Failed to unmarshal expected body")

					assert.Equal(t, expectedBody["status"], status, "Status field does not match")
				} else {
					assert.Fail(t, "Response body does not contain 'status' field")
				}
			}
		})
	}
}

func TestUpdateFeatureBrief(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()

	fHandler := NewFeatureHandler(db.TestDB)

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	workspace := db.Workspace{
		Uuid:        uuid.New().String(),
		Name:        "test-workspace" + uuid.New().String(),
		OwnerPubKey: person.OwnerPubKey,
		Github:      "https://github.com/test",
		Website:     "https://www.testwebsite.com",
		Description: "test-description",
	}
	db.TestDB.CreateOrEditWorkspace(workspace)
	workspace = db.TestDB.GetWorkspaceByUuid(workspace.Uuid)

	feature := db.WorkspaceFeatures{
		ID:            1,
		Uuid:          uuid.New().String(),
		WorkspaceUuid: workspace.Uuid,
		Name:          "test-feature1",
		Brief:         "brief",
		Architecture:  "architecture",
		Requirements:  "requirements",
		Url:           "https://github.com/test-feature",
		Priority:      0,
	}
	db.TestDB.CreateOrEditFeature(feature)

	tests := []struct {
		name           string
		contextKey     interface{}
		contextValue   interface{}
		uuid           string
		body           map[string]interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Request with Existing Brief",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           feature.Uuid,
			body:           map[string]interface{}{"output": map[string]interface{}{"featureBrief": "Updated feature brief", "audioLink": "http://example.com/audio", "featureUUID": feature.Uuid}},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"brief":"brief\n\n* Generated Feature Brief *\n\nUpdated feature brief"}`,
		},
		{
			name:           "Valid Request with Empty Brief",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           feature.Uuid,
			body:           map[string]interface{}{"output": map[string]interface{}{"featureBrief": "", "audioLink": "http://example.com/audio", "featureUUID": feature.Uuid}},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "Missing UUID Parameter",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           "",
			body:           map[string]interface{}{"output": map[string]interface{}{"featureBrief": "Updated feature brief", "audioLink": "http://example.com/audio", "featureUUID": ""}},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "Empty Request Body",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           feature.Uuid,
			body:           nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "Valid Request with Non-Existent UUID",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           "nonExistentUUID",
			body:           map[string]interface{}{"output": map[string]interface{}{"featureBrief": "Updated feature brief", "audioLink": "http://example.com/audio", "featureUUID": "nonExistentUUID"}},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			name:         "Unauthorized Access (No Context Value)",
			contextKey:   "",
			contextValue: "",
			uuid:         feature.Uuid,
			body: map[string]interface{}{
				"output": map[string]interface{}{
					"featureBrief": "Unauthorized request",
					"audioLink":    "http://example.com/audio",
					"featureUUID":  feature.Uuid,
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name:         "FeatureUUID with Excessively Long String",
			contextKey:   auth.ContextKey,
			contextValue: person.OwnerPubKey,
			uuid:         strings.Repeat("a", 500),
			body: map[string]interface{}{
				"output": map[string]interface{}{
					"featureBrief": "Excessively long UUID",
					"audioLink":    "http://example.com/audio",
					"featureUUID":  strings.Repeat("a", 500),
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			name:         "FeatureBrief with JSON Special Characters",
			contextKey:   auth.ContextKey,
			contextValue: person.OwnerPubKey,
			uuid:         feature.Uuid,
			body: map[string]interface{}{
				"output": map[string]interface{}{
					"featureBrief": `{"key":"value", "nested":{"key":"value"}}`,
					"audioLink":    "http://example.com/audio",
					"featureUUID":  feature.Uuid,
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"brief":"brief\n\n* Generated Feature Brief *\n\nUpdated feature brief\n\n* Generated Feature Brief *\n\n{\"key\":\"value\", \"nested\":{\"key\":\"value\"}}"}`,
		},
		{
			name:         "Simultaneous Update by Multiple Users",
			contextKey:   auth.ContextKey,
			contextValue: person.OwnerPubKey,
			uuid:         feature.Uuid,
			body: map[string]interface{}{
				"output": map[string]interface{}{
					"featureBrief": "Simultaneous updates test",
					"audioLink":    "http://example.com/audio",
					"featureUUID":  feature.Uuid,
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"brief":"brief\n\n* Generated Feature Brief *\n\nUpdated feature brief\n\n* Generated Feature Brief *\n\n{\"key\":\"value\", \"nested\":{\"key\":\"value\"}}\n\n* Generated Feature Brief *\n\nSimultaneous updates test"}`,
		},
		{
			name:           "Empty Request Body with Valid Context",
			contextKey:     auth.ContextKey,
			contextValue:   person.OwnerPubKey,
			uuid:           feature.Uuid,
			body:           nil, // Empty body
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:         "Request with Missing Required Fields in Body",
			contextKey:   auth.ContextKey,
			contextValue: person.OwnerPubKey,
			uuid:         feature.Uuid,
			body: map[string]interface{}{
				"output": map[string]interface{}{
					"featureUUID": feature.Uuid,
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/features/brief", bytes.NewBuffer(reqBody))
			req = req.WithContext(context.WithValue(req.Context(), tt.contextKey, tt.contextValue))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("uuid", tt.uuid)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			http.HandlerFunc(fHandler.UpdateFeatureBrief).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedBody != "" {
				var responseBody db.WorkspaceFeatures
				err := json.NewDecoder(rr.Body).Decode(&responseBody)
				assert.NoError(t, err, "Failed to decode response body")

				if brief := responseBody.Brief; brief != "" {
					var expectedBody map[string]interface{}
					err := json.Unmarshal([]byte(tt.expectedBody), &expectedBody)
					assert.NoError(t, err, "Failed to unmarshal expected body")

					assert.Equal(t, expectedBody["brief"], brief, "Brief field does not match")
				} else {
					assert.Fail(t, "Response body does not contain 'featureBrief' field")
				}
			}
		})
	}
}

func TestBriefSend(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	db.CleanTestData()

	person := db.Person{
		Uuid:        uuid.New().String(),
		OwnerAlias:  "test-alias",
		UniqueName:  "test-unique-name",
		OwnerPubKey: "test-pubkey",
		PriceToMeet: 0,
		Description: "test-description",
	}
	db.TestDB.CreateOrEditPerson(person)

	fHandler := NewFeatureHandler(db.TestDB)

	tests := []struct {
		name           string
		contextKey     interface{}
		contextValue   string
		body           string
		envHost        string
		envSWWFKEY     string
		mockError      error
		expectedStatus int
		expectedBody   string
		expectedPanic  string
	}{
		{
			name:           "Valid Request with All Required Fields",
			contextKey:     auth.ContextKey,
			contextValue:   "test-pubkey",
			body:           `{"audioLink":"link","featureUUID":"uuid","source":"source","examples":["example1"]}`,
			envHost:        "http://localhost",
			envSWWFKEY:     "validKey",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "{\"success\":false,\"error\":{\"message\":\"Unauthorized\"}}",
		},
		{
			name:           "Empty JSON Body",
			contextKey:     auth.ContextKey,
			contextValue:   "test-pubkey",
			body:           ``,
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   "Invalid JSON format\n",
		},
		{
			name:           "Missing Required JSON Fields",
			contextKey:     auth.ContextKey,
			contextValue:   "test-pubkey",
			body:           `{"audioLink":"link"}`,
			envHost:        "http://localhost",
			envSWWFKEY:     "",
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   "",
		},
		{
			name:           "No pubKeyFromAuth in Context",
			contextKey:     "",
			contextValue:   "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name:           "Invalid JSON Format",
			contextKey:     auth.ContextKey,
			contextValue:   "test-pubkey",
			body:           `{"audioLink":}`,
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   "Invalid JSON format\n",
		},
		{
			name:           "Environment Variable HOST Not Set",
			contextKey:     auth.ContextKey,
			contextValue:   "test-pubkey",
			body:           `{"audioLink":"link","featureUUID":"uuid","source":"source","examples":["example1"]}`,
			envHost:        "",
			envSWWFKEY:     "validKey",
			expectedStatus: http.StatusNotAcceptable,
		},
		{
			name:           "Environment Variable SWWFKEY Not Set",
			contextKey:     auth.ContextKey,
			contextValue:   "test-pubkey",
			body:           `{"audioLink":"link","featureUUID":"uuid","source":"source","examples":["example1"]}`,
			envHost:        "",
			envSWWFKEY:     "validKey",
			expectedStatus: http.StatusNotAcceptable,
		},
		{
			name:           "Stakwork the API Request Creation Failure",
			contextKey:     auth.ContextKey,
			contextValue:   "test-pubkey",
			body:           `{"audioLink":"link","featureUUID":"uuid","source":"source","examples":["example1"]}`,
			envHost:        "http://localhost",
			envSWWFKEY:     "validKey",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "{\"success\":false,\"error\":{\"message\":\"Unauthorized\"}}",
		},
		{
			name:           "Stakwork API Request Sending Failure",
			contextKey:     auth.ContextKey,
			contextValue:   "test-pubkey",
			body:           `{"audioLink":"link","featureUUID":"uuid","source":"source","examples":["example1"]}`,
			envHost:        "http://localhost",
			envSWWFKEY:     "validKey",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "{\"success\":false,\"error\":{\"message\":\"Unauthorized\"}}",
		},
		{
			name:           "Stakwork API Response Reading Failure",
			contextKey:     auth.ContextKey,
			contextValue:   "test-pubkey",
			body:           `{"audioLink":"link","featureUUID":"uuid","source":"source","examples":["example1"]}`,
			envHost:        "http://localhost",
			envSWWFKEY:     "validKey",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "{\"success\":false,\"error\":{\"message\":\"Unauthorized\"}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envHost != "" {
				os.Setenv("HOST", tt.envHost)
			}
			if tt.envSWWFKEY != "" {
				os.Setenv("SWWFKEY", tt.envSWWFKEY)
			}

			req := httptest.NewRequest(http.MethodPost, "/brief/send", bytes.NewBufferString(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), tt.contextKey, tt.contextValue))
			w := httptest.NewRecorder()

			if tt.expectedPanic != "" {
				assert.PanicsWithValue(t, tt.expectedPanic, func() {
					fHandler.BriefSend(w, req)
				})
			} else {
				fHandler.BriefSend(w, req)
				resp := w.Result()
				body, _ := io.ReadAll(resp.Body)
				assert.Equal(t, tt.expectedStatus, resp.StatusCode)
				assert.Equal(t, tt.expectedBody, string(body))
			}

			if tt.envHost != "" {
				os.Setenv("HOST", "")
			}
			if tt.envSWWFKEY != "" {
				os.Setenv("SWWFKEY", "")
			}

		})
	}
}
