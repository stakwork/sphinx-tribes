package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}

	conn.WriteJSON("hello This")

	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()

		input := string(message)
		fmt.Println("Message ==", input)

		if err != nil {
			log.Println("read failed:", err)
			break
		}

		message = []byte("Hello Now ======")
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write failed:", err)
			break
		}

	}
}
