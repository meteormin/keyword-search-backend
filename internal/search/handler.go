package search

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/miniyus/gofiber/apierrors"
	"github.com/miniyus/gofiber/config"
	"github.com/miniyus/gofiber/pagination"
	"github.com/miniyus/gofiber/utils"
	"github.com/miniyus/keyword-search-backend/internal/auth"
	"os"
	"path"
	"strconv"
)

type Handler interface {
	Create(c *fiber.Ctx) error
	All(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Patch(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	UploadImage(c *fiber.Ctx) error
	GetImage(c *fiber.Ctx) error
}

type HandlerStruct struct {
	service Service
}

func NewHandler(s Service) Handler {
	return &HandlerStruct{
		service: s,
	}
}

// Create
// @Summary create search
// @description create search
// @Tags Search
// @param request body CreateSearch true "create search"
// @Success 201 {object} Response
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Accept json
// @Produce json
// @Router /api/search [post]
// @Security BearerAuth
func (h *HandlerStruct) Create(c *fiber.Ctx) error {
	dto := &CreateSearch{}

	err := c.BodyParser(dto)
	if err != nil {
		return fiber.ErrBadRequest
	}

	errRes := utils.HandleValidate(c, dto)
	if errRes != nil {
		return errRes.Response()
	}

	user, err := auth.GetAuthUser(c)
	if err != nil {
		return err
	}

	created, err := h.service.Create(user.Id, dto)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// All
// @Summary create search
// @description create search
// @Tags Search
// @Param page query int true "page number"
// @Param page_size query int true "page size"
// @Success 200 {object} ResponseAll
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Failure 404 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/search/all [get]
// @Security BearerAuth
func (h *HandlerStruct) All(c *fiber.Ctx) error {
	page, err := pagination.GetPageFromCtx(c)
	if err != nil {
		return err
	}

	all, err := h.service.All(page)

	if err != nil {
		return err
	}

	return c.JSON(ResponseAll{
		Paginator: all,
		Data:      all.Data,
	})
}

// Get
// @Summary get search
// @description get search
// @Tags Search
// @Param id path int true "search pk"
// @Success 200 {object} Response
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Failure 404 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/search/{id} [get]
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

	search, err := h.service.Find(uint(pk), user.Id)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(search)
}

// Update
// @Summary get search
// @description get search
// @Tags Search
// @Param id path int true "search pk"
// @param request body UpdateSearch true "update search"
// @Success 200 {object} Response
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Failure 404 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/search/{id} [put]
// @Security BearerAuth
func (h *HandlerStruct) Update(c *fiber.Ctx) error {
	params := c.AllParams()
	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	user, err := auth.GetAuthUser(c)
	if err != nil {
		return err
	}

	dto := UpdateSearch{}

	errRes := utils.HandleValidate(c, &dto)
	if errRes != nil {
		return errRes.Response()
	}

	exists, err := h.service.Find(uint(pk), user.Id)
	if exists == nil || exists.Id == 0 {
		return fiber.ErrNotFound
	}

	updated, err := h.service.Update(uint(pk), user.Id, &dto)
	if err != nil {
		return err
	}

	return c.JSON(updated)
}

// Patch
// @Summary get search
// @description get search
// @Tags Search
// @Param id path int true "search pk"
// @param request body PatchSearch true "update search"
// @Success 200 {object} Response
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Failure 404 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/search/{id} [patch]
// @Security BearerAuth
func (h *HandlerStruct) Patch(c *fiber.Ctx) error {
	params := c.AllParams()
	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}
	user, err := auth.GetAuthUser(c)
	if err != nil {
		return err
	}

	dto := PatchSearch{}

	errRes := utils.HandleValidate(c, &dto)
	if errRes != nil {
		return errRes.Response()
	}

	_, err = h.service.Find(uint(pk), user.Id)
	if err != nil {
		return err
	}

	patch, err := h.service.Patch(uint(pk), user.Id, &dto)
	if err != nil {
		return err
	}

	return c.JSON(patch)
}

// Delete
// @Summary get search
// @description get search
// @Tags Search
// @Param id path int true "search pk"
// @Success 204 {object} utils.StatusResponse
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Failure 404 {object} apierrors.ErrorResponse
// @Accept json
// @Produce json
// @Router /api/search/{id} [delete]
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

	rs, err := h.service.Delete(uint(pk), user.Id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).JSON(utils.StatusResponse{
		Status: rs,
	})
}

// UploadImage
// @Summary upload image
// @Description upload image
// @Tags Search
// @Param id path int true "search pk"
// @Success 201 {object} Response
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Failure 403 {object} apierrors.ErrorResponse
// @Failure 404 {object} apierrors.ErrorResponse
// @Accept multipart/form-data
// @produce json
// @Router /api/search/{id} [post]
// @Security BearerAuth
func (h *HandlerStruct) UploadImage(c *fiber.Ctx) error {
	params := c.AllParams()
	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	user, err := auth.GetAuthUser(c)
	if err != nil {
		return err
	}

	file := form.File["image"][0]
	res, err := h.service.UploadImage(uint(pk), user.Id, file)
	if err != nil {
		return err
	}

	dataPath := config.GetConfigs().Path.DataPath
	savePath := fmt.Sprintf("images/%s", file.Filename)
	if _, err = os.Stat(path.Join(dataPath, "images")); err != nil {
		err = os.Mkdir(path.Join(dataPath, "images"), 0775)
		if err != nil {
			return err
		}
	}
	err = c.SaveFile(file, path.Join(dataPath, savePath))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *HandlerStruct) GetImage(c *fiber.Ctx) error {
	params := c.AllParams()
	pk, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		return err
	}

	user, err := auth.GetAuthUser(c)
	if err != nil {
		return err
	}

	imagePath, err := h.service.FindImagePath(uint(pk), user.Id)
	if err != nil {
		return err
	}

	dataPath := config.GetConfigs().Path.DataPath

	return c.Download(path.Join(dataPath, imagePath), path.Base(imagePath))
}
