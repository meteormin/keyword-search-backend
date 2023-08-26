package repo

import (
	"github.com/miniyus/go-restapi"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

type GalleryRepository struct {
	restapi.Repository[entity.Gallery]
}

func NewGalleryRepository(db *gorm.DB) *GalleryRepository {
	return &GalleryRepository{
		restapi.NewRepository[entity.Gallery](db, entity.Gallery{}),
	}
}
