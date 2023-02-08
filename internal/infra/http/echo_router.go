package http

import (
	"github.com/BogdanStaziyev/NIX_Junior/config"
	"github.com/BogdanStaziyev/NIX_Junior/config/container"
	"github.com/BogdanStaziyev/NIX_Junior/internal/infra/http/validators"
	"github.com/labstack/echo/v4"
	MW "github.com/labstack/echo/v4/middleware"
)

func EchoRouter(e *echo.Echo, cont container.Container) {

	e.Use(MW.Logger())
	e.Validator = validators.NewValidator()

	u := e.Group("user")
	u.POST("/register", cont.RegisterHandler.Register)
	u.POST("/login", cont.RegisterHandler.Login)

	v1 := e.Group("/api/v1")
	v1.GET("", PingHandler)

	auth := cont.AuthMiddleware.JWT(config.GetConfiguration().AccessSecret)
	valid := cont.AuthMiddleware.ValidateJWT()
	ws := e.Group("ws")
	ws.Use(auth, valid)
	ws.GET("/", cont.Handlers.Socket)
}
