package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/utils"
	"github.com/stretchr/testify/assert"
)

func TestSendEdgeListToJarvis(t *testing.T) {
	originalURL := config.JarvisUrl
	originalToken := config.JarvisToken
	defer func() {
		config.JarvisUrl = originalURL
		config.JarvisToken = originalToken
	}()

	tests := []struct {
		name          string
		edgeList      utils.EdgeList
		setupMock     func() *httptest.Server
		setupConfig   func()
		expectedError bool
	}{
		{
			name: "Basic Functionality",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{
					{
						Edge: utils.EdgeInfo{
							EdgeType: "test",
							Weight:   1.0,
						},
						Source: utils.Node{
							NodeType: "test",
							NodeData: map[string]interface{}{"test": "data"},
						},
						Targets: []utils.Node{
							{
								NodeType: "test",
								NodeData: map[string]interface{}{"test": "data"},
							},
						},
					},
				},
			},
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
					assert.Equal(t, "test-token", r.Header.Get("x-api-token"))
					w.WriteHeader(http.StatusOK)
				}))
			},
			setupConfig: func() {
				config.JarvisToken = "test-token"
			},
			expectedError: false,
		},
		{
			name: "Empty EdgeList",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{},
			},
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			setupConfig: func() {
				config.JarvisToken = "test-token"
			},
			expectedError: false,
		},
		{
			name: "Missing JarvisUrl",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{
					{
						Edge: utils.EdgeInfo{
							EdgeType: "test",
							Weight:   1.0,
						},
					},
				},
			},
			setupMock: nil,
			setupConfig: func() {
				config.JarvisUrl = ""
				config.JarvisToken = "test-token"
			},
			expectedError: false,
		},
		{
			name: "Missing JarvisToken",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{
					{
						Edge: utils.EdgeInfo{
							EdgeType: "test",
							Weight:   1.0,
						},
					},
				},
			},
			setupMock: nil,
			setupConfig: func() {
				config.JarvisUrl = "http://test.com"
				config.JarvisToken = ""
			},
			expectedError: false,
		},
		{
			name: "Invalid JSON Marshalling",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{
					{
						Edge: utils.EdgeInfo{
							EdgeType: "test",
							Weight:   1.0,
						},
						Source: utils.Node{
							NodeType: "test",
							NodeData: make(chan int),
						},
					},
				},
			},
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			setupConfig: func() {
				config.JarvisToken = "test-token"
			},
			expectedError: false,
		},
		{
			name: "HTTP Request Creation Failure",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{
					{
						Edge: utils.EdgeInfo{
							EdgeType: "test",
							Weight:   1.0,
						},
					},
				},
			},
			setupMock: nil,
			setupConfig: func() {
				config.JarvisUrl = "http://[::1]:namedport"
				config.JarvisToken = "test-token"
			},
			expectedError: false,
		},
		{
			name: "Network Failure",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{
					{
						Edge: utils.EdgeInfo{
							EdgeType: "test",
							Weight:   1.0,
						},
					},
				},
			},
			setupMock: func() *httptest.Server {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic("Force connection close")
				}))
				server.Close()
				return server
			},
			setupConfig: func() {
				config.JarvisToken = "test-token"
			},
			expectedError: true,
		},
		{
			name: "Non-2xx HTTP Response",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{
					{
						Edge: utils.EdgeInfo{
							EdgeType: "test",
							Weight:   1.0,
						},
					},
				},
			},
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{
						"error": "Internal Server Error",
					})
				}))
			},
			setupConfig: func() {
				config.JarvisToken = "test-token"
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupConfig != nil {
				tt.setupConfig()
			}

			if tt.setupMock != nil {
				server := tt.setupMock()
				defer server.Close()
				config.JarvisUrl = server.URL
			}

			err := sendEdgeListToJarvis(tt.edgeList)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
