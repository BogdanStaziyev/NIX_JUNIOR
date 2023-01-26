package domain

import (
	"github.com/labstack/gommon/log"
)

type Hub struct {
	ClientsChat  map[string]*ChatRoom
	Broadcast    chan []byte
	RegisterChat chan *ChatUser
	Unregister   chan *ChatUser
}

func NewHub() *Hub {
	return &Hub{
		ClientsChat:  make(map[string]*ChatRoom),
		Broadcast:    make(chan []byte),
		RegisterChat: make(chan *ChatUser),
		Unregister:   make(chan *ChatUser),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.RegisterChat:
			h.registerClient(client)
			log.Errorf("HUB Find new client %s", client.Name)
		case client := <-h.Unregister:
			h.unregisterClient(client)
			log.Errorf("HUB Unregister client %d", client.ID)
		case message := <-h.Broadcast:
			h.broadcastToClients(message)
		}
	}
}

func (h *Hub) registerClient(client *ChatUser) {
	if ch, ok := h.ClientsChat[client.ChatName]; ok {
		client.Chat = ch
		ch.Connect <- client
	} else if !ok {
		newChat := client.Chat.New(client.ChatName)
		go newChat.Run()
		h.ClientsChat[client.ChatName] = newChat
		newChat.Connect <- client
		log.Errorf("HUB Create new chat room %s", client.ChatName)
	}
}

func (h *Hub) unregisterClient(client *ChatUser) {
	h.ClientsChat[client.ChatName].Disconnect <- client
}

func (h *Hub) broadcastToClients(message []byte) {
	for client, chat := range h.ClientsChat {
		select {
		case chat.Broadcast <- message:
		default:
			close(chat.Broadcast)
			delete(h.ClientsChat, client)
		}
	}
}
