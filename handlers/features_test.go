package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
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

func TestGetFeaturesByWorkspaceUuid(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	oHandler := NewWorkspaceHandler(db.TestDB)

	t.Run("should return error if a user is not authorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetFeaturesByWorkspaceUuid)

		ctx := context.WithValue(context.Background(), auth.ContextKey, "")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/forworkspace/"+workspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("created feature should be present in the returned array", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(oHandler.GetFeaturesByWorkspaceUuid)

		person := db.Person{
			Uuid:        uuid.New().String(),
			OwnerAlias:  "alias",
			UniqueName:  "unique_name",
			OwnerPubKey: "pubkey",
			PriceToMeet: 0,
			Description: "description",
		}
		db.TestDB.CreateOrEditPerson(person)
		workspace := db.Workspace{
			Uuid:        uuid.New().String(),
			Name:        "unique_workspace_name" + uuid.New().String(),
			OwnerPubKey: person.OwnerPubKey,
			Github:      "gtihub",
			Website:     "website",
			Description: "description",
		}
		db.TestDB.CreateOrEditWorkspace(workspace)
		feature := db.WorkspaceFeatures{
			Uuid:          uuid.New().String(),
			WorkspaceUuid: workspace.Uuid,
			Name:          "feature_name",
			Url:           "https://www.bountieswebsite.com",
			Priority:      0,
		}
		db.TestDB.CreateOrEditFeature(feature)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("workspace_uuid", workspace.Uuid)
		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/forworkspace/"+workspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var returnedWorkspaceFeatures []db.WorkspaceFeatures
		err = json.Unmarshal(rr.Body.Bytes(), &returnedWorkspaceFeatures)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		// Verify that the created feature is present in the returned array
		found := false
		for _, f := range returnedWorkspaceFeatures {
			if f.Uuid == feature.Uuid {
				assert.Equal(t, feature.Name, f.Name)
				assert.Equal(t, feature.Url, f.Url)
				assert.Equal(t, feature.Priority, f.Priority)
				found = true
				break
			}
		}
		assert.True(t, found, "The created feature should be present in the returned array")
	})

}

func TestGetWorkspaceFeaturesCount(t *testing.T) {
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

	ctx := context.WithValue(context.Background(), auth.ContextKey, workspace.OwnerPubKey)

	t.Run("Should test that it throws a 401 error if a user is not authorized", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/workspace/count/"+workspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetWorkspaceFeaturesCount).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that the features count returned from the API response for the workspace is equal to the number of features created for the workspace", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", workspace.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/workspace/count/"+workspace.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(oHandler.GetWorkspaceFeaturesCount).ServeHTTP(rr, req)

		var returnedWorkspaceFeatures int64
		err = json.Unmarshal(rr.Body.Bytes(), &returnedWorkspaceFeatures)
		assert.NoError(t, err)

		featureCount := db.TestDB.GetWorkspaceFeaturesCount(workspace.Uuid)

		assert.Equal(t, returnedWorkspaceFeatures, featureCount)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestGetFeatureByUuid(t *testing.T) {
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
	feature = db.TestDB.GetFeatureByUuid(feature.Uuid)

	t.Run("should return error if not authorized", func(t *testing.T) {
		featureUUID := feature.Uuid
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.GetFeatureByUuid)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", featureUUID)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, "/features/"+featureUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return feature if user is authorized", func(t *testing.T) {
		featureUUID := feature.Uuid

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(fHandler.GetFeatureByUuid)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uuid", featureUUID)
		ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, "/features/"+featureUUID, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var returnedFeature db.WorkspaceFeatures
		err = json.Unmarshal(rr.Body.Bytes(), &returnedFeature)
		assert.NoError(t, err)
		assert.Equal(t, feature.Name, returnedFeature.Name)
		assert.Equal(t, feature.Url, returnedFeature.Url)
		assert.Equal(t, feature.Priority, returnedFeature.Priority)
	})
}

func TestCreateOrEditFeaturePhase(t *testing.T) {

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
}

func TestGetFeaturePhaseByUUID(t *testing.T) {
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
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), http.MethodGet, feature.Uuid+"/phase/"+featurePhase.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeaturePhaseByUUID).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Should test that the workspace features phases returned from the API has the feature phases created", func(t *testing.T) {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("phase_uuid", featurePhase.Uuid)
		rctx.URLParams.Add("feature_uuid", feature.Uuid)
		req, err := http.NewRequestWithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx), http.MethodGet, feature.Uuid+"/phase/"+featurePhase.Uuid, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		http.HandlerFunc(fHandler.GetFeaturePhaseByUUID).ServeHTTP(rr, req)

		var returnedFeaturePhases db.FeaturePhase
		err = json.Unmarshal(rr.Body.Bytes(), &returnedFeaturePhases)
		assert.NoError(t, err)

		featurePhase.Created = returnedFeaturePhases.Created
		featurePhase.Updated = returnedFeaturePhases.Updated

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, featurePhase, returnedFeaturePhases)
	})

}

func TestDeleteFeaturePhase(t *testing.T) {

}

func TestCreateOrEditStory(t *testing.T) {

}

func TestGetStoriesByFeatureUuid(t *testing.T) {

}

func TestGetStoryByUuid(t *testing.T) {

}

func TestDeleteStory(t *testing.T) {

}

func TestGetBountiesByFeatureAndPhaseUuid(t *testing.T) {

}

func TestGetBountiesCountByFeatureAndPhaseUuid(t *testing.T) {

}
