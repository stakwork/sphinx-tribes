package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildV2KeysendBodyData(t *testing.T) {
	tests := []struct {
		name           string
		amount         uint
		receiverPubkey string
		routeHint      string
		memo           string
		expected       string
		expectError    bool
	}{
		{
			name:           "Standard Input with Route Hint",
			amount:         100,
			receiverPubkey: "abcdef123456",
			routeHint:      "hint123",
			memo:           "Test transaction",
			expected:       `{"amt_msat": 100000, "dest": "abcdef123456", "route_hint": "hint123", "data": "Test transaction", "wait": true}`,
		},
		{
			name:           "Standard Input without Route Hint",
			amount:         100,
			receiverPubkey: "abcdef123456",
			routeHint:      "",
			memo:           "Test transaction",
			expected:       `{"amt_msat": 100000, "dest": "abcdef123456", "route_hint": "", "data": "Test transaction", "wait": true}`,
		},
		{
			name:           "Minimum Amount",
			amount:         0,
			receiverPubkey: "abcdef123456",
			routeHint:      "hint123",
			memo:           "Test transaction",
			expected:       `{"amt_msat": 0, "dest": "abcdef123456", "route_hint": "hint123", "data": "Test transaction", "wait": true}`,
		},
		{
			name:           "Empty Strings for All String Parameters",
			amount:         100,
			receiverPubkey: "",
			routeHint:      "",
			memo:           "",
			expected:       `{"amt_msat": 100000, "dest": "", "route_hint": "", "data": "", "wait": true}`,
		},
		{
			name:           "Large Amount",
			amount:         4294967295,
			receiverPubkey: "abcdef123456",
			routeHint:      "hint123",
			memo:           "Test transaction",
			expected:       `{"amt_msat": 4294967295000, "dest": "abcdef123456", "route_hint": "hint123", "data": "Test transaction", "wait": true}`,
		},
		{
			name:           "Long Strings",
			amount:         100,
			receiverPubkey: strings.Repeat("a", 1000),
			routeHint:      strings.Repeat("b", 1000),
			memo:           strings.Repeat("c", 1000),
			expected:       fmt.Sprintf(`{"amt_msat": 100000, "dest": "%s", "route_hint": "%s", "data": "%s", "wait": true}`, strings.Repeat("a", 1000), strings.Repeat("b", 1000), strings.Repeat("c", 1000)),
		},
		{
			name:           "Special Characters in Strings",
			amount:         100,
			receiverPubkey: "abc!@#123",
			routeHint:      "hint$%^",
			memo:           "Test &*() transaction",
			expected:       `{"amt_msat": 100000, "dest": "abc!@#123", "route_hint": "hint$%^", "data": "Test &*() transaction", "wait": true}`,
		},
		{
			name:           "Whitespace in Strings",
			amount:         100,
			receiverPubkey: "abcdef123456",
			routeHint:      "hint123  ",
			memo:           "  Test transaction  ",
			expected:       `{"amt_msat": 100000, "dest": "abcdef123456", "route_hint": "hint123", "data": "Test transaction", "wait": true}`,
		},
		{
			name:           "Non-ASCII Characters in Strings",
			amount:         100,
			receiverPubkey: "abcñ123",
			routeHint:      "hintü123",
			memo:           "Test transaction with emoji 😊",
			expected:       `{"amt_msat": 100000, "dest": "abcñ123", "route_hint": "hintü123", "data": "Test transaction with emoji 😊", "wait": true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildV2KeysendBodyData(tt.amount, tt.receiverPubkey, tt.routeHint, tt.memo)
			assert.JSONEq(t, tt.expected, result)
		})
	}
}

func TestGetPaginationParams(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		expectedOffset int
		expectedLimit  int
		expectedSortBy string
		expectedDir    string
		expectedSearch string
	}{
		{
			name:           "Standard Input with All Parameters Present",
			query:          "page=2&limit=10&sortBy=created&direction=asc&search=test",
			expectedOffset: 10,
			expectedLimit:  10,
			expectedSortBy: "created",
			expectedDir:    "asc",
			expectedSearch: "test",
		},
		{
			name:           "Standard Input with Default Parameters",
			query:          "page=1&limit=5",
			expectedOffset: 0,
			expectedLimit:  5,
			expectedSortBy: "created",
			expectedDir:    "desc",
			expectedSearch: "",
		},
		{
			name:           "Minimum Values for Page and Limit",
			query:          "page=1&limit=1",
			expectedOffset: 0,
			expectedLimit:  1,
			expectedSortBy: "created",
			expectedDir:    "desc",
			expectedSearch: "",
		},
		{
			name:           "Zero Values for Page and Limit",
			query:          "page=0&limit=0",
			expectedOffset: 0,
			expectedLimit:  1,
			expectedSortBy: "created",
			expectedDir:    "desc",
			expectedSearch: "",
		},
		{
			name:           "Negative Values for Page and Limit",
			query:          "page=-1&limit=-10",
			expectedOffset: 0,
			expectedLimit:  -10,
			expectedSortBy: "created",
			expectedDir:    "desc",
			expectedSearch: "",
		},
		{
			name:           "Non-Integer Values for Page and Limit",
			query:          "page=abc&limit=xyz",
			expectedOffset: 0,
			expectedLimit:  1,
			expectedSortBy: "created",
			expectedDir:    "desc",
			expectedSearch: "",
		},
		{
			name:           "Null Request",
			query:          "",
			expectedOffset: 0,
			expectedLimit:  1,
			expectedSortBy: "created",
			expectedDir:    "desc",
			expectedSearch: "",
		},
		{
			name:           "Large Values for Page and Limit",
			query:          "page=10000&limit=1000",
			expectedOffset: 9999000,
			expectedLimit:  1000,
			expectedSortBy: "created",
			expectedDir:    "desc",
			expectedSearch: "",
		},
		{
			name:           "Empty Query Parameters",
			query:          "",
			expectedOffset: 0,
			expectedLimit:  1,
			expectedSortBy: "created",
			expectedDir:    "desc",
			expectedSearch: "",
		},
		{
			name:           "Only Search Parameter Present",
			query:          "search=example",
			expectedOffset: 0,
			expectedLimit:  1,
			expectedSortBy: "created",
			expectedDir:    "desc",
			expectedSearch: "example",
		},
		{
			name:           "Invalid Direction Value",
			query:          "direction=upwards",
			expectedOffset: 0,
			expectedLimit:  1,
			expectedSortBy: "created",
			expectedDir:    "desc",
			expectedSearch: "",
		},
		{
			name:           "Invalid SortBy Value",
			query:          "sortBy=unknown",
			expectedOffset: 0,
			expectedLimit:  1,
			expectedSortBy: "unknown",
			expectedDir:    "desc",
			expectedSearch: "",
		},
		{
			name:           "Mixed Valid and Invalid Parameters",
			query:          "page=3&limit=abc&sortBy=unknown&direction=upwards",
			expectedOffset: 2,
			expectedLimit:  1,
			expectedSortBy: "unknown",
			expectedDir:    "desc",
			expectedSearch: "",
		},
		{
			name:           "Whitespace in Parameters",
			query:          "page= 2 &limit= 10&sortBy= created &direction= asc ",
			expectedOffset: 10,
			expectedLimit:  10,
			expectedSortBy: "created",
			expectedDir:    "asc",
			expectedSearch: "",
		},
		{
			name:           "Case Sensitivity in Parameters",
			query:          "sortBy=CREATED&direction=ASC",
			expectedOffset: 0,
			expectedLimit:  1,
			expectedSortBy: "created",
			expectedDir:    "asc",
			expectedSearch: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				URL: &url.URL{
					RawQuery: tt.query,
				},
			}
			offset, limit, sortBy, direction, search := GetPaginationParams(req)
			assert.Equal(t, tt.expectedOffset, offset)
			assert.Equal(t, tt.expectedLimit, limit)
			assert.Equal(t, tt.expectedSortBy, sortBy)
			assert.Equal(t, tt.expectedDir, direction)
			assert.Equal(t, tt.expectedSearch, search)
		})
	}
}
