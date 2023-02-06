package utils

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type ContextKey string

// AuthUserKey context constants
const AuthUserKey ContextKey = "authUser"

func AddContext(localsKey ContextKey, value interface{}) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals(localsKey, value)

		return ctx.Next()
	}
}

func GetContext[T interface{}](ctx *fiber.Ctx, localsKey ContextKey) (T, error) {
	getCtx, ok := ctx.Locals(localsKey).(T)
	if !ok {
		return getCtx, errors.New(fmt.Sprintf("Can not get context: %s", localsKey))
	}

	return getCtx, nil
}
