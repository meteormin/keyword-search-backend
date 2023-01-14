package host_search

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/app/search"
	_ "github.com/miniyus/keyword-search-backend/internal/core/api_error"
	"github.com/miniyus/keyword-search-backend/internal/core/logger"
	"github.com/miniyus/keyword-search-backend/internal/core/register/resolver"
	"github.com/miniyus/keyword-search-backend/internal/utils"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"strconv"
)

type Handler interface {
	GetByHostId(c *fiber.Ctx) error
	GetDescriptionsByHostId(c *fiber.Ctx) error
	BatchCreate(c *fiber.Ctx) error
	logger.HasLogger
}

type HandlerStruct struct {
	service search.Service
	logger.HasLoggerStruct
}

func NewHandler(s search.Service) Handler {
	return &HandlerStruct{
		service: s,
		HasLoggerStruct: logger.HasLoggerStruct{
			Logger: s.GetLogger(),
		},
	}
}

// GetByHostId
// @Summary get by host id
// @description get by host id
// @Tags Hosts
// @Param id path int true "host id"
// @Param page query int true "page number"
// @Param page_size query int true "page size"
// @Success 200 {object} search.ResponseByHost
// @Failure 400 {object} api_error.ValidationErrorResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts/{id}/search [get]
// @Security BearerAuth
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

	return c.Status(fiber.StatusOK).JSON(search.ResponseByHost{
		Paginator: utils.Paginator[search.Response]{
			Page:       data.Page,
			TotalCount: data.TotalCount,
		},
		Data: data.Data,
	})
}

// GetDescriptionsByHostId
// @Summary get by host id
// @description get by host id
// @Tags Hosts
// @Param id path int true "host id"
// @Param page query int true "page number"
// @Param page_size query int true "page size"
// @Success 200 {object} search.DescriptionWithPage
// @Failure 400 {object} api_error.ValidationErrorResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts/{id}/search/descriptions [get]
// @Security BearerAuth
func (h *HandlerStruct) GetDescriptionsByHostId(c *fiber.Ctx) error {
	page, err := utils.GetPageFromCtx(c)
	if err != nil {
		return err
	}

	params := c.AllParams()
	hostId, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	data, err := h.service.GetDescriptionsByHostId(uint(hostId), page)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(search.DescriptionWithPage{
		Paginator: utils.Paginator[search.Description]{
			Page:       data.Page,
			TotalCount: data.TotalCount,
		},
		Data: data.Data,
	})
}

// BatchCreate
// @Summary batch create search by host id
// @description batch create search by host id
// @Tags Hosts
// @Param id path int true "host id"
// @Param request body search.MultiCreateSearch true "multi create search"
// @Success 200 {object} search.Response
// @Failure 400 {object} api_error.ValidationErrorResponse
// @Failure 403 {object} api_error.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts/{id}/search [post]
// @Security BearerAuth
func (h *HandlerStruct) BatchCreate(c *fiber.Ctx) error {
	params := c.AllParams()
	hostId, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	dto := &search.MultiCreateSearch{}
	err = c.BodyParser(dto)
	if err != nil {
		return fiber.ErrBadRequest
	}

	errRes := utils.HandleValidate(c, dto)
	if errRes != nil {
		return errRes.Response()
	}

	var jobDispatcher worker.Dispatcher
	err = resolver.ResolveContext(c, &jobDispatcher)
	if err != nil {
		return err
	}

	jobDispatcher.SelectWorker(string(config.DefaultWorker))

	searchCollection := utils.NewCollection(dto.Search)
	searchCollection.Chunk(100, func(v []*search.CreateSearch, i int) {
		jobId := fmt.Sprintf("hosts.%s", strconv.Itoa(int(hostId)))
		err = jobDispatcher.Dispatch(jobId, func(job worker.Job) error {
			create, err := h.service.BatchCreate(uint(hostId), dto.Search)
			if err != nil {
				return err
			}

			if create != nil {
				return nil
			}

			return nil
		})
	})

	if err != nil {
		return err
	}

	//create, err := h.service.BatchCreate(uint(hostId), dto.Search)
	//if err != nil {
	//	return err
	//}

	return c.Status(fiber.StatusCreated).JSON(utils.StatusResponse{
		Status: true,
	})
}
