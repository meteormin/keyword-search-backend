package hosts

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/miniyus/gofiber/apierrors"
	"github.com/miniyus/gofiber/pagination"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/internal/auth"
	"strconv"
)

type Handler interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Patch(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	GetSubjects(c *fiber.Ctx) error
	All(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type HandlerStruct struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &HandlerStruct{service: service}
}

// Create
// @Summary create host
// @Description create host
// @Tags Hosts
// @Param request body CreateHost true "create host"
// @Success 201 {object} HostResponse
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts [post]
// @Security BearerAuth
func (h *HandlerStruct) Create(c *fiber.Ctx) error {
	user, err := auth.GetAuthUser(c)
	if err != nil {
		return err
	}

	dto := CreateHost{}

	errRes := utils.HandleValidate(c, &dto)
	if errRes != nil {
		return errRes.Response()
	}

	dto.UserId = user.Id

	result, err := h.service.Create(&dto)
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
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Failure 404 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts/{id} [put]
// @Security BearerAuth
func (h *HandlerStruct) Update(c *fiber.Ctx) error {
	dto := UpdateHost{}
	params := c.AllParams()

	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	errRes := utils.HandleValidate(c, &dto)
	if errRes != nil {
		return errRes.Response()
	}

	user, err := auth.GetAuthUser(c)
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

	result, err := h.service.Update(uint(pk), user.Id, &dto)
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
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Failure 404 {object} apierrors.ErrorResponse
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
	dto := PatchHost{}

	errRes := utils.HandleValidate(c, &dto)
	if errRes != nil {
		return errRes.Response()
	}

	user, err := auth.GetAuthUser(c)
	if err != nil {
		return err
	}

	result, err := h.service.Patch(uint(pk), user.Id, &dto)
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
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Failure 404 {object} apierrors.ErrorResponse
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

	user, err := auth.GetAuthUser(c)
	if err != nil {
		return err
	}

	result, err := h.service.Find(uint(pk), user.Id)

	if err != nil {
		return err
	}

	return c.JSON(result)
}

// GetSubjects
// @Summary get subjects by group id
// @Description get subjects by group id
// @Tags Hosts
// @Param page query int true "page number"
// @Param page_size query int true "page size"
// @Success 200 {object} HostSubjectsResponse
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Failure 404 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts/subjects [get]
// @Security BearerAuth
func (h *HandlerStruct) GetSubjects(c *fiber.Ctx) error {
	user, err := auth.GetAuthUser(c)
	if err != nil {
		return err
	}

	if user.GroupId == nil {
		return fiber.ErrForbidden
	}

	page, err := pagination.GetPageFromCtx(c)
	if err != nil {
		return err
	}

	result, err := h.service.GetSubjectsByGroupId(*user.GroupId, page)

	if err != nil {
		return err
	}

	return c.JSON(HostSubjectsResponse{
		Paginator: result,
		Data:      result.Data,
	})
}

// All
// @Summary get all hosts
// @description get all hosts
// @Tags Hosts
// @Param page query int true "page number"
// @Param page_size query int true "page size"
// @Success 200 {object} HostListResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts [get]
// @Security BearerAuth
func (h *HandlerStruct) All(c *fiber.Ctx) error {
	page, err := pagination.GetPageFromCtx(c)
	if err != nil {
		return err
	}

	user, err := auth.GetAuthUser(c)

	if err != nil {
		return err
	}

	results, err := h.service.GetByUserId(user.Id, page)
	if err != nil {
		return err
	}

	return c.JSON(HostListResponse{
		Paginator: results,
		Data:      results.Data,
	})
}

// Delete
// @Summary Delete host
// @description Delete host
// @Tags Hosts
// @param id path int true "host pk"
// @Success 204 {object} utils.StatusResponse
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Failure 404 {object} apierrors.ErrorResponse
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

	user, err := auth.GetAuthUser(c)
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
