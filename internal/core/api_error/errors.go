package api_error

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type ErrorInterface interface {
	Response()
}

type ErrorResponse struct {
	ctx          *fiber.Ctx
	Status       string            `json:"status"`
	Code         int               `json:"code"`
	Message      string            `json:"message"`
	FailedFields map[string]string `json:"failed_fields,omitempty"`
}

func NewFromError(ctx *fiber.Ctx, err error) *ErrorResponse {
	if err == nil {
		return nil
	}

	var errRes *ErrorResponse

	if vErr, ok := err.(*fiber.Error); ok {
		errRes = &ErrorResponse{ctx: ctx, Status: "error", Code: vErr.Code, Message: vErr.Message}
	} else if vErr, ok := err.(error); ok {
		errRes = &ErrorResponse{ctx: ctx, Status: "error", Code: fiber.StatusInternalServerError, Message: vErr.Error()}
	} else {
		errRes = &ErrorResponse{ctx: ctx, Status: "error", Code: fiber.StatusInternalServerError, Message: "Unknown Error"}
	}

	return errRes
}

func NewErrorResponse(ctx *fiber.Ctx, code int, message string) ErrorResponse {
	return ErrorResponse{
		ctx:     ctx,
		Status:  "error",
		Code:    code,
		Message: message,
	}
}

func NewValidationError(ctx *fiber.Ctx) ErrorResponse {
	return ErrorResponse{
		ctx:     ctx,
		Status:  "error",
		Code:    fiber.StatusBadRequest,
		Message: utils.StatusMessage(fiber.StatusBadRequest),
	}
}

func (er *ErrorResponse) Response() error {
	if er.Code == 0 {
		er.Code = fiber.StatusInternalServerError
	}

	if er.Message == "" {
		er.Message = utils.StatusMessage(er.Code)
	}

	return er.ctx.Status(er.Code).JSON(er)
}
