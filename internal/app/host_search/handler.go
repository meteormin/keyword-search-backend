package host_search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/app/search"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"github.com/miniyus/go-fiber/internal/utils"
	"strconv"
)

type Handler interface {
	GetByHostId(c *fiber.Ctx) error
	BatchCreate(c *fiber.Ctx) error
}

type HandlerStruct struct {
	service search.Service
}

func NewHandler(s search.Service) Handler {
	return &HandlerStruct{service: s}
}

func (h *HandlerStruct) GetByHostId(c *fiber.Ctx) error {
	page, err := utils.GetPageFromCtx(c)
	if err != nil {
		return err
	}

	params := c.AllParams()
	hostId, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	data, err := h.service.GetByHostId(uint(hostId), page)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(data)
}

func (h *HandlerStruct) BatchCreate(c *fiber.Ctx) error {
	params := c.AllParams()
	hostId, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	dto := &search.MultiCreateSearch{}
	err = c.BodyParser(dto)
	if err != nil {
		res := api_error.NewValidationError(c)
		return res.Response()
	}

	errRes := utils.HandleValidate(c, dto)
	if errRes != nil {
		return errRes.Response()
	}

	create, err := h.service.BatchCreate(uint(hostId), dto.Search)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(create)
}
