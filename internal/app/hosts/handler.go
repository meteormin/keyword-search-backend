package hosts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-fiber/internal/core/api_error"
	"github.com/miniyus/go-fiber/internal/core/logger"
	"github.com/miniyus/go-fiber/internal/utils"
	"strconv"
)

type Handler interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Patch(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	All(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	logger.HasLogger
}

type HandlerStruct struct {
	service Service
	logger.HasLoggerStruct
}

func NewHandler(service Service) Handler {
	return &HandlerStruct{service: service, HasLoggerStruct: logger.HasLoggerStruct{Logger: service.GetLogger()}}
}

// Create
// @Summary create host
// @description create host
// @Tags Hosts
// @Param request body CreateHost true "create host"
// @Success 201 {object} HostResponse
// @Failure 400 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts [post]
// @Security BearerAuth
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

// Update
// @Summary update host
// @description update host
// @Tags Hosts
// @Param id path int true "host pk"
// @Param request body UpdateHost true "update host"
// @Success 200 {object} HostResponse
// @Failure 400 {object} api_error.ErrorResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Failure 404 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts/{id} [put]
// @Security BearerAuth
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
		return fiber.ErrNotFound
	}

	if exists.UserId != user.Id || err != nil {
		return fiber.ErrForbidden
	}

	result, err := h.service.Update(uint(pk), user.Id, dto)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// Patch
// @Summary patch host
// @description patch host
// @Tags Hosts
// @Param id path int true "host pk"
// @Param request body PatchHost true "patch host"
// @Success 200 {object} HostResponse
// @Failure 400 {object} api_error.ErrorResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Failure 404 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts/{id} [patch]
// @Security BearerAuth
func (h *HandlerStruct) Patch(c *fiber.Ctx) error {
	params := c.AllParams()
	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}
	dto := &PatchHost{}

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

	result, err := h.service.Patch(uint(pk), user.Id, dto)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// Get
// @Summary patch host
// @description patch host
// @Tags Hosts
// @Param id path int true "host pk"
// @Success 200 {object} HostResponse
// @Failure 400 {object} api_error.ErrorResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Failure 404 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts/{id} [get]
// @Security BearerAuth
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

	if err != nil {
		return err
	}

	return c.JSON(result)
}

// All
// @Summary get all hosts
// @description get all hosts
// @Tags Hosts
// @Param id path int true "host pk"
// @Success 200 {object} utils.DataResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts [get]
// @Security BearerAuth
func (h *HandlerStruct) All(c *fiber.Ctx) error {
	user, err := utils.GetAuthUser(c)

	if err != nil {
		return err
	}

	results, err := h.service.GetByUserId(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(utils.DataResponse{
		Data: results,
	})
}

// Delete
// @Summary Delete host
// @description Delete host
// @Tags Hosts
// @param id path int true "host pk"
// @Success 204 {object} utils.StatusResponse
// @Failure 400 {object} api_error.ErrorResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Failure 404 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts/{id} [delete]
// @Security BearerAuth
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

	return c.Status(fiber.StatusNoContent).JSON(utils.StatusResponse{
		Status: result,
	})
}
