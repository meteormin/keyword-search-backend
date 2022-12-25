package core

import (
	"github.com/gofiber/fiber/v2"
	configure "github.com/miniyus/go-fiber/config"
	"github.com/miniyus/go-fiber/internal/core/container"
	"github.com/miniyus/go-fiber/internal/core/database"
	"github.com/miniyus/go-fiber/internal/core/register"
)

func New() container.Container {
	config := configure.GetConfigs()

	wrapper := container.New(
		fiber.New(config.App),
		database.DB(config.Database),
		config,
	)

	register.Resister(wrapper)

	return wrapper
}
