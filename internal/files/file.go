package files

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-restapi"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/repo"
	"gorm.io/gorm"
	"strconv"
)

type Response struct {
	restapi.ResponseDTO[entity.File] `json:"-"`
	Id                               uint   `json:"id"`
	MimeType                         string `json:"mime_type"`
	Size                             int64  `json:"size"`
	Path                             string `json:"path"`
	Extension                        string `json:"extension"`
}

func (res *Response) FromEntity(ent entity.File) error {
	err := restapi.Map(res, ent)
	if err != nil {
		return err
	}
	res.Id = ent.ID
	return nil
}

type ResponseList struct {
	Files []Response `json:"files"`
}

type Handler struct {
	repository *repo.FileRepository
}

func (h *Handler) GetFiles(ctx *fiber.Ctx) error {
	all, err := h.repository.All()
	if err != nil {
		return err
	}

	resSlice := make([]Response, 0)
	for _, file := range all {
		res := Response{}
		err = res.FromEntity(file)
		if err != nil {
			return err
		}
		resSlice = append(resSlice, res)
	}

	return ctx.Status(fiber.StatusOK).JSON(ResponseList{
		resSlice,
	})
}

func (h *Handler) GetFile(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	pk, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return err
	}

	find, err := h.repository.Find(uint(pk))
	if err != nil {
		return err
	}

	res := Response{}
	err = res.FromEntity(*find)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

const Prefix = "/files"

func New(db *gorm.DB) app.SubRouter {
	r := repo.NewFileRepository(db)
	h := &Handler{r}

	return func(router fiber.Router) {
		router.Get("/", h.GetFiles)
		router.Get("/:id", h.GetFile)
	}
}
