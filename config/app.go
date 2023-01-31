package config

import (
	"github.com/gofiber/fiber/v2"
	"os"
)

func app() fiber.Config {
	prefork := false
	appEnv := os.Getenv("APP_ENV")
	if appEnv == string(PRD) {
		prefork = true
	}

	return fiber.Config{
		AppName: os.Getenv("APP_NAME"),
		Prefork: prefork,
	}
}
