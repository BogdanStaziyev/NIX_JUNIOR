package handlers

import (
	"github.com/BogdanStaziyev/NIX_Junior/internal/app"
	"github.com/BogdanStaziyev/NIX_Junior/internal/domain"
	s "github.com/BogdanStaziyev/NIX_Junior/internal/infra/http"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
)

type WebsocketConn struct {
	server        *s.Server
	clientService app.ClientService
}

func NewWebsocketConn(s *s.Server, c app.ClientService) WebsocketConn {
	return WebsocketConn{
		server:        s,
		clientService: c,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (cli *WebsocketConn) Socket(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	chatName := c.QueryParam("chat")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	user := token.Claims.(*app.JwtTokenClaim)

	client := domain.NewClient(conn, cli.server.Hub)

	client.ID = user.ID
	client.Name = user.Name
	client.ChatName = chatName

	client.Hub.RegisterChat <- client

	go cli.clientService.WritePump(client)
	go cli.clientService.ReadPump(client)
	return err
}
