package handlers

import (
	"bytes"
	"context"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateOrEditBounty(t *testing.T) {

	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")

	t.Run("should return error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateOrEditBounty)

		invalidJson := []byte(`{"key": "value"`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code, "invalid status received")
	})

	t.Run("missing required field, bounty type", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateOrEditBounty)

		invalidBody := []byte(`{"type": ""}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("missing required field, bounty title", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateOrEditBounty)

		invalidBody := []byte(`{"type": "bounty_type", "title": ""}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("missing required field, bounty description", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateOrEditBounty)

		invalidBody := []byte(`{"type": "bounty_type", "title": "first bounty", "description": ""}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(invalidBody))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("return error if trying to update other user bounty", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateOrEditBounty)

		body := []byte(`{"id": 1, "type": "bounty_type", "title": "first bounty", "description": "my first bounty", "tribe": "random-value", "assignee": "john-doe", "owner_id": "second-user"}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		bounties, _ := db.DB.GetBountyById("1")
		assert.Equal(t, 0, len(bounties))
	})

	t.Run("return error if user does not have required roles", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateOrEditBounty)

		mockOrg := db.Organization{
			ID:          1,
			Uuid:        "org-1",
			Name:        "custom org",
			OwnerPubKey: "org-key",
		}
		_, _ = db.DB.CreateOrEditOrganization(mockOrg)
		_ = db.DB.CreateUserRoles(nil, mockOrg.Uuid, "test-key")
		body := []byte(`{"id": 1, "type": "bounty_type", "title": "first bounty", "description": "my first bounty", "tribe": "random-value", "assignee": "john-doe", "owner_id": "second-user", "org_uuid": "org-1"}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		bounties, _ := db.DB.GetBountyById("1")
		assert.Equal(t, 0, len(bounties))
	})

	t.Run("should allow to add or edit bounty if user has role", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateOrEditBounty)

		mockOrg := db.Organization{
			ID:          1,
			Uuid:        "org-1",
			Name:        "custom org",
			OwnerPubKey: "org-key",
		}
		roles := []db.UserRoles{
			{Role: db.AddBounty, OwnerPubKey: "test-key", OrgUuid: "org-1"},
			{Role: db.UpdateBounty, OwnerPubKey: "test-key", OrgUuid: "org-1"},
			{Role: db.DeleteBounty, OwnerPubKey: "test-key", OrgUuid: "org-1"},
			{Role: db.PayBounty, OwnerPubKey: "test-key", OrgUuid: "org-1"},
		}
		_, _ = db.DB.CreateOrEditOrganization(mockOrg)
		_ = db.DB.CreateUserRoles(roles, mockOrg.Uuid, "test-key")
		body := []byte(`{"id": 1, "type": "bounty_type", "title": "first bounty", "description": "my first bounty", "tribe": "random-value", "assignee": "john-doe", "owner_id": "second-user", "org_uuid": "org-1"}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		bounties, _ := db.DB.GetBountyById("1")
		assert.Equal(t, 1, len(bounties))
		assert.Equal(t, "bounty_type", bounties[0].Type)
		assert.Equal(t, "first bounty", bounties[0].Title)
	})

	t.Run("add bounty if not present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateOrEditBounty)

		body := []byte(`{"id": 0, "type": "bounty_type", "title": "first bounty", "description": "my first bounty", "tribe": "random-value", "assignee": "john-doe", "owner_id": "test-key"}`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		bounties, _ := db.DB.GetBountyById("0")
		assert.Equal(t, 1, len(bounties))
		assert.Equal(t, "bounty_type", bounties[0].Type)
		assert.Equal(t, "first bounty", bounties[0].Title)
	})
}
