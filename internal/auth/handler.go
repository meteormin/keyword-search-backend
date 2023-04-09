package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/miniyus/gofiber/apierrors"
	"github.com/miniyus/gofiber/log"
	"github.com/miniyus/gofiber/utils"
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
}

func NewHandler(service Service) Handler {
	return &HandlerStruct{
		service: service,
	}
}

func validateSignUp(ctx *fiber.Ctx, signUp *SignUp) (bool, *apierrors.ValidationErrorResponse) {
	errRes := utils.HandleValidate(ctx, signUp)
	if errRes != nil {
		return false, errRes
	}

	return true, nil
}

// SignUp
// @Summary Sign up
// @Description sign up
// @Tags Auth
// @Success 201 {object} SignUpResponse
// @Failure 400 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Param request body SignUp true "sign up body"
// @Router /api/auth/register [post]
func (h *HandlerStruct) SignUp(ctx *fiber.Ctx) error {
	signUp := &SignUp{}
	err := ctx.BodyParser(signUp)
	if err != nil {
		return fiber.ErrBadRequest
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

func validateSignIn(ctx *fiber.Ctx, in *SignIn) (bool, *apierrors.ValidationErrorResponse) {
	errRes := utils.HandleValidate(ctx, in)
	if errRes != nil {
		return false, errRes
	}

	return true, nil
}

// SignIn
// @Summary login
// @Description login
// @Tags Auth
// @Success 200 {object} TokenInfo
// @Failure 400 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Param request body SignIn true "login  body"
// @Router /api/auth/token [post]
func (h *HandlerStruct) SignIn(ctx *fiber.Ctx) error {
	signIn := &SignIn{}
	err := ctx.BodyParser(signIn)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if isValid, errRes := validateSignIn(ctx, signIn); !isValid {
		return errRes.Response()
	}

	result, err := h.service.SignIn(signIn)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(TokenInfo{
		Token:     result.Token,
		ExpiresAt: result.ExpiresAt,
		ExpiresIn: h.service.GetTokenExp(),
	})
}

// Me
// @Summary get my info
// @description get login user info
// @Tags Auth
// @Success 200 {object} auth.User
// @Accept json
// @Produce json
// @Router /api/auth/me [get]
// @Security BearerAuth
func (h *HandlerStruct) Me(ctx *fiber.Ctx) error {
	user, err := GetAuthUser(ctx)
	if err != nil {
		log.GetLogger().Error(user)
		return err
	}

	return ctx.JSON(user)
}

// ResetPassword
// @Summary reset password
// @description reset login user's password
// @Tags Auth
// @Param request body ResetPasswordStruct true "reset password body"
// @Success 200 {object} auth.User
// @Failure 400 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/auth/password [patch]
// @Security BearerAuth
func (h *HandlerStruct) ResetPassword(ctx *fiber.Ctx) error {
	user, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	dto := &ResetPasswordStruct{}

	err = ctx.BodyParser(dto)
	if err != nil {
		return fiber.ErrBadRequest
	}

	errRes := utils.HandleValidate(ctx, dto)
	if errRes != nil {
		return errRes.Response()
	}

	rs, err := h.service.ResetPassword(user.Id, dto)
	if err != nil {
		errRes := apierrors.NewFromError(ctx, err)
		return errRes.Response()
	}

	return ctx.JSON(rs)
}

// RevokeToken
// @Summary revoke token
// @description revoke current jwt token
// @Tags Auth
// @Success 200 {object} utils.StatusResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/auth/revoke [delete]
// @Security BearerAuth
func (h *HandlerStruct) RevokeToken(ctx *fiber.Ctx) error {
	user, err := GetAuthUser(ctx)
	if err != nil {
		return err
	}

	tokenInfo, ok := ctx.Locals("user").(*jwt.Token)
	if !ok {
		return fiber.ErrForbidden
	}

	token := tokenInfo.Raw

	rs, err := h.service.RevokeToken(user.Id, token)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusNoContent).JSON(utils.StatusResponse{
		Status: rs,
	})
}
