package http

import (
	"github.com/BogdanStaziyev/NIX_Junior/config/container"
	"github.com/BogdanStaziyev/NIX_Junior/internal/infra/http/validators"
	MW "github.com/labstack/echo/v4/middleware"
)

func EchoRouter(s *Server, cont container.Container) {

	e := s.Echo
	e.Use(MW.Logger())
	e.Validator = validators.NewValidator()

	e.POST("/register", cont.RegisterHandler.Register)
	e.POST("/login", cont.RegisterHandler.Login)

	v1 := e.Group("/api/v1")
	v1.GET("", PingHandler)
}
