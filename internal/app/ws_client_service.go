package app

import (
	"encoding/json"
	"errors"
	"github.com/BogdanStaziyev/NIX_Junior/internal/domain"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type ClientService interface {
	ProcessEvents(rawMessage []byte, client *domain.ChatUser) error
	ReadPump(client *domain.ChatUser)
	WritePump(client *domain.ChatUser)
}

type clientService struct {
	e EventService
}

func NewClientService(event EventService) ClientService {
	return clientService{e: event}
}

func (c clientService) ProcessEvents(rawMessage []byte, client *domain.ChatUser) error {
	var baseMessage domain.Base
	err := json.Unmarshal(rawMessage, &baseMessage)
	if err != nil {
		return err
	}

	if baseMessage.Action == "" {
		return errors.New("error deserializing message")
	}

	switch baseMessage.Action {
	case domain.ActionSandMessageToAllUsers:
		var message domain.SendMessageToAll
		err = json.Unmarshal(rawMessage, &message)
		if err != nil {
			log.Println(err)
			return err
		}
		if err = c.e.SendMessageToAllUsers(client, message); err != nil {
			return err
		}
	case domain.ActionSendPrivate:
		var message domain.SendMessageToOne
		err = json.Unmarshal(rawMessage, &message)
		if err != nil {
			log.Println(err)
			return err
		}
		if err = c.e.SendMessageToOne(client, message); err != nil {
			return err
		}
	case domain.ActionSandMessageToChat:
		var message domain.SendMessageToChat
		err = json.Unmarshal(rawMessage, &message)
		if err != nil {
			log.Println(err)
			return err
		}
		if err = c.e.SendMessageInChat(client, message); err != nil {
			log.Println(err)
			return err
		}
	}

	return err
}

// ReadPump pumps events from the websocket connection to the hub.
//
// The application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c clientService) ReadPump(client *domain.ChatUser) {

	defer func() {
		client.Disconnect()
	}()
	client.Conn.SetReadLimit(maxMessageSize)
	if err := client.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	client.Conn.SetPongHandler(func(string) error {
		return client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error 1: %v", err)
			}
			log.Printf("error 2: %v", err)
			break
		}
		err = c.ProcessEvents(message, client)
		if err != nil {
			log.Printf("error 3: %v", err)
		}
	}

}

// WritePump pumps messages from the hub to the websocket connection.
//
// A goroutine running WritePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c clientService) WritePump(client *domain.ChatUser) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.Send:
			//todo ???
			err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if !ok {
				// The hub closed the channel.
				if err := client.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}
			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Println("connection closed error: ", err)
				return
			}
			if err := w.Close(); err != nil {
				log.Println("connection closed error: ", err)
				return
			}
		case <-ticker.C:
			err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if err = client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("write msg: ", err)
				return
			}
		}
	}
}
