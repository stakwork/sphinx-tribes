package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/db"
	dbMocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newFeatureFlagEndpointRequest(method string, body string, params map[string]string) (*http.Request, *httptest.ResponseRecorder) {
	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}

	req := httptest.NewRequest(method, "/feature-flags", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	return req, httptest.NewRecorder()
}

func TestAddFeatureFlagEndpointAcceptsArrayBody(t *testing.T) {
	flagUUID := uuid.New()
	mockDB := dbMocks.NewDatabase(t)
	handler := NewFeatureFlagHandler(mockDB)

	mockDB.On("GetFeatureFlagByUUID", flagUUID).Return(db.FeatureFlag{UUID: flagUUID}, nil).Once()
	mockDB.On("GetEndpointsByFeatureFlag", flagUUID).Return([]db.Endpoint{}, nil).Once()
	mockDB.On("AddEndpoint", mock.MatchedBy(func(endpoint *db.Endpoint) bool {
		return endpoint.UUID != uuid.Nil &&
			endpoint.Path == "/api/select" &&
			endpoint.FeatureFlagUUID == flagUUID
	})).Return(db.Endpoint{
		UUID:            uuid.New(),
		Path:            "/api/select",
		FeatureFlagUUID: flagUUID,
	}, nil).Once()
	mockDB.On("GetFeatureFlagByUUID", flagUUID).Return(db.FeatureFlag{
		UUID: flagUUID,
		Endpoints: []db.Endpoint{{
			UUID:            uuid.New(),
			Path:            "/api/select",
			FeatureFlagUUID: flagUUID,
		}},
	}, nil).Once()

	req, rr := newFeatureFlagEndpointRequest(http.MethodPost, `["/api/select"]`, map[string]string{
		"feature_flag_id": flagUUID.String(),
	})

	handler.AddFeatureFlagEndpoint(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "Endpoints processed successfully")
}

func TestAddFeatureFlagEndpointRejectsInvalidPath(t *testing.T) {
	flagUUID := uuid.New()
	mockDB := dbMocks.NewDatabase(t)
	handler := NewFeatureFlagHandler(mockDB)

	req, rr := newFeatureFlagEndpointRequest(http.MethodPost, `["api/select"]`, map[string]string{
		"feature_flag_id": flagUUID.String(),
	})

	handler.AddFeatureFlagEndpoint(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid input")
}

func TestAddFeatureFlagEndpointRejectsDuplicatePath(t *testing.T) {
	flagUUID := uuid.New()
	mockDB := dbMocks.NewDatabase(t)
	handler := NewFeatureFlagHandler(mockDB)

	mockDB.On("GetFeatureFlagByUUID", flagUUID).Return(db.FeatureFlag{UUID: flagUUID}, nil).Once()
	mockDB.On("GetEndpointsByFeatureFlag", flagUUID).Return([]db.Endpoint{{
		UUID:            uuid.New(),
		Path:            "/api/select",
		FeatureFlagUUID: flagUUID,
	}}, nil).Once()

	req, rr := newFeatureFlagEndpointRequest(http.MethodPost, `{"endpoints":["/api/select"]}`, map[string]string{
		"feature_flag_id": flagUUID.String(),
	})

	handler.AddFeatureFlagEndpoint(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Contains(t, rr.Body.String(), "Endpoint path already exists")
}

func TestUpdateFeatureFlagEndpointRejectsMismatchedFeatureFlag(t *testing.T) {
	flagUUID := uuid.New()
	endpointUUID := uuid.New()
	mockDB := dbMocks.NewDatabase(t)
	handler := NewFeatureFlagHandler(mockDB)

	mockDB.On("GetFeatureFlagByUUID", flagUUID).Return(db.FeatureFlag{UUID: flagUUID}, nil).Once()
	mockDB.On("GetEndpointByUUID", endpointUUID).Return(db.Endpoint{
		UUID:            endpointUUID,
		Path:            "/api/old",
		FeatureFlagUUID: uuid.New(),
	}, nil).Once()

	req, rr := newFeatureFlagEndpointRequest(http.MethodPut, `{"new_endpoint_path":"/api/new"}`, map[string]string{
		"feature_flag_id": flagUUID.String(),
		"endpoint_id":     endpointUUID.String(),
	})

	handler.UpdateFeatureFlagEndpoint(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Endpoint does not belong")
}

func TestUpdateFeatureFlagEndpointRejectsDuplicatePath(t *testing.T) {
	flagUUID := uuid.New()
	endpointUUID := uuid.New()
	mockDB := dbMocks.NewDatabase(t)
	handler := NewFeatureFlagHandler(mockDB)

	mockDB.On("GetFeatureFlagByUUID", flagUUID).Return(db.FeatureFlag{UUID: flagUUID}, nil).Once()
	mockDB.On("GetEndpointByUUID", endpointUUID).Return(db.Endpoint{
		UUID:            endpointUUID,
		Path:            "/api/old",
		FeatureFlagUUID: flagUUID,
	}, nil).Once()
	mockDB.On("GetEndpointsByFeatureFlag", flagUUID).Return([]db.Endpoint{
		{
			UUID:            endpointUUID,
			Path:            "/api/old",
			FeatureFlagUUID: flagUUID,
		},
		{
			UUID:            uuid.New(),
			Path:            "/api/new",
			FeatureFlagUUID: flagUUID,
		},
	}, nil).Once()

	req, rr := newFeatureFlagEndpointRequest(http.MethodPut, `{"new_endpoint_path":"/api/new"}`, map[string]string{
		"feature_flag_id": flagUUID.String(),
		"endpoint_id":     endpointUUID.String(),
	})

	handler.UpdateFeatureFlagEndpoint(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Contains(t, rr.Body.String(), "Endpoint path already exists")
}
