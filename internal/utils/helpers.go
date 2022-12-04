package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"github.com/miniyus/go-fiber/internal/core/auth"
	"github.com/miniyus/go-fiber/internal/core/context"
	"net/http"
	"time"
)

func HandleValidate(c *fiber.Ctx, data interface{}) error {
	failed := api_error.Validate(data)
	if failed != nil {
		errRes := api_error.NewValidationError(c)
		errRes.FailedFields = failed
		return errRes.Response()
	}

	return nil
}

func GetAuthUser(c *fiber.Ctx) (*auth.User, error) {
	user, ok := c.Locals(context.AuthUser).(auth.User)
	if !ok {
		status := fiber.StatusUnauthorized
		errRes := api_error.NewErrorResponse(c, status, http.StatusText(status))
		return nil, errRes.Response()
	}

	return &user, nil
}

func TimeIn(t time.Time, tz string) time.Time {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}

	return t.In(loc)
}
