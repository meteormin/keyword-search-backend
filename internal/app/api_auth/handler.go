package api_auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"github.com/miniyus/go-fiber/internal/core/auth"
	"github.com/miniyus/go-fiber/internal/core/context"
	"github.com/miniyus/go-fiber/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

type Handler interface {
	SignUp(ctx *fiber.Ctx) error
	SignIn(ctx *fiber.Ctx) error
	Me(ctx *fiber.Ctx) error
	ResetPassword(ctx *fiber.Ctx) error
	RevokeToken(ctx *fiber.Ctx) error
}

type HandlerStruct struct {
	service Service
	logger  *zap.SugaredLogger
}

func NewHandler(service Service, logger *zap.SugaredLogger) *HandlerStruct {
	return &HandlerStruct{service: service, logger: logger}
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
			"PasswordConfirm": "패스워드와 패스워드확인 필드가 일치하지않습니다.",
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
		"expires_at": result.ExpiresAt.Format("2006-01-02 15:04:05"),
	})
}

func (h *HandlerStruct) Me(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals(context.AuthUser).(*auth.User)
	if !ok {
		logger, ok := ctx.Locals(context.Logger).(*zap.SugaredLogger)

		if !ok {
			return fiber.NewError(500, "Can't Load Context Logger")
		}
		logger.Error(user)
		return fiber.NewError(500, "Can't Load Context AuthUser")
	}

	return ctx.JSON(user)
}

func (h *HandlerStruct) ResetPassword(ctx *fiber.Ctx) error {

	user, err := utils.GetAuthUser(ctx)
	if err != nil {
		h.logger.Error(err)
		errRes := api_error.NewFromError(ctx, err)
		return errRes.Response()
	}

	dto := &ResetPasswordStruct{}
	h.logger.Debug("parse body...")
	err = ctx.BodyParser(dto)
	if err != nil {
		h.logger.Error(err)
		errRes := api_error.NewValidationError(ctx)
		return errRes.Response()
	}

	failedFields := api_error.Validate(dto)
	if failedFields != nil {
		h.logger.Error(err)
		errRes := api_error.NewValidationError(ctx)
		errRes.FailedFields = failedFields
		return errRes.Response()
	}

	if dto.Password != dto.PasswordConfirm {
		errRes := api_error.NewValidationError(ctx)
		errRes.FailedFields = map[string]string{
			"PasswordConfirm": "패스워드와 패스워드확인 필드가 일치하지않습니다.",
		}
		return errRes.Response()
	}

	rs, err := h.service.ResetPassword(user.Id, dto)
	if err != nil {
		h.logger.Error(err)
		errRes := api_error.NewFromError(ctx, err)
		return errRes.Response()
	}

	return ctx.JSON(rs)
}

func (h *HandlerStruct) RevokeToken(ctx *fiber.Ctx) error {
	user, err := utils.GetAuthUser(ctx)
	if err != nil {
		return err
	}

	tokenInfo, ok := ctx.Locals("user").(*jwt.Token)
	if !ok {
		errRes := api_error.NewErrorResponse(ctx, fiber.StatusForbidden, http.StatusText(fiber.StatusForbidden))
		return errRes.Response()
	}

	token := tokenInfo.Raw

	rs, err := h.service.RevokeToken(user.Id, token)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"status": rs,
	})
}
