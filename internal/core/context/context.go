package context

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Key string

// context constants
// ctx.Locals() 메서드에서 주로 사용됨
const (
	Container    Key = "container"
	App          Key = "app"
	DB           Key = "db"
	Config       Key = "config"
	Logger       Key = "logger"
	AuthUser     Key = "authUser"
	JwtGenerator Key = "jwtGenerator"
	Permissions  Key = "permissions"
	Redis        Key = "redis"
)

func AddContext(localsKey Key, value interface{}) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals(localsKey, value)

		return ctx.Next()
	}
}

func GetContext[T interface{}](ctx *fiber.Ctx, localsKey Key) (T, error) {
	getCtx, ok := ctx.Locals(localsKey).(T)
	if !ok {
		return getCtx, errors.New(fmt.Sprintf("Can not get context: %s", localsKey))
	}

	return getCtx, nil
}
