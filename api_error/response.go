package api_error

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	fUtils "github.com/gofiber/fiber/v2/utils"
	"gorm.io/gorm"
)

type ErrorInterface interface {
	Response() error
}

type ErrorResponse struct {
	ctx     *fiber.Ctx
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	ErrorResponse
	FailedFields map[string]string `json:"failed_fields,omitempty"`
}

func NewFromError(ctx *fiber.Ctx, err error) ErrorInterface {
	if err == nil {
		return nil
	}

	var errRes ErrorInterface

	if vErr, ok := err.(*fiber.Error); ok {
		errRes = NewErrorResponse(ctx, vErr.Code, vErr.Message)
	} else if vErr, ok := err.(error); ok {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errRes = NewErrorResponse(ctx, fiber.StatusNotFound, vErr.Error())
		} else {
			errRes = NewErrorResponse(ctx, fiber.StatusInternalServerError, vErr.Error())
		}

	} else {

		errRes = NewErrorResponse(ctx, fiber.StatusInternalServerError, vErr.Error())
	}

	return errRes
}

func NewErrorResponse(ctx *fiber.Ctx, code int, message string) *ErrorResponse {
	return &ErrorResponse{
		ctx:     ctx,
		Status:  "error",
		Code:    code,
		Message: message,
	}
}

func NewValidationErrorResponse(ctx *fiber.Ctx, failedFields map[string]string) *ValidationErrorResponse {
	code := fiber.StatusBadRequest
	return &ValidationErrorResponse{
		ErrorResponse: ErrorResponse{ctx: ctx, Status: "error", Code: code, Message: fUtils.StatusMessage(code)},
		FailedFields:  failedFields,
	}
}

func (er *ErrorResponse) Response() error {
	if er.Code == 0 {
		er.Code = fiber.StatusInternalServerError
	}

	if er.Message == "" {
		er.Message = fUtils.StatusMessage(er.Code)
	}

	return er.ctx.Status(er.Code).JSON(er)
}
