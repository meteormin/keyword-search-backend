package galleries

import (
	"github.com/miniyus/go-restapi"
	"github.com/miniyus/gofiber/app"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/repo"
	"gorm.io/gorm"
)

const Prefix = "/galleries"

func New(db *gorm.DB) app.SubRouter {
	r := repo.NewGalleryRepository(db)
	service := NewService(r)
	h := NewHandler(service)

	return restapi.Route[entity.Gallery, *GalleryRequest, *GalleryResponse](h)
}
