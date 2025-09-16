package ws

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Role string // examer / student
	Send chan []byte
}
