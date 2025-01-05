package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/db"
	dbmocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFeatureFlag(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		setupMock      func(*dbmocks.Database)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Feature enabled - should allow request",
			path: "/api/test",
			setupMock: func(mockDb *dbmocks.Database) {
				flagUUID := uuid.New()
				mockDb.On("GetAllEndpoints").Return([]db.Endpoint{
					{
						UUID:            uuid.New(),
						Path:            "/api/test",
						FeatureFlagUUID: flagUUID,
					},
				}, nil)
				mockDb.On("GetFeatureFlagByUUID", flagUUID).Return(db.FeatureFlag{
					UUID:    flagUUID,
					Enabled: true,
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Feature disabled - should block request",
			path: "/api/test",
			setupMock: func(mockDb *dbmocks.Database) {
				flagUUID := uuid.New()
				mockDb.On("GetAllEndpoints").Return([]db.Endpoint{
					{
						UUID:            uuid.New(),
						Path:            "/api/test",
						FeatureFlagUUID: flagUUID,
					},
				}, nil)
				mockDb.On("GetFeatureFlagByUUID", flagUUID).Return(db.FeatureFlag{
					UUID:    flagUUID,
					Enabled: false,
				}, nil)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "This feature is currently unavailable.",
			},
		},
		{
			name: "Path not feature flagged - should allow request",
			path: "/api/unrestricted",
			setupMock: func(mockDb *dbmocks.Database) {
				mockDb.On("GetAllEndpoints").Return([]db.Endpoint{
					{
						UUID:            uuid.New(),
						Path:            "/api/test",
						FeatureFlagUUID: uuid.New(),
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Path with parameters - should match correctly",
			path: "/api/users/123",
			setupMock: func(mockDb *dbmocks.Database) {
				flagUUID := uuid.New()
				mockDb.On("GetAllEndpoints").Return([]db.Endpoint{
					{
						UUID:            uuid.New(),
						Path:            "/api/users/:id",
						FeatureFlagUUID: flagUUID,
					},
				}, nil)
				mockDb.On("GetFeatureFlagByUUID", flagUUID).Return(db.FeatureFlag{
					UUID:    flagUUID,
					Enabled: true,
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := dbmocks.NewDatabase(t)

			tt.setupMock(mockDb)

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			middleware := FeatureFlag(mockDb)
			middleware(nextHandler).ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestMatchPath(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		requestPath string
		expected    bool
	}{
		{
			name:        "Exact match",
			pattern:     "/api/test",
			requestPath: "/api/test",
			expected:    true,
		},
		{
			name:        "Parameter match",
			pattern:     "/api/users/:id",
			requestPath: "/api/users/123",
			expected:    true,
		},
		{
			name:        "No match - different path",
			pattern:     "/api/test",
			requestPath: "/api/other",
			expected:    false,
		},
		{
			name:        "Multiple parameters",
			pattern:     "/api/:resource/:id",
			requestPath: "/api/users/123",
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchPath(tt.pattern, tt.requestPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}
