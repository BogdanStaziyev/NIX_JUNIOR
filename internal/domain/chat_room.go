package domain

import "C"
import (
	"fmt"
	"github.com/labstack/gommon/log"
)

type ChatRoom struct {
	ChatName   string
	Users      map[*ChatUser]bool
	Connect    chan *ChatUser
	Disconnect chan *ChatUser
	Broadcast  chan []byte
}

func (c *ChatRoom) New(name string) *ChatRoom {
	return &ChatRoom{
		ChatName:   name,
		Users:      make(map[*ChatUser]bool),
		Connect:    make(chan *ChatUser),
		Disconnect: make(chan *ChatUser),
		Broadcast:  make(chan []byte),
	}
}

func (c *ChatRoom) Run() {
	for {
		select {
		case client := <-c.Connect:
			c.registerClient(client)
			str := fmt.Sprintf("CHATROOM User %d entered the chat", client.ID)
			c.broadcastToClients([]byte(str))
			log.Infof("CHATROOM Register new client %d", client.ID)
		case client := <-c.Disconnect:
			c.unregisterClient(client)
			str := fmt.Sprintf("CHATROOM User %d left the chat", client.ID)
			c.broadcastToClients([]byte(str))
			log.Infof("CHATROOM Unregister new client %d", client.ID)
		case message := <-c.Broadcast:
			c.broadcastToClients(message)
		}
	}
}

func (c *ChatRoom) registerClient(client *ChatUser) {
	c.Users[client] = true
	log.Info(c.Users)
}

func (c *ChatRoom) unregisterClient(client *ChatUser) {
	if _, ok := c.Users[client]; ok {
		delete(c.Users, client)
		close(client.Send)
	}
	log.Info(c.Users)
}

func (c *ChatRoom) broadcastToClients(message []byte) {
	for client := range c.Users {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(c.Users, client)
		}
	}
}
