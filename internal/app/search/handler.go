package search

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"github.com/miniyus/go-fiber/internal/utils"
	"strconv"
)

type Handler interface {
	Create(c *fiber.Ctx) error
	All(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	GetByHostId(c *fiber.Ctx) error
	BatchCreate(c *fiber.Ctx) error
}

type HandlerStruct struct {
	service Service
}

func NewHandler(s Service) Handler {
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

func (h *HandlerStruct) Create(c *fiber.Ctx) error {
	dto := &CreateSearch{}

	err := c.BodyParser(dto)
	if err != nil {
		res := api_error.NewValidationError(c)
		return res.Response()
	}

	errRes := utils.HandleValidate(c, dto)
	if errRes != nil {
		return errRes.Response()
	}

	created, err := h.service.Create(dto)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *HandlerStruct) All(c *fiber.Ctx) error {
	page, err := utils.GetPageFromCtx(c)
	if err != nil {
		return err
	}

	all, err := h.service.All(page)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(all)
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

	search, err := h.service.Find(uint(pk), user.Id)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(search)
}

func (h *HandlerStruct) BatchCreate(c *fiber.Ctx) error {
	params := c.AllParams()
	hostId, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	dto := &MultiCreateSearch{}
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
