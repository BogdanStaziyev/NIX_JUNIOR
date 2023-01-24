package main

import (
	"github.com/BogdanStaziyev/NIX_Junior/config"
	"github.com/BogdanStaziyev/NIX_Junior/config/container"
	"github.com/BogdanStaziyev/NIX_Junior/internal/infra/database"
	"github.com/BogdanStaziyev/NIX_Junior/internal/infra/http"
	"github.com/BogdanStaziyev/NIX_Junior/internal/infra/http/router"
	"log"
)

// @title 		NIX TRAINEE PROGRAM Demo App
// @version 	V1.echo
// @description REST service for NIX TRAINEE PROGRAM

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host 		localhost:8080
// @BasePath 	/
func main() {
	var conf = config.GetConfiguration()

	err := database.Migrate(conf)
	if err != nil {
		log.Fatalf("Unable to apply migrations: %q\n", err)
	}

	// Echo Server
	srv := http.NewServer()

	cont := container.New(conf, *srv)

	router.EchoRouter(srv, cont)

	err = srv.Start()
	if err != nil {
		log.Fatal("Port already used")
	}
}
