package test_api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	configure "github.com/miniyus/go-fiber/config"
)

type Handler interface {
	GetTest(ctx *fiber.Ctx) error
}

type HandlerImpl struct {
	service Service
}

func NewHandler(service Service) HandlerImpl {
	return HandlerImpl{service: service}
}

func (h HandlerImpl) GetTest(ctx *fiber.Ctx) error {
	var data map[string]any

	config := ctx.Locals("configs")
	if config != nil {
		return errors.New("JUST")
	}

	if config, ok := config.(*configure.Configs); ok {
		data = config.Test
	}

	return ctx.JSON(fiber.Map{
		"data": data,
	})
}
