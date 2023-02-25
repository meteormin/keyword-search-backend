package hosts

import (
	"github.com/miniyus/gofiber/pagination"
	"github.com/miniyus/keyword-search-backend/entity"
)

type CreateHost struct {
	UserId      uint   `json:"user_id"`
	Host        string `json:"host" validate:"required,url"`
	Subject     string `json:"subject" validate:"required"`
	Description string `json:"description" validate:"required"`
	Path        string `json:"path" validate:"required"`
	Publish     bool   `json:"publish" validate:"required,boolean"`
}

func (ch CreateHost) ToEntity() entity.Host {
	return entity.Host{
		UserId:      ch.UserId,
		Host:        ch.Host,
		Subject:     ch.Subject,
		Description: ch.Description,
		Path:        ch.Path,
		Publish:     ch.Publish,
	}
}

type UpdateHost struct {
	Host        string `json:"host" validate:"required"`
	Subject     string `json:"subject" validate:"required"`
	Description string `json:"description" validate:"required"`
	Path        string `json:"path" validate:"required"`
	Publish     bool   `json:"publish" validate:"required,boolean"`
}

func (uh UpdateHost) ToEntity() entity.Host {
	return entity.Host{
		Subject:     uh.Subject,
		Description: uh.Description,
		Path:        uh.Path,
		Publish:     uh.Publish,
	}
}

type PatchHost struct {
	Host        *string `json:"host,omitempty" validate:"omitempty,url"`
	Subject     *string `json:"subject,omitempty" validate:"omitempty"`
	Description *string `json:"description,omitempty" validate:"omitempty"`
	Path        *string `json:"path,omitempty" validate:"omitempty"`
	Publish     *bool   `json:"publish,omitempty" validate:"omitempty,boolean"`
}

func (ph PatchHost) ToEntity() entity.Host {
	var ent entity.Host

	if ph.Host != nil {
		ent.Host = *ph.Host
	}

	if ph.Subject != nil {
		ent.Subject = *ph.Subject
	}

	if ph.Description != nil {
		ent.Description = *ph.Description
	}

	if ph.Path != nil {
		ent.Path = *ph.Path
	}

	if ph.Publish != nil {
		ent.Publish = *ph.Publish
	}

	return ent
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
	pagination.Paginator[HostResponse]
	Data []HostResponse `json:"data"`
}

type HostResponseAll struct {
	pagination.Paginator[entity.Host]
	Data []entity.Host `json:"data"`
}

type HostSubjectsResponse struct {
	pagination.Paginator[Subjects]
	Data []Subjects `json:"data"`
}

func (hr HostResponse) FromEntity(host entity.Host) HostResponse {
	return HostResponse{
		Id:          host.ID,
		Host:        host.Host,
		Path:        host.Path,
		UserId:      host.UserId,
		Subject:     host.Subject,
		Description: host.Description,
		Publish:     host.Publish,
	}
}
