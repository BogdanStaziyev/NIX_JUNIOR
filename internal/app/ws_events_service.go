package app

import (
	"encoding/json"
	"github.com/BogdanStaziyev/NIX_Junior/internal/domain"
	"log"
)

type EventService interface {
	SendMessageToAllUsers(c *domain.ChatUser, message domain.SendMessageToAll) error
	SendMessageToOne(c *domain.ChatUser, message domain.SendMessageToOne) error
	SendMessageInChat(c *domain.ChatUser, message domain.SendMessageToChat) error
}

type eventService struct {
}

func NewEventService() EventService {
	return eventService{}
}

func (e eventService) SendMessageToAllUsers(c *domain.ChatUser, message domain.SendMessageToAll) error {
	byt, err := json.Marshal(c.Name + ": " + message.Message)
	if err != nil {
		log.Println(err)
		return err
	}
	c.Hub.Broadcast <- byt
	return nil
}

func (e eventService) SendMessageToOne(c *domain.ChatUser, message domain.SendMessageToOne) error {
	byt, err := json.Marshal(c.Name + ": " + message.Message)
	if err != nil {
		log.Println(err)
		return err
	}
	//todo find in db after register
	for _, chat := range c.Hub.ClientsChat {
		for ch := range chat.Users {
			if ch.ID == message.UserID {
				ch.Send <- byt
				c.Send <- byt
				return nil
			}
		}
	}
	c.Send <- []byte("client does not exist")
	return nil
}

func (e eventService) SendMessageInChat(c *domain.ChatUser, message domain.SendMessageToChat) error {
	byt, err := json.Marshal(c.Name + ": " + message.Message)
	if err != nil {
		log.Println(err)
		return err
	}
	c.Chat.Broadcast <- byt
	return nil
}
