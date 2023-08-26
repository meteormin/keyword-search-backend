package galleries

import (
	"github.com/miniyus/go-restapi"
	"github.com/miniyus/keyword-search-backend/entity"
)

type Service struct {
	restapi.Service[entity.Gallery, *GalleryRequest, *GalleryResponse]
}

func NewService(repository restapi.Repository[entity.Gallery]) *Service {
	return &Service{
		restapi.NewService[entity.Gallery, *GalleryRequest, *GalleryResponse](repository, &GalleryResponse{}),
	}
}
