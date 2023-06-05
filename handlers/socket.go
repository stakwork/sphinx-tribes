package handlers

import (
	"net/http"

	"github.com/stakwork/sphinx-tribes/websocket"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	pool := websocket.NewPool()
	go pool.Start()

	websocket.ServeWs(pool, w, r)
}
