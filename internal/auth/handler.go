package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/api_error"
	"net/http"
)

type Handler interface {
	SignUp(ctx *fiber.Ctx) error
	SignIn(ctx *fiber.Ctx) error
	//Revoke(ctx *fiber.Ctx) error
}

type HandlerStruct struct {
	service Service
}

func NewHandler(service Service) *HandlerStruct {
	return &HandlerStruct{service: service}
}

func (h *HandlerStruct) SignUp(ctx *fiber.Ctx) error {
	signUp := &SignUp{}
	err := ctx.BodyParser(signUp)
	if err != nil {
		errRes := api_error.ErrorResponse{
			Code:         fiber.StatusBadRequest,
			Message:      http.StatusText(fiber.StatusBadRequest),
			FailedFields: api_error.Validate(err),
		}
		return errRes.Response(ctx)
	}

	result, err := h.service.SignUp(signUp)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"token":      result.Token,
		"expires_in": result.ExpiresAt,
	})
}

func (h *HandlerStruct) SignIn(ctx *fiber.Ctx) error {
	signIn := &SignIn{}
	err := ctx.BodyParser(signIn)
	if err != nil {
		errRes := api_error.ErrorResponse{
			Code:         fiber.StatusBadRequest,
			Message:      http.StatusText(fiber.StatusBadRequest),
			FailedFields: api_error.Validate(err),
		}

		return errRes.Response(ctx)
	}

	result, err := h.service.SignIn(signIn)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"token":      result.Token,
		"expires_in": result.ExpiresAt,
	})
}

//func (h *HandlerStruct) Revoke(ctx *fiber.Ctx) error {
//
//}
