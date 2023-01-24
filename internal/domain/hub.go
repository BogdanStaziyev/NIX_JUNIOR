package domain

import (
	"fmt"
	"github.com/labstack/gommon/log"
)

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
			str := fmt.Sprintf("User %s entered the chat", client.Name)
			h.broadcastToClients([]byte(str))
			log.Infof("Register new client", client.ID)
		case client := <-h.Unregister:
			h.unregisterClient(client)
			str := fmt.Sprintf("User %s left the chat", client.Name)
			h.broadcastToClients([]byte(str))
			log.Infof("Unregister new client", client.ID)
		case message := <-h.Broadcast:
			h.broadcastToClients(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.Clients[client] = true
}

func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.Clients[client]; ok {
		delete(h.Clients, client)
		close(client.Send)
	}
}

func (h *Hub) broadcastToClients(message []byte) {
	for client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.Clients, client)
		}
	}
}
