package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/stakwork/sphinx-tribes/db"
)

type Client struct {
	Host string
	Conn *websocket.Conn
	Pool *Pool
}

type ClientData struct {
	Client *Client
	Status bool
}

type Message struct {
	Type int    `json:"type"`
	Msg  string `json:"msg"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
		db.Store.DeleteCache(c.Host)
	}()

	for {
		var socketMsg db.LnHost
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		err = json.Unmarshal(p, &socketMsg)
		if err != nil {
			fmt.Println("Message Decode Error", err, string(p))
		}
		message := Message{Type: messageType, Body: string(p)}

		fmt.Printf("Message Received: %+v\n", message)
		c.Pool.Broadcast <- message
	}
}
