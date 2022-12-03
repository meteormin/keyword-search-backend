package hosts

import (
	"github.com/gofiber/fiber/v2"
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

func NewHandler(service Service) *HandlerStruct {
	return &HandlerStruct{service: service}
}

func (h *HandlerStruct) Create(c *fiber.Ctx) error {
	dto := &CreateHost{}

	err := c.BodyParser(dto)
	if err != nil {
		errRes := api_error.NewValidationError(c)
		return errRes.Response()
	}

	err = utils.HandleValidate(c, dto)
	if err != nil {
		return err
	}

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

	err = utils.HandleValidate(c, dto)
	if err != nil {
		return err
	}

	result, err := h.service.Update(uint(pk), dto)

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

	result, err := h.service.Find(uint(pk))
	if err != nil {
		return err
	}

	return c.JSON(result)
}

func (h *HandlerStruct) All(c *fiber.Ctx) error {
	results, err := h.service.All()
	if err != nil {
		return err
	}

	return c.JSON(results)
}

func (h *HandlerStruct) Delete(c *fiber.Ctx) error {
	params := c.AllParams()

	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	result, err := h.service.Delete(uint(pk))

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"status": result,
	})
}
