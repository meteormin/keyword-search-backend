package core

import (
	"github.com/gofiber/fiber/v2"
	configure "github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/core/container"
	"github.com/miniyus/keyword-search-backend/internal/core/database"
	"github.com/miniyus/keyword-search-backend/internal/core/register"
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
