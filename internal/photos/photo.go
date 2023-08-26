package photos

import (
	"github.com/gofiber/fiber/v2"
	"github.com/miniyus/go-restapi"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/repo"
	"gorm.io/gorm"
)

const Prefix = "/galleries/:gallery_id/photo"

func New(db *gorm.DB) app.SubRouter {
	s := NewService(
		repo.NewPhotoRepository(db),
		repo.NewGalleryRepository(db),
	)
	h := NewHandler(s)
	defaultRouter := restapi.Route[entity.Photo, *Request, *Response](h)
	return func(router fiber.Router) {
		defaultRouter(router)
		router.Post("/", h.Create)
		router.Get("/:photo_id", h.FindImage)
	}
}
