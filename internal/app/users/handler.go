package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"strconv"
)

type Handler interface {
	All(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	ResetPassword(ctx *fiber.Ctx) error
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
	entities, err := h.service.All()

	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"data": entities,
	})
}

func (h *HandlerStruct) Get(ctx *fiber.Ctx) error {
	prams := ctx.AllParams()
	userId, err := strconv.ParseUint(prams["id"], 0, 64)
	if err != nil {
		return err
	}

	user, err := h.service.Get(uint(userId))

	if err != nil {
		return err
	}

	return ctx.JSON(user)
}

func (h *HandlerStruct) Create(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (h *HandlerStruct) Update(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (h *HandlerStruct) ResetPassword(ctx *fiber.Ctx) error {
	password := &ResetPasswordStruct{}
	params := ctx.AllParams()
	err := ctx.BodyParser(password)
	if err != nil {
		errRes := api_error.NewValidationError(ctx)
		return errRes.Response()
	}

	failedFields := api_error.Validate(password)

	if failedFields != nil {
		errRes := api_error.NewValidationError(ctx)
		errRes.FailedFields = failedFields

		return errRes.Response()
	}
	id, err := strconv.ParseUint(params["id"], 0, 64)
	if err != nil {
		return err
	}

	rs, err := h.service.ResetPassword(uint(id), *password)
	if err != nil {
		return err
	}
	return ctx.JSON(rs)
}

func (h *HandlerStruct) Patch(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

func (h *HandlerStruct) Delete(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}
