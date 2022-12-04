package api_auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/context"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"github.com/miniyus/go-fiber/internal/core/auth"
	"go.uber.org/zap"
	"log"
)

type Handler interface {
	SignUp(ctx *fiber.Ctx) error
	SignIn(ctx *fiber.Ctx) error
	Me(ctx *fiber.Ctx) error
	//Revoke(ctx *fiber.Ctx) error
}

type HandlerStruct struct {
	service Service
}

func NewHandler(service Service) *HandlerStruct {
	return &HandlerStruct{service: service}
}

func validateSignUp(ctx *fiber.Ctx, signUp *SignUp) (bool, *api_error.ErrorResponse) {
	if err := api_error.Validate(signUp); err != nil {
		errRes := api_error.NewValidationError(ctx)
		errRes.FailedFields = err

		return false, &errRes
	}

	if signUp.Password != signUp.PasswordConfirm {
		errRes := api_error.NewValidationError(ctx)
		errRes.FailedFields = map[string]string{
			"PasswordConfirm": "패스워드와 패스워드 확인 필드가 같지 않습니다.",
		}

		return false, &errRes
	}

	return true, nil
}

func (h *HandlerStruct) SignUp(ctx *fiber.Ctx) error {
	signUp := &SignUp{}
	err := ctx.BodyParser(signUp)
	if err != nil {
		errRes := api_error.NewValidationError(ctx)
		return errRes.Response()
	}

	if isValid, errRes := validateSignUp(ctx, signUp); !isValid {
		return errRes.Response()
	}

	result, err := h.service.SignUp(signUp)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(result)
}

func validateSignIn(ctx *fiber.Ctx, in *SignIn) (bool, *api_error.ErrorResponse) {
	if err := api_error.Validate(in); err != nil {
		errRes := api_error.NewValidationError(ctx)
		errRes.FailedFields = err

		return false, &errRes
	}

	return true, nil
}

func (h *HandlerStruct) SignIn(ctx *fiber.Ctx) error {
	signIn := &SignIn{}
	err := ctx.BodyParser(signIn)
	if err != nil {
		errRes := api_error.NewValidationError(ctx)

		return errRes.Response()
	}

	if isValid, errRes := validateSignIn(ctx, signIn); !isValid {
		return errRes.Response()
	}

	result, err := h.service.SignIn(signIn)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"token":      result.Token,
		"expires_at": result.ExpiresAt,
	})
}

func (h *HandlerStruct) Me(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals(context.AuthUser).(*auth.User)
	if !ok {
		logger, ok := ctx.Locals(context.Logger).(*zap.SugaredLogger)
		log.Print(logger)
		if !ok {
			return fiber.NewError(500, "...")
		}
		logger.Debug(user)
	}

	return ctx.JSON(user)
}

//func (h *HandlerStruct) Revoke(ctx *fiber.Ctx) error {
//
//}
