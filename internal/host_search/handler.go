package host_search

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/miniyus/gofiber/apierrors"
	"github.com/miniyus/gofiber/jobqueue"
	"github.com/miniyus/gofiber/pagination"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/gollection"
	worker "github.com/miniyus/goworker"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/internal/auth"
	"github.com/miniyus/keyword-search-backend/internal/search"
	"strconv"
)

type Handler interface {
	GetByHostId(c *fiber.Ctx) error
	GetDescriptionsByHostId(c *fiber.Ctx) error
	BatchCreate(c *fiber.Ctx) error
}

type HandlerStruct struct {
	service    search.Service
	dispatcher worker.Dispatcher
}

func NewHandler(s search.Service, dispatcher worker.Dispatcher) Handler {
	return &HandlerStruct{
		service:    s,
		dispatcher: dispatcher,
	}
}

// GetByHostId
// @Summary get by host id
// @description get by host id
// @Tags Hosts
// @Param id path int true "host id"
// @Param filter query search.Query false "filter query"
// @Success 200 {object} search.ResponseByHost
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts/{id}/search [get]
// @Security BearerAuth
func (h *HandlerStruct) GetByHostId(c *fiber.Ctx) error {
	params := c.AllParams()
	hostId, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	user, err := auth.GetAuthUser(c)
	if err != nil {
		return err
	}

	var queryFilter search.Query
	err = c.QueryParser(&queryFilter)
	if err != nil {
		return err
	}

	page, err := pagination.GetPageFromCtx(c)
	queryFilter.Page = page
	data, err := h.service.GetByHostId(uint(hostId), user.Id, queryFilter)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(search.ResponseByHost{
		Paginator: pagination.Paginator[search.Response]{
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
// @Param filter query search.Query false "filter query"
// @Success 200 {object} search.DescriptionWithPage
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/hosts/{id}/search/descriptions [get]
// @Security BearerAuth
func (h *HandlerStruct) GetDescriptionsByHostId(c *fiber.Ctx) error {
	params := c.AllParams()
	hostId, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	user, err := auth.GetAuthUser(c)
	if err != nil {
		return err
	}

	var queryFilter search.Query
	err = c.QueryParser(&queryFilter)
	if err != nil {
		return err
	}

	data, err := h.service.GetDescriptionsByHostId(uint(hostId), user.Id, queryFilter)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(search.DescriptionWithPage{
		Paginator: pagination.Paginator[search.Description]{
			Page:       data.Page,
			TotalCount: data.TotalCount,
		},
		Data: data.Data,
	})
}

// BatchCreate
// @Summary batch create search by host id
// @Description batch create search by host id
// @Tags Hosts
// @Param id path int true "host id"
// @Param request body search.MultiCreateSearch true "multi create search"
// @Success 200 {object} search.Response
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
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

	dto := search.MultiCreateSearch{}

	errRes := utils.HandleValidate(c, &dto)
	if errRes != nil {
		return errRes.Response()
	}

	if err != nil {
		return err
	}

	user, err := auth.GetAuthUser(c)
	if err != nil {
		return err
	}

	h.dispatcher.SelectWorker(string(config.DefaultWorker))

	jobId := fmt.Sprintf("hosts.%s", strconv.Itoa(int(hostId)))

	searchCollection := gollection.NewCollection(dto.Search)
	searchCollection.Chunk(100, func(v []*search.CreateSearch, i int) {

		err = h.dispatcher.Dispatch(jobId, func(job *worker.Job) error {
			job.Meta["user_id"] = user.Id

			create, batchCreateErr := h.service.BatchCreate(uint(hostId), user.Id, dto.Search)
			if batchCreateErr != nil {
				return batchCreateErr
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

	foundJob, err := jobqueue.FindJob(jobId, string(config.DefaultWorker))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(utils.DataResponse[*worker.Job]{
		Data: foundJob,
	})
}
