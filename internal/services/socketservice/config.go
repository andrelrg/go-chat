package socketservice

import (
	"github.com/gorilla/websocket"
	"net/http"
)

func Socket() {
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	Connected = make(map[int64][]*SocketUtils)
}

var Upgrader websocket.Upgrader

var Connected map[int64][]*SocketUtils
