package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
)

func TestGetTicket(t *testing.T) {
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
		Name:        "test-phase",
		Priority:    0,
	}
	db.TestDB.CreateOrEditFeaturePhase(featurePhase)

	ticket := db.Tickets{
		UUID:        uuid.New(),
		FeatureUUID: feature.Uuid,
		PhaseUUID:   featurePhase.Uuid,
		Name:        "Test Ticket",
		Sequence:    1,
		Description: "Test Description",
		Status:      db.DraftTicket,
	}
	createdTicket, _ := db.TestDB.UpdateTicket(ticket)

	t.Run("should return 400 if UUID is empty", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GetTicket)

		req, err := http.NewRequest(http.MethodGet, "/tickets/", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 404 if ticket doesn't exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GetTicket)

		nonExistentUUID := uuid.New()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", nonExistentUUID.String())
		req, err := http.NewRequest(http.MethodGet, "/tickets/"+nonExistentUUID.String(), nil)
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("should return ticket if exists", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.GetTicket)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", createdTicket.UUID.String())
		req, err := http.NewRequest(http.MethodGet, "/tickets/"+createdTicket.UUID.String(), nil)
		if err != nil {
			t.Fatal(err)
		}
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var returnedTicket db.Tickets
		err = json.Unmarshal(rr.Body.Bytes(), &returnedTicket)
		assert.NoError(t, err)
		assert.Equal(t, createdTicket.Name, returnedTicket.Name)
		assert.Equal(t, createdTicket.Description, returnedTicket.Description)
		assert.Equal(t, createdTicket.Status, returnedTicket.Status)
		assert.Equal(t, createdTicket.FeatureUUID, returnedTicket.FeatureUUID)
		assert.Equal(t, createdTicket.PhaseUUID, returnedTicket.PhaseUUID)
	})
}

func TestUpdateTicket(t *testing.T) {
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
		Name:        "test-phase",
		Priority:    0,
	}
	db.TestDB.CreateOrEditFeaturePhase(featurePhase)

	ticket := db.Tickets{
		UUID:        uuid.New(),
		FeatureUUID: feature.Uuid,
		PhaseUUID:   featurePhase.Uuid,
		Name:        "Test Ticket",
		Sequence:    1,
		Description: "Test Description",
		Status:      db.DraftTicket,
	}
	createdTicket, _ := db.TestDB.UpdateTicket(ticket)

	t.Run("should return 401 if no auth token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.UpdateTicket)

		req, err := http.NewRequest(http.MethodPost, "/tickets/", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return 400 if UUID is empty", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.UpdateTicket)

		req, err := http.NewRequest(http.MethodPost, "/tickets/", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 400 if UUID is invalid", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.UpdateTicket)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", "invalid-uuid")

		req, err := http.NewRequest(http.MethodPost, "/tickets/invalid-uuid", bytes.NewReader([]byte("{}")))
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 400 if body is not valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.UpdateTicket)

		invalidJson := []byte(`{"key": "value"`)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", createdTicket.UUID.String())
		req, err := http.NewRequest(http.MethodPost, "/tickets/"+createdTicket.UUID.String(), bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should update ticket with only UUID and optional fields", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.UpdateTicket)

		updateTicket := createdTicket
		updateTicket.Description = "Updated Description"
		updateTicket.Status = db.ReadyTicket

		updateRequest := UpdateTicketRequest{
			Metadata: struct {
				Source string `json:"source"`
				ID     string `json:"id"`
			}{
				Source: "test-source",
				ID:     "test-id",
			},
			Ticket: &updateTicket,
		}

		requestBody, _ := json.Marshal(updateRequest)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", updateTicket.UUID.String())
		req, err := http.NewRequest(http.MethodPost, "/tickets/"+updateTicket.UUID.String(), bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var returnedTicket db.Tickets
		err = json.Unmarshal(rr.Body.Bytes(), &returnedTicket)
		assert.NoError(t, err)

		// Verify that only the provided fields were updated
		assert.Equal(t, updateTicket.Description, returnedTicket.Description)
		assert.Equal(t, updateTicket.Status, returnedTicket.Status)
		// Original fields should remain unchanged
		assert.Equal(t, createdTicket.FeatureUUID, returnedTicket.FeatureUUID)
		assert.Equal(t, createdTicket.PhaseUUID, returnedTicket.PhaseUUID)
		assert.Equal(t, createdTicket.Name, returnedTicket.Name)
	})

	t.Run("should update ticket successfully", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.UpdateTicket)

		updatedTicket := createdTicket
		updatedTicket.Name = "Updated Test Ticket"
		updatedTicket.Description = "Updated Description"
		updatedTicket.Status = db.CompletedTicket

		updateRequest := UpdateTicketRequest{
			Metadata: struct {
				Source string `json:"source"`
				ID     string `json:"id"`
			}{
				Source: "test-source",
				ID:     "test-id",
			},
			Ticket: &updatedTicket,
		}

		requestBody, _ := json.Marshal(updateRequest)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", createdTicket.UUID.String())

		req, err := http.NewRequest(http.MethodPost, "/tickets/"+createdTicket.UUID.String(), bytes.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var returnedTicket db.Tickets
		err = json.Unmarshal(rr.Body.Bytes(), &returnedTicket)
		assert.NoError(t, err)
		assert.Equal(t, updatedTicket.Name, returnedTicket.Name)
		assert.Equal(t, updatedTicket.Description, returnedTicket.Description)
		assert.Equal(t, updatedTicket.Status, returnedTicket.Status)
		assert.Equal(t, updatedTicket.FeatureUUID, returnedTicket.FeatureUUID)
		assert.Equal(t, updatedTicket.PhaseUUID, returnedTicket.PhaseUUID)
	})

}

func TestDeleteTicket(t *testing.T) {
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
		Name:        "test-phase",
		Priority:    0,
	}
	db.TestDB.CreateOrEditFeaturePhase(featurePhase)

	ticket := db.Tickets{
		UUID:        uuid.New(),
		FeatureUUID: feature.Uuid,
		PhaseUUID:   featurePhase.Uuid,
		Name:        "Test Ticket",
		Sequence:    1,
		Description: "Test Description",
		Status:      db.DraftTicket,
	}
	createdTicket, _ := db.TestDB.UpdateTicket(ticket)

	t.Run("should return 401 if no auth token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.DeleteTicket)

		req, err := http.NewRequest(http.MethodDelete, "/tickets/", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return 400 if UUID is empty", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.DeleteTicket)

		req, err := http.NewRequest(http.MethodDelete, "/tickets/", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 404 if ticket doesn't exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.DeleteTicket)

		nonExistentUUID := uuid.New()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", nonExistentUUID.String())

		req, err := http.NewRequest(http.MethodDelete, "/tickets/"+nonExistentUUID.String(), nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("should delete ticket successfully", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.DeleteTicket)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", createdTicket.UUID.String())

		req, err := http.NewRequest(http.MethodDelete, "/tickets/"+createdTicket.UUID.String(), nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Verify ticket was deleted
		_, err = db.TestDB.GetTicket(createdTicket.UUID.String())
		assert.Error(t, err)
		assert.Equal(t, "ticket not found", err.Error())
	})
}

func TestTicketToBounty(t *testing.T) {
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
		Name:        "test-phase",
		Priority:    0,
	}
	db.TestDB.CreateOrEditFeaturePhase(featurePhase)

	ticket := db.Tickets{
		UUID:        uuid.New(),
		FeatureUUID: feature.Uuid,
		PhaseUUID:   featurePhase.Uuid,
		Name:        "Test Ticket",
		Sequence:    1,
		Description: "Test Description",
		Status:      db.DraftTicket,
	}
	createdTicket, _ := db.TestDB.UpdateTicket(ticket)

	t.Run("should return 401 if no auth token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.TicketToBounty)

		req, err := http.NewRequest(http.MethodPost, "/tickets/bounty", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return 400 if ticket UUID is empty", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.TicketToBounty)

		req, err := http.NewRequest(http.MethodPost, "/tickets/bounty", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
		req = req.WithContext(ctx)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 404 if ticket doesn't exist", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.TicketToBounty)

		nonExistentUUID := uuid.New()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("ticket_uuid", nonExistentUUID.String())

		req, err := http.NewRequest(http.MethodPost, "/tickets/"+nonExistentUUID.String()+"/bounty", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("should create bounty successfully", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tHandler.TicketToBounty)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("ticket_uuid", createdTicket.UUID.String())

		req, err := http.NewRequest(http.MethodPost, "/tickets/"+createdTicket.UUID.String()+"/bounty", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(req.Context(), auth.ContextKey, person.OwnerPubKey)
		req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response CreateBountyResponse
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotZero(t, response.BountyID)
		assert.Equal(t, "Bounty created successfully", response.Message)

		// Verify bounty was created with correct data
		bounty := db.TestDB.GetBounty(response.BountyID)
		assert.Equal(t, createdTicket.Name, bounty.Title)
		assert.Equal(t, createdTicket.Description, bounty.Description)
		assert.Equal(t, createdTicket.PhaseUUID, bounty.PhaseUuid)
		assert.True(t, bounty.Show)
		assert.Equal(t, uint(21), bounty.Price)
		assert.Equal(t, "Other", bounty.Type)
	})
}
