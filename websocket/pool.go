package websocket

import (
	"fmt"

	"github.com/stakwork/sphinx-tribes/db"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]*ClientData
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*ClientData),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			// ceck to acoid nil pointer
			if pool.Clients == nil {
				pool.Clients = make(map[string]*ClientData)
			}

			pool.Clients[client.Host] = &ClientData{
				Client: client,
				Status: true,
			}
			fmt.Println("Size of Websocket Connection Pool: ", len(pool.Clients))
			err := db.Store.SetSocketConnections(db.Client{
				Host: client.Host,
				Conn: client.Conn,
			})
			if err == nil {
				pool.Clients[client.Host].Client.Conn.WriteJSON(Message{Type: 1, Msg: "user_connect", Body: client.Host})
				go client.Read()
			} else {
				fmt.Println("Websocket pool client save error")
			}
		case client := <-pool.Unregister:
			// ceck to acoid nil pointer
			if pool.Clients[client.Host] != nil {
				pool.Clients[client.Host].Client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
				delete(pool.Clients, client.Host)
				fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			}

		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			// ceck to acoid nil pointer
			if pool.Clients != nil {
				for client, _ := range pool.Clients {
					if err := pool.Clients[client].Client.Conn.WriteJSON(message); err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		}
	}
}

func (pool *Pool) SendTicketMessage(message TicketMessage) error {

	if pool == nil {
		return fmt.Errorf("pool is nil")
	}

	if message.BroadcastType == "direct" {
		if message.SourceSessionID == "" {
			return fmt.Errorf("client not found")
		}
		// check if client
		if client, ok := pool.Clients[message.SourceSessionID]; ok {
			return client.Client.Conn.WriteJSON(message)
		}
		return fmt.Errorf("client not found: %s", message.SourceSessionID)
	}

	return nil
}


func (pool *Pool) SendTicketPlanMessage(message TicketPlanMessage) error {
    if pool == nil {
        return fmt.Errorf("pool is nil")
    }

    if message.BroadcastType == "direct" {
        if message.SourceSessionID == "" {
            return fmt.Errorf("client not found")
        }

		if client, ok := pool.Clients[message.SourceSessionID]; ok {
            return client.Client.Conn.WriteJSON(message)
        }
        return fmt.Errorf("client not found: %s", message.SourceSessionID)
    }

    return nil
}