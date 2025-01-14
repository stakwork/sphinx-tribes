package feeds_test

import (
	"github.com/stakwork/sphinx-tribes/feeds"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindGenerator(t *testing.T) {

	tests := []struct {
		name           string
		url            string
		serverResponse http.HandlerFunc
		expectedGen    int
		expectedBody   []byte
		expectError    bool
		errorContains  string
	}{
		{
			name: "Valid URL with Known Generator",
			url:  "http://example.com",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<rss><channel><generator>wordpress</generator></channel></rss>`))
			},
			expectedGen:  1,
			expectedBody: []byte(`<rss><channel><generator>wordpress</generator></channel></rss>`),
			expectError:  false,
		},
		{
			name: "Valid URL with Unknown Generator",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<rss><channel><generator>unknown</generator></channel></rss>`))
			},
			expectedGen:  0,
			expectedBody: []byte(`<rss><channel><generator>unknown</generator></channel></rss>`),
			expectError:  false,
		},
		{
			name: "URL with No Generator Tag",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<rss><channel></channel></rss>`))
			},
			expectedGen:  0,
			expectedBody: []byte(`<rss><channel></channel></rss>`),
			expectError:  false,
		},
		{
			name: "URL with Empty Generator Tag",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<rss><channel><generator></generator></channel></rss>`))
			},
			expectedGen:  0,
			expectedBody: []byte(`<rss><channel><generator></generator></channel></rss>`),
			expectError:  false,
		},
		{
			name: "Non-XML Response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`Not an XML`))
			},
			expectedGen:  0,
			expectedBody: []byte(`Not an XML`),
			expectError:  false,
		},
		{
			name: "Malformed XML Response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<rss><channel><generator>wordpress</generator>`))
			},
			expectedGen:  0,
			expectedBody: []byte(`<rss><channel><generator>wordpress</generator>`),
			expectError:  false,
		},
		{
			name: "Large XML Feed",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<rss><channel><generator>wordpress</generator></channel></rss>`))
			},
			expectedGen:  1,
			expectedBody: []byte(`<rss><channel><generator>wordpress</generator></channel></rss>`),
			expectError:  false,
		},
		{
			name: "Multiple Generator Tags",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<rss><channel><generator>unknown</generator><generator>wordpress</generator></channel></rss>`))
			},
			expectedGen:  1,
			expectedBody: []byte(`<rss><channel><generator>unknown</generator><generator>wordpress</generator></channel></rss>`),
			expectError:  false,
		},
		{
			name: "Case Sensitivity in Generator Tag",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<rss><channel><generator>WordPress</generator></channel></rss>`))
			},
			expectedGen:  0,
			expectedBody: []byte(`<rss><channel><generator>WordPress</generator></channel></rss>`),
			expectError:  false,
		},
		{
			name: "Generator Tag with Additional Text",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<rss><channel><generator>powered by wordpress</generator></channel></rss>`))
			},
			expectedGen:  1,
			expectedBody: []byte(`<rss><channel><generator>powered by wordpress</generator></channel></rss>`),
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

			gen, body, err := feeds.FindGenerator(url)

			assert.Equal(t, tt.expectedGen, gen)
			assert.Equal(t, tt.expectedBody, body)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
