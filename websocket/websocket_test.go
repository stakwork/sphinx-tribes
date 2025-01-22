package websocket

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stretchr/testify/assert"
)

func TestServeWs(t *testing.T) {
	t.Run("Valid Connection with Unique ID", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pool := NewPool()
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=test123"

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		defer ws.Close()

		assert.NotNil(t, ws)
	})

	t.Run("Connection with Empty Unique ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pool := NewPool()
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		defer ws.Close()

		assert.NotNil(t, ws)
	})

	t.Run("Connection with 'null' Unique ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pool := NewPool()
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=null"
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		defer ws.Close()

		assert.NotNil(t, ws)
	})

	t.Run("Connection with 'undefined' Unique ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pool := NewPool()
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=undefined"
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		defer ws.Close()

		assert.NotNil(t, ws)
	})

	t.Run("Multiple Connections with Same Unique ID", func(t *testing.T) {
		pool := NewPool()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=same123"

		ws1, _, err1 := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err1)
		defer ws1.Close()

		ws2, _, err2 := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err2)
		defer ws2.Close()

		assert.NotNil(t, ws1)
		assert.NotNil(t, ws2)
	})

	t.Run("Connection with Special Characters in Unique ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pool := NewPool()
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		uniqueID := url.QueryEscape("test@123!#$%")
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=" + uniqueID

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		defer ws.Close()

		assert.NotNil(t, ws)
	})

	t.Run("Connection with Various Special Characters in Unique ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pool := NewPool()
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		testCases := []string{
			"test@email.com",
			"test#123",
			"test&user",
			"test+plus",
			"test space",
			"test/slash",
			"test?query",
			"test=equals",
			"test:colon",
		}

		for _, testID := range testCases {
			t.Run(testID, func(t *testing.T) {
				uniqueID := url.QueryEscape(testID)
				wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=" + uniqueID

				ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
				assert.NoError(t, err)
				if ws != nil {
					defer ws.Close()
				}
				assert.NotNil(t, ws)
			})
		}
	})

	t.Run("Connection with Unicode Characters in Unique ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pool := NewPool()
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		testCases := []string{
			"test🚀rocket",
			"test👨‍💻coder",
			"test❤️heart",
			"testñáéíóú",
			"test中文",
			"testالعربية",
			"testРусский",
		}

		for _, testID := range testCases {
			t.Run(testID, func(t *testing.T) {
				uniqueID := url.QueryEscape(testID)
				wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=" + uniqueID

				ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
				assert.NoError(t, err)
				if ws != nil {
					defer ws.Close()
				}
				assert.NotNil(t, ws)
			})
		}
	})

	t.Run("Connection with Maximum Length Unique ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pool := NewPool()
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		longID := strings.Repeat("test@123!#$%", 20)
		uniqueID := url.QueryEscape(longID)
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=" + uniqueID

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		defer ws.Close()

		assert.NotNil(t, ws)
	})

	t.Run("Connection with Very Long Unique ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pool := NewPool()
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		longID := strings.Repeat("a", 1000)
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=" + longID
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		defer ws.Close()

		assert.NotNil(t, ws)
	})

	t.Run("Origin Check for Production Host", func(t *testing.T) {
		originalHost := config.Host
		config.Host = "https://people.sphinx.chat"
		defer func() { config.Host = originalHost }()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pool := NewPool()
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=test789"

		header := http.Header{}
		header.Add("Host", "people.sphinx.chat")

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, header)
		assert.NoError(t, err)
		defer ws.Close()

		assert.NotNil(t, ws)
	})

	t.Run("Valid uniqueId Provided", func(t *testing.T) {
		pool := NewPool()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=valid-test-id"
		ws, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		assert.NotNil(t, ws)
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		defer ws.Close()
	})

	t.Run("No uniqueId Provided", func(t *testing.T) {
		pool := NewPool()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
		ws, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		assert.NotNil(t, ws)
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		defer ws.Close()
	})

	t.Run("uniqueId is null", func(t *testing.T) {
		pool := NewPool()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=null"
		ws, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		assert.NotNil(t, ws)
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		defer ws.Close()
	})

	t.Run("uniqueId is undefined", func(t *testing.T) {
		pool := NewPool()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=undefined"
		ws, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		assert.NotNil(t, ws)
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		defer ws.Close()
	})

	t.Run("Empty uniqueId After Random Generation", func(t *testing.T) {
		pool := NewPool()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId="
		ws, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		assert.NotNil(t, ws)
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		defer ws.Close()
	})

	t.Run("WebSocket Upgrade Failure", func(t *testing.T) {
		pool := NewPool()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Connection", "close")
			w.WriteHeader(http.StatusBadRequest)
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=test-id"
		ws, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.Error(t, err)
		assert.Nil(t, ws)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Invalid HTTP Method", func(t *testing.T) {
		pool := NewPool()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		url := server.URL + "?uniqueId=test-id"

		req, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Various Invalid HTTP Methods", func(t *testing.T) {
		pool := NewPool()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		invalidMethods := []string{
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodOptions,
			http.MethodHead,
		}

		for _, method := range invalidMethods {
			t.Run(method, func(t *testing.T) {
				url := server.URL + "?uniqueId=test-id"

				req, err := http.NewRequest(method, url, nil)
				assert.NoError(t, err)

				resp, err := http.DefaultClient.Do(req)
				assert.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
			})
		}
	})

	t.Run("Valid GET Method", func(t *testing.T) {
		pool := NewPool()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?uniqueId=test-id"
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		if ws != nil {
			defer ws.Close()
		}
		assert.NotNil(t, ws)
	})

	t.Run("Method Override Attempt", func(t *testing.T) {
		pool := NewPool()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		url := server.URL + "?uniqueId=test-id"

		req, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		req.Header.Set("X-HTTP-Method-Override", "GET")
		req.Header.Set("X-Method-Override", "GET")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Custom Method String", func(t *testing.T) {
		pool := NewPool()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		url := server.URL + "?uniqueId=test-id"

		req, err := http.NewRequest("CUSTOM", url, nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Concurrent Connections with Different IDs", func(t *testing.T) {
		pool := NewPool()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ServeWs(pool, w, r)
		}))
		defer server.Close()

		for i := 0; i < 5; i++ {
			wsURL := fmt.Sprintf("ws%s?uniqueId=test-id-%d",
				strings.TrimPrefix(server.URL, "http"), i)
			ws, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
			assert.NoError(t, err)
			assert.NotNil(t, ws)
			assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
			defer ws.Close()
		}
	})

}
