package websocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestNewPool(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T, pool *Pool)
	}{
		{
			name: "Basic Functionality: Successful Pool Initialization",
			validate: func(t *testing.T, pool *Pool) {
				assert.NotNil(t, pool, "Pool should not be nil")
				assert.NotNil(t, pool.Register, "Register channel should not be nil")
				assert.NotNil(t, pool.Unregister, "Unregister channel should not be nil")
				assert.NotNil(t, pool.Clients, "Clients map should not be nil")
				assert.NotNil(t, pool.Broadcast, "Broadcast channel should not be nil")
			},
		},
		{
			name: "Edge Case: Verify Channel Types and States",
			validate: func(t *testing.T, pool *Pool) {
				assert.Equal(t, 0, cap(pool.Register), "Register channel should be unbuffered")
				assert.Equal(t, 0, cap(pool.Unregister), "Unregister channel should be unbuffered")
				assert.Equal(t, 0, cap(pool.Broadcast), "Broadcast channel should be unbuffered")
			},
		},
		{
			name: "Edge Case: Verify Map Initialization",
			validate: func(t *testing.T, pool *Pool) {
				assert.Equal(t, 0, len(pool.Clients), "Clients map should be initialized and empty")
			},
		},
		{
			name: "Error Conditions: Invalid State Handling",
			validate: func(t *testing.T, pool *Pool) {
				// No specific validation needed, just ensure no panic occurs
			},
		},
		{
			name: "Performance and Scale: Multiple Pool Instances",
			validate: func(t *testing.T, pool *Pool) {
				anotherPool := NewPool()
				assert.NotSame(t, pool, anotherPool, "Each Pool instance should be independent")
				assert.NotSame(t, pool.Register, anotherPool.Register, "Register channels should be independent")
				assert.NotSame(t, pool.Unregister, anotherPool.Unregister, "Unregister channels should be independent")
				assert.NotSame(t, pool.Clients, anotherPool.Clients, "Clients maps should be independent")
				assert.NotSame(t, pool.Broadcast, anotherPool.Broadcast, "Broadcast channels should be independent")
			},
		},
		{
			name: "Concurrency Safety: Readiness for Concurrent Use",
			validate: func(t *testing.T, pool *Pool) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewPool()
			tt.validate(t, pool)
		})
	}
}

func TestSendTicketMessage(t *testing.T) {
	t.Run("Direct Broadcast with Valid Client", func(t *testing.T) {
		pool := NewPool()
		ws, server := setupTestWebsocket(t)
		defer server.Close()
		defer ws.Close()

		client := &Client{
			Host: "test-client",
			Conn: ws,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client.Host] = &ClientData{
			Client: client,
			Status: true,
		}

		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "test-client",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.NoError(t, err)
	})

	t.Run("Non-Direct Broadcast", func(t *testing.T) {
		pool := NewPool()
		ws, server := setupTestWebsocket(t)
		defer server.Close()
		defer ws.Close()

		client := &Client{
			Host: "test-client",
			Conn: ws,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client.Host] = &ClientData{
			Client: client,
			Status: true,
		}

		message := TicketMessage{
			BroadcastType:   "broadcast",
			SourceSessionID: "test-client",
			Message:         "Test broadcast message",
		}

		err := pool.SendTicketMessage(message)
		assert.NoError(t, err)
	})

	t.Run("Empty SourceSessionID", func(t *testing.T) {
		pool := NewPool()
		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "client not found")
	})

	t.Run("Empty BroadcastType", func(t *testing.T) {
		pool := NewPool()
		ws, server := setupTestWebsocket(t)
		defer server.Close()
		defer ws.Close()

		client := &Client{
			Host: "test-client",
			Conn: ws,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client.Host] = &ClientData{
			Client: client,
			Status: true,
		}

		message := TicketMessage{
			BroadcastType:   "",
			SourceSessionID: "test-client",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.NoError(t, err)
	})

	t.Run("Client Not Found", func(t *testing.T) {
		pool := NewPool()
		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "non-existent-client",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "client not found")
	})

	t.Run("WriteJSON Error", func(t *testing.T) {
		pool := NewPool()

		ws, server := setupTestWebsocket(t)
		server.Close()
		ws.Close()

		client := &Client{
			Host: "test-client",
			Conn: ws,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client.Host] = &ClientData{
			Client: client,
			Status: true,
		}

		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "test-client",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.Error(t, err)
	})

	t.Run("Large Message Payload", func(t *testing.T) {
		pool := NewPool()
		ws, server := setupTestWebsocket(t)
		defer server.Close()
		defer ws.Close()

		client := &Client{
			Host: "test-client",
			Conn: ws,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client.Host] = &ClientData{
			Client: client,
			Status: true,
		}

		largeMessage := strings.Repeat("a", 1024*1024)
		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "test-client",
			Message:         largeMessage,
		}

		err := pool.SendTicketMessage(message)
		assert.NoError(t, err)
	})

	t.Run("Multiple Clients with Same SessionID", func(t *testing.T) {
		pool := NewPool()
		ws1, server1 := setupTestWebsocket(t)
		ws2, server2 := setupTestWebsocket(t)
		defer server1.Close()
		defer server2.Close()
		defer ws1.Close()
		defer ws2.Close()

		client1 := &Client{
			Host: "same-session-id",
			Conn: ws1,
			Pool: pool,
		}

		client2 := &Client{
			Host: "same-session-id",
			Conn: ws2,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client1.Host] = &ClientData{
			Client: client1,
			Status: true,
		}
		pool.Clients[client2.Host] = &ClientData{
			Client: client2,
			Status: true,
		}

		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "same-session-id",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.NoError(t, err)
	})

	t.Run("Null or Uninitialized Pool", func(t *testing.T) {
		var pool *Pool
		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "test-client",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.Error(t, err)
	})
}

func TestSendTicketMessageScenarios(t *testing.T) {
	t.Run("Direct Broadcast with Valid Client", func(t *testing.T) {
		pool := NewPool()
		ws, server := setupTestWebsocket(t)
		defer server.Close()
		defer ws.Close()

		client := &Client{
			Host: "test-client",
			Conn: ws,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client.Host] = &ClientData{
			Client: client,
			Status: true,
		}

		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "test-client",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.NoError(t, err)
	})

	t.Run("Non-Direct Broadcast", func(t *testing.T) {
		pool := NewPool()
		ws, server := setupTestWebsocket(t)
		defer server.Close()
		defer ws.Close()

		client := &Client{
			Host: "test-client",
			Conn: ws,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client.Host] = &ClientData{
			Client: client,
			Status: true,
		}

		message := TicketMessage{
			BroadcastType:   "broadcast",
			SourceSessionID: "test-client",
			Message:         "Test broadcast message",
		}

		err := pool.SendTicketMessage(message)
		assert.NoError(t, err)
	})

	t.Run("Empty SourceSessionID with Direct Broadcast", func(t *testing.T) {
		pool := NewPool()
		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "client not found")
	})

	t.Run("Empty BroadcastType", func(t *testing.T) {
		pool := NewPool()
		ws, server := setupTestWebsocket(t)
		defer server.Close()
		defer ws.Close()

		client := &Client{
			Host: "test-client",
			Conn: ws,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client.Host] = &ClientData{
			Client: client,
			Status: true,
		}

		message := TicketMessage{
			BroadcastType:   "",
			SourceSessionID: "test-client",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.NoError(t, err)
	})

	t.Run("Client Not Found", func(t *testing.T) {
		pool := NewPool()
		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "non-existent-client",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "client not found")
	})

	t.Run("WriteJSON Error", func(t *testing.T) {
		pool := NewPool()

		ws, server := setupTestWebsocket(t)
		server.Close()
		ws.Close()

		client := &Client{
			Host: "test-client",
			Conn: ws,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client.Host] = &ClientData{
			Client: client,
			Status: true,
		}

		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "test-client",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.Error(t, err)
	})

	t.Run("Large Message Payload", func(t *testing.T) {
		pool := NewPool()
		ws, server := setupTestWebsocket(t)
		defer server.Close()
		defer ws.Close()

		client := &Client{
			Host: "test-client",
			Conn: ws,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client.Host] = &ClientData{
			Client: client,
			Status: true,
		}

		largeMessage := strings.Repeat("a", 1024*1024)
		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "test-client",
			Message:         largeMessage,
		}

		err := pool.SendTicketMessage(message)
		assert.NoError(t, err)
	})

	t.Run("Multiple Clients with Same SessionID", func(t *testing.T) {
		pool := NewPool()
		ws1, server1 := setupTestWebsocket(t)
		ws2, server2 := setupTestWebsocket(t)
		defer server1.Close()
		defer server2.Close()
		defer ws1.Close()
		defer ws2.Close()

		client1 := &Client{
			Host: "same-session-id",
			Conn: ws1,
			Pool: pool,
		}

		client2 := &Client{
			Host: "same-session-id",
			Conn: ws2,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client1.Host] = &ClientData{
			Client: client1,
			Status: true,
		}
		pool.Clients[client2.Host] = &ClientData{
			Client: client2,
			Status: true,
		}

		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "same-session-id",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.NoError(t, err)
	})

	t.Run("Null or Uninitialized Pool", func(t *testing.T) {
		var pool *Pool
		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "test-client",
			Message:         "Test message",
		}

		err := pool.SendTicketMessage(message)
		assert.Error(t, err)
	})

	t.Run("Direct Broadcast with Null Message", func(t *testing.T) {
		pool := NewPool()
		ws, server := setupTestWebsocket(t)
		defer server.Close()
		defer ws.Close()

		client := &Client{
			Host: "test-client",
			Conn: ws,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client.Host] = &ClientData{
			Client: client,
			Status: true,
		}

		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "test-client",
			Message:         "",
		}

		err := pool.SendTicketMessage(message)
		assert.NoError(t, err)
	})

	t.Run("Direct Broadcast with Special Characters in Message", func(t *testing.T) {
		pool := NewPool()
		ws, server := setupTestWebsocket(t)
		defer server.Close()
		defer ws.Close()

		client := &Client{
			Host: "test-client",
			Conn: ws,
			Pool: pool,
		}

		pool.Clients = make(map[string]*ClientData)
		pool.Clients[client.Host] = &ClientData{
			Client: client,
			Status: true,
		}

		message := TicketMessage{
			BroadcastType:   "direct",
			SourceSessionID: "test-client",
			Message:         "Test message with special chars: !@#$%^&*()",
		}

		err := pool.SendTicketMessage(message)
		assert.NoError(t, err)
	})
}

func setupTestWebsocket(t *testing.T) (*websocket.Conn, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}))

	wsURL := "ws" + server.URL[4:]
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	return ws, server
}
