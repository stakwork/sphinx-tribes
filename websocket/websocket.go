package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/utils"
)

var WebsocketPool = NewPool()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		if config.Host == "https://people.sphinx.chat" {
			if r.Host != "people.sphinx.chat" && r.Host != "people-test.sphinx.chat" {
				return false
			} else {
				return true
			}
		}

		return true
	}}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	websocketToken := utils.GetRandomToken(40)

	conn, err := Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &Client{
		Host: websocketToken,
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}
