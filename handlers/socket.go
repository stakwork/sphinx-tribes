package handlers

import (
	"net/http"

	"github.com/stakwork/sphinx-tribes/websocket"
)

// var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
// 	if config.Host == "https://people.sphinx.chat" {
// 		if r.Host != "people.sphinx.chat" && r.Host != "people-test.sphinx.chat" {
// 			return false
// 		} else {
// 			return true
// 		}
// 	}

// 	return true
// }}
// var Socket *websocket.Conn

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	pool := websocket.NewPool()
	go pool.Start()

	websocket.ServeWs(pool, w, r)

	// conn, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	log.Print("upgrade failed: ", err)
	// 	return
	// }

	// Socket = conn
	// defer conn.Close()

	// for {
	// 	_, message, err := conn.ReadMessage()
	// 	input := string(message)

	// 	fmt.Println("Websocket message ==", input)

	// 	if err != nil {
	// 		log.Println("read failed:", err)
	// 		break
	// 	}
	// }
}
