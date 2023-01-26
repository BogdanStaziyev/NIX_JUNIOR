package domain

import (
	"github.com/gorilla/websocket"
)

type ChatUser struct {
	ID       int64
	ChatName string
	Name     string
	Conn     *websocket.Conn
	Chat     *ChatRoom
	Hub      *Hub
	Send     chan []byte
}

func NewClient(conn *websocket.Conn, hub *Hub) *ChatUser {
	c := &ChatUser{
		Conn: conn,
		Hub:  hub,
		Send: make(chan []byte, 1024),
	}
	return c
}

func (c *ChatUser) Disconnect() {
	c.Hub.Unregister <- c
	c.Conn.Close()
}
