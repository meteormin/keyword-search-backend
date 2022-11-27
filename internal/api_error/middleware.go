package api_error

import (
	"github.com/gofiber/fiber/v2"
)

func OverrideDefaultErrorHandler(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	errRes := NewFromError(err)

	return errRes.Response(ctx)
}

func ErrorHandler(ctx *fiber.Ctx) error {
	err := ctx.Next()

	if err == nil {
		return nil
	}

	errRes := NewFromError(err)

	return errRes.Response(ctx)
}
