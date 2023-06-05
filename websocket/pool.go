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
			fmt.Print("Host ===", client.Host)
			pool.Clients[client.Host] = &ClientData{
				Client: client,
				Status: true,
			}
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			pool.Clients[client.Host].Client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			db.Store.SetSocketConnections(db.Client{
				Host: client.Host,
				Conn: client.Conn,
			})
			db.Store.SetSocketConnections(db.Client{
				Host: "c058-102-88-63-132.eu.ngrok.io",
				Conn: client.Conn,
			})
			db.Store.SetSocketConnections(db.Client{
				Host: "localhost:5005",
				Conn: client.Conn,
			})
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client.Host)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := pool.Clients[client].Client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
