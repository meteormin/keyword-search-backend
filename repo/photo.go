package repo

import (
	"github.com/miniyus/go-restapi"
	"github.com/miniyus/keyword-search-backend/entity"
	"gorm.io/gorm"
)

type PhotoRepository struct {
	restapi.Repository[entity.Photo]
}

func NewPhotoRepository(db *gorm.DB) *PhotoRepository {
	return &PhotoRepository{
		restapi.NewRepository[entity.Photo](db, entity.Photo{}),
	}
}
