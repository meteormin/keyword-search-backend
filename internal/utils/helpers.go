package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/config"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"github.com/miniyus/go-fiber/internal/core/auth"
	"net/http"
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
	user, ok := c.Locals(config.AuthUser).(auth.User)
	if !ok {
		status := fiber.StatusUnauthorized
		errRes := api_error.NewErrorResponse(c, status, http.StatusText(status))
		return nil, errRes.Response()
	}

	return &user, nil
}
