package api_error

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type ErrorInterface interface {
	Response()
}

type ErrorResponse struct {
	Status       string
	Code         int
	Message      string
	FailedFields map[string]string
}

func NewFromError(err error) *ErrorResponse {
	if err == nil {
		return nil
	}

	var errRes *ErrorResponse

	if vErr, ok := err.(*fiber.Error); ok {
		errRes = &ErrorResponse{Status: "error", Code: vErr.Code, Message: vErr.Message}
	} else if vErr, ok := err.(error); ok {
		errRes = &ErrorResponse{Status: "error", Code: fiber.StatusInternalServerError, Message: vErr.Error()}
	} else {
		errRes = &ErrorResponse{Status: "error", Code: fiber.StatusInternalServerError, Message: "Unknown Error"}
	}

	return errRes
}

func (er *ErrorResponse) Response(ctx *fiber.Ctx) error {
	if er.Code == 0 {
		er.Code = fiber.StatusInternalServerError
	}

	if er.Message == "" {
		er.Message = http.StatusText(er.Code)
	}

	if er.Code == fiber.StatusBadRequest && er.FailedFields != nil {
		return ctx.Status(er.Code).JSON(fiber.Map{
			"status":        er.Status,
			"code":          er.Code,
			"message":       er.Message,
			"failed_fields": er.FailedFields,
		})
	}

	return ctx.Status(er.Code).JSON(fiber.Map{
		"status":  er.Status,
		"code":    er.Code,
		"message": er.Message,
	})
}
