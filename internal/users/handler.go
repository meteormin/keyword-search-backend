package users

import "github.com/gofiber/fiber/v2"

type Handler interface {
	All(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Patch(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type HandlerStruct struct {
	service Service
}

func NewHandler(service Service) *HandlerStruct {
	return &HandlerStruct{service: service}
}

func (h *HandlerStruct) All(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (h *HandlerStruct) Get(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (h *HandlerStruct) Create(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (h *HandlerStruct) Update(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (h *HandlerStruct) Patch(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (h *HandlerStruct) Delete(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}
