package container

import (
	"fmt"
	"github.com/BogdanStaziyev/NIX_Junior/config"
	"github.com/BogdanStaziyev/NIX_Junior/internal/app"
	"github.com/BogdanStaziyev/NIX_Junior/internal/infra/database"
	"github.com/BogdanStaziyev/NIX_Junior/internal/infra/http/handlers"
	"github.com/BogdanStaziyev/NIX_Junior/middleware"
	"github.com/go-redis/redis/v7"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Container struct {
	Services
	Handlers
	Middleware
}

type Services struct {
	app.UserService
	app.AuthService
}

type Handlers struct {
	handlers.RegisterHandler
}

type Middleware struct {
	middleware.AuthMiddleware
}

func New(conf config.Configuration) Container {
	sess := getDbSess(conf)
	newRedis := getRedis(conf)

	userRepository := database.NewUSerRepo(sess)
	passwordGenerator := app.NewGeneratePasswordHash(bcrypt.DefaultCost)
	userService := app.NewUserService(userRepository, passwordGenerator)
	authService := app.NewAuthService(userService, conf, newRedis)
	registerController := handlers.NewRegisterHandler(authService)

	authMiddleware := middleware.NewMiddleware(authService, newRedis)

	return Container{
		Services: Services{
			userService,
			authService,
		},
		Handlers: Handlers{
			registerController,
		},
		Middleware: Middleware{
			authMiddleware,
		},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	sess, err := mysql.Open(
		mysql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}

func getRedis(conf config.Configuration) *redis.Client {
	addr := fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort)
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
