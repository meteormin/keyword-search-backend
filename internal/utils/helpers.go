package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/api_error"
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
