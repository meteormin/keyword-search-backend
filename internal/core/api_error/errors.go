package api_error

import (
	"github.com/gofiber/fiber/v2"
	fUtils "github.com/gofiber/fiber/v2/utils"
)

type errorMaker func(ctx *fiber.Ctx) ErrorInterface

func newErrorMaker(code int) errorMaker {
	return func(ctx *fiber.Ctx) ErrorInterface {
		return NewErrorResponse(ctx, code, fUtils.StatusMessage(code))
	}
}

func NewForbiddenError(ctx *fiber.Ctx) ErrorInterface {
	code := fiber.StatusForbidden
	return newErrorMaker(code)(ctx)
}

func NewBadRequestError(ctx *fiber.Ctx) ErrorInterface {
	code := fiber.StatusBadRequest
	return newErrorMaker(code)(ctx)
}

func NewServerError(ctx *fiber.Ctx) ErrorInterface {
	code := fiber.StatusInternalServerError
	return newErrorMaker(code)(ctx)
}
