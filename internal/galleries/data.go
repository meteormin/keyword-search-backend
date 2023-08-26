package galleries

import (
	"github.com/miniyus/go-restapi"
	"github.com/miniyus/keyword-search-backend/entity"
)

type GalleryRequest struct {
	restapi.RequestDTO[entity.Gallery]
	UserId      uint
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

func (gr *GalleryRequest) ToEntity(ent *entity.Gallery) error {
	ent.Subject = gr.Subject
	ent.Description = gr.Description

	return nil
}

type GalleryResponse struct {
	restapi.ResponseDTO[entity.Gallery]
}

func (gr *GalleryResponse) FromEntity(ent entity.Gallery) error {
	return nil
}
