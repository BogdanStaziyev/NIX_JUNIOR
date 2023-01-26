package app

import (
	"encoding/json"
	"github.com/BogdanStaziyev/NIX_Junior/internal/domain"
	"log"
)

type EventService interface {
	SendMessageToAll(c *domain.ChatUser, message domain.SendMessageToAll) error
	SendMessageToOne(c *domain.ChatUser, message domain.SendMessageToOne) error
}

type eventService struct {
}

func NewEventService() EventService {
	return eventService{}
}

func (e eventService) SendMessageToAll(c *domain.ChatUser, message domain.SendMessageToAll) error {
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
	for ind := range c.Chat.Users {
		if ind.ID == message.UserID {
			ind.Send <- byt
			c.Send <- byt
			return nil
		}
	}
	c.Send <- []byte("client does not exist")
	return nil
}
