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
	ctx          *fiber.Ctx
	Status       string            `json:"status"`
	Code         int               `json:"code"`
	Message      string            `json:"message"`
	FailedFields map[string]string `json:"failed_fields,omitempty"`
}

func NewFromError(ctx *fiber.Ctx, err error) ErrorInterface {
	if err == nil {
		return nil
	}

	var errRes *ErrorResponse

	if vErr, ok := err.(*fiber.Error); ok {
		errRes = &ErrorResponse{ctx: ctx, Status: "error", Code: vErr.Code, Message: vErr.Message}
	} else if vErr, ok := err.(error); ok {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errRes = &ErrorResponse{ctx: ctx, Status: "error", Code: fiber.StatusNotFound, Message: vErr.Error()}
		}

		errRes = &ErrorResponse{ctx: ctx, Status: "error", Code: fiber.StatusInternalServerError, Message: vErr.Error()}
	} else {
		errRes = &ErrorResponse{ctx: ctx, Status: "error", Code: fiber.StatusInternalServerError, Message: "Unknown Error"}
	}

	return errRes
}

func NewErrorResponse(ctx *fiber.Ctx, code int, message string) ErrorInterface {
	return &ErrorResponse{
		ctx:     ctx,
		Status:  "error",
		Code:    code,
		Message: message,
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
