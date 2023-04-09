package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/gofiber/apierrors"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/internal/auth"
	"strconv"
)

type Handler interface {
	All(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	PatchForAdm(ctx *fiber.Ctx) error
	PatchForMe(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type HandlerStruct struct {
	service Service
}

func NewHandler(service Service) *HandlerStruct {
	return &HandlerStruct{service: service}
}

func (h *HandlerStruct) All(ctx *fiber.Ctx) error {
	entities, err := h.service.All()

	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"data": entities,
	})
}

func (h *HandlerStruct) Get(ctx *fiber.Ctx) error {
	prams := ctx.AllParams()
	userId, err := strconv.ParseUint(prams["id"], 10, 64)
	if err != nil {
		return err
	}

	user, err := h.service.Get(uint(userId))

	if err != nil {
		return err
	}

	return ctx.JSON(user)
}

func (h *HandlerStruct) Create(ctx *fiber.Ctx) error {
	var dto CreateUser

	err := utils.HandleValidate(ctx, &dto)
	if err != nil {
		return err.Response()
	}

	create, err2 := h.service.Create(dto)
	if err2 != nil {
		return err2
	}

	return ctx.JSON(create)
}

func (h *HandlerStruct) PatchForAdm(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return apierrors.NewValidationErrorResponse(ctx, map[string]string{
			"id": "id는 정수형이여야 합니다.",
		}).Response()
	}

	cu, err := auth.GetAuthUser(ctx)
	if err != nil {
		return err
	}

	if cu.Role != string(entity.Admin) {
		return fiber.ErrForbidden
	}

	var dto PatchUser
	validateError := utils.HandleValidate(ctx, &dto)
	if err != nil {
		return validateError.Response()
	}

	update, err := h.service.Update(uint(pk), dto)
	if err != nil {
		return err
	}

	return ctx.JSON(update)
}

func (h *HandlerStruct) PatchForMe(ctx *fiber.Ctx) error {
	cu, err := auth.GetAuthUser(ctx)
	if err != nil {
		return err
	}

	pk := cu.Id

	var dto PatchUser
	validateError := utils.HandleValidate(ctx, &dto)
	if err != nil {
		return validateError.Response()
	}

	update, err := h.service.Update(pk, dto)
	if err != nil {
		return err
	}

	return ctx.JSON(update)
}

func (h *HandlerStruct) Delete(ctx *fiber.Ctx) error {
	cu, err := auth.GetAuthUser(ctx)
	if err != nil {
		return err
	}

	pk := cu.Id

	rs, err := h.service.Delete(pk)
	if err != nil {
		return err
	}

	return ctx.JSON(utils.StatusResponse{
		Status: rs,
	})
}
