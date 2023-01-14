package hosts

import (
	"github.com/miniyus/keyword-search-backend/internal/entity"
	"github.com/miniyus/keyword-search-backend/internal/utils"
)

type CreateHost struct {
	UserId      uint   `json:"user_id"`
	Host        string `json:"host" validate:"required,url"`
	Subject     string `json:"subject" validate:"required"`
	Description string `json:"description" validate:"required"`
	Path        string `json:"path" validate:"required"`
	Publish     bool   `json:"publish" validate:"required,boolean"`
}

type UpdateHost struct {
	Host        string `json:"host" validate:"required"`
	Subject     string `json:"subject" validate:"required"`
	Description string `json:"description" validate:"required"`
	Path        string `json:"path" validate:"required,dir"`
	Publish     bool   `json:"publish" validate:"required,boolean"`
}

type PatchHost struct {
	Host        *string `json:"host,omitempty" validate:"omitempty,url"`
	Subject     *string `json:"subject,omitempty" validate:"omitempty"`
	Description *string `json:"description,omitempty" validate:"omitempty"`
	Path        *string `json:"path,omitempty" validate:"omitempty,dir"`
	Publish     *bool   `json:"publish,omitempty" validate:"omitempty,boolean"`
}

type HostResponse struct {
	Id          uint   `json:"id"`
	UserId      uint   `json:"user_id"`
	Host        string `json:"host"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Publish     bool   `json:"publish"`
}

type Subjects struct {
	Id      uint   `json:"id"`
	Subject string `json:"subject"`
}

type HostListResponse struct {
	utils.Paginator[HostResponse]
	Data []HostResponse `json:"data"`
}

type HostResponseAll struct {
	utils.Paginator[entity.Host]
	Data []entity.Host `json:"data"`
}

type HostSubjectsResponse struct {
	utils.Paginator[Subjects]
	Data []Subjects `json:"data"`
}

func ToHostResponse(host *entity.Host) *HostResponse {
	return &HostResponse{
		Id:          host.ID,
		Host:        host.Host,
		Path:        host.Path,
		UserId:      host.UserId,
		Subject:     host.Subject,
		Description: host.Description,
		Publish:     host.Publish,
	}
}
