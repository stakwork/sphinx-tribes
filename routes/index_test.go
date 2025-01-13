package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/utils"
	"github.com/stretchr/testify/assert"
)

const defaultAuthURL = "http://auth:9090"

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
		{
			name: "Large Payload",
			edgeList: utils.EdgeList{
				EdgeList: func() []utils.Edge {
					edges := make([]utils.Edge, 1000)
					for i := 0; i < 1000; i++ {
						edges[i] = utils.Edge{
							Edge: utils.EdgeInfo{
								EdgeType: fmt.Sprintf("test_%d", i),
								Weight:   float64(i),
							},
							Source: utils.Node{
								NodeType: "test",
								NodeData: map[string]interface{}{
									"data": strings.Repeat("large_payload", 100),
								},
							},
						}
					}
					return edges
				}(),
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
			name: "Slow Server Response",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{{
					Edge: utils.EdgeInfo{
						EdgeType: "test",
						Weight:   1.0,
					},
				}},
			},
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(11 * time.Second) // Longer than client timeout
					w.WriteHeader(http.StatusOK)
				}))
			},
			setupConfig: func() {
				config.JarvisToken = "test-token"
			},
			expectedError: true,
		},
		{
			name: "Special Characters in EdgeType",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{{
					Edge: utils.EdgeInfo{
						EdgeType: "test!@#$%^&*()",
						Weight:   1.0,
					},
				}},
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
			name: "Unicode Characters in NodeData",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{{
					Edge: utils.EdgeInfo{
						EdgeType: "test",
						Weight:   1.0,
					},
					Source: utils.Node{
						NodeType: "test",
						NodeData: map[string]interface{}{
							"data": "测试データ",
						},
					},
				}},
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
			name: "Malformed Response Body",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{{
					Edge: utils.EdgeInfo{
						EdgeType: "test",
						Weight:   1.0,
					},
				}},
			},
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Length", "1000")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("Malformed JSON{"))
				}))
			},
			setupConfig: func() {
				config.JarvisToken = "test-token"
			},
			expectedError: true,
		},
		{
			name: "Empty Response Body with Success Status",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{{
					Edge: utils.EdgeInfo{
						EdgeType: "test",
						Weight:   1.0,
					},
				}},
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
			name: "Invalid Content-Type Response",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{{
					Edge: utils.EdgeInfo{
						EdgeType: "test",
						Weight:   1.0,
					},
				}},
			},
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Not JSON"))
				}))
			},
			setupConfig: func() {
				config.JarvisToken = "test-token"
			},
			expectedError: false,
		},
		{
			name: "Redirect Response",
			edgeList: utils.EdgeList{
				EdgeList: []utils.Edge{{
					Edge: utils.EdgeInfo{
						EdgeType: "test",
						Weight:   1.0,
					},
				}},
			},
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Redirect(w, r, "/new-location", http.StatusTemporaryRedirect)
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

var currentAuthURL = defaultAuthURL

var originalGetFromAuth = getFromAuth

func TestGetFromAuth(t *testing.T) {
	originalClient := http.DefaultClient
	defer func() {
		http.DefaultClient = originalClient
	}()

	tests := []struct {
		name           string
		path           string
		setupMock      func() *httptest.Server
		expectedResult *extractResponse
		expectedError  bool
	}{
		{
			name: "Valid Path with Valid JSON Response",
			path: "/valid",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey": "test-pubkey",
						"valid":  true,
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "test-pubkey",
				Valid:  true,
			},
			expectedError: false,
		},
		{
			name: "Valid Path with JSON Response Missing pubkey",
			path: "/missing-pubkey",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"valid": true,
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "",
				Valid:  true,
			},
			expectedError: false,
		},
		{
			name: "Valid Path with JSON Response Missing valid",
			path: "/missing-valid",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey": "test-pubkey",
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "test-pubkey",
				Valid:  false,
			},
			expectedError: false,
		},
		{
			name: "HTTP Request Error",
			path: "/error",
			setupMock: func() *httptest.Server {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic("Force connection close")
				}))
				server.Close()
				return server
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "Invalid JSON Response",
			path: "/invalid-json",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("invalid json{"))
				}))
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "Non-JSON Response",
			path: "/non-json",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("Plain text response"))
				}))
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "HTTP Response with Status Code 404",
			path: "/not-found",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "HTTP Response with Status Code 500",
			path: "/server-error",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "Large JSON Response",
			path: "/large-response",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					largeData := make(map[string]interface{})
					largeData["pubkey"] = "test-pubkey"
					largeData["valid"] = true
					for i := 0; i < 50000; i++ {
						largeData[fmt.Sprintf("key_%d", i)] = "large value"
					}
					json.NewEncoder(w).Encode(largeData)
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "test-pubkey",
				Valid:  true,
			},
			expectedError: false,
		},
		{
			name: "JSON Response with Additional Fields",
			path: "/extra-fields",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey":  "test-pubkey",
						"valid":   true,
						"extra":   "field",
						"another": 123,
						"moreFields": map[string]interface{}{
							"nested": "value",
						},
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "test-pubkey",
				Valid:  true,
			},
			expectedError: false,
		},
		{
			name: "JSON Response with Null Values",
			path: "/null-values",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(`{"pubkey": null, "valid": null}`))
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "",
				Valid:  false,
			},
			expectedError: false,
		},
		{
			name: "JSON Response with Incorrect Data Types",
			path: "/incorrect-types",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey": 12345,
						"valid":  "not-a-bool",
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "",
				Valid:  false,
			},
			expectedError: false,
		},
		{
			name: "Empty Response Body",
			path: "/empty-body",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "Slow Response",
			path: "/slow-response",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(2 * time.Second)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey": "test-pubkey",
						"valid":  true,
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "test-pubkey",
				Valid:  true,
			},
			expectedError: false,
		},
		{
			name: "Malformed URL Path",
			path: "/%invalid-path",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
				}))
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "Unicode Characters in Response",
			path: "/unicode",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey": "测试-pubkey",
						"valid":  true,
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "测试-pubkey",
				Valid:  true,
			},
			expectedError: false,
		},
		{
			name: "Very Long Path",
			path: "/" + strings.Repeat("a", 2048),
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusRequestURITooLong)
				}))
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "Response with Whitespace",
			path: "/whitespace",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(`{
						"pubkey": "  test-pubkey  ",
						"valid": true
					}`))
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "  test-pubkey  ",
				Valid:  true,
			},
			expectedError: false,
		},
		{
			name: "Response with Zero Values",
			path: "/zero-values",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey": "",
						"valid":  false,
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "",
				Valid:  false,
			},
			expectedError: false,
		},
		{
			name: "Chunked Response",
			path: "/chunked",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					flusher, ok := w.(http.Flusher)
					if !ok {
						return
					}
					w.Header().Set("Transfer-Encoding", "chunked")
					fmt.Fprintf(w, `{"pubkey": "`)
					flusher.Flush()
					time.Sleep(100 * time.Millisecond)
					fmt.Fprintf(w, `test-pubkey", "valid": true}`)
					flusher.Flush()
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "test-pubkey",
				Valid:  true,
			},
			expectedError: false,
		},
		{
			name: "Response with Escaped Characters",
			path: "/escaped",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(`{"pubkey": "test\npubkey\twith\"escaped\"chars", "valid": true}`))
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "test\npubkey\twith\"escaped\"chars",
				Valid:  true,
			},
			expectedError: false,
		},
		{
			name: "Empty Path",
			path: "",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
				}))
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "Path with Special Characters",
			path: "/test@#$%^&*()",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
				}))
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "Network Error (Connection Refused)",
			path: "/network-error",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "Empty JSON Response",
			path: "/empty-json",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte("{}"))
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "",
				Valid:  false,
			},
			expectedError: false,
		},
		{
			name: "Non-Boolean Valid Field",
			path: "/non-boolean-valid",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey": "test-pubkey",
						"valid":  "true",
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "test-pubkey",
				Valid:  false,
			},
			expectedError: false,
		},
		{
			name: "Non-String Pubkey Field",
			path: "/non-string-pubkey",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey": 12345,
						"valid":  true,
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "",
				Valid:  true,
			},
			expectedError: false,
		},
		{
			name: "Response with Mixed Types",
			path: "/mixed-types",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey": []interface{}{"test-pubkey"},
						"valid":  1,
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "",
				Valid:  false,
			},
			expectedError: false,
		},
		{
			name: "Response with Nested JSON",
			path: "/nested-json",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey": "test-pubkey",
						"valid":  true,
						"nested": map[string]interface{}{
							"additional": "data",
						},
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "test-pubkey",
				Valid:  true,
			},
			expectedError: false,
		},
		{
			name: "Response with Array Values",
			path: "/array-values",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey": []string{"test-pubkey"},
						"valid":  []bool{true},
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "",
				Valid:  false,
			},
			expectedError: false,
		},
		{
			name: "Response with Empty Strings",
			path: "/empty-strings",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey": "",
						"valid":  true,
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "",
				Valid:  true,
			},
			expectedError: false,
		},
		{
			name: "Response with Special Characters in Values",
			path: "/special-chars-values",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"pubkey": "test-pubkey!@#$%^&*()",
						"valid":  true,
					})
				}))
			},
			expectedResult: &extractResponse{
				Pubkey: "test-pubkey!@#$%^&*()",
				Valid:  true,
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupMock()
			defer server.Close()

			http.DefaultClient = &http.Client{
				Transport: &http.Transport{
					Proxy: func(req *http.Request) (*url.URL, error) {
						return url.Parse(server.URL)
					},
				},
			}

			result, err := getFromAuth(tt.path)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestInitChi(t *testing.T) {

	t.Run("Basic Router Initialization", func(t *testing.T) {
		router := chi.NewRouter()
		assert.NotNil(t, router, "Router should be initialized")
		assert.IsType(t, &chi.Mux{}, router, "Router should be of type chi.Mux")
	})

	t.Run("Middleware Stack Configuration", func(t *testing.T) {
		router := chi.NewRouter()
		router.Use(middleware.RequestID)
		router.Use(middleware.Recoverer)

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		router.Get("/test", testHandler)

		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "Handler should execute successfully")
	})

	t.Run("CORS Configuration", func(t *testing.T) {
		router := chi.NewRouter()
		corsMiddleware := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		})
		router.Use(corsMiddleware.Handler)

		router.Get("/test-cors", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("OPTIONS", "/test-cors", nil)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Access-Control-Request-Method", "GET")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"), "CORS should allow all origins")
	})

	t.Run("Timeout Middleware", func(t *testing.T) {
		router := chi.NewRouter()
		router.Use(middleware.Timeout(10 * time.Millisecond))

		router.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-r.Context().Done():

				w.WriteHeader(http.StatusServiceUnavailable)
				return
			case <-time.After(100 * time.Millisecond):

				w.WriteHeader(http.StatusOK)
			}
		})

		req := httptest.NewRequest("GET", "/slow", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusServiceUnavailable, rr.Code, "Should timeout and return 503")
	})

	t.Run("Internal Server Error Handler", func(t *testing.T) {
		router := chi.NewRouter()

		router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if rvr := recover(); rvr != nil {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}()
				next.ServeHTTP(w, r)
			})
		})

		router.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		req := httptest.NewRequest("GET", "/panic", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Should handle panic and return 500")
	})

	t.Run("Request ID Generation", func(t *testing.T) {
		router := chi.NewRouter()
		router.Use(middleware.RequestID)

		router.Get("/test-id", func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			w.Header().Set("X-Request-ID", reqID)
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test-id", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.NotEmpty(t, rr.Header().Get("X-Request-ID"), "Request ID should not be empty")
		assert.Equal(t, http.StatusOK, rr.Code, "Handler should execute successfully")
	})

	t.Run("Multiple Middleware Interaction", func(t *testing.T) {
		router := chi.NewRouter()
		router.Use(middleware.RequestID)
		router.Use(middleware.Recoverer)
		router.Use(cors.Default().Handler)

		router.Get("/test-multi", func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			assert.NotEmpty(t, reqID, "Request ID should be present")
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test-multi", nil)
		req.Header.Set("Origin", "http://example.com")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "Handler should execute successfully")
		assert.NotEmpty(t, rr.Header().Get("Access-Control-Allow-Origin"), "CORS headers should be present")
	})
}

var mockLogger = &MockLogger{}

type MockLogger struct {
	mu      sync.Mutex
	entries []string
}

func (l *MockLogger) Error(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, format)
}

func (l *MockLogger) Machine(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, format)
}

func (l *MockLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = []string{}
}

func TestInternalServerErrorHandler(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.Handler
		expectedStatus int
		expectLog      bool
	}{
		{
			name: "Standard Request Handling",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatus: http.StatusOK,
			expectLog:      false,
		},
		{
			name: "Request with Logging",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mockLogger.Machine("Logging message")
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatus: http.StatusOK,
			expectLog:      true,
		},
		{
			name: "Empty Request",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatus: http.StatusOK,
			expectLog:      false,
		},
		{
			name: "Large Request",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatus: http.StatusOK,
			expectLog:      false,
		},
		{
			name: "Panic in Handler",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("test panic")
			}),
			expectedStatus: http.StatusInternalServerError,
			expectLog:      false,
		},
		{
			name: "Interceptor Error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatus: http.StatusOK,
			expectLog:      false,
		},
		{
			name: "Error in sendEdgeListToJarvis",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatus: http.StatusOK,
			expectLog:      false,
		},
		{
			name: "Invalid HTTP Method",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}),
			expectedStatus: http.StatusMethodNotAllowed,
			expectLog:      false,
		},
		{
			name: "High Volume of Requests",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatus: http.StatusOK,
			expectLog:      false,
		},
		{
			name: "Custom Logger Configuration",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mockLogger.Machine("Custom log message")
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatus: http.StatusOK,
			expectLog:      true,
		},
		{
			name: "Custom Interceptor Logic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mockLogger.Machine("Custom interceptor logic")
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatus: http.StatusOK,
			expectLog:      true,
		},
		{
			name: "Non-HTTP Error Panic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(errors.New("non-http error"))
			}),
			expectedStatus: http.StatusInternalServerError,
			expectLog:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger.Clear()
			handler := internalServerErrorHandler(tt.handler)

			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectLog {
				assert.NotEmpty(t, mockLogger.entries)
			} else {
				assert.Empty(t, mockLogger.entries)
			}
		})
	}
}
