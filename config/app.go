package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/internal/core/api_error"
	"os"
)

func app() fiber.Config {
	return fiber.Config{
		AppName:      os.Getenv("APP_NAME"),
		Prefork:      true,
		ErrorHandler: api_error.OverrideDefaultErrorHandler,
	}
}
