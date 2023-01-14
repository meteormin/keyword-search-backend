package resolver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/internal/core/container"
	"github.com/miniyus/keyword-search-backend/internal/core/context"
	"github.com/miniyus/keyword-search-backend/internal/core/permission"
	"github.com/miniyus/keyword-search-backend/pkg/jwt"
	"go.uber.org/zap"
)

func AddContext(localsKey context.Key, value interface{}) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals(localsKey, value)

		return ctx.Next()
	}
}

func GetContext(ctx *fiber.Ctx, localsKey context.Key) interface{} {
	instance := ctx.Locals(localsKey)
	if instance != nil {
		return instance
	}

	wrapper, ok := ctx.Locals(context.Container).(container.Container)
	if !ok {
		statusCode := fiber.StatusInternalServerError
		return fiber.NewError(statusCode, "Failed Get Container in Ctx")
	}

	return wrapper.Get(localsKey)
}

func ResolveContext(ctx *fiber.Ctx, dest interface{}) error {
	wrapper, ok := ctx.Locals(context.Container).(container.Container)
	if !ok {
		statusCode := fiber.StatusInternalServerError
		return fiber.NewError(statusCode, "Failed Get Container in Ctx")
	}

	result := wrapper.Resolve(dest)
	if result == nil {
		statusCode := fiber.StatusInternalServerError
		return fiber.NewError(statusCode, "Failed Resolve...")
	}

	return nil
}

func TokenGenerator(c container.Container) jwt.Generator {
	var tokenGenerator jwt.Generator
	c.Resolve(&tokenGenerator)
	return tokenGenerator
}

func Logger(c container.Container) *zap.SugaredLogger {
	var logger *zap.SugaredLogger
	c.Resolve(&logger)
	return logger
}

func PermissionCollection(c container.Container) permission.Collection {
	var permCollect permission.Collection
	c.Resolve(&permCollect)
	return permCollect
}