package galleries

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-restapi"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/internal/auth"
)

type Handler struct {
	restapi.Handler[entity.Gallery, *GalleryRequest, *GalleryResponse]
}

func NewHandler(service restapi.Service[entity.Gallery, *GalleryRequest, *GalleryResponse]) *Handler {
	h := &Handler{
		Handler: restapi.NewHandler[entity.Gallery, *GalleryRequest, *GalleryResponse](
			&GalleryRequest{},
			service,
		),
	}
	return h
}

func (h *Handler) Create(ctx *fiber.Ctx) error {
	hook := h.Handler.Hook()
	hook.Create().BeforeCallService(func(dto *GalleryRequest) error {
		user, err := auth.GetAuthUser(ctx)
		if err != nil {
			return err
		}
		dto.UserId = user.Id

		return nil
	})

	return h.Handler.Create(ctx)
}
