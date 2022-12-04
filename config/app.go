package config

import (
	"github.com/gofiber/fiber/v2"
	"os"
)

func app() fiber.Config {
	return fiber.Config{
		AppName: os.Getenv("APP_NAME"),
		Prefork: true,
	}
}
