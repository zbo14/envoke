package types

import (
	ws "github.com/gorilla/websocket"
	"net/http"
)

func Upgrader() *ws.Upgrader {
	return &ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(req *http.Request) bool {
			return true
		},
	}
}
