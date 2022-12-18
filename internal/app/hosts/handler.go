package hosts

import (
	"github.com/gofiber/fiber/v2"
	fUtils "github.com/gofiber/fiber/v2/utils"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"github.com/miniyus/go-fiber/internal/utils"
	"strconv"
)

type Handler interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	All(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type HandlerStruct struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &HandlerStruct{service: service}
}

func (h *HandlerStruct) Create(c *fiber.Ctx) error {
	user, err := utils.GetAuthUser(c)
	if err != nil {
		return err
	}

	dto := &CreateHost{}
	err = c.BodyParser(dto)
	if err != nil {
		errRes := api_error.NewValidationError(c)
		return errRes.Response()
	}

	errRes := utils.HandleValidate(c, dto)
	if errRes != nil {
		return errRes.Response()
	}

	dto.UserId = user.Id

	result, err := h.service.Create(dto)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *HandlerStruct) Update(c *fiber.Ctx) error {
	dto := &UpdateHost{}
	params := c.AllParams()

	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	err = c.BodyParser(dto)
	if err != nil {
		errRes := api_error.NewValidationError(c)
		return errRes.Response()
	}

	errRes := utils.HandleValidate(c, dto)
	if errRes != nil {
		return errRes.Response()
	}

	user, err := utils.GetAuthUser(c)
	if err != nil {
		return err
	}

	exists, err := h.service.Find(uint(pk), user.Id)
	if exists == nil || exists.Id == 0 {
		status := fiber.StatusNotFound
		message := fUtils.StatusMessage(status)
		errRes := api_error.NewErrorResponse(c, status, message)
		return errRes.Response()
	}

	if exists.UserId != user.Id || err != nil {
		status := fiber.StatusForbidden
		message := fUtils.StatusMessage(status)
		errRes := api_error.NewErrorResponse(c, status, message)
		return errRes.Response()
	}

	result, err := h.service.Update(uint(pk), user.Id, dto)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

func (h *HandlerStruct) Get(c *fiber.Ctx) error {
	params := c.AllParams()

	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	user, err := utils.GetAuthUser(c)
	if err != nil {
		return err
	}

	result, err := h.service.Find(uint(pk), user.Id)

	if result == nil || result.Id == 0 {
		status := fiber.StatusNotFound
		message := fUtils.StatusMessage(status)
		errRes := api_error.NewErrorResponse(c, status, message)
		return errRes.Response()
	}

	if result.UserId != user.Id || err != nil {
		status := fiber.StatusForbidden
		message := fUtils.StatusMessage(status)
		errRes := api_error.NewErrorResponse(c, status, message)
		return errRes.Response()
	}

	return c.JSON(result)
}

func (h *HandlerStruct) All(c *fiber.Ctx) error {
	user, err := utils.GetAuthUser(c)

	if err != nil {
		return err
	}

	results, err := h.service.GetByUserId(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": results,
	})
}

func (h *HandlerStruct) Delete(c *fiber.Ctx) error {
	params := c.AllParams()

	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	user, err := utils.GetAuthUser(c)
	if err != nil {
		return err
	}

	result, err := h.service.Delete(uint(pk), user.Id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"status": result,
	})
}
