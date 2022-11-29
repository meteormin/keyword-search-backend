package api_error

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type ErrorInterface interface {
	Response()
}

type ErrorResponse struct {
	ctx          *fiber.Ctx
	Status       string            `json:"status"`
	Code         int               `json:"code"`
	Message      string            `json:"message"`
	FailedFields map[string]string `json:"failed_fields"`
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

func NewValidationError(ctx *fiber.Ctx) ErrorResponse {
	return ErrorResponse{
		ctx:     ctx,
		Status:  "error",
		Code:    fiber.StatusBadRequest,
		Message: http.StatusText(fiber.StatusBadRequest),
	}
}

func (er *ErrorResponse) Response() error {
	if er.Code == 0 {
		er.Code = fiber.StatusInternalServerError
	}

	if er.Message == "" {
		er.Message = http.StatusText(er.Code)
	}

	if er.Code == fiber.StatusBadRequest && er.FailedFields != nil {
		return er.ctx.Status(er.Code).JSON(er)
	}

	return er.ctx.Status(er.Code).JSON(er)
}
