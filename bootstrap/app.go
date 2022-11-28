package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	configure "github.com/miniyus/go-fiber/config"
	"github.com/miniyus/go-fiber/core/container"
	"github.com/miniyus/go-fiber/core/register"
)

func Run() {
	config := configure.GetConfigs()

	wrapper := container.NewContainer(fiber.New(config.App), config)

	register.Resister(wrapper)

	wrapper.Run()
}
