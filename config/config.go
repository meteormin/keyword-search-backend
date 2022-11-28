package config

import (
	"github.com/gofiber/fiber/v2"
	loggerMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
	"strconv"
)

type Configs struct {
	AppEnv   string
	AppPort  int
	Locale   string
	App      fiber.Config
	Logger   loggerMiddleware.Config
	Database DB
	Path     Path
	Auth     Auth
}

func GetConfigs() *Configs {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))

	if err != nil {
		log.Printf("App Port is not numeric... %v", err)
		port = 8000
	}

	return &Configs{
		AppEnv:   os.Getenv("APP_ENV"),
		AppPort:  port,
		Locale:   os.Getenv("LOCALE"),
		App:      app(),
		Logger:   logger(),
		Database: database(),
		Path:     getPath(),
		Auth:     auth(),
	}
}

func InjectConfigContext(ctx *fiber.Ctx) error {
	ctx.Locals(Config, GetConfigs())

	return ctx.Next()
}
