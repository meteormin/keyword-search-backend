package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/core/api_error"
	"github.com/miniyus/keyword-search-backend/internal/core/auth"
	"github.com/miniyus/keyword-search-backend/internal/core/container"
	"github.com/miniyus/keyword-search-backend/internal/core/context"
	"net/http"
	"time"
)

const (
	DefaultDateLayout = "2006-01-02 15:04:05"
)

func HandleValidate(c *fiber.Ctx, data interface{}) *api_error.ErrorResponse {
	failed := Validate(data)
	if failed != nil {
		errRes := api_error.NewValidationError(c)
		errRes.FailedFields = failed
		return &errRes
	}

	return nil
}

func GetAuthUser(c *fiber.Ctx) (*auth.User, error) {
	user, ok := c.Locals(context.AuthUser).(*auth.User)
	if !ok {
		status := fiber.StatusUnauthorized
		errRes := api_error.NewErrorResponse(c, status, http.StatusText(status))
		return nil, errRes.Response()
	}

	return user, nil
}

func TimeIn(t time.Time, tz string) time.Time {
	if tz == "" {
		cfg := config.GetConfigs()
		tz = cfg.TimeZone
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}

	return t.In(loc)
}

func AddContext(localsKey context.Key, value interface{}) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals(localsKey, value)

		return ctx.Next()
	}
}

func FindContext(ctx *fiber.Ctx, dest interface{}) error {
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
