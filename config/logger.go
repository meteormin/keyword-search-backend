package config

import (
	loggerMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	"os"
)

func logger() loggerMiddleware.Config {

	return loggerMiddleware.Config{
		Format:     "[${time}] ${ip}:${port} | (${pid}) ${status} - ${method} ${path}\n",
		TimeZone:   os.Getenv("TIME_ZONE"),
		TimeFormat: "2006-01-02 15:04:05",
	}
}
