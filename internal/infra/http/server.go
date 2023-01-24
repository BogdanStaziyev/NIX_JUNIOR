package http

import (
	"github.com/BogdanStaziyev/NIX_Junior/internal/domain"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo *echo.Echo
	Hub  *domain.Hub
}

func NewServer() *Server {
	s := &Server{
		Echo: echo.New(),
		Hub:  domain.NewHub(),
	}
	go s.Hub.Run()
	return s
}

func (s Server) Start() error {
	return s.Echo.Start(":8080")
}
