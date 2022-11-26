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

type HandlerImpl struct {
	service Service
}

func NewHandler(service Service) *HandlerImpl {
	return &HandlerImpl{service: service}
}

func (h *HandlerImpl) All(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (h *HandlerImpl) Get(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (h *HandlerImpl) Create(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (h *HandlerImpl) Update(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (h *HandlerImpl) Patch(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (h *HandlerImpl) Delete(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}
