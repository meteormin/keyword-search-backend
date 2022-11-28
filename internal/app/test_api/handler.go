package test_api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	configure "github.com/miniyus/go-fiber/config"
)

type Handler interface {
	GetTest(ctx *fiber.Ctx) error
}

type HandlerStruct struct {
	service Service
}

func NewHandler(service Service) HandlerStruct {
	return HandlerStruct{service: service}
}

func (h HandlerStruct) GetTest(ctx *fiber.Ctx) error {
	var data *configure.Configs

	config := ctx.Locals(configure.Config)
	if config != nil {
		return errors.New("JUST")
	}

	if config, ok := config.(*configure.Configs); ok {
		data = config
	}

	return ctx.JSON(fiber.Map{
		"data": data,
	})
}
