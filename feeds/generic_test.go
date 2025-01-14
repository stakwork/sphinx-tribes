package feeds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttpGet(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		url            string
		expectedBody   string
		expectError    bool
		errorContains  string
		setup          func(*httptest.Server) string
	}{
		{
			name: "Successful GET request with JSON response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"status": "success"}`)
			},
			expectedBody: `{"status": "success"}`,
			expectError:  false,
		},
		{
			name: "Server returns 404",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectedBody: "",
			expectError:  false,
		},
		{
			name: "Server returns 500",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "Internal Server Error")
			},
			expectedBody: "Internal Server Error",
			expectError:  false,
		},
		{
			name: "Large response body",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				largeBody := strings.Repeat("a", 1024*1024)
				fmt.Fprint(w, largeBody)
			},
			expectedBody: strings.Repeat("a", 1024*1024),
			expectError:  false,
		},
		{
			name: "Slow response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(100 * time.Millisecond)
				fmt.Fprint(w, "Delayed response")
			},
			expectedBody: "Delayed response",
			expectError:  false,
		},
		{
			name: "Empty response body",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
			},
			expectedBody: "",
			expectError:  false,
		},
		{
			name: "Response with special characters",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "Special chars: ñ, é, 漢字, 🌟")
			},
			expectedBody: "Special chars: ñ, é, 漢字, 🌟",
			expectError:  false,
		},
		{
			name:          "Invalid URL with invalid port",
			url:           "http://invalid.localhost:99999",
			expectedBody:  "",
			expectError:   true,
			errorContains: "dial tcp",
		},
		{
			name:          "Invalid URL with non-existent domain",
			url:           "http://this-domain-does-not-exist-123456789.com",
			expectedBody:  "",
			expectError:   true,
			errorContains: "no such host",
		},
		{
			name:          "Invalid URL with invalid protocol",
			url:           "invalid://example.com",
			expectedBody:  "",
			expectError:   true,
			errorContains: "unsupported protocol",
		},
		{
			name: "Connection closed by server",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				hj, ok := w.(http.Hijacker)
				if !ok {
					t.Fatal("webserver doesn't support hijacking")
				}
				conn, _, err := hj.Hijack()
				if err != nil {
					t.Fatal(err)
				}
				conn.Close()
			},
			expectedBody:  "",
			expectError:   true,
			errorContains: "EOF",
		},
		{
			name: "Chunked response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				flusher, ok := w.(http.Flusher)
				if !ok {
					t.Fatal("webserver doesn't support flushing")
				}
				fmt.Fprint(w, "chunk1")
				flusher.Flush()
				fmt.Fprint(w, "chunk2")
			},
			expectedBody: "chunk1chunk2",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server
			if tt.serverResponse != nil {
				server = httptest.NewServer(http.HandlerFunc(tt.serverResponse))
				defer server.Close()
			}

			url := tt.url
			if server != nil {
				url = server.URL
			}

			body, err := httpget(url)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, string(body))
			}
		})
	}
}

func TestAddedValue(t *testing.T) {
	tests := []struct {
		name             string
		value            *Value
		tribeOwnerPubkey string
		expected         *Value
	}{
		{
			name: "Empty tribe owner pubkey",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "existing_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "existing_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
		},
		{
			name:             "Nil value with tribe owner",
			value:            nil,
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.00000015000",
				},
				Destinations: []Destination{
					{
						Address: "tribe_owner_key",
						Type:    "node",
						Split:   100,
					},
				},
			},
		},
		{
			name: "Single destination with split 1 and different address",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "different_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "different_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
					{
						Address: "tribe_owner_key",
						Split:   99,
						Type:    "node",
					},
				},
			},
		},
		{
			name: "Single destination with split 1 and same address",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "tribe_owner_key",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "tribe_owner_key",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
		},
		{
			name: "Multiple destinations",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("50"),
						Type:    "node",
					},
					{
						Address: "address2",
						Split:   json.Number("50"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("50"),
						Type:    "node",
					},
					{
						Address: "address2",
						Split:   json.Number("50"),
						Type:    "node",
					},
				},
			},
		},
		{
			name: "Empty destinations",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{},
			},
		},
		{
			name: "Nil destinations",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: nil,
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: nil,
			},
		},
		{
			name: "Split as Zero",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("0"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("0"),
						Type:    "node",
					},
				},
			},
		},
		{
			name: "Split as Negative Number",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("-1"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("-1"),
						Type:    "node",
					},
				},
			},
		},
		{
			name: "tribeOwnerPubkey as Special Characters",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "different_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "!@#$%^&*()_+",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "different_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
					{
						Address: "!@#$%^&*()_+",
						Split:   99,
						Type:    "node",
					},
				},
			},
		},
		{
			name: "Split as a Large Number",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("999999999999"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "address1",
						Split:   json.Number("999999999999"),
						Type:    "node",
					},
				},
			},
		},
		{
			name: "tribeOwnerPubkey as a Very Long String",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "different_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
				},
			},
			tribeOwnerPubkey: strings.Repeat("a", 1000),
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: []Destination{
					{
						Address: "different_address",
						Split:   json.Number("1"),
						Type:    "node",
					},
					{
						Address: strings.Repeat("a", 1000),
						Split:   99,
						Type:    "node",
					},
				},
			},
		},
		{
			name: "Large Number of Destinations",
			value: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: func() []Destination {
					dests := make([]Destination, 100)
					for i := 0; i < 100; i++ {
						dests[i] = Destination{
							Address: fmt.Sprintf("address_%d", i),
							Split:   json.Number(fmt.Sprintf("%d", i+1)),
							Type:    "node",
						}
					}
					return dests
				}(),
			},
			tribeOwnerPubkey: "tribe_owner_key",
			expected: &Value{
				Model: Model{
					Type:      "lightning",
					Suggested: "0.0001",
				},
				Destinations: func() []Destination {
					dests := make([]Destination, 100)
					for i := 0; i < 100; i++ {
						dests[i] = Destination{
							Address: fmt.Sprintf("address_%d", i),
							Split:   json.Number(fmt.Sprintf("%d", i+1)),
							Type:    "node",
						}
					}
					return dests
				}(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddedValue(tt.value, tt.tribeOwnerPubkey)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)

			assert.Equal(t, tt.expected.Model.Type, result.Model.Type)
			assert.Equal(t, tt.expected.Model.Suggested, result.Model.Suggested)

			assert.Equal(t, len(tt.expected.Destinations), len(result.Destinations))
			for i, expectedDest := range tt.expected.Destinations {
				assert.Equal(t, expectedDest.Address, result.Destinations[i].Address)
				assert.Equal(t, expectedDest.Type, result.Destinations[i].Type)
				assert.Equal(t, expectedDest.Split, result.Destinations[i].Split)
				assert.Equal(t, expectedDest.CustomKey, result.Destinations[i].CustomKey)
				assert.Equal(t, expectedDest.CustomValue, result.Destinations[i].CustomValue)
			}
		})
	}
}
