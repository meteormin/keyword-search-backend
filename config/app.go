package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/app"
	"log"
	"os"
	"strconv"
)

func appConfig() app.Config {
	prefork := false
	appEnv := os.Getenv("APP_ENV")
	if appEnv == string(app.PRD) {
		prefork = true
	}

	appPortStr := os.Getenv("APP_PORT")
	appPort, err := strconv.Atoi(appPortStr)
	if err != nil {
		log.Printf("App Port is not numeric... %v", err)
		log.Print("App Port set 8000...")
		appPort = 8000
	}

	return app.Config{
		Env:      app.Env(appEnv),
		Port:     appPort,
		Locale:   os.Getenv("LOCALE"),
		TimeZone: os.Getenv("TIME_ZONE"),
		FiberConfig: fiber.Config{
			AppName: os.Getenv("APP_NAME"),
			Prefork: prefork,
		},
	}
}
