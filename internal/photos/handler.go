package photos

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-restapi"
	"github.com/miniyus/keyword-search-backend/config"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/internal/auth"
	"github.com/miniyus/keyword-search-backend/utils"
	"os"
	"path"
	"strconv"
)

type Handler struct {
	restapi.Handler[entity.Photo, *Request, *Response]
	Service *Service
}

func NewHandler(service restapi.Service[entity.Photo, *Request, *Response]) *Handler {
	return &Handler{
		Handler: restapi.NewHandler[entity.Photo, *Request, *Response](&Request{}, service),
		Service: service.(*Service),
	}
}

type BatchCreateRequest struct {
	Photos []Request `json:"photos"`
}

// Create
// @Summary create photo
// @Description create photo
// @Tags Hosts
// @Param request body BatchCreateRequest true "create photos"
// @Success 201 {object} Response
// @Failure 400 {object} apierrors.ValidationErrorResponse
// @Accept json
// @Produce json
// @Router /api/gallery{gallery_id}/photo [post]
// @Security BearerAuth
func (h *Handler) Create(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	pk, err := strconv.ParseUint(params["gallery_id"], 10, 64)
	if err != nil {
		return err
	}

	user, err := auth.GetAuthUser(ctx)
	if err != nil {
		return err
	}

	dtos := make([]Request, 0)
	err = ctx.BodyParser(&dtos)
	if err != nil {
		return err
	}

	files := form.File["image"]

	for i, file := range files {
		dto := dtos[i]
		dataPath := config.GetConfigs().Path.DataPath
		savePath := fmt.Sprintf("images/%s", file.Filename)
		if _, err = os.Stat(path.Join(dataPath, "images")); err != nil {
			err = os.Mkdir(path.Join(dataPath, "images"), 0775)
			if err != nil {
				return err
			}
		}
		info, err := utils.GetFileInfo(file)
		if err != nil {
			return err
		}
		dto.GalleryId = uint(pk)
		dto.FileInfo = FileInfo{
			Path:      savePath,
			Extension: info.Extension,
			MimType:   info.MimeType,
			Size:      info.Size,
		}

		dtos[i] = dto
	}

	create, err := h.Service.BatchCreate(dtos, user.Id)
	if err != nil {
		return err
	}

	for i, photo := range create {
		file := files[i]
		err = ctx.SaveFile(file, photo.FileInfo.Path)
		if err != nil {
			return err
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(create)
}

func (h *Handler) FindImage(ctx *fiber.Ctx) error {
	params := ctx.AllParams()

	pk, err := strconv.ParseUint(params["photo_id"], 10, 64)
	if err != nil {
		return err
	}

	user, err := auth.GetAuthUser(ctx)
	if err != nil {
		return err
	}

	image, err := h.Service.FindImage(uint(pk), user.Id)
	if err != nil {
		return err
	}

	imagePath := image.FileInfo.Path

	dataPath := config.GetConfigs().Path.DataPath

	return ctx.Download(
		path.Join(dataPath, imagePath),
		path.Base(imagePath),
	)
}
