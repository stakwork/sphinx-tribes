package main

import (
	"github.com/olahol/melody"
)

var socket = melody.New()

func init() {
	socket.HandleMessage(func(s *melody.Session, msg []byte) {
		socket.Broadcast(msg)
	})
}
