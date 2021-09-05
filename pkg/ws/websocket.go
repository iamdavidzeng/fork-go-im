package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	WebSocket *websocket.Conn
	err       error
)

func App(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	WebSocket, err = upgrader.Upgrade(w, r, nil)
	return WebSocket, err
}
